//go:build k8s

package config

var Config = config{
	Db: DbConfig{
		DSN: "root:root@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Add: "webook-redis:6380",
	},
	EnableLimit: true,
}
