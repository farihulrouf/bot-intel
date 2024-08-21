package controllers

import (
	"bot_intel/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	clients        = make(map[string]*whatsmeow.Client)
	dataClient     = make(map[string]*whatsmeow.Client)
	mutex          = &sync.Mutex{}
	storeContainer *sqlstore.Container
	clientLog      waLog.Logger
)

func SetStoreContainer(container *sqlstore.Container) {
	storeContainer = container
}

func AddClient(id string, client *whatsmeow.Client) {
	mutex.Lock()
	defer mutex.Unlock()

	if client == nil {
		log.Printf("Failed to add client: client is nil for id %s\n", id)
		return
	}

	clients[id] = client
	log.Printf("Client added successfully: %s\n", id)
}

func GetClient(deviceStore *store.Device) *whatsmeow.Client {
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(EventHandler)
	return client
}

func initialClient() {
	mutex.Lock()
	defer mutex.Unlock()

	for key, value := range dataClient {
		clients[key] = value
	}
}

func setClientData(key string, client *whatsmeow.Client) {
	mutex.Lock()
	defer mutex.Unlock()

	// Clear existing data
	for k := range dataClient {
		delete(dataClient, k)
	}
	// Set new client
	dataClient[key] = client
}

func CreateDevice(w http.ResponseWriter, r *http.Request) {
	deviceStore := storeContainer.NewDevice()
	client := GetClient(deviceStore)
	deviceID := GenerateRandomString("Device", 3)
	setClientData(deviceID, client)
	qrCode, jid := connectClient(client)

	var response []models.ClientInfo

	fmt.Println("Client data after adding:", jid)

	// Iterate through the `clients` map to create response
	for key, client := range clients {
		response = append(response, models.ClientInfo{
			ID:     key,
			Number: client.Store.ID.String(),
			Busy:   true,
			QR:     "",
			Status: "connected",
			Name:   client.Store.PushName,
		})
	}

	// Add the new client to the response
	if qrCode != "" {
		response = append(response, models.ClientInfo{
			ID:     "",
			Number: "",
			Busy:   false,
			QR:     qrCode,
			Status: "pairing",
			Name:   "",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if len(response) > 0 {
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Failed to connect the client", http.StatusInternalServerError)
	}
}

func GenerateRandomString(prefix string, length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("%s-%s", prefix, string(b))
}

func connectClient(client *whatsmeow.Client) (string, *types.JID) {
	qrChan := make(chan string)

	// Disconnect client if it's already connected
	if client.IsConnected() {
		client.Disconnect()
	}

	// Generate new QR code for new login session
	qrChannel, _ := client.GetQRChannel(context.Background())
	go func() {
		for evt := range qrChannel {
			switch evt.Event {
			case "code":
				qrChan <- evt.Code
			case "login":
				close(qrChan)
			}
		}
	}()
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	qrCode := <-qrChan
	return qrCode, client.Store.ID
}
