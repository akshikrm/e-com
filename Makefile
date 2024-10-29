build:
	@go build -o dist cmd/main.go 

run: build
	@clear
	@./dist

initdb: build
	@./dist --initdb 

seed-users: build
	@./dist --seed-users

nuke-db: build
	@./dist --nuke-db

