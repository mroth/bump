package main

import (
	"testing"
)

func Test_parsePreStr(t *testing.T) {
	tests := []struct {
		pre         string
		wantPrefix  string
		wantCounter uint
		wantErr     bool
	}{
		{"pre1", "pre", 1, false},
		{"pre.1", "pre.", 1, false},
		{"rc32", "rc", 32, false},
		{"", "", 0, true},
		{"foobar", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.pre, func(t *testing.T) {
			gotPrefix, gotCounter, err := parsePreStr(tt.pre)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePreStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPrefix != tt.wantPrefix {
				t.Errorf("parsePreStr() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotCounter != tt.wantCounter {
				t.Errorf("parsePreStr() gotCounter = %v, want %v", gotCounter, tt.wantCounter)
			}
		})
	}
}

func Test_incPreStr(t *testing.T) {
	tests := []struct {
		pre     string
		want    string
		wantErr bool
	}{
		{"pre1", "pre2", false},
		{"pre.1", "pre.2", false},
		{"dev01", "dev02", false}, // preserve leading zeros
		{"dev.01", "dev.02", false},
		{"rc32", "rc33", false}, // no sep
		{"", "", true},
		{"foobar", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.pre, func(t *testing.T) {
			got, err := incPreStr(tt.pre)
			if (err != nil) != tt.wantErr {
				t.Errorf("incPreStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("incPreStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
