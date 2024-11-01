package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/cluster"
	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	"github.com/msanath/mrds/internal/ledger/metainstance"
	"github.com/msanath/mrds/internal/ledger/node"
	"github.com/msanath/mrds/internal/sqlstorage"
	"github.com/msanath/mrds/pkg/grpcservers"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/msanath/gondolf/pkg/simplesql/test"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type serverOptions struct {
	testMode bool
}

func main() {
	so := serverOptions{}
	cmd := cobra.Command{
		Use: "mrds-apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewTextHandler(os.Stdout, ctxslog.NewCustomHandler(slog.LevelInfo)))
			ctx := ctxslog.NewContext(cmd.Context(), logger)
			return so.Run(ctx)
		},
	}
	cmd.Flags().BoolVar(&so.testMode, "test-mode", false, "Uses in-memory database. Data will be lost after server restart.")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func (o serverOptions) Run(ctx context.Context) error {
	log := ctxslog.FromContext(ctx)
	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		return err
	}

	gServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	var dbConn *sqlx.DB
	if o.testMode {
		dbConn, err = test.NewTestSQLiteDB()
		if err != nil {
			return err
		}
	} else {
		dbConn, err = newMySQLConn()
		if err != nil {
			return err
		}
	}

	storage, err := sqlstorage.NewSQLStorage(dbConn, false)
	if err != nil {
		return err
	}
	clusterLedger := cluster.NewLedger(storage.Cluster)
	mrdspb.RegisterClustersServer(
		gServer,
		grpcservers.NewClusterService(clusterLedger),
	)
	computeCapabilityLedger := computecapability.NewLedger(storage.ComputeCapability)
	mrdspb.RegisterComputeCapabilitiesServer(
		gServer,
		grpcservers.NewComputeCapabilityService(computeCapabilityLedger),
	)

	nodeLedger := node.NewLedger(storage.Node)
	mrdspb.RegisterNodesServer(
		gServer,
		grpcservers.NewNodeService(nodeLedger),
	)

	metaInstanceLedger := metainstance.NewLedger(storage.MetaInstance)
	mrdspb.RegisterMetaInstancesServer(
		gServer,
		grpcservers.NewMetaInstanceService(metaInstanceLedger),
	)

	deploymentPlanLedger := deploymentplan.NewLedger(storage.DeploymentPlan)
	mrdspb.RegisterDeploymentPlansServer(
		gServer,
		grpcservers.NewDeploymentPlanService(deploymentPlanLedger),
	)

	log.Info("Starting MRDS API server")
	return gServer.Serve(lis)
}

func newMySQLConn() (*sqlx.DB, error) {
	mysqlConfig := &mysql.Config{
		User:                 "root",
		Addr:                 "127.0.0.1:3306",
		DBName:               "mrds",
		Passwd:               "",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	// dsn := "root:@tcp(127.0.0.1:3306)/?parseTime=true"
	dsn := mysqlConfig.FormatDSN()
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
