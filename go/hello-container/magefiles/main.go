//go:build mage

package main

import (
	"context"
	"fmt"

	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

const (
	golangImage    = "golang:latest"
	baseImage      = "alpine:latest"
	publishAddress = "ghcr.io/kpenfound/hello-container:latest"
)

func Run(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		// Load image
		golang := core.Container().From(golangImage)

		// Set workdir
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
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
		golang := core.Container().From(golangImage)

		// Set workdir
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
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

func Push(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		// Load image
		builder := core.Container().From(golangImage)

		// Set workdir
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
		builder = builder.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// Execute Command
		builder = builder.Exec(api.ContainerExecOpts{
			Args: []string{"go", "build", "-o", "hello"},
		})

		// Get built binary
		helloBin, err := builder.File("/src/hello").ID(ctx)
		if err != nil {
			return err
		}

		// Get base image for publishing
		base := core.Container().From(baseImage)
		// Add built binary to /bin
		base = base.WithMountedFile("/tmp/hello", helloBin)
		// Copy mounted file to rootfs
		base = base.Exec(api.ContainerExecOpts{
			Args: []string{"cp", "/tmp/hello", "/bin/hello"},
		})
		// Set entrypoint
		base = base.WithEntrypoint([]string{"/bin/hello"})
		// Publish image
		addr, err := base.Publish(ctx, publishAddress)
		if err != nil {
			return err
		}

		fmt.Println(addr)

		return nil
	}); err != nil {
		panic(err)
	}
}
