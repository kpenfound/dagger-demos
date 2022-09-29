# hello-extension

A simple hello world app with a Dagger extension

This expands on the basic version [here](../hello/).

Try it out!

- `mage run`
- `mage test`

Dagger is executed in the magefile `magefiles/main.go`

The `Run` and `Test` tasks each load the dagger engine, and run our extension located in `dagger/`.

The extension performs the same operations as the more basic example, but bundles them in an extension instead of requiring the magefile to do all the work.