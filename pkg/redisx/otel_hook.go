package redisx

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net"
	"strings"
)

type OtelHook struct {
	tracer trace.Tracer
}

func NewOtelHook(tracer trace.Tracer) *OtelHook {
	return &OtelHook{
		tracer: tracer,
	}
}

func (o *OtelHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (o *OtelHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		ctx, span := o.tracer.Start(ctx, "redis.process")
		var err error
		defer func() {
			span.AddEvent(cmd.Name())
			keyExists := errors.Is(err, redis.Nil)
			strBuilder := strings.Builder{}
			for _, arg := range cmd.Args() {
				strBuilder.WriteString(fmt.Sprintf("%v", arg))
				strBuilder.WriteString(" ")
			}
			span.SetAttributes(attribute.String("command.full", strBuilder.String()))
			if err != nil && !keyExists {
				span.SetAttributes(attribute.String("err", err.Error()))
				span.End()
			}
			span.SetAttributes(attribute.Bool("key.exists", !keyExists))
			span.End()
		}()
		err = next(ctx, cmd)
		return err
	}
}

func (o *OtelHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
