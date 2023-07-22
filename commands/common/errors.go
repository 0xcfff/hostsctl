package common

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyArguments   = errors.New("too many arguments")
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrIpOrAliasExpected  = fmt.Errorf("IP or alias expected %w", ErrNotEnoughArguments)
	ErrTooManyEntries     = errors.New("too many entries found")
	ErrBlockNotFound      = errors.New("block not found")
	ErrAliasNotFound      = errors.New("alias not found")
)