# go-cli

## The demo

At a lot of organizations, the DevOps/Infra/RelEng/DevEx team will maintain and distribute a binary that handles the CICD operations for that organization.

This is an example of what that looks like with Dagger.

At Kyle Inc., there are several teams that work on static sites which are built with `yarn build` and deployed to Netlify. To manage these deployments, the `kyle-ci` binary has a `netlify-deploy` subcommand. The subcommand takes arguments such as the git repository and netlify site name.

When a developer is ready to see their branch/release/etc on Netlify, they use the `kyle-ci netlify-deploy` command:

```
./kyle-ci netlify-deploy -r https://github.com/dagger/todoapp -s kyle-todoapp

Deploying...
...
#45 DONE 12.1s
https://github.com/dagger/todoapp#main deployed to http://634068a4528a422ed31c8f7e--kyle-todoapp.netlify.app
```

And that happened in Dagger!

Going a step further, that means the whole organization can use a shared buildkit environment to take advantage of a shared cache!

## Try it yourself

1. Build the binary with `go build -o kyle-ci`
2. export your netlify token `export NETLIFY_TOKEN=xxxx`
3. Run `./kyle-ci netlify-deploy -r https://github.com/dagger/todoapp -s $USER-todoapp`

The demo does have some guardrails. It assumes the repo passed in will produce a `/build` output from a `yarn build`.
