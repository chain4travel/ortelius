// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************

package rewards

import (
	"context"
	"fmt"
	"time"

	"github.com/chain4travel/caminogo/api"
	"github.com/chain4travel/caminogo/ids"
	"github.com/chain4travel/caminogo/utils/formatting"
	caminoGoAvax "github.com/chain4travel/caminogo/vms/components/avax"
	"github.com/chain4travel/caminogo/vms/platformvm"
	"github.com/chain4travel/magellan/db"
	"github.com/chain4travel/magellan/models"
	"github.com/chain4travel/magellan/services"
	"github.com/chain4travel/magellan/services/indexes/avax"
	"github.com/chain4travel/magellan/servicesctrl"
	"github.com/chain4travel/magellan/utils"
)

type Handler struct {
	client      platformvm.Client
	conns       *utils.Connections
	perist      db.Persist
	avaxAssetID ids.ID
	writer      *avax.Writer
	cid         ids.ID
	doneCh      chan struct{}
}

func (r *Handler) Start(sc *servicesctrl.Control) error {
	conns, err := sc.Database()
	if err != nil {
		return err
	}
	go r.runTicker(sc, conns)
	return nil
}

func (r *Handler) Close() {
	close(r.doneCh)
}

func (r *Handler) runTicker(sc *servicesctrl.Control, conns *utils.Connections) {
	sc.Log.Info("start")
	defer func() {
		sc.Log.Info("stop")
	}()

	ticker := time.NewTicker(24 * time.Hour)

	r.doneCh = make(chan struct{}, 1)

	r.conns = conns
	r.client = platformvm.NewClient(sc.ServicesCfg.CaminoNode)
	r.perist = db.NewPersist()

	r.avaxAssetID = sc.GenesisContainer.AvaxAssetID

	r.cid = ids.Empty
	r.writer = avax.NewWriter(r.cid.String(), r.avaxAssetID)

	defer func() {
		ticker.Stop()
		_ = conns.Close()
	}()

	for {
		select {
		case <-ticker.C:
			err := r.processRewards()
			if err != nil {
				sc.Log.Error("process rewards %s", err)
			}
		case <-r.doneCh:
			return
		}
	}
}

func (r *Handler) processRewards() error {
	job := r.conns.Stream().NewJob("rewards-handler")
	sess := r.conns.DB().NewSessionForEventReceiver(job)

	ctx := context.Background()

	var err error

	type RewardTx struct {
		ID        string
		Txid      string
		Type      models.BlockType
		CreatedAt time.Time
	}
	var reardsTxs []RewardTx
	_, err = sess.Select(
		db.TableRewards+".id",
		db.TableRewards+".txid",
		db.TablePvmBlocks+".type",
		db.TableRewards+".created_at",
	).
		From(db.TableRewards).
		Join(db.TablePvmBlocks, db.TableRewards+".block_id = "+db.TablePvmBlocks+".parent_id").
		Where(db.TableRewards+".processed = ? and "+db.TableRewards+".created_at < ?", 0, time.Now().Add(-3*time.Second)).
		LoadContext(ctx, &reardsTxs)
	if err != nil {
		return err
	}
	if len(reardsTxs) == 0 {
		return nil
	}

	for _, rewardTx := range reardsTxs {
		if rewardTx.Type == models.BlockTypeAbort {
			err = r.markRewardProcessed(rewardTx.ID)
			if err != nil {
				return err
			}
			continue
		}

		id, err := ids.FromString(rewardTx.Txid)
		if err != nil {
			return err
		}
		var rewardsUtxos [][]byte
		arg := &api.GetTxArgs{TxID: id, Encoding: formatting.Hex}
		rewardsUtxos, err = r.client.GetRewardUTXOs(ctx, arg)
		if err != nil {
			return err
		}

		if len(rewardsUtxos) == 0 {
			return fmt.Errorf("no rewards %s", rewardTx.Txid)
		}

		err = r.processRewardUtxos(rewardsUtxos, rewardTx.CreatedAt)
		if err != nil {
			return err
		}

		err = r.markRewardProcessed(rewardTx.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Handler) processRewardUtxos(rewardsUtxos [][]byte, createdAt time.Time) error {
	job := r.conns.Stream().NewJob("rewards-handler-persist")
	sess := r.conns.DB().NewSessionForEventReceiver(job)

	dbTx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer dbTx.RollbackUnlessCommitted()

	ctx := context.Background()

	for _, reawrdUtxo := range rewardsUtxos {
		var utxo *caminoGoAvax.UTXO
		_, err = platformvm.Codec.Unmarshal(reawrdUtxo, &utxo)
		if err != nil {
			return err
		}

		cCtx := services.NewConsumerContext(ctx, sess, createdAt.Unix(), int64(createdAt.Nanosecond()), r.perist, r.cid.String())

		_, _, err = r.writer.ProcessStateOut(
			cCtx,
			utxo.Out,
			utxo.TxID,
			utxo.OutputIndex,
			utxo.AssetID(),
			0,
			0,
			r.cid.String(),
			false,
			false,
		)
		if err != nil {
			return err
		}
	}

	return dbTx.Commit()
}

func (r *Handler) markRewardProcessed(id string) error {
	job := r.conns.Stream().NewJob("rewards-handler")
	sess := r.conns.DB().NewSessionForEventReceiver(job)

	ctx := context.Background()

	reward := &db.Rewards{
		ID:        id,
		Processed: 1,
	}

	return r.perist.UpdateRewardsProcessed(ctx, sess, reward)
}
