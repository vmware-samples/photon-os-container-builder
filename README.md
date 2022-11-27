### photon-os-container-builder
----
***cntrctl*** is CLI tool which spawns Photon OS in a light-weight container. It uses `systemd-nspawn` to start Photon OS containers. The primary
use case for ***cntrctl*** is to run Photon OS test cases in a isolated environment.

Photon OS package manager ***tdnf*** integrated with ***cntrctl***. Hence it allows to prepare a root fs consisting packages depending on the user choice. It automatically prepares the root fs and boots into the container quickly. VMDK images can be automatically deployed via ***cntrctl*** and tested.

```bash
> ❯ sudo cntrctl spawn photon4
Refreshing metadata for: 'VMware Photon Linux 4.0 (x86_64)'
Refreshing metadata for: 'VMware Photon Linux 4.0 (x86_64) Updates'
Refreshing metadata for: 'VMware Photon Extras 4.0 (x86_64)'

Installing:
filesystem               x86_64       1.1-4.ph4        photon-release   7.19k 7363

Total installed size:   7.19k 7363

Downloading:

Testing transaction
Running transaction
Installing/Updating: filesystem-1.1-4.ph4.x86_64

Complete!
```

```bash
> ❯ sudo cntrctl dir photon4
Spawning container photon4 on /var/lib/machines/photon4.
Press ^] three times within 1s to kill container.
root@photon4 [ ~ ]# passwd
New password:
BAD PASSWORD: The password is shorter than 8 characters
Retype new password:
passwd: password updated successfully
```

```bash
> ❯ sudo cntrctl boot photon4

Spawning container photon4 on /var/lib/machines/photon4.
Press ^] three times within 1s to kill container.
systemd v247.11-4.ph4 running in system mode. (+PAM -AUDIT +SELINUX +IMA -APPARMOR +SMACK +SYSVINIT +UTMP -LIBCRYPTSETUP +GCRYPT +GNUTLS +ACL +XZ +LZ4 +ZSTD +SECCOMP +BLKID +ELFUTILS +KMOD -IDN2 -IDN -PCRE2 default-hierarchy=hybrid)
Detected virtualization systemd-nspawn.
Detected architecture x86-64.

Welcome to VMware Photon OS/Linux!

Welcome to VMware Photon OS/Linux!

[  OK  ] Finished Permit User Sessions.
[  OK  ] Started Console Getty.
[  OK  ] Reached target Login Prompts.
[  OK  ] Started Network Service.
[  OK  ] Reached target Multi-User System.
         Starting Update UTMP about System Runlevel Changes...
[  OK  ] Finished Update UTMP about System Runlevel Changes.
[  OK  ] Started OpenSSH Daemon.

Welcome to Photon 4.0 (x86_64) - Kernel 5.10.142-2.ph4 (console)
photon4 login:
```

```bash
> cntrctl
NAME:
   cntrctl - Controls state of containers

USAGE:
   cntrctl [global options] command [command options] [arguments...]

VERSION:
   v0.1

DESCRIPTION:
   Compose and deploy photon OS containers

COMMANDS:
   spawn, s      [NAME] Spawn a container
   boot, b       [NAME] Boot a container
   dir, d        [NAME] Directory to use as file system root for the container
   start         [NAME] start container as a systemd service unit (use host networking)
   stop          [NAME] stop container as a systemd service unit
   restart       [NAME] restart container as a systemd service unit

```

#### Build

```bash
❯  make build
❯  sudo make install
❯  sudo tdnf install systemd-container
```

cntrctl spawn [command options] [arguments...]

OPTIONS:

   `--packages value, -p`
      If specified, the list of packages will be used to compose the container.

   `--release value, -r`
      If specified, the Photon OS release version will be used. Defaults to 4.0.

   `--ephemeral, -x`
      If specified, a systemd service unit will be created with ephemeral flag.

   `--dir, -d`
      If specified, Once installation is finished, chroot into the container,

   `--network value, -n`
       If specified, enables kind of network (macvlan, ipvlan) and also enable systemd-networkd inside container

   `--link value, -l`
      If specified, the parent physical interface that is to be associated with a MACVLAN/IPVLAN to container. This
      should be used with combination of `--network` option.

   `--machine value, -m`
       If specified, sets the machine name for this container during runtime.


#### Contributing
----

The ***photon-os-container-builder*** project team welcomes contributions from the community. If you wish to contribute code and you have not signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For any questions about the CLA process, please refer to our [FAQ](https://cla.vmware.com/faq).

slack channel [#photon](https://code.vmware.com/web/code/join).

#### License
----

[BSD 2-Clause](https://spdx.org/licenses/BSD-2-Clause.html)
