package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mbndr/figlet4go"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	ascii := figlet4go.NewAsciiRender()
	containerHeader, _ := ascii.Render("CONTAINERS")
	fmt.Print(containerHeader)
	for _, container := range containers {
		fmt.Printf("%s %s %s [%s] %s \n",
			container.ID[:10],
			container.Image,
			strings.TrimPrefix(container.Names[0], "/"),
			container.State,
			container.Status)
	}

	// Images
	imageHeader, _ := ascii.Render("Images")
	fmt.Print(imageHeader)
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("%s\t%s \n", strings.TrimPrefix(image.ID, "sha256:")[:5], image.RepoTags[0])
	}
}
