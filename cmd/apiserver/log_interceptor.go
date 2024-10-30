package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/msanath/gondolf/pkg/printer"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Unary Logging Interceptor for logging gRPC requests and responses
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	logger := slog.New(slog.NewTextHandler(os.Stdout, ctxslog.NewCustomHandler(slog.LevelInfo)))
	ctx = ctxslog.NewContext(ctx, logger)

	jsonReq, err := marshalJSON(req.(proto.Message))
	if err != nil {
		logger.Error("Failed to marshal request", "error", err)
	}

	logger.Info("Incoming request",
		"method", info.FullMethod)
	// "request", jsonReq)

	// Custom printing for testing
	fmt.Println(printer.BlueText(jsonReq))

	// Call the handler to process the request
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error("Failed to process request", "error", err)
		return resp, err
	}

	jsonResp, err := marshalJSON(resp.(proto.Message))
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		return resp, err
	}

	logger.Info("Outgoing response", "method", info.FullMethod, "response", err, "duration", time.Since(start))
	fmt.Println(printer.GreenText(jsonResp))

	return resp, err
}

func marshalJSON(pb proto.Message) (string, error) {
	b, err := protojson.MarshalOptions{EmitUnpopulated: true, Multiline: true}.Marshal(pb)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
