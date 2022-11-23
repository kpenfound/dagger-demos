package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
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
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get working directory on host
	src := client.Host().Directory(".")

	// initialize new container from image
	golang := client.Container().
		From(golangImage).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "main.go"})

	// get command output
	out, err := golang.Stdout(ctx)
	if err != nil {
		panic(err)
	}

	// print output to console
	fmt.Println(out)
}

func test(ctx context.Context) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get working directory on host
	src := client.Host().Directory(".")

	// initialize new container from image
	golang := client.Container().
		From(golangImage).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test"})

	// get command output
	out, err := golang.Stdout(ctx)
	if err != nil {
		panic(err)
	}

	// print output to console
	fmt.Println(out)
}

func publish(ctx context.Context) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get working directory on host
	src := client.Host().Directory(".")

	// initialize new container from image
	builder := client.Container().
		From(golangImage).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "hello"})

	// initialize new container for publishing from image
	base := client.Container().From(baseImage)

	// mount binary file at container path
	base = base.WithRootfs(
		base.Rootfs().WithFile(
			"/bin/hello",
			builder.File("/src/hello"),
		),
	).WithEntrypoint([]string{"/bin/hello"})

	// publish image
	addr, err := base.Publish(ctx, publishAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}
