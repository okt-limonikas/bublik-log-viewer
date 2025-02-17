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
