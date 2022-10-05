//go:build mage

package main

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger"
)

type GoArgs struct {
	Version *string  `json:"version"`
	Arch    *string  `json:"arch"`
	Os      *string  `json:"os"`
	Args    []string `json:"args"`
}

func Run(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		work, err := workdir(ctx, client)
		if err != nil {
			return err
		}

		args := GoArgs{
			Args: []string{"go", "run", "main.go"},
		}
		out, err := exec(
			ctx,
			client,
			work,
			args,
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

		work, err := workdir(ctx, client)
		if err != nil {
			return err
		}

		args := GoArgs{
			Args: []string{"go", "test"},
		}
		out, err := exec(
			ctx,
			client,
			work,
			args,
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

func Benchmark(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		work, err := workdir(ctx, client)
		if err != nil {
			return err
		}

		args := GoArgs{
			Args: []string{"go", "test"},
		}
		out, err := exec(
			ctx,
			client,
			work,
			args,
		)
		if err != nil {
			return err
		}

		fmt.Println(out)

		args = GoArgs{
			Args: []string{"go", "run", "main.go"},
		}
		out, err = exec(
			ctx,
			client,
			work,
			args,
		)
		if err != nil {
			return err
		}

		fmt.Println(out)

		args = GoArgs{
			Args: []string{"go", "build"},
		}
		out, err = exec(
			ctx,
			client,
			work,
			args,
		)
		if err != nil {
			return err
		}

		fmt.Println(out)

		args = GoArgs{
			Args: []string{"go", "install"},
		}
		out, err = exec(
			ctx,
			client,
			work,
			args,
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

func exec(ctx context.Context, client graphql.Client, mount dagger.FSID, args GoArgs) (string, error) {
	req := &graphql.Request{
		Query: `
query ($mount: FSID!, $args: GoArgs!) {
	golang {
		do(project: $mount, opts: $args) {
			id
		}
	}
}
`,
		Variables: map[string]any{
			"mount": mount,
			"args":  args,
		},
	}
	resp := struct {
		Golang struct {
			Do struct {
				Fileystem struct {
					ID string
				}
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Golang.Do.Fileystem.ID, nil
}
