package main

import (
	"createChain/config"
	"createChain/generate"
	"createChain/log"
	"createChain/stringConst"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var configPath string
var loglevel int

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")
	flag.IntVar(&loglevel, "loglevel", 2, "loglevel [1: debug, 2: info]")
	flag.Parse()

}
func main() {
	log.InitLog(loglevel, log.Stdout)
	config.LoadConfig(configPath)
	log.Info("加载配置文件成功！")
	nodes, extra := generate.GenerateNodes(config.Conf.NodeNum, config.Conf.Machines) //生成节点和extradata
	log.Info("生成节点和extra data成功！")
	err := generate.GenerateConfig(nodes, config.Conf.ChainPath) //生成节点配置信息

	if err != nil {
		panic(err)
	}
	log.Info("生成节点配置信息成功！")
	genesis := stringConst.GenerateGenesis(big.NewInt(int64(config.Conf.ChainID)), common.Hex2Bytes(extra[2:]), nodes) //生成genesis文件
	static := stringConst.GenerateStatic(nodes)
	log.Info("生成genesis和static文件成功！")
	generate.MakeDir(nodes) //生成文件夹和脚本文件
	log.Info("创建文件夹成功！")
	err = generate.WriteGenesis(nodes, genesis, config.Conf.ChainPath) //写入genesis文件
	if err != nil {
		panic(err)
	}
	log.Info("写入genesis文件成功！")
	err = generate.WriteStatic(nodes, static, config.Conf.ChainPath) //写入static-nodes文件
	if err != nil {
		panic(err)
	}
	log.Info("写入static-nodes文件成功！")
	err = generate.WriteKey(nodes, config.Conf.ChainPath) //写入公私钥对
	if err != nil {
		panic(err)
	}
	log.Info("写入公私钥对成功！")
	err = generate.WriteBash(nodes, config.Conf)
	if err != nil {
		panic(err)
	}
	log.Info("写入脚本文件成功！")
}
