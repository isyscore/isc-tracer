package main

import (
	c "github.com/isyscore/isc-tracer/config"
)

func newTracerClient(config *c.Config) {
	c.ServerConfig = config
	if c.ServerConfig.Enable {
		//todo 初始化

	}

}
