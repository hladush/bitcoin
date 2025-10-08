package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/ihladush/bitcoin/internal/clients"
	"github.com/ihladush/bitcoin/internal/handlers"
	"github.com/ihladush/bitcoin/internal/repository"
	"github.com/ihladush/bitcoin/internal/services"
)

func main() {
	// Initialize database
	repo, err := repository.NewSQLiteRepository("bitcoin_tracker.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	// Initialize Bitcoin client
	client := clients.NewBlockchairClient()

	// Initialize service
	service := services.NewBitcoinService(repo, client)

	// Initialize handlers
	handler := handlers.NewBitcoinHandler(service)

	// Setup routes
	router := setupRoutes(handler)

	// Start background sync worker
	go startBackgroundSync(service)

	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("🚀 Bitcoin Tracker API starting on port 8080")
		log.Println("📋 API Documentation:")
		log.Println("   GET    /health                        - Health check")
		log.Println("   GET    /addresses                     - List all tracked addresses")
		log.Println("   POST   /addresses                     - Add new address")
		log.Println("   GET    /addresses/{address}           - Get address details")
		log.Println("   DELETE /addresses/{address}           - Remove address")
		log.Println("   GET    /addresses/{address}/balance   - Get address balance")
		log.Println("   GET    /addresses/{address}/transactions - Get address transactions")
		log.Println("   POST   /addresses/{address}/sync      - Sync specific address")
		log.Println("   POST   /sync                          - Sync all addresses")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")
}

// setupRoutes configures all API routes
func setupRoutes(handler *handlers.BitcoinHandler) *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Address management
	router.HandleFunc("/addresses", handler.GetAllAddresses).Methods("GET")
	router.HandleFunc("/addresses", handler.AddAddress).Methods("POST")
	router.HandleFunc("/addresses/{address}", handler.GetAddress).Methods("GET")
	router.HandleFunc("/addresses/{address}", handler.RemoveAddress).Methods("DELETE")

	// Balance and transactions
	router.HandleFunc("/addresses/{address}/balance", handler.GetBalance).Methods("GET")
	router.HandleFunc("/addresses/{address}/transactions", handler.GetTransactions).Methods("GET")

	// Synchronization
	router.HandleFunc("/addresses/{address}/sync", handler.SyncAddress).Methods("POST")
	router.HandleFunc("/sync", handler.SyncAllAddresses).Methods("POST")

	// Add CORS middleware
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	return router
}

// startBackgroundSync runs periodic synchronization
func startBackgroundSync(service *services.BitcoinService) {
	ticker := time.NewTicker(5 * time.Minute) // Sync every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🔄 Starting background sync...")
		if err := service.SyncAllAddresses(); err != nil {
			log.Printf("❌ Background sync failed: %v", err)
		} else {
			log.Println("✅ Background sync completed")
		}
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
