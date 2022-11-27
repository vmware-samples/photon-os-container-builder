// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package system

import (
	"os"
	"path"

	"github.com/photon-os-container-builder/pkg/keyfile"
)

func CreateUnitFile(container string, network string, link string, machine string, ephemeral bool) error {
	m, err := keyfile.Create(path.Join("/lib/systemd/system", container+".service"))
	if err != nil {
		return err
	}

	m.SetKeySectionString("Unit", "Description", "Photon OS container")
	m.SetKeySectionString("Unit", "Documentation", "man:cntrctl(1)")
	m.SetKeySectionString("Unit", "Wants", "modprobe@tun.service modprobe@loop.service modprobe@dm-mod.service")
	m.SetKeySectionString("Unit", "PartOf", "machines.target")
	m.SetKeySectionString("Unit", "Before", "machines.target")
	m.SetKeySectionString("Unit", "After", "network.target systemd-resolved.service modprobe@tun.service modprobe@loop.service modprobe@dm-mod.service")

	execStart := "/usr/bin/cntrctl boot "
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

	m.SetKeySectionString("Service", "ExecStart", execStart+container)
	m.SetKeySectionString("Service", "KillMode", "mixed")
	m.SetKeySectionString("Service", "Type", "notify")
	m.SetKeySectionString("Service", "RestartForceExitStatus", "133")
	m.SetKeySectionString("Service", "SuccessExitStatus", "133")
	m.SetKeySectionString("Service", "Slice", "machine.slice")
	m.SetKeySectionString("Service", "Delegate", "yes")
	m.SetKeySectionString("Service", "TasksMax", "16384")
	m.SetKeySectionString("Service", "DevicePolicy", "closed")
	m.SetKeySectionString("Service", "DeviceAllow", "/dev/net/tun rwm")
	m.SetKeySectionString("Service", "DeviceAllow", "char-pts rw")
	m.SetKeySectionString("Service", "DeviceAllow", "/dev/loop-control rw")
	m.SetKeySectionString("Service", "DeviceAllow", "block-loop rw")
	m.SetKeySectionString("Service", "DeviceAllow", "block-blkext rw")
	m.SetKeySectionString("Service", "DeviceAllow", "block-device-mapper rw")
	m.SetKeySectionString("Service", "DeviceAllow", "/dev/mapper/control rw")

	m.SetKeySectionString("Install", "WantedBy", "machines.target")
	return m.Save()
}

func RemoveUnitFile(container string) error {
	unit := path.Join("/lib/systemd/system", container+".service")
	if err := os.Remove(unit); err != nil {
		return err
	}

	return nil
}

func CreateNetworkUnitFile(container string, network string, link string) error {
	u := "10-" + container + ".network"
	m, err := keyfile.Create(path.Join("/var/lib/machines/"+container+"/lib/systemd/network", u))
	if err != nil {
		return err
	}

	if network == "ipvlan" {
		m.SetKeySectionString("Match", "Name", "iv*")
		m.SetKeySectionString("DHCP4", "RequestBroadcast", "yes")
	} else {
		m.SetKeySectionString("Match", "Name", "mv*")
	}

	m.SetKeySectionString("Network", "DHCP", "yes*")
	m.Save()

	usr, err := GetUserCredentials("systemd-network")
	if err != nil {
		return err
	}

	return os.Chown(m.Path, int(usr.Uid), int(usr.Gid))
}
