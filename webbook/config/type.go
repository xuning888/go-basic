package config

type config struct {
	Db    DbConfig
	Redis RedisConfig
}

type DbConfig struct {
	DSN string
}

type RedisConfig struct {
	Add string
}
