package services

import (
    "database/sql"
    "os"
    "sync"
	"log"

    _ "github.com/go-sql-driver/mysql"
)

type DatabaseService struct {
    mu          sync.RWMutex
    db          *sql.DB
    schemaQuery string
}

func NewDatabaseService() *DatabaseService {
    s := &DatabaseService{}
    _ = s.ConnectDB()
    return s
}

func (s *DatabaseService) ConnectDB() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.db != nil {
        _ = s.db.Close()
    }

    // Ambil DSN dari Environment Variable (os.Getenv)
    dsn := os.Getenv("MYSQL_DSN")
    if dsn == "" {
        // Fallback default jika ENV belum diset
        dsn = "Admin:Adminjkt168#$@tcp(10.60.192.101:3306)/go_shipping_db?parseTime=true"
    }

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    if err := db.Ping(); err != nil {
        return err
    }
    s.db = db

	log.Println("Koneksi database MySQL berhasil terhubung!")

    // Otomatis jalankan sql/schema.sql saat aplikasi terkoneksi ke MySQL
    if schemaBytes, err := os.ReadFile("sql/schema.sql"); err == nil {
        _, _ = s.db.Exec(string(schemaBytes))
    }
    return nil
}