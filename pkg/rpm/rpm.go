// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package rpm

import (
	"fmt"
	"os"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/set"
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	RPMCli  = "/usr/bin/rpm"
	TDNFCli = "/usr/bin/tdnf"
)

func ConstructOSTree(release string, target string, packages set.Set) error {
	if err := initRPMDB(target); err != nil {
		return err
	}

	if err := importGPGKey(target); err != nil {
		return err
	}

	if err := installPackages(release, target, packages); err != nil {
		return err
	}

	if err := system.RecursesiveChmod(target, 0755); err != nil {
		return err
	}

	return nil
}

func installPackages(release string, target string, packages set.Set) error {
	if release == "" {
		release = "--releasever=" + conf.DefaultReleaseVersion
	} else {
		release = "--releasever=" + release
	}
	for pkg := range packages.M {
		system.ExecAndShowProgess(TDNFCli, release, "--installroot", target, "install", pkg, "-y")
	}

	return nil
}

func importGPGKey(target string) error {
	keys, err := system.FilePathWalkDir(conf.DefaultGPGDir)
	if err != nil {
		fmt.Printf("Failed to find any RPM GPG keys in '%s': %+v\n", conf.DefaultGPGDir, err)
		return err
	}

	for key := range keys {
		if err := system.ExecAndDisplay(os.Stdout, RPMCli, "--root", target, "--import", key); err != nil {
			fmt.Printf("Failed to import GPG key '%s'': %+v\n", key, err)
			return err
		}
	}

	return nil
}

func initRPMDB(target string) error {
	if err := system.ExecAndDisplay(os.Stdout, RPMCli, "--root", target, "--initdb"); err != nil {
		fmt.Printf("Failed to init RPM db': %+v\n", err)
		return err
	}

	return nil
}
