package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	errorTypes "github.com/linqcod/transaction-system/publisher_service/internal/common"
	"github.com/linqcod/transaction-system/publisher_service/internal/http/dto"
	"github.com/linqcod/transaction-system/publisher_service/internal/jetstream"
	"github.com/linqcod/transaction-system/publisher_service/internal/model"
	"github.com/linqcod/transaction-system/publisher_service/pkg/currencyapi"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"net/http"
)

type AccountRepository interface {
	UpdateAccountFrozenBalance(card string, newFrozenBalance float64) error
	GetAccount(card string) (*model.Account, error)
	GetAccountsBalance() ([]*model.AccountBalance, error)
}

type TransactionHandler struct {
	logger     *zap.SugaredLogger
	js         nats.JetStreamContext
	repository AccountRepository
}

func NewTransactionHandler(logger *zap.SugaredLogger, js nats.JetStreamContext, repo AccountRepository) *TransactionHandler {
	return &TransactionHandler{
		logger:     logger,
		js:         js,
		repository: repo,
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

		if transactionType == model.InvoiceTransactionType {
			account, err := h.repository.GetAccount(transaction.CardNumber)
			if err != nil {
				h.logger.Errorf("error while getting account: %v", err)
				c.JSON(http.StatusBadRequest, dto.ErrorDTO{
					Error: errorTypes.ErrAccountNotFound.Error(),
				})
				return
			}

			currencyCoefficient, err := currencyapi.ConvertCurrencies(transaction.Currency, "RUB")
			if err != nil {
				h.logger.Errorf("error while converting currencies: %v", err)
				c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
					Error: errorTypes.ErrConvertingCurrencies.Error(),
				})
				return
			}

			newFrozenBalance := account.FrozenBalance + transaction.Amount*currencyCoefficient

			err = h.repository.UpdateAccountFrozenBalance(transaction.CardNumber, newFrozenBalance)
			if err != nil {
				h.logger.Errorf("error while updating account frozen balance: %v", err)
				c.JSON(http.StatusBadRequest, dto.ErrorDTO{
					Error: errorTypes.ErrUpdatingAccountFrozenBalance.Error(),
				})
				return
			}
		}

		transaction.Type = transactionType
		transaction.Status = model.CreatedTransactionStatus

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

func (h TransactionHandler) GetBalances(c *gin.Context) {
	balances, err := h.repository.GetAccountsBalance()
	if err != nil {
		h.logger.Errorf("error while getting accounts: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Error: errorTypes.ErrGettingAccountsBalance.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, balances)
}
