// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"errors"
	"time"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/webhook/webhookpg"
)

type executionTrends interface {
	Trends(context.Context, execution.TrendQuery) (execution.TrendPage, error)
}
type webhookTrends interface {
	Trends(context.Context, webhookpg.TrendQuery) (webhookpg.TrendPage, error)
}

// Real seed events happen during startup. Include the open hour so the demo
// starts with visible bars, and fail early if either event source is missing.
func verifySeededThroughput(ctx context.Context, executions executionTrends, webhooks webhookTrends) error {
	now := time.Now().UTC()
	from := now.Truncate(time.Hour).Add(-24 * time.Hour)
	for _, namespace := range []string{"payments", "fulfillment"} {
		page, err := executions.Trends(ctx, execution.TrendQuery{Namespace: namespace, From: from, Until: now, Bucket: time.Hour})
		if err != nil {
			return err
		}
		var started, finished int64
		for _, b := range page.Buckets {
			if n := b.Counts["started"]; n != nil {
				started += *n
			}
			for _, key := range []string{"completed", "failed"} {
				if n := b.Counts[key]; n != nil {
					finished += *n
				}
			}
		}
		if page.Incomplete || started == 0 || finished == 0 {
			return errors.New("execution seed has no complete throughput shape")
		}
	}
	for _, direction := range []string{"in", "out"} {
		page, err := webhooks.Trends(ctx, webhookpg.TrendQuery{Direction: direction, From: from, Until: now, Bucket: time.Hour})
		if err != nil {
			return err
		}
		var enqueued, finished int64
		for _, b := range page.Buckets {
			if n := b.Counts["enqueued"]; n != nil {
				enqueued += *n
			}
			for key, n := range b.Counts {
				if key != "enqueued" && n != nil {
					finished += *n
				}
			}
		}
		if page.Incomplete || enqueued == 0 || finished == 0 {
			return errors.New("webhook seed has no complete throughput shape")
		}
	}
	return nil
}
