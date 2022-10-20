package config

import (
	"createChain/files"
	"createChain/types"
	"encoding/json"
	"fmt"
)

var (
	Conf       *Config
	Validators []types.Account
	Signers    []types.Account
)

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
	AutoAlloc bool
}

func LoadConfig(filepath string) {
	data, err := files.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &Conf); err != nil {
		panic(err)
	}
	if Conf.AutoAlloc {
		mcs := make([]MachineConfig, 0)
		for i := 0; i < Conf.NodeNum; i++ {
			mc := MachineConfig{
				MachineIP: "localhost",
				P2PPort:   fmt.Sprintf("3030%d", i),
				RPCPort:   fmt.Sprintf("2200%d", i),
			}
			mcs = append(mcs, mc)
		}
		Conf.Machines = mcs
	} else {
		if Conf.NodeNum != len(Conf.Machines) {
			panic("node number must be equal to length of machines!")
		}
	}
}

func LoadValidators(filepath string) {
	data, err := files.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &Validators); err != nil {
		panic(err)
	}
}

func LoadSigners(filepath string) {
	data, err := files.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &Signers); err != nil {
		panic(err)
	}
}
