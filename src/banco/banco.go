package banco

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //Driver
)

func ConectarBanco() (*sql.DB, error) {
	StringConexaoBanco := "root:6947lhblRMTY@/campanha_golang?charset=utf8&parseTime=True&loc=Local"
	db, erro := sql.Open("mysql", StringConexaoBanco)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}