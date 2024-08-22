package controllers

import (
	"bot_intel/helpers"
	"bot_intel/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mau.fi/whatsmeow/types/events"
)

const (
	apiUrl       = "https://data-grivy.vercel.app/api/business"
	defaultPage  = 1
	defaultLimit = 10
)

// Fetch business data from the API based on name, page, and limit
func fetchBusinessData(name string, page int, limit int) ([]models.Business, error) {
	// Use default values if page or limit are zero
	if page <= 0 {
		page = defaultPage
	}
	if limit <= 0 {
		limit = defaultLimit
	}

	// Construct the API URL with query parameters
	url := fmt.Sprintf("%s?page=%d&limit=%d&name=%s", apiUrl, page, limit, name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Businesses []models.Business `json:"businesses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Businesses, nil
}

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.PairSuccess:
		fmt.Println("Pair success:", v.ID.User)
		initialClient() // Ensure this initializes `whatsmeowClient`
	case *events.Message:
		text := v.Message.GetConversation()
		senderID := v.Info.Sender.String() // Mendapatkan ID pengirim
		fmt.Printf("Message received from: %s\n", senderID)

		// Fetch business data based on the message text
		businesses, err := fetchBusinessData(text, defaultPage, defaultLimit)
		if err != nil {
			log.Printf("Failed to fetch business data: %v", err)
			return
		}

		// Construct the response message
		var replyMessage string
		if len(businesses) == 0 {
			replyMessage = "No business information found."
		} else {
			replyMessage = "Here is the business information:\n"
			for _, b := range businesses {
				replyMessage += fmt.Sprintf("Name: %s\nBusiness Name: %s\nAddress: %s\nSocial Media: %s\nWhatsApp: %s\nCategory: %s\n\n",
					b.Name, b.BusinessName, b.FullAddress, b.SocialMediaUrl, b.WhatsappNumber, b.Category)
			}
			client, ok := clients["silver12"] // Assuming senderID is used as a key
			if !ok {
				log.Printf("No client found for sender ID: %s", senderID)
				return
			}

			// Send auto-reply
			err = helpers.SendMessageToPhoneNumber(client, senderID, replyMessage)
			if err != nil {
				log.Printf("Failed to send auto-reply: %v", err)
			}
		}

	default:
		fmt.Printf("Unhandled event: %T\n", v)
	}
}
