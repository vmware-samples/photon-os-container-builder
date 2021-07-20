// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.
package network

import (
	"crypto/rand"
	"net"

	"github.com/vishvananda/netlink"
)

func AcquireLinks(n *Network) error {
	linkList, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, link := range linkList {
		if link.Attrs().Name == "lo" {
			continue
		}

		n.LinksByName[link.Attrs().Name] = link.Attrs().Index
		n.LinksByIndex[link.Attrs().Index] = link.Attrs().Name

	}

	return nil
}

func SetLinkOperStateUp(dev string) error {
	link, err := netlink.LinkByName(dev)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	return nil
}

func RandomMAC() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	buf[0] &= 0xfc
	return buf, nil
}
