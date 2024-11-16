package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auction by id = %s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (repo *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return auctionsEntity, nil
}

func (ar *AuctionRepository) FindExpiredAuctions(ctx context.Context, auctionTimeoutSeconds int64) ([]auction_entity.Auction, *internal_error.InternalError) {
	// Tempo atual
	currentTime := time.Now().Unix()

	// Filtro para buscar leilões com status aberto (0) e timestamp expirado
	filter := bson.M{
		"status": 0,
		"timestamp": bson.M{
			"$lt": currentTime - auctionTimeoutSeconds, // Leilão expirado
		},
	}

	// Buscando leilões expirados no banco de dados
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		return nil, internal_error.NewInternalServerError("Erro ao buscar leilões expirados")
	}
	defer cursor.Close(ctx)

	var auctions []auction_entity.Auction
	for cursor.Next(ctx) {
		var auctionEntityMongo AuctionEntityMongo
		if err := cursor.Decode(&auctionEntityMongo); err != nil {
			return nil, internal_error.NewInternalServerError("Erro ao decodificar leilão")
		}

		// Convertendo a entidade MongoDB para a entidade Auction
		auction := auction_entity.Auction{
			Id:          auctionEntityMongo.Id,
			ProductName: auctionEntityMongo.ProductName,
			Category:    auctionEntityMongo.Category,
			Description: auctionEntityMongo.Description,
			Condition:   auctionEntityMongo.Condition,
			Status:      auctionEntityMongo.Status,
			Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
		}

		auctions = append(auctions, auction)
	}

	if err := cursor.Err(); err != nil {
		return nil, internal_error.NewInternalServerError("Erro ao iterar sobre os leilões")
	}

	return auctions, nil
}

func (ar *AuctionRepository) UpdateAuctionStatus(ctx context.Context, id string, status int) *internal_error.InternalError {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Erro ao atualizar status do leilão com ID %s: %v\n", id, err)
		return internal_error.NewInternalServerError("Erro ao atualizar status do leilão")
	}

	return nil
}
