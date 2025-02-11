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

package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/chain4travel/magellan/models"
	"github.com/gocraft/dbr/v2"
)

const (
	TestDB  = "mysql"
	TestDSN = "root:password@tcp(127.0.0.1:3306)/magellan_test?parseTime=true"
)

func TestTransaction(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &Transactions{}
	v.ID = "id"
	v.ChainID = "cid1"
	v.Type = "txtype"
	v.Memo = []byte("memo")
	v.CanonicalSerialization = []byte("cs")
	v.Txfee = 1
	v.Genesis = true
	v.CreatedAt = tm
	v.NetworkID = 1
	v.Status = 1

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableTransactions).Exec()

	err = p.InsertTransactions(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryTransactions(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}

	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.ChainID = "cid2"
	v.Type = "txtype1"
	v.Memo = []byte("memo1")
	v.CanonicalSerialization = []byte("cs1")
	v.Txfee = 2
	v.Genesis = false
	v.NetworkID = 2
	err = p.InsertTransactions(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryTransactions(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}

	if fv.NetworkID != 2 {
		t.Fatal("compare fail")
	}
	if fv.Txfee != 2 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputsRedeeming(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &OutputsRedeeming{}
	v.ID = "id1"
	v.RedeemedAt = tm
	v.RedeemingTransactionID = "rtxid"
	v.Amount = 100
	v.OutputIndex = 1
	v.Intx = "intx1"
	v.AssetID = "aid1"
	v.ChainID = "cid1"
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputsRedeeming).Exec()

	err = p.InsertOutputsRedeeming(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputsRedeeming(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.RedeemingTransactionID = "rtxid1"
	v.Amount = 102
	v.OutputIndex = 3
	v.Intx = "intx2"
	v.AssetID = "aid2"
	v.ChainID = "cid2"

	err = p.InsertOutputsRedeeming(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputsRedeeming(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Intx != "intx2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputs(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &Outputs{}
	v.ID = "id1"
	v.ChainID = "cid1"
	v.TransactionID = "txid1"
	v.OutputIndex = 1
	v.AssetID = "aid1"
	v.OutputType = models.OutputTypesSECP2556K1Transfer
	v.Amount = 2
	v.Locktime = 3
	v.Threshold = 4
	v.GroupID = 5
	v.Payload = []byte("payload")
	v.StakeLocktime = 6
	v.Stake = true
	v.Frozen = true
	v.Stakeableout = true
	v.Genesisutxo = true
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputs).Exec()

	err = p.InsertOutputs(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputs(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.ChainID = "cid2"
	v.TransactionID = "txid2"
	v.OutputIndex = 2
	v.AssetID = "aid2"
	v.OutputType = models.OutputTypesSECP2556K1Mint
	v.Amount = 3
	v.Locktime = 4
	v.Threshold = 5
	v.GroupID = 6
	v.Payload = []byte("payload2")
	v.StakeLocktime = 7
	v.Stake = false
	v.Frozen = false
	v.Stakeableout = false
	v.Genesisutxo = false

	err = p.InsertOutputs(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputs(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Amount != 3 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAssets(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &Assets{}
	v.ID = "id1"
	v.ChainID = "cid1"
	v.Name = "name1"
	v.Symbol = "symbol1"
	v.Denomination = 0x1
	v.Alias = "alias1"
	v.CurrentSupply = 1
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAssets).Exec()

	err = p.InsertAssets(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAssets(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.ChainID = "cid2"
	v.Name = "name2"
	v.Symbol = "symbol2"
	v.Denomination = 0x2
	v.Alias = "alias2"
	v.CurrentSupply = 2
	v.CreatedAt = tm

	err = p.InsertAssets(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryAssets(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Name != "name2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAddresses(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)
	tmu := time.Now().UTC().Truncate(1 * time.Second)

	basebin := [33]byte{}
	for cnt := 0; cnt < len(basebin); cnt++ {
		basebin[cnt] = byte(cnt + 1)
	}

	v := &Addresses{}
	v.Address = "id1"
	v.PublicKey = make([]byte, len(basebin))
	copy(v.PublicKey, basebin[:])
	v.CreatedAt = tm
	v.UpdatedAt = tmu

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAddresses).Exec()

	err = p.InsertAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	basebin[0] = 0xF
	basebin[5] = 0xE
	copy(v.PublicKey, basebin[:])
	v.CreatedAt = tm
	v.UpdatedAt = tmu.Add(1 * time.Minute)

	err = p.InsertAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.PublicKey[0] != 0xF {
		t.Fatal("compare fail")
	}
	if fv.PublicKey[5] != 0xE {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAddressChain(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)
	tmu := time.Now().UTC().Truncate(1 * time.Second)

	v := &AddressChain{}
	v.Address = "id1"
	v.ChainID = "ch1"
	v.CreatedAt = tm
	v.UpdatedAt = tmu

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAddressChain).Exec()

	err = p.InsertAddressChain(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAddressChain(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.ChainID = "ch2"
	v.CreatedAt = tm
	v.UpdatedAt = tmu.Add(1 * time.Minute)

	err = p.InsertAddressChain(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryAddressChain(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.ChainID != "ch2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputAddresses(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)
	tmu := time.Now().UTC().Truncate(1 * time.Second)

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputAddresses).Exec()

	v := &OutputAddresses{}
	v.OutputID = "oid1"
	v.Address = "id1"
	v.CreatedAt = tm
	v.UpdatedAt = tmu

	err = p.InsertOutputAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.RedeemingSignature != nil {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.OutputID = "oid1"
	v.Address = "id1"
	v.RedeemingSignature = []byte("rd1")
	v.CreatedAt = tm
	v.UpdatedAt = tmu.Add(1 * time.Minute)

	err = p.InsertOutputAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.UpdatedAt != tmu.Add(1*time.Minute) {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.RedeemingSignature = []byte("rd2")
	v.CreatedAt = tm
	v.UpdatedAt = tmu.Add(2 * time.Minute)

	err = p.InsertOutputAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if string(v.RedeemingSignature) != "rd2" {
		t.Fatal("compare fail")
	}
	if fv.UpdatedAt != tmu.Add(2*time.Minute) {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.RedeemingSignature = []byte("rd3")
	v.CreatedAt = tm
	v.UpdatedAt = tmu.Add(3 * time.Minute)

	err = p.UpdateOutputAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("update fail", err)
	}
	fv, err = p.QueryOutputAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if string(v.RedeemingSignature) != "rd3" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestTransactionsEpoch(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &TransactionsEpoch{}
	v.ID = "id1"
	v.Epoch = 10
	v.VertexID = "vid1"
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableTransactionsEpochs).Exec()

	err = p.InsertTransactionsEpoch(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryTransactionsEpoch(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Epoch = 11
	v.VertexID = "vid2"
	v.CreatedAt = tm

	err = p.InsertTransactionsEpoch(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryTransactionsEpoch(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.VertexID != "vid2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestCvmBlocks(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &CvmBlocks{}
	v.Block = "1"
	v.Serialization = []byte("{}")
	v.Hash = "0x"
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableCvmBlocks).Exec()

	err = p.InsertCvmBlocks(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryCvmBlock(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestCvmAddresses(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &CvmAddresses{}
	v.ID = "id1"
	v.Type = models.CChainIn
	v.Idx = 1
	v.TransactionID = "tid1"
	v.Address = "addr1"
	v.AssetID = "assid1"
	v.Amount = 2
	v.Nonce = 3
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableCvmAddresses).Exec()

	err = p.InsertCvmAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryCvmAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Type = models.CchainOut
	v.Idx = 2
	v.TransactionID = "tid2"
	v.Address = "addr2"
	v.AssetID = "assid2"
	v.Amount = 3
	v.Nonce = 4
	v.CreatedAt = tm

	err = p.InsertCvmAddresses(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryCvmAddresses(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Idx != 2 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestCvmTransactions(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second).Add(-1 * time.Hour)

	v := &CvmTransactionsAtomic{}
	v.TransactionID = "trid1"
	v.Type = models.CChainIn
	v.ChainID = "bid1"
	v.Block = "1"
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableCvmTransactionsAtomic).Exec()

	err = p.InsertCvmTransactionsAtomic(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryCvmTransactionsAtomic(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	tm = time.Now().UTC().Truncate(1 * time.Second).Add(-1 * time.Hour)

	v.Type = models.CchainOut
	v.TransactionID = "trid2"
	v.ChainID = "bid2"
	v.Block = "2"
	v.CreatedAt = tm

	err = p.InsertCvmTransactionsAtomic(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryCvmTransactionsAtomic(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !fv.CreatedAt.Equal(tm) {
		t.Fatal("compare fail")
	}
	if fv.TransactionID != "trid2" {
		t.Fatal("compare fail")
	}
	if fv.Block != "2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestCvmTransactionsTxdata(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &CvmTransactionsTxdata{}
	v.Hash = "h1"
	v.Block = "1"
	v.Idx = 1
	v.CreatedAt = tm
	v.Serialization = []byte("test123")

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableCvmTransactionsTxdata).Exec()

	err = p.InsertCvmAccount(ctx, rawDBConn.NewSession(stream), &CvmAccount{}, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	err = p.InsertCvmTransactionsTxdata(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryCvmTransactionsTxdata(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Idx = 7
	v.CreatedAt = tm
	v.Serialization = []byte("test456")

	err = p.InsertCvmTransactionsTxdata(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryCvmTransactionsTxdata(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if string(fv.Serialization) != "test456" {
		t.Fatal("compare fail")
	}
	if fv.Hash != "h1" {
		t.Fatal("compare fail")
	}
	if fv.Idx != 7 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestPvmBlocks(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &PvmBlocks{}
	v.ID = "id1"
	v.ChainID = "cid1"
	v.Type = models.BlockTypeAbort
	v.ParentID = "pid1"
	v.Serialization = []byte("ser1")
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TablePvmBlocks).Exec()

	err = p.InsertPvmBlocks(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryPvmBlocks(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.ChainID = "cid2"
	v.Type = models.BlockTypeCommit
	v.ParentID = "pid2"
	v.Serialization = []byte("ser2")
	v.CreatedAt = tm

	err = p.InsertPvmBlocks(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryPvmBlocks(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if string(fv.Serialization) != "ser2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestTransactionsValidator(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &TransactionsValidator{}
	v.ID = "id1"
	v.NodeID = "nid1"
	v.Start = 1
	v.End = 2
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableTransactionsValidator).Exec()

	err = p.InsertTransactionsValidator(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryTransactionsValidator(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.NodeID = "nid2"
	v.Start = 2
	v.End = 3
	v.CreatedAt = tm

	err = p.InsertTransactionsValidator(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryTransactionsValidator(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if v.NodeID != "nid2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestTransactionsBlock(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(1 * time.Second)

	v := &TransactionsBlock{}
	v.ID = "id1"
	v.TxBlockID = "txb1"
	v.CreatedAt = tm

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableTransactionsBlock).Exec()

	err = p.InsertTransactionsBlock(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryTransactionsBlock(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.TxBlockID = "txb2"
	v.CreatedAt = tm

	err = p.InsertTransactionsBlock(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryTransactionsBlock(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if v.TxBlockID != "txb2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAddressBech32(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	tmu := time.Now().UTC().Truncate(1 * time.Second)

	v := &AddressBech32{}
	v.Address = "adr1"
	v.Bech32Address = "badr1"
	v.UpdatedAt = tmu

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAddressBech32).Exec()

	err = p.InsertAddressBech32(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAddressBech32(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Bech32Address = "badr2"
	v.UpdatedAt = tmu.Add(1 * time.Minute)

	err = p.InsertAddressBech32(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryAddressBech32(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if v.Bech32Address != "badr2" {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputAddressAccumulateOut(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &OutputAddressAccumulate{}
	v.OutputID = "out1"
	v.Address = "adr1"
	v.TransactionID = "txid1"
	v.OutputIndex = 1
	v.CreatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputAddressAccumulateOut).Exec()

	err = p.InsertOutputAddressAccumulateOut(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputAddressAccumulateOut(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.OutputIndex = 3
	v.TransactionID = "tr3"

	err = p.InsertOutputAddressAccumulateOut(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputAddressAccumulateOut(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if v.OutputIndex != 3 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputAddressAccumulateIn(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &OutputAddressAccumulate{}
	v.OutputID = "out1"
	v.Address = "adr1"
	v.TransactionID = "txid1"
	v.OutputIndex = 1
	v.CreatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputAddressAccumulateIn).Exec()

	err = p.InsertOutputAddressAccumulateIn(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputAddressAccumulateIn(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.OutputIndex = 3
	v.TransactionID = "tr3"

	err = p.InsertOutputAddressAccumulateIn(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryOutputAddressAccumulateIn(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if v.OutputIndex != 3 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestOutputTxsAccumulate(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &OutputTxsAccumulate{}
	v.ChainID = "ch1"
	v.AssetID = "asset1"
	v.Address = "adr1"
	v.TransactionID = "tr1"
	v.CreatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableOutputTxsAccumulate).Exec()

	err = p.InsertOutputTxsAccumulate(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryOutputTxsAccumulate(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAccumulateBalancesReceived(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &AccumulateBalancesAmount{}
	v.ChainID = "ch1"
	v.AssetID = "asset1"
	v.Address = "adr1"
	v.TotalAmount = "0"
	v.UtxoCount = "0"
	v.UpdatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAccumulateBalancesReceived).Exec()

	err = p.InsertAccumulateBalancesReceived(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAccumulateBalancesReceived(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAccumulateBalancesSent(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &AccumulateBalancesAmount{}
	v.ChainID = "ch1"
	v.AssetID = "asset1"
	v.Address = "adr1"
	v.TotalAmount = "0"
	v.UtxoCount = "0"
	v.UpdatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAccumulateBalancesSent).Exec()

	err = p.InsertAccumulateBalancesSent(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAccumulateBalancesSent(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestAccumulateBalancesTransactions(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &AccumulateBalancesTransactions{}
	v.ChainID = "ch1"
	v.AssetID = "asset1"
	v.Address = "adr1"
	v.TransactionCount = "0"
	v.UpdatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableAccumulateBalancesTransactions).Exec()

	err = p.InsertAccumulateBalancesTransactions(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryAccumulateBalancesTransactions(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestTxPool(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &TxPool{}
	v.NetworkID = 1
	v.ChainID = "ch1"
	v.Serialization = []byte("hello")
	v.Topic = "topic1"
	v.MsgKey = "key1"
	v.CreatedAt = time.Now().UTC().Truncate(1 * time.Second)

	v.ComputeID()

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableTxPool).Exec()

	err = p.InsertTxPool(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryTxPool(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestKeyValueStore(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &KeyValueStore{}
	v.K = "k"
	v.V = "v"

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableKeyValueStore).Exec()

	err = p.InsertKeyValueStore(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryKeyValueStore(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestNodeIndex(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()

	v := &NodeIndex{}
	v.Instance = "def"
	v.Topic = "top"
	v.Idx = 1

	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableNodeIndex).Exec()

	err = p.InsertNodeIndex(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err := p.QueryNodeIndex(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Idx = 2

	err = p.InsertNodeIndex(ctx, rawDBConn.NewSession(stream), v, true)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryNodeIndex(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Idx != 2 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}

	v.Idx = 3
	err = p.UpdateNodeIndex(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	fv, err = p.QueryNodeIndex(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if fv.Idx != 3 {
		t.Fatal("compare fail")
	}
	if !reflect.DeepEqual(*v, *fv) {
		t.Fatal("compare fail")
	}
}

func TestInsertMultisigAlias(t *testing.T) {
	p := NewPersist()
	ctx := context.Background()
	stream := &dbr.NullEventReceiver{}

	rawDBConn, err := dbr.Open(TestDB, TestDSN, stream)
	if err != nil {
		t.Fatal("db fail", err)
	}
	_, _ = rawDBConn.NewSession(stream).DeleteFrom(TableMultisigAliases).Exec()

	v := &MultisigAlias{}
	v.Alias = "abcdefghijklmnopqrstABCDEF1234567"
	v.Owner = "ABCDEFghijklmnopqrstabcdef1234567"
	v.Memo = "Memo"
	v.Bech32Address = "kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvy"
	v.TransactionID = "abcdefghijklmnopqrstABCDEF1234567abcdefghijklmnop"
	v.CreatedAt = time.Now().UTC().Truncate(1 * time.Second)

	err = p.InsertMultisigAlias(ctx, rawDBConn.NewSession(stream), v)
	if err != nil {
		t.Fatal("insert fail", err)
	}
	err = p.InsertAddressBech32(ctx, rawDBConn.NewSession(stream), &AddressBech32{Address: v.Alias, Bech32Address: "kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvy", UpdatedAt: time.Now().UTC().Truncate(1 * time.Second)}, false)
	if err != nil {
		t.Fatal("insert address bech32 fail", err)
	}

	owners := []string{v.Owner}
	fv, err := p.QueryMultisigAliasesForOwners(ctx, rawDBConn.NewSession(stream), owners)
	if err != nil {
		t.Fatal("query fail", err)
	}
	if !reflect.DeepEqual(v.Bech32Address, (*fv)[0].Bech32Address) {
		t.Fatal("compare fail")
	}
	err = p.DeleteMultisigAlias(ctx, rawDBConn.NewSession(stream), v.Alias)
	if err != nil {
		t.Fatal("delete fail", err)
	}
}
