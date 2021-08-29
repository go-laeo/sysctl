// Package sysctl inspired by systemd's sysctl source code and kubernetes's sysctl subpackage.
//
// Example:
// 		sysctl.Set("net.ipv6.conf.all.disable_ipv6", "1") // disables all ipv6
// 		value, err := sysctl.Get("net.ipv6.conf.all.disable_ipv6") // retrieves kernel settings value
package sysctl

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func setupPropertyFile(t *testing.T, sysPath, property string) {
	target := path.Join(sysPath, Normalize(property))
	if err := os.MkdirAll(path.Dir(target), os.ModePerm|os.ModeDir); err != nil {
		t.Errorf("Setup target dir error = %v", err)
	}
}

func TestCustomSet(t *testing.T) {
	type args struct {
		procSysPath string
		property    string
		value       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Disable ipv6 sould succeed",
			args: args{
				procSysPath: t.TempDir(),
				property:    "net.ipv6.conf.all.disable_ipv6",
				value:       "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPropertyFile(t, tt.args.procSysPath, tt.args.property)

			if err := CustomSet(tt.args.procSysPath, tt.args.property, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CustomSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		property string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normalize standard property style should succeed",
			args: args{
				property: "net.ipv6.conf.all.disable_ipv6",
			},
			want: "net/ipv6/conf/all/disable_ipv6",
		},
		{
			name: "Normalize if first separator is slash should succeed",
			args: args{
				property: "net/ipv6.conf.all.disable_ipv6",
			},
			want: "net/ipv6.conf.all.disable_ipv6",
		},
		{
			name: "Normalize returns should contains a slash",
			args: args{
				property: "net.ipv6.conf.all/disable_ipv6",
			},
			want: "net/ipv6/conf/all.disable_ipv6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.property); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomGet(t *testing.T) {
	type args struct {
		procSysPath string
		property    string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "CustomGet should read ipv6 disable state",
			args: args{
				procSysPath: t.TempDir(),
				property:    "net.ipv6.conf.all.disable_ipv6",
			},
			want: []byte("1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPropertyFile(t, tt.args.procSysPath, tt.args.property)
			CustomSet(tt.args.procSysPath, tt.args.property, "1")

			got, err := CustomGet(tt.args.procSysPath, tt.args.property)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
