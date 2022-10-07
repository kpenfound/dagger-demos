package ci

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger"
)

type NetlifyDeployConfig struct {
	NetlifyToken  string
	NetlifySite   string
	Repository    string
	RepositoryRef string
}

func NetlifyDeploy(cfg *NetlifyDeployConfig) (string, error) {
	return netlifyDeployInDagger(cfg)
}

func netlifyDeployInDagger(cfg *NetlifyDeployConfig) (string, error) {
	ctx := context.Background()
	url := ""
	engineErr := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		secret, err := addSecret(ctx, client, cfg.NetlifyToken)
		if err != nil {
			return err
		}

		url, err = netlifyBuildAndDeploy(ctx, client, cfg.Repository, cfg.RepositoryRef, cfg.NetlifySite, secret)
		return err
	})
	return url, engineErr
}

func addSecret(ctx context.Context, client graphql.Client, secret string) (dagger.SecretID, error) {
	req := &graphql.Request{
		Query: `
		query ($secret: String!) {
			core {
			  addSecret(plaintext: $secret)
			}
		  }
		`,
		Variables: map[string]any{
			"secret": secret,
		},
	}
	resp := struct {
		Core struct {
			AddSecret dagger.SecretID
		}
	}{}

	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})

	return resp.Core.AddSecret, err
}

func netlifyBuildAndDeploy(ctx context.Context, client graphql.Client, remote, ref, site string, secret dagger.SecretID) (string, error) {
	req := &graphql.Request{
		Query: `
		  query NetlifyBuildAndDeploy($remote: String!, $ref: String!, $site: String!, $secret: SecretID!) {
			core {
			  git(remote: $remote, ref: $ref) {
				yarn(runArgs: ["build"]) {
				  netlifyDeploy(
					token: $secret
					subdir: "build"
					siteName: $site
				  ) {
					deployURL
				  }
				}
			  }
			}
		  }`,
		Variables: map[string]any{
			"remote": remote,
			"ref":    ref,
			"site":   site,
			"secret": secret,
		},
	}
	resp := struct {
		Core struct {
			Git struct {
				Yarn struct {
					NetlifyDeploy struct {
						DeployURL string
					}
				}
			}
		}
	}{}

	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})

	return resp.Core.Git.Yarn.NetlifyDeploy.DeployURL, err
}
