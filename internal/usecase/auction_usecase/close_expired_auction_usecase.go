package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"log"
)

func (au *AuctionUseCase) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	// Cria uma goroutine para processar o fechamento dos leilões expirados
	go func() {
		// Busca os leilões expirados usando o método existente
		expiredAuctions, err := au.FindExpiredAuctions(ctx)
		if err != nil {
			log.Printf("Erro ao buscar leilões expirados: %v", err)
			return
		}

		// Se não houver leilões expirados, encerra o processamento
		if len(expiredAuctions) == 0 {
			log.Println("Nenhum leilão expirado encontrado.")
			return
		}

		// Itera pelos leilões expirados e atualiza o status para "fechado"
		for _, auction := range expiredAuctions {
			err := au.auctionRepositoryInterface.UpdateAuctionStatus(ctx, auction.Id, 1) // 1 = Fechado
			if err != nil {
				log.Printf("Erro ao fechar leilão com ID %s: %v\n", auction.Id, err)
				// Continua tentando fechar os outros leilões mesmo em caso de erro
			} else {
				log.Printf("Leilão com ID %s fechado com sucesso.\n", auction.Id)
			}
		}
	}()

	// Retorna nil para indicar que a goroutine foi iniciada com sucesso
	return nil
}
