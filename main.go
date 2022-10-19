package main

import (
	"createChain/config"
	"createChain/generate"
	"createChain/stringConst"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")
	flag.Parse()

}
func main() {
	config.LoadConfig(configPath)
	nodes, extra := generate.GenerateNodes(config.Conf.NodeNum, config.Conf.Machines) //生成节点和extradata
	err := generate.GenerateConfig(nodes, config.Conf.ChainPath)                      //生成节点配置信息
	if err != nil {
		panic(err)
	}
	genesis := stringConst.GenerateGenesis(big.NewInt(int64(config.Conf.ChainID)), common.Hex2Bytes(extra[2:]), nodes) //生成genesis文件
	static := stringConst.GenerateStatic(nodes)
	generate.MakeDir(nodes)                                            //生成文件夹和脚本文件
	err = generate.WriteGenesis(nodes, genesis, config.Conf.ChainPath) //写入genesis文件
	if err != nil {
		panic(err)
	}
	err = generate.WriteStatic(nodes, static, config.Conf.ChainPath) //写入static-nodes文件
	if err != nil {
		panic(err)
	}
	err = generate.WriteKey(nodes, config.Conf.ChainPath) //写入公私钥对
	if err != nil {
		panic(err)
	}
	err = generate.WriteBash(nodes, config.Conf)
	if err != nil {
		panic(err)
	}
}
