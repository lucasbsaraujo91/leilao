package rest_err_test

import (
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/internal_error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertError_BadRequest(t *testing.T) {
	// Simulando um erro "bad_request"
	internalErr := &internal_error.InternalError{
		Err:     "bad_request",
		Message: "Invalid input",
	}
	restErr := rest_err.ConvertError(internalErr)

	// Verificando se a conversão gerou um erro do tipo "bad_request"
	assert.NotNil(t, restErr)
	assert.Equal(t, "Invalid input", restErr.Message)
	assert.Equal(t, "bad_request", restErr.Err)
	assert.Equal(t, 400, restErr.Code) // HTTP Status BadRequest
}

func TestConvertError_NotFound(t *testing.T) {
	// Simulando um erro "not_found"
	internalErr := &internal_error.InternalError{
		Err:     "not_found",
		Message: "Resource not found",
	}
	restErr := rest_err.ConvertError(internalErr)

	// Verificando se a conversão gerou um erro do tipo "not_found"
	assert.NotNil(t, restErr)
	assert.Equal(t, "Resource not found", restErr.Message)
	assert.Equal(t, "not_found", restErr.Err)
	assert.Equal(t, 404, restErr.Code) // HTTP Status NotFound
}

func TestConvertError_InternalServerError(t *testing.T) {
	// Simulando um erro que não seja "bad_request" ou "not_found" (gerando um "internal_server")
	internalErr := &internal_error.InternalError{
		Err:     "unknown_error",
		Message: "Something went wrong",
	}
	restErr := rest_err.ConvertError(internalErr)

	// Verificando se a conversão gerou um erro do tipo "internal_server"
	assert.NotNil(t, restErr)
	assert.Equal(t, "Something went wrong", restErr.Message)
	assert.Equal(t, "internal_server", restErr.Err)
	assert.Equal(t, 500, restErr.Code) // HTTP Status InternalServerError
}

func TestNewBadRequestError(t *testing.T) {
	// Criando um erro de BadRequest manualmente
	restErr := rest_err.NewBadRequestError("Invalid input", rest_err.Causes{Field: "field1", Message: "This field is required"})

	// Verificando se o erro foi criado corretamente
	assert.NotNil(t, restErr)
	assert.Equal(t, "Invalid input", restErr.Message)
	assert.Equal(t, "bad_request", restErr.Err)
	assert.Equal(t, 400, restErr.Code) // HTTP Status BadRequest
	assert.Len(t, restErr.Causes, 1)
	assert.Equal(t, "field1", restErr.Causes[0].Field)
	assert.Equal(t, "This field is required", restErr.Causes[0].Message)
}

func TestNewInternalServerError(t *testing.T) {
	// Criando um erro de InternalServerError manualmente
	restErr := rest_err.NewInternalServerError("Something went wrong")

	// Verificando se o erro foi criado corretamente
	assert.NotNil(t, restErr)
	assert.Equal(t, "Something went wrong", restErr.Message)
	assert.Equal(t, "internal_server", restErr.Err)
	assert.Equal(t, 500, restErr.Code) // HTTP Status InternalServerError
	assert.Nil(t, restErr.Causes)
}

func TestNewNotFoundError(t *testing.T) {
	// Criando um erro de NotFoundError manualmente
	restErr := rest_err.NewNotFoundError("Resource not found")

	// Verificando se o erro foi criado corretamente
	assert.NotNil(t, restErr)
	assert.Equal(t, "Resource not found", restErr.Message)
	assert.Equal(t, "not_found", restErr.Err)
	assert.Equal(t, 404, restErr.Code) // HTTP Status NotFound
	assert.Nil(t, restErr.Causes)
}
