# getting-started

A simple getting started app

## Requirements

- Go - [install here](https://go.dev/doc/install)
- Docker - [install here](https://docs.docker.com/engine/install/)

## Create a simple Go app

- Create a go project following [instructions here](https://go.dev/doc/code#Command)
- You should now have a go.mod, go.sum, and main.go
- Run your app

```
> go run main.go
Hello
```

## Add Dagger

- `go get dagger.io/dagger@<tag>`
- `go mod edit -replace github.com/docker/docker=github.com/docker/docker@v20.10.3-0.20220414164044-61404de7df1a+incompatible`
- `go mod tidy`
- Create ci/main.go and build out `func main()` and `func run()` as in [ci/main.go](ci/main.go)
- Try it out! `go run ci/main.go run`
- Build out `func test()` and `func push()`
- `go run ci/main.go test`
- `go run ci/main.go push`
