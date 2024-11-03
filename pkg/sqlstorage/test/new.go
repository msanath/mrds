// Code generated by ledger-builder. DO NOT EDIT.

package test

import (
	"testing"

	"github.com/msanath/mrds/pkg/sqlstorage"

	simplesqltest "github.com/msanath/gondolf/pkg/simplesql/test"
	"github.com/stretchr/testify/require"
)

func TestSQLStorage(t *testing.T) *sqlstorage.SQLStorage {
	// db, err := simplesqltest.NewTestMySQLDB()
	db, err := simplesqltest.NewTestSQLiteDB()
	require.NoError(t, err)
	// storage, err := sqlstorage.NewSQLStorage(db, false)
	storage, err := sqlstorage.NewSQLStorage(db, true)
	require.NoError(t, err)
	return storage
}