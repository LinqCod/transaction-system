package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linqcod/transaction-system/publisher_service/internal/http/dto"
	"github.com/linqcod/transaction-system/publisher_service/internal/http/handler"
	"github.com/linqcod/transaction-system/publisher_service/internal/jetstream"
	"go.uber.org/zap"
)

func InitRouter(logger *zap.SugaredLogger) *gin.Engine {
	router := gin.Default()

	//init nats, handler and group endpoints
	js, err := jetstream.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to nats jetstream: %v", err)
	}
	if err = jetstream.CreateStream(js); err != nil {
		logger.Fatalf("error while creating stream: %v", err)
	}

	transactionHandler := handler.NewTransactionHandler(logger, js)

	api := router.Group("/api/v1")
	{
		transactions := api.Group("/accounts")
		{
			transactions.POST("/invoice", transactionHandler.HandleTransaction(dto.InvoiceTransactionType))
			transactions.POST("/withdraw", transactionHandler.HandleTransaction(dto.WithdrawTransactionType))
		}
	}

	return router
}
