package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	errorTypes "github.com/linqcod/transaction-system/publisher_service/internal/common"
	"github.com/linqcod/transaction-system/publisher_service/internal/http/dto"
	"github.com/linqcod/transaction-system/publisher_service/internal/jetstream"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"net/http"
)

type TransactionHandler struct {
	logger *zap.SugaredLogger
	js     nats.JetStreamContext
}

func NewTransactionHandler(logger *zap.SugaredLogger, js nats.JetStreamContext) *TransactionHandler {
	return &TransactionHandler{
		logger: logger,
		js:     js,
	}
}

func (h TransactionHandler) HandleTransaction(transactionType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction dto.TransactionDTO

		if err := json.NewDecoder(c.Request.Body).Decode(&transaction); err != nil {
			h.logger.Errorf("error while unmarshaling transaction body: %v", err)
			c.JSON(http.StatusBadRequest, dto.ErrorDTO{
				Error: errorTypes.ErrJSONUnmarshalling.Error(),
			})
			return
		}

		if transaction.Amount < 0 {
			h.logger.Error(errorTypes.ErrNegativeRefillAmount)
			c.JSON(http.StatusBadRequest, dto.ErrorDTO{
				Error: errorTypes.ErrNegativeRefillAmount.Error(),
			})
			return
		}

		transaction.Type = transactionType
		transaction.Status = dto.CreatedTransactionStatus

		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			h.logger.Errorf("error while marshaling transaction: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
				Error: errorTypes.ErrMarshallingToJSON.Error(),
			})
			return
		}

		_, err = h.js.Publish(jetstream.SubjectName, transactionJSON)
		if err != nil {
			h.logger.Errorf("error while publishing transaction: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
				Error: errorTypes.ErrNatsPublishing.Error(),
			})
			return
		}

		h.logger.Infof("transaction sent successfully: %v", transaction)

		c.JSON(http.StatusOK, transactionType+" transaction published successfully")
	}
}
