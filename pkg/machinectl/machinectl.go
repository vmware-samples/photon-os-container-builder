// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package machinectl

import (
	"os"

	"github.com/photon-os-container-builder/pkg/system"
)

const (
	machineCtl = "/usr/bin/machinectl"
)

func ListRunningContainers() error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "list"); err != nil {
		return err
	}

	return nil
}

func ListImages() error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "list-images"); err != nil {
		return err
	}

	return nil
}

func ListImageStatus(image string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "list-status", image); err != nil {
		return err
	}

	return nil
}

func RemoveImage(image string) error {
	err := system.ExecAndDisplay(os.Stdout, machineCtl, "remove", image)
	if err != nil {
		return err
	}

	return nil
}

func CloneImage(image string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "clone", image); err != nil {
		return err
	}

	return nil
}

func RenameImage(image string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "rename", image); err != nil {
		return err
	}

	return nil
}

func CleanImage() error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "clean"); err != nil {
		return err
	}

	return nil
}

func Poweroff(c string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "poweroff", c); err != nil {
		return err
	}

	return nil
}

func Reboot(c string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "reboot", c); err != nil {
		return err
	}

	return nil
}

func Terminate(c string) error {
	if err := system.ExecAndDisplay(os.Stdout, machineCtl, "terminate", c); err != nil {
		return err
	}

	return nil
}

func Start(c string) error {
	if err := system.ExecAndRenounce(machineCtl, "start", c); err != nil {
		return err
	}

	return nil
}

func Login(c string) error {
	if err := system.ExecAndRenounce(machineCtl, "login", c); err != nil {
		return err
	}

	return nil
}

func Shell(c string, cmd string) error {
	if err := system.ExecAndRenounce(machineCtl, "shell", c, cmd); err != nil {
		return err
	}

	return nil
}

func CopyTo(c string, to string, from string) error {
	if err := system.ExecAndRenounce(machineCtl, "copy-to", c, to, from); err != nil {
		return err
	}

	return nil
}

func CopyFrom(c string, to string, from string) error {
	if err := system.ExecAndRenounce(machineCtl, "copy-to", c, from, to); err != nil {
		return err
	}

	return nil
}

func PullTar(url string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "pull-tar", url, d); err != nil {
		return err
	}

	return nil
}

func PullRaw(url string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "pull-raw", url, d); err != nil {
		return err
	}

	return nil
}

func ImportTar(t string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "import-tar", t, d); err != nil {
		return err
	}

	return nil
}

func ImportRaw(t string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "import-raw", t, d); err != nil {
		return err
	}

	return nil
}

func ImportFS(t string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "import-fs", t, d); err != nil {
		return err
	}

	return nil
}

func ExportRaw(t string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "export-raw", t, d); err != nil {
		return err
	}

	return nil
}

func ExportTar(t string, d string) error {
	if err := system.ExecAndShowProgess(machineCtl, "export-tar", t, d); err != nil {
		return err
	}

	return nil
}
