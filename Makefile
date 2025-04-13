# Fichier : Makefile

# Variables
DB_NAME=macrotracker
DB_USER=macro
DB_HOST=localhost
DB_PORT=5432

.PHONY: run test test-data reset-db sql shell help

## Lance l'application
run:
	go run main.go

## Lance tous les tests
test:
	go test ./... -v

test-data:
	go test ./data/test -v

## Réinitialise complètement la base (destructif)
reset-db:
	psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f ./data/config/schema.sql

## Ouvre un terminal SQL sur la base
shell:
	psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME)

## Affiche l'aide
help:
	@echo "Commandes disponibles :"
	@grep -E '^##' Makefile | sed -E 's/## //' | column -t -s ':'
