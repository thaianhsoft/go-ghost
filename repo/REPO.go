package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/thaianhsoft/go-ghost/changeset"
	"github.com/thaianhsoft/go-ghost/edge"
)

var dbInternalPointer *sql.DB

func NewRepoMYSQL(driver string, config *mysql.Config) *sql.DB {
	if dbInternalPointer != nil {
		return dbInternalPointer
	}
	dbInternalPointer, err := sql.Open("mysql", config.FormatDSN())
	if err == nil {
		return dbInternalPointer
	}
	return nil
}

type RepoService interface{
	OpenTransaction(ctx context.Context) (*sql.Tx, error)
	Save(ctx context.Context, schemaClass edge.ISchema, tx ...*sql.Tx) error
	Update(ctx context.Context, schemaClass edge.ISchema, tx...*sql.Tx) error
	GetById(ctx context.Context, id interface{}, tx ...*sql.Tx) error
	RemoveById(ctx context.Context, id interface{}, tx ...*sql.Tx) error
}

type CRUDRepoService struct {

}

func (C *CRUDRepoService) OpenTransaction(ctx context.Context) (*sql.Tx, error) {
	if dbInternalPointer != nil {
		newTx, err := dbInternalPointer.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err == nil {
			return newTx, nil
		} else {
			return nil, err
		}
	}
	return nil, fmt.Errorf("db connection isn't opened, retry")
}

func (C *CRUDRepoService) Save(ctx context.Context, cs changeset.ChangeSet, tx ...*sql.Tx) error {
	var sessionDb interface{} = dbInternalPointer // default

}

func (C *CRUDRepoService) Update(ctx context.Context, schemaClass edge.ISchema, tx ...*sql.Tx) error {
	panic("implement me")
}

func (C *CRUDRepoService) GetById(ctx context.Context, id interface{}, tx ...*sql.Tx) error {
	panic("implement me")
}

func (C *CRUDRepoService) RemoveById(ctx context.Context, id interface{}, tx ...*sql.Tx) error {
	panic("implement me")
}



