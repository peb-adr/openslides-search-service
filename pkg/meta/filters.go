package meta

import (
	log "github.com/sirupsen/logrus"

	"github.com/goccy/go-yaml"
)

// Filter is part of the meta model.
type Filter struct {
	Name        string
	Items       []string
	ItemsConfig map[string]*CollectionSearchableConfig
	Additional  []string
	Contains    map[string]struct{}
	Relations   map[string]*CollectionRelation
}

// Filters is a list of filters.
type Filters []Filter

// UnmarshalYAML Parses yaml to Filters
func (fs *Filters) UnmarshalYAML(node []byte) error {
	var fsm map[string]CollectionDescription
	if err := yaml.Unmarshal(node, &fsm); err != nil {
		return err
	}

	*fs = make(Filters, 0, len(fsm))
	for k := range fsm {
		relations := map[string]*CollectionRelation{}
		for k, r := range fsm[k].Relations {
			relations[k] = r
		}

		contains := make(map[string]struct{}, len(fsm[k].Contains))
		for _, c := range fsm[k].Contains {
			contains[c] = struct{}{}
		}

		*fs = append(*fs, Filter{
			Name:        k,
			Items:       fsm[k].Searchable,
			ItemsConfig: fsm[k].SearchableConfig,
			Additional:  fsm[k].Additional,
			Relations:   relations,
			Contains:    contains,
		})
	}
	return nil
}

// ContainmentMap returns a map with info which collections can be found within another collection
func (fs Filters) ContainmentMap() map[string]map[string]struct{} {
	containment := map[string]map[string]struct{}{}
	for _, f := range fs {
		containment[f.Name] = map[string]struct{}{}
		containment[f.Name][f.Name] = struct{}{}
		for _, g := range fs {
			if _, ok := g.Contains[f.Name]; ok {
				containment[f.Name][g.Name] = struct{}{}
			}
		}
	}

	return containment
}

// Retain returns a keep function for [Retain] which also updates
// if Members are searchable and adds their relation informations
func (fs Filters) Retain(verbose bool) func(string, string, *Member) bool {
	type key struct {
		rel   string
		field string
	}
	keep := map[key]struct{}{}
	additional := map[key]struct{}{}
	relations := map[key]*CollectionRelation{}
	config := map[key]*CollectionSearchableConfig{}
	for _, m := range fs {
		for _, f := range m.Items {
			keep[key{rel: m.Name, field: f}] = struct{}{}
		}

		for f, data := range m.ItemsConfig {
			config[key{rel: m.Name, field: f}] = data
		}

		for _, f := range m.Additional {
			additional[key{rel: m.Name, field: f}] = struct{}{}
		}

		for f, data := range m.Relations {
			relations[key{rel: m.Name, field: f}] = data
		}
	}
	return func(rk, fk string, m *Member) bool {
		if _, ok := relations[key{rel: rk, field: fk}]; ok {
			m.Relation = relations[key{rel: rk, field: fk}]
		}

		if c, ok := config[key{rel: rk, field: fk}]; ok {
			if c.Type != nil {
				m.Type = *c.Type
			}

			m.Analyzer = c.Analyzer
		}

		if _, ok := additional[key{rel: rk, field: fk}]; ok {
			m.Searchable = false
			return true
		}

		_, ok := keep[key{rel: rk, field: fk}]
		if !ok && verbose {
			log.Printf("removing filtered %s.%s\n", rk, fk)
		} else {
			m.Searchable = true
		}

		return ok
	}
}
