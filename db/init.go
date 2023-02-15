package db

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/rroy233/logger"
	"github.com/rroy233/tg-stickers-dl/config"
	"github.com/rroy233/tg-stickers-dl/languages"
)

// var db *sqlx.DB
var rdb *redis.Client
var ctx = context.Background()

func Init() {
	var tlsConfig *tls.Config
	if config.Get().Redis.TLS == true {
		tlsConfig = &tls.Config{
			ServerName: config.Get().Redis.Server,
		}
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:      config.Get().Redis.Server + ":" + config.Get().Redis.Port,
		Password:  config.Get().Redis.Password,
		DB:        config.Get().Redis.DB,
		TLSConfig: tlsConfig,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		logger.FATAL.Fatalln(languages.Get(nil).System.DbRedisStartFailed, err)
		return
	}
	logger.Info.Println(languages.Get(nil).System.DbRedisConnected)

	//queue
	initQueue(config.Get().General.ProcessWaitQueueMaxSize)

	//cache
	initCache()

	return
}

func Close() {
	if err := rdb.Close(); err != nil {
		logger.Error.Println("Redis close Error:", err)
	}
}
