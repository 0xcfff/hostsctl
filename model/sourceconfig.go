package model

import (
	"crypto/sha1"
)

const (
	// synchronization source
	CFGPROP_SOURCE = "source"
	// synchronization source
	TYPE = "type"
	// refresh interval
	INTERVAL = "interval"
)

type SourceConfig interface {

	// verifies configuration schema
	VerifySchema(schema *ConfigSchema) error

	// returns value of one single source configuration property
	Property(name string) (value string, ok bool)

	// return all known properties
	Properties(names []string) map[string]string

	// method calculates hash based on the source configuration
	// if configuration changes, the hash should change as well
	ConfigHash() []byte
}

type sourceConfig struct {
	properties map[string]string
}

func (s *sourceConfig) Property(name string) (value string, ok bool) {
	val, okk := s.properties[name]
	return val, okk
}

func (s *sourceConfig) Properties(names []string) map[string]string {
	res := make(map[string]string, len(names))
	for _, k := range names {
		val, ok := s.properties[k]
		if ok {
			res[k] = val
		}
	}
	return res
}

func (s *sourceConfig) VerifySchema(schema *ConfigSchema) error {
	return nil
}

func (s *sourceConfig) ConfigHash() []byte {
	configProps := s.properties
	hasher := sha1.New()
	for key := range configProps {
		hasher.Write([]byte(key))
		hasher.Write([]byte(configProps[key]))
	}
	binHash := hasher.Sum(nil)
	return binHash
}

func NewSourceConfig(properties map[string]string) SourceConfig {
	cfg := sourceConfig{properties: properties}
	return &cfg
}
