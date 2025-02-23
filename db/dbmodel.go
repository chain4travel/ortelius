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
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/chain4travel/magellan/models"
	"github.com/chain4travel/magellan/utils"
	"github.com/gocraft/dbr/v2"
)

const (
	TableTransactions                   = "avm_transactions"
	TableOutputsRedeeming               = "avm_outputs_redeeming"
	TableOutputs                        = "avm_outputs"
	TableAssets                         = "avm_assets"
	TableAddresses                      = "addresses"
	TableAddressChain                   = "address_chain"
	TableOutputAddresses                = "avm_output_addresses"
	TableTransactionsEpochs             = "transactions_epoch"
	TableCvmAddresses                   = "cvm_addresses"
	TableCvmBlocks                      = "cvm_blocks"
	TableCvmTransactionsAtomic          = "cvm_transactions_atomic"
	TableCvmTransactionsTxdata          = "cvm_transactions_txdata"
	TableCvmAccounts                    = "cvm_accounts"
	TablePvmBlocks                      = "pvm_blocks"
	TableTransactionsValidator          = "transactions_validator"
	TableTransactionsBlock              = "transactions_block"
	TableAddressBech32                  = "addresses_bech32"
	TableOutputAddressAccumulateOut     = "output_addresses_accumulate_out"
	TableOutputAddressAccumulateIn      = "output_addresses_accumulate_in"
	TableOutputTxsAccumulate            = "output_txs_accumulate"
	TableAccumulateBalancesReceived     = "accumulate_balances_received"
	TableAccumulateBalancesSent         = "accumulate_balances_sent"
	TableAccumulateBalancesTransactions = "accumulate_balances_transactions"
	TableTxPool                         = "tx_pool"
	TableKeyValueStore                  = "key_value_store"
	TableNodeIndex                      = "node_index"
	TableCamLastBlockCache              = "cam_last_block_cache"
	TableMultisigAliases                = "multisig_aliases"
	TableReward                         = "reward"
	TableRewardOwner                    = "reward_owner"
)

type Persist interface {
	QueryTransactions(
		context.Context,
		dbr.SessionRunner,
		*Transactions,
	) (*Transactions, error)
	InsertTransactions(
		context.Context,
		dbr.SessionRunner,
		*Transactions,
		bool,
	) error

	QueryOutputsRedeeming(
		context.Context,
		dbr.SessionRunner,
		*OutputsRedeeming,
	) (*OutputsRedeeming, error)
	InsertOutputsRedeeming(
		context.Context,
		dbr.SessionRunner,
		*OutputsRedeeming,
		bool,
	) error

	QueryOutputs(
		context.Context,
		dbr.SessionRunner,
		*Outputs,
	) (*Outputs, error)
	InsertOutputs(
		context.Context,
		dbr.SessionRunner,
		*Outputs,
		bool,
	) error

	QueryAssets(
		context.Context,
		dbr.SessionRunner,
		*Assets,
	) (*Assets, error)
	InsertAssets(
		context.Context,
		dbr.SessionRunner,
		*Assets,
		bool,
	) error

	QueryAddresses(
		context.Context,
		dbr.SessionRunner,
		*Addresses,
	) (*Addresses, error)
	InsertAddresses(
		context.Context,
		dbr.SessionRunner,
		*Addresses,
		bool,
	) error

	QueryAddressChain(
		context.Context,
		dbr.SessionRunner,
		*AddressChain,
	) (*AddressChain, error)

	InsertAddressChain(
		context.Context,
		dbr.SessionRunner,
		*AddressChain,
		bool,
	) error

	QueryOutputAddresses(
		context.Context,
		dbr.SessionRunner,
		*OutputAddresses,
	) (*OutputAddresses, error)
	InsertOutputAddresses(
		context.Context,
		dbr.SessionRunner,
		*OutputAddresses,
		bool,
	) error
	UpdateOutputAddresses(
		context.Context,
		dbr.SessionRunner,
		*OutputAddresses,
	) error

	QueryTransactionsEpoch(
		context.Context,
		dbr.SessionRunner,
		*TransactionsEpoch,
	) (*TransactionsEpoch, error)
	InsertTransactionsEpoch(
		context.Context,
		dbr.SessionRunner,
		*TransactionsEpoch,
		bool,
	) error

	QueryCvmBlock(
		context.Context,
		dbr.SessionRunner,
		*CvmBlocks,
	) (*CvmBlocks, error)
	InsertCvmBlocks(
		context.Context,
		dbr.SessionRunner,
		*CvmBlocks,
	) error

	QueryCountLastBlockCache(
		context.Context,
		dbr.SessionRunner,
		*CamLastBlockCache,
	) (*CountLastBlockCache, error)
	QueryCamLastBlockCache(
		context.Context,
		dbr.SessionRunner,
		*CamLastBlockCache,
	) (*CamLastBlockCache, error)
	InsertCamLastBlockCache(
		context.Context,
		dbr.SessionRunner,
		*CamLastBlockCache,
		bool,
	) error

	QueryCvmAddresses(
		context.Context,
		dbr.SessionRunner,
		*CvmAddresses,
	) (*CvmAddresses, error)
	InsertCvmAddresses(
		context.Context,
		dbr.SessionRunner,
		*CvmAddresses,
		bool,
	) error

	QueryCvmTransactionsAtomic(
		context.Context,
		dbr.SessionRunner,
		*CvmTransactionsAtomic,
	) (*CvmTransactionsAtomic, error)
	InsertCvmTransactionsAtomic(
		context.Context,
		dbr.SessionRunner,
		*CvmTransactionsAtomic,
		bool,
	) error

	QueryCvmTransactionsTxdata(
		context.Context,
		dbr.SessionRunner,
		*CvmTransactionsTxdata,
	) (*CvmTransactionsTxdata, error)
	InsertCvmTransactionsTxdata(
		context.Context,
		dbr.SessionRunner,
		*CvmTransactionsTxdata,
		bool,
	) error

	QueryCvmAccount(
		ctx context.Context,
		sess dbr.SessionRunner,
		q *CvmAccount,
	) (*CvmAccount, error)

	InsertCvmAccount(
		ctx context.Context,
		sess dbr.SessionRunner,
		v *CvmAccount,
		upd bool,
	) error

	QueryPvmBlocks(
		context.Context,
		dbr.SessionRunner,
		*PvmBlocks,
	) (*PvmBlocks, error)
	InsertPvmBlocks(
		context.Context,
		dbr.SessionRunner,
		*PvmBlocks,
		bool,
	) error

	QueryTransactionsValidator(
		context.Context,
		dbr.SessionRunner,
		*TransactionsValidator,
	) (*TransactionsValidator, error)

	InsertTransactionsValidator(
		context.Context,
		dbr.SessionRunner,
		*TransactionsValidator,
		bool,
	) error

	QueryTransactionsBlock(
		context.Context,
		dbr.SessionRunner,
		*TransactionsBlock,
	) (*TransactionsBlock, error)

	InsertTransactionsBlock(
		context.Context,
		dbr.SessionRunner,
		*TransactionsBlock,
		bool,
	) error

	QueryAddressBech32(
		context.Context,
		dbr.SessionRunner,
		*AddressBech32,
	) (*AddressBech32, error)
	InsertAddressBech32(
		context.Context,
		dbr.SessionRunner,
		*AddressBech32,
		bool,
	) error

	QueryOutputAddressAccumulateOut(
		context.Context,
		dbr.SessionRunner,
		*OutputAddressAccumulate,
	) (*OutputAddressAccumulate, error)
	InsertOutputAddressAccumulateOut(
		context.Context,
		dbr.SessionRunner,
		*OutputAddressAccumulate,
		bool,
	) error

	QueryOutputAddressAccumulateIn(
		context.Context,
		dbr.SessionRunner,
		*OutputAddressAccumulate,
	) (*OutputAddressAccumulate, error)
	InsertOutputAddressAccumulateIn(
		context.Context,
		dbr.SessionRunner,
		*OutputAddressAccumulate,
		bool,
	) error
	UpdateOutputAddressAccumulateInOutputsProcessed(
		context.Context,
		dbr.SessionRunner,
		string,
	) error

	QueryOutputTxsAccumulate(
		context.Context,
		dbr.SessionRunner,
		*OutputTxsAccumulate,
	) (*OutputTxsAccumulate, error)
	InsertOutputTxsAccumulate(
		context.Context,
		dbr.SessionRunner,
		*OutputTxsAccumulate,
	) error

	QueryAccumulateBalancesReceived(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesAmount,
	) (*AccumulateBalancesAmount, error)
	InsertAccumulateBalancesReceived(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesAmount,
	) error

	QueryAccumulateBalancesSent(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesAmount,
	) (*AccumulateBalancesAmount, error)
	InsertAccumulateBalancesSent(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesAmount,
	) error

	QueryAccumulateBalancesTransactions(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesTransactions,
	) (*AccumulateBalancesTransactions, error)
	InsertAccumulateBalancesTransactions(
		context.Context,
		dbr.SessionRunner,
		*AccumulateBalancesTransactions,
	) error

	QueryTxPool(
		context.Context,
		dbr.SessionRunner,
		*TxPool,
	) (*TxPool, error)
	InsertTxPool(
		context.Context,
		dbr.SessionRunner,
		*TxPool,
	) error
	RemoveTxPool(
		context.Context,
		dbr.SessionRunner,
		*TxPool,
	) error

	QueryKeyValueStore(
		context.Context,
		dbr.SessionRunner,
		*KeyValueStore,
	) (*KeyValueStore, error)
	InsertKeyValueStore(
		context.Context,
		dbr.SessionRunner,
		*KeyValueStore,
	) error

	QueryNodeIndex(
		context.Context,
		dbr.SessionRunner,
		*NodeIndex,
	) (*NodeIndex, error)
	InsertNodeIndex(
		context.Context,
		dbr.SessionRunner,
		*NodeIndex,
		bool,
	) error
	UpdateNodeIndex(
		context.Context,
		dbr.SessionRunner,
		*NodeIndex,
	) error

	InsertMultisigAlias(
		context.Context,
		dbr.SessionRunner,
		*MultisigAlias,
	) error

	QueryMultisigAlias(
		context.Context,
		dbr.SessionRunner,
		string,
	) (*[]MultisigAlias, error)

	QueryMultisigAliasesForOwners(
		context.Context,
		dbr.SessionRunner,
		[]string,
	) (*[]MultisigAlias, error)

	DeleteMultisigAlias(
		context.Context,
		dbr.SessionRunner,
		string,
	) error

	InsertRewardOwner(
		context.Context,
		dbr.SessionRunner,
		*RewardOwner,
	) error

	InsertReward(
		context.Context,
		dbr.SessionRunner,
		*Reward,
	) error
}

type persist struct{}

func NewPersist() Persist {
	return &persist{}
}

func EventErr(t string, upd bool, err error) error {
	updmsg := ""
	if upd {
		updmsg = " upd"
	}
	return fmt.Errorf("%w (%s%s)", err, t, updmsg)
}

func PrintDbr(d *dbr.SelectStmt) *dbr.SelectStmt {
	buffer := dbr.NewBuffer()
	if err := d.Build(d.Dialect, buffer); err == nil {
		fmt.Println(buffer.String())
	}
	return d
}

type Transactions struct {
	ID                     string
	ChainID                string
	Type                   string
	Memo                   []byte
	CanonicalSerialization []byte
	Txfee                  uint64
	NetworkID              uint32
	Genesis                bool
	Status                 uint8
	CreatedAt              time.Time
}

func (p *persist) QueryTransactions(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *Transactions,
) (*Transactions, error) {
	v := &Transactions{}
	err := sess.Select(
		"id",
		"chain_id",
		"type",
		"memo",
		"created_at",
		"canonical_serialization",
		"txfee",
		"genesis",
		"network_id",
		"status",
	).From(TableTransactions).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertTransactions(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *Transactions,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableTransactions).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("type", v.Type).
		Pair("memo", v.Memo).
		Pair("created_at", v.CreatedAt).
		Pair("canonical_serialization", v.CanonicalSerialization).
		Pair("txfee", v.Txfee).
		Pair("genesis", v.Genesis).
		Pair("network_id", v.NetworkID).
		Pair("status", v.Status).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTransactions, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableTransactions).
			Set("chain_id", v.ChainID).
			Set("type", v.Type).
			Set("memo", v.Memo).
			Set("canonical_serialization", v.CanonicalSerialization).
			Set("txfee", v.Txfee).
			Set("genesis", v.Genesis).
			Set("network_id", v.NetworkID).
			Set("created_at", v.CreatedAt).
			Set("Status", v.Status).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableTransactions, true, err)
		}
	}
	return nil
}

type OutputsRedeeming struct {
	ID                     string
	RedeemedAt             time.Time
	RedeemingTransactionID string
	Amount                 uint64
	OutputIndex            uint32
	Intx                   string
	AssetID                string
	ChainID                string
	CreatedAt              time.Time
}

func (p *persist) QueryOutputsRedeeming(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *OutputsRedeeming,
) (*OutputsRedeeming, error) {
	v := &OutputsRedeeming{}
	err := sess.Select(
		"id",
		"redeemed_at",
		"redeeming_transaction_id",
		"amount",
		"output_index",
		"intx",
		"asset_id",
		"chain_id",
		"created_at",
	).From(TableOutputsRedeeming).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputsRedeeming(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputsRedeeming,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableOutputsRedeeming).
		Pair("id", v.ID).
		Pair("redeemed_at", v.RedeemedAt).
		Pair("redeeming_transaction_id", v.RedeemingTransactionID).
		Pair("amount", v.Amount).
		Pair("output_index", v.OutputIndex).
		Pair("intx", v.Intx).
		Pair("asset_id", v.AssetID).
		Pair("created_at", v.CreatedAt).
		Pair("chain_id", v.ChainID).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputsRedeeming, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableOutputsRedeeming).
			Set("redeeming_transaction_id", v.RedeemingTransactionID).
			Set("amount", v.Amount).
			Set("output_index", v.OutputIndex).
			Set("intx", v.Intx).
			Set("asset_id", v.AssetID).
			Set("chain_id", v.ChainID).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableOutputsRedeeming, true, err)
		}
	}
	return nil
}

type Outputs struct {
	ID            string
	ChainID       string
	TransactionID string
	OutputIndex   uint32
	OutputType    models.OutputType
	AssetID       string
	Amount        uint64
	Locktime      uint64
	Threshold     uint32
	GroupID       uint32
	Payload       []byte
	StakeLocktime uint64
	Stake         bool
	Frozen        bool
	Stakeableout  bool
	Genesisutxo   bool
	CreatedAt     time.Time
}

func (p *persist) QueryOutputs(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *Outputs,
) (*Outputs, error) {
	v := &Outputs{}
	err := sess.Select(
		"id",
		"chain_id",
		"transaction_id",
		"output_index",
		"asset_id",
		"output_type",
		"amount",
		"locktime",
		"threshold",
		"group_id",
		"payload",
		"stake_locktime",
		"stake",
		"frozen",
		"stakeableout",
		"genesisutxo",
		"created_at",
	).From(TableOutputs).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputs(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *Outputs,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableOutputs).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("transaction_id", v.TransactionID).
		Pair("output_index", v.OutputIndex).
		Pair("asset_id", v.AssetID).
		Pair("output_type", v.OutputType).
		Pair("amount", v.Amount).
		Pair("locktime", v.Locktime).
		Pair("threshold", v.Threshold).
		Pair("group_id", v.GroupID).
		Pair("payload", v.Payload).
		Pair("stake_locktime", v.StakeLocktime).
		Pair("stake", v.Stake).
		Pair("frozen", v.Frozen).
		Pair("stakeableout", v.Stakeableout).
		Pair("genesisutxo", v.Genesisutxo).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputs, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableOutputs).
			Set("chain_id", v.ChainID).
			Set("transaction_id", v.TransactionID).
			Set("output_index", v.OutputIndex).
			Set("asset_id", v.AssetID).
			Set("output_type", v.OutputType).
			Set("amount", v.Amount).
			Set("locktime", v.Locktime).
			Set("threshold", v.Threshold).
			Set("group_id", v.GroupID).
			Set("payload", v.Payload).
			Set("stake_locktime", v.StakeLocktime).
			Set("stake", v.Stake).
			Set("frozen", v.Frozen).
			Set("stakeableout", v.Stakeableout).
			Set("genesisutxo", v.Genesisutxo).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableOutputs, true, err)
		}
	}
	return nil
}

type Assets struct {
	ID            string
	ChainID       string
	Name          string
	Symbol        string
	Denomination  byte
	Alias         string
	CurrentSupply uint64
	CreatedAt     time.Time
}

func (p *persist) QueryAssets(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *Assets,
) (*Assets, error) {
	v := &Assets{}
	err := sess.Select(
		"id",
		"chain_id",
		"name",
		"symbol",
		"denomination",
		"alias",
		"current_supply",
		"created_at",
	).From(TableAssets).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAssets(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *Assets,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAssets).
		Pair("id", v.ID).
		Pair("chain_Id", v.ChainID).
		Pair("name", v.Name).
		Pair("symbol", v.Symbol).
		Pair("denomination", v.Denomination).
		Pair("alias", v.Alias).
		Pair("current_supply", v.CurrentSupply).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAssets, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableAssets).
			Set("chain_Id", v.ChainID).
			Set("name", v.Name).
			Set("symbol", v.Symbol).
			Set("denomination", v.Denomination).
			Set("alias", v.Alias).
			Set("current_supply", v.CurrentSupply).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAssets, true, err)
		}
	}
	return nil
}

type Addresses struct {
	Address   string
	PublicKey []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *persist) QueryAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *Addresses,
) (*Addresses, error) {
	v := &Addresses{}
	err := sess.Select(
		"address",
		"public_key",
		"created_at",
		"updated_at",
	).From(TableAddresses).
		Where("address=?", q.Address).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *Addresses,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAddresses).
		Pair("address", v.Address).
		Pair("public_key", v.PublicKey).
		Pair("created_at", v.CreatedAt).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAddresses, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableAddresses).
			Set("public_key", v.PublicKey).
			Set("updated_at", v.UpdatedAt).
			Set("created_at", v.CreatedAt).
			Where("address = ?", v.Address).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAddresses, true, err)
		}
	}

	return nil
}

type AddressChain struct {
	Address   string
	ChainID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *persist) QueryAddressChain(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *AddressChain,
) (*AddressChain, error) {
	v := &AddressChain{}
	err := sess.Select(
		"address",
		"chain_id",
		"created_at",
		"updated_at",
	).From(TableAddressChain).
		Where("address=? and chain_id=?", q.Address, q.ChainID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAddressChain(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *AddressChain,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAddressChain).
		Pair("address", v.Address).
		Pair("chain_id", v.ChainID).
		Pair("created_at", v.CreatedAt).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAddressChain, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableAddressChain).
			Set("updated_at", v.UpdatedAt).
			Set("created_at", v.CreatedAt).
			Where("address = ? and chain_id=?", v.Address, v.ChainID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAddressChain, true, err)
		}
	}
	return nil
}

type OutputAddresses struct {
	OutputID           string
	Address            string
	RedeemingSignature []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (p *persist) QueryOutputAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *OutputAddresses,
) (*OutputAddresses, error) {
	v := &OutputAddresses{}
	err := sess.Select(
		"output_id",
		"address",
		"redeeming_signature",
		"created_at",
		"updated_at",
	).From(TableOutputAddresses).
		Where("output_id=? and address=?", q.OutputID, q.Address).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputAddresses,
	upd bool,
) error {
	var err error
	stmt := sess.
		InsertInto(TableOutputAddresses).
		Pair("output_id", v.OutputID).
		Pair("address", v.Address).
		Pair("created_at", v.CreatedAt).
		Pair("updated_at", v.UpdatedAt)
	if v.RedeemingSignature != nil {
		stmt = stmt.Pair("redeeming_signature", v.RedeemingSignature)
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputAddresses, false, err)
	}
	if upd {
		stmt := sess.
			Update(TableOutputAddresses)
		if v.RedeemingSignature != nil {
			stmt = stmt.Set("redeeming_signature", v.RedeemingSignature)
		}
		_, err = stmt.
			Set("updated_at", v.UpdatedAt).
			Set("created_at", v.CreatedAt).
			Where("output_id = ? and address=?", v.OutputID, v.Address).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableOutputAddresses, true, err)
		}
	}
	return nil
}

func (p *persist) UpdateOutputAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputAddresses,
) error {
	var err error
	_, err = sess.
		Update(TableOutputAddresses).
		Set("redeeming_signature", v.RedeemingSignature).
		Set("updated_at", v.UpdatedAt).
		Where("output_id = ? and address=?", v.OutputID, v.Address).
		ExecContext(ctx)
	if err != nil {
		return EventErr(TableOutputAddresses, true, err)
	}
	return nil
}

type TransactionsEpoch struct {
	ID        string
	Epoch     uint32
	VertexID  string
	CreatedAt time.Time
}

func (p *persist) QueryTransactionsEpoch(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *TransactionsEpoch,
) (*TransactionsEpoch, error) {
	v := &TransactionsEpoch{}
	err := sess.Select(
		"id",
		"epoch",
		"vertex_id",
		"created_at",
	).From(TableTransactionsEpochs).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertTransactionsEpoch(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *TransactionsEpoch,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableTransactionsEpochs).
		Pair("id", v.ID).
		Pair("epoch", v.Epoch).
		Pair("vertex_id", v.VertexID).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTransactionsEpochs, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableTransactionsEpochs).
			Set("epoch", v.Epoch).
			Set("vertex_id", v.VertexID).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableTransactionsEpochs, true, err)
		}
	}

	return nil
}

type CvmBlocks struct {
	Block         string
	Hash          string
	ChainID       string
	EvmTx         int16
	AtomicTx      int16
	Serialization []byte
	CreatedAt     time.Time
	Proposer      string
	ProposerTime  *time.Time
	Size          float64
}

func (p *persist) QueryCvmBlock(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CvmBlocks,
) (*CvmBlocks, error) {
	v := &CvmBlocks{}
	err := sess.Select(
		"block",
		"hash",
		"chain_id",
		"evm_tx",
		"atomic_tx",
		"serialization",
		"created_at",
		"proposer",
		"proposer_time",
	).From(TableCvmBlocks).
		Where("block="+q.Block).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCvmBlocks(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CvmBlocks,
) error {
	var err error
	_, err = sess.
		InsertBySql("insert into "+TableCvmBlocks+" (block,hash,chain_id,evm_tx,atomic_tx,serialization,created_at, proposer,proposer_time,size) values("+v.Block+",?,?,?,?,?,?,?,?,?)",
			v.Hash, v.ChainID, v.EvmTx, v.AtomicTx, v.Serialization, v.CreatedAt, v.Proposer, v.ProposerTime, v.Size).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableCvmBlocks, false, err)
	}
	return nil
}

type CamLastBlockCache struct {
	CurrentBlock string
	ChainID      string
}

type CountLastBlockCache struct {
	Cnt uint64
}

func (p *persist) QueryCountLastBlockCache(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CamLastBlockCache,
) (*CountLastBlockCache, error) {
	v := &CountLastBlockCache{}
	err := sess.Select(
		"count(*) as cnt",
	).From(TableCamLastBlockCache).
		Where("chainid=?", q.ChainID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) QueryCamLastBlockCache(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CamLastBlockCache,
) (*CamLastBlockCache, error) {
	v := &CamLastBlockCache{}
	err := sess.Select(
		"current_block",
		"chainid",
	).From(TableCamLastBlockCache).
		Where("chainid=?", q.ChainID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCamLastBlockCache(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CamLastBlockCache,
	upd bool,
) error {
	var err error
	if upd {
		_, err = sess.
			Update(TableCamLastBlockCache).
			Set("current_block", v.CurrentBlock).
			Where("chainid = ?", v.ChainID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableCamLastBlockCache, true, err)
		}
	} else {
		_, err = sess.
			InsertBySql("insert into "+TableCamLastBlockCache+" (current_block,chainid) values("+v.CurrentBlock+",?)",
				v.ChainID).
			ExecContext(ctx)
		if err != nil && !utils.ErrIsDuplicateEntryError(err) {
			return EventErr(TableCamLastBlockCache, false, err)
		}
	}
	return nil
}

type CvmAddresses struct {
	ID            string
	Type          models.CChainType
	Idx           uint64
	TransactionID string
	Address       string
	AssetID       string
	Amount        uint64
	Nonce         uint64
	CreatedAt     time.Time
}

func (p *persist) QueryCvmAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CvmAddresses,
) (*CvmAddresses, error) {
	v := &CvmAddresses{}
	err := sess.Select(
		"id",
		"type",
		"idx",
		"transaction_id",
		"address",
		"asset_id",
		"amount",
		"nonce",
		"created_at",
	).From(TableCvmAddresses).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCvmAddresses(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CvmAddresses,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableCvmAddresses).
		Pair("id", v.ID).
		Pair("type", v.Type).
		Pair("idx", v.Idx).
		Pair("transaction_id", v.TransactionID).
		Pair("address", v.Address).
		Pair("asset_id", v.AssetID).
		Pair("amount", v.Amount).
		Pair("nonce", v.Nonce).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableCvmAddresses, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableCvmAddresses).
			Set("type", v.Type).
			Set("idx", v.Idx).
			Set("transaction_id", v.TransactionID).
			Set("address", v.Address).
			Set("asset_id", v.AssetID).
			Set("amount", v.Amount).
			Set("nonce", v.Nonce).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableCvmAddresses, true, err)
		}
	}
	return nil
}

type CvmTransactionsAtomic struct {
	TransactionID string
	Block         string
	ChainID       string
	Type          models.CChainType
	CreatedAt     time.Time
}

func (p *persist) QueryCvmTransactionsAtomic(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CvmTransactionsAtomic,
) (*CvmTransactionsAtomic, error) {
	v := &CvmTransactionsAtomic{}
	err := sess.Select(
		"transaction_id",
		"cast(block as char) as block",
		"chain_id",
		"type",
		"created_at",
	).From(TableCvmTransactionsAtomic).
		Where("transaction_id=?", q.TransactionID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCvmTransactionsAtomic(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CvmTransactionsAtomic,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertBySql("insert into "+TableCvmTransactionsAtomic+" (transaction_id,block,chain_id,type,created_at) values(?,"+v.Block+",?,?,?)",
			v.TransactionID, v.ChainID, v.Type, v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableCvmTransactionsAtomic, false, err)
	}
	if upd {
		_, err = sess.
			UpdateBySql("update "+TableCvmTransactionsAtomic+" set block="+v.Block+",chain_id=?,type=?,created_at=? where transaction_id=?",
				v.ChainID, v.Type, v.CreatedAt, v.TransactionID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableCvmTransactionsAtomic, true, err)
		}
	}
	return nil
}

type CvmTransactionsTxdata struct {
	Hash          string
	Block         string
	Idx           uint64
	FromAddr      string
	ToAddr        string
	Nonce         uint64
	Amount        uint64
	Status        uint16
	GasUsed       uint64
	GasPrice      uint64
	Serialization []byte
	Receipt       []byte
	CreatedAt     time.Time
}

func (p *persist) QueryCvmTransactionsTxdata(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CvmTransactionsTxdata,
) (*CvmTransactionsTxdata, error) {
	v := &CvmTransactionsTxdata{}
	err := sess.Select(
		"hash",
		"block",
		"idx",
		"F.address",
		"T.address",
		"nonce",
		"amount",
		"status",
		"gas_used",
		"gas_price",
		"serialization",
		"receipt",
		"created_at",
	).From(TableCvmTransactionsTxdata).
		Join(dbr.I(TableCvmAccounts).As("F"), "F.id=id_from_addr").
		Join(dbr.I(TableCvmAccounts).As("T"), "T.id=id_to_addr").
		Where("hash=?", q.Hash).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCvmTransactionsTxdata(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CvmTransactionsTxdata,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableCvmTransactionsTxdata).
		Pair("hash", v.Hash).
		Pair("block", v.Block).
		Pair("idx", v.Idx).
		Pair("id_from_addr", dbr.Select("id").From(TableCvmAccounts).Where("address=?", v.FromAddr)).
		Pair("id_to_addr", dbr.Select("id").From(TableCvmAccounts).Where("address=?", v.ToAddr)).
		Pair("nonce", v.Nonce).
		Pair("amount", v.Amount).
		Pair("status", v.Status).
		Pair("gas_used", v.GasUsed).
		Pair("gas_price", v.GasPrice).
		Pair("serialization", v.Serialization).
		Pair("receipt", v.Receipt).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableCvmTransactionsTxdata, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableCvmTransactionsTxdata).
			Set("block", v.Block).
			Set("idx", v.Idx).
			Set("id_from_addr", dbr.Select("id").From(TableCvmAccounts).Where("address=?", v.FromAddr)).
			Set("id_to_addr", dbr.Select("id").From(TableCvmAccounts).Where("address=?", v.ToAddr)).
			Set("nonce", v.Nonce).
			Set("amount", v.Amount).
			Set("status", v.Status).
			Set("gas_used", v.GasUsed).
			Set("gas_price", v.GasPrice).
			Set("serialization", v.Serialization).
			Set("receipt", v.Receipt).
			Set("created_at", v.CreatedAt).
			Where("hash=?", v.Hash).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableCvmTransactionsTxdata, true, err)
		}
	}
	return nil
}

type CvmAccount struct {
	ID         uint64
	Address    string
	TxCount    uint64
	CreationTx *string
}

func (p *persist) QueryCvmAccount(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *CvmAccount,
) (*CvmAccount, error) {
	v := &CvmAccount{}
	err := sess.Select(
		"id",
		"tx_count",
		"creation_tx",
	).From(TableCvmAccounts).
		Where("address=?", q.Address).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertCvmAccount(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *CvmAccount,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableCvmAccounts).
		Pair("address", v.Address).
		Pair("tx_count", v.TxCount).
		Pair("creation_tx", v.CreationTx).
		ExecContext(ctx)
	if err == nil {
		return nil
	} else if !upd || !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableCvmAccounts, false, err)
	}
	_, err = sess.
		Update(TableCvmAccounts).
		IncrBy("tx_count", v.TxCount).
		Where("address=?", v.Address).
		ExecContext(ctx)
	if err != nil {
		return EventErr(TableCvmAccounts, true, err)
	}
	return nil
}

type PvmBlocks struct {
	ID            string
	ChainID       string
	Type          models.BlockType
	ParentID      string
	Serialization []byte
	CreatedAt     time.Time
	Height        uint64
	Proposer      string
	ProposerTime  *time.Time
}

func (p *persist) QueryPvmBlocks(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *PvmBlocks,
) (*PvmBlocks, error) {
	v := &PvmBlocks{}
	err := sess.Select(
		"id",
		"chain_id",
		"type",
		"parent_id",
		"serialization",
		"created_at",
		"height",
		"proposer",
		"proposer_time",
	).From(TablePvmBlocks).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertPvmBlocks(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *PvmBlocks,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TablePvmBlocks).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("type", v.Type).
		Pair("parent_id", v.ParentID).
		Pair("created_at", v.CreatedAt).
		Pair("serialization", v.Serialization).
		Pair("height", v.Height).
		Pair("proposer", v.Proposer).
		Pair("proposer_time", v.ProposerTime).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TablePvmBlocks, false, err)
	}
	if upd {
		_, err = sess.
			Update(TablePvmBlocks).
			Set("chain_id", v.ChainID).
			Set("type", v.Type).
			Set("parent_id", v.ParentID).
			Set("serialization", v.Serialization).
			Set("height", v.Height).
			Set("created_at", v.CreatedAt).
			Set("proposer", v.Proposer).
			Set("proposer_time", v.ProposerTime).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TablePvmBlocks, true, err)
		}
	}

	return nil
}

type TransactionsValidator struct {
	ID        string
	NodeID    string
	Start     uint64
	End       uint64
	CreatedAt time.Time
}

func (p *persist) QueryTransactionsValidator(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *TransactionsValidator,
) (*TransactionsValidator, error) {
	v := &TransactionsValidator{}
	err := sess.Select(
		"id",
		"node_id",
		"start",
		"end",
		"created_at",
	).From(TableTransactionsValidator).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertTransactionsValidator(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *TransactionsValidator,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableTransactionsValidator).
		Pair("id", v.ID).
		Pair("node_id", v.NodeID).
		Pair("start", v.Start).
		Pair("end", v.End).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTransactionsValidator, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableTransactionsValidator).
			Set("node_id", v.NodeID).
			Set("start", v.Start).
			Set("end", v.End).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableTransactionsValidator, true, err)
		}
	}
	return nil
}

type TransactionsBlock struct {
	ID        string
	TxBlockID string
	CreatedAt time.Time
}

func (p *persist) QueryTransactionsBlock(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *TransactionsBlock,
) (*TransactionsBlock, error) {
	v := &TransactionsBlock{}
	err := sess.Select(
		"id",
		"tx_block_id",
		"created_at",
	).From(TableTransactionsBlock).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertTransactionsBlock(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *TransactionsBlock,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableTransactionsBlock).
		Pair("id", v.ID).
		Pair("tx_block_id", v.TxBlockID).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTransactionsBlock, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableTransactionsBlock).
			Set("tx_block_id", v.TxBlockID).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableTransactionsBlock, true, err)
		}
	}
	return nil
}

type AddressBech32 struct {
	Address       string
	Bech32Address string
	UpdatedAt     time.Time
}

func (p *persist) QueryAddressBech32(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *AddressBech32,
) (*AddressBech32, error) {
	v := &AddressBech32{}
	err := sess.Select(
		"address",
		"bech32_address",
		"updated_at",
	).From(TableAddressBech32).
		Where("address=?", q.Address).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAddressBech32(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *AddressBech32,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAddressBech32).
		Pair("address", v.Address).
		Pair("bech32_address", v.Bech32Address).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAddressBech32, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableAddressBech32).
			Set("bech32_address", v.Bech32Address).
			Set("updated_at", v.UpdatedAt).
			Where("address = ?", v.Address).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAddressBech32, true, err)
		}
	}
	return nil
}

type OutputAddressAccumulate struct {
	ID              string
	OutputID        string
	Address         string
	Processed       int
	OutputProcessed int
	TransactionID   string
	OutputIndex     uint32
	CreatedAt       time.Time
}

func (b *OutputAddressAccumulate) ComputeID() {
	idsv := fmt.Sprintf("%s:%s", b.OutputID, b.Address)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

func (p *persist) QueryOutputAddressAccumulateOut(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *OutputAddressAccumulate,
) (*OutputAddressAccumulate, error) {
	v := &OutputAddressAccumulate{}
	err := sess.Select(
		"id",
		"output_id",
		"address",
		"processed",
		"transaction_id",
		"output_index",
		"created_at",
	).From(TableOutputAddressAccumulateOut).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputAddressAccumulateOut(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputAddressAccumulate,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableOutputAddressAccumulateOut).
		Pair("id", v.ID).
		Pair("output_id", v.OutputID).
		Pair("address", v.Address).
		Pair("transaction_id", v.TransactionID).
		Pair("output_index", v.OutputIndex).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputAddressAccumulateOut, false, err)
	}

	if upd {
		_, err = sess.
			Update(TableOutputAddressAccumulateOut).
			Set("output_id", v.OutputID).
			Set("address", v.Address).
			Set("transaction_id", v.TransactionID).
			Set("output_index", v.OutputIndex).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAddressBech32, true, err)
		}
	}

	return nil
}

func (p *persist) QueryOutputAddressAccumulateIn(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *OutputAddressAccumulate,
) (*OutputAddressAccumulate, error) {
	v := &OutputAddressAccumulate{}
	err := sess.Select(
		"id",
		"output_id",
		"address",
		"processed",
		"transaction_id",
		"output_index",
		"created_at",
	).From(TableOutputAddressAccumulateIn).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputAddressAccumulateIn(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputAddressAccumulate,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableOutputAddressAccumulateIn).
		Pair("id", v.ID).
		Pair("output_id", v.OutputID).
		Pair("address", v.Address).
		Pair("transaction_id", v.TransactionID).
		Pair("output_index", v.OutputIndex).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputAddressAccumulateIn, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableOutputAddressAccumulateIn).
			Set("output_id", v.OutputID).
			Set("address", v.Address).
			Set("transaction_id", v.TransactionID).
			Set("output_index", v.OutputIndex).
			Set("created_at", v.CreatedAt).
			Where("id = ?", v.ID).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableAddressBech32, true, err)
		}
	}
	return nil
}

func (p *persist) UpdateOutputAddressAccumulateInOutputsProcessed(
	ctx context.Context,
	sess dbr.SessionRunner,
	id string,
) error {
	var err error
	_, err = sess.
		Update(TableOutputAddressAccumulateIn).
		Set("output_processed", 1).
		Where("output_id=? and output_processed <> ?", id, 1).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputAddressAccumulateIn, false, err)
	}

	return nil
}

type OutputTxsAccumulate struct {
	ID            string
	ChainID       string
	AssetID       string
	Address       string
	TransactionID string
	Processed     int
	CreatedAt     time.Time
}

func (b *OutputTxsAccumulate) ComputeID() {
	idsv := fmt.Sprintf("%s:%s:%s:%s", b.ChainID, b.AssetID, b.Address, b.TransactionID)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

func (p *persist) QueryOutputTxsAccumulate(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *OutputTxsAccumulate,
) (*OutputTxsAccumulate, error) {
	v := &OutputTxsAccumulate{}
	err := sess.Select(
		"id",
		"chain_id",
		"asset_id",
		"address",
		"transaction_id",
		"processed",
		"created_at",
	).From(TableOutputTxsAccumulate).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertOutputTxsAccumulate(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *OutputTxsAccumulate,
) error {
	var err error
	_, err = sess.
		InsertInto(TableOutputTxsAccumulate).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("asset_id", v.AssetID).
		Pair("address", v.Address).
		Pair("transaction_id", v.TransactionID).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableOutputTxsAccumulate, false, err)
	}

	return nil
}

type AccumulateBalancesAmount struct {
	ID          string
	ChainID     string
	AssetID     string
	Address     string
	TotalAmount string
	UtxoCount   string
	UpdatedAt   time.Time
}

func (b *AccumulateBalancesAmount) ComputeID() {
	idsv := fmt.Sprintf("%s:%s:%s", b.ChainID, b.AssetID, b.Address)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

func (p *persist) QueryAccumulateBalancesReceived(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *AccumulateBalancesAmount,
) (*AccumulateBalancesAmount, error) {
	v := &AccumulateBalancesAmount{}
	err := sess.Select(
		"id",
		"chain_id",
		"asset_id",
		"address",
		"cast(total_amount as char) total_amount",
		"cast(utxo_count as char) utxo_count",
		"updated_at",
	).From(TableAccumulateBalancesReceived).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAccumulateBalancesReceived(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *AccumulateBalancesAmount,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAccumulateBalancesReceived).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("asset_id", v.AssetID).
		Pair("address", v.Address).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAccumulateBalancesReceived, false, err)
	}

	return nil
}

func (p *persist) QueryAccumulateBalancesSent(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *AccumulateBalancesAmount,
) (*AccumulateBalancesAmount, error) {
	v := &AccumulateBalancesAmount{}
	err := sess.Select(
		"id",
		"chain_id",
		"asset_id",
		"address",
		"cast(total_amount as char) total_amount",
		"cast(utxo_count as char) utxo_count",
		"updated_at",
	).From(TableAccumulateBalancesSent).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAccumulateBalancesSent(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *AccumulateBalancesAmount,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAccumulateBalancesSent).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("asset_id", v.AssetID).
		Pair("address", v.Address).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAccumulateBalancesSent, false, err)
	}

	return nil
}

type AccumulateBalancesTransactions struct {
	ID               string
	ChainID          string
	AssetID          string
	Address          string
	TransactionCount string
	UpdatedAt        time.Time
}

func (b *AccumulateBalancesTransactions) ComputeID() {
	idsv := fmt.Sprintf("%s:%s:%s", b.ChainID, b.AssetID, b.Address)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

func (p *persist) QueryAccumulateBalancesTransactions(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *AccumulateBalancesTransactions,
) (*AccumulateBalancesTransactions, error) {
	v := &AccumulateBalancesTransactions{}
	err := sess.Select(
		"id",
		"chain_id",
		"asset_id",
		"address",
		"cast(transaction_count as char) transaction_count",
		"updated_at",
	).From(TableAccumulateBalancesTransactions).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertAccumulateBalancesTransactions(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *AccumulateBalancesTransactions,
) error {
	var err error
	_, err = sess.
		InsertInto(TableAccumulateBalancesTransactions).
		Pair("id", v.ID).
		Pair("chain_id", v.ChainID).
		Pair("asset_id", v.AssetID).
		Pair("address", v.Address).
		Pair("updated_at", v.UpdatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableAccumulateBalancesTransactions, false, err)
	}

	return nil
}

type TxPool struct {
	ID            string
	NetworkID     uint32
	ChainID       string
	MsgKey        string
	Serialization []byte
	Topic         string
	CreatedAt     time.Time
}

func (b *TxPool) ComputeID() {
	idsv := fmt.Sprintf("%s:%s", b.MsgKey, b.Topic)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

func (p *persist) QueryTxPool(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *TxPool,
) (*TxPool, error) {
	v := &TxPool{}
	err := sess.Select(
		"id",
		"network_id",
		"chain_id",
		"msg_key",
		"serialization",
		"topic",
		"created_at",
	).From(TableTxPool).
		Where("id=?", q.ID).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertTxPool(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *TxPool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableTxPool).
		Pair("id", v.ID).
		Pair("network_id", v.NetworkID).
		Pair("chain_id", v.ChainID).
		Pair("msg_key", v.MsgKey).
		Pair("serialization", v.Serialization).
		Pair("topic", v.Topic).
		Pair("created_at", v.CreatedAt).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTxPool, false, err)
	}

	return nil
}

func (p *persist) RemoveTxPool(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *TxPool,
) error {
	var err error
	_, err = sess.
		DeleteFrom(TableTxPool).
		Where("id=?", v.ID).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableTxPool, false, err)
	}

	return nil
}

type KeyValueStore struct {
	K string
	V string
}

func (p *persist) QueryKeyValueStore(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *KeyValueStore,
) (*KeyValueStore, error) {
	v := &KeyValueStore{}
	err := sess.Select(
		"k",
		"v",
	).From(TableKeyValueStore).
		Where("k=?", q.K).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertKeyValueStore(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *KeyValueStore,
) error {
	var err error
	_, err = sess.
		InsertInto(TableKeyValueStore).
		Pair("k", v.K).
		Pair("v", v.V).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableKeyValueStore, false, err)
	}

	return nil
}

type NodeIndex struct {
	Instance string
	Topic    string
	Idx      uint64
}

func (p *persist) QueryNodeIndex(
	ctx context.Context,
	sess dbr.SessionRunner,
	q *NodeIndex,
) (*NodeIndex, error) {
	v := &NodeIndex{}
	err := sess.Select(
		"instance",
		"topic",
		"idx",
	).From(TableNodeIndex).
		Where("instance=? and topic=?", q.Instance, q.Topic).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) InsertNodeIndex(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *NodeIndex,
	upd bool,
) error {
	var err error
	_, err = sess.
		InsertInto(TableNodeIndex).
		Pair("instance", v.Instance).
		Pair("topic", v.Topic).
		Pair("idx", v.Idx).
		ExecContext(ctx)
	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableNodeIndex, false, err)
	}
	if upd {
		_, err = sess.
			Update(TableNodeIndex).
			Set("idx", v.Idx).
			Where("instance=? and topic=?", v.Instance, v.Topic).
			ExecContext(ctx)
		if err != nil {
			return EventErr(TableNodeIndex, true, err)
		}
	}
	return nil
}

func (p *persist) UpdateNodeIndex(
	ctx context.Context,
	sess dbr.SessionRunner,
	v *NodeIndex,
) error {
	var err error
	_, err = sess.
		Update(TableNodeIndex).
		Set("idx", v.Idx).
		Where("instance=? and topic=?", v.Instance, v.Topic).
		ExecContext(ctx)
	if err != nil {
		return EventErr(TableNodeIndex, true, err)
	}
	return nil
}

type CvmLogs struct {
	ID            string
	BlockHash     string
	TxHash        string
	LogIndex      uint64
	FirstTopic    string
	Block         string
	Removed       bool
	CreatedAt     time.Time
	Serialization []byte
}

func (b *CvmLogs) ComputeID() {
	idsv := fmt.Sprintf("%s:%s:%d", b.BlockHash, b.TxHash, b.LogIndex)
	id := ids.ID(hashing.ComputeHash256Array([]byte(idsv)))
	b.ID = id.String()
}

type MultisigAlias struct {
	Alias         string
	Memo          string
	Bech32Address string
	Owner         string
	TransactionID string
	CreatedAt     time.Time
}

func (p *persist) InsertMultisigAlias(ctx context.Context, session dbr.SessionRunner, alias *MultisigAlias) error {
	var err error
	_, err = session.
		InsertInto(TableMultisigAliases).
		Pair("alias", alias.Alias).
		Pair("memo", alias.Memo).
		Pair("owner", alias.Owner).
		Pair("transaction_id", alias.TransactionID).
		Pair("created_at", alias.CreatedAt).
		ExecContext(ctx)

	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableMultisigAliases, false, err)
	}
	return nil
}

func (p *persist) QueryMultisigAlias(
	ctx context.Context,
	session dbr.SessionRunner,
	alias string,
) (*[]MultisigAlias, error) {
	v := &[]MultisigAlias{}
	err := session.Select(
		"alias",
		"memo",
		"owner",
		"transaction_id",
		"created_at",
	).From(TableMultisigAliases).
		Where("alias=?", alias).
		LoadOneContext(ctx, v)
	return v, err
}

func (p *persist) QueryMultisigAliasesForOwners(
	ctx context.Context,
	session dbr.SessionRunner,
	owners []string,
) (*[]MultisigAlias, error) {
	v := &[]MultisigAlias{}
	_, err := session.Select(
		"bech32_address",
	).From(TableMultisigAliases).Join(TableAddressBech32, "alias=address").
		Where("owner IN ?", owners).
		GroupBy("alias").
		LoadContext(ctx, v)
	return v, err
}

func (p *persist) DeleteMultisigAlias(
	ctx context.Context,
	session dbr.SessionRunner,
	alias string,
) error {
	_, err := session.DeleteFrom(TableMultisigAliases).Where("alias=?", alias).ExecContext(ctx)
	if err != nil {
		return EventErr(TableMultisigAliases, false, err)
	}
	return nil
}

type RewardOwner struct {
	Address   string
	Hash      string
	CreatedAt time.Time
}

func (p *persist) InsertRewardOwner(ctx context.Context, session dbr.SessionRunner, owner *RewardOwner) error {
	var err error
	_, err = session.
		InsertInto(TableRewardOwner).
		Pair("address", owner.Address).
		Pair("hash", owner.Hash).
		Pair("created_at", owner.CreatedAt).
		ExecContext(ctx)

	if err != nil && !utils.ErrIsDuplicateEntryError(err) {
		return EventErr(TableRewardOwner, false, err)
	}
	return nil
}

type RewardType int

const (
	Deposit RewardType = iota
	Validator
)

type Reward struct {
	RewardOwnerBytes []byte
	RewardOwnerHash  string
	TxID             string
	Type             RewardType
	CreatedAt        time.Time
}

func (p *persist) InsertReward(ctx context.Context, session dbr.SessionRunner, reward *Reward) error {
	var err error
	_, err = session.
		InsertBySql("INSERT INTO "+TableReward+" (reward_owner_bytes, reward_owner_hash, tx_id, type, updated_at) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE updated_at=?",
			reward.RewardOwnerBytes,
			reward.RewardOwnerHash,
			reward.TxID,
			reward.Type,
			reward.CreatedAt,
			reward.CreatedAt,
		).ExecContext(ctx)
	if err != nil {
		return EventErr(TableReward, false, err)
	}
	return nil
}

func (p *persist) DeactivateReward(ctx context.Context, session dbr.SessionRunner, reward *Reward) error {
	// Need to fetch rewardOwner and type
	v := &Reward{}
	err := session.Select(
		"reward_owner_bytes",
		"reward_owner_hash",
		"type",
	).From(TableReward).
		Where("tx_id=?", reward.TxID).
		LoadOneContext(ctx, v)
	if err != nil {
		return EventErr(TableReward, false, err)
	}
	v.CreatedAt = reward.CreatedAt

	_, err = session.DeleteFrom(TableReward).
		Where("tx_id=?", reward.TxID).
		ExecContext(ctx)
	if err != nil {
		return EventErr(TableReward, false, err)
	}

	return p.InsertReward(ctx, session, v)
}
