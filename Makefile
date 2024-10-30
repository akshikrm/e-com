build:
	@go build -o dist cmd/main.go 

run: build
	@clear
	@./dist

init-db: build
	@./dist --init-db 

seed-users: build
	@./dist --seed-users

seed-roles: build
	@./dist --seed-roles

nuke-db: build
	@./dist --nuke-db

