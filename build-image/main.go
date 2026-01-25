package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"github.com/moby/go-archive"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	imageName := os.Getenv("IMAGE_NAME")

	err = buildImage(imageName)
	if err != nil {
		slog.Error(err.Error())
	}
}

func buildImage(imageName string) error {

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer client.Close()

	err = removeImageIfExists(client, imageName)
	if err != nil {
		return err
	}

	tar, err := archive.TarWithOptions("context/", &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := build.ImageBuildOptions{
		Dockerfile:  "Dockerfile",
		Tags:        []string{imageName},
		Remove:      true,
		ForceRemove: true,
	}

	ctx := context.Background()
	res, err := client.ImageBuild(ctx, tar, opts)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Info(line)
	}

	return nil
}

func removeImageIfExists(client *client.Client, imageName string) error {

	ctx := context.Background()
	args := filters.NewArgs(filters.KeyValuePair{Key: "reference", Value: imageName})
	images, err := client.ImageList(ctx, image.ListOptions{Filters: args})
	if err != nil {
		return err
	}

	opts := image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	}
	for _, summary := range images {
		_, err = client.ImageRemove(ctx, summary.ID, opts)
		if err != nil {
			return err
		}
		slog.Info("Deleted an image", "id", summary.ID)
	}

	return nil
}
