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

func (s *PostgresStore) Connect() error {
	db, err := sql.Open("postgres", s.getConnectionString())

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Print("🗃️ connected to database")
	s.DB = db

	initdb := flag.Bool("initdb", false, "initialze db if true")
	seedUsers := flag.Bool("seed-users", false, "seed db if true")
	seedRoles := flag.Bool("seed-roles", false, "seed db if true")
	nukeDb := flag.Bool("nuke-db", false, "clear everything in the database")

	flag.Parse()
	if *initdb {
		s.Init()
		os.Exit(0)
	}

	if *seedRoles {
		s.seedRoles()
		os.Exit(0)
	}

	if *seedUsers {
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
		fmt.Printf("Failed to drop %s due to %s\n", table, err)
	} else {
		fmt.Printf("drop %s\n", table)
	}
}

func dropTrigger(store *sql.DB, trigger string, table string) {
	query := fmt.Sprintf("DROP TRIGGER IF EXISTS %s on %s", trigger, table)
	if _, err := store.Exec(query); err != nil {
		fmt.Printf("Failed to drop trigger %s  on table %s, due to %s\n", trigger, table, err)
	} else {
		fmt.Printf("drop trigger %s on table %s\n", trigger, table)
	}
}

func dropFunction(store *sql.DB, function string) {
	query := fmt.Sprintf("DROP FUNCTION IF EXISTS %s\n", function)
	if _, err := store.Exec(query); err != nil {
		fmt.Printf("Failed to drop function %s due to %s", function, err)
	} else {
		fmt.Printf("drop trigger %s\n", function)
	}
}

func (s *PostgresStore) NukeDB() {
	dropTrigger(s.DB, "update_user_task_updated_on", "roles")
	dropTrigger(s.DB, "update_user_task_updated_on", "users")
	dropTrigger(s.DB, "update_user_task_updated_on", "profiles")
	dropTables(s.DB, "roles")
	dropTables(s.DB, "profiles")
	dropTables(s.DB, "users")
	dropFunction(s.DB, "update_updated_on_user_task")
}

func (s *PostgresStore) seedRoles() {
	fmt.Println("seeding roles")
	roleService := services.NewRoleService(s.DB)
	role := types.CreateRoleRequest{
		Name:        "Admin",
		Code:        "admin",
		Description: "Role assigned to admin",
	}
	err := roleService.Create(&role)
	if err != nil {
		log.Printf("Failed to create role %s due to %s", role.Name, err)
	}
	log.Printf("Successfully created role %s", role.Name)
}

func (s *PostgresStore) seedUsers() {
	fmt.Println("seeding users")
	userService := services.NewUserService(s.DB)
	userFile, err := os.Open("./seed/users.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer userFile.Close()

	byteValue, err := io.ReadAll(userFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	users := []types.CreateUserRequest{}
	json.Unmarshal(byteValue, &users)
	for i, element := range users {
		if _, err := userService.Create(element); err != nil {
			fmt.Printf("Failed to add user %s due to %s", element.Email, err)
			continue
		}
		fmt.Printf("Inserting %d", i)
	}
}

func (s *PostgresStore) Init() {
	CreateRoleTable(s.DB)
	CreateUserTable(s.DB)
	CreateProfileTable(s.DB)
	log.Println("successfully created all tables")

	CreateUpdatedAtFunction(s.DB)
	log.Println("successfully created all functions")

	CreateUpdatedAtTrigger(s.DB, "users")
	CreateUpdatedAtTrigger(s.DB, "profiles")
	CreateUpdatedAtTrigger(s.DB, "roles")
	log.Println("successfully created all triggers")
}

func CreateRoleTable(db *sql.DB) {
	log.Println("Creating roles table")
	query := `CREATE TABLE IF NOT EXISTS roles (
	id serial primary key,
	code varchar(10) NOT NULL,
	Name varchar(20) NOT NULL,
	Description varchar(120) NOT NULL,
	created_at timestamp DEFAULT NOW() NOT NULL,
	updated_at timestamp DEFAULT NOW() NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create roles table %s", err)
		os.Exit(1)
	}
	log.Println("Created roles table")

}

func CreateUserTable(db *sql.DB) {
	log.Println("Creating users table")
	query := `CREATE TABLE IF NOT EXISTS users (
	id serial primary key,
	password varchar NOT NULL,
	created_at timestamp DEFAULT NOW() NOT NULL,
	updated_at timestamp DEFAULT NOW() NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create users table %s", err)
		os.Exit(1)
	}
	log.Println("Created users table")
}

func CreateProfileTable(db *sql.DB) {
	log.Println("Creating profiles table")
	query := `CREATE TABLE IF NOT EXISTS profiles (
	id serial primary key,
	user_id int UNIQUE,
	first_name varchar(50) DEFAULT '' NOT NULL,
	last_name varchar(50) DEFAULT '' NOT NULL,
	email varchar(50) UNIQUE DEFAULT '' NOT NULL,
	pincode varchar(10) DEFAULT '' NOT NULL,
	address_one varchar(100) DEFAULT '' NOT NULL,
	address_two varchar(100) DEFAULT '' NOT NULL,
	phone_number varchar(15) DEFAULT '' NOT NULL,
	created_at timestamp DEFAULT NOW() NOT NULL,
	updated_at timestamp DEFAULT NOW() NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create profiles table %s", err)
		os.Exit(1)
	}
	log.Println("Created profiles table")
}

func CreateUpdatedAtFunction(db *sql.DB) {
	log.Println("Creating updated at function")
	query := `CREATE  FUNCTION update_updated_on_user_task() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ language 'plpgsql';`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create function update_updated_on_user_task %s", err)
		os.Exit(1)
	}
	log.Println("Created function update_updated_on_user_task")
}

func CreateUpdatedAtTrigger(db *sql.DB, table string) {
	log.Printf("Creating trigger update_user_task_updated_on on %s", table)
	query := fmt.Sprintf(`CREATE TRIGGER update_user_task_updated_on BEFORE UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE update_updated_on_user_task();`, table)

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create trigger update_user_task_updated_on on %s due to %s", table, err)
		os.Exit(1)
	}
	log.Printf("Created trigger update_user_task_updated_on on %s", table)
}
func (s *PostgresStore) getConnectionString() string {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", db_host, db_port, db_user, db_name, db_password)

}
