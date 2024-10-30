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

seed-resources: build
	@./dist --seed-resources

seed-permission: build
		@./dist --seed-permission


nuke-db: build
	@./dist --nuke-db

refresh-db: build 
	@./dist --refresh-db
