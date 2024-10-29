package sqlstorage

import (
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage/tables"

	"github.com/msanath/gondolf/pkg/simplesql"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/metainstance"
	"github.com/msanath/mrds/internal/ledger/node"

	"github.com/msanath/mrds/internal/ledger/cluster"

	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	// ++ledgerbuilder:Imports

	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	ComputeCapability computecapability.Repository
	Node              node.Repository
	MetaInstance      metainstance.Repository
	Cluster           cluster.Repository
	DeploymentPlan    deploymentplan.Repository
	// ++ledgerbuilder:RepositoryInterface
}

func NewSQLStorage(
	db *sqlx.DB,
	is_sqlite bool,
) (*SQLStorage, error) {

	var errHandler simplesql.ErrHandler
	if is_sqlite {
		errHandler = simplesql.SQLiteErrHandler
	} else {
		errHandler = simplesql.MySQLErrHandler
	}

	simpleDB := simplesql.NewDatabase(
		db, simplesql.WithErrHandler(errHandler),
	)
	err := tables.Initialize(simpleDB)
	if err != nil {
		return nil, err
	}

	return &SQLStorage{
		Cluster:           newClusterStorage(simpleDB),
		ComputeCapability: newComputeCapabilityStorage(simpleDB),
		Node:              newNodeStorage(simpleDB),
		MetaInstance:      newMetaInstanceStorage(simpleDB),
		DeploymentPlan:    newDeploymentPlanStorage(simpleDB),
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
