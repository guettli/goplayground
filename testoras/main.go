package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	credentials "github.com/oras-project/oras-credentials-go"
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
	//fmt.Printf("url: %s tag: %s\n", url, tag)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Need env var GITHUB_TOKEN")
		os.Exit(1)
	}

	ctx := context.Background()

	reg, err := remote.NewRegistry("ghcr.io")
	if err != nil {
		log.Fatalf("NewRegistry failed: %q", err.Error())
	}

	repo, err := reg.Repository(ctx, "ghcr.io")
	if err != nil {
		log.Fatal("Repository failed", err)
	}

	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		log.Fatalf("NewStoreFromDocker failed: %q", err.Error())
	}

	err = credentials.Login(ctx, credStore, reg, auth.Credential{
		Username: "guettli",
		Password: token,
	})
	if err != nil {
		fmt.Printf("Login failed: %q\n", err.Error())
		os.Exit(1)
	}
	err = reg.Ping(ctx)
	if err != nil {
		log.Fatalf("Ping failed: %q", err.Error())
	}
	///////

	descriptor, err := repo.Resolve(ctx, tag)
	if err != nil {
		log.Fatalf("failed to resolve tag: %q", err.Error())
	}
	exist, err := repo.Exists(ctx, descriptor)
	if err != nil {
		log.Fatal(err)
	}
	if exist {
		fmt.Printf("%s exists\n", url)
	} else {
		fmt.Printf("%s does not exist\n", url)
	}
}
