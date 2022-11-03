//go:build mage

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type GoArgs struct {
	Version *string  `json:"version"`
	Arch    *string  `json:"arch"`
	Os      *string  `json:"os"`
	Args    []string `json:"args"`
}

func Run(ctx context.Context) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	work := client.Host().Workdir()

	args := []string{"go", "run", "main.go"}
	out, err := exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func Test(ctx context.Context) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	work := client.Host().Workdir()

	args := []string{"go", "test"}
	out, err := exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func Benchmark(ctx context.Context) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	work := client.Host().Workdir()

	args := []string{"go", "test"}
	out, err := exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)

	args = []string{"go", "run", "main.go"}
	out, err = exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)

	args = []string{"go", "build"}
	out, err = exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)

	args = []string{"go", "install"}
	out, err = exec(ctx, client, work, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func exec(ctx context.Context, client *dagger.Client, source dagger.Directory, args []string) (string, error) {
	container := client.Container().From("golang:latest")
	container = container.WithMountedDirectory("/src", source).WithWorkdir("/src")

	// Enable or disable mod caching with CACHING_ENABLED=1
	if shouldCache() == "1" {
		cacheKey := "gomods"
		cache := client.CacheVolume(cacheKey)

		container = container.WithMountedCache("/cache", cache)
		container = container.WithEnvVariable("GOMODCACHE", "/cache")
	}

	container = container.Exec(dagger.ContainerExecOpts{
		Args: args,
	})
	return container.Stdout().Contents(ctx)
}

func shouldCache() string {
	return os.Getenv("CACHING_ENABLED")
}
