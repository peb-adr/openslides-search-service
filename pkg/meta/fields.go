package meta

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

// Fields is part of the meta model.
type Fields struct {
	Type string `yaml:"type"`
	To   string `yaml:"to"`
}

// UnmarshalYAML Parses yaml to Fields
func (fs *Fields) UnmarshalYAML(node []byte) error {
	var s string
	if err := yaml.Unmarshal(node, &s); err == nil {
		fs.Type = s
		return nil
	}
	var field struct {
		Type string `yaml:"type"`
		To   string `yaml:"to"`
	}
	if err := yaml.Unmarshal(node, &field); err != nil {
		return fmt.Errorf("fields object without type: %w", err)
	}
	fs.Type = field.Type
	fs.To = field.To
	return nil
}

// Clone returns a deep copy.
func (fs *Fields) Clone() *Fields {
	if fs == nil {
		return nil
	}
	return &Fields{
		Type: fs.Type,
		To:   fs.To,
	}
}
