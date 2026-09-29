package main

import (
	"bytes"
	"debug/buildinfo"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mtk177a/agent-harness-retro/skills/outer-harness-retrospective/scripts/retro-pilot/internal/pilot"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "pilot failed:", err)
		os.Exit(1)
	}
}

func run() error {
	var binary, configPath string
	flag.StringVar(&binary, "agent-sessions-bin", "", "absolute path to the pinned agent-sessions executable")
	flag.StringVar(&configPath, "config", "", "optional agent-sessions configuration path")
	flag.Parse()
	if binary == "" || !filepath.IsAbs(binary) || flag.NArg() != 0 {
		return errors.New("an absolute --agent-sessions-bin is required")
	}
	if configPath != "" && !filepath.IsAbs(configPath) {
		return errors.New("--config must be an absolute path")
	}
	if err := checkBuild(binary); err != nil {
		return err
	}
	caller := func(args ...string) ([]byte, error) {
		if configPath != "" {
			args = append([]string{args[0], "--config", configPath}, args[1:]...)
		}
		command := exec.Command(binary, args...)
		var output bytes.Buffer
		command.Stdout = &output
		command.Stderr = io.Discard
		err := command.Run()
		if output.Len() > 8<<20 || output.Len() == 0 {
			return nil, errors.New("invalid agent-sessions output")
		}
		var exit *exec.ExitError
		if err != nil && !errors.As(err, &exit) {
			return nil, errors.New("agent-sessions execution failed")
		}
		return output.Bytes(), nil
	}
	bundle, err := pilot.Run(caller, time.Now().UTC())
	if err != nil {
		return errors.New("pilot processing failed")
	}
	directory, err := os.MkdirTemp("", "agent-harness-retro-pilot-")
	if err != nil {
		return errors.New("could not create private output directory")
	}
	working, err := os.Getwd()
	if err != nil {
		os.RemoveAll(directory)
		return errors.New("could not verify private output location")
	}
	outside, err := outsideProtectedRoot(directory, working)
	if err != nil || !outside {
		os.RemoveAll(directory)
		return errors.New("private output must be outside the repository")
	}
	path := filepath.Join(directory, "evidence.json")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		os.Remove(directory)
		return errors.New("could not create private output")
	}
	encodeErr := json.NewEncoder(file).Encode(bundle)
	closeErr := file.Close()
	if encodeErr != nil || closeErr != nil {
		os.RemoveAll(directory)
		return errors.New("could not write private output")
	}
	public := struct {
		Status     string                         `json:"status"`
		Reason     string                         `json:"reason,omitempty"`
		Counts     map[string]pilot.ProviderCount `json:"counts,omitempty"`
		EventCount int                            `json:"event_count"`
	}{bundle.Status, bundle.Reason, bundle.Counts, bundle.EventCount}
	if err := json.NewEncoder(os.Stdout).Encode(public); err != nil {
		os.RemoveAll(directory)
		return errors.New("could not report pilot result")
	}
	fmt.Fprintln(os.Stderr, "private metadata file:", path)
	return nil
}

func outsideProtectedRoot(directory, working string) (bool, error) {
	protected := filepath.Clean(filepath.Join(working, "..", ".."))
	for path := working; ; path = filepath.Dir(path) {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			protected = path
			break
		}
		if filepath.Dir(path) == path {
			break
		}
	}
	protected, err := filepath.EvalSymlinks(protected)
	if err != nil {
		return false, err
	}
	resolvedDirectory, err := filepath.EvalSymlinks(directory)
	if err != nil {
		return false, err
	}
	relative, err := filepath.Rel(protected, resolvedDirectory)
	if err != nil {
		return false, err
	}
	return relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)), nil
}

func checkBuild(binary string) error {
	info, err := buildinfo.ReadFile(binary)
	if err != nil {
		return errors.New("agent-sessions build information is unavailable")
	}
	if info.Path != "github.com/mtk177a/agent-sessions/cmd/agent-sessions" {
		return errors.New("unexpected agent-sessions executable")
	}
	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}
	if revision != pilot.AgentSessionsCommit || modified != "false" {
		return errors.New("agent-sessions executable does not match the pinned clean commit")
	}
	return nil
}
