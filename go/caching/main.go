package main

import (
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/aws/aws-sdk-go-v2/aws"
	vault "github.com/hashicorp/vault/api"
)

func main() {
	say := greeting()

	// Pull in a bunch of sdks
	_ = aws.Bool(false)
	_ = &compute.InstancesClient{}
	_ = azcore.AccessToken{}
	_ = vault.DefaultConfig()

	fmt.Println(say)
}

func greeting() string {
	return "Hello!"
}
