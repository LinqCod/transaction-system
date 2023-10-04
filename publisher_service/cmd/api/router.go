package api

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/linqcod/transaction-system/publisher_service/internal/http/handler"
	"github.com/linqcod/transaction-system/publisher_service/internal/jetstream"
	"github.com/linqcod/transaction-system/publisher_service/internal/model"
	"github.com/linqcod/transaction-system/publisher_service/internal/repository"
	"go.uber.org/zap"
)

func InitRouter(db *sql.DB, logger *zap.SugaredLogger) *gin.Engine {
	router := gin.Default()

	//init nats, handler and group endpoints
	js, err := jetstream.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to nats jetstream: %v", err)
	}
	if err = jetstream.CreateStream(js); err != nil {
		logger.Fatalf("error while creating stream: %v", err)
	}

	accountRepo := repository.NewAccountRepository(context.Background(), db)

	transactionHandler := handler.NewTransactionHandler(logger, js, accountRepo)

	api := router.Group("/api/v1")
	{
		accounts := api.Group("/accounts")
		{
			accounts.POST("/invoice", transactionHandler.HandleTransaction(model.InvoiceTransactionType))
			accounts.POST("/withdraw", transactionHandler.HandleTransaction(model.WithdrawTransactionType))
			accounts.GET("/", transactionHandler.GetBalances)
		}
	}

	return router
}
