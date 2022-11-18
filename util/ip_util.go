package util

import (
	"fmt"
	"log"
	"net"
)

type LocalIp struct {
	LocalIp string
	Ip4     []byte
}

var li *LocalIp

func GetLocalIp4() []byte {
	if li != nil {
		return li.Ip4
	}
	findAddr()
	return li.Ip4
}
func GetLocalIp() string {
	if li != nil {
		return li.LocalIp
	}
	findAddr()
	return li.LocalIp
}

func findAddr() (string, bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("获取本地地址异常", err)
		return "", false
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := fmt.Sprintf(ipnet.IP.String())
				li = &LocalIp{LocalIp: ip, Ip4: ipnet.IP.To4()}
				return li.LocalIp, true
			}
		}
	}
	return "", false
}
