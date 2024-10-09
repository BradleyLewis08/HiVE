package imager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	dockerv1 "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type Imager struct {
	dockerClient *dockerv1.Client
}

func NewImager() *Imager {
	docker, err := dockerv1.NewClientWithOpts(dockerv1.FromEnv, dockerv1.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating docker client: %s\n", err)
		return nil
	}
	return &Imager{dockerClient: docker}
}

func (i* Imager) createImage(
	courseName string,
	contents string,
) (string, error) {
	tempDir, err := os.MkdirTemp("", "docker-build")

	if err != nil {
		fmt.Printf("Error creating temp dir: %s\n", err)
		return "", err
	}

	// defer os.RemoveAll(tempDir)

	dockerFilePath := fmt.Sprintf("%s/Dockerfile", tempDir)

	// fullDockerfileContents := contents
	fullDockerfileContents := fmt.Sprintf("FROM bradleylewis08/ssh-pod:latest\n%s", contents)
	fmt.Printf("Dockerfile contents: %s\n", fullDockerfileContents)
	err = os.WriteFile(dockerFilePath, []byte(fullDockerfileContents), 0644)
	if err != nil {
		fmt.Printf("Error writing Dockerfile: %s\n", err)
		return "", err
	}

	// Create tar archive of build context
	buildContextTar, err := archive.TarWithOptions(tempDir, &archive.TarOptions{})
	if err != nil {
		fmt.Printf("Error creating tar file: %s\n", err)
		return "", err
	}

	imageName := fmt.Sprintf("bradleylewis08/course-environments:%s", strings.ToLower(courseName))

	// Build image
	resp, err := i.dockerClient.ImageBuild(context.Background(), buildContextTar, types.ImageBuildOptions{
		Tags: []string{imageName},
		Dockerfile: "Dockerfile",
		Remove: true,
	})

	if err != nil {
		fmt.Printf("error building image: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	// Print build output
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Printf("Error copying build output: %s\n", err)
		return "", err
	}

	return imageName, nil
}

func (i* Imager) pushImage (
	fullImageName string,
) error {
	authConfig := registry.AuthConfig{
		Username: os.Getenv("DOCKER_REGISTRY_USERNAME"),
		Password: os.Getenv("DOCKER_REGISTRY_TOKEN"),
	}

	authConfigBytes, _ := json.Marshal(authConfig)
	authStr := base64.URLEncoding.EncodeToString(authConfigBytes)

	pushResp, err := i.dockerClient.ImagePush(
		context.Background(),
		fullImageName,
		image.PushOptions {
			RegistryAuth: authStr,

		},
	)

	if err != nil {
		return fmt.Errorf("error pushing image: %s", err)
	}

	defer pushResp.Close()

	_, err = io.Copy(os.Stdout, pushResp)
	if err != nil {
		return fmt.Errorf("failed to read push output: %v", err)
	}

	return nil
}

func(i* Imager) CreateAndPushImage(
	courseName string,
	contents string,
) (string, error) {
	fullImageName, err := i.createImage(courseName, contents)
	if err != nil {
		fmt.Printf("Error creating image: %s\n", err)
		return "", err
	}

	err = i.pushImage(fullImageName)

	if err != nil {
		fmt.Printf("Error pushing image: %s\n", err)
		return "", err
	}

	fmt.Printf("Image pushed successfully\n")
	return fullImageName, nil
}

