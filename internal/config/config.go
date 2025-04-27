package config

import (
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    Port         string `envconfig:"PORT" default:"8080"`
    PostgresURL  string `envconfig:"POSTGRES_URL" required:"true"`
    MongoURL     string `envconfig:"MONGO_URL" required:"true"`
    MongoDBName  string `envconfig:"MONGO_DB" default:"webhooks"`
    RedisURL     string `envconfig:"REDIS_URL" required:"true"`
    RabbitURL    string `envconfig:"RABBIT_URL" required:"true"`
    DeliveryTTL  int    `envconfig:"DELIVERY_TTL" default:"259200"`
}

func Load() (*Config, error) {
    var cfg Config
    err := envconfig.Process("", &cfg)
    return &cfg, err
}
