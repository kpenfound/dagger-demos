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

		req := &graphql.Request{
			Query: `
query {
	hello {
		run(opts: {})
	}
}
`,
		}
		resp := struct {
			Core struct {
				Hello struct {
					Run string
				}
			}
		}{}
		err = client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
		if err != nil {
			return err
		}

		fmt.Printf(resp.Core.Hello.Run)
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

		req := &graphql.Request{
			Query: `
query {
	hello {
		test(opts: {})
	}
}
`,
		}
		resp := struct {
			Core struct {
				Hello struct {
					Test string
				}
			}
		}{}
		err = client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
		if err != nil {
			return err
		}

		fmt.Printf(resp.Core.Hello.Test)
		return nil
	}); err != nil {
		panic(err)
	}
}
