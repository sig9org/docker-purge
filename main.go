package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type ImageInfo struct {
	id string
	repository string
	size int64
	tag string
}

func (m ImageInfo) ID() string {	
	return strings.Split(m.id, ":")[1]
}

func (m ImageInfo) Repository() string {
	return strings.Split(m.repository, "@")[0]
}

func (m ImageInfo) Size() string {
	//return strconv.FormatInt(m.size / 1024 / 1024, 10) + "MB"
	return strconv.FormatInt(m.size / 1000 / 1000, 10) + "MB"
}

func (m ImageInfo) Tag() string {
	return strings.Split(m.tag, ":")[1]
}

func main() {
	pluginMain()
}

func pluginMain() {
	var (
		rmi, preRun bool
	)
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command {
			Use:   "purge",
			Short: "Stop and delete all containers",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				if err := plugin.PersistentPreRunE(cmd, args); err != nil {
					return err
				}
				if preRun {
					fmt.Fprintf(dockerCli.Err(), "Plugin PersistentPreRunE called")
				}
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := context.Background()
				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					panic(err)
				}
				cli.NegotiateAPIVersion(ctx)

				containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
				if err != nil {
					panic(err)
				}

				col2 := len("IMAGE")
				col3 := len("STATUS")
				col4 := len("NAMES")
				for _, container := range containers {
					col2 = getLonger(col2, container.Image)
					col3 = getLonger(col3, container.Status)
					col4 = getLonger(col4, container.Names[0][1:])
				}

				fmt.Fprintf(dockerCli.Out(), "%12s  %-" + strconv.Itoa(col2) + "s  %-" + strconv.Itoa(col3) + "s  %-" + strconv.Itoa(col4) + "s  %s\n", "CONTAINER ID", "IMAGE", "STATUS", "NAMES", "PURGED")
				for _, container := range containers {
					fmt.Fprintf(dockerCli.Out(), "%12s  %-" + strconv.Itoa(col2) + "s  %-" + strconv.Itoa(col3) +"s  %-" + strconv.Itoa(col4) + "s  ", container.ID[:12], container.Image, container.Status, container.Names[0][1:])
					timeout := 0 * time.Millisecond
					if err := cli.ContainerStop(ctx, container.ID, &timeout); err != nil {
						panic(err)
					}
					if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
						panic(err)
					}
					fmt.Fprintf(dockerCli.Out(), "%s", "Done\n")
				}

				if rmi {
					images, err := cli.ImageList(ctx, types.ImageListOptions{})
					if err != nil {
						panic(err)
					}

					icol1 := len("REPOSITORY")
					icol2 := len("TAG")
					icol4 := len("SIZE")
					for _, image := range images {
						var img ImageInfo
						img.repository = image.RepoDigests[0]
						img.tag = image.RepoTags[0]
						img.id = image.ID
						img.size = image.Size
						icol1 = getLonger(icol1, img.Repository())
						icol2 = getLonger(icol2, img.Tag())
						icol4 = getLonger(icol4, img.Size())
					}

					fmt.Fprintf(dockerCli.Out(), "\n%-" + strconv.Itoa(icol1) + "s  %-" + strconv.Itoa(icol2) + "s  %-12s  %-" + strconv.Itoa(icol4) + "s  %s\n", "REPOSITORY", "TAG", "IMAGE ID", "SIZE", "PURGED")
					for _, image := range images {
						var img ImageInfo
						img.repository = image.RepoDigests[0]
						img.tag = image.RepoTags[0]
						img.id = image.ID
						img.size = image.Size
						fmt.Fprintf(dockerCli.Out(), "%-" + strconv.Itoa(icol1) + "s  %-" + strconv.Itoa(icol2) + "s  %-12s  %-" + strconv.Itoa(icol4) + "s  ", img.Repository(), img.Tag(), img.ID()[:12], img.Size())

						if _, err := cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force: true}); err != nil {
							panic(err)
						}
						fmt.Fprintf(dockerCli.Out(), "%s", "Done\n")
					}
				}

				return dockerCli.ConfigFile().Save()
			},
		}
		flags := cmd.Flags()
		flags.BoolVarP(&rmi, "images", "i", false, "Delete all Docker container images.")
		return cmd
	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Docker Inc.",
			Version:       "0.0.1",
		})
}

func getLonger(index int, value string) (int) {
	if index < len(value) {
		return len(value)
	} else {
		return index
	}
}

