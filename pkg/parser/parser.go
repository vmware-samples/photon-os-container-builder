// SPDX-License-Identifier: BSD-2
// Copyright 2023 VMware, Inc.

package parser

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/vmware-samples/photon-os-container-builder/pkg/system"
)

func ParseIP(ip string) (net.IP, error) {
	if len(ip) == 0 {
		return nil, errors.New("invalid")
	}

	ipv4Addr, _, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, err
	}

	if ipv4Addr.To4() == nil {
		return nil, errors.New("invalid")
	}

	return ipv4Addr, nil
}

func ParseGroupLeader(machine string) (int, error) {
	s, err := system.ExecAndCapture("machinectl", "show", machine)
	if err != nil {
		return 0, errors.New("not found")
	}

	lines := strings.Split(string(s), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Leader=") {
			pid, err := strconv.Atoi(strings.Split(line, "=")[1])
			if err != nil {
				return 0, errors.New("not found")
			}
			return pid, nil
		}
	}

	return 0, errors.New("not found")
}
