#!/bin/sh
set -eu

revision=8bc9542a15ed85f0b9b341886e920c76b1572d20

if [ "$#" -ne 1 ]; then
    printf '%s\n' 'Usage: sh install-agent-sessions.sh <absolute-bin-directory>' >&2
    exit 2
fi
case "$1" in
    /*) bin_dir=$1 ;;
    *) printf '%s\n' 'The bin directory must be an absolute path.' >&2; exit 2 ;;
esac
target=$bin_dir/agent-sessions
if [ -e "$target" ] || [ -L "$target" ]; then
    printf 'Installation target already exists: %s\n' "$target" >&2
    exit 1
fi

command -v git >/dev/null
GOTOOLCHAIN=local go version >&2
build=$(mktemp -d "${TMPDIR:-/tmp}/retro-agent-sessions.XXXXXXXX")
staged=
cleanup() {
    if [ -n "$staged" ]; then rm -rf -- "$staged"; fi
    rm -rf -- "$build"
}
trap cleanup EXIT
trap 'exit 1' HUP INT TERM

git init -q "$build/source"
git -C "$build/source" fetch --depth=1 https://github.com/mtk177a/agent-sessions.git "$revision" >&2
git -C "$build/source" checkout -q --detach FETCH_HEAD
test "$(git -C "$build/source" rev-parse HEAD)" = "$revision"
(
    cd "$build/source"
    export GOCACHE="$build/cache" GOMODCACHE="$build/modules"
    export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOFLAGS=
    go build -o "$build/agent-sessions" ./cmd/agent-sessions
)

mkdir -p "$bin_dir"
staged=$(mktemp -d "$bin_dir/.agent-sessions.XXXXXXXX")
install -m 755 "$build/agent-sessions" "$staged/agent-sessions"
# A hard link publishes the finished file without replacing an existing path.
ln -h -- "$staged/agent-sessions" "$bin_dir/"
printf '%s\n' "$target"
