package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/utils"
)

func (au *AuctionUseCase) FindExpiredAuctions(
	ctx context.Context) ([]AuctionOutputDTO, *internal_error.InternalError) {

	// Buscando leilões expirados
	auctions, err := au.auctionRepositoryInterface.FindExpiredAuctions(ctx, utils.GetAuctionTimeoutSeconds())
	if err != nil {
		return nil, err
	}

	// Criando a lista de leilões expirados no formato adequado para a saída
	var expiredAuctions []AuctionOutputDTO
	for _, auction := range auctions {
		auctionDTO := AuctionOutputDTO{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   ProductCondition(auction.Condition),
			Status:      AuctionStatus(auction.Status),
			Timestamp:   auction.Timestamp,
		}
		expiredAuctions = append(expiredAuctions, auctionDTO)
	}

	//fmt.Printf("Leilões expirados: %v\n", expiredAuctions)

	// Retorna os leilões expirados
	return expiredAuctions, nil
}
