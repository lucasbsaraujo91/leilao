package mongodb

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "MONGODB_URL"
	MONGODB_DB  = "MONGODB_DB"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	// Recuperando variáveis de ambiente
	mongoURL := os.Getenv(MONGODB_URL)
	mongoDatabase := os.Getenv(MONGODB_DB)

	// Log para garantir que a URL e o DB são carregados corretamente
	logger.Info(fmt.Sprintf("Connecting to MongoDB at %s, using database: %s", mongoURL, mongoDatabase))

	// Verifica se a URL do MongoDB está configurada
	if mongoURL == "" {
		err := fmt.Errorf("environment variable %s is not set", MONGODB_URL)
		logger.Error(err.Error(), err)
		return nil, err
	}

	// Verifica se o nome do banco de dados está configurado
	if mongoDatabase == "" {
		err := fmt.Errorf("environment variable %s is not set", MONGODB_DB)
		logger.Error(err.Error(), err)
		return nil, err
	}

	// Configura o cliente do MongoDB com timeout
	clientOptions := options.Client().ApplyURI(mongoURL).SetConnectTimeout(10 * time.Second)

	// Conectando com o MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Error trying to connect to MongoDB", err)
		return nil, err
	}

	// Verifica a conexão com o MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error trying to ping MongoDB", err)
		return nil, err
	}

	// Log de sucesso na conexão
	logger.Info("Successfully connected to MongoDB")

	// Retorna o banco de dados
	return client.Database(mongoDatabase), nil
}
