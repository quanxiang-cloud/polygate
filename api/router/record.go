package router

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	client "github.com/quanxiang-cloud/polygate/pkg/kafka"
)

// Record Record
type Record struct {
	ctx        context.Context
	client     *client.Client
	recordChan chan *client.Record
}

// NewRecord NewRecord
func NewRecord(conf *config.Config) *Record {
	producer, err := client.NewSyncProducer(conf.Kafka)
	if err != nil {
		panic(err)
	}
	r := &Record{
		ctx:        context.Background(),
		client:     client.New(producer, conf.Handler),
		recordChan: make(chan *client.Record, conf.Handler.Buffer),
	}

	for i := 0; i < conf.Handler.NumOfProcessor; i++ {
		go r.process(r.ctx)
	}
	return r
}

func (r *Record) record() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		record := &client.Record{
			RequestID:     c.Request.Header.Get(header.HeaderRequestID),
			UserID:        c.Request.Header.Get(header.HeaderUserID),
			UserName:      c.Request.Header.Get(header.HeaderUserName),
			OperationTime: start.UnixNano() / 1e6,
			OperationUA:   c.Request.Header.Get("User-Agent"),
			OperationType: c.Request.Method,
			IP:            c.ClientIP(),
			Detail:        fmt.Sprintf("status: %d, path: %s", c.Writer.Status(), path),
		}
		r.recordChan <- record
	}
}

func (r *Record) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case entity := <-r.recordChan:
			err := r.client.Send(entity)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
