package database

import (
    "context"
    "database/sql"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context, dsn string) (*sql.DB,error){
	db,err := sql.Open("pgx",dsn)

	if err != nil {
		return nil,err
	}

	err = db.PingContext(ctx)


	if err != nil{
		db.Close()
		return nil,err

	}
	return db,nil
}

