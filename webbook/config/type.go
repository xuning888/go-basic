package config

type config struct {
	Db          DbConfig
	Redis       RedisConfig
	EnableLimit bool
}

type DbConfig struct {
	DSN string
}

type RedisConfig struct {
	Add string
}
