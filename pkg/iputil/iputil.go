package iputil

import (
	"net"
)

func Ip2Long(ip net.IP) uint32 {
	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])
	return uint32(a<<24 | b<<16 | c<<8 | d)
}

func Long2Ip(ip uint32) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}

var (
	L_10_0_0_0        = Ip2Long(net.ParseIP("10.0.0.0"))
	L_10_255_255_255  = Ip2Long(net.ParseIP("10.255.255.255"))
	L_172_16_0_0      = Ip2Long(net.ParseIP("172.16.0.0"))
	L_172_31_255_255  = Ip2Long(net.ParseIP("172.31.255.255"))
	L_192_168_0_0     = Ip2Long(net.ParseIP("192.168.0.0"))
	L_192_168_255_255 = Ip2Long(net.ParseIP("192.168.255.255"))
	L_127_0_0_1       = Ip2Long(net.ParseIP("127.0.0.1"))
)

func IsInternalIp(ip net.IP) bool {
	l := Ip2Long(ip)
	if (l >= L_10_0_0_0 && l <= L_10_255_255_255) || (l >= L_172_16_0_0 && l <= L_172_31_255_255) || (l >= L_192_168_0_0 && l <= L_192_168_255_255) || l == L_127_0_0_1 {
		return true
	}
	return false
}

func GetServerIpList() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var list []net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				list = append(list, ipnet.IP)
			}
		}
	}
	return list
}

func GetServerIpListByName(name string) []net.IP {
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		return nil
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil
	}
	var list []net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				list = append(list, ipnet.IP)
			}
		}
	}
	return list
}
