package stringConst

import (
	"createChain/config"
	"createChain/types"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func GenerateGenesis(Conf *config.Config, extra string, nodes types.Nodes) string {
	extra1 := common.Hex2Bytes(extra[2:])
	chainID := big.NewInt(int64(Conf.ChainID))
	var rawGenesis types.Genesis
	rawGenesisAlloc := make(map[common.Address]types.GenesisAccount, 8)
	rawString := "{\n    \"config\": {\n        \"chainId\": 60801,\n        \"homesteadBlock\": 0,\n        \"eip150Block\": 0,\n        \"eip155Block\": 0,\n        \"eip158Block\": 0,\n        \"byzantiumBlock\": 0,\n        \"constantinopleBlock\": 0,\n        \"petersburgBlock\": 0,\n        \"istanbulBlock\": 0,\n        \"berlinBlock\": 0,\n        \"londonBlock\": 0,\n        \"hotstuff\": {\n            \"protocol\": \"basic\"\n        }\n    },\n    \"alloc\": {\n        \"0x09f4E484D43B3D6b20957F7E1760beE3C6F62186\": {\"balance\": \"30000000000000000000000000\"},\n        \"0x294b8211E7010f457d85942aC874d076D739E32a\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x9deAD91D8632DCEEC701710bAF7922324DD45F58\": {\"balance\": \"10000000000000000000000000\"},\n        \"0xc5e2344b875e236b3475e9e4E70448525cA5210F\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x059178D8cD466c628b5F3291aCEDc97aa260104A\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x30876777cEc10593b2580C4A70FEfDF2f209E4D3\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x3B62008a8d5D3a8a2A539BbEe87D529d14dd6309\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x258af48e28E4A6846E931dDfF8e1Cdf8579821e5\": {\"balance\": \"10000000000000000000000000\"}\n    },\n    \"governance\": [],\n    \"community_rate\": 2000,\n    \"community_address\": \"0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D\",\n    \"coinbase\": \"0x0000000000000000000000000000000000000000\",\n    \"difficulty\": \"0x1\",\n    \"extraData\": \"\",\n    \"gasLimit\": \"0x1fffff\",\n    \"nonce\": \"0x4510809143055965\",\n    \"mixhash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"parentHash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"timestamp\": \"0x00\"\n}\n"
	rawData := []byte(rawString)
	err := rawGenesis.UnmarshalJSON(rawData)
	if err != nil {
		panic(err)
	}
	rawGenesis.Config.ChainID = chainID
	rawGenesis.ExtraData = extra1
	for _, node := range nodes {
		rawAccount := &types.GovernanceAccount{
			Validator: node.Validator.Addr,
			Signer:    node.Signer.Addr,
		}
		rawGenesis.Governance = append(rawGenesis.Governance, *rawAccount)
	}
	for i := 0; i < 20; i++ {
		var raw types.GenesisAccount
		raw = types.GenesisAccount{Balance: "5000000000000000000000000"}
		rawGenesisAlloc[config.Signers[i].Addr] = raw
	}
	rawGenesis.Alloc = rawGenesisAlloc
	s, err := rawGenesis.MarshalJSON()
	return string(s)

}

func GenerateStatic(nodes types.Nodes) string {

	staticNodes := make([]string, 0)
	for _, node := range nodes {
		staticNodes = append(staticNodes, node.Static)
	}
	staticNodesEnc, err := json.MarshalIndent(staticNodes, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(staticNodesEnc)

}
