package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/msanath/mrds/controlplane"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/runtime/kind"
	temporalclient "go.temporal.io/sdk/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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

	tc, err := temporalclient.Dial(temporalclient.Options{
		HostPort:  "localhost:7233",
		Namespace: "mrds",
		Logger:    log,
	})
	if err != nil {
		return err
	}

	// Create a kind runtime activity.
	kubeconfig := clientcmd.RecommendedHomeFile
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	// Create a clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	kindRuntime := kind.NewKindRuntime(
		mrdspb.NewMetaInstancesClient(conn),
		mrdspb.NewDeploymentPlansClient(conn),
		mrdspb.NewNodesClient(conn),
		clientset,
	)

	cp := controlplane.NewTemporalControlPlane(conn, tc, kindRuntime)

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
