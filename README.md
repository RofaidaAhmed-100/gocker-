# GOCKER

A minimal container runtime in Go. Shows that Docker is just Linux features combined together.

## What it does?

Runs a process in an isolated environment using:
- **Namespaces** — isolated process tree, filesystem, and network
- **cgroups v2** — limits RAM to 512MB and CPU weight to 50
- **chroot** — Alpine Linux as the container filesystem
- **veth + bridge** — virtual network interface


## Setup

**1. Clone the repo:**
```bash
git clone https://github.com/RofaidaAhmed-100/gocker-.git
cd gocker
```

**2. Install Alpine rootfs:**
```bash
mkdir -p /tmp/gocker/alpine && cd /tmp/gocker/alpine
wget https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz
tar xzf alpine-minirootfs-3.19.0-x86_64.tar.gz
rm alpine-minirootfs-3.19.0-x86_64.tar.gz
cd -
```

**3. Run:**
```bash
sudo go run main.go run /bin/sh
```



## Requirements

- Linux
- Go 1.21+
- Alpine rootfs at `/tmp/gocker/alpine`
- `ip` tool (`sudo apt install iproute2`)
