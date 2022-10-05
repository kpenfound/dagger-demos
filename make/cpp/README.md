# hello-cpp

This is a silly example that builds postgresql from source to demonstrate build caching

Note: run `make delete` to cleanup `bin/` directory

Note: The caching is currently broken between the build and copy steps. It runs the build twice.

## Timings

`make postgresql`: about 2m30s every execution on m1
`./dagger.sh` fully uncached: 
`./dagger.sh` uncached workdir: 2m51s
`./dagger.sh` fully cached: 4s
