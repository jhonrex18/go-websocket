package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/denisenkom/go-mssqldb"
    "github.com/joho/godotenv"
)

var db *sql.DB

func init() {
    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    initDB()
}

func initDB() {
    connString := os.Getenv("MSSQL_CONN_STRING")
    var err error
    db, err = sql.Open("sqlserver", connString)
    if err != nil {
        log.Fatal("Error creating connection pool: ", err.Error())
    }

    // Check if the connection is alive
    err = db.Ping()
    if err != nil {
        log.Fatal("Error connecting to database: ", err)
    }

    // Set database connection pooling limits
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)

    fmt.Println("Connected to MSSQL!")
}

func saveDataToDB(cameraName, licensePlate, visitTimestamp, vehicleType, confidenceLevel, photoURL string) error {
    query := `INSERT INTO Requests (camera_name, license_plate, visit_timestamp, vehicle_type, confidence_level, photo) 
              VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`
    stmt, err := db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing statement: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(cameraName, licensePlate, visitTimestamp, vehicleType, confidenceLevel, photoURL)
    if err != nil {
        return fmt.Errorf("error executing statement: %v", err)
    }
    return nil
}

func catchRequest(w http.ResponseWriter, r *http.Request) {
    // Set security headers
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")

    // Limit request size to 1MB
    r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    cameraName := r.FormValue("camera_name")
    licensePlate := r.FormValue("license_plate")
    visitTimestamp := r.FormValue("visit_timestamp")
    vehicleType := r.FormValue("vehicle_type")
    confidenceLevel := r.FormValue("confidence_level")
    photoURL := r.FormValue("photo")

    // Save data to the database
    err = saveDataToDB(cameraName, licensePlate, visitTimestamp, vehicleType, confidenceLevel, photoURL)
    if err != nil {
        log.Printf("Error saving data: %v", err)
        http.Error(w, "Error saving data", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Data saved!")
}

func main() {
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      http.DefaultServeMux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  15 * time.Second,
    }

    http.HandleFunc("/", catchRequest)
    log.Fatal(srv.ListenAndServe())
}
