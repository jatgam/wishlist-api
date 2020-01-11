package microservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
		aggLogBuffer := newAggregateLogBuffer()
		reqLogger := &logrus.Logger{
			Out:       &aggLogBuffer,
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
			aggLogBuffer.StoreHeader("request-info", fields)
			fmt.Fprintf(os.Stdout, aggLogBuffer.String())
		}
	}
}

type aggregatelogBuffer struct {
	Buff     strings.Builder
	header   map[string]interface{}
	headerMU *sync.RWMutex
	MaxSize  uint
}

func newAggregateLogBuffer() aggregatelogBuffer {
	buffer := aggregatelogBuffer{
		headerMU: &sync.RWMutex{},
		MaxSize:  100000,
	}
	return buffer
}

func (b *aggregatelogBuffer) StoreHeader(key string, value interface{}) {
	b.headerMU.Lock()
	if b.header == nil {
		b.header = make(map[string]interface{})
	}
	b.header[key] = value
	b.headerMU.Unlock()
}

func (b *aggregatelogBuffer) Write(data []byte) (n int, err error) {
	newEntry := bytes.TrimSuffix(data, []byte("\n"))

	if len(newEntry)+b.Buff.Len() > int(b.MaxSize) {
		return 0, fmt.Errorf("write failed: buffer MaxSize = %d, current len = %d, attempted to write len = %d, data == %s", b.MaxSize, b.Buff.Len(), len(newEntry), newEntry)
	}
	return b.Buff.Write(append(newEntry, []byte(",")...))
}

func (b *aggregatelogBuffer) String() string {
	var str strings.Builder
	str.WriteString("{")
	if b.header != nil && len(b.header) != 0 {
		b.headerMU.RLock()
		hdr, err := json.Marshal(b.header)
		b.headerMU.RUnlock()
		if err != nil {
			fmt.Println("Error Marshaling aggregateLogBuffer JSON")
		}
		str.Write(hdr[1 : len(hdr)-1])
		str.WriteString(",")
	}
	str.WriteString("\"logEntries\":[" + strings.TrimSuffix(b.Buff.String(), ",") + "]")
	str.WriteString("}\n")
	return str.String()
}
