package main

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	store := db.NewStorage()
	db.Connect(store)
	Seed(store)
}

func Seed(s *db.Storage) error {
	initdb := flag.Bool("init-db", false, "initialize db if true")
	seedUsers := flag.Bool("seed-users", false, "seed db if true")
	seedResources := flag.Bool("seed-resources", false, "seed db if true")
	seedPermission := flag.Bool("seed-permission", false, "seed db if true")
	nukeDb := flag.Bool("nuke-db", false, "clear everything in the database")
	refreshDb := flag.Bool("refresh-db", false, "clear everything in the database")

	flag.Parse()
	if *initdb {
		Init(s)
		os.Exit(0)
	}

	if *seedResources {
		seedResourcesFunc(s)
		os.Exit(0)
	}

	if *seedUsers {
		seedUsersFunc(s)
		os.Exit(0)
	}
	if *seedPermission {
		seedPermissionFunc(s)
		os.Exit(0)
	}

	if *refreshDb {
		NukeDB(s)
		Init(s)
		seedResourcesFunc(s)
		seedPermissionFunc(s)
		seedUsersFunc(s)
		os.Exit(0)
	}

	if *nukeDb {
		NukeDB(s)
		os.Exit(0)
	}

	return nil
}

const (
	CREATE_ROLE             = "CREATE TABLE IF NOT EXISTS roles (id SERIAL PRIMARY KEY, name VARCHAR(20) NOT NULL, code VARCHAR(10) UNIQUE NOT NULL, description VARCHAR(120) NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL)"
	CREATE_RESOURCE         = "CREATE TABLE IF NOT EXISTS resources (id SERIAL PRIMARY KEY, name VARCHAR(10) NOT NULL, code VARCHAR(10) UNIQUE NOT NULL, description VARCHAR(120) NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL)"
	CREATE_PERMISSION       = "CREATE TABLE IF NOT EXISTS permissions (id SERIAL PRIMARY KEY, role_code VARCHAR(10) NOT NULL, resource_code VARCHAR(10) NOT NULL, r BOOLEAN DEFAULT false NOT NULL, w BOOLEAN DEFAULT false NOT NULL, u BOOLEAN DEFAULT false NOT NULL, d BOOLEAN DEFAULT false NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL,deleted_at TIMESTAMP DEFAULT NULL, CONSTRAINT fk_role FOREIGN KEY(role_code) REFERENCES roles(code), CONSTRAINT fk_resource FOREIGN KEY(role_code) REFERENCES resources(code))"
	CREATE_USERS            = "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, password VARCHAR NOT NULL, role_code VARCHAR(10) DEFAULT user NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL, CONSTRAINT fk_role FOREIGN KEY(role_code) REFERENCES roles(code))"
	CREATE_PROFILES         = "CREATE TABLE IF NOT EXISTS profiles (id SERIAL PRIMARY KEY, user_id int UNIQUE, first_name VARCHAR(50) DEFAULT '' NOT NULL, last_name VARCHAR(50) DEFAULT '' NOT NULL, email VARCHAR(50) UNIQUE DEFAULT '' NOT NULL, pincode VARCHAR(10) DEFAULT '' NOT NULL, address_one VARCHAR(100) DEFAULT '' NOT NULL, address_two VARCHAR(100) DEFAULT '' NOT NULL, phone_number VARCHAR(15) DEFAULT '' NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id))"
	CREATE_PRODUCT          = "CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY, name VARCHAR(30), slug VARCHAR(30), price INTEGER NOT NULL DEFAULT 0, image VARCHAR(100),  description VARCHAR(300) NOT NULL, category_id INTEGER NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL, CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES product_categories(id))"
	CREATE_PRODUCT_CATEGORY = "create table if not exists product_categories(id SERIAL PRIMARY KEY, name VARCHAR(30) NOT NULL, slug VARCHAR(30) NOT NULL,enabled BOOLEAN DEFAULT true, description VARCHAR(120) NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL)"
	CREATE_CART             = "CREATE TABLE IF NOT EXISTS carts (id SERIAL PRIMARY KEY, user_id INTEGER NOT NULL, product_id INTEGER NOT NULL, quantity INTEGER DEFAULT 1 NOT NULL, created_at TIMESTAMP DEFAULT NOW() NOT NULL, updated_at TIMESTAMP DEFAULT NOW() NOT NULL, deleted_at TIMESTAMP DEFAULT NULL, CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id), CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id))"
)

func dropTables(store *sql.DB, table string) {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	if _, err := store.Exec(query); err != nil {
		log.Printf("Failed to drop %s due to %s\n", table, err)
	} else {
		log.Printf("drop %s\n", table)
	}
}

func dropTrigger(store *sql.DB, trigger string, table string) {
	query := fmt.Sprintf("DROP TRIGGER IF EXISTS %s on %s", trigger, table)
	if _, err := store.Exec(query); err != nil {
		log.Printf("Failed to drop trigger %s  on table %s, due to %s\n", trigger, table, err)
	} else {
		log.Printf("drop trigger %s on table %s\n", trigger, table)
	}
}

func dropFunction(store *sql.DB, function string) {
	query := fmt.Sprintf("DROP FUNCTION IF EXISTS %s\n", function)
	if _, err := store.Exec(query); err != nil {
		log.Printf("Failed to drop function %s due to %s", function, err)
	} else {
		log.Printf("drop function %s\n", function)
	}
}

func NukeDB(s *db.Storage) {
	dropTrigger(s.DB, "update_user_task_updated_on", "roles")
	dropTrigger(s.DB, "update_user_task_updated_on", "users")
	dropTrigger(s.DB, "update_user_task_updated_on", "profiles")
	dropTrigger(s.DB, "update_user_task_updated_on", "resources")
	dropTrigger(s.DB, "update_user_task_updated_on", "permissions")
	dropTrigger(s.DB, "update_user_task_updated_on", "products")
	dropTrigger(s.DB, "update_user_task_updated_on", "carts")
	dropTrigger(s.DB, "update_user_task_updated_on", "product_categories")

	dropTables(s.DB, "permissions")
	dropTables(s.DB, "profiles")
	dropTables(s.DB, "carts")
	dropTables(s.DB, "users")
	dropTables(s.DB, "roles")
	dropTables(s.DB, "products")
	dropTables(s.DB, "product_categories")
	dropTables(s.DB, "resources")
	dropFunction(s.DB, "update_updated_on_user_task")
}

func seedRolesFunc(s *db.Storage, role *types.CreateRoleRequest) {
	log.Println("seeding roles")
	roleService := services.NewRoleService(s)
	err := roleService.Create(role)
	if err != nil {
		log.Printf("Failed to seed role %s due to %s\n", role.Name, err)
	}
	log.Printf("Successfully seed role %s\n", role.Name)
}

func seedResourcesFunc(s *db.Storage) {
	log.Println("seeding Resource")
	resourceService := services.NewResourceService(s)
	resource := types.CreateResourceRequest{
		Name:        "User",
		Code:        "user",
		Description: "resource assigned to admin",
	}
	err := resourceService.Create(&resource)
	if err != nil {
		log.Printf("Failed to seed resource %s due to %s\n", resource.Name, err)
	}
	log.Printf("Successfully seed resource %s\n", resource.Name)
}

func seedPermissionFunc(s *db.Storage) {
	log.Println("seeding permission")
	permissionService := services.NewPermissionService(s)
	permission := types.CreateNewPermission{
		RoleCode:     "admin",
		ResourceCode: "user",
		R:            true,
		U:            true,
		D:            true,
	}
	err := permissionService.Create(&permission)
	if err != nil {
		log.Printf("Failed to seed permission due to %s\n", err)
	}
	log.Println("Successfully seed permission")
}

func seedAdminFunc(s *db.Storage) {
	log.Println("seeding admin")
	userModel := storage.NewUserStorage(s.DB)
	profileModel := storage.NewProfileStorage(s.DB)
	userService := services.NewUserService(userModel, profileModel)
	user := types.CreateUserRequest{
		FirstName: "Admin",
		LastName:  "Me",
		Email:     "admin@me.com",
		Password:  "root",
		Role:      "admin",
	}
	_, err := userService.Create(user)
	if err != nil {
		log.Printf("Failed to seed admin due to %s\n", err)
	}
	log.Println("Successfully seed admin")
}

func seedUsersFunc(s *db.Storage) {
	log.Println("seeding users")
	userModel := storage.NewUserStorage(s.DB)
	profileModel := storage.NewProfileStorage(s.DB)
	userService := services.NewUserService(userModel, profileModel)
	userFile, err := os.Open("./seed/users.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer userFile.Close()

	byteValue, err := io.ReadAll(userFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	users := []types.CreateUserRequest{}
	json.Unmarshal(byteValue, &users)
	for i, element := range users {
		if _, err := userService.Create(element); err != nil {
			log.Printf("Failed to add user %s due to %s\n", element.Email, err)
			continue
		}
		log.Printf("Inserting %d\n", i)
	}
	log.Println("Successfully seed users")
}

func Init(s *db.Storage) {
	CreateTable(s.DB, CREATE_ROLE, "roles")
	CreateTable(s.DB, CREATE_RESOURCE, "resources")
	CreateTable(s.DB, CREATE_PERMISSION, "permissions")
	CreateTable(s.DB, CREATE_USERS, "users")
	CreateTable(s.DB, CREATE_PROFILES, "profiles")
	CreateTable(s.DB, CREATE_PRODUCT_CATEGORY, "product_categories")
	CreateTable(s.DB, CREATE_PRODUCT, "products")
	CreateTable(s.DB, CREATE_CART, "carts")
	log.Println("successfully created all tables")

	CreateUpdatedAtFunction(s.DB)
	log.Println("successfully created all functions")

	CreateUpdatedAtTrigger(s.DB, "users")
	CreateUpdatedAtTrigger(s.DB, "profiles")
	CreateUpdatedAtTrigger(s.DB, "permissions")
	CreateUpdatedAtTrigger(s.DB, "roles")
	CreateUpdatedAtTrigger(s.DB, "resources")
	CreateUpdatedAtTrigger(s.DB, "products")
	CreateUpdatedAtTrigger(s.DB, "carts")
	CreateUpdatedAtTrigger(s.DB, "product_categories")
	log.Println("successfully created all triggers")

	adminRole := types.CreateRoleRequest{
		Name:        "Admin",
		Code:        "admin",
		Description: "Role assigned to admin",
	}
	seedRolesFunc(s, &adminRole)

	userRole := types.CreateRoleRequest{
		Name:        "User",
		Code:        "user",
		Description: "Role assigned to user",
	}
	seedRolesFunc(s, &userRole)
	seedAdminFunc(s)

}

func CreateTable(store *sql.DB, query string, table string) {
	log.Printf("Creating table %s\n", table)
	_, err := store.Exec(query)
	if err != nil {
		log.Printf("Failed to create %s table due to %s\n", table, err)
		os.Exit(1)
	}
	log.Printf("Created %s table\n", table)

}

func CreateUpdatedAtFunction(db *sql.DB) {
	log.Println("Creating updated at function")
	query := `CREATE  FUNCTION update_updated_on_user_task() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ language 'plpgsql';`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create function update_updated_on_user_task %s\n", err)
		os.Exit(1)
	}
	log.Printf("Created function update_updated_on_user_task\n")
}

func CreateUpdatedAtTrigger(db *sql.DB, table string) {
	log.Printf("Creating trigger update_user_task_updated_on on %s\n", table)
	query := fmt.Sprintf(`CREATE TRIGGER update_user_task_updated_on BEFORE UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE update_updated_on_user_task();`, table)

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create trigger update_user_task_updated_on on %s due to %s\n", table, err)
		os.Exit(1)
	}
	log.Printf("Created trigger update_user_task_updated_on on %s\n", table)
}
