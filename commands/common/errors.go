package common

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyArguments         = errors.New("too many arguments")
	ErrNotEnoughArguments       = errors.New("not enough arguments")
	ErrIpOrAliasExpected        = fmt.Errorf("IP or alias expected %w", ErrNotEnoughArguments)
	ErrEntryAlreadyExists       = errors.New("entry already exists")
	ErrTooManyEntries           = errors.New("too many entries found")
	ErrWrongArgumentValue       = errors.New("wrong argument value")
	ErrBlockNotFound            = errors.New("block not found")
	ErrAliasNotFound            = errors.New("alias not found")
	ErrNotSupportedOutputFormat = errors.New("not supported output format")
)
