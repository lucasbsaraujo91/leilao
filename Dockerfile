FROM golang:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/auction cmd/auction/main.go

# Instala o cron
RUN apt-get update && apt-get install -y cron

# Adiciona o cron job para rodar a função de fechamento de leilão
RUN echo "* * * * * root /bin/sh -c '/app/auction close_expired'" >> /etc/crontab

EXPOSE 8080

# Inicia o cron e a aplicação
ENTRYPOINT cron && /app/auction