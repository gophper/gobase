package gobase

import (
	"fmt"
	"testing"
)

func Test_GetIpAddrs(t *testing.T) {
	ip0, _ := GetIpAddrs("39.128.0.0", "10", false)
	fmt.Println("TestGetIpAddrs0", ip0.IpEnd, ip0.IpStart)
	ip1, ip2 := GetIpAddrs("192.168.0.1", "28", true)
	fmt.Println("TestGetIpAddrs1", ip1.IpEnd, ip1.IpStart)
	for _, v := range ip2 {
		fmt.Println("TestGetIpAddrs2", v)
	}
}

func Test_GetIpMask(t *testing.T) {
	ip1 := GetIpMask("211.103.0.0", "211.103.127.255")
	fmt.Println("TestGetIpMask", ip1)
	ip2 := GetIpMask("192.168.0.1", "192.175.0.100")
	fmt.Println("TestGetIpMask", ip2)

}
