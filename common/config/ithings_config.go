package config

import (
	"sync"

	_ "k8s.io/klog/v2"
)

const ()

var ITHINGS_CONFIG *Config
var once = sync.Once{}

func ForceInit() {
	ITHINGS_CONFIG = NewYamlConfig("config.yaml")
}

func init() {
	once.Do(func() {
		//load the config.yaml from conf/
		ITHINGS_CONFIG = NewYamlConfig("config.yaml")
	})
}
