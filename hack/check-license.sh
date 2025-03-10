#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

echo "Checking for license header..."
allfiles=$(find pkg -name '*.go')
licRes=""
for file in $allfiles; do
  if ! head -n3 "${file}" | grep -Eq "(Copyright|generated|GENERATED|Licensed)" ; then
    licRes="${licRes}\n"$(echo -e "  ${file}")
  fi
done
if [ -n "${licRes}" ]; then
  echo -e "license header checking failed:\n${licRes}"
  exit 255
fi