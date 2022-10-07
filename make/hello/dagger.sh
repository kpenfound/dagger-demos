#!/bin/bash
touch cloak.yaml

RUN="dagger do -f ./queries.graphql"

WORK=$($RUN Workdir)
WORKFS=$(echo -n $WORK | jq -r '.host.workdir.read.id')

TASK="hello"
HELLO=$($RUN --set "workfs=$WORKFS" --set "task=$TASK" Make)
OUT=$(echo -n $HELLO | jq -r '.core.image.exec.fs.exec.stdout')
echo $OUT

#cleanup
rm cloak.yaml