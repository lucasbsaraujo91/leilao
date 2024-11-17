package auction_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

func CreateAuction(
	productName, category, description string,
	condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (au *Auction) Validate() *internal_error.InternalError {
	// Validação de campos obrigatórios e valores válidos
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 2 ||
		len(au.Description) <= 10 || // Descrição deve ter mais de 10 caracteres
		(au.Condition != New && au.Condition != Refurbished && au.Condition != Used) { // Condição deve ser New, Refurbished ou Used
		return internal_error.NewBadRequestError("invalid auction object")
	}

	return nil
}

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota + 1
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionEntity *Auction) *internal_error.InternalError

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]Auction, *internal_error.InternalError)

	FindAuctionById(
		ctx context.Context, id string) (*Auction, *internal_error.InternalError)

	FindExpiredAuctions(ctx context.Context, auctionTimeoutSeconds int64) ([]Auction, *internal_error.InternalError)

	UpdateAuctionStatus(ctx context.Context, id string, status int) *internal_error.InternalError
}
