package microservice

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Microservice struct {
	*gin.Engine
	HealthCheck   healthcheck.Handler
	Prometheus    *ginprometheus.Prometheus
	MetricsRouter *gin.Engine
	MetricsPort   string
	HealthPort    string
}

func NewMicroservice(metricsPort string, healthPort string, requestLimit int, debug bool) *Microservice {
	router := gin.New()
	metricsRouter := gin.New()

	microService := Microservice{router, healthcheck.NewHandler(), ginprometheus.NewPrometheus("gin_metrics"), metricsRouter, metricsPort, healthPort}
	if requestLimit != 0 {
		router.Use(limit.MaxAllowed(requestLimit))
	}
	microService.Use(initAggregatedLogging())
	microService.Use(gin.Recovery())
	logrus.SetFormatter(&logrus.JSONFormatter{})

	microService.MetricsRouter.Use(gin.Recovery())
	microService.Prometheus.SetListenAddressWithRouter(fmt.Sprintf(":%s", microService.MetricsPort), metricsRouter)
	microService.Prometheus.Use(router)

	return &microService
}

func (ms *Microservice) StartHealthRouter() {
	go http.ListenAndServe(":"+ms.HealthPort, ms.HealthCheck)
}

func GetLogger(c *gin.Context) *logrus.Entry {
	ctxLogger, ok := c.Get("ctxLogger")
	if ok {
		return ctxLogger.(*logrus.Entry)
	}
	var logger *logrus.Entry
	log, found := c.Get("aggregate-logger")
	if found {
		logger = logrus.WithFields(logrus.Fields{})
		logger.Logger = log.(*logrus.Logger)
	}
	c.Set("ctxLogger", logger)
	return logger
}

func initAggregatedLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		logBuffer := ginlogrus.NewLogBuffer()
		reqLogger := &logrus.Logger{
			Out:       &logBuffer,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		}
		start := time.Now()
		path := c.Request.URL.Path
		c.Set("aggregate-logger", reqLogger)
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fields := logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency-ms": float64(latency) / float64(time.Millisecond),
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(time.RFC3339),
			"comment":    comment,
		}

		if len(c.Errors) > 0 {
			entry := logrus.StandardLogger().WithFields(fields)
			entry.Error(c.Errors.String())
		} else {
			logBuffer.StoreHeader("request-summary-info", fields)
			fmt.Fprintf(os.Stdout, logBuffer.String())
		}
	}
}
