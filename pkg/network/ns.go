// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package network

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/c-robinson/iplib"
	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func joinOriginNS(origins netns.NsHandle) {
	if err := netns.Set(origins); err != nil {
		log.Errorf("Failed to return back to origin machine network namespace: +%v", err)
		return
	}
}

func ConfigureNSNetwork(n *Network, c *conf.Config, machine string, pid int) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origins, _ := netns.Get()
	defer origins.Close()

	log.Debugf("Switching to Machine=%s GroupLeader=%d", machine, pid)

	machineNs, err := netns.GetFromPid(pid)
	if err != nil {
		log.Errorf("Could not determine network namespace for (Machine='%s') (GroupLeader=%d): %+v", machine, pid, err)
		return err
	}

	if err = netns.Set(machineNs); err != nil {
		log.Errorf("Failed to connect machine network (namespace='%s'): +%v", err)
		return err
	}
	defer joinOriginNS(origins)

	var l string
	if c.Network.Kind == "macvlan" {
		l = fmt.Sprintf("mv-%s", "photon")
	} else {
		l = "host0"
	}

	link, err := netlink.LinkByName(l)
	if err != nil {
		log.Debugf("Failed to acquire link='%s'", err)
		return err
	}

	log.Debugf("Acquired link='%s' MAC='%s'", link.Attrs().Name, link.Attrs().HardwareAddr.String())

	hwAddr, _ := RandomMAC()
	netlink.LinkSetHardwareAddr(link, hwAddr)

	if err := netlink.LinkSetUp(link); err != nil {
		log.Errorf("Failed to bring link='%s' up in network ns machine='%s': %+v", link.Attrs().Name, machine, err)
	}

	a := iplib.NextIP(n.Pool.IP)
	n.Pool.IP = a
	size, _ := a.To4().DefaultMask().Size()

	prefix := strconv.Itoa(size)

	addr := a.String() + "/" + prefix

	log.Debugf("Setting Address='%s' on link='%s'", addr, link.Attrs().Name)

	if err := AddAddress(link.Attrs().Index, addr); err != nil {
		log.Errorf("Failed to set address='%s' link='%s' in network ns machine='%s': %+v", addr, link.Attrs().Name, machine, err)
		return err
	}

	log.Debugf("Setting GW='%s' on link='%s'", n.GW, link.Attrs().Name)

	r := Route{
		IfIndex: link.Attrs().Index,
		Gw:      n.GW,
	}

	if err := r.addRoute(); err != nil {
		log.Errorf("Failed to set GW='%s' link='%s' in network ns machine='%s': %+v", n.GW, link.Attrs().Name, machine, err)
		return err
	}

	return nil
}
