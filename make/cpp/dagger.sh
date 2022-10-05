#!/bin/bash
touch cloak.yaml

echo "Loading workdir"
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

echo "Preparing Builder"
BUILDER=$(cloak do<<'EOF'
query {
  core {
    image(ref: "debian:stable") {
      exec(input: {
        args: ["apt-get", "update"]
      }) {
        fs {
          exec(input: {
            args: ["apt-get", "install", "-y", "build-essential", "openssl", "curl", "libreadline-dev", "libpython-all-dev", "libperl-dev", "zlib1g-dev"]
          }) {
            fs {
              id
            }
          }
        }
      }
    }
  }
}
EOF
)
BUILDERID=$(echo -n $BUILDER | jq -r '.core.image.exec.fs.exec.fs.id')

echo "Doing postgresql build"
BUILD=$(cloak --set "workfs=$WORKFS" --set "builderid=$BUILDERID" do<<'EOF'
query ($builderid: FSID!, $workfs: FSID!) {
  core {
    filesystem(id: $builderid) {
      exec(input: {
        args: ["make", "postgresql"]
        workdir: "/src"
        mounts: [
          {
            path: "/src"
            fs: $workfs
          }
        ]
      }) {
        mount(path: "/src") {
          id
        }
      }
    }
  }
}
EOF
)
BUILDFS=$(echo -n $BUILD | jq -r '.core.filesystem.exec.mount.id')

echo "Copying build outputs"
CP=$(cloak --local-dir bin=./bin --set "buildfs=$BUILDFS" do<<'EOF'
query($buildfs: FSID!) {
  host {
    dir(id: "bin") {
      write(contents: $buildfs, path: ".")
    }
  }
}
EOF
)

ls -la bin

#cleanup
rm cloak.yaml
