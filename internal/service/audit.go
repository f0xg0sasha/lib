package service

import (
	"context"

	"github.com/f0xg0sasha/audit_logger/pkg/domain/audit"
)

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}
