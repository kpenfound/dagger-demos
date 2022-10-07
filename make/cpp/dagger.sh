#!/bin/bash
touch cloak.yaml

RUN="dagger do -f ./queries.graphql"

echo "Loading workdir"
WORKFS=$($RUN Workdir | jq -r '.host.workdir.read.id')

echo "Preparing Builder"
BUILDER=$($RUN Builder | jq -r '.core.image.exec.fs.exec.fs.id')

echo "Doing postgresql build"
BUILDFS=$($RUN Build --set "workfs=$WORKFS" --set "builderid=$BUILDER" | jq -r '.core.filesystem.exec.mount.id')

echo "Copying build outputs"
CP=$($RUN Copy --local-dir bin=./bin --set "buildfs=$BUILDFS")

ls -la bin

#cleanup
rm cloak.yaml
