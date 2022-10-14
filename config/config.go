package config

import (
	"createChain/files"
	"encoding/json"
)

var Conf *Config

type Config struct {
	NodeNum   int
	ChainID   uint64
	Machines  []string
	P2PPort   uint64
	RPCPort   uint64
	ChainPath string
}

func LoadConfig(filepath string, nodenum int) {
	data, err := files.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &Conf); err != nil {
		panic(err)
	}
	if nodenum > 0 {
		Conf.NodeNum = nodenum
	}

}
