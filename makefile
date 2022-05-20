
.PHONY : run
dep:
	go mod download
run:
	@echo "Start server..."
	go run src/main.go
build:
	go build -o bin/manga_spider src/main.go
	@echo "build success!"


