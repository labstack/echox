// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"sync"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/audit/auditadmin"
)

type demoAuditStore interface {
	audit.Recorder
	audit.ReadBackend
}

// memoryAuditSink keeps bounded semantic idempotency for the development demo.
// It does not claim durable audit guarantees or survive process exit.
type memoryAuditSink struct {
	mu        sync.Mutex
	store     *audit.MemoryStore
	semantics map[string]string
}

func (s *memoryAuditSink) Record(ctx context.Context, r admin.AuditRecord) error {
	return s.record(ctx, "emission:"+r.EmissionID, r)
}
func (s *memoryAuditSink) RecordAnnotation(ctx context.Context, r admin.AuditRecord) error {
	return s.record(ctx, "annotation:"+r.OperationID, r)
}
func (s *memoryAuditSink) record(ctx context.Context, key string, r admin.AuditRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p, err := auditadmin.Prepare(r)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if previous, ok := s.semantics[key]; ok {
		if previous != r.SemanticDigest {
			return admin.ErrAuditConflict
		}
		return nil
	}
	if len(s.semantics) >= 1000 {
		return &audit.ConfigError{Field: "demo audit capacity"}
	}
	if _, err = s.store.Record(ctx, p); err != nil {
		return err
	}
	s.semantics[key] = r.SemanticDigest
	return nil
}
