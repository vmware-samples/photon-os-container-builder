// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vmware-samples/photon-os-container-builder/pkg/conf"
	"github.com/vmware-samples/photon-os-container-builder/pkg/container"
	"github.com/vmware-samples/photon-os-container-builder/pkg/systemd"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s\n", c.App.Version)
	}

	cfg, _ := conf.Parse()

	app := &cli.App{
		Name:        "cntrctl",
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
				&cli.StringFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Enable kind of network (macvlan, ipvlan) and also enable systemd-networkd inside container",
				},
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "Specifies the parent physical interface that is to be associated with a MACVLAN/IPVLAN to container",
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
				network := c.String("network")
				link := c.String("link")

				dir := c.Bool("dir")
				ephemeral := c.Bool("ephemeral")

				if network != "" && link == "" {
					fmt.Printf("network='%s' is specified but link is missing", network)
					os.Exit(1)
				} else if link != "" && network == "" {
					fmt.Printf("link='%s' is specified but network is missing", link)
					os.Exit(1)
				}

				if c.String("packages") == "" {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, conf.DefaultPackages, network, link, machine, dir, ephemeral); err != nil {
						os.Exit(1)
					}
				} else {
					if err := container.Spawn(conf.DefaultStorageDir, c.Args().First(), release, c.String("packages"), network, link, machine, dir, ephemeral); err != nil {
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
				&cli.StringFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Enable kind of network (MACVLAN/IPVLAN ) and also enable systemd-networkd inside container",
				},
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "Specifies the parent physical interface that is to be associated with a MACVLAN/IPVLAN to container",
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

				network := c.String("network")
				link := c.String("link")
				machine := c.String("machine")

				if network != "" && link == "" {
					fmt.Printf("network='%s' is specified but link is missing", network)
					os.Exit(1)
				} else if link != "" && network == "" {
					fmt.Printf("link='%s' is specified but network is missing", link)
					os.Exit(1)
				}

				if err := container.Boot(cfg, conf.DefaultStorageDir, c.Args().First(), network, link, machine, c.Bool("ephemeral")); err != nil {
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
				&cli.StringFlag{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Enable kind of network (MACVLAN/IPVLAN ) and also enable systemd-networkd inside container",
				},
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "Specifies the parent physical interface that is to be associated with a MACVLAN/IPVLAN to container",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 1)
				}

				network := c.String("network")
				link := c.String("link")
				machine := c.String("machine")

				if network != "" && link == "" {
					fmt.Printf("network='%s' is specified but link is missing", network)
					os.Exit(1)
				} else if link != "" && network == "" {
					fmt.Printf("link='%s' is specified but network is missing", link)
					os.Exit(1)
				}

				if err := container.JumpStart(cfg, conf.DefaultStorageDir, c.Args().First(), network, link, machine, c.Bool("ephemeral")); err != nil {
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
				u := systemd.Unit{
					Name:    c.Args().First(),
					Command: "start",
				}

				if err := u.ApplyCommand(); err != nil {
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

				u := systemd.Unit{
					Name:    c.Args().First(),
					Command: "stop",
				}

				if err := u.ApplyCommand(); err != nil {
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

				u := systemd.Unit{
					Name:    c.Args().First(),
					Command: "restart",
				}

				if err := u.ApplyCommand(); err != nil {
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
