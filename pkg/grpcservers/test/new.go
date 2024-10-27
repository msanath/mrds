package test

import (
	"context"
	"fmt"
	"net"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/cluster"
	"github.com/msanath/mrds/internal/sqlstorage"
	"github.com/msanath/mrds/pkg/grpcservers"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/node"
	// ++ledgerbuilder:Imports

	"github.com/msanath/gondolf/pkg/simplesql/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type TestServer struct {
	server *grpc.Server
	conn   *grpc.ClientConn
}

var testDb = test.NewTestSQLiteDB

// var testDb = test.NewTestMySQLDB

func NewTestServer() (*TestServer, error) {
	gServer := grpc.NewServer()

	db, err := testDb()
	if err != nil {
		return nil, fmt.Errorf("failed to create test sqlite db: %w", err)
	}
	storage, err := sqlstorage.NewSQLStorage(db, false)
	if err != nil {
		return nil, err
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
	// ++ledgerbuilder:TestServerRegister

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := gServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	//nolint:staticcheck
	conn, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %w", err)
	}

	return &TestServer{
		server: gServer,
		conn:   conn,
	}, nil
}

func (s *TestServer) Conn() *grpc.ClientConn {
	return s.conn
}

func (s *TestServer) Close() {
	s.conn.Close()
	s.server.Stop()
}
