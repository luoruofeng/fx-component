package srv

import (
	"context"
	"time"

	c "github.com/luoruofeng/fx-component/conf"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// "go.uber.org/fx"

type RedisSrv struct {
	Cli *redis.Client
}

func NewRedisSrv(lc fx.Lifecycle, logger *zap.Logger) RedisSrv {
	var result RedisSrv = RedisSrv{}
	// 创建Redis客户端
	config := c.GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DbNumber,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
	})
	result.Cli = rdb

	// 设置Redis连接池
	// rdbOpts := &redis.Options{
	// 	Addr:     rdbOptions.Addr,
	// 	Password: rdbOptions.Password,
	// 	DB:       rdbOptions.DB,
	// }
	// rdbOpts.PoolSize = rdbOptions.PoolSize
	// rdbOpts.MinIdleConns = rdbOptions.MinIdleConns
	// rdbOpts.IdleTimeout = rdbOptions.IdleTimeout
	// rdbOpts.MaxConnAge = rdbOptions.MaxConnAge
	// rdbOpts.ReadTimeout = rdbOptions.ReadTimeout
	// rdbOpts.WriteTimeout = rdbOptions.WriteTimeout
	// rdbOpts.MaxRedirects = rdbOptions.MaxRedirects
	// rdbOpts.MinRetryBackoff = rdbOptions.MinRetryBackoff
	// rdbOpts.MaxRetryBackoff = rdbOptions.MaxRetryBackoff
	// rdbOpts.DialTimeout = rdbOptions.DialTimeout
	// rdbOpts.UseTLS = rdbOptions.UseTLS
	// rdbOpts.TLSConfig = rdbOptions.TLSConfig
	// rdbOpts.SkipVerify = rdbOptions.SkipVerify

	// pool := redis.NewClient(rdbOpts).Pool()
	// defer pool.Close()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Redis client...", zap.Any("rdb", rdb))
			s := rdb.Ping(ctx)
			if r, err := s.Result(); err != nil {
				logger.Error("Redis client start failed", zap.Any("err", err))
				panic(err)
			} else {
				logger.Info("Redis client start success", zap.Any("r", r))
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down Redis client...")
			return rdb.Close()
		},
	})
	return result
}
