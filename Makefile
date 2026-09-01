run:
	go run ./cmd/mathua

build:
	go build -o bin/mathua ./cmd/mathua

build-all:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -o dist/mathua-windows-amd64.exe ./cmd/mathua
	GOOS=linux GOARCH=amd64 go build -o dist/mathua-linux-amd64 ./cmd/mathua
	GOOS=darwin GOARCH=amd64 go build -o dist/mathua-darwin-amd64 ./cmd/mathua
	GOOS=darwin GOARCH=arm64 go build -o dist/mathua-darwin-arm64 ./cmd/mathua

test:
	go test ./...

validate:
	go run scripts/validate_graph.go
	python3 scripts/audit_lessons.py

fuzz:
	go test ./internal/generator/... -run TestFuzz -count 1000

tidy:
	go mod tidy

serve:
	go run ./cmd/mathua --serve --port 8080
