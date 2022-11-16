package util

import (
	"fmt"
	"log"
	"net"
)

type LocalIp struct {
	LocalIp string
}

var li *LocalIp

func GetLocalIp() string {
	if li != nil {
		return li.LocalIp
	}
	li = &LocalIp{
		LocalIp: "127.0.0.1",
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("获取本地地址异常", err)
		return li.LocalIp
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := fmt.Sprintf(ipnet.IP.String())
				li = &LocalIp{ip}
				return li.LocalIp
			}

		}
	}
	return li.LocalIp
}
