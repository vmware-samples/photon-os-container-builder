// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/photon-os-container-builder/pkg/conf"
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/photon-os-container-builder/pkg/machined"
	"github.com/photon-os-container-builder/pkg/network"
)

func main() {
	log.Init()

	c, err := conf.Parse()
	if err != nil {
		log.Warnf("Failed to parse configuration: %v", err)
	}

	n := network.New(c)
	if n == nil {
		log.Fatalln("Failed to create network. Aborting ...")
		os.Exit(1)
	}

	if err := network.SetupNetwork(n, c); err != nil {
		log.Fatalln("Failed to setup network. Aborting ...")
		os.Exit(1)
	}

	finished := make(chan bool)
	go machined.Watch(n, c, finished)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		os.Exit(0)
	}()

	<-finished
}
