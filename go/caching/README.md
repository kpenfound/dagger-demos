# caching

A simple hello world app with a Golang to demonstrate mod caching

Try it out!

- `mage run`
- `mage test`

Dagger is executed in the magefile `magefiles/main.go`

The `Run`, `Test`, and `Benchmark` tasks each connect to the dagger engine and run a pipeline.

## Performance:
In both cases, a fresh workdir cache was used

_Without_ module caching, `mage benchmark` took 2m57s.

_With_ module caching, `CACHING_ENABLED=1 mage benchmark` took 1m23s.

This is accomplished through `Container.WithMountedCache` and setting the Go module cache dir to the cache.
