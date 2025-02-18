## Install

### Toolchain

1. Make sure you installed go https://go.dev/doc/install
2. Install `go install github.com/okt-limonikas/bublik-log-viewer/cmd/blv@latest`
3. Run `blv --help` to get started

This will install binary to `$GOPATH`

### Script

At the root of the project is an `install.sh` script to download and install the binary.

```shell
curl -fsSL \
 https://raw.githubusercontent.com/okt-limonikas/bublik-log-viewer/master/install.sh |\
 sh
```

- The default output directory is `${HOME}/.local/bin`, but can be changed by setting `LOG_INSTALL`.
- **Do not include `/bin`, it is added by the script.**

### Docker
If you have docker installed, you can use the following commands to run the binary.

1. Local Logs
```bash
docker run -it -p 5050:5050 -v <path/to/logs>:/root/json ghcr.io/okt-limonikas/blv:latest
```

2. Remote Logs
```bash
docker run -it -p 5050:5050 ghcr.io/okt-limonikas/blv:latest serve https://ts-factory.io/logs/dpdk-ethdev-ts/2025/02/18/arm-ampere-xl710-root-18:44:06 --remote
```
