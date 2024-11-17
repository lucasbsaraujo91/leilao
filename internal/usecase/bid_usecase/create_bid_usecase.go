package bid_usecase

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository              bid_entity.BidEntityRepository
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	maxBatchSize               int
	batchInsertInterval        time.Duration
	bidChannel                 chan bid_entity.Bid
	wg                         sync.WaitGroup
}

var bidBatch []bid_entity.Bid

type BidUseCaseInterface interface {
	CreateBid(
		ctx context.Context,
		bidInputDTO BidInputDTO) *internal_error.InternalError

	FindWinningBidByAuctionId(
		ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError)

	FindBidByAuctionId(
		ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
}

func NewBidUseCase(bidRepository bid_entity.BidEntityRepository, auctionRepositoryInterface auction_entity.AuctionRepositoryInterface) BidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		BidRepository:              bidRepository,
		maxBatchSize:               maxBatchSize,
		batchInsertInterval:        maxSizeInterval,
		bidChannel:                 make(chan bid_entity.Bid, 100), // Canal com buffer maior
		auctionRepositoryInterface: auctionRepositoryInterface,
	}

	// Inicia a goroutine para processar bids de forma contínua
	go bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

func (bu *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	ticker := time.NewTicker(bu.batchInsertInterval)
	defer ticker.Stop()

	for {
		select {
		case bidEntity := <-bu.bidChannel: // Consome lances do canal
			bidBatch = append(bidBatch, bidEntity)
			log.Printf("Bid added to batch. Current batch size: %d", len(bidBatch))

			// Processa o lote se atingir o tamanho máximo
			if len(bidBatch) >= bu.maxBatchSize {
				log.Println("Max batch size reached. Processing batch.")
				bu.processBids(ctx)
			}

		case <-ticker.C: // Processa a cada intervalo, mesmo que a batch esteja incompleta
			if len(bidBatch) > 0 {
				log.Println("Time interval reached. Processing batch.")
				bu.processBids(ctx)
			}
		}
	}
}

func (bu *BidUseCase) processBids(ctx context.Context) {
	if len(bidBatch) == 0 {
		return
	}

	log.Printf("Processing batch of %d bids...", len(bidBatch))

	// Processa os lances no repositório
	err := bu.BidRepository.CreateBid(ctx, bidBatch)
	if err != nil {
		logger.Error("Error processing batch:", err)
	} else {
		log.Printf("Successfully processed batch of %d bids", len(bidBatch))
	}

	// Limpa o batch após processamento
	bidBatch = nil
}

func (bu *BidUseCase) CreateBid(
	ctx context.Context,
	bidInputDTO BidInputDTO) *internal_error.InternalError {

	// Buscar o leilão
	auction, err := bu.auctionRepositoryInterface.FindAuctionById(ctx, bidInputDTO.AuctionId)
	if err != nil {
		log.Printf("Erro ao buscar leilão com ID %s: %v", bidInputDTO.AuctionId, err)
		return err
	}

	if auction == nil {
		log.Printf("Leilão %s não encontrado", bidInputDTO.AuctionId)
		return internal_error.NewInternalServerError("auction not found")
	}

	// Verificar o status do leilão
	if auction.Status == 1 { // Supondo que 1 seja "encerrado"
		log.Printf("Leilão %s encerrado, não é possível aceitar novos lances", bidInputDTO.AuctionId)
		return internal_error.NewInternalServerError("Leilão encerrado. Não é possível aceitar novos lances.")
	}

	// Cria a entidade do lance
	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.AuctionId, bidInputDTO.Amount)
	if err != nil {
		return err
	}

	// Adiciona ao canal para processamento
	select {
	case bu.bidChannel <- *bidEntity:
		log.Println("Bid successfully added to channel")
	default:
		// Se o canal estiver cheio, força o processamento do lote atual
		log.Println("Channel is full. Forcing batch processing.")
		bu.processBids(ctx)
		bu.bidChannel <- *bidEntity // Tenta novamente adicionar após esvaziar o batch
	}

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	if batchInsertInterval == "" {
		return 1 * time.Minute // Ajuste o intervalo conforme necessário
	}
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 1 * time.Minute // Caso o valor não esteja presente ou esteja incorreto, use um valor padrão
	}
	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 10 // Valor padrão para o tamanho do batch
	}
	return value
}
