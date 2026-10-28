// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"errors"
	"fmt"

	"echo.labstack.com/pro/quota"
)

// docs:start ExampleLimiter_Generations
func ExampleLimiter_Generations() {
	ctx := context.Background()
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{quota.PerDay("calls", 10)}}}})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	request := quota.ResetRequest{Identity: quota.Identity{Tenant: "acme"}, Plan: "team"}
	generations, err := l.Generations(ctx, request)
	if err != nil {
		panic(err)
	}
	request.Expected = generations
	result, err := l.Reset(ctx, request)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Generations["calls"])
	_, err = l.Reset(ctx, request)
	fmt.Println(errors.Is(err, quota.ErrResetConflict))
	// Output:
	// 1
	// true
}

// docs:end ExampleLimiter_Generations
