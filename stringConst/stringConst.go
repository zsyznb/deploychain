package stringConst

import (
	"createChain/types"
	"encoding/json"
	"math/big"
)

func GenerateGenesis(chainID *big.Int, extra []byte, nodes types.Nodes) string {
	var rawGenesis types.Genesis
	rawString := "{\n    \"config\": {\n        \"chainId\": 60801,\n        \"homesteadBlock\": 0,\n        \"eip150Block\": 0,\n        \"eip155Block\": 0,\n        \"eip158Block\": 0,\n        \"byzantiumBlock\": 0,\n        \"constantinopleBlock\": 0,\n        \"petersburgBlock\": 0,\n        \"istanbulBlock\": 0,\n        \"berlinBlock\": 0,\n        \"londonBlock\": 0,\n        \"hotstuff\": {\n            \"protocol\": \"basic\"\n        }\n    },\n    \"alloc\": {\n        \"0x1eeDc8e7A4c708acF64106205F79CA4CDe11Ce3A\": {\"balance\": \"30000000000000000000000000\"},\n        \"0xb06Bf71465eA8071e83a87f6b914f8FFA6829f4b\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x328C78eb1b265F380381879e8F332fbB622EBAe5\": {\"balance\": \"10000000000000000000000000\"},\n        \"0xfEb056933BE960a183a0178249f7195AEbB74A4C\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x51d352512e592f7c50880a1fE280259fD68B07a2\": {\"balance\": \"10000000000000000000000000\"},\n        \"0xF6AE2D66B8a30104360a0188A7E8Da98eb336075\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x347Da040f1079C40AA76A84DeD08028D3425CC9D\": {\"balance\": \"10000000000000000000000000\"},\n        \"0x3d7ADdF663a30Ef02E086aE9eC37c311C892AFF5\": {\"balance\": \"10000000000000000000000000\"}\n    },\n    \"governance\": [],\n    \"community_rate\": 2000,\n    \"community_address\": \"0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D\",\n    \"coinbase\": \"0x0000000000000000000000000000000000000000\",\n    \"difficulty\": \"0x1\",\n    \"extraData\": \"\",\n    \"gasLimit\": \"0x1fffff\",\n    \"nonce\": \"0x4510809143055965\",\n    \"mixhash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"parentHash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"timestamp\": \"0x00\"\n}\n"
	rawData := []byte(rawString)
	err := rawGenesis.UnmarshalJSON(rawData)
	if err != nil {
		panic(err)
	}
	rawGenesis.Config.ChainID = chainID
	rawGenesis.ExtraData = extra
	for _, node := range nodes {
		rawAccount := &types.GovernanceAccount{
			Validator: node.Validator.Addr,
			Signer:    node.Signer.Addr,
		}
		rawGenesis.Governance = append(rawGenesis.Governance, *rawAccount)
	}
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
