//go:build !k8s

package config

var Config = config{
	Db: DbConfig{
		DSN: "root:09161549@tcp(localhost:3306)/webook",
	},
	Redis: RedisConfig{
		Add: "localhost:6379",
	},
	EnableLimit: false,
}
