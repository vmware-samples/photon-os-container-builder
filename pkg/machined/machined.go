// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package machined

import (
	"github.com/fsnotify/fsnotify"
	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/photon-os-container-builder/pkg/network"
	"github.com/photon-os-container-builder/pkg/parser"
	"github.com/photon-os-container-builder/pkg/system"
)

const (
	MachinesStateDir = "/run/systemd/machines"
)

func ConfigureOneMachine(n *network.Network, c *conf.Config, machine string) {
	leader, err := parser.ParseGroupLeader(machine)
	if err != nil {
		return
	}

	if err := network.ConfigureNSNetwork(n, c, machine, leader); err != nil {
		log.Errorf("Failed to configure network in NS for machine '%s': %+v", machine, err)
	}
}

func TaskMachine(n *network.Network, c *conf.Config) error {
	machines, err := system.FilePathWalkDir(MachinesStateDir)
	if err != nil {
		log.Errorf("Failed to fetch machines: %+v", err)
		return err
	}

	if len(n.Machines.M) <= 0 {
		for machine := range machines {
			ConfigureOneMachine(n, c, machine)
		}
		n.Machines.M = machines
		return nil
	}

	for machine := range machines {
		_, ok := n.Machines.M[machine]
		if ok {
			continue
		}

		ConfigureOneMachine(n, c, machine)
	}

	n.Machines.M = machines
	return nil
}

func TaskMachineRemove(n *network.Network) error {
	machines, err := system.ParseMachines(MachinesStateDir)
	if err != nil {
		log.Errorf("Failed to fetch machines: %+v", err)
		return err
	}

	n.Machines.M = machines

	return nil
}

func Watch(n *network.Network, c *conf.Config, done chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("Failed to watch Machines : %+v", err)
	}
	defer watcher.Close()

	log.Infoln("Watching Machines ...")

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debugf("Received event: %s", event.Op.String())

				switch event.Op.String() {
				case "CREATE":
					TaskMachine(n, c)
				case "REMOVE":
					TaskMachineRemove(n)
				}

			case err := <-watcher.Errors:
				log.Errorf("Failed to watch Machines: %+v", err)
			}
		}
	}()

	if err := watcher.Add(MachinesStateDir); err != nil {
		log.Errorf("Failed to watch Machines: %+v", err)
	}

	<-done
}
