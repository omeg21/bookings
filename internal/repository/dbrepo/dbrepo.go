package dbrepo

import (
	"database/sql"

	"github.com/omeg21/project-repo/internal/config"
	"github.com/omeg21/project-repo/internal/repository"
)


type postgresDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo (conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo{
	return &postgresDBRepo{
		App: a,
		DB: conn,
	}
}