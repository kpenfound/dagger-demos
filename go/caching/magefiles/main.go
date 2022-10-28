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

	work, err := workdir(ctx, client)
	if err != nil {
		panic(err)
	}

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

	work, err := workdir(ctx, client)
	if err != nil {
		panic(err)
	}

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

	work, err := workdir(ctx, client)
	if err != nil {
		panic(err)
	}

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

func workdir(ctx context.Context, client *dagger.Client) (dagger.DirectoryID, error) {
	return client.Host().Workdir().Read().ID(ctx)
}

func exec(ctx context.Context, client *dagger.Client, source dagger.DirectoryID, args []string) (string, error) {
	container := client.Container().From("golang:latest")
	container = container.WithMountedDirectory("/src", source).WithWorkdir("/src")

	// Enable or disable mod caching with CACHING_ENABLED=1 sdfsdfsdfvsdljksdfljsdf
	if shouldCache() == "1" {
		cacheKey := "gomods"
		cacheID, err := client.CacheVolume(cacheKey).ID(ctx)
		if err != nil {
			return "", err
		}

		container = container.WithMountedCache(cacheID, "/cache")
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
