#!/bin/bash
echo "name: hello" >> cloak.yaml

RUN="dagger do -f ./queries.graphql"

WORKFS=$($RUN Workdir | jq -r '.host.workdir.read.id')

TASK="hello"
HELLO=$($RUN --set "workfs=$WORKFS" --set "task=$TASK" Make | jq -r '.. | .contents? | select(. != null)')
echo $HELLO

#cleanup
rm cloak.yaml