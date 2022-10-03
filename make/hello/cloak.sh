#!/bin/bash
touch cloak.yaml

WORK=$(cloak do<<'EOF'
query {
    host {
        workdir {
            read {
                id
            }
        }
    }
}
EOF
)
WORKFS=$(echo -n $WORK | jq -r '.host.workdir.read.id')

HELLO=$(cloak --set "workfs=$WORKFS" do<<'EOF'
query ($workfs: FSID!) {
  core {
    image(ref: "alpine") {
      exec(input: {
        args: ["apk", "add", "make"]
      }) {
        fs {
            exec(input: {
            args: ["make", "hello"]
            workdir: "/src"
            mounts: [
                {
                    path: "/src"
                    fs: $workfs
                }
            ]
        }) {
            stdout
        }
        }
      }
    }
  }
}
EOF
)
OUT=$(echo -n $HELLO | jq -r '.core.image.exec.fs.exec.stdout')
echo $OUT

#cleanup
rm cloak.yaml