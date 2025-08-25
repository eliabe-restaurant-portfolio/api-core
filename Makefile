include setup/dev/Makefile

generate-keys:
	openssl req -x509 -nodes -newkey rsa:2048 -keyout storage/private_key.pem -out storage/public_key.pem -sha256 -days 365 -subj "/CN=localhost"

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-force:
	go run cmd/migrate/main.go force