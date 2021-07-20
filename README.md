### photon-os-container-builder
----
***containerctl*** is CLI tool which spawns Photon OS in a light-weight container. It uses `systemd-nspawn` to jump start Photon OS containers. The primary
use case for ***containerctl*** is to run Photon OS test cases in isolated environment.


```
> sudo containerctl spawn test55

Refreshing metadata for: 'VMware Photon Linux 4.0 (x86_64) Updates'
Refreshing metadata for: 'VMware Photon Linux 4.0 (x86_64)'
Refreshing metadata for: 'VMware Photon Extras 4.0 (x86_64)'

Installing:
nettle                   x86_64       3.7.2-1.ph4      photon-updates 731.22k 748769
glib                     x86_64       2.68.0-1.ph4     photon-updates   3.53M 3697027
coreutils-selinux        x86_64       8.32-2.ph4       photon-release   6.85M 7181874
lz4                      x86_64       1.9.2-2.ph4      photon-updates 464.87k 476022


Complete!
```

```
> sudo containerctl dir test55
Spawning container test55 on /var/lib/machines/test55.
Press ^] three times within 1s to kill container.
root@test55 [ ~ ]#
```


```
> > containerctl
NAME:
   containerctl - A new cli application

USAGE:
   containerctl [global options] command [command options] [arguments...]

VERSION:
   v0.1

DESCRIPTION:
   Spawns Photon OS containers

COMMANDS:
   spawn, s      [NAME] Spawn a container
   boot, b       [NAME] Boot a container
   dir, d        [NAME] Directory to use as file system root for the container

```

#### Build
```
❯  make build
❯  make install
```


#### Contributing
----

The ***photon-os-container-builder*** project team welcomes contributions from the community. If you wish to contribute code and you have not signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For any questions about the CLA process, please refer to our [FAQ](https://cla.vmware.com/faq).

slack channel [#photon](https://code.vmware.com/web/code/join).

#### License
----

[BSD 2-Clause](https://spdx.org/licenses/BSD-2-Clause.html)
