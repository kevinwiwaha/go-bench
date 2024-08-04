package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Database configuration
const (
	DB_USER     = "savings_service"
	DB_PASSWORD = "oyCucmfCx7M978cGhKou"
	DB_NAME     = "savings_service_integration"
	DB_HOST     = "localhost"
	DB_PORT     = 5433
)

const (
	numRecords          = 5000000
	numGoroutines       = 50
	recordsPerGoroutine = numRecords / numGoroutines
	numCustomers        = 1000
)

func main() {
	rand.Seed(time.Now().UnixNano())

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	log.Println("Connecting to database with the following details:")
	log.Printf("Host: %s, Port: %d, User: %s, DBName: %s", DB_HOST, DB_PORT, DB_USER, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	var wg sync.WaitGroup

	customerIDs := generateCustomerIDs(numCustomers)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go seedData(db, &wg, customerIDs, i*recordsPerGoroutine, (i+1)*recordsPerGoroutine)
	}

	wg.Wait()

	fmt.Println("Seeding complete.")
}

func generateCustomerIDs(num int) []string {
	customerIDs := make(map[string]struct{})
	customerIDs["3252604"] = struct{}{}

	for len(customerIDs) < num {
		customerID := fmt.Sprintf("%d", rand.Intn(9000000)+1000000)
		customerIDs[customerID] = struct{}{}
	}

	keys := make([]string, 0, len(customerIDs))
	for k := range customerIDs {
		keys = append(keys, k)
	}

	return keys
}

func seedData(db *sql.DB, wg *sync.WaitGroup, customerIDs []string, start, end int) {
	defer wg.Done()

	stmt, err := db.Prepare(`INSERT INTO savings (savings_id, service_type, savings_amount, order_number, customer_id, plan_start_date, event_timestamp, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	serviceTypes := []string{"GO_FOOD", "GO_CAR", "GO_RIDE"}
	totalRecords := end - start

	for i := start; i < end; i++ {
		time.Sleep(10 * time.Millisecond)
		savingsID := uuid.New().String()
		serviceType := serviceTypes[rand.Intn(len(serviceTypes))]
		savingsAmount := rand.Intn(10)*1000 + 1000
		orderNumber := uuid.New().String()
		customerID := customerIDs[rand.Intn(len(customerIDs))]
		currentTime := time.Now()

		_, err = stmt.Exec(savingsID, serviceType, savingsAmount, orderNumber, customerID, currentTime, currentTime, currentTime)
		if err != nil {
			log.Printf("Failed to execute statement for record %d: %v", i, err)
			continue
		}

		progress := float64(i-start+1) / float64(totalRecords) * 100
		fmt.Printf("Goroutine %d: %.2f%% complete\n", start/recordsPerGoroutine, progress)
	}

	fmt.Printf("Goroutine %d finished inserting records\n", start/recordsPerGoroutine)
}
