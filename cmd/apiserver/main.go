package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/cluster"
	"github.com/msanath/mrds/internal/sqlstorage"
	"github.com/msanath/mrds/pkg/grpcservers"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/msanath/gondolf/pkg/simplesql/test"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type serverOptions struct{}

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

	gServer := grpc.NewServer()
	// dbConn, err := newDBConn()
	// if err != nil {
	// 	return err
	// }
	dbConn, err := test.NewTestSQLiteDB()
	if err != nil {
		return err
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

	log.Info("Starting MRDS API server")
	return gServer.Serve(lis)
}

// func newDBConn() (*sqlx.DB, error) {
// 	mysqlConfig := &mysql.Config{
// 		User:                 "root",
// 		Addr:                 "127.0.0.1:3306",
// 		DBName:               "mrds",
// 		Passwd:               "",
// 		ParseTime:            true,
// 		AllowNativePasswords: true,
// 	}
// 	// dsn := "root:@tcp(127.0.0.1:3306)/?parseTime=true"
// 	dsn := mysqlConfig.FormatDSN()
// 	db, err := sqlx.Connect("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
