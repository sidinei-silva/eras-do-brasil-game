# Roda o servidor em modo de desenvolvimento
run-server:
	cd server/cmd/game && go run main.go

# Compila o binário final para a pasta bin/
build-server:
	go build -o bin/game ./server/cmd/game

# Roda os testes (quando você tiver)
test-server:
	go test ./server/...