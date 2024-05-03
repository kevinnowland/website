package main

import (
	"context"
	"net"
)

const (
	LogStrings = "LogStrings"
	Errors     = "Errors"
)

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	ctx = context.WithValue(ctx, Errors, []error{})
	ctx = context.WithValue(ctx, LogStrings, map[string]string{})
	return ctx
}
