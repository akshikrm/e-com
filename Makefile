build:
	@go build -o dist cmd/main.go 

run: build
	@./dist

initdb: build
		@./dist --initdb 

