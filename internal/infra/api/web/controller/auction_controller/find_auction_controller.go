package auction_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	auctionData, err := u.auctionUseCase.FindAuctionById(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (u *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("productName")

	statusNumber, errConv := strconv.Atoi(status)
	if errConv != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate auction status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.auctionUseCase.FindAuctions(context.Background(),
		auction_usecase.AuctionStatus(statusNumber), category, productName)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (u *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	auctionData, err := u.auctionUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (u *AuctionController) FindExpiredAuctions(c *gin.Context) {

	// Chama o UseCase para buscar os leilões expirados
	auctions, err := u.auctionUseCase.FindExpiredAuctions(context.Background())
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	// Retorna os leilões expirados no formato de resposta
	c.JSON(http.StatusOK, auctions)
}

func (u *AuctionController) CloseExpiredAuctions(c *gin.Context) {
	log.Println("Iniciando fechamento de leilões expirados...")

	// Chama o UseCase para fechar os leilões expirados
	err := u.auctionUseCase.CloseExpiredAuctions(context.Background())
	if err != nil {
		log.Printf("Erro ao fechar leilões expirados: %v\n", err)
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	log.Println("Fechamento de leilões expirados concluído com sucesso.")

	// Retorna uma mensagem indicando que os leilões expirados foram processados
	c.JSON(http.StatusOK, gin.H{
		"message": "Expired auctions processed successfully.",
	})
}
