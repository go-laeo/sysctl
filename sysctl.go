// Package sysctl inspired by systemd's sysctl source code and kubernetes's sysctl subpackage.
//
// Example:
// 		sysctl.Set("net.ipv6.conf.all.disable_ipv6", "1") // disables all ipv6
// 		value, err := sysctl.Get("net.ipv6.conf.all.disable_ipv6") // retrieves kernel settings value
package sysctl

import (
	"bytes"
	"os"
	"path"
	"strconv"
)

const (
	defaultsProcSysPath = "/proc/sys"
)

// Set writes given value and a special char "\n" to /proc/sys.
func Set(property, value string) error {
	return CustomSet(defaultsProcSysPath, property, value)
}

func CustomSet(procSysPath, property, value string) error {
	return os.WriteFile(path.Join(procSysPath, Normalize(property)), []byte(value+"\n"), 0640)
}

// Get reads property value from in /proc/sys.
func Get(property string) ([]byte, error) {
	return CustomGet(defaultsProcSysPath, property)
}

func GetInt(property string) (int, error) {
	b, err := Get(property)
	if err != nil {
		return -1, err
	}

	if len(b) == 0 {
		return 0, nil
	}

	return strconv.Atoi(string(b))
}

func GetBool(property string) (bool, error) {
	b, err := Get(property)
	if err != nil {
		return false, err
	}

	if len(b) == 0 {
		return false, nil
	}

	return b[0] == '1', nil
}

func CustomGet(procSysPath, property string) ([]byte, error) {
	b, err := os.ReadFile(path.Join(procSysPath, Normalize(property)))
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

// Normalize will checks if the first separator is a slash,
// the path is assumed to be normalized and slashes remain slashes
// and dots remains dots.
func Normalize(property string) string {
	b := []byte(property)
	pos := bytes.IndexAny(b, "/.")
	if pos != -1 && b[pos] == '.' {
		for ; pos < len(b)-1; pos++ {
			if b[pos] == '/' || b[pos] == '.' {
				if b[pos] == '.' {
					b[pos] = '/'
				} else {
					b[pos] = '.'
				}
			}
		}
	}

	return path.Clean(string(b))
}
