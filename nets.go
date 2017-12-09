// A library for the nets...that is, the internets. Actually the networks.

package nets

import (
	"math"
	"net"
)

// sources:
// https://github.com/freb/netaddr/blob/master/netaddr/ip.go
// https://github.com/docker/libnetwork/blob/master/netutils/utils.go

// I need to rewrite all the functions I want to use
// I should also check the speed of IPAdd against the code by russ cox that everyone uses:
//
//  http://play.golang.org/p/m8TNTtygK0
// func inc(ip net.IP) {
// 	for j := len(ip) - 1; j >= 0; j-- {
// 		ip[j]++
// 		if ip[j] > 0 {
// 			break
// 		}
// 	}
// }

// if it is faster, just use it, and don't worry about all this guys functions.
// otherwise, use these, because they make more sense, but add descriptions and remove functions that
// i dont need like the check on ipv4/ipv6

func isZeros(p net.IP) bool {
	for _, b := range p {
		if b != 0 {
			return false
		}
	}
	return true
}

// IsIPv4 returns true if ip is IPv4 address.
func IsIPv4(ip net.IP) bool {
	return len(ip) == net.IPv4len ||
		isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff
}

func IPToI32(ip net.IP) int32 {
	ip = ip.To4()
	return int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3])
}

func I32ToIP(a int32) net.IP {
	return net.IPv4(byte(a>>24), byte(a>>16), byte(a>>8), byte(a))
}

func IPToU64(ip net.IP) uint64 {
	return uint64(ip[0])<<56 | uint64(ip[1])<<48 | uint64(ip[2])<<40 |
		uint64(ip[3])<<32 | uint64(ip[4])<<24 | uint64(ip[5])<<16 |
		uint64(ip[6])<<8 | uint64(ip[7])
}

func u64ToIP(ip net.IP, a uint64) {
	ip[0] = byte(a >> 56)
	ip[1] = byte(a >> 48)
	ip[2] = byte(a >> 40)
	ip[3] = byte(a >> 32)
	ip[4] = byte(a >> 24)
	ip[5] = byte(a >> 16)
	ip[6] = byte(a >> 8)
	ip[7] = byte(a)
}

// IPAdd adds offset to ip
func IPAdd(ip net.IP, offset int) net.IP {
	if IsIPv4(ip) {
		a := int(IPToI32(ip[len(ip)-4:]))
		return I32ToIP(int32(a + offset))
	}
	a := IPToU64(ip[:net.IPv6len/2])
	b := IPToU64(ip[net.IPv6len/2:])
	o := uint64(offset)
	if math.MaxUint64-b < o {
		a++
	}
	b += o
	if offset < 0 {
		a += math.MaxUint64
	}
	ip = make(net.IP, net.IPv6len)
	u64ToIP(ip[:net.IPv6len/2], a)
	u64ToIP(ip[net.IPv6len/2:], b)
	return ip
}

// this would help the a UsableNetRange func
// func IPSub

// IPMod calculates ip % d
func IPMod(ip net.IP, d uint) uint {
	if IsIPv4(ip) {
		return uint(IPToI32(ip[len(ip)-4:])) % d
	}
	b := uint64(d)
	hi := IPToU64(ip[:net.IPv6len/2])
	lo := IPToU64(ip[net.IPv6len/2:])
	return uint(((hi%b)*((0-b)%b) + lo%b) % b)
}

// my function
// IPMaskCount returns the number of addresses in a subnet with that mask
func IPMaskCount(m net.IPMask) int {
	ones, bits := m.Size()
	return 1 << (uint(bits - ones))
}

// my function (not really but its ok, simple and renamed)
// slices point to an underlying array, so creating a new variable alone will not create a copy
func CopyIP(from net.IP) net.IP {
	if from == nil {
		return nil
	}
	to := make(net.IP, len(from))
	copy(to, from)
	return to
}

func NetsOverlap(netX *net.IPNet, netY *net.IPNet) bool {
	return netX.Contains(netY.IP) || netY.Contains(netX.IP)
}

// changed name, docker function
// // NetRange
func NetRange(network *net.IPNet) (net.IP, net.IP) {
	if network == nil {
		return nil, nil
	}

	firstIP := network.IP.Mask(network.Mask)
	lastIP := CopyIP(firstIP)
	for i := 0; i < len(firstIP); i++ {
		lastIP[i] = firstIP[i] | ^network.Mask[i] // use xor to create lastIP
	}

	if network.IP.To4() != nil {
		firstIP = firstIP.To4()
		lastIP = lastIP.To4()
	}

	return firstIP, lastIP
}

// func UsableNetRange(network *net.IPNet) (net.IP, net.IP) {
// 	first, last := netRange
// }

// my function
func IPLessThan(a, b net.IP) bool {
	a = a.To16()
	b = b.To16()
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// myfunction
// Given an IPNet return an IPNet where IP is the first address
// on the network. E.g. when given 10.0.0.20/24, return 10.0.0.0/24.
func IPNetNet(ipnet *net.IPNet) *net.IPNet {
	ip, _ := NetRange(ipnet)
	return &net.IPNet{IP: ip, Mask: ipnet.Mask}
}
