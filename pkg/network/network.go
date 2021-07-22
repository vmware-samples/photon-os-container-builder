// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package network

import (
	"io/ioutil"
	"strings"

	"github.com/c-robinson/iplib"
	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/photon-os-container-builder/pkg/set"
	"github.com/vishvananda/netlink"
)

const (
	ProcFSIpv4Forward = "/proc/sys/net/ipv4/ip_forward"
)

type Network struct {
	LinksByName  map[string]int
	LinksByIndex map[int]string

	GW       string
	Pool     iplib.Net
	Machines set.Set
}

func New(c *conf.Config) *Network {
	n := &Network{
		LinksByName:  make(map[string]int),
		LinksByIndex: make(map[int]string),
	}

	n.Machines = set.New()

	_, ip, _ := iplib.ParseCIDR(c.Network.AddressPool)
	n.Pool = ip
	a := iplib.IncrementIP4By(ip.IP, conf.DefaultPoolOffSet)
	n.Pool.IP = a

	return n
}

func SetupNetwork(n *Network, c *conf.Config) error {
	EnsureIPv4Forward()

	if c.Network.Kind == "macvlan" {
		if err := CreateMACVLan("photon", c.Network.ParentLink); err != nil {
			log.Fatalf("Failed to create MACVlan='photon' on link='%s': %+v", err)
			return err
		}
	} else {
		if err := CreateBridge("photon"); err != nil {
			log.Fatalf("Failed to create Bridge='photon': %+v", err)
			return err
		}
	}

	if err := SetLinkOperStateUp("photon"); err != nil {
		log.Fatalf("Failed to bring up link='photon': %+v", err)
		return err
	}

	link, err := netlink.LinkByName(c.Network.ParentLink)
	if err != nil {
		log.Fatalf("Failed to find parent link='%s': %+v", c.Network.ParentLink, err)
		return err
	}

	gw, err := getIpv4Gateway(link.Attrs().Index)
	if err != nil {
		log.Errorf("Failed to find GW for link='%s' in network ns machine='%s': %+v", c.Network.ParentLink, err)
		return err
	}

	n.GW = gw

	if err := AcquireLinks(n); err != nil {
		log.Fatalf("Failed to acquire link information. Unable to continue: %v", err)
		return err
	}

	// Watch network
	go WatchNetwork(n)

	return nil
}

func EnsureIPv4Forward() error {
	d, err := ioutil.ReadFile(ProcFSIpv4Forward)
	if err != nil {
		log.Errorf("failed to read '%s'", ProcFSIpv4Forward)
		return err
	}

	if strings.Contains(string(d), "1") {
		return nil
	}

	if err = ioutil.WriteFile(ProcFSIpv4Forward, []byte("1"), 0); err != nil {
		log.Errorf("failed to write to '%s': %v", ProcFSIpv4Forward, err)
		return err
	}

	return nil
}
