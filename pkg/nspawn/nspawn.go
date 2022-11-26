// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package nspawn

import (
	"errors"
	"fmt"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	nspawn     = "/usr/bin/systemd-nspawn"
	capability = "--capability=all"
)

func determineNetworking(network string, link string) (string, error) {
	var netDev string

	if network == "macvlan" {
		netDev = "--network-macvlan=" + link
	} else if network == "ipvlan" {
		netDev = "--network-ipvlan=" + link
	} else {
		return "", errors.New("unsupported networking")
	}

	return netDev, nil
}

func Spawn(c string, dir bool) (err error) {
	if dir {
		if err = system.ExecAndRenounce(nspawn, "-D", c); err != nil {
			fmt.Printf("Failed to execute systemd-nspawn: %+v\n", err)
			return err
		}
	}

	return nil
}

func ThunderBolt(c *conf.Config, container string, network string, link string, machine string, ephemeral bool) (err error) {
	var netDev string

	if network != "" {
		netDev, err = determineNetworking(network, link)
		if err != nil {
			return err
		}
	}

	if network != "" {
		if ephemeral {
			if machine != "" {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", machine, netDev)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, netDev)
			}
		} else {
			if machine != "" {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", machine, netDev)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, netDev)
			}
		}
	} else {
		if ephemeral {
			if machine != "" {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", machine)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container)
			}
		} else {
			if machine != "" {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", machine)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container)
			}
		}

	}
	if err != nil {
		fmt.Printf("Failed to start existing container '%s': %+v\n", container, err)
		return err
	}

	return nil
}

func Boot(c *conf.Config, container string, network string, link string, machine string, ephemeral bool) (err error) {
	var netDev string

	if network != "" {
		netDev, err = determineNetworking(network, link)
		if err != nil {
			return err
		}
	}

	fmt.Println(network, link, machine)
	if network != "" {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, netDev, "--link-journal=try-guest", "-M", machine)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, netDev, "--link-journal=try-guest", "-M", machine)
		}
	} else {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, "-M", machine)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, "--link-journal=try-guest", "-M", machine)
		}
	}

	if err != nil {
		fmt.Printf("Failed to boot container '%s': %+v\n", container, err)
		return err
	}

	return nil
}
