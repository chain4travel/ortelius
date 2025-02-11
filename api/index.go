// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code whose
// original notices appear below.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************
// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"encoding/json"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/chain4travel/magellan/models"
)

const (
	AVMName     = "avm"
	XChainAlias = "x"
	PVMName     = "pvm"
	PChainAlias = "p"
	CVMName     = "cvm"
	CChainAlias = "c"
)

func newIndexResponse(networkID uint32, xChainID, cChainID, avaxAssetID ids.ID) ([]byte, error) {
	return json.Marshal(&struct {
		NetworkID uint32                      `json:"network_id"`
		Chains    map[string]models.ChainInfo `json:"chains"`
	}{
		NetworkID: networkID,
		Chains: map[string]models.ChainInfo{
			xChainID.String(): {
				VM:          AVMName,
				Alias:       XChainAlias,
				NetworkID:   networkID,
				AVAXAssetID: models.StringID(avaxAssetID.String()),
				ID:          models.StringID(xChainID.String()),
			},
			ids.Empty.String(): {
				VM:          PVMName,
				Alias:       PChainAlias,
				NetworkID:   networkID,
				AVAXAssetID: models.StringID(avaxAssetID.String()),
				ID:          models.StringID(ids.Empty.String()),
			},
			cChainID.String(): {
				VM:          CVMName,
				Alias:       CChainAlias,
				NetworkID:   networkID,
				AVAXAssetID: models.StringID(avaxAssetID.String()),
				ID:          models.StringID(cChainID.String()),
			},
		},
	})
}

func newLegacyIndexResponse(networkID uint32, xChainID ids.ID, avaxAssetID ids.ID) ([]byte, error) {
	return json.Marshal(&models.ChainInfo{
		VM:          AVMName,
		NetworkID:   networkID,
		Alias:       XChainAlias,
		AVAXAssetID: models.StringID(avaxAssetID.String()),
		ID:          models.StringID(xChainID.String()),
	})
}
