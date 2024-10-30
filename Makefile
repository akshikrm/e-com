build:
	@go build -o dist cmd/main.go 

run: build
	@clear
	@./dist

initdb: build
	@./dist --initdb 

seed-users: build
	@./dist --seed-users

seed-roles: build
	@./dist --seed-roles

nuke-db: build
	@./dist --nuke-db

