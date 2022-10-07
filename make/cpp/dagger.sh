#!/bin/bash
touch cloak.yaml

RUN="dagger do -f ./queries.graphql"

echo "Loading workdir"
WORK=$($RUN Workdir)
WORKFS=$(echo -n $WORK | jq -r '.host.workdir.read.id')

echo "Preparing Builder"
BUILDER=$($RUN Builder)
BUILDERID=$(echo -n $BUILDER | jq -r '.core.image.exec.fs.exec.fs.id')

echo "Doing postgresql build"
BUILD=$($RUN Build --set "workfs=$WORKFS" --set "builderid=$BUILDERID")
BUILDFS=$(echo -n $BUILD | jq -r '.core.filesystem.exec.mount.id')

echo "Copying build outputs"
CP=$($RUN Copy --local-dir bin=./bin --set "buildfs=$BUILDFS")

ls -la bin

#cleanup
rm cloak.yaml
