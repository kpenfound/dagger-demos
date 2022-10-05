# make-hello

This is a project that uses a Makefile to run tasks

`make hello`

To onboard this to dagger, an easy first step is to have dagger run make. This is a good first step because it doesn't require rewriting any existing tooling, but you immediately get the benefits of local ci, strongly defined environments, and caching.

In [dagger.sh](./dagger.sh) we create a simple query to load the working directory, create a make runtime, and run the make target.

`./dagger.sh`