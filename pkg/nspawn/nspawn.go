// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package nspawn

import (
	"fmt"
	"os"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	nspawn = "/usr/bin/systemd-nspawn"
)

func Spawn(c string, dir bool) (err error) {
	err = system.ExecAndDisplay(os.Stdout, "/usr/bin/systemd-machine-id-setup", "--root", c)
	if err != nil {
		fmt.Printf("Failed to execute systemd-machine-id-setup in '%s': %+v\n", c, err)
		return err
	}

	if dir {
		err = system.ExecAndRenounce(nspawn, "-D", c)
		if err != nil {
			fmt.Printf("Failed to execute systemd-nspawn: %+v\n", err)
			return err
		}
	}

	return nil
}

func ThunderBolt(c *conf.Config, container string, ephemeral bool, machine bool, network bool) (err error) {
	capability := "--capability=CAP_SYS_ADMIN,CAP_NET_ADMIN,CAP_MKNOD"
	var netDev string

	if network {
		if c.Network.Kind == "macvlan" {
			netDev = "--network-macvlan=photon"
		} else {
			netDev = "--network-bridge=photon"
		}
	}

	if network {
		if ephemeral {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", container, netDev)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, netDev)
			}
		} else {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", container, netDev)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, netDev)
			}
		}
	} else {
		if ephemeral {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container, "-M", container)
			} else {
				err = system.ExecAndRenounce(nspawn, capability, "-xD", container)
			}
		} else {
			if machine {
				err = system.ExecAndRenounce(nspawn, capability, "-D", container, "-M", container)
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

func Boot(c *conf.Config, container string, ephemeral bool, network bool) (err error) {
	capability := "--capability=CAP_SYS_ADMIN,CAP_NET_ADMIN,CAP_MKNOD"
	var netDev string

	if network {
		if c.Network.Kind == "macvlan" {
			netDev = "--network-macvlan=photon"
		} else {
			netDev = "--network-bridge=photon"
		}
	}

	if network {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, netDev)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, netDev)
		}
	} else {
		if ephemeral {
			err = system.ExecAndRenounce(nspawn, capability, "-xbD", container, netDev)
		} else {
			err = system.ExecAndRenounce(nspawn, capability, "-bD", container, netDev)
		}
	}
	if err != nil {
		fmt.Printf("Failed to boot container '%s': %+v\n", container, err)
	}

	return nil
}
