package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s", "ID", "Image", "Name", "State", "Status")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s", "----", "----", "----", "----", "----")

	for _, container := range containers {
		fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s",
			container.ID[:10],
			container.Image,
			strings.TrimPrefix(container.Names[0], "/"),
			container.State,
			container.Status)
	}
	w.Flush()
	fmt.Printf("\n\n")

	// Images

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	imageHeader, _ := ascii.Render("Images")
	fmt.Print(imageHeader)

	fmt.Fprintf(w, "\n%s\t%s", "ID", "RepoTags")
	fmt.Fprintf(w, "\n%s\t%s", "----", "----")

	for _, image := range images {
		fmt.Fprintf(w, "\n%s\t%s", strings.TrimPrefix(image.ID, "sha256:")[:5], image.RepoTags[0])
	}
	w.Flush()
	fmt.Printf("\n\n")
	// Volumes

	volumes, err := cli.VolumeList(context.Background(), filters.Args{})
	if err != nil {
		panic(err)
	}

	volumeHeader, _ := ascii.Render("Volumes")
	fmt.Print(volumeHeader)

	fmt.Fprintf(w, "\n%s\t%s\t%s", "Driver", "Name", "refCount")
	fmt.Fprintf(w, "\n%s\t%s\t%s", "----", "----", "----")
	for _, volume := range volumes.Volumes {
		var refCount int64
		if volume.UsageData != nil {
			refCount = volume.UsageData.RefCount
		}

		fmt.Fprintf(w, "%s\t%s\t%v\n", volume.Driver, volume.Name, refCount)
	}
	w.Flush()
	fmt.Printf("\n\n")

	// Networks

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}
	networkHeader, _ := ascii.Render("Networks")
	fmt.Print(networkHeader)
	fmt.Fprintf(w, "\n%s\t%s\t%s", "ID", "Name", "Driver")
	fmt.Fprintf(w, "\n%s\t%s\t%s", "----", "----", "----")
	for _, network := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", strings.TrimPrefix(network.ID, "sha256:")[:8], network.Name, network.Driver)
	}
	w.Flush()
	fmt.Printf("\n\n")
}
