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
	"fmt"

	"go.uber.org/zap"

	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/gocraft/dbr/v2"
	"github.com/palantir/stacktrace"
)

type EventRcvr struct {
	logger logging.Logger
}

func (e *EventRcvr) NewJob(s string) dbr.EventReceiver {
	return &event{name: s, logger: e.logger}
}

func (e *EventRcvr) SetLog(logger logging.Logger) {
	e.logger = logger
}

type event struct {
	name      string
	eventName string
	logger    logging.Logger
}

func (e *event) Event(eventName string) {
	e.eventName = eventName
}

func (e *event) EventKv(eventName string, kvs map[string]string) {
	e.eventName = eventName
}

func (e *event) EventErr(eventName string, err error) error {
	return e.EventErrKv(eventName, err, nil)
}

func (e *event) EventErrKv(eventName string, err error, kvs map[string]string) error {
	if ErrIsDuplicateEntryError(err) || ErrIsLockError(err) {
		return err
	}
	e.eventName = eventName
	e.logger.Warn("event error",
		zap.String("name", e.name),
		zap.String("eventName", e.eventName),
		zap.Error(err),
	)
	return stacktrace.Propagate(err, fmt.Sprintf("%s %s", e.name, e.eventName))
}

func (e *event) Timing(eventName string, nanoseconds int64) {
}

func (e *event) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
}
