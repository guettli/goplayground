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
	if strings.HasPrefix(url, "http") {
		fmt.Println("url should not start with http. Example: ghcr.io/foo/bar:tag")
		os.Exit(1)
	}
	parts := strings.Split(url, ":")
	if len(parts) == 1 {
		fmt.Println("Missing :tag in url")
		os.Exit(1)
	}
	url = strings.Join(parts[0:len(parts)-1], ":")
	tag := parts[len(parts)-1]

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Need env var GITHUB_TOKEN")
		os.Exit(1)
	}

	ctx := context.Background()

	parts = strings.Split(url, "/")
	if len(parts) == 1 {
		fmt.Printf("failed to parse url %q\n", url)
		os.Exit(1)
	}
	regName := parts[0]

	repo, err := remote.NewRepository(url)
	if err != nil {
		log.Fatal("Repository failed", err)
	}

	repo.Client = &auth.Client{
		// expectedHostAddress is of form ipaddr:port
		Credential: auth.StaticCredential(regName, auth.Credential{
			Username: "github",
			Password: token,
		}),
	}

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
