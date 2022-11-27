// SPDX-License-Identifier: BSD-2
// Copyright 2022 VMware, Inc.

package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel       = "info"
	DefaultReleaseVersion = "4.0"
	DefaultStorageDir     = "/var/lib/machines"
	DefaultUnitFilePath   = "/etc/systemd/system"
	DefaultGPGDir         = "/etc/pki/rpm-gpg"
	DefaultPackages       = "systemd,dbus,iproute2,tdnf,photon-release,photon-repos,curl,shadow,ncurses-terminfo,iputils,glibc,zlib," +
		"filesystem,pkg-config,bash,bzip2,procps-ng,iana-etc,coreutils,bc,libtool,net-tools,findutils,xz,util-linux," +
		"ca-certificates,Linux-PAM,file,e2fsprogs,rpm,openssh,gdbm,python3,python3-libs,python3-xml,sed,grep,cpio,gzip," +
		"vim,open-vm-tools,cloud-init,krb5,which,tzdata,"

	Version  = "0.1"
	ConfPath = "/etc/photon-os-container/"
	ConfFile = "photon-os-container"
)

type System struct {
	Packages string `mapstructure:"Packages"`
	Release  string `mapstructure:"Release"`
}
type Config struct {
	System System `mapstructure:"System"`
}

func Parse() (*Config, error) {
	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("Failed to read configuration: %+v", err)
	}

	viper.SetDefault("System.Release", DefaultReleaseVersion)
	viper.SetDefault("System.Packages", DefaultPackages)

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println("Failed to parse config file: '/etc/photon-os-container/photon-os-container.toml'")
	}

	return &c, nil
}
