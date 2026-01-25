package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	imageName := os.Getenv("IMAGE_NAME")
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

	auth, err := authFromPassword(username, password)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = pushImage(imageName, auth)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func pushImage(imageName string, auth string) error {

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	opts := image.PushOptions{
		RegistryAuth: auth,
	}
	reader, err := client.ImagePush(ctx, imageName, opts)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Info(line)
	}

	return nil
}

func authFromPassword(username string, password string) (string, error) {

	authConfig := registry.AuthConfig{
		Username: username,
		Password: password,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
