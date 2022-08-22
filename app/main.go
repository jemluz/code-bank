package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/codeedu/codebank/infrastructure/grpc/server"
	"github.com/codeedu/codebank/infrastructure/kafka"
	"github.com/codeedu/codebank/infrastructure/repository"
	"github.com/codeedu/codebank/usecase"
)

func main() {
	db := setupDB()
	defer db.Close()

	// cc := domain.NewCreditCard()
	// cc.Number = "1234"
	// cc.Name = "Jemima"
	// cc.ExpirationMonth = 7
	// cc.ExpirationYear = 2021
	// cc.CVV = 123
	// cc.Limit = 1000
	// cc.Balance = 0

	// repo := repository.NewTransactionRepositoryDB(db)
	// repo.CreateCreditCard(*cc)

	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGRPC(processTransactionUseCase)
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDB(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer("host.docker.internal:9094")
	return producer
}

func setupDB() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database", err)
	}

	return db
}

func serveGRPC(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	fmt.Println("Rodando grpc server...")
	grpcServer.Serve()
}
