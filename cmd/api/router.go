package api

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(ctx context.Context, logger *zap.SugaredLogger, db *sql.DB) *gin.Engine {
	router := gin.Default()

	// TODO: init service, repo, handler and group endpoints

	return router
}
