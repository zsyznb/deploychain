package main

import (
	"bytes"
	"createChain/config"
	"crypto/ecdsa"
	"crypto/elliptic"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"os"
)

var (
	HotstuffExtraSeal    = 65
	HotstuffExtraVanity  = 32
	GenesisBlockPerEpoch = new(big.Int).SetUint64(30)
	nodeNum              int
	configPath           string
)

type Discv5NodeID [64]byte

type HotstuffExtra struct {
	StartHeight   uint64           // denote the epoch start height
	EndHeight     uint64           // the epoch end height
	Validators    []common.Address // consensus participants address for next epoch, and in the first block, it contains all genesis validators. keep empty if no epoch change.
	Seal          []byte           // proposer signature
	CommittedSeal [][]byte         // consensus participants signatures and it's size should be greater than 2/3 of validators
	Salt          []byte           // omit empty
}

type Node struct {
	Addr common.Address

	NodeKey string
	PubKey  string
	Static  string
}

type Nodes []*Node

func init() {
	flag.IntVar(&nodeNum, "nodenum", 4, "number of nodes")
	flag.StringVar(&configPath, "config", "config.json", "config path")
	flag.Parse()

}
func main() {
	config.LoadConfig(configPath, nodeNum)
	_, extra := generateNodes(config.Conf.NodeNum, config.Conf.Machines, config.Conf.P2PPort)
	log.Info(extra)
	makeDir(config.Conf.NodeNum)
}

func createBash(path string) *os.File { //创建脚本
	bash, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	return bash
}

func createFile(path string) *os.File { //创建文件
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	return file
}

func makeDir(n int) []*os.File { //创建文件目录
	dirs := make([]*os.File, 0)
	for i := 0; i < n; i++ {
		str := fmt.Sprintf("/node%d", i)
		err := os.Mkdir(config.Conf.ChainPath+str, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	log.Info("创建目录成功")
	return dirs
}

func generateNodes(n int, machines []string, port uint64) (Nodes, string) { //生成节点
	nodes := make([]*Node, 0)
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		nodekey := hexutil.Encode(crypto.FromECDSA(key))
		node := &Node{
			Addr:    crypto.PubkeyToAddress(key.PublicKey),
			NodeKey: nodekey,
			PubKey:  hexutil.Encode(crypto.CompressPubkey(&key.PublicKey)),
			Static:  fmt.Sprintf("enode://%s@%s:%v?discport=0", PubkeyID(&key.PublicKey), machines[i], port),
		}
		nodes = append(nodes, node)
	}
	log.Info("生成节点成功！")
	sortedNodes := sortNodes(nodes)
	log.Info("节点排序成功！")
	extra := generateExtra(NodesAddress(sortedNodes))
	log.Info("genesis extra创建成功！")
	return sortedNodes, extra
}

func PubkeyID(pub *ecdsa.PublicKey) Discv5NodeID {
	var id Discv5NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

func sortNodes(nodes Nodes) Nodes {
	oriAddrs := make([]common.Address, len(nodes))
	Nodeindex := make(map[common.Address]int)
	for index, node := range nodes {
		oriAddrs = append(oriAddrs, node.Addr)
		Nodeindex[node.Addr] = index
	}

	sort(oriAddrs)
	list := make([]*Node, 0)
	for _, addr := range oriAddrs {
		list = append(list, nodes[Nodeindex[addr]])
	}
	return list

}

func sort(addrs []common.Address) {
	for n := 0; n < len(addrs); n++ {
		for i := 0; i < len(addrs)-1; i++ {
			if Compare(addrs[i], addrs[i+1]) == 1 {
				Swap(addrs[i], addrs[i+1])
			}

		}
	}

}

//当addr1 > addr2时，返回值为1；
//当addr1 < addr2时，返回值为0；
func Compare(addr1 common.Address, addr2 common.Address) int {
	var flag1 int = 2
	length := len(addr1)
	for i := 0; i < int(length); i++ {
		if addr1[i] > addr2[i] {
			flag1 = 1
			break
		}
		if addr1[i] < addr2[i] {
			flag1 = 0
			break
		}
	}
	return flag1
}

func Swap(a interface{}, b interface{}) {
	var c interface{}
	c = a
	a = b
	b = c
}

func NodesAddress(src []*Node) []common.Address {
	list := make([]common.Address, 0)
	for _, v := range src {
		list = append(list, v.Addr)
	}
	return list
}

func generateExtra(addrs []common.Address) string {
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, HotstuffExtraVanity)...)
	ist := &HotstuffExtra{
		StartHeight:   0,
		EndHeight:     GenesisBlockPerEpoch.Uint64(),
		Validators:    addrs,
		Seal:          make([]byte, HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
	}
	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return ""
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...))
}
