/* SPDX-License-Identifier: Apache-2.0
 * Copyright © 2021 VMware, Inc.
 */

package network

import (
	"net"

	"github.com/photon-os-container-builder/pkg/log"
	"github.com/vishvananda/netlink"
)

func CreateMACVLan(netDev string, l string) error {
	link, err := netlink.LinkByName(l)
	if err != nil {
		return err
	}

	dev := netlink.Macvlan{
		LinkAttrs: netlink.LinkAttrs{Name: netDev, ParentIndex: link.Attrs().Index},
		Mode:      netlink.MACVLAN_MODE_BRIDGE,
	}

	if err := netlink.LinkAdd(&dev); err != nil && err.Error() != "file exists" {
		return err
	}

	return nil
}

func RemoveMACFromMACVLan(netdev string, mac string) error {
	link, err := netlink.LinkByName(netdev)
	if err != nil {
		return err
	}

	hwdr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}

	if err := netlink.MacvlanMACAddrDel(link, hwdr); err != nil {
		return err
	}

	log.Debugf("Successfully removed MACAddress='%s' from '%s'", mac, netdev)

	return nil
}

func CreateBridge(netDev string) error {
	dev := netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{Name: netDev},
	}

	if err := netlink.LinkAdd(&dev); err != nil && err.Error() != "file exists" {
		return err
	}

	return nil
}

