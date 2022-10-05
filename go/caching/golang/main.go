package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/sdk/go/dagger"
)

func (r *golang) do(ctx context.Context, mount dagger.FSID, opts GoArgs) (*dagger.Filesystem, error) {
	client, err := dagger.Client(ctx)
	if err != nil {
		return nil, err
	}

	version := "latest"
	if opts.Version != nil {
		version = *opts.Version
	}

	imageRef := fmt.Sprintf("golang:%s", version)

	fsid, err := image(ctx, client, imageRef)
	if err != nil {
		fmt.Printf("cant load image: %v", err)
		return nil, err
	}

	execid, err := exec(ctx, client, fsid, mount, opts)
	if err != nil {
		fmt.Printf("cant execute command: %v", err)
		return nil, err
	}

	return &dagger.Filesystem{ID: execid}, nil
}

func image(ctx context.Context, client graphql.Client, ref string) (dagger.FSID, error) {
	req := &graphql.Request{
		Query: `
query Image ($ref: String!) {
	core {
		image(ref: $ref) {
			id
		}
	}
}
`,
		Variables: map[string]any{
			"ref": ref,
		},
	}
	resp := struct {
		Core struct {
			Image struct {
				ID dagger.FSID
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Core.Image.ID, nil
}

func exec(ctx context.Context, client graphql.Client, root dagger.FSID, mount dagger.FSID, opts GoArgs) (dagger.FSID, error) {
	os := runtime.GOOS
	if opts.Os != nil {
		os = *opts.Os
	}
	arch := runtime.GOARCH
	if opts.Arch != nil {
		arch = *opts.Arch
	}

	req := &graphql.Request{
		Query: `
query ($root: FSID!, $mount: FSID!, $args: [String!]!, $os: String!, $arch: String!) {
	core {
		filesystem(id: $root) {
			exec(input: {
				args: $args,
				workdir: "/src",
				mounts: [
					{
						path: "/src",
						fs: $mount
					}
				],
				env: [
					{
						name: "GOARCH",
						value: $arch
					},
					{
						name: "GOOS",
						value: $os
					},
					{
						name: "GOMODCACHE",
						value: "/cache"
					}
				],
				cacheMounts:{name:"gomod", path:"/cache", sharingMode:"locked"},
			}) {
				fs {
					id
				}
			}
		}
	}
}
`,
		Variables: map[string]any{
			"root":  root,
			"mount": mount,
			"args":  opts.Args,
			"os":    os,
			"arch":  arch,
		},
	}
	resp := struct {
		Core struct {
			Filesystem struct {
				Exec struct {
					FS struct {
						ID dagger.FSID
					}
				}
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Core.Filesystem.Exec.FS.ID, nil
}
