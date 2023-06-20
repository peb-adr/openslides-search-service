// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

// Package config implements the configuration of the search service.
package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Default configuration.
const (
	DefaultWebPort       = 9050
	DefaultWebHost       = ""
	DefaultMaxQueue      = 5
	DefaultIndexAge      = 100 * time.Millisecond
	DefaultIndexFile     = "search.bleve"
	DefaultIndexUpdate   = 2 * time.Minute
	DefaultIndexBatch    = 4096
	DefaultModels        = "models.yml"
	DefaultSearch        = "search.yml"
	DefaultDB            = "openslides"
	DefaultDBUser        = "openslides"
	DefaultSecretsPath   = "/run/secrets"
	DefaultDBPassword    = "secret:postgres_password"
	DefaultDBHost        = "localhost"
	DefaultDBPort        = 5432
	DefaultRestricterURL = ""
)

// Web are the parameters for the web server.
type Web struct {
	Port     int
	Host     string
	MaxQueue int
}

// Index are the parameters for the indexer.
type Index struct {
	File   string
	Age    time.Duration
	Update time.Duration
	Batch  int
}

// Models are the paths to the YAML files containing the models
// and the searched collections.
type Models struct {
	Models string
	Search string
}

// Database are the credentials for the datavbase.
type Database struct {
	Database string
	User     string
	Password string
	Host     string
	Port     int
}

// Config is the configuration of the search service.
type Config struct {
	SecretsPath string
	Web         Web
	Index       Index
	Models      Models
	Database    Database
	Restricter  Restricter
}

// Restricter is the URL of the restricter to filter content by user id.
type Restricter struct {
	URL string
}

// GetConfig returns the configuration overwritten with env vars.
func GetConfig() (*Config, error) {
	cfg := &Config{
		SecretsPath: DefaultSecretsPath,
		Web: Web{
			Port: DefaultWebPort,
			Host: DefaultWebHost,
		},
		Index: Index{
			File:   DefaultIndexFile,
			Age:    DefaultIndexAge,
			Update: DefaultIndexUpdate,
			Batch:  DefaultIndexBatch,
		},
		Models: Models{
			Models: DefaultModels,
			Search: DefaultSearch,
		},
		Database: Database{
			Database: DefaultDB,
			User:     DefaultDBUser,
			Password: DefaultDBPassword,
			Host:     DefaultDBHost,
			Port:     DefaultDBPort,
		},
		Restricter: Restricter{
			URL: DefaultRestricterURL,
		},
	}
	if err := cfg.fromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// fromEnv fills the config from env vars.
func (cfg *Config) fromEnv() error {
	var (
		storeString   = store(noparse)
		storeInt      = store(strconv.Atoi)
		storeDuration = store(parseDuration)
		storeSecret   = store(parseSecrets(&cfg.SecretsPath))
	)
	return storeFromEnv([]storeEnv{
		{"SECRETS_PATH", storeString(&cfg.SecretsPath)},
		{"OPENSLIDES_SEARCH_PORT", storeInt(&cfg.Web.Port)},
		{"OPENSLIDES_SEARCH_HOST", storeString(&cfg.Web.Host)},
		{"OPENSLIDES_SEARCH_MAX_QUEUED", storeInt(&cfg.Web.MaxQueue)},
		{"OPENSLIDES_SEARCH_INDEX_AGE", storeDuration(&cfg.Index.Age)},
		{"OPENSLIDES_SEARCH_INDEX_FILE", storeString(&cfg.Index.File)},
		{"OPENSLIDES_SEARCH_INDEX_BATCH", storeInt(&cfg.Index.Batch)},
		{"OPENSLIDES_SEARCH_INDEX_UPDATE_INTERVAL", storeDuration(&cfg.Index.Update)},
		{"OPENSLIDES_MODELS_YML", storeString(&cfg.Models.Models)},
		{"OPENSLIDES_SEARCH_YML", storeString(&cfg.Models.Search)},
		{"OPENSLIDES_DB", storeString(&cfg.Database.Database)},
		{"OPENSLIDES_DB_USER", storeString(&cfg.Database.User)},
		{"OPENSLIDES_DB_PASSWORD", storeSecret(&cfg.Database.Password)},
		{"OPENSLIDES_DB_HOST", storeString(&cfg.Database.Host)},
		{"OPENSLIDES_DB_PORT", storeInt(&cfg.Database.Port)},
		{"OPENSLIDES_RESTRICTER", storeString(&cfg.Restricter.URL)},
	})
}

// pgEncode encodes a string to be used in the postgres key value style.
// See: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func pgEncode(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}

// ConnectionConfig returns a postgres Keyword/Value Connection String.
func (db *Database) ConnectionConfig() string {
	return fmt.Sprintf("user='%s' password='%s' host='%s' port='%d' dbname='%s'",
		pgEncode(db.User),
		pgEncode(db.Password),
		pgEncode(db.Host),
		db.Port,
		pgEncode(db.Database))
}
