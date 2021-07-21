// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package nspawn

import (
	"fmt"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	nspawn = "/usr/bin/systemd-nspawn"
)

func determineNetworking(c *conf.Config, link string) (n string) {
	var netDev string

	if c.Network.Kind == "macvlan" {
		netDev = "--network-macvlan=photon"
	} else if c.Network.Kind == "bridge" {
		netDev = "--network-bridge=photon"
	} else {
		if c.Network.Link == "" {
			netDev = "--network-interface=eth1"
		} else {
			if link == "" {
				netDev = "--network-interface=" + c.Network.Link
			} else {
				netDev = "--network-interface=" + link
			}
		}
	}

	return netDev
}

func Spawn(c string, dir bool) (err error) {
	if dir {
		err = system.ExecAndRenounce(nspawn, "-D", c)
		if err != nil {
			fmt.Printf("Failed to execute systemd-nspawn: %+v\n", err)
			return err
		}
	}

	return nil
}

func ThunderBolt(c *conf.Config, container string, link string, ephemeral bool, machine bool, network bool) (err error) {
	capability := "--capability=CAP_SYS_ADMIN,CAP_NET_ADMIN,CAP_MKNOD"
	machineArg := "--machine=" + container
	var netDev string

	if network {
		netDev = determineNetworking(c, link)
	}

	if network {
		if ephemeral {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			}
		} else {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			}
		}
	} else {
		if ephemeral {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", container, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container)
			}
		} else {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", container, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "--link-journal=try-guest", "U", "--settings=override", machineArg)
			}
		}

	}
	if err != nil {
		fmt.Printf("Failed to start existing container '%s': %+v\n", container, err)
		return err
	}

	return nil
}

func Boot(c *conf.Config, container string, link string, ephemeral bool, network bool) (err error) {
	capability := "--capability=CAP_SYS_ADMIN,CAP_NET_ADMIN,CAP_MKNOD"
	machineArg := "--machine=" + container
	var netDev string

	if network {
		netDev = determineNetworking(c, link)
	}

	if network {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
		}
	} else {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, netDev, "--link-journal=try-guest", "U", "--settings=override", machineArg)
		}
	}
	if err != nil {
		fmt.Printf("Failed to boot container '%s': %+v\n", container, err)
	}

	return nil
}
