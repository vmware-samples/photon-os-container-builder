### photon-os-container-builder
----
***containerctl*** is CLI tool which spawns Photon OS in a light-weight container. It uses `systemd-nspawn` to jump start Photon OS containers. The primary
use case for ***scontainerctl*** is to run Photon OS test cases in isolated environment. 


```
> sudo containerctl spawn -p systemd,iproute -r 4.0  my-build

Installing:
nettle                   x86_64       3.7.2-1.ph4      photon-updates 731.22k 748769
glib                     x86_64       2.68.0-1.ph4     photon-updates   3.53M 3697027
coreutils-selinux        x86_64       8.32-2.ph4       photon-release   6.85M 7181874
expat-libs               x86_64       2.2.9-3.ph4      photon-release 185.24k 189688
expat                    x86_64       2.2.9-3.ph4      photon-release  30.04k 30760
e2fsprogs-libs           x86_64       1.45.6-2.ph4     photon-release  58.37k 59768

Complete!

Spawning container my-build on /var/lib/machines/my-build.
Press ^] three times within 1s to kill container.
root@my-build [ ~ ]# 
```

```
> sudo containerctl --help
NAME:
   containerctl - A new cli application

USAGE:
   containerctl[global options] command [command options] [arguments...]

VERSION:
   v0.1

DESCRIPTION:
   Spawns nspawn containers

COMMANDS:
   spawn, s      [NAME] Spawn a container
   boot, b       [NAME] Boot a container
   dir, d        [NAME] Directory to use as file system root for the container
   login         [NAME] Login prompt in a container
   shell         [USER@NAME] [COMMAND] Invoke a shell (or other command) in a container
   copy-to       [NAME] [PATH] PATH Copy files from the host to a container
   copy-from     [NAME] PATH [PATH] Copy files from a container to the host
   list          List running VMs and containers
   list-images   List container images
   image-status  [NAME] Display image status
   remove        [NAME] Remove an image
   clone         [NAME] Clone an image
   rename        [NAME] Rename an image
   clean         [NAME] Remove all images
   poweroff      [NAME] poweroff a container
   terminate     [NAME] Terminate a container
   reboot        [NAME] Reboot a container
   pull-tar      URL [NAME] Download a TAR container image
   pull-raw      URL [NAME] Download a RAW container or VM image
   import-tar    FILE [NAME] Import a local TAR container image
   import-raw    FILE [NAME] Import a local RAW container or VM image
   import-fs     DIRECTORY [NAME] Import a local directory container image
   export-tar    NAME [FILE] Export a TAR container image locally
   export-raw    AME [FILE] Export a RAW container or VM image locally
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

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

[Apache-2.0](https://spdx.org/licenses/Apache-2.0.html)
