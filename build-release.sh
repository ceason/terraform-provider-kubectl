#!/usr/bin/env bash
[ "$DEBUG" = "1" ] && set -x
set -euo pipefail
err_report() { echo "errexit on line $(caller)" >&2; }
trap err_report ERR


tag=$(git describe --tags --abbrev=0)

rm -rf gh-release
mkdir gh-release

for platform in linux windows darwin; do
	export GOOS=$platform
	export GOARCH=amd64
	outfile=gh-release/terraform-provider-kubectl-$tag-${GOOS}_${GOARCH}
	if [ $platform = "windows" ]; then
		outfile="$outfile.exe"
	fi
	go build -ldflags="-s -w" -o "$outfile" .
done

tree gh-release

