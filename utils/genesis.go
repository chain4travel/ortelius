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

package utils

import (
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
)

type GenesisContainer struct {
	NetworkID       uint32
	Time            uint64
	XChainGenesisTx *txs.Tx
	XChainID        ids.ID
	AvaxAssetID     ids.ID
	GenesisBytes    []byte
}

func NewGenesisContainer(networkID uint32) (*GenesisContainer, error) {
	gc := &GenesisContainer{NetworkID: networkID}
	config := genesis.GetConfig(gc.NetworkID)
	var err error
	gc.GenesisBytes, gc.AvaxAssetID, err = genesis.FromConfig(config)
	if err != nil {
		return nil, err
	}

	gc.XChainGenesisTx, err = genesis.VMGenesis(gc.GenesisBytes, constants.AVMID)
	if err != nil {
		return nil, err
	}

	gc.XChainID = gc.XChainGenesisTx.ID()

	gc.Time = config.StartTime
	return gc, nil
}
