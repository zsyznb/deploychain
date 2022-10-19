package types

import (
	"createChain/params"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

var (
	HotstuffExtraSeal    = 65
	HotstuffExtraVanity  = 32
	GenesisBlockPerEpoch = new(big.Int).SetUint64(400000)
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

type Account struct {
	Addr    common.Address
	NodeKey string
	PubKey  string
}
type Node struct {
	Name      string
	Validator Account
	Signer    Account
	Static    string
}

type Nodes []*Node

type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64              `json:"nonce"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash         `json:"mixHash"`
	Coinbase   common.Address      `json:"coinbase"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`
	Governance GenesisGovernance   `json:"governance" gencodec:"required"`
	// config of community pool
	CommunityRate    *big.Int       `json:"community_rate" gencodec:"required"`
	CommunityAddress common.Address `json:"community_address" gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}
type GenesisAlloc map[common.Address]GenesisAccount
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    string                      `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}
type GenesisGovernance []GovernanceAccount
type GovernanceAccount struct {
	Validator common.Address `json:"validator" gencodec:"required"`
	Signer    common.Address `json:"signer" gencodec:"required"`
}

func (n Discv5NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

func (g *Genesis) UnmarshalJSON(input []byte) error {
	type Genesis struct {
		Config           *params.ChainConfig                         `json:"config"`
		Nonce            *math.HexOrDecimal64                        `json:"nonce"`
		Timestamp        *math.HexOrDecimal64                        `json:"timestamp"`
		ExtraData        *hexutil.Bytes                              `json:"extraData"`
		GasLimit         *math.HexOrDecimal64                        `json:"gasLimit"   gencodec:"required"`
		Difficulty       *math.HexOrDecimal256                       `json:"difficulty" gencodec:"required"`
		Mixhash          *common.Hash                                `json:"mixHash"`
		Coinbase         *common.Address                             `json:"coinbase"`
		Alloc            map[common.UnprefixedAddress]GenesisAccount `json:"alloc"      gencodec:"required"`
		Governance       *GenesisGovernance                          `json:"governance" gencodec:"required"`
		CommunityRate    *big.Int                                    `json:"community_rate" gencodec:"required"`
		CommunityAddress *common.Address                             `json:"community_address" gencodec:"required"`
		Number           *math.HexOrDecimal64                        `json:"number"`
		GasUsed          *math.HexOrDecimal64                        `json:"gasUsed"`
		ParentHash       *common.Hash                                `json:"parentHash"`
	}
	var dec Genesis
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Config != nil {
		g.Config = dec.Config
	}
	if dec.Nonce != nil {
		g.Nonce = uint64(*dec.Nonce)
	}
	if dec.Timestamp != nil {
		g.Timestamp = uint64(*dec.Timestamp)
	}
	if dec.ExtraData != nil {
		g.ExtraData = *dec.ExtraData
	}
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Genesis")
	}
	g.GasLimit = uint64(*dec.GasLimit)
	if dec.Difficulty == nil {
		return errors.New("missing required field 'difficulty' for Genesis")
	}
	g.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.Mixhash != nil {
		g.Mixhash = *dec.Mixhash
	}
	if dec.Coinbase != nil {
		g.Coinbase = *dec.Coinbase
	}
	if dec.Alloc == nil {
		return errors.New("missing required field 'alloc' for Genesis")
	}
	g.Alloc = make(GenesisAlloc, len(dec.Alloc))
	for k, v := range dec.Alloc {
		g.Alloc[common.Address(k)] = v
	}
	if dec.Governance == nil {
		return errors.New("missing required field 'governance' for Genesis")
	}
	g.Governance = *dec.Governance
	if dec.CommunityRate == nil {
		return errors.New("missing required field 'community_rate' for Genesis")
	}
	g.CommunityRate = dec.CommunityRate
	if dec.CommunityAddress == nil {
		return errors.New("missing required field 'community_address' for Genesis")
	}
	g.CommunityAddress = *dec.CommunityAddress
	if dec.Number != nil {
		g.Number = uint64(*dec.Number)
	}
	if dec.GasUsed != nil {
		g.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.ParentHash != nil {
		g.ParentHash = *dec.ParentHash
	}
	return nil
}

func (g Genesis) MarshalJSON() ([]byte, error) {
	type Genesis struct {
		Config           *params.ChainConfig                         `json:"config"`
		Nonce            math.HexOrDecimal64                         `json:"nonce"`
		Timestamp        math.HexOrDecimal64                         `json:"timestamp"`
		ExtraData        hexutil.Bytes                               `json:"extraData"`
		GasLimit         math.HexOrDecimal64                         `json:"gasLimit"   gencodec:"required"`
		Difficulty       *math.HexOrDecimal256                       `json:"difficulty" gencodec:"required"`
		Mixhash          common.Hash                                 `json:"mixHash"`
		Coinbase         common.Address                              `json:"coinbase"`
		Alloc            map[common.UnprefixedAddress]GenesisAccount `json:"alloc"      gencodec:"required"`
		Governance       GenesisGovernance                           `json:"governance" gencodec:"required"`
		CommunityRate    *big.Int                                    `json:"community_rate" gencodec:"required"`
		CommunityAddress common.Address                              `json:"community_address" gencodec:"required"`
		Number           math.HexOrDecimal64                         `json:"number"`
		GasUsed          math.HexOrDecimal64                         `json:"gasUsed"`
		ParentHash       common.Hash                                 `json:"parentHash"`
	}
	var enc Genesis
	enc.Config = g.Config
	enc.Nonce = math.HexOrDecimal64(g.Nonce)
	enc.Timestamp = math.HexOrDecimal64(g.Timestamp)
	enc.ExtraData = g.ExtraData
	enc.GasLimit = math.HexOrDecimal64(g.GasLimit)
	enc.Difficulty = (*math.HexOrDecimal256)(g.Difficulty)
	enc.Mixhash = g.Mixhash
	enc.Coinbase = g.Coinbase
	if g.Alloc != nil {
		enc.Alloc = make(map[common.UnprefixedAddress]GenesisAccount, len(g.Alloc))
		for k, v := range g.Alloc {
			enc.Alloc[common.UnprefixedAddress(k)] = v
		}
	}
	enc.Governance = g.Governance
	enc.CommunityRate = g.CommunityRate
	enc.CommunityAddress = g.CommunityAddress
	enc.Number = math.HexOrDecimal64(g.Number)
	enc.GasUsed = math.HexOrDecimal64(g.GasUsed)
	enc.ParentHash = g.ParentHash
	return json.Marshal(&enc)
}
