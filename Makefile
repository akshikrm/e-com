build:
	@go build -o dist cmd/main.go 
		
run: build
	@clear
	@./dist

build-seeder:
	@go build -o seeder seed/main.go 

init-db: build-seeder
	@./seeder --init-db 

seed-users: build-seeder
	@./seeder --seed-users

seed-groups: build-seeder
	@./seeder --seed-groups

seed-roles: build-seeder
	@./seeder --seed-roles

seed-resources: build-seeder
	@./seeder --seed-resources

seed-permission: build-seeder
		@./seeder --seed-permission

nuke-db: build-seeder
	@./seeder --nuke-db

refresh-db: build-seeder 
	@./seeder --refresh-db
