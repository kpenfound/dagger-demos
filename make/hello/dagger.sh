#!/bin/bash
echo "name: hello" >> cloak.yaml

RUN="dagger do -f ./queries.graphql"

WORKFS=$($RUN Workdir | jq -r '.host.workdir.read.id')

TASK="hello"
HELLO=$($RUN --set "workfs=$WORKFS" --set "task=$TASK" Make | jq -r '.container.from.exec.withMountedDirectory.withWorkdir.exec.stdout.contents')
echo $HELLO

#cleanup
rm cloak.yaml