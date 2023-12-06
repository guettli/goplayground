package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("One argument needed: url. Example ghcr.io/foo/bar:tag")
		os.Exit(1)
	}
	url := os.Args[1]
	parts := strings.Split(url, ":")
	if len(parts) == 1 {
		fmt.Println("Missing :tag in url")
		os.Exit(1)
	}
	url = strings.Join(parts[0:len(parts)-1], ":")
	tag := parts[len(parts)-1]
	fmt.Printf("url: %s tag: %s\n", url, tag)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Need env var GITHUB_TOKEN")
		os.Exit(1)
	}

	ctx := context.Background()

	orasClient := &auth.Client{
		// expectedHostAddress is of form ipaddr:port
		Credential: auth.StaticCredential("ghcr.io", auth.Credential{
			Username:    "github",
			AccessToken: token,
		}),
	}

	repo, err := remote.NewRepository(url)
	// fmt.Printf("repo: %s/%s\n", hetznerNodeImageRegistry, strings.TrimSuffix(fileName, "-release.json"))
	if err != nil {
		log.Fatal(err)
	}
	repo.Client = orasClient

	descriptor, err := repo.Resolve(ctx, tag)
	if err != nil {
		log.Fatalf("failed to resolve tag: %q", err.Error())
	}
	exist, err := repo.Exists(ctx, descriptor)
	if err != nil {
		log.Fatal(err)
	}
	if exist {
		log.Println("already exists, Refusing to overwrite it")
	}
}
