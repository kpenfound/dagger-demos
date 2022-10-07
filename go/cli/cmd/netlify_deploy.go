package cmd

import (
	"fmt"
	"os"

	"github.com/kpenfound/dagger-demos/ci"
	"github.com/spf13/cobra"
)

var netlifyCmd = &cobra.Command{
	Use: "netlify-deploy",
	Run: NetlifyDeploy,
}

func NetlifyDeploy(cmd *cobra.Command, args []string) {
	fmt.Printf("Deploying...\n")

	netlifyToken := os.Getenv("NETLIFY_TOKEN")
	if netlifyToken == "" {
		fmt.Printf("NETLIFY_TOKEN must be set to deploy to Netlify")
		os.Exit(1)
	}

	if site == "" {
		site = "kyleci"
	}

	if repoRef == "" {
		repoRef = "main"
	}

	cfg := &ci.NetlifyDeployConfig{
		NetlifyToken:  netlifyToken,
		NetlifySite:   site,
		Repository:    repository,
		RepositoryRef: repoRef,
	}

	url, err := ci.NetlifyDeploy(cfg)
	if err != nil {
		fmt.Errorf("unable to deploy to netlify: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("%s#%s deployed to %s\n", repository, repoRef, url)
}
