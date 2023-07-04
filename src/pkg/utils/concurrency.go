// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
// forked from https://www.socketloop.com/tutorials/golang-byte-format-example

// Package utils provides generic helper functions.
package utils

import (
	"context"
	"sync"

	"github.com/defenseunicorns/zarf/src/pkg/message"
)

type ConcurrencyTools struct {
	ProgressChan chan string
	ErrorChan    chan message.ErrorWithMessage
	Context      context.Context
	Cancel       context.CancelFunc
	WaitGroup    *sync.WaitGroup
}

func NewConcurrencyTools(length int) *ConcurrencyTools {
	ctx, cancel := context.WithCancel(context.Background())

	progressChan := make(chan string, length)

	errorChan := make(chan message.ErrorWithMessage, length)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(length)

	concurrencyTools :=  ConcurrencyTools{
		ProgressChan: progressChan,
		ErrorChan:    errorChan,
		Context:      ctx,
		Cancel:       cancel,
		WaitGroup:    &waitGroup,
	}

	return &concurrencyTools
}


func ContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
