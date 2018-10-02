#!/usr/bin/env bash
[ "$DEBUG" = "1" ] && set -x
set -euo pipefail
err_report() { echo "errexit on line $(caller)" >&2; }
trap err_report ERR

# get
PRERELEASE_TESTS=(
  //...
  //tests:create_delete_integration_test
)

cd $(dirname "$0")

bazel build //...
bazel test "${PRERELEASE_TESTS[@]}"

if git diff-index --cached --quiet HEAD --; then
	echo "No staged changes"
else
	git commit
fi
commit=$(git rev-parse --verify HEAD)

git fetch --tags
tag=$(git describe --tags --abbrev=0)
major=$(cut -f1 -d'.' <<<"$tag")
minor=$(cut -f2 -d'.' <<<"$tag")
patch=$(git tag|grep "^$major.$minor."|sort --version-sort|tail -1|cut -f3 -d'.')
tag="$major.$minor.$((patch+1))"

rm -rf gh-release
mkdir -p gh-release
trap "rm -rf $PWD/gh-release" EXIT

for platform in linux windows darwin; do
	export GOOS=$platform
	export GOARCH=amd64
	outfile=gh-release/terraform-provider-kubectl-$tag-${GOOS}_${GOARCH}
	if [ $platform = "windows" ]; then
		outfile="$outfile.exe"
	fi
	go build -ldflags="-s -w" -o "$outfile" .
done

git push
hub release create \
  --commitish="$commit" \
  --attach=gh-release/terraform-provider-kubectl-$tag-linux_amd64 \
  --attach=gh-release/terraform-provider-kubectl-$tag-darwin_amd64 \
  --attach=gh-release/terraform-provider-kubectl-$tag-windows_amd64.exe \
  --message="$tag" \
  --browse \
  "$tag"

git fetch --tags

