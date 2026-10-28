// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"fmt"
	"time"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
)

func seedQuota(ctx context.Context) (*quota.Limiter, *quotaadmin.Panel, error) {
	denials, err := quotaadmin.NewDenialRecorder(256)
	if err != nil {
		return nil, nil, err
	}
	top := quota.NewTopK(10)
	//nolint:contextcheck // Construction has no context API; no Outbox is configured.
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Hooks: quota.JoinHooks(top.Hooks(), denials.Hooks()), Plans: []quota.Plan{
		{ID: "team", Limits: []quota.Limit{quota.PerDay("daily", 100).Meter("requests"), quota.Credits("credits").Meter("tokens")}},
		{ID: "growth_preview", Limits: []quota.Limit{quota.PerDay("daily", 45).Meter("requests")}},
		{ID: "enterprise", Limits: []quota.Limit{quota.PerDay("daily", 10000).Meter("requests"), quota.Credits("credits").Meter("tokens")}},
	}})
	if err != nil {
		return nil, nil, err
	}
	ok := false
	defer func() {
		if !ok {
			closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
			defer cancel()
			if closeErr := lim.Close(closeCtx); closeErr != nil {
				fmt.Println("quota seed cleanup failed")
			}
		}
	}()
	for i, tenant := range demoTenants {
		identity := quota.Identity{Tenant: tenant}
		for _, plan := range []string{"team", "enterprise"} {
			if _, err = lim.AddCredits(ctx, quota.CreditRequest{Identity: identity, Plan: plan, Limit: "credits", Amount: []int64{1000, 600, 2000, 2500}[i], ID: "opening-balance"}); err != nil {
				return nil, nil, err
			}
		}
		for n := range 65 + i*25 {
			d, allowErr := lim.Allow(ctx, quota.Request{Identity: identity, Plan: "team", Shadow: "growth_preview", Operation: "POST /api/orders", Cost: quota.Cost{"requests": 1, "tokens": int64(8 + i*3)}, CorrelationID: fmt.Sprintf("demo-%s-%03d", tenant, n)})
			if allowErr != nil {
				return nil, nil, allowErr
			}
			if d.Err != nil {
				return nil, nil, d.Err
			}
		}
		if _, err = lim.Allow(ctx, quota.Request{Identity: identity, Plan: "enterprise", Cost: quota.Cost{"requests": int64(1200 + 350*i), "tokens": 300}}); err != nil {
			return nil, nil, err
		}
	}
	p, err := quotaadmin.New(quotaadmin.Config{Limiter: lim, TopK: top, Denials: denials, EnableCredits: true, EnableReset: true})
	ok = err == nil
	return lim, p, err
}
