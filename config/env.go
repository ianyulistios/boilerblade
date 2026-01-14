package config

type Env struct {
	MODE           string `envconfig:"MODE" default:"development"`
	FIBER_PORT     string `envconfig:"FIBER_PORT" default:"3000"`
	FIBER_APP_NAME string `envconfig:"FIBER_APP_NAME" default:"boilerblade"`
	APP_KEY        string `envconfig:"APP_KEY" default:""`
	SERVER_MODE    string `envconfig:"SERVER_MODE" default:"both"` // http, amqp, or both

	// Connection enable flags
	ENABLE_DB    bool `envconfig:"ENABLE_DB" default:"true"`
	ENABLE_REDIS bool `envconfig:"ENABLE_REDIS" default:"true"`
	ENABLE_AMQP  bool `envconfig:"ENABLE_AMQP" default:"true"`

	DB_TYPE               string `envconfig:"DB_TYPE" default:"postgres"` // postgres or mysql
	DB_HOST               string `envconfig:"DB_HOST" default:"localhost"`
	DB_PORT               string `envconfig:"DB_PORT" default:"5432"`
	DB_USER               string `envconfig:"DB_USER" default:"postgres"`
	DB_PASSWORD           string `envconfig:"DB_PASSWORD" default:"postgres"`
	DB_NAME               string `envconfig:"DB_NAME" default:"boilerblade"`
	DB_MAX_OPEN_CONNS     int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	DB_MAX_IDLE_CONNS     int    `envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
	DB_MAX_LIFETIME_CONNS int    `envconfig:"DB_MAX_LIFETIME_CONNS" default:"10"`

	REDIS_HOST     string `envconfig:"REDIS_HOST" default:"localhost"`
	REDIS_PORT     string `envconfig:"REDIS_PORT" default:"6379"`
	REDIS_PASSWORD string `envconfig:"REDIS_PASSWORD" default:""`
	REDIS_DB       int    `envconfig:"REDIS_DB" default:"0"`

	AMQP_HOST     string `envconfig:"AMQP_HOST" default:"localhost"`
	AMQP_PORT     string `envconfig:"AMQP_PORT" default:"5672"`
	AMQP_USER     string `envconfig:"AMQP_USER" default:"guest"`
	AMQP_PASSWORD string `envconfig:"AMQP_PASSWORD" default:"guest"`
}
