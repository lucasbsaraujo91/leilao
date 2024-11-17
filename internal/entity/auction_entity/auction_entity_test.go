package auction_entity_test

import (
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAuction_ValidData(t *testing.T) {
	// Teste de criação de leilão com dados válidos
	auction, err := auction_entity.CreateAuction("Produto Teste", "Categoria Teste", "Descrição do Produto", auction_entity.New)

	// Verificando se o erro é nil, ou seja, criação bem-sucedida
	assert.Nil(t, err)
	// Verificando se o ID foi gerado
	assert.NotEmpty(t, auction.Id)
	// Verificando os valores dos campos
	assert.Equal(t, "Produto Teste", auction.ProductName)
	assert.Equal(t, "Categoria Teste", auction.Category)
	assert.Equal(t, "Descrição do Produto", auction.Description)
	assert.Equal(t, auction_entity.New, auction.Condition)
	assert.Equal(t, auction_entity.Active, auction.Status)
}

func TestCreateAuction_InvalidProductName(t *testing.T) {
	// Teste de criação de leilão com nome de produto inválido
	auction, err := auction_entity.CreateAuction("", "Categoria Teste", "Descrição do Produto", auction_entity.New)

	// Verificando se o erro é retornado devido ao nome do produto inválido
	assert.NotNil(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, internal_error.NewBadRequestError("invalid auction object"), err)
}

func TestCreateAuction_InvalidCategory(t *testing.T) {
	// Teste de criação de leilão com categoria inválida
	auction, err := auction_entity.CreateAuction("Produto Teste", "Ca", "Descrição do Produto", auction_entity.New)

	// Verificando se o erro é retornado devido à categoria inválida
	assert.NotNil(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, internal_error.NewBadRequestError("invalid auction object"), err)
}

func TestCreateAuction_InvalidDescription(t *testing.T) {
	// Teste de criação de leilão com descrição inválida
	auction, err := auction_entity.CreateAuction("Produto Teste", "Categoria Teste", "Desc", auction_entity.New)

	// Verificando se o erro é retornado devido à descrição inválida
	assert.NotNil(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, internal_error.NewBadRequestError("invalid auction object"), err)
}

func TestCreateAuction_InvalidCondition(t *testing.T) {
	// Teste de criação de leilão com condição inválida
	auction, err := auction_entity.CreateAuction("Produto Teste", "Categoria Teste", "Descrição do Produto", 100) // Condição inválida

	// Verificando se o erro é retornado devido à condição inválida
	assert.NotNil(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, internal_error.NewBadRequestError("invalid auction object"), err)
}

func TestValidateAuction_Valid(t *testing.T) {
	// Teste de validação com dados válidos
	auction := &auction_entity.Auction{
		ProductName: "Produto Teste",
		Category:    "Categoria Teste",
		Description: "Descrição do Produto",
		Condition:   auction_entity.New,
	}

	// Verificando se não há erro na validação
	err := auction.Validate()
	assert.Nil(t, err)
}

func TestValidateAuction_Invalid(t *testing.T) {
	// Teste de validação com dados inválidos
	auction := &auction_entity.Auction{
		ProductName: "",
		Category:    "Categoria Teste",
		Description: "Descrição do Produto",
		Condition:   auction_entity.New,
	}

	// Verificando se o erro é retornado devido ao nome do produto inválido
	err := auction.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, internal_error.NewBadRequestError("invalid auction object"), err)
}
