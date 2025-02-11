// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package pvm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gocraft/dbr/v2"

	"github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/multisig"
	"github.com/ava-labs/avalanchego/vms/components/verify"
	"github.com/ava-labs/avalanchego/vms/platformvm/blocks"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/proposervm/block"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/chain4travel/magellan/cfg"
	"github.com/chain4travel/magellan/db"
	"github.com/chain4travel/magellan/models"
	"github.com/chain4travel/magellan/services"
	avaxIndexer "github.com/chain4travel/magellan/services/indexes/avax"
	"github.com/chain4travel/magellan/utils"
)

var (
	MaxSerializationLen = (16 * 1024 * 1024) - 1

	ChainID = ids.ID{}

	ErrUnknownBlockType = errors.New("unknown block type")
)

type Writer struct {
	chainID     string
	networkID   uint32
	avaxAssetID ids.ID

	avax *avaxIndexer.Writer
	ctx  *snow.Context
}

func NewWriter(networkID uint32, chainID string) (*Writer, error) {
	_, avaxAssetID, err := genesis.FromConfig(genesis.GetConfig(networkID))
	if err != nil {
		return nil, err
	}

	bcLookup := ids.NewAliaser()
	id, err := ids.FromString(chainID)
	if err != nil {
		return nil, err
	}
	if err = bcLookup.Alias(id, "P"); err != nil {
		return nil, err
	}

	ctx := &snow.Context{
		NetworkID: networkID,
		ChainID:   id,
		Log:       logging.NoLog{},
		Metrics:   metrics.NewOptionalGatherer(),
		BCLookup:  bcLookup,
	}

	return &Writer{
		chainID:     chainID,
		networkID:   networkID,
		avaxAssetID: avaxAssetID,
		avax:        avaxIndexer.NewWriter(chainID, avaxAssetID),
		ctx:         ctx,
	}, nil
}

func (*Writer) Name() string { return "pvm-index" }

type PtxDataModel struct {
	Tx        *txs.Tx               `json:"tx,omitempty"`
	TxType    *string               `json:"txType,omitempty"`
	Block     *blocks.Block         `json:"block,omitempty"`
	BlockID   *string               `json:"blockID,omitempty"`
	BlockType *string               `json:"blockType,omitempty"`
	Proposer  *models.BlockProposal `json:"proposer,omitempty"`
}

func (w *Writer) ParseJSON(b []byte, proposer *models.BlockProposal) ([]byte, error) {
	// Try and parse as a tx
	tx, err := txs.Parse(blocks.GenesisCodec, b)
	if err == nil {
		tx.Unsigned.InitCtx(w.ctx)
		// TODO: Should we be reporting the type of [tx.Unsigned] rather than
		//       `tx`?
		txtype := reflect.TypeOf(tx)
		txtypeS := txtype.String()
		return json.Marshal(&PtxDataModel{
			Tx:     tx,
			TxType: &txtypeS,
		})
	}

	// Try and parse as block
	blk, err := blocks.Parse(blocks.GenesisCodec, b)
	if err == nil {
		blk.InitCtx(w.ctx)
		blkID := blk.ID()
		blkIDStr := blkID.String()
		btype := reflect.TypeOf(blk)
		btypeS := btype.String()
		return json.Marshal(&PtxDataModel{
			BlockID:   &blkIDStr,
			Block:     &blk,
			BlockType: &btypeS,
		})
	}

	// Try and parse as proposervm block
	proposerBlock, err := block.Parse(b)
	if err != nil {
		return nil, err
	}

	blk, err = blocks.Parse(blocks.GenesisCodec, proposerBlock.Block())
	if err != nil {
		return nil, err
	}

	blk.InitCtx(w.ctx)
	blkID := blk.ID()
	blkIDStr := blkID.String()
	btype := reflect.TypeOf(blk)
	btypeS := btype.String()
	return json.Marshal(&PtxDataModel{
		BlockID:   &blkIDStr,
		Block:     &blk,
		BlockType: &btypeS,
		Proposer:  proposer,
	})
}

func (w *Writer) ConsumeConsensus(_ context.Context, _ *utils.Connections, _ services.Consumable, _ db.Persist) error {
	return nil
}

func (w *Writer) Consume(ctx context.Context, conns *utils.Connections, c services.Consumable, persist db.Persist) error {
	job := conns.Stream().NewJob("pvm-index")
	sess := conns.DB().NewSessionForEventReceiver(job)

	dbTx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer dbTx.RollbackUnlessCommitted()

	// Consume the tx and commit
	err = w.indexBlock(services.NewConsumerContext(ctx, dbTx, c.Timestamp(), c.Nanosecond(), persist, c.ChainID()), c.Body())
	if err != nil {
		return err
	}
	return dbTx.Commit()
}

func (w *Writer) Bootstrap(ctx context.Context, conns *utils.Connections, persist db.Persist, gc *utils.GenesisContainer) error {
	txDupCheck := set.NewSet[ids.ID](2*len(gc.Genesis.Camino.AddressStates) +
		2*len(gc.Genesis.Camino.ConsortiumMembersNodeIDs))

	addressStateTx := func(addr ids.ShortID, state txs.AddressStateBit) *txs.Tx {
		tx := &txs.Tx{
			Unsigned: &txs.AddressStateTx{
				BaseTx: txs.BaseTx{
					BaseTx: avax.BaseTx{
						NetworkID:    gc.NetworkID,
						BlockchainID: ChainID,
					},
				},
				Address: addr,
				State:   state,
				Remove:  false,
			},
		}
		if tx.Sign(txs.GenesisCodec, nil) != nil || txDupCheck.Contains(tx.ID()) {
			return nil
		}
		txDupCheck.Add(tx.ID())
		return tx
	}

	var (
		job  = conns.Stream().NewJob("bootstrap")
		db   = conns.DB().NewSessionForEventReceiver(job)
		cCtx = services.NewConsumerContext(ctx, db, int64(gc.Time), 0, persist, w.chainID)
	)

	for _, utxo := range gc.Genesis.UTXOs {
		select {
		case <-ctx.Done():
		default:
		}

		_, _, err := w.avax.ProcessStateOut(
			cCtx,
			utxo.Out,
			utxo.TxID,
			utxo.OutputIndex,
			utxo.AssetID(),
			0,
			0,
			w.chainID,
			false,
			true,
		)
		if err != nil {
			return err
		}
	}

	platformTx := gc.Genesis.Validators
	platformTx = append(platformTx, gc.Genesis.Chains...)
	for _, tx := range platformTx {
		select {
		case <-ctx.Done():
		default:
		}

		err := w.indexTransaction(cCtx, ChainID, tx, true)
		if err != nil {
			return err
		}
	}

	for _, as := range gc.Genesis.Camino.AddressStates {
		select {
		case <-ctx.Done():
		default:
		}

		if as.State&txs.AddressStateKYCVerified != 0 {
			if tx := addressStateTx(as.Address, txs.AddressStateBitKYCVerified); tx != nil {
				err := w.indexTransaction(cCtx, ChainID, tx, true)
				if err != nil {
					return err
				}
			}
		}
		if as.State&txs.AddressStateConsortiumMember != 0 {
			if tx := addressStateTx(as.Address, txs.AddressStateBitConsortium); tx != nil {
				err := w.indexTransaction(cCtx, ChainID, tx, true)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, cm := range gc.Genesis.Camino.ConsortiumMembersNodeIDs {
		select {
		case <-ctx.Done():
		default:
		}

		tx := &txs.Tx{
			Unsigned: &txs.RegisterNodeTx{
				BaseTx: txs.BaseTx{
					BaseTx: avax.BaseTx{
						NetworkID:    gc.NetworkID,
						BlockchainID: ChainID,
					},
				},
				OldNodeID:        ids.EmptyNodeID,
				NewNodeID:        cm.NodeID,
				NodeOwnerAddress: cm.ConsortiumMemberAddress,
				NodeOwnerAuth:    &secp256k1fx.Input{},
			},
		}

		if tx.Sign(txs.GenesisCodec, nil) == nil && !txDupCheck.Contains(tx.ID()) {
			txDupCheck.Add(tx.ID())
			err := w.indexTransaction(cCtx, ChainID, tx, true)
			if err != nil {
				return err
			}
		}
	}

	for _, ma := range gc.Genesis.Camino.MultisigAliases {
		tx := &txs.Tx{
			Unsigned: &txs.MultisigAliasTx{
				BaseTx: txs.BaseTx{
					BaseTx: avax.BaseTx{
						NetworkID:    gc.NetworkID,
						BlockchainID: ChainID,
					},
				},
				MultisigAlias: *ma,
				Auth:          &secp256k1fx.Input{},
			},
		}
		if tx.Sign(txs.GenesisCodec, nil) == nil {
			err := w.indexTransaction(cCtx, ChainID, tx, true)
			if err != nil {
				return err
			}
		}
	}

	parent := ChainID
	blockIDs, err := genesis.GetGenesisBlocksIDs(gc.GenesisBytes, gc.Genesis)
	if err != nil {
		return err
	}
	for index, block := range gc.Genesis.Camino.Blocks {
		cCtx = services.NewConsumerContext(ctx, db, int64(block.Timestamp), 0, persist, w.chainID)
		if err := w.indexCommonBlock(
			cCtx,
			blockIDs[index],
			models.BlockTypeStandard,
			blocks.CommonBlock{
				PrntID: parent,
				Hght:   uint64(index + 1),
			},
			&models.BlockProposal{},
			nil,
		); err != nil {
			return err
		}
		parent = blockIDs[index]

		platformTx = block.Txs()
		for _, tx := range platformTx {
			select {
			case <-ctx.Done():
			default:
			}

			err := w.indexTransaction(cCtx, blockIDs[index], tx, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Writer) indexBlock(ctx services.ConsumerCtx, blockBytes []byte) error {
	proposerBlock, err := block.Parse(blockBytes)
	var innerBlockBytes []byte
	if err != nil {
		innerBlockBytes = blockBytes
		// We use the "nil"ness below, so we explicitly empty the value here to
		// avoid unexpected errors
		proposerBlock = nil
	} else {
		innerBlockBytes = proposerBlock.Block()
	}

	blk, err := blocks.Parse(blocks.GenesisCodec, innerBlockBytes)
	if err != nil {
		return err
	}

	blkID := blk.ID()
	ctxTime := ctx.Time()
	pvmProposer := models.NewBlockProposal(proposerBlock, &ctxTime)

	adjustCtxTime := func(tm uint64) {
		ctxTime := time.Unix(int64(tm), 0)
		ctx.SetTime(ctxTime)
		pvmProposer.TimeStamp = &ctxTime
	}

	errs := wrappers.Errs{}
	switch blk := blk.(type) {
	case *blocks.ApricotProposalBlock:
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeProposal, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.ApricotStandardBlock:
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeStandard, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.ApricotAtomicBlock:
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeProposal, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.ApricotAbortBlock:
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeAbort, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.ApricotCommitBlock:
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeCommit, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.BanffProposalBlock:
		adjustCtxTime(blk.Time)
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeStandard, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.BanffStandardBlock:
		adjustCtxTime(blk.Time)
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeStandard, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.BanffAbortBlock:
		adjustCtxTime(blk.Time)
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeAbort, blk.CommonBlock, pvmProposer, innerBlockBytes))
	case *blocks.BanffCommitBlock:
		adjustCtxTime(blk.Time)
		errs.Add(w.indexCommonBlock(ctx, blkID, models.BlockTypeCommit, blk.CommonBlock, pvmProposer, innerBlockBytes))
	default:
		return fmt.Errorf("unknown type %T", blk)
	}
	for _, tx := range blk.Txs() {
		errs.Add(w.indexTransaction(ctx, blkID, tx, false))
	}

	return errs.Err
}

func (w *Writer) indexCommonBlock(
	ctx services.ConsumerCtx,
	blkID ids.ID,
	blkType models.BlockType,
	blk blocks.CommonBlock,
	proposer *models.BlockProposal,
	blockBytes []byte,
) error {
	if len(blockBytes) > MaxSerializationLen {
		blockBytes = []byte("")
	}

	pvmBlocks := &db.PvmBlocks{
		ID:            blkID.String(),
		ChainID:       w.chainID,
		Type:          blkType,
		ParentID:      blk.Parent().String(),
		Serialization: blockBytes,
		CreatedAt:     ctx.Time(),
		Height:        blk.Height(),
		Proposer:      proposer.Proposer,
		ProposerTime:  proposer.TimeStamp,
	}
	return ctx.Persist().InsertPvmBlocks(ctx.Ctx(), ctx.DB(), pvmBlocks, cfg.PerformUpdates)
}

//nolint:gocyclo
func (w *Writer) indexTransaction(ctx services.ConsumerCtx, blkID ids.ID, tx *txs.Tx, genesis bool) error {
	var (
		txID   = tx.ID()
		baseTx avax.BaseTx
		typ    models.TransactionType
		ins    *avaxIndexer.AddInsContainer
		outs   *avaxIndexer.AddOutsContainer
	)
	switch castTx := tx.Unsigned.(type) {
	case *txs.AddValidatorTx:
		baseTx = castTx.BaseTx.BaseTx
		outs = &avaxIndexer.AddOutsContainer{
			Outs:    castTx.StakeOuts,
			Stake:   true,
			ChainID: w.chainID,
		}
		typ = models.TransactionTypeAddValidator
		err := w.InsertTransactionValidator(ctx, txID, castTx.Validator)
		if err != nil {
			return err
		}
	case *txs.AddSubnetValidatorTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddSubnetValidator
	case *txs.CreateSubnetTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeCreateSubnet
	case *txs.CreateChainTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeCreateChain
	case *txs.ImportTx:
		baseTx = castTx.BaseTx.BaseTx
		ins = &avaxIndexer.AddInsContainer{
			Ins:     castTx.ImportedInputs,
			ChainID: castTx.SourceChain.String(),
		}
		typ = models.TransactionTypePVMImport
	case *txs.ExportTx:
		baseTx = castTx.BaseTx.BaseTx
		outs = &avaxIndexer.AddOutsContainer{
			Outs:    castTx.ExportedOutputs,
			ChainID: castTx.DestinationChain.String(),
		}
		typ = models.TransactionTypePVMExport
	case *txs.AdvanceTimeTx:
		return nil
	case *txs.RemoveSubnetValidatorTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeRemoveSubnetValidator
	case *txs.TransformSubnetTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeTransformSubnet
	case *txs.AddPermissionlessValidatorTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddPermissionlessValidator

		// TODO: Handle this for all subnetIDs
		if castTx.Subnet != constants.PrimaryNetworkID {
			break
		}

		outs = &avaxIndexer.AddOutsContainer{
			Outs:    castTx.StakeOuts,
			Stake:   true,
			ChainID: w.chainID,
		}
		err := w.InsertTransactionValidator(ctx, txID, castTx.Validator)
		if err != nil {
			return err
		}
	case *txs.AddPermissionlessDelegatorTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddPermissionlessDelegator

		// TODO: Handle this for all subnetIDs
		if castTx.Subnet != constants.PrimaryNetworkID {
			break
		}

		outs = &avaxIndexer.AddOutsContainer{
			Outs:    castTx.StakeOuts,
			Stake:   true,
			ChainID: w.chainID,
		}
		err := w.InsertTransactionValidator(ctx, txID, castTx.Validator)
		if err != nil {
			return err
		}
	case *txs.CaminoAddValidatorTx:
		innerTx := castTx.AddValidatorTx
		baseTx = innerTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddValidator
		err := w.InsertTransactionValidator(ctx, txID, innerTx.Validator)
		if err != nil {
			return err
		}
		if castTx.RewardsOwner != nil {
			err = w.insertReward(ctx, txID, castTx.RewardsOwner, db.Validator)
			if err != nil {
				return err
			}
		}
	case *txs.DepositTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeDeposit
		if castTx.RewardsOwner != nil {
			err := w.insertReward(ctx, txID, castTx.RewardsOwner, db.Deposit)
			if err != nil {
				return err
			}
		}
	case *txs.UnlockDepositTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeUnlockDeposit
	case *txs.AddressStateTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddAddressState
	case *txs.RegisterNodeTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeRegisterNodeTx
	case *txs.BaseTx:
		baseTx = castTx.BaseTx
		typ = models.TransactionTypePvmBase
	case *txs.MultisigAliasTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeMultisigAlias
		err := w.InsertMultisigAlias(ctx, &castTx.MultisigAlias, castTx.Auth, txID)
		if err != nil {
			return err
		}
	case *txs.ClaimTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeClaimReward
	case *txs.RewardsImportTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeRewardsImport
	case *txs.CaminoRewardValidatorTx:
		baseTx = avax.BaseTx{
			NetworkID:    w.networkID,
			BlockchainID: w.ctx.ChainID,
			Ins:          castTx.Ins,
			Outs:         castTx.Outs,
		}
		typ = models.TransactionTypeCaminoRewardValidator
	case *txs.AddDepositOfferTx:
		baseTx = castTx.BaseTx.BaseTx
		typ = models.TransactionTypeAddDepositOffer
	default:
		return fmt.Errorf("unknown tx type %T", castTx)
	}

	err := w.InsertTransactionBlock(ctx, txID, blkID)
	if err != nil {
		return err
	}

	return w.avax.InsertTransaction(
		ctx,
		tx.Bytes(),
		tx.ID(),
		tx.Unsigned.Bytes(),
		&baseTx,
		tx.Creds,
		typ,
		ins,
		outs,
		0,
		genesis,
	)
}

func (w *Writer) insertReward(ctx services.ConsumerCtx, txID ids.ID, rewardOwner verify.Verifiable, rewardType db.RewardType) error {
	owner, ok := rewardOwner.(*secp256k1fx.OutputOwners)
	if !ok {
		return fmt.Errorf("rewardOwner %T", rewardOwner)
	}
	ownerID, err := txs.GetOwnerID(rewardOwner)
	if err != nil {
		return fmt.Errorf("rewardOwner hash %v", err)
	}
	ownerIDStr := ownerID.String()
	ownerBytes, err := blocks.GenesisCodec.Marshal(txs.Version, rewardOwner)
	if err != nil {
		return fmt.Errorf("rewardOwner bytes %v", err)
	}

	err = ctx.Persist().InsertReward(ctx.Ctx(), ctx.DB(), &db.Reward{
		RewardOwnerBytes: ownerBytes,
		RewardOwnerHash:  ownerIDStr,
		TxID:             txID.String(),
		Type:             rewardType,
		CreatedAt:        ctx.Time(),
	})
	if err != nil {
		return fmt.Errorf("rewardOwner insertReward %v", err)
	}

	// Ingest each Output Address
	for _, addr := range owner.Addresses() {
		addrStr, err := cb58.Encode(addr)
		if err != nil {
			return fmt.Errorf("rewardOwner %v", err)
		}

		err = ctx.Persist().InsertRewardOwner(ctx.Ctx(), ctx.DB(), &db.RewardOwner{
			Address:   addrStr,
			Hash:      ownerIDStr,
			CreatedAt: ctx.Time(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) InsertTransactionValidator(ctx services.ConsumerCtx, txID ids.ID, validator txs.Validator) error {
	transactionsValidator := &db.TransactionsValidator{
		ID:        txID.String(),
		NodeID:    validator.NodeID.String(),
		Start:     validator.Start,
		End:       validator.End,
		CreatedAt: ctx.Time(),
	}
	return ctx.Persist().InsertTransactionsValidator(ctx.Ctx(), ctx.DB(), transactionsValidator, cfg.PerformUpdates)
}

func (w *Writer) InsertTransactionBlock(ctx services.ConsumerCtx, txID ids.ID, blkTxID ids.ID) error {
	transactionsBlock := &db.TransactionsBlock{
		ID:        txID.String(),
		TxBlockID: blkTxID.String(),
		CreatedAt: ctx.Time(),
	}
	return ctx.Persist().InsertTransactionsBlock(ctx.Ctx(), ctx.DB(), transactionsBlock, cfg.PerformUpdates)
}

func (w *Writer) InsertMultisigAlias(
	ctx services.ConsumerCtx,
	alias *multisig.Alias,
	auth verify.Verifiable,
	txID ids.ID,
) error {
	var err error

	// If alias.ID is an empty ID, then it's a new alias, and we need to generate the aliasID from the txID
	if alias.ID == ids.ShortEmpty {
		alias.ID = multisig.ComputeAliasID(txID)
	}

	_, err = ctx.Persist().QueryMultisigAlias(ctx.Ctx(), ctx.DB(), alias.ID.String())
	if err != nil && err != dbr.ErrNotFound {
		return err
	}

	// if there is an already existing alias with this aliasID or auth is nil, then we need to delete it
	if auth == nil || err == nil {
		err = ctx.Persist().DeleteMultisigAlias(ctx.Ctx(), ctx.DB(), alias.ID.String())
		if err != nil {
			return err
		}
	}

	// add alias to bech32 address mapping table
	err = persistMultisigAliasAddresses(ctx, alias.ID, w.chainID)
	if err != nil {
		return err
	}

	// Get owner addresses
	owner, ok := alias.Owners.(*secp256k1fx.OutputOwners)
	if !ok {
		return fmt.Errorf("could not parse Multisig owners %T", alias.Owners)
	}

	// Loop over owner addresses and insert an entry for each
	for _, addr := range owner.Addresses() {
		addrid, err := ids.ToShortID(addr)
		if err != nil {
			return err
		}
		multisigAlias := &db.MultisigAlias{
			Alias:         alias.ID.String(),
			Memo:          string(alias.Memo),
			Owner:         addrid.String(),
			TransactionID: txID.String(),
			CreatedAt:     ctx.Time(),
		}

		err = ctx.Persist().InsertMultisigAlias(ctx.Ctx(), ctx.DB(), multisigAlias)
		if err != nil {
			return err
		}

		// add owner address to bech32 address mapping table
		err = persistMultisigAliasAddresses(ctx, addrid, w.chainID)
		if err != nil {
			return err
		}
	}
	return nil
}

func persistMultisigAliasAddresses(ctx services.ConsumerCtx, addr ids.ShortID, chainID string) error {
	var err error

	// add alias and owners to address table
	addressChain := &db.AddressChain{
		Address:   addr.String(),
		ChainID:   chainID,
		CreatedAt: ctx.Time(),
		UpdatedAt: time.Now().UTC(),
	}
	err = ctx.Persist().InsertAddressChain(ctx.Ctx(), ctx.DB(), addressChain, cfg.PerformUpdates)
	if err != nil {
		return err
	}

	bech32Addr, err := address.FormatBech32(models.Bech32HRP, addr.Bytes())
	if err != nil {
		return err
	}

	addressBech32 := &db.AddressBech32{
		Address:       addr.String(),
		Bech32Address: bech32Addr,
		UpdatedAt:     time.Now().UTC(),
	}

	err = ctx.Persist().InsertAddressBech32(ctx.Ctx(), ctx.DB(), addressBech32, cfg.PerformUpdates)
	if err != nil {
		return err
	}

	return nil
}
