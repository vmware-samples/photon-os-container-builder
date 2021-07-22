// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/container"
	"github.com/photon-os-container-builder/pkg/machinectl"
	"github.com/photon-os-container-builder/pkg/systemd"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s\n", c.App.Version)
	}

	cfg, _ := conf.Parse()

	app := &cli.App{
		Name:        "containerctl",
		Version:     "v0.1",
		Description: "Compose and deploy photon OS containers",
		Usage:       "Controls state of containers",
		HideHelp:    false,
		HideVersion: false,
	}

	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true

	app.Commands = []*cli.Command{
		{
			Name:    "spawn",
			Aliases: []string{"s"},
			Usage:   "[NAME] Spawn a container",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "packages",
					Aliases: []string{"p"},
					Usage:   "List of packages to install separated by ','",
				},
				&cli.StringFlag{
					Name:    "release",
					Aliases: []string{"r"},
					Usage:   "Photon OS release version",
				},
				&cli.BoolFlag{
					Name:    "ephemeral",
					Aliases: []string{"x"},
					Usage:   "Run with a temporary snapshot of its file system that is removed immediately when the container terminates",
				},
				&cli.BoolFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "Once installation is finished, chroot into the container",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				release := c.String("release")
				dir := c.Bool("dir")

				if c.String("packages") == "" {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, conf.DefaultPackages, dir); err != nil {
						os.Exit(1)
					}
				} else {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, c.String("packages"), dir); err != nil {
						os.Exit(1)
					}
				}

				return nil
			},
		},
		{
			Name:    "boot",
			Aliases: []string{"b"},
			Usage:   "[NAME] Boot a container",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "ephemeral",
					Aliases: []string{"x"},
					Usage:   "Run with a temporary snapshot of its file system that is removed immediately when the container terminates",
				},
				&cli.BoolFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Disconnect networking of the container from the host (private network)",
				},
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "Assign the specified network interface to the container",
				},
			},

			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := container.Boot(cfg, conf.DefaultStorageDir, c.Args().First(), c.String("link"), c.Bool("ephemeral"), c.Bool("network")); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:    "dir",
			Aliases: []string{"d"},
			Usage:   "[NAME] Directory to use as file system root for the container",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "ephemeral",
					Aliases: []string{"x"},
					Usage:   "Run with a temporary snapshot of its file system that is removed immediately when the container terminates",
				},
				&cli.BoolFlag{
					Name:    "machine",
					Aliases: []string{"m"},
					Usage:   "Sets the machine name for this container. This name may be used to identify this container during its runtime",
				},
				&cli.BoolFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Disconnect networking of the container from the host (private network)",
				},
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "Assign the specified network interface to the container",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := container.JumpStart(cfg, conf.DefaultStorageDir, c.Args().First(), c.String("link"), c.Bool("ephemeral"), c.Bool("machine"), c.Bool("network")); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "[NAME] start container as a systemd service unit (use host networking)",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := systemd.Start(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "[NAME] stop container as a systemd service unit",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := systemd.Stop(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "restart",
			Usage: "[NAME] restart container as a systemd service unit",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := systemd.Stop(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "login",
			Usage: "[NAME] Login prompt in a container",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.Login(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "shell",
			Usage: "[USER@NAME]",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.Shell(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "copy-to",
			Usage: "[NAME] [PATH] PATH Copy files from the host to a container",
			Action: func(c *cli.Context) error {
				if c.NArg() < 3 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.CopyFrom(c.Args().First(), c.Args().Get(2), c.Args().Get(3)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "copy-from",
			Usage: "[NAME] PATH [PATH] Copy files from a container to the host",
			Action: func(c *cli.Context) error {
				if c.NArg() < 3 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				err := machinectl.CopyFrom(c.Args().First(), c.Args().Get(2), c.Args().Get(3))
				if err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List running VMs and containers",
			Action: func(c *cli.Context) error {
				err := machinectl.ListRunningContainers()
				if err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "list-images",
			Usage: "List container images",
			Action: func(c *cli.Context) error {

				if err := machinectl.ListImages(); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "image-status",
			Usage: "[NAME] Display image status",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ListImageStatus(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "[NAME] Remove an image",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.RemoveImage(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "clone",
			Usage: "[NAME] Clone an image",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.CloneImage(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "rename",
			Usage: "[NAME] Rename an image",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.RenameImage(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "clean",
			Usage: "Remove all images",
			Action: func(c *cli.Context) error {
				if err := machinectl.CleanImage(); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "poweroff",
			Usage: "[NAME] poweroff a container",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.RemoveImage(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "terminate",
			Usage: "[NAME] Terminate a container",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.Terminate(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "reboot",
			Usage: "[NAME] Reboot a container",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.Reboot(c.Args().First()); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "pull-tar",
			Usage: "URL [NAME] Download a TAR container image",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.PullTar(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "pull-raw",
			Usage: "URL [NAME] Download a RAW container or VM image",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.PullRaw(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "import-tar",
			Usage: "FILE [NAME] Import a local TAR container image",
			Action: func(c *cli.Context) error {
				if c.NArg() <= 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ImportTar(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "import-raw",
			Usage: "FILE [NAME] Import a local RAW container or VM image",
			Action: func(c *cli.Context) error {
				if c.NArg() <= 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ImportRaw(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "import-fs",
			Usage: "DIRECTORY [NAME] Import a local directory container image",
			Action: func(c *cli.Context) error {
				if c.NArg() <= 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ImportFS(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "export-tar",
			Usage: "NAME [FILE] Export a TAR container image locally",
			Action: func(c *cli.Context) error {
				if c.NArg() <= 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ExportTar(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "export-raw",
			Usage: "NAME [FILE] Export a RAW container or VM image locally",
			Action: func(c *cli.Context) error {
				if c.NArg() <= 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				if err := machinectl.ExportRaw(c.Args().First(), c.Args().Get(2)); err != nil {
					os.Exit(1)
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
