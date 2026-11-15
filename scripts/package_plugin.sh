#!/bin/sh
set -eu

if [ "$#" -ne 1 ] || [ -z "$1" ]; then
    printf '%s\n' 'Usage: sh scripts/package_plugin.sh <new-marketplace-directory>' >&2
    exit 2
fi

repo=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
case "$1" in
    /*) marketplace=$1 ;;
    *) marketplace=$PWD/$1 ;;
esac
package=$marketplace/plugins/agent-harness-retro

mkdir "$marketplace"
mkdir -p "$package/skills"
mkdir -p "$marketplace/.agents/plugins"
cp "$repo/plugin/marketplace.json" "$marketplace/.agents/plugins/marketplace.json"
cp "$repo/plugin/plugin.json" "$repo/LICENSE" "$package/"

for name in agent-harness-retro-setup outer-harness-retrospective; do
    source=$repo/skills/$name
    destination=$package/skills/$name
    mkdir "$destination"
    cp "$source/SKILL.md" "$source/SKILL-ja.md" "$destination/"
    for directory in agents references scripts; do
        if [ -d "$source/$directory" ]; then
            cp -R "$source/$directory" "$destination/"
        fi
    done
done

printf '%s\n' "$marketplace"
