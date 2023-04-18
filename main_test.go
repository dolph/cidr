package main

import (
	"bytes"
	"net"
	"testing"
)

func TestToIP(t *testing.T) {
	tests := []struct {
		val      uint32
		expected net.IP
	}{
		{0, net.IPv4(0, 0, 0, 0)},
		{3232235521, net.IPv4(192, 168, 0, 1)},
		{4294967295, net.IPv4(255, 255, 255, 255)},
	}

	for _, test := range tests {
		result := ToIP(test.val)
		if !result.Equal(test.expected) {
			t.Errorf("ToIP(%v) = %v; expected %v", test.val, result, test.expected)
		}
	}
}

func TestToUint32(t *testing.T) {
	tests := []struct {
		ip       net.IP
		expected uint32
	}{
		{net.IPv4(0, 0, 0, 0), 0},
		{net.IPv4(127, 0, 0, 1), 2130706433},
		{net.IPv4(192, 168, 0, 1), 3232235521},
		{net.IPv4(255, 255, 255, 255), 4294967295},
	}

	for _, test := range tests {
		result := ToUint32(test.ip)
		if result != test.expected {
			t.Errorf("ToUint32(%v) = %v; expected %v", test.ip, result, test.expected)
		}
	}
}

func TestToCIDRs(t *testing.T) {
	tests := []struct {
		firstIP  net.IP
		lastIP   net.IP
		expected []string
	}{
		{
			net.ParseIP("192.168.0.1"),
			net.ParseIP("192.168.0.10"),
			[]string{"192.168.0.1/32", "192.168.0.2/31", "192.168.0.4/30", "192.168.0.8/31", "192.168.0.10/32"},
		},
		{
			net.ParseIP("10.0.0.1"),
			net.ParseIP("10.0.0.255"),
			[]string{"10.0.0.1/32", "10.0.0.2/31", "10.0.0.4/30", "10.0.0.8/29", "10.0.0.16/28", "10.0.0.32/27", "10.0.0.64/26", "10.0.0.128/25"},
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		PrintRange = func(cidr string) {
			buf.WriteString(cidr)
			buf.WriteByte('\n')
		}

		ToCIDRs(test.firstIP, test.lastIP)

		got := buf.String()
		for i, expected := range test.expected {
			if !bytes.Contains([]byte(got), []byte(expected)) {
				t.Errorf("TestToCIDRs[%d]:\nExpected CIDR range: %q\nActual CIDR range: %q\n", i, expected, got)
			}
		}
	}
}
