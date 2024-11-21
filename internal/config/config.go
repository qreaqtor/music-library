package config

import "time"

type Config struct {
	Addr Adress `yaml:"address"`

	Env string `yaml:"env"`

	Api ApiConfig `yaml:"api"`

	Postgres PostgresConfig `yaml:"postgres"`
}

type ApiConfig struct {
	Version int `yaml:"version"`
}

type PostgresConfig struct {
	User       string             `yaml:"user"`
	Password   string             `yaml:"password"`
	Database   string             `yaml:"database"`
	Addr       Adress             `yaml:"address"`
	Connection PostgresConnection `yaml:"connection"`
}

type PostgresConnection struct {
	Attempts int           `yaml:"attempts"`
	Delay    time.Duration `yaml:"delay"`
}
