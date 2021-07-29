// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package systemd

import (
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	systemCtl = "/usr/bin/systemctl"
)

func Stop(c string) error {
	if err := system.ExecAndRenounce(systemCtl, "stop", c); err != nil {
		return err
	}
	return nil
}

func Start(c string) error {
	if err := system.ExecAndRenounce(systemCtl, "start", c); err != nil {
		return err
	}

	return nil
}

func Restart(c string) error {
	if err := system.ExecAndRenounce(systemCtl, "restart", c); err != nil {
		return err
	}

	return nil
}

func SetupContainerService(container string, network string, link string, machine string, ephemeral bool) error {
	if err := system.CreateUnitFile(container, network, link, machine, ephemeral); err != nil {
		return err
	}

	if err := system.CreateNetworkdUnitFile(container, network, link); err != nil {
		return err
	}

	if err := system.ExecAndRenounce(systemCtl, "daemon-reload"); err != nil {
		return err
	}

	return nil
}
