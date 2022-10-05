# hello-extension

A simple hello world app with a Golang extension to demonstrate mod caching

This expands on the basic version [here](../hello/).

Try it out!

- `mage run`
- `mage test`

Dagger is executed in the magefile `magefiles/main.go`

The `Run` and `Test` tasks each load the dagger engine, and run our extension located in `golang/`.

## Performance:
In both cases, a fresh workdir cache was used

Using the golang extension _without_ module caching, `mage benchmark` took 1m12s.

Using the golang extension _with_ module caching, `mage benchmark` took 50s.


Roughly 30% time saved just by adding this to the `fs.exec`:

`cacheMounts:{name:"gomod", path:"/cache", sharingMode:"locked"},`
