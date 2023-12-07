#!/bin/bash
set -euo pipefail
echo "test" >testfile.txt
echo "$GITHUB_TOKEN" | oras login ghcr.io -u github --password-stdin
echo
echo "Trying dummy upload to container registry"
echo
dummy=ghcr.io/syself/autopilot/staging/dummy-upload-test:0.0.1
if ! oras push "$dummy" --artifact-type text/plain testfile.txt; then
    echo "push to container registry failed"
    echo "First chars of GITHUB_TOKEN: ${GITHUB_TOKEN:0:10}..."
    exit 1
fi
