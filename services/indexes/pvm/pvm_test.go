// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************

package pvm

import (
	"context"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/vms/platformvm/blocks"

	"github.com/ava-labs/avalanchego/vms/platformvm/txs"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/chain4travel/magellan/cfg"
	"github.com/chain4travel/magellan/db"
	"github.com/chain4travel/magellan/models"
	"github.com/chain4travel/magellan/services"
	"github.com/chain4travel/magellan/services/indexes/avax"
	"github.com/chain4travel/magellan/services/indexes/params"
	"github.com/chain4travel/magellan/servicesctrl"
	"github.com/chain4travel/magellan/utils"
)

var testXChainID = ids.ID([32]byte{7, 193, 50, 215, 59, 55, 159, 112, 106, 206, 236, 110, 229, 14, 139, 125, 14, 101, 138, 65, 208, 44, 163, 38, 115, 182, 177, 179, 244, 34, 195, 120})

func TestBootstrap(t *testing.T) {
	networkID := uint32(12345)
	conns, w, r, closeFn := newTestIndex(t, networkID, ChainID)
	defer closeFn()

	persist := db.NewPersist()

	genesis, err := utils.NewInternalGenesisContainer(networkID)
	if err != nil {
		t.Fatal("Failed to create genesis:", err.Error())
	}

	if err := w.Bootstrap(context.Background(), conns, persist, genesis); err != nil {
		t.Fatal(err)
	}

	txList, err := r.ListTransactions(context.Background(), &params.ListTransactionsParams{
		ChainIDs: []string{ChainID.String()},
	}, ids.Empty)
	if err != nil {
		t.Fatal("Failed to list transactions:", err.Error())
	}

	if txList == nil || *txList.Count < 1 {
		if txList.Count == nil {
			t.Fatal("Incorrect number of transactions:", txList.Count)
		} else {
			t.Fatal("Incorrect number of transactions:", *txList.Count)
		}
	}
}

func newTestIndex(t *testing.T, networkID uint32, chainID ids.ID) (*utils.Connections, *Writer, *avax.Reader, func()) {
	logConf := logging.Config{
		DisplayLevel: logging.Info,
		LogLevel:     logging.Debug,
	}
	conf := cfg.Services{
		Logging: logConf,
		DB: &cfg.DB{
			Driver: "mysql",
			DSN:    "root:password@tcp(127.0.0.1:3306)/magellan_test?parseTime=true",
		},
	}

	sc := &servicesctrl.Control{Log: logging.NoLog{}, Services: conf}
	conns, err := sc.Database()
	if err != nil {
		t.Fatal("Failed to create connections:", err.Error())
	}

	// Create index
	writer, err := NewWriter(networkID, chainID.String())
	if err != nil {
		t.Fatal("Failed to create writer:", err.Error())
	}

	cmap := make(map[string]services.Consumer)
	reader, _ := avax.NewReader(networkID, conns, cmap, sc)
	return conns, writer, reader, func() {
		_ = conns.Close()
	}
}

func TestInsertTxInternal(t *testing.T) {
	conns, writer, _, closeFn := newTestIndex(t, 5, testXChainID)
	defer closeFn()
	ctx := context.Background()

	tx := txs.Tx{}
	validatorTx := &txs.AddValidatorTx{}
	tx.Unsigned = validatorTx

	persist := db.NewPersistMock()
	session, _ := conns.DB().NewSession("pvm_test_tx", cfg.RequestTimeout)
	cCtx := services.NewConsumerContext(ctx, session, time.Now().Unix(), 0, persist, testXChainID.String())
	err := writer.indexTransaction(cCtx, tx.ID(), &tx, false)
	if err != nil {
		t.Fatal("insert failed", err)
	}
	if len(persist.Transactions) != 1 {
		t.Fatal("insert failed")
	}
	if len(persist.TransactionsBlock) != 1 {
		t.Fatal("insert failed")
	}
	if len(persist.TransactionsValidator) != 1 {
		t.Fatal("insert failed")
	}
}

func TestCommonBlock(t *testing.T) {
	conns, writer, _, closeFn := newTestIndex(t, 5, testXChainID)
	defer closeFn()
	ctx := context.Background()

	tx := blocks.CommonBlock{}
	blkid := ids.ID{}

	persist := db.NewPersistMock()
	session, _ := conns.DB().NewSession("pvm_test_tx", cfg.RequestTimeout)
	cCtx := services.NewConsumerContext(ctx, session, time.Now().Unix(), 0, persist, testXChainID.String())
	err := writer.indexCommonBlock(cCtx, blkid, models.BlockTypeCommit, tx, &models.BlockProposal{}, []byte(""))
	if err != nil {
		t.Fatal("insert failed", err)
	}
	if len(persist.PvmBlocks) != 1 {
		t.Fatal("insert failed")
	}
}
