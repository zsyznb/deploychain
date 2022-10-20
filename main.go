package main

import (
	"createChain/config"
	"createChain/generate"
	"createChain/log"
	"createChain/stringConst"
	"flag"
)

var (
	configPath    string
	loglevel      int
	ValidatorPath string
	SignerPath    string
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")
	flag.StringVar(&ValidatorPath, "validator", "validators.json", "validator config path")
	flag.StringVar(&SignerPath, "signer", "signers.json", "signer config path")
	flag.IntVar(&loglevel, "loglevel", 2, "loglevel [1: debug, 2: info]")
	flag.Parse()

}
func main() {
	//初始化log,加载配置文件
	log.InitLog(loglevel, log.Stdout)
	config.LoadConfig(configPath)
	config.LoadValidators(ValidatorPath)
	config.LoadSigners(SignerPath)
	log.Info("加载配置文件成功！")

	//生成节点和extra data
	nodes, extra := generate.GenerateNodes(config.Conf, config.Validators, config.Signers)
	log.Info("生成节点和extra data成功！")

	//生成genesis文件和static文件
	log.Info("生成节点配置信息成功！")
	genesis := stringConst.GenerateGenesis(config.Conf, extra, nodes) //生成genesis文件
	static := stringConst.GenerateStatic(nodes)                       //生成static文件
	log.Info("生成genesis和static文件成功！")

	//生成文件夹和脚本文件
	generate.MakeDir(nodes, config.Conf) //生成文件夹和脚本文件
	log.Info("创建文件夹成功！")
	err := generate.WriteGenesis(nodes, genesis, config.Conf) //写入genesis文件
	if err != nil {
		panic(err)
	}
	log.Info("写入genesis文件成功！")
	err = generate.WriteStatic(nodes, static, config.Conf) //写入static-nodes文件
	if err != nil {
		panic(err)
	}
	log.Info("写入static-nodes文件成功！")
	err = generate.WriteKey(nodes, config.Conf) //写入公私钥对
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
