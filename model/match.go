package model

import "bytes"

type MatchResult int

const (
	MATCH_NONE MatchResult = iota
	MATCH_FULL
	MATCH_PARTIAL
)

func matchConfigs(c1, c2 SourceConfig) MatchResult {
	firstHash := c1.ConfigHash()
	secondHash := c2.ConfigHash()

	if bytes.Equal(firstHash, secondHash) {
		return MATCH_FULL
	}

	firstSource, fsFound := c1.Property(CFGPROP_SOURCE)
	secondSource, ssFound := c2.Property(CFGPROP_SOURCE)

	if !fsFound || !ssFound || firstSource != secondSource {
		return MATCH_NONE
	}

	return MATCH_PARTIAL
}

func matchData(d1, d2 SourceData) MatchResult {
	firstHash := d1.DataHash()
	secondHash := d2.DataHash()

	if bytes.Equal(firstHash, secondHash) {
		return MATCH_FULL
	}
	return MATCH_NONE
}

func (this SourceState) Matches(other SourceState) MatchResult {

	matchConfigs := matchConfigs(this.Config, other.Config)

	if matchConfigs == MATCH_FULL {
		return MATCH_FULL
	}

	matchData := matchData(this.Data, other.Data)

	if matchData == MATCH_FULL {
		return MATCH_FULL
	}

	if matchConfigs == MATCH_PARTIAL || matchData == MATCH_PARTIAL {
		return MATCH_PARTIAL
	}

	return MATCH_NONE
}
