// internal/utils/auction_config.go
package utils

import (
	"log"
	"os"
	"strconv"
)

// GetAuctionTimeoutSeconds lê o tempo de expiração do leilão a partir das variáveis de ambiente
// e retorna o valor em segundos. Em caso de erro, retorna o valor padrão de 3600 segundos.
func GetAuctionTimeoutSeconds() int64 {
	auctionTimeout := os.Getenv("AUCTION_TIMEOUT_SECONDS")

	// Verificar se a variável de ambiente foi definida
	if auctionTimeout == "" {
		log.Println("Variável AUCTION_TIMEOUT_SECONDS não definida. Usando valor padrão de 3600 segundos (1 hora).")
		return 3600 // valor padrão de 1 hora
	}

	timeoutSeconds, errConv := strconv.ParseInt(auctionTimeout, 10, 64)
	if errConv != nil {
		log.Printf("Erro ao interpretar AUCTION_TIMEOUT_SECONDS: %v. Usando valor padrão de 3600 segundos (1 hora).", errConv)
		return 3600 // valor padrão de 1 hora
	}

	return timeoutSeconds
}
