### photon-os-container-builder
----
***cntrctl*** spawns Photon OS image as a light-weight container. It uses `systemd-nspawn` to start Photon OS containers. The primary
use case for ***cntrctl*** is to run Photon OS test cases in a isolated environment.

Photon OS package manager ***tdnf*** integrated with ***cntrctl***. Hence it allows to prepare a root fs consisting packages depending on the user choice. It automatically prepares the root fs and boots into the container quickly. VMDK images can be automatically deployed via ***cntrctl*** and tested.

Install systemd-container
```bash
❯ sudo tdnf install systemd-container -y
Installing:
systemd-container                                                x86_64                     5.10-1.ph4                          photon-release              1.65M 1728537

Total installed size:   1.65M 1728537

Downloading:
systemd-container                                  478107 100%
Testing transaction
Running transaction
Installing/Updating: systemd-container-247.11-4.ph4

Complete!
```

```bash
❯ sudo cntrctl spawn photon4
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

#### Run container as systemd service
```bash
❯ sudo cntrctl start photon4
❯ sudo systemctl status photon4
● photon4.service - Photon OS container photon4
   ● photon4.service - Photon OS container
     Loaded: loaded (8;;file://zeus/usr/lib/systemd/system/photon4.12.service^G/usr/lib/systemd/system/photon4.12.service8;;^G; disabled; preset: enabled)
     Active: active (running) since Sun 2022-11-27 13:16:28 UTC; 16s ago
       Docs: 8;;man:cntrctl(1)^Gman:cntrctl(1)8;;^G
   Main PID: 194027 (systemd-nspawn)
     Status: "Container running: Startup finished in 4.458s."
      Tasks: 1 (limit: 16384)
     Memory: 1.1M
     CGroup: /machine.slice/photon4.12.service
             └─194027 /usr/bin/systemd-nspawn --capability=all -bD /var/lib/machines/photon4.12 --link-journal=try-guest -M

Nov 27 13:16:32 zeus cntrctl[194027]: ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBOql3dIN0eJ/4NLKSaOV7aPc9luOtpDcRWs5xs9+13vS8qVR6XIBshv3TwmUu+8NP+>
Nov 27 13:16:32 zeus cntrctl[194027]: ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOq6/QRG07DMxIzy/7/gTB0hsdJfNP5FVZyvyO5agJyq root@photon4
Nov 27 13:16:32 zeus cntrctl[194027]: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSE0byx3ZW75bAdhKNUohEBppQuxjWzQSoFTvZ9K6tfnoyV+qsFgh3nViT+XJWoE6fONNpdWRolNWYiyhRiu1JIsITQZTVbIM5kXN>
Nov 27 13:16:32 zeus cntrctl[194027]: -----END SSH HOST KEY KEYS-----
Nov 27 13:16:32 zeus cntrctl[194027]: [  OK  ] Finished Execute cloud user/final scripts.
Nov 27 13:16:32 zeus cntrctl[194027]: [  OK  ] Reached target Cloud-init target.
Nov 27 13:16:32 zeus cntrctl[194027]: [  OK  ] Stopped OpenSSH Daemon.
Nov 27 13:16:32 zeus cntrctl[194027]: [  OK  ] Started OpenSSH Daemon.
Nov 27 13:16:33 zeus cntrctl[194027]:
Nov 27 13:16:33 zeus cntrctl[194027]: Welcome to Photon 4.0 (x86_64) - Kernel 5.10.142-2.ph4 (console)
```
#### Login to the container
```bash
❯ sudo machinectl login photon4
Connected to machine photon4. Press ^] three times within 1s to exit session.

Welcome to Photon 4.0 (x86_64) - Kernel 5.10.142-2.ph4 (pts/1)
photon4 login: root
Password:
root@photon4 [ ~ ]#
```

#### Creating container with macvlan network
```
❯ cntrctl spawn --network macvlan --link eth0 ph4-macvlan
```
```
❯ systemctl start ph4-macvlan
```
```
❯ systemctl status ph4-macvlan
● ph4-macvlan.service - Photon OS container
     Loaded: loaded (8;;file://zeus/usr/lib/systemd/system/ph4-macvlan.service^G/usr/lib/systemd/system/ph4-macvlan.service8;;^G; disabled; preset: enabled)
     Active: active (running) since Sun 2022-11-27 14:04:29 UTC; 6s ago
       Docs: 8;;man:cntrctl(1)^Gman:cntrctl(1)8;;^G
   Main PID: 2818 (systemd-nspawn)
     Status: "Container running: Startup finished in 4.307s."
      Tasks: 1 (limit: 16384)
     Memory: 1.2M
     CGroup: /machine.slice/ph4-macvlan.service
             └─2818 /usr/bin/systemd-nspawn --capability=all -bD /var/lib/machines/ph4-macvlan --network-macvlan=eth0 --link-journal=try-guest -M

Nov 27 14:04:33 zeus cntrctl[2818]: [  OK  ] Reached target Login Prompts.
Nov 27 14:04:33 zeus cntrctl[2818]: [  OK  ] Reached target Multi-User System.
Nov 27 14:04:33 zeus cntrctl[2818]:          Starting Update UTMP about System Runlevel Changes...
Nov 27 14:04:33 zeus cntrctl[2818]: [  OK  ] Finished Update UTMP about System Runlevel Changes.
Nov 27 14:04:33 zeus cntrctl[2818]: [  OK  ] Finished Apply the settings specified in cloud-config.
Nov 27 14:04:33 zeus cntrctl[2818]:          Starting Execute cloud user/final scripts...
Nov 27 14:04:34 zeus cntrctl[2818]: [  OK  ] Finished Execute cloud user/final scripts.
Nov 27 14:04:34 zeus cntrctl[2818]: [  OK  ] Reached target Cloud-init target.
Nov 27 14:04:35 zeus cntrctl[2818]:
Nov 27 14:04:35 zeus cntrctl[2818]: Welcome to Photon 4.0 (x86_64) - Kernel 5.10.152-3.ph4 (console)
```

```
❯ machinectl login  ph4-macvlan
Connected to machine ph4-macvlan. Press ^] three times within 1s to exit session.

Welcome to Photon 4.0 (x86_64) - Kernel 5.10.152-3.ph4 (pts/1)
ph4-macvlan login: root
Password:

root@ph4-macvlan [ ~ ]# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: mv-eth0@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 16:d3:88:53:c7:be brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 192.168.103.177/23 brd 192.168.103.255 scope global dynamic mv-eth0
       valid_lft 7174sec preferred_lft 7174sec
    inet6 fe80::14d3:88ff:fe53:c7be/64 scope link
       valid_lft forever preferred_lft forever
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
