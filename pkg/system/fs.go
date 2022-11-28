// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package system

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	NetworkdService = "/etc/systemd/system/multi-user.target.wants/systemd-networkd.service"
	NetworkdDBus    = "/etc/systemd/system/dbus-org.freedesktop.network1.service"
	NetworkdOnline  = "/etc/systemd/system/network-online.target.wants/systemd-networkd-wait-online.service"
	NetworkdSocket  = "/etc/systemd/system/sockets.target.wants/systemd-networkd.socket"
)

func DisableNetworkd(c string) {
	os.Remove(path.Join(c, NetworkdService))
	os.Remove(path.Join(c, NetworkdDBus))
	os.Remove(path.Join(c, NetworkdOnline))
	os.Remove(path.Join(c, NetworkdSocket))
}

func RecursiveChmod(path string, mode os.FileMode) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		os.Chmod(path, mode)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func PathExists(path string) bool {
	_, r := os.Stat(path)
	return !os.IsNotExist(r)
}

func CreateDirectory(parent string, dir string) (err error) {
	d := path.Join(parent, dir)
	if PathExists(d) {
		return errors.New("dir exists")
	}

	return os.MkdirAll(d, 0755)
}

func RemoveDir(d string) {
	os.RemoveAll(d)
}

func FilePathWalkDir(root string) (map[string]bool, error) {
	files := make(map[string]bool)
	err := filepath.Walk(root, func(f string, info os.FileInfo, err error) error {

		if !info.IsDir() && !strings.HasPrefix(f, ".") {
			files[f] = true
		}
		return nil
	})
	return files, err
}

func ParseMachines(root string) (map[string]bool, error) {
	files := make(map[string]bool)
	err := filepath.Walk(root, func(f string, info os.FileInfo, err error) error {

		base := path.Base(f)
		if !info.IsDir() && !strings.HasSuffix(base, "scope") && !strings.HasPrefix(base, ".") {
			files[base] = true
		}
		return nil
	})
	return files, err
}

func ParseMachine(path string) (map[string]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.Contains(line, "=") {
			continue
		}

		c := strings.Split(line, "=")
		m[c[0]] = c[1]
	}

	_, ok := m["LEADER"]
	if !ok {
		return nil, errors.New("group leader not found")
	}

	_, ok = m["NAME"]
	if !ok {
		return nil, errors.New("machine name not found")
	}

	return m, nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

		lines = append(lines, scanner.Text())

	}

	return lines, scanner.Err()
}
