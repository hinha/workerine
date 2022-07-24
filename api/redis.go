package api

import (
	"crypto/tls"
	"strings"

	"github.com/hibiken/asynq"

	"github.com/hinha/workerine/server/config"
)

func redisConnOpt(cfg *config.Config) (asynq.RedisConnOpt, error) {
	// Connecting to redis-cluster
	if len(cfg.Redis.RedisClusterNodes) > 0 {
		return asynq.RedisClusterClientOpt{
			Addrs:     strings.Split(cfg.Redis.RedisClusterNodes, ","),
			Password:  cfg.Redis.Password,
			TLSConfig: makeTLSConfig(cfg.Redis),
		}, nil
	}

	// Connecting to redis-sentinels
	if strings.HasPrefix(cfg.Redis.RedisURL, "redis-sentinel") {
		res, err := asynq.ParseRedisURI(cfg.Redis.RedisURL)
		if err != nil {
			return nil, err
		}
		connOpt := res.(asynq.RedisFailoverClientOpt) // safe to type-assert
		connOpt.TLSConfig = makeTLSConfig(cfg.Redis)
		return connOpt, nil
	}

	// Connecting to single redis server
	var connOpt asynq.RedisClientOpt
	if len(cfg.Redis.RedisURL) > 0 {
		res, err := asynq.ParseRedisURI(cfg.Redis.RedisURL)
		if err != nil {
			return nil, err
		}
		connOpt = res.(asynq.RedisClientOpt) // safe to type-assert
	} else {
		connOpt.Addr = cfg.Redis.Addr
		connOpt.DB = int(cfg.Redis.DB)
		connOpt.Password = cfg.Redis.Password
	}
	if connOpt.TLSConfig == nil {
		connOpt.TLSConfig = makeTLSConfig(cfg.Redis)
	}
	return connOpt, nil
}

func makeTLSConfig(cfg config.Redis) *tls.Config {
	if cfg.RedisTLS == "" && !cfg.RedisInsecureTLS {
		return nil
	}
	return &tls.Config{
		ServerName:         cfg.RedisTLS,
		InsecureSkipVerify: cfg.RedisInsecureTLS,
	}
}
