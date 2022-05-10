+++
draft = false
date = 2022-05-04T21:11:24+01:00
title = "Temporal.io errGroup"
description = "Synchronise cancellation on error in Temporal.io"
slug = "temporalio errGroup"
tags = ["Temporal"]
categories = []
externalLink = ""
series = []
+++

Recently during develop a temporal workflow, I found I need an `errGroup`
implementation - an easy way to synchronise the cancellation of all Temporal
coroutines when one of them returned error. I adapted the
[`x/sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup). Here is the
code:

```
package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context) (string, error) {
	g, cc := withErrGroup(ctx)

	for i := 0; i < 3; i++ {
		g.Go(cc, func(ctx workflow.Context) error {
			workflow.Sleep(ctx, 3*time.Second)
			if ctx.Err() != nil {
				println("ctx error", ctx.Err().Error())
			} else {
				panic("shouldn't be here")
			}
			return nil
		})
	}
	g.Go(cc, func(ctx workflow.Context) error {
		return errors.New("foo error")
	})

	err := g.Wait(cc)
	if err == nil {
		return "", errors.New("expected error")
	}
	return "expected error received: " + err.Error(), nil
}

// adapted from https://cs.opensource.google/go/x/sync/+/036812b2:errgroup/errgroup.go
type errGroup struct {
	wg      workflow.WaitGroup
	cancel  func()
	errOnce sync.Once
	err     error
}

func withErrGroup(ctx workflow.Context) (*errGroup, workflow.Context) {
	cc, cancel := workflow.WithCancel(ctx)
	eg := &errGroup{
		cancel: cancel,
		wg:     workflow.NewWaitGroup(ctx),
	}
	return eg, cc
}

func (g *errGroup) Wait(ctx workflow.Context) error {
	g.wg.Wait(ctx)
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *errGroup) Go(ctx workflow.Context, f func(workflow.Context) error) {
	g.wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer g.wg.Done()
		if err := f(ctx); err != nil {
			g.errOnce.Do(func() { // feels like errOnce is not needed
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}

func main() {
	s := &testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	env.ExecuteWorkflow(Workflow)
	err := env.GetWorkflowError()
	if err != nil {
		log.Fatalf("[ERROR] not expecting workflow error: %v", err)
	}
	var result any
	err = env.GetWorkflowResult(&result)
	if err != nil {
		log.Fatalf("[ERROR] failed to get workflow result: %v", err)
	}
	log.Printf("result: %v", result)
}
```

Result is:

```
2022/05/04 22:26:45 DEBUG RequestCancelTimer TimerID 1
2022/05/04 22:26:45 DEBUG RequestCancelTimer TimerID 2
2022/05/04 22:26:45 DEBUG RequestCancelTimer TimerID 3
ctx error canceled
ctx error canceled
ctx error canceled
2022/05/04 22:26:45 result: expected error received: foo error
```
