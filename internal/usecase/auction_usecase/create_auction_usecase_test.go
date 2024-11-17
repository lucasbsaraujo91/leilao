package auction_usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
)

type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	args := m.Called(ctx, auctionEntity)
	// Garantir que o erro retornado seja ou *internal_error.InternalError ou nil
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*internal_error.InternalError)
}

func (m *MockAuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, id)
	return args.Get(0).(*auction_entity.Auction), args.Get(1).(*internal_error.InternalError)
}

func (m *MockAuctionRepository) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, status, category, productName)
	return args.Get(0).([]auction_entity.Auction), args.Get(1).(*internal_error.InternalError)
}

func (m *MockAuctionRepository) FindExpiredAuctions(ctx context.Context, timestamp int64) ([]auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, timestamp)
	return args.Get(0).([]auction_entity.Auction), args.Get(1).(*internal_error.InternalError)
}

func (m *MockAuctionRepository) UpdateAuctionStatus(ctx context.Context, id string, status int) *internal_error.InternalError {
	args := m.Called(ctx, id, status)
	return args.Get(0).(*internal_error.InternalError)
}

func TestCreateAuction_Success(t *testing.T) {
	mockRepo := new(MockAuctionRepository)
	auctionUC := auction_usecase.NewAuctionUseCase(mockRepo, nil)

	// Define test input
	input := auction_usecase.AuctionInputDTO{
		ProductName: "Test Product",
		Category:    "Electronics",
		Description: "A very nice product.",
		Condition:   1,
	}

	// Mock behavior to simulate a successful creation (nil error)
	mockRepo.On("CreateAuction", mock.Anything, mock.Anything).Return(nil)

	// Call the method
	err := auctionUC.CreateAuction(context.Background(), input)

	// Assert that no error is returned
	assert.Nil(t, err)
}

func TestCreateAuction_Failure(t *testing.T) {
	mockRepo := new(MockAuctionRepository)
	auctionUC := auction_usecase.NewAuctionUseCase(mockRepo, nil)

	// Define test input
	input := auction_usecase.AuctionInputDTO{
		ProductName: "Test Product",
		Category:    "Electronics",
		Description: "A very nice product.",
		Condition:   1,
	}

	// Mock behavior to simulate a failure (non-nil error)
	mockRepo.On("CreateAuction", mock.Anything, mock.Anything).Return(&internal_error.InternalError{Message: "error creating auction"})

	// Call the method
	err := auctionUC.CreateAuction(context.Background(), input)

	// Assert that the error returned matches the expected one
	assert.NotNil(t, err)
	assert.Equal(t, "error creating auction", err.Message)
}
