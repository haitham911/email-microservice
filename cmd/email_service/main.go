package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/email-microservice/config"
	gw "github.com/email-microservice/internal/email/proto"
	"github.com/email-microservice/internal/server"
	"github.com/email-microservice/pkg/jaeger"
	"github.com/email-microservice/pkg/logger"
	"github.com/email-microservice/pkg/mailer"
	"github.com/email-microservice/pkg/postgres"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	ot "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	prefixTracerState  = "x-b3-"
	zipkinTraceID      = prefixTracerState + "traceid"
	zipkinSpanID       = prefixTracerState + "spanid"
	zipkinParentSpanID = prefixTracerState + "parentspanid"
	zipkinSampled      = prefixTracerState + "sampled"
	zipkinFlags        = prefixTracerState + "flags"
)

var otHeaders = []string{
	zipkinTraceID,
	zipkinSpanID,
	zipkinParentSpanID,
	zipkinSampled,
	zipkinFlags}

func main() {
	log.Println("Starting server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer psqlDB.Close()

	appLogger.Infof("PostgreSQL connected: %#v", psqlDB.Stats())

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	ot.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")
	go func() {
		Gateway()
	}()
	mailDialer := mailer.NewMailDialer(cfg)
	appLogger.Info("Mail dialer connected")

	s := server.NewEmailsServer(appLogger, cfg, mailDialer, psqlDB)

	appLogger.Fatal(s.Run())
}

func Gateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	annotators := []annotator{injectHeadersIntoMetadata}

	ropts := []runtime.ServeMuxOption{
		runtime.WithMetadata(chainGrpcAnnotators(annotators...)),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
	}

	mux := runtime.NewServeMux(ropts...)
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_opentracing.StreamClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer()))))
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_opentracing.UnaryClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer()))))
	opts = append(opts, grpc.WithInsecure())
	// users
	echoEndpoint := "localhost:5000"
	err := gw.RegisterEmailServiceHandlerFromEndpoint(ctx, mux, echoEndpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting gateway server on port 8001")
	http.ListenAndServe(":8001", mux)
}

type annotator func(context.Context, *http.Request) metadata.MD

func chainGrpcAnnotators(annotators ...annotator) annotator {
	return func(c context.Context, r *http.Request) metadata.MD {
		mds := []metadata.MD{}
		for _, a := range annotators {
			mds = append(mds, a(c, r))
		}
		return metadata.Join(mds...)
	}
}

func injectHeadersIntoMetadata(ctx context.Context, req *http.Request) metadata.MD {
	pairs := []string{}
	for _, h := range otHeaders {
		if v := req.Header.Get(h); len(v) > 0 {
			pairs = append(pairs, h, v)
		}

	}
	return metadata.Pairs(pairs...)
}
