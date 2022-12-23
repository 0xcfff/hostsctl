package common

import (
	"context"

	"github.com/spf13/afero"
)

type commandContextValue int

const (
	ctxCustomFileSystem commandContextValue = iota
)

// Overrides filesystem used by commands
func WithCustomFilesystem(ctx context.Context, fs afero.Fs) context.Context {
	return context.WithValue(ctx, ctxCustomFileSystem, fs)
}

// Returns filesystem override if any
func FileSystem(ctx context.Context) afero.Fs {
	return ctx.Value(ctxCustomFileSystem).(afero.Fs)
}
