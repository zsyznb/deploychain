package generate

import (
	"bytes"
	"createChain/config"
	"createChain/ip"
	"createChain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"os"
	"sort"
	"strconv"
)

func createBash(path string) *os.File { //创建脚本
	bash, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	return bash
}

func createFile(path string) *os.File { //创建文件
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

func MakeDir(nodes types.Nodes, Conf *config.Config) { //创建文件目录
	n := len(nodes)
	machines := Conf.Machines

	for i := 0; i < n; i++ { //创建node文件夹
		fmt.Println(machines[i].MachineIP, ip.IP)
		if machines[i].MachineIP == ip.IP {
			err := os.Mkdir(config.Conf.ChainPath+nodes[i].Name, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

	}

	for i := 0; i < n; i++ { //创建setup和node文件夹
		if machines[i].MachineIP == ip.IP {
			err := os.MkdirAll(config.Conf.ChainPath+nodes[i].Name+"/setup/node", os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}

	bashs := []string{"start.sh", "stop.sh", "init.sh", "build.sh"}
	for i := 0; i < n; i++ { //创建脚本
		if machines[i].MachineIP == ip.IP {
			for _, bash := range bashs {
				createBash(config.Conf.ChainPath + nodes[i].Name + "/" + bash)
			}
		}
	}

	jsons := []string{"genesis.json", "static-nodes.json"}
	for i := 0; i < n; i++ { //创建genesis.json和static.json
		if machines[i].MachineIP == ip.IP {
			for _, json := range jsons {
				createFile(config.Conf.ChainPath + nodes[i].Name + "/setup/" + json)
			}
		}
	}

	keys := []string{"nodekey", "pubkey"}
	for i := 0; i < n; i++ { //创建nodekey和pubkey文件
		if machines[i].MachineIP == ip.IP {
			for _, key := range keys {
				createFile(config.Conf.ChainPath + nodes[i].Name + "/setup/node/" + key)
			}
		}
	}
}

func GenerateNodes(Conf *config.Config, validators []types.Account, signers []types.Account) (types.Nodes, string) { //生成节点
	nodes := make([]*types.Node, 0)
	n := Conf.NodeNum
	machines := Conf.Machines

	//生成Validator账户和Signer账户
	v := validators[:n]
	s := signers[:n]
	Signers := sortAccounts(s)
	Validators := sortAccounts(v)
	for i := 0; i < n; i++ {
		key, _ := crypto.HexToECDSA(Validators[i].NodeKey[2:])
		node := &types.Node{
			Name:      fmt.Sprintf("node%d", i),
			Signer:    Signers[i],
			Validator: Validators[i],
			Static:    fmt.Sprintf("enode://%s@%s:%s?discport=0", PubkeyID(&key.PublicKey), machines[i].MachineIP, machines[i].P2PPort),
		}
		nodes = append(nodes, node)
	}

	extra := generateExtra(NodesAddress(nodes))
	return nodes, extra
}

func PubkeyID(pub *ecdsa.PublicKey) types.Discv5NodeID {
	var id types.Discv5NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

func GenerateAccounts(n int) []types.Account {
	var Accs []types.Account
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		account := types.Account{
			Addr:    crypto.PubkeyToAddress(key.PublicKey),
			NodeKey: hexutil.Encode(crypto.FromECDSA(key)),
			PubKey:  hexutil.Encode(crypto.CompressPubkey(&key.PublicKey)),
		}
		Accs = append(Accs, account)
	}
	Accs = sortAccounts(Accs)
	return Accs
}

func sortAccounts(accs []types.Account) []types.Account { //对账号排序
	oriAddrs := make([]common.Address, 0)
	AccIndex := make(map[common.Address]int)
	for index, node := range accs {
		oriAddrs = append(oriAddrs, node.Addr)
		AccIndex[node.Addr] = index
	}

	sort.Slice(oriAddrs, func(i, j int) bool {
		var flag1 bool
		for n := 0; n < 20; n++ {
			if oriAddrs[i][n] > oriAddrs[j][n] {
				flag1 = false
				break
			}
			if oriAddrs[i][n] < oriAddrs[j][n] {
				flag1 = true
				break
			}
			if oriAddrs[i][n] == oriAddrs[j][n] {
				continue
			}
		}
		return flag1
	})
	list := make([]types.Account, 0)
	for _, addr := range oriAddrs {
		list = append(list, accs[AccIndex[addr]])
	}

	return list

}

func NodesAddress(src []*types.Node) []common.Address {
	list := make([]common.Address, 0)
	for _, v := range src {
		list = append(list, v.Validator.Addr)
	}
	return list
}

func generateExtra(addrs []common.Address) string {
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)
	ist := &types.HotstuffExtra{
		StartHeight:   0,
		EndHeight:     types.GenesisBlockPerEpoch.Uint64(),
		Validators:    addrs,
		Seal:          make([]byte, types.HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
	}
	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return ""
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...))
}

func WriteGenesis(nodes types.Nodes, genesis string, Conf *config.Config) error {
	filepath := Conf.ChainPath
	for i := range nodes {
		if Conf.Machines[i].MachineIP == ip.IP {
			file, err := os.OpenFile(filepath+nodes[i].Name+"/setup/genesis.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}
			_, err = file.WriteString(genesis)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteStatic(nodes types.Nodes, static string, Conf *config.Config) error {
	filepath := Conf.ChainPath
	for i := range nodes {
		if Conf.Machines[i].MachineIP == ip.IP {
			file, err := os.OpenFile(filepath+nodes[i].Name+"/setup/static-nodes.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			_, err = file.WriteString(static)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func WriteKey(nodes types.Nodes, Conf *config.Config) error {
	filepath := Conf.ChainPath
	for i := range nodes {
		if Conf.Machines[i].MachineIP == ip.IP {
			file, err := os.OpenFile(filepath+nodes[i].Name+"/setup/node/nodekey", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}
			_, err = file.WriteString(nodes[i].Validator.NodeKey[2:])
			if err != nil {
				return err
			}
		}

	}
	for i := range nodes {
		if Conf.Machines[i].MachineIP == ip.IP {
			file, err := os.OpenFile(filepath+nodes[i].Name+"/setup/node/pubkey", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}
			_, err = file.WriteString(nodes[i].Validator.PubKey[2:])
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func WriteBash(nodes types.Nodes, Conf *config.Config) error { //写入脚本文件
	filepath := Conf.ChainPath
	for index, node := range nodes {
		if Conf.Machines[index].MachineIP == ip.IP {

			//写入build.sh
			build1 := "cd /data/gohome/src/Zion\n"
			build := "#!/bin/bash\n\nworkdir=$PWD\n" + build1 + "make geth\ncp build/bin/geth $workdir\n\ncd $workdir\nmd5sum geth\n"
			fileBuild, err := os.OpenFile(filepath+node.Name+"/build.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			_, err = fileBuild.WriteString(build)
			if err != nil {
				return err
			}

			//写入stop.sh
			stop := "#!/bin/bash\n        \nkill -s SIGINT $(ps aux|grep geth|grep node|grep -v grep|awk '{print $2}');\nps -ef|grep geth"
			fileStop, err := os.OpenFile(filepath+node.Name+"/stop.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			_, err = fileStop.WriteString(stop)
			if err != nil {
				return err
			}

			//写入init.sh
			init := "#!/bin/bash\n\nif [ ! -f node/genesis.json ]\nthen\nmkdir -p node/geth/\ncp setup/genesis.json node/\nfi\n\nif [ ! -f node/static-nodes.json ]\nthen\ncp setup/static-nodes.json node/\nfi\n\nif [ ! -f node/geth/nodekey ]\nthen\ncp setup/node/nodekey node/geth/\nfi\n\n./geth init node/genesis.json --datadir node\n"
			fileInit, err := os.OpenFile(filepath+node.Name+"/init.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			_, err = fileInit.WriteString(init)
			if err != nil {
				return err
			}

			//写入start.sh
			startP2P := fmt.Sprintf("startP2PPort=%s\n", Conf.Machines[index].P2PPort)
			startRPC := fmt.Sprintf("startRPCPort=%s\n", Conf.Machines[index].RPCPort)
			startChainID := fmt.Sprintf("chainID=%s\n", strconv.Itoa(int(Conf.ChainID)))
			startCoinbase := "coinbases=("
			for index, node := range nodes {
				if index != len(nodes)-1 {
					str := fmt.Sprintf("%s\t", node.Validator.Addr.String())
					startCoinbase = startCoinbase + str
				}
				if index == len(nodes)-1 {
					str := fmt.Sprintf("%s)\n", node.Validator.Addr.String())
					startCoinbase = startCoinbase + str
				}

			}
			startMiner := fmt.Sprintf("miner=%s\n", node.Validator.Addr.String())
			startEcho := "echo \"node" + strconv.Itoa(index) + " and miner is " + node.Validator.Addr.String() + "\""
			startGeth := "\nnohup ./geth --mine --miner.threads 1 \\\n--miner.etherbase=$miner \\\n--identity=node \\\n--maxpeers=100 \\\n--light.serve 20 \\\n--syncmode full \\\n--gcmode archive \\\n--allow-insecure-unlock \\\n--datadir node \\\n--networkid $chainID \\\n--http.api admin,eth,debug,miner,net,txpool,personal,web3 \\\n--http --http.addr 0.0.0.0 --http.port $startRPCPort --http.vhosts \"*\" \\\n--rpc.allow-unprotected-txs \\\n--nodiscover \\\n--port $startP2PPort \\\n--verbosity 5 >> node/node.log 2>&1 &\nsleep 1s\nps -ef|grep geth "
			start := "#!/bin/bash\n" + startP2P + startRPC + startChainID + startCoinbase + startMiner + startEcho + startGeth
			fileStart, err := os.OpenFile(filepath+node.Name+"/start.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			_, err = fileStart.WriteString(start)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
