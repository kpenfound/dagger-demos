#!/bin/bash
touch cloak.yaml

RUN="dagger do -f ./queries.graphql"

WORKFS=$($RUN Workdir | jq -r '.host.workdir.read.id')

TASK="hello"
HELLO=$($RUN --set "workfs=$WORKFS" --set "task=$TASK" Make | jq -r '.core.image.exec.fs.exec.stdout')
echo $HELLO

#cleanup
rm cloak.yaml