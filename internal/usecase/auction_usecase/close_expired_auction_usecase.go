package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"log"
	"time"
)

func (au *AuctionUseCase) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	go func() {
		for {
			// Obtém os leilões expirados
			expiredAuctions, err := au.FindExpiredAuctions(ctx)
			if err != nil {
				log.Printf("Erro ao buscar leilões expirados: %v\n", err)
				return
			}

			// Fecha os leilões expirados
			for _, auction := range expiredAuctions {
				err := au.auctionRepositoryInterface.UpdateAuctionStatus(ctx, auction.Id, 1) // 1 = Fechado
				if err != nil {
					log.Printf("Erro ao fechar leilão com ID %s: %v\n", auction.Id, err)
				} else {
					log.Printf("Leilão com ID %s fechado com sucesso.\n", auction.Id)
				}
			}

			// Espera o intervalo antes de verificar novamente
			time.Sleep(time.Minute) // intervalo de 1 minuto
		}
	}()
	return nil
}
