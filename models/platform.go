// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package models

import (
	"time"
)

// ChainInfo represents an overview of data about a given chain
type ChainInfo struct {
	ID          StringID `json:"chainID"`
	Alias       string   `json:"chainAlias"`
	VM          string   `json:"vm"`
	AVAXAssetID StringID `json:"avaxAssetID"`
	NetworkID   uint32   `json:"networkID"`
}

type Block struct {
	ID        StringID  `json:"id"`
	ParentID  StringID  `json:"parentID"`
	ChainID   StringID  `json:"chainID"`
	Type      BlockType `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type ControlSignature []byte

type BlockList struct {
	ListMetadata
	Blocks []*Block `json:"blocks"`
}

type BodyRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type CurrentValidatorsParams struct {
	SubnetID *string  `json:"subnetID"`
	NodeIDs  []string `json:"nodeIDs"`
}

type PeersParams struct {
	NodeIDs []string `json:"nodeIDs"`
}
