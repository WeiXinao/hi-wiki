package ioc

import (
	"fmt"
	"github.com/WeiXinao/hi-wiki/config"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/hi-wiki/pkg/gormx"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	prom "github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

func InitDB(appConfig *config.AppConfig, l logger.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(appConfig.DB.DSN), &gorm.Config{
		Logger: gormLogger.New(gormLoggerFunc(l.Debug), gormLogger.Config{
			SlowThreshold:             time.Microsecond * 10,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  gormLogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "hi_wiki",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	cb := gormx.NewCallback(prom.SummaryOpts{
		Namespace: "xiaoxin",
		Subsystem: "hi_wiki",
		Name:      "gin_db",
		Help:      "统计 GORM 的数据库查询",
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

	err = db.Use(cb)
	if err != nil {
		panic(err)
	}

	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics(),
		tracing.WithDBName("hi_wiki")))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g("gorm", logger.String("message", fmt.Sprintf(msg, args...)))
}
