package main

import (
	"context"
	"fmt"
	"os"

	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

const (
	golangImage    = "golang:latest"
	baseImage      = "alpine:latest"
	publishAddress = "ghcr.io/kpenfound/hello-container:latest"
)

func main() {
	ctx := context.Background()

	task := os.Args[1]

	if len(os.Args) != 2 {
		fmt.Println("Please pass a task as an argument")
		os.Exit(1)
	}

	switch task {
	case "run":
		run(ctx)
	case "test":
		test(ctx)
	case "publish":
		publish(ctx)
	default:
		fmt.Printf("Unknown task %s\n", task)
		os.Exit(1)
	}
}

func run(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)
		
		// get working directory on host
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
		
		// initialize new container from image
		golang := core.Container().From(golangImage)

		// mount working directory to /src
		golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// execute command
		cmd := golang.Exec(api.ContainerExecOpts{
			Args: []string{"go", "run", "main.go"},
		})

		// get command output
		out, err := cmd.Stdout().Contents(ctx)
		if err != nil {
			return err
		}

		// print output to console
		fmt.Println(out)

		return nil
	}); err != nil {
		panic(err)
	}
}

func test(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		// get working directory on host
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
		
		// initialize new container from image
		golang := core.Container().From(golangImage)
		
		// mount working directory to /src
		golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// execute command
		cmd := golang.Exec(api.ContainerExecOpts{
			Args: []string{"go", "test"},
		})

		// get command output
		out, err := cmd.Stdout().Contents(ctx)
		if err != nil {
			return err
		}

		// print output to console
		fmt.Println(out)

		return nil
	}); err != nil {
		panic(err)
	}
}

func publish(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		// get working directory on host
		src, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}
		
		// initialize new container from image
		builder := core.Container().From(golangImage)

		// mount working directory to /src
		builder = builder.WithMountedDirectory("/src", src).WithWorkdir("/src")

		// execute build command
		builder = builder.Exec(api.ContainerExecOpts{
			Args: []string{"go", "build", "-o", "hello"},
		})

		// get built binary file
		helloBin, err := builder.File("/src/hello").ID(ctx)
		if err != nil {
			return err
		}

		// initialize new container for publishing from image
		base := core.Container().From(baseImage)
		
		// mount binary file at container path
		base = base.WithMountedFile("/tmp/hello", helloBin)
		
		// copy mounted file to container filesystem
		base = base.Exec(api.ContainerExecOpts{
			Args: []string{"cp", "/tmp/hello", "/bin/hello"},
		})
		
		// set container entrypoint
		base = base.WithEntrypoint([]string{"/bin/hello"})
		
		// publish image
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
