package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	err := listImages()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func listImages() error {

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	images, err := client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return err
	}

	for _, summary := range images {
		for _, tag := range summary.RepoTags {
			fmt.Printf("%s (size: %d)\n", tag, summary.Size)
		}
	}

	return nil
}
