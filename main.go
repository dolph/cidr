package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
)

func formatRange(cidr string, firstIP net.IP, lastIP net.IP, length int) string {
	plurality := ""
	if length != 1 {
		plurality = "es"
	}
	return fmt.Sprintf("%s: %s-%s (%s address%s)", cidr, firstIP, lastIP, formatWithCommas(length), plurality)
}

func formatWithCommas(number int) string {
	negative := false
	if number < 0 {
		negative = true
		number = -number
	}

	numStr := strconv.Itoa(number)
	length := len(numStr)
	commaIndexes := (length - 1) / 3
	totalLength := length + commaIndexes

	if negative {
		totalLength++
	}

	var result bytes.Buffer
	result.Grow(totalLength)

	if negative {
		result.WriteByte('-')
	}

	for i, digit := range numStr {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteByte(byte(digit))
	}

	return result.String()
}

func main() {
	if len(os.Args) == 2 {
		cidr := os.Args[1]
		err := PrintRange(cidr)
		if err != nil {
			fmt.Println("Invalid CIDR notation:", err)
			os.Exit(1)
		}
	} else if len(os.Args) == 3 {
		firstIP := net.ParseIP(os.Args[1])
		lastIP := net.ParseIP(os.Args[2])
		if firstIP == nil || lastIP == nil || firstIP.To4() == nil || lastIP.To4() == nil {
			fmt.Println("Invalid IP addresses provided.")
			os.Exit(1)
		}

		ToCIDRs(firstIP, lastIP)
	} else {
		fmt.Println("Usage: cidr [<CIDR> | <First-IP> <Last-IP>]")
		fmt.Println()
		fmt.Println("cidr <CIDR>: Calculates the first and last IP in a CIDR range.")
		fmt.Println()
		fmt.Println("cidr <First-IP> <Last-IP>: Calculates the minimal number of CIDRs to express the given range.")
		os.Exit(2)
	}
}

func PrintRange(cidr string) error {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	maskSize, _ := ipNet.Mask.Size()
	length := 1 << (32 - uint(maskSize))

	firstIP := ip.Mask(ipNet.Mask)
	lastIP := make(net.IP, len(firstIP))
	copy(lastIP, firstIP)

	for i := range firstIP {
		lastIP[i] = firstIP[i] | ^ipNet.Mask[i]
	}

	fmt.Println(formatRange(cidr, firstIP, lastIP, length))
	return nil
}

func ToUint32(ip net.IP) uint32 {
	ipBytes := ip.To4()
	return uint32(ipBytes[0])<<24 | uint32(ipBytes[1])<<16 | uint32(ipBytes[2])<<8 | uint32(ipBytes[3])
}

func ToIP(val uint32) net.IP {
	return net.IPv4(byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func ToCIDRs(firstIP, lastIP net.IP) {
	start := ToUint32(firstIP)
	end := ToUint32(lastIP)

	if start > end {
		// Swap the first and last values
		temp := start
		start = end
		end = temp
	}

	for start <= end {
		maskSize := 0
		for maskSize < 32 {
			mask := uint32(0xffffffff) << uint(32-maskSize)
			if (start&mask) == start && (start|^mask) <= end {
				break
			}
			maskSize++
		}

		cidr := fmt.Sprintf("%s/%d", ToIP(start), maskSize)
		PrintRange(cidr)

		start += 1 << uint(32-maskSize)
	}
}
