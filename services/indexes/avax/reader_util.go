// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************

package avax

import (
	"math/big"

	"github.com/chain4travel/magellan/models"
)

func addAssetInfoMap(assets map[models.StringID]models.AssetInfo, assets2 map[models.StringID]models.AssetInfo) map[models.StringID]models.AssetInfo {
	addStringsFunc := func(t string, f string) string {
		tbi := new(big.Int)
		tbi.SetString(t, 10)
		fbi := new(big.Int)
		fbi.SetString(f, 10)
		rbi := new(big.Int)
		return rbi.Add(tbi, fbi).String()
	}

	addAssetInfos := func(t models.AssetInfo, f models.AssetInfo) models.AssetInfo {
		t.TransactionCount += f.TransactionCount
		t.UTXOCount += t.UTXOCount
		t.Balance = models.TokenAmount(addStringsFunc(string(t.Balance), string(f.Balance)))
		t.TotalReceived = models.TokenAmount(addStringsFunc(string(t.TotalReceived), string(f.TotalReceived)))
		t.TotalSent = models.TokenAmount(addStringsFunc(string(t.TotalSent), string(f.TotalSent)))
		return t
	}

	for k, v := range assets2 {
		if assetInfo, ok := assets[k]; ok {
			assets[k] = addAssetInfos(assetInfo, v)
		} else {
			assets[k] = v
		}
	}
	return assets
}
