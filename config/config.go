package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	Database      DatabaseConfig
	Elasticsearch ElasticsearchConfig
	Kafka         KafkaConfig
}

type AppConfig struct {
	Name    string `validate:"required"`
	Env     string `validate:"required,oneof=production development"`
	Port    int    `validate:"required"`
	BaseURL string `validate:"required"`
}

type DatabaseConfig struct {
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Pool     PoolConfig
}

type PoolConfig struct {
	MaxIdleConns int `validate:"required"`
	MaxOpenConns int `validate:"required"`
	ConnIdleTime int `validate:"required"`
	ConnLifeTime int `validate:"required"`
}

type ElasticsearchConfig struct {
	Addresses []string `validate:"required,dive,required,http_url"`

	// Index Names
	UserIndex    string `validate:"required"`
	ProductIndex string `validate:"required"`
	OrderIndex   string `validate:"required"`
}

type KafkaConfig struct {
	GroupId string   `validate:"required"`
	Brokers []string `validate:"required,dive,required,hostname_port"`
	Topic   string   `validate:"required"`
}

func Load(validate *validator.Validate, log *logrus.Logger) *Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath("./..")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Panic("Error reading .env file")
	}

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", 3000)

	viper.SetDefault("DATABASE_PORT", 3036)

	cfg := new(Config)
	cfg.App.Name = viper.GetString("APP_NAME")
	cfg.App.Env = viper.GetString("APP_ENV")
	cfg.App.Port = viper.GetInt("APP_PORT")
	cfg.App.BaseURL = viper.GetString("APP_BASE_URL")

	cfg.Database.Name = viper.GetString("DATABASE_NAME")
	cfg.Database.User = viper.GetString("DATABASE_USER")
	cfg.Database.Password = viper.GetString("DATABASE_PASSWORD")
	cfg.Database.Host = viper.GetString("DATABASE_HOST")
	cfg.Database.Port = viper.GetInt("DATABASE_PORT")

	cfg.Database.Pool.MaxIdleConns = viper.GetInt("POOL_MAX_IDLE_CONNS")
	cfg.Database.Pool.MaxOpenConns = viper.GetInt("POOL_MAX_OPEN_CONNS")
	cfg.Database.Pool.ConnIdleTime = viper.GetInt("POOL_CONN_IDLE_TIME")
	cfg.Database.Pool.ConnLifeTime = viper.GetInt("POOL_CONN_LIFE_TIME")

	cfg.Elasticsearch.Addresses = viper.GetStringSlice("ELASTICSEARCH_ADDRESSES")
	cfg.Elasticsearch.UserIndex = viper.GetString("ELASTICSEARCH_USER_INDEX")
	cfg.Elasticsearch.ProductIndex = viper.GetString("ELASTICSEARCH_PRODUCT_INDEX")
	cfg.Elasticsearch.OrderIndex = viper.GetString("ELASTICSEARCH_ORDER_INDEX")

	cfg.Kafka.GroupId = viper.GetString("KAFKA_GROUP_ID")
	cfg.Kafka.Brokers = viper.GetStringSlice("KAFKA_BROKERS")
	cfg.Kafka.Topic = viper.GetString("KAFKA_TOPIC")

	err = validate.Struct(cfg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			log.Error(validationError.Error())
		}
		os.Exit(1)
	}

	return cfg
}
