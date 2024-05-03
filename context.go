package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/google/uuid"
)

const (
	LogAttrs  = "LogAttrs"
	Errors    = "Errors"
	RequestId = "RequestId"
)

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	ctx = context.WithValue(ctx, RequestId, uuid.New())
	ctx = context.WithValue(ctx, Errors, []error{})
	ctx = context.WithValue(ctx, LogAttrs, []slog.Attr{})
	return ctx
}
