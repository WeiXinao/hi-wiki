package ioc

import (
	"github.com/WeiXinao/hi-wiki/config"
	"github.com/WeiXinao/hi-wiki/pkg/redisx"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

func InitRedis(config *config.AppConfig,
	promHook *redisx.PrometheusHook,
	otelHook *redisx.OtelHook) redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis.Addr,
	})
	redisClient.AddHook(promHook)
	redisClient.AddHook(otelHook)
	return redisClient
}

func InitRedisPrometheus() *redisx.PrometheusHook {
	return redisx.NewPrometheusHook(prom.SummaryOpts{
		Namespace: "xiaoxin",
		Subsystem: "hi_wiki",
		Name:      "gin_redis",
		Help:      "统计redis的缓存命令",
		ConstLabels: map[string]string{
			"instance_id": "hi_wiki_1",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
}

func InitRedisOtel() *redisx.OtelHook {
	tracer := otel.Tracer("github.com/WeiXinao/hi-wiki/redisx")
	return redisx.NewOtelHook(tracer)
}
