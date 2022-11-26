// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package system

import (
	"os"
	"path"
)

const (
	defaultStorageDir           = "/var/lib/machines"
	defaultUnitFilePath         = "/lib/systemd/system"
	defaultNetworkdUnitFilePath = "/lib/systemd/network"
)

func CreateUnitFile(container string, network string, link string, machine string, ephemeral bool) error {
	unit := path.Join(defaultUnitFilePath, container) + ".service"

	file, err := os.Create(unit)
	if err != nil {
		return err
	}
	defer file.Close()

	line := "[Unit]\nDescription=Photon OS container " + container + "\n" +
		"Documentation=man:containerctl(1)\n" +
		"Wants=modprobe@tun.service modprobe@loop.service modprobe@dm-mod.service\n" +
		"PartOf=machines.target\n" +
		"Before=machines.target\n" +
		"After=network.target systemd-resolved.service modprobe@tun.service modprobe@loop.service modprobe@dm-mod.service\n\n"

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	execStart := "ExecStart=/usr/bin/containerctl boot "
	if machine != "" {
		execStart += "-m " + machine + " "
	}
	if ephemeral {
		execStart += "-x "
	}
	if network != "" {
		execStart += "-n " + network + " "
	}
	if link != "" {
		execStart += "-l " + link + " "
	}

	line = "[Service]\n" +
		execStart + container + "\n" +
		"KillMode=mixed\n" +
		"Type=notify \n" +
		"RestartForceExitStatus=133 \n" +
		"SuccessExitStatus=133 \n" +
		"Slice=machine.slice \n" +
		"Delegate=yes \n" +
		"TasksMax=16384\n\n" +
		"DevicePolicy=closed \n" +
		"DeviceAllow=/dev/net/tun rwm \n" +
		"DeviceAllow=char-pts rw \n" +

		"# nspawn itself needs access to /dev/loop-control and /dev/loop, to implement \n" +
		"# the --image= option. Add these here, too. \n" +
		"DeviceAllow=/dev/loop-control rw \n" +
		"DeviceAllow=block-loop rw \n" +
		"DeviceAllow=block-blkext rw \n\n" +

		"# nspawn can set up LUKS encrypted loopback files, in which case it needs \n" +
		"# access to /dev/mapper/control and the block devices /dev/mapper/*. \n" +
		"DeviceAllow=/dev/mapper/control rw \n" +
		"DeviceAllow=block-device-mapper rw \n\n"

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	line = "[Install]\nWantedBy=machines.target\n"
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func RemoveUnitFile(container string) error {
	unit := path.Join(defaultUnitFilePath, container) + ".service"

	if err := os.Remove(unit); err != nil {
		return err
	}

	return nil
}

func CreateNetworkdUnitFile(container string, network string, link string) error {
	unit := "10-" + container + ".network"
	unit = path.Join(defaultNetworkdUnitFilePath, unit)
	unit = path.Join(container, unit)
	unit = path.Join(defaultStorageDir, unit)

	file, err := os.Create(unit)
	if err != nil {
		return err
	}
	defer file.Close()

	line := "[Match]\nName="

	if network == "ipvlan" {
		line += "iv*\n\n[Network]\nDHCP=ipv4\n\n"
		line += "[DHCPv4]\nRequestBroadcast=yes\n"
	} else {
		line += "mv*\n\n[Network]\nDHCP=ipv4\n\n"
	}

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	if err := os.Chmod(unit, 0644); err != nil {
		return err
	}

	return nil
}
