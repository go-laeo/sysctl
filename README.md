# sysctl

![build.yaml](https://github.com/go-laeo/sysctl/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/sysctl.svg)](https://pkg.go.dev/github.com/go-laeo/sysctl) ![golangci.yaml](https://github.com/go-laeo/sysctl/actions/workflows/golangci-lint.yaml/badge.svg)

Package sysctl inspired by [systemd's sysctl source code](https://github.com/systemd/systemd/blob/main/src/sysctl/sysctl.c) and [kubernetes's sysctl subpackage](https://github.com/kubernetes/kubernetes/blob/v1.22.1/pkg/util/sysctl/sysctl.go).

## Install

```bash
go get github.com/go-laeo/sysctl
```

## Example

```golang
// disables all ipv6
sysctl.Set("net.ipv6.conf.all.disable_ipv6", "1")

// check if ipv6 has disabled
disabled, err := sysctl.GetBool("net.ipv6.conf.all.disable_ipv6")
```

## License

The MIT License.
