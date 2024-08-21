package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bot_intel/controllers"
	"bot_intel/db"

	//"bot_intel/middleware"
	router "bot_intel/routers"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	// Memuat nilai dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Gagal memuat file .env: %v", err)
	}

	// Mengambil nilai dari variabel lingkungan
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Default port jika variabel lingkungan tidak ada
	}
	dbPath := os.Getenv("DB_PATH")

	// Mengatur logging untuk database
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Menginisialisasi SQL store
	storeContainer, err := sqlstore.New("sqlite3", dbPath, dbLog)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Menetapkan storeContainer ke variabel package controller
	controllers.SetStoreContainer(storeContainer)

	// Mengambil semua perangkat dari database
	devices, err := storeContainer.GetAllDevices()
	if err != nil {
		log.Fatalf("Gagal mengambil perangkat dari database: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		controllers.AddClient(controllers.GenerateRandomString("Device", 3), client)
		client.AddEventHandler(controllers.EventHandler)

		if client.Store.ID == nil {
			// Login baru
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				log.Fatalf("Gagal menghubungkan klien: %v", err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					// Menampilkan QR code
					fmt.Println("QR code:", evt.Code)
				} else {
					fmt.Println("Event login:", evt.Event)
				}
			}
		} else {
			// Sudah login, langsung hubungkan
			err = client.Connect()
			if err != nil {
				log.Fatalf("Gagal menghubungkan klien: %v", err)
			}
		}
	}

	// Inisialisasi database
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	defer db.CloseDB()

	// Mengatur router
	r := router.SetupRouter()

	// Menyajikan file statis jika diperlukan
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Mengaktifkan CORS jika diperlukan
	//corsHandler := middleware.SetupCORS(r)

	log.Printf("Server berjalan di port %s\n", port)

	// Jalankan server
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}

	// Shutdown yang aman
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		client.Disconnect()
	}
}
