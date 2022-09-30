package model

import (
	"bytes"
)

type DiffResult int

const (
	DIFF_NONE DiffResult = iota
	DIFF_CONFIG
	DIFF_DATA
	DIFF_ALL
)

func (this SourceState) Differs(other SourceState) DiffResult {
	configsEqual := bytes.Equal(this.config.ConfigHash(), other.config.ConfigHash())
	dataEqual := bytes.Equal(this.data.DataHash(), other.data.DataHash())

	if configsEqual && dataEqual {
		return DIFF_NONE
	}

	type compRes struct{ config, data bool }

	rawCompareResults := compRes{configsEqual, dataEqual}

	switch rawCompareResults {
	case compRes{true, true}:
		return DIFF_NONE
	case compRes{true, false}:
		return DIFF_DATA
	case compRes{false, true}:
		return DIFF_CONFIG
	default:
		return DIFF_ALL
	}
}
