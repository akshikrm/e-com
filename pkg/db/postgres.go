package db

import (
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

const (
	CREATE_ROLE       = "CREATE TABLE IF NOT EXISTS roles (id SERIAL PRIMARY KEY, code varchar(10) NOT NULL, Name varchar(20) NOT NULL, Description varchar(120) NOT NULL, created_at timestamp DEFAULT NOW() NOT NULL, updated_at timestamp DEFAULT NOW() NOT NULL)"
	CREATE_RESOURCE   = "CREATE TABLE IF NOT EXISTS resources (id SERIAL PRIMARY KEY, code varchar(10) NOT NULL, Name varchar(20) NOT NULL, Description varchar(120) NOT NULL, created_at timestamp DEFAULT NOW() NOT NULL, updated_at timestamp DEFAULT NOW() NOT NULL)"
	CREATE_PERMISSION = "CREATE TABLE IF NOT EXISTS permissions (id SERIAL PRIMARY KEY, role_code INT NOT NULL, resource_code INT NOT NULL, r BOOLEAN DEFAULT false NOT NULL, w BOOLEAN DEFAULT false NOT NULL, u BOOLEAN DEFAULT false NOT NULL, d BOOLEAN DEFAULT false NOT NULL, created_at timestamp DEFAULT NOW() NOT NULL, updated_at timestamp DEFAULT NOW() NOT NULL)"
	CREATE_USERS      = "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, password varchar NOT NULL, created_at timestamp DEFAULT NOW() NOT NULL, updated_at timestamp DEFAULT NOW() NOT NULL)"
	CREATE_PROFILES   = "CREATE TABLE IF NOT EXISTS profiles (id SERIAL PRIMARY KEY, user_id int UNIQUE, first_name varchar(50) DEFAULT '' NOT NULL, last_name varchar(50) DEFAULT '' NOT NULL, email varchar(50) UNIQUE DEFAULT '' NOT NULL, pincode varchar(10) DEFAULT '' NOT NULL, address_one varchar(100) DEFAULT '' NOT NULL, address_two varchar(100) DEFAULT '' NOT NULL, phone_number varchar(15) DEFAULT '' NOT NULL, created_at timestamp DEFAULT NOW() NOT NULL, updated_at timestamp DEFAULT NOW() NOT NULL, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id))"
)

func (s *PostgresStore) Connect() error {
	db, err := sql.Open("postgres", s.getConnectionString())

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Println("🗃️ connected to database")
	s.DB = db

	initdb := flag.Bool("init-db", false, "initialize db if true")
	seedUsers := flag.Bool("seed-users", false, "seed db if true")
	seedRoles := flag.Bool("seed-roles", false, "seed db if true")
	seedResources := flag.Bool("seed-resources", false, "seed db if true")
	seedPermission := flag.Bool("seed-permission", false, "seed db if true")
	nukeDb := flag.Bool("nuke-db", false, "clear everything in the database")
	refreshDb := flag.Bool("refresh-db", false, "clear everything in the database")

	flag.Parse()
	if *initdb {
		s.Init()
		os.Exit(0)
	}

	if *seedRoles {
		s.seedRoles()
		os.Exit(0)
	}

	if *seedResources {
		s.seedResources()
		os.Exit(0)
	}

	if *seedUsers {
		s.seedUsers()
		os.Exit(0)
	}
	if *seedPermission {
		s.seedPermission()
		os.Exit(0)
	}

	if *refreshDb {
		s.NukeDB()
		s.Init()
		s.seedRoles()
		s.seedResources()
		s.seedPermission()
		s.seedUsers()
		os.Exit(0)
	}

	if *nukeDb {
		s.NukeDB()
		os.Exit(0)
	}

	return nil
}

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
		log.Printf("drop trigger %s\n", function)
	}
}

func (s *PostgresStore) NukeDB() {
	dropTrigger(s.DB, "update_user_task_updated_on", "roles")
	dropTrigger(s.DB, "update_user_task_updated_on", "users")
	dropTrigger(s.DB, "update_user_task_updated_on", "profiles")
	dropTables(s.DB, "roles")
	dropTables(s.DB, "resources")
	dropTables(s.DB, "profiles")
	dropTables(s.DB, "users")
	dropTables(s.DB, "permissions")
	dropFunction(s.DB, "update_updated_on_user_task")
}

func (s *PostgresStore) seedRoles() {
	log.Println("seeding roles")
	roleService := services.NewRoleService(s.DB)
	role := types.CreateRoleRequest{
		Name:        "Admin",
		Code:        "admin",
		Description: "Role assigned to admin",
	}
	err := roleService.Create(&role)
	if err != nil {
		log.Printf("Failed to seed role %s due to %s\n", role.Name, err)
	}
	log.Printf("Successfully seed role %s\n", role.Name)
}

func (s *PostgresStore) seedResources() {
	log.Println("seeding Resource")
	resourceService := services.NewResourceService(s.DB)
	resource := types.CreateResourceRequest{
		Name:        "Product",
		Code:        "product",
		Description: "resource assigned to admin",
	}
	err := resourceService.Create(&resource)
	if err != nil {
		log.Printf("Failed to seed resource %s due to %s\n", resource.Name, err)
	}
	log.Printf("Successfully seed resource %s\n", resource.Name)
}

func (s *PostgresStore) seedPermission() {
	log.Println("seeding permission")
	permissionService := services.NewPermissionService(s.DB)
	permission := types.CreateNewPermission{
		RoleCode:     1,
		ResourceCode: 1,
		R:            true,
	}
	err := permissionService.Create(&permission)
	if err != nil {
		log.Printf("Failed to seed permission due to %s\n", err)
	}
	log.Println("Successfully seed permission")
}

func (s *PostgresStore) seedUsers() {
	log.Println("seeding users")
	userService := services.NewUserService(s.DB)
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

func (s *PostgresStore) Init() {
	CreateTable(s.DB, CREATE_ROLE, "roles")
	CreateTable(s.DB, CREATE_RESOURCE, "resources")
	CreateTable(s.DB, CREATE_PERMISSION, "permissions")
	CreateTable(s.DB, CREATE_USERS, "users")
	CreateTable(s.DB, CREATE_PROFILES, "profiles")
	log.Println("successfully created all tables")

	CreateUpdatedAtFunction(s.DB)
	log.Println("successfully created all functions")

	CreateUpdatedAtTrigger(s.DB, "users")
	CreateUpdatedAtTrigger(s.DB, "profiles")
	CreateUpdatedAtTrigger(s.DB, "permissions")
	CreateUpdatedAtTrigger(s.DB, "roles")
	CreateUpdatedAtTrigger(s.DB, "resources")
	log.Println("successfully created all triggers")
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
func (s *PostgresStore) getConnectionString() string {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", db_host, db_port, db_user, db_name, db_password)

}
