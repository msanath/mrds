package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/msanath/mrds/pkg/controlplane"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serverOptions struct{}

func main() {
	so := serverOptions{}
	cmd := cobra.Command{
		Use: "mrds-controlplane",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewTextHandler(os.Stdout, ctxslog.NewCustomHandler(slog.LevelInfo)))
			ctx := ctxslog.NewContext(cmd.Context(), logger)
			return so.Run(ctx)
		},
	}

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func (o serverOptions) Run(ctx context.Context) error {
	log := ctxslog.FromContext(ctx)

	log.Info("Starting control plane")
	conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return err
	}

	cp := controlplane.NewTemporalControlPlane(controlplane.ControlPlaneOptions{
		MRDSConn: conn,
	})

	cpErrChan := make(chan error)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := cp.Start(ctx)
		if err != nil {
			cpErrChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutting down")
	case err := <-cpErrChan:
		log.Error("Control plane error", "error", err)
	}

	cancel()
	return nil
}
