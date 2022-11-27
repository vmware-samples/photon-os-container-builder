// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package systemd

import (
	"context"
	"errors"
	"strings"
	"time"

	sd "github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"

	"github.com/photon-os-container-builder/pkg/system"
)

const (
	defaultRequestTimeout = 10 * time.Second
)

type Unit struct {
	Command string
	Name    string
}

func (u *Unit) appendSuffixIfMissing() {
	ok := strings.HasSuffix(u.Name, ".service")
	if !ok {
		u.Name += ".service"
	}
}

func SetupContainerService(container string, network string, link string, machine string, ephemeral bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	c, err := sd.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Errorf("Failed to establish connection with system bus: %v", err)
		return err
	}
	defer c.Close()

	if err := system.CreateUnitFile(container, network, link, machine, ephemeral); err != nil {
		return err
	}

	if err := system.CreateNetworkdUnitFile(container, network, link); err != nil {
		return err
	}

	if err := c.ReloadContext(ctx); err != nil {
		return err
	}

	return nil
}

func (u *Unit) ApplyCommand() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	c, err := sd.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Errorf("Failed to establish connection with system bus: %v", err)
		return err
	}
	defer c.Close()

	u.appendSuffixIfMissing()

	ch := make(chan string)
	switch u.Command {
	case "start":
		jid, err := c.StartUnitContext(ctx, u.Name, "replace", ch)
		if err != nil {
			log.Errorf("Failed to start systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'start' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "stop":
		jid, err := c.StopUnitContext(ctx, u.Name, "fail", ch)
		if err != nil {
			log.Errorf("Failed to stop systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'stop' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "restart":
		jid, err := c.RestartUnitContext(ctx, u.Name, "replace", ch)
		if err != nil {
			log.Errorf("Failed to restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "try-restart":
		jid, err := c.TryRestartUnitContext(ctx, u.Name, "replace", ch)
		if err != nil {
			log.Errorf("Failed to try restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'try-restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "reload-or-restart":
		jid, err := c.ReloadOrRestartUnitContext(ctx, u.Name, "replace", ch)
		if err != nil {
			log.Errorf("Failed to reload or restart systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'reload-or-restart' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "reload":
		jid, err := c.ReloadUnitContext(ctx, u.Name, "replace", ch)
		if err != nil {
			log.Errorf("Failed to reload systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'reload' on systemd unit='%s' job_id='%d'", u.Name, jid)

	case "enable":
		install, changes, err := c.EnableUnitFilesContext(ctx, []string{u.Name}, false, true)
		if err != nil {
			log.Errorf("Failed to enable systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'enable' on systemd unit='%s' install='%t' changes='%s'", u.Name, install, changes)

	case "disable":
		changes, err := c.DisableUnitFilesContext(ctx, []string{u.Name}, false)
		if err != nil {
			log.Errorf("Failed to disable systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'disable' on systemd unit='%s' changes='%s'", u.Name, changes)

	case "mask":
		changes, err := c.MaskUnitFilesContext(ctx, []string{u.Name}, false, true)
		if err != nil {
			log.Errorf("Failed to mask systemd unit='%s': %v", u.Name, err)
			return err
		}

		log.Debugf("Successfully executed 'mask' on systemd unit='%s' changes='%s'", u.Name, changes)

	case "unmask":
		changes, err := c.UnmaskUnitFilesContext(ctx, []string{u.Name}, false)
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
