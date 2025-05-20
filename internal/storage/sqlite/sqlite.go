package sqlite

import (
	"database/sql"

	"github.com/Shreya20002/students-go/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// Sqlite represents the SQLite storage with a database connection.
type Sqlite struct {
	Db *sql.DB
}

// New is like a constructor in other languages - naming convention in go
// for pointers: if * on type of struct or primitive type, it means its an sddress pointer (refercence). If * on var name, it is deferencing - providing val at that address
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

// to make the application plugable, we will implement interface as sqlite , so that in future we can change sqlite to
// any other database like mysql, postgres, etc. without changing the code
// to make this method inside out Sqlite struct , we use func(s *Sqlite) , the things inside brscket
// pointer to struct includes this method inside our Sqlite struct
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	// press ctrl first + then click on Exec to see the definition of the function
	// as well as extra function and returns it can perform , like a mini contextual doc
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}
