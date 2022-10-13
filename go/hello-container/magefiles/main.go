//go:build mage

package main

import (
	"context"
	"fmt"

	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

func Run(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		// Load image
		golang := core.Container().From("golang:latest")

		// Set workdir
		src := api.DirectoryID(ctx.Workdir) // hacky cast cannot use ctx.Workdir (variable of type core.DirectoryID) as type api.DirectoryID
		golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// Execute Command
		cmd := golang.Exec(api.ContainerExecOpts{
			Args: []string{"go", "run", "main.go"},
		})

		// Get Command Output
		out, err := cmd.Stdout().Contents(ctx)
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
		core := api.New(ctx.Client)

		// Load image
		golang := core.Container().From("golang:latest")

		// Set workdir
		src := api.DirectoryID(ctx.Workdir) // hacky cast cannot use ctx.Workdir (variable of type core.DirectoryID) as type api.DirectoryID
		golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// Execute Command
		cmd := golang.Exec(api.ContainerExecOpts{
			Args: []string{"go", "test"},
		})

		// Get Command Output
		out, err := cmd.Stdout().Contents(ctx)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	}); err != nil {
		panic(err)
	}
}
