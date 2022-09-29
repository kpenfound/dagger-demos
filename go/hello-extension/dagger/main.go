package main

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/sdk/go/dagger"
)

const goimg = "golang:latest"
const gobin = "/usr/local/go/bin/go"

func (g *hello) build(ctx context.Context, opts GoOpts) (*dagger.Filesystem, error) {
	client, err := dagger.Client(ctx)
	if err != nil {
		return nil, err
	}

	goImage, err := image(ctx, client, goimg)
	if err != nil {
		return nil, err
	}

	work, err := workdir(ctx, client)
	if err != nil {
		return nil, err
	}

	fs, err := execfs(
		ctx,
		client,
		goImage,
		work,
		[]string{gobin, "build"},
	)
	if err != nil {
		return nil, err
	}

	return &dagger.Filesystem{ID: fs}, nil
}

func (g *hello) run(ctx context.Context, opts GoOpts) (string, error) {
	out := ""
	client, err := dagger.Client(ctx)
	if err != nil {
		return out, err
	}

	goImage, err := image(ctx, client, goimg)
	if err != nil {
		return out, err
	}

	work, err := workdir(ctx, client)
	if err != nil {
		return out, err
	}

	return execout(
		ctx,
		client,
		goImage,
		work,
		[]string{gobin, "run", "main.go"},
	)
}

func (g *hello) test(ctx context.Context, opts GoOpts) (string, error) {
	out := ""
	client, err := dagger.Client(ctx)
	if err != nil {
		return out, err
	}

	goImage, err := image(ctx, client, goimg)
	if err != nil {
		return out, err
	}

	work, err := workdir(ctx, client)
	if err != nil {
		return out, err
	}

	return execout(
		ctx,
		client,
		goImage,
		work,
		[]string{gobin, "test"},
	)
}

func workdir(ctx context.Context, client graphql.Client) (dagger.FSID, error) {
	req := &graphql.Request{
		Query: `
query {
	host {
		workdir {
			read {
				id
			}
		}
	}	
}
`,
	}
	resp := struct {
		Host struct {
			Workdir struct {
				Read struct {
					ID dagger.FSID
				}
			}
		}
	}{}

	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})

	return resp.Host.Workdir.Read.ID, err
}

func image(ctx context.Context, client graphql.Client, ref string) (dagger.FSID, error) {
	req := &graphql.Request{
		Query: `
query ($ref: String!) {
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

func execout(ctx context.Context, client graphql.Client, root dagger.FSID, mount dagger.FSID, args []string) (string, error) {
	req := &graphql.Request{
		Query: `
query ($root: FSID!, $mount: FSID!, $args: [String!]!) {
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
			}) {
				stdout
			}
		}
	}
}
`,
		Variables: map[string]any{
			"root":  root,
			"mount": mount,
			"args":  args,
		},
	}
	resp := struct {
		Core struct {
			Filesystem struct {
				Exec struct {
					Stdout string
				}
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Core.Filesystem.Exec.Stdout, nil
}

func execfs(ctx context.Context, client graphql.Client, root dagger.FSID, mount dagger.FSID, args []string) (dagger.FSID, error) {
	req := &graphql.Request{
		Query: `
query ($root: FSID!, $mount: FSID!, $args: [String!]!) {
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
			"args":  args,
		},
	}
	resp := struct {
		Core struct {
			Filesystem struct {
				Exec struct {
					Fs struct {
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

	return resp.Core.Filesystem.Exec.Fs.ID, nil
}
