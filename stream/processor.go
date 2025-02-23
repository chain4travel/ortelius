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

package stream

import (
	"context"
	"errors"
	"time"

	"github.com/chain4travel/magellan/cfg"
	"github.com/chain4travel/magellan/db"
	"github.com/chain4travel/magellan/servicesctrl"
	"github.com/chain4travel/magellan/utils"
)

var (
	processorFailureRetryInterval = 200 * time.Millisecond

	// ErrNoMessage is no message
	ErrNoMessage = errors.New("no message")
)

type (
	ProcessorFactoryChainDB func(*servicesctrl.Control, cfg.Config, string, string) (ProcessorDB, error)
	ProcessorFactoryInstDB  func(*servicesctrl.Control, cfg.Config) (ProcessorDB, error)
)

type ProcessorDB interface {
	Process(*utils.Connections, *db.TxPool) error
	Close() error
	ID() string
	Topic() []string
}

func UpdateTxPool(
	ctxTimeout time.Duration,
	conns *utils.Connections,
	persist db.Persist,
	txPool *db.TxPool,
	sc *servicesctrl.Control,
) error {
	sess := conns.DB().NewSessionForEventReceiver(conns.Stream().NewJob("update-tx-pool"))

	ctx, cancelCtx := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancelCtx()

	err := persist.InsertTxPool(ctx, sess, txPool)
	if err == nil {
		sc.Enqueue(txPool)
	}
	return err
}
