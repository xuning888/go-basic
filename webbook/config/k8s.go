//go:build k8s

package config

var Config = config{
	Db: DbConfig{
		DSN: "root:root@tcp(webook-mysql:13309)/webook",
	},
	Redis: RedisConfig{
		Add: "webook-redis:11479",
	},
}
