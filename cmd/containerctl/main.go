// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/container"
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
					Usage:   "create systemd service unit with ephemeral flag",
				},
				&cli.BoolFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "Once installation is finished, chroot into the container",
				},
				&cli.BoolFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Enable systemd-networkd inside container",
				},
				&cli.StringFlag{
					Name:    "machine",
					Aliases: []string{"m"},
					Usage:   "Sets the machine name for this container",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				release := c.String("release")
				machine := c.String("machine")
				dir := c.Bool("dir")
				network := c.Bool("network")
				ephemeral := c.Bool("ephemeral")

				if c.String("packages") == "" {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, conf.DefaultPackages, machine, dir, network, ephemeral); err != nil {
						os.Exit(1)
					}
				} else {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, c.String("packages"), machine, dir, network, ephemeral); err != nil {
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
				&cli.StringFlag{
					Name:    "machine",
					Aliases: []string{"m"},
					Usage:   "Sets the machine name for this container",
				},
			},

			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				link := c.String("link")
				machine := c.String("machine")

				if err := container.Boot(cfg, conf.DefaultStorageDir, c.Args().First(), link, machine, c.Bool("ephemeral"), c.Bool("network")); err != nil {
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
				&cli.StringFlag{
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

				if err := container.JumpStart(cfg, conf.DefaultStorageDir, c.Args().First(), c.String("link"), c.String("machine"), c.Bool("ephemeral"), c.Bool("network")); err != nil {
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
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
