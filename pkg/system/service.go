// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package system

import (
	"os"
	"path"
)

const (
	systemctlCli        = "/usr/bin/systemctl"
	defaultUnitFilePath = "/usr/lib/systemd/system"
)

func CreateUnitFile(container string) error {
	unit := path.Join(defaultUnitFilePath, container) + ".service"

	file, err := os.Create(unit)
	if err != nil {
		return err
	}
	defer file.Close()

	line := "[Unit]\nDescription=" + container + "\n\n"
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	line = "[Service]\nLimitNOFILE=100000\nExecStart=/usr/bin/containerctl boot " + container + "\nRestart=always\n\n"
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	line = "[Install]\nAlso=dbus.service\n"
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func SetupContainerService(container string) error {
	if err := CreateUnitFile(container); err != nil {
		return err
	}

	if err := ExecAndDisplay(os.Stdout, systemctlCli, "daemon-reload"); err != nil {
		return err
	}

	return nil
}
