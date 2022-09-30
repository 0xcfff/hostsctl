package model

type ConfigProperty struct {
	name        string
	description string
	required    bool
}

type ConfigSchema struct {
	properties []ConfigProperty
}

func NewRequiredProperty(name string, description string) ConfigProperty {
	return ConfigProperty{required: true, name: name, description: description}
}
func NewOptionalProperty(name string, description string) ConfigProperty {
	return ConfigProperty{required: false, name: name, description: description}
}

func NewConfigSchema(props ...ConfigProperty) ConfigSchema {
	return ConfigSchema{properties: props}
}
