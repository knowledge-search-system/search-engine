package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env           string              `yaml:"env"`
	GRPC          GRPCConfig          `yaml:"grpc"`
	HTTP          HTTPConfig          `yaml:"http"`
	Postgres      PostgresConfig      `yaml:"postgres"`
	Redis         RedisConfig         `yaml:"redis"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	Logger        LoggerConfig        `yaml:"logger"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN      string `yaml:"dsn"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type RedisConfig struct {
	Addr     string        `yaml:"addr"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	CacheTTL time.Duration `yaml:"cache_ttl"`
}

type ElasticsearchConfig struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	IndexName string   `yaml:"index_name"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Overload()

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SEARCH_ENGINE_GRPC_PORT"); v != "" {
		cfg.GRPC.Port = atoiOrDefault(v, cfg.GRPC.Port)
	}
	if v := os.Getenv("SEARCH_ENGINE_HTTP_PORT"); v != "" {
		cfg.HTTP.Port = atoiOrDefault(v, cfg.HTTP.Port)
	}
	if v := os.Getenv("SEARCH_ENGINE_POSTGRES_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
	if v := os.Getenv("SEARCH_ENGINE_REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("SEARCH_ENGINE_REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("SEARCH_ENGINE_ELASTICSEARCH_ADDRESSES"); v != "" {
		cfg.Elasticsearch.Addresses = strings.Split(v, ",")
	}
	if v := os.Getenv("SEARCH_ENGINE_ELASTICSEARCH_USERNAME"); v != "" {
		cfg.Elasticsearch.Username = v
	}
	if v := os.Getenv("SEARCH_ENGINE_ELASTICSEARCH_PASSWORD"); v != "" {
		cfg.Elasticsearch.Password = v
	}
	if v := os.Getenv("SEARCH_ENGINE_LOG_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
}

func atoiOrDefault(s string, def int) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return def
		}
		n = n*10 + int(c-'0')
	}
	return n
}
