package sqlstorage

import (
	"github.com/msanath/mrds/internal/ledger/cluster"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"

	"github.com/msanath/gondolf/pkg/simplesql"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/node"
	// ++ledgerbuilder:Imports

	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	Cluster           cluster.Repository
	ComputeCapability computecapability.Repository
	Node              node.Repository
	// ++ledgerbuilder:RepositoryInterface
}

func NewSQLStorage(
	db *sqlx.DB,
	is_sqlite bool,
) (*SQLStorage, error) {
	var schemaMigrations = []simplesql.Migration{}

	schemaMigrations = append(schemaMigrations, clusterTableMigrations...)
	schemaMigrations = append(schemaMigrations, computeCapabilityTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeTableMigrations...)
	// ++ledgerbuilder:Migrations

	var errHandler simplesql.ErrHandler
	if is_sqlite {
		errHandler = simplesql.SQLiteErrHandler
	} else {
		errHandler = simplesql.MySQLErrHandler
	}

	simpleDB := simplesql.NewDatabase(
		db, simplesql.WithErrHandler(errHandler),
	)
	err := simpleDB.ApplyMigrations(schemaMigrations)
	if err != nil {
		return nil, err
	}

	return &SQLStorage{
		Cluster:           newClusterStorage(simpleDB),
		ComputeCapability: newComputeCapabilityStorage(simpleDB),
		Node:              newNodeStorage(simpleDB),
		// ++ledgerbuilder:RepoInstance
	}, nil
}

func errHandler(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case simplesql.ErrRecordNotFound:
		return ledgererrors.NewLedgerError(ledgererrors.ErrRecordNotFound, "Record not found.")
	case simplesql.ErrInsertConflict:
		return ledgererrors.NewLedgerError(ledgererrors.ErrRecordInsertConflict, "Duplicate entry, record already exists.")
	case simplesql.ErrInternal:
		return ledgererrors.NewLedgerError(ledgererrors.ErrRepositoryInternal, "Internal error.")
	default:
		return err
	}
}
