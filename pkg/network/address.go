// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package network

import (
	"github.com/vishvananda/netlink"
)

func AddAddress(ifIndex int, address string) error {
	link, err := netlink.LinkByIndex(ifIndex)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	if err := netlink.AddrAdd(link, addr); err != nil && err.Error() != "file exists" {
		return err
	}

	return nil
}
