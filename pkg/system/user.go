// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package system

import (
	"os/user"
	"strconv"
	"syscall"
)

func GetUserCredentials(usr string) (*syscall.Credential, error) {
	var u *user.User
	var err error

	if usr != "" {
		u, err = user.Lookup(usr)

	} else {
		u, err = user.Current()
	}
	if err != nil {
		return nil, err
	}

	i, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return nil, err
	}
	uid := uint32(i)

	i, err = strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return nil, err
	}
	gid := uint32(i)

	return &syscall.Credential{Uid: uid, Gid: gid}, nil
}