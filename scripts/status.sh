#!/usr/bin/env bash

# This command is used by bazel as the workspace_status_command
# to implement build stamping with git information.

RELEASE_VERSION=$(git describe --tags --dirty --always)
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

cat <<EOF
STABLE_VERSION ${RELEASE_VERSION}
STABLE_BUILD_TIME ${BUILD_TIME}
STABLE_GIT_COMMIT ${GIT_COMMIT}
STABLE_GIT_BRANCH ${GIT_BRANCH}
EOF
