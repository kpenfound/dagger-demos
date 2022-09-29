//go:build mage

package main

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger"
)

func Run(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		goImage, err := image(ctx, client, "golang:latest")
		if err != nil {
			return err
		}

		work, err := workdir(ctx, client)
		if err != nil {
			return err
		}

		out, err := exec(
			ctx,
			client,
			goImage,
			work,
			[]string{"go", "run", "main.go"},
		)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	}); err != nil {
		panic(err)
	}
}

func Test(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		goImage, err := image(ctx, client, "golang:latest")
		if err != nil {
			return err
		}

		work, err := workdir(ctx, client)
		if err != nil {
			return err
		}

		out, err := exec(
			ctx,
			client,
			goImage,
			work,
			[]string{"go", "test"},
		)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	}); err != nil {
		panic(err)
	}
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

func exec(ctx context.Context, client graphql.Client, root dagger.FSID, mount dagger.FSID, args []string) (string, error) {
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
