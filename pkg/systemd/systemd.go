// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package systemd

import (
	"context"
	"errors"
	"time"

	sd "github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"

	"github.com/photon-os-container-builder/pkg/system"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

type Unit struct {
	Command string
	Name    string
}

func SetupContainerService(container string, network string, link string, machine string, ephemeral bool) error {
	if err := system.CreateUnitFile(container, network, link, machine, ephemeral); err != nil {
		return err
	}

	if err := system.CreateNetworkdUnitFile(container, network, link); err != nil {
		return err
	}

	if err := system.ExecAndRenounce("/usr/bin/systemctl", "daemon-reload"); err != nil {
		return err
	}

	return nil
}

func (u *Unit) ApplyCommand() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	conn, err := sd.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Errorf("Failed to establish connection with system bus: %v", err)
		return err
	}
	defer conn.Close()

	c := make(chan string)
	switch u.Command {
	case "start":
		jid, err := conn.StartUnitContext(ctx, u.Name, "replace", c)
		if err != nil {
			log.Errorf("Failed to start systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'start' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "stop":
		jid, err := conn.StopUnitContext(ctx, u.Name, "fail", c)
		if err != nil {
			log.Errorf("Failed to stop systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'stop' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "restart":
		jid, err := conn.RestartUnitContext(ctx, u.Name, "replace", c)
		if err != nil {
			log.Errorf("Failed to restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "try-restart":
		jid, err := conn.TryRestartUnitContext(ctx, u.Name, "replace", c)
		if err != nil {
			log.Errorf("Failed to try restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'try-restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "reload-or-restart":
		jid, err := conn.ReloadOrRestartUnitContext(ctx, u.Name, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload or restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'reload-or-restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "reload":
		jid, err := conn.ReloadUnitContext(ctx, u.Name, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'reload' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "enable":
		install, changes, err := conn.EnableUnitFilesContext(ctx, []string{u.Name}, false, true)
		if err != nil {
			log.Errorf("Failed to enable systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'enable' on systemd unit='%s' install='%t' changes='%s'", u.Name, install, changes)

	case "disable":
		changes, err := conn.DisableUnitFilesContext(ctx, []string{u.Name}, false)
		if err != nil {
			log.Errorf("Failed to disable systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'disable' on systemd unit='%s' changes='%s'", u.Name, changes)

	case "mask":
		changes, err := conn.MaskUnitFilesContext(ctx, []string{u.Name}, false, true)
		if err != nil {
			log.Errorf("Failed to mask systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'mask' on systemd unit='%s' changes='%s'", u.Name, changes)

	case "unmask":
		changes, err := conn.UnmaskUnitFilesContext(ctx, []string{u.Name}, false)
		if err != nil {
			log.Errorf("Failed to unmask systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'unmask' on systemd unit='%s' changes='%s'", u.Name, changes)

	default:
		log.Errorf("Unknown unit command='%s' for systemd unit='%s'", u.Command, u.Name)
		return errors.New("unknown unit command")
	}

	return nil
}
