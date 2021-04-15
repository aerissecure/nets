// A library for the nets...that is, the internets. Actually the networks.

package nets

import (
	"math"
	"net"
)

// IPInc increments the ip, updating the value in place. This method should be
// used when all IP address values in a range are needed. To only get the final
// value after x increments, use IPAdd instead as it is much faster.
//
// Be careful when calling with an IPv4 address in 16-byte form, such as that
// returned by net.IPv4() as callin git on 255.255.255.255 will produce
// unexpected results. Call on a 4-byte IPv4 address or an IPv6 address will
// predicably roll over to 0.0.0.0 or its equivalent.
func IPInc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

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

// IPToI32 converts a net.IP to an int32.
func IPToI32(ip net.IP) int32 {
	ip = ip.To4()
	return int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3])
}

// I32ToIP converts an int32 to a net.IP.
func I32ToIP(a int32) net.IP {
	return net.IPv4(byte(a>>24), byte(a>>16), byte(a>>8), byte(a))
}

// IPToU64 converts a net.IP to a uint64.
func IPToU64(ip net.IP) uint64 {
	return uint64(ip[0])<<56 | uint64(ip[1])<<48 | uint64(ip[2])<<40 |
		uint64(ip[3])<<32 | uint64(ip[4])<<24 | uint64(ip[5])<<16 |
		uint64(ip[6])<<8 | uint64(ip[7])
}

// u64ToIP converts a uint64 to a net.IP.
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

// IPAdd increments the offset amount. A new IP address is returned. Use this
// function when only the resulting IP address is nedded. If intermediate values
// are needed, such as when calling with offset=1 in a loop, use IPInc instead
// as it is much faster.
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

// IPMaskCount returns the number of addresses in a subnet with that mask
func IPMaskCount(m net.IPMask) int {
	ones, bits := m.Size()
	return 1 << (uint(bits - ones))
}

// CopyIP makes a newly allocated copy of the provided ip address.
func CopyIP(ip net.IP) net.IP {
	if ip == nil {
		ip = net.IPv4zero
	}
	dst := make(net.IP, len(ip))
	copy(dst, ip)
	return dst
}

// NetsOverlap returns true if the two networks have overlapping ip addresses.
func NetsOverlap(netA net.IPNet, netB net.IPNet) bool {
	return netA.Contains(netB.IP) || netB.Contains(netA.IP)
}

// NetRange returns the first (network) and last (broadcast) IP addresses from
// a net.IPNet
func NetRange(ipnet net.IPNet) (net.IP, net.IP) {
	first := ipnet.IP.Mask(ipnet.Mask)
	last := make(net.IP, len(first))
	for i, b := range first {
		last[i] = b | ^ipnet.Mask[i]
	}

	if ipnet.IP.To4() != nil {
		first = first.To4()
		last = last.To4()
	}
	return first, last
}

// IPLessThan returns true is IP a is smaller than IP b.
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

// IPNetNet takes an IPNet and returns an IPNet.IP that is now the first address
// on the network. E.g. when given 10.0.0.20/24, it returns 10.0.0.0/24.
func IPNetNet(ipnet net.IPNet) *net.IPNet {
	ip, _ := NetRange(ipnet)
	return &net.IPNet{IP: ip, Mask: ipnet.Mask}
}
