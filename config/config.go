package config

import (
	"createChain/files"
	"encoding/json"
)

var Conf *Config

type MachineConfig struct {
	MachineIP string
	P2PPort   string
	RPCPort   string
}
type Config struct {
	NodeNum   int
	ChainID   uint64
	Machines  []MachineConfig
	ChainPath string
}

func LoadConfig(filepath string) {
	data, err := files.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &Conf); err != nil {
		panic(err)
	}
	if Conf.NodeNum != len(Conf.Machines) {
		panic("node number must be equal to length of machines!")
	}

}
