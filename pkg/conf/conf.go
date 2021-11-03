// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package conf

import (
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/photon-os-container-builder/pkg/parser"
	"github.com/spf13/viper"
)

const (
	DefaultLogLevel       = "info"
	DefaultReleaseVersion = "4.0"
	DefaultStorageDir     = "/var/lib/machines"
	DefaultUnitFilePath   = "/etc/systemd/system"
	DefaultGPGDir         = "/etc/pki/rpm-gpg"
	DefaultPackages       = "systemd,dbus,iproute2,tdnf,photon-release,photon-repos,curl,shadow,ncurses-terminfo,iputils,glibc,zlib,filesystem,pkg-config,bash,bzip2,procps-ng,iana-etc,coreutils,bc,libtool,net-tools,findutils,xz,util-linux,kmod,linux-esx,ca-certi    ficates,iptables,Linux-PAM,file,e2fsprogs,rpm,openssh,gdbm,python3,python3-libs,python3-xml,sed,grep,cpio,gzip,vim,open-vm-tools,docker,bridge-utils,cloud-init,krb5,which,tzdata,initramfs"

	DefaultParentLink  = "eth0"
	DefaultNetworkKind = "link"
	DefaultLink        = "eth1"
	DefaultAddressPool = "172.16.85.50/24 "
	DefaultPoolOffSet  = 64

	Version  = "0.1"
	ConfPath = "/etc/photon-os-container/"
	ConfFile = "photon-os-container"
)

// Config file key value
type Network struct {
	Kind        string `mapstructure:"Kind"`
	Link        string `mapstructure:"Link"`
	ParentLink  string `mapstructure:"ParentLink"`
	AddressPool string `mapstructure:"AddressPool"`
	PoolOffset  int    `mapstructure:"PoolOffset "`
}

type System struct {
	Packages string `mapstructure:"Packages"`

	Release  string `mapstructure:"Release"`
	LogLevel string `mapstructure:"LogLevel"`
}
type Config struct {
	Network Network `mapstructure:"Network"`
	System  System  `mapstructure:"System"`
}

func Parse() (*Config, error) {
	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("%+v", err)
	}

	viper.SetDefault("Network.AddressPool", DefaultAddressPool)
	viper.SetDefault("Network.PoolOffset", DefaultAddressPool)
	viper.SetDefault("Network.Link", DefaultLink)
	viper.SetDefault("Network.ParentLink", DefaultParentLink)
	viper.SetDefault("Network.Kind", DefaultNetworkKind)

	viper.SetDefault("System.LogLevel", DefaultLogLevel)
	viper.SetDefault("System.Release", DefaultReleaseVersion)
	viper.SetDefault("System.Packages", DefaultPackages)

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		log.Warnf("Failed to parse config file: '/etc/photon-os-container/photon-os-container.toml'")
	}

	log.SetLevel(c.System.LogLevel)

	if _, err := parser.ParseIP(c.Network.AddressPool); err != nil {
		log.Debugf("Failed to parse address pool. Default will be used", err)
		c.Network.AddressPool = DefaultAddressPool
	}

	return &c, nil
}
