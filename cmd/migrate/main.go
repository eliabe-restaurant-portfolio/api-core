package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
	)

	migrationPath := "file://internal/connections/postgres/migrations"

	m, err := migrate.New(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Falha ao criar a instância do migrate: %v", err)
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		runUp(m)
	case "down":
		runDown(m)
	case "force":
		runForce(m)
	default:
		log.Fatalf("Comando desconhecido: %s. Use 'up', 'down', ou 'force'.", command)
	}
}

func runUp(m *migrate.Migrate) {
	log.Println("Aplicando migrações 'up'...")
	if err := m.Up(); err != nil {

		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Falha ao aplicar migrações 'up': %v", err)
		}
	}
	log.Println("Migrações 'up' aplicadas com sucesso!")
}

func runDown(m *migrate.Migrate) {
	log.Println("Revertendo uma migração 'down'...")
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Falha ao aplicar migração 'down': %v", err)
	}
	log.Println("Migração 'down' aplicada com sucesso!")
}

func runForce(m *migrate.Migrate) {
	if len(os.Args) < 3 {
		log.Fatal("Uso: go run main.go force <VERSION>")
	}
	version, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("A versão deve ser um número inteiro. Erro: %v", err)
	}

	log.Printf("Forçando o banco para a versão %d...", version)
	if err := m.Force(version); err != nil {
		log.Fatalf("Falha ao forçar a versão: %v", err)
	}
	log.Printf("Banco de dados forçado para a versão %d com sucesso!", version)
}
