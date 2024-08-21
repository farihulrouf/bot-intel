package controllers

import (
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

// EventHandler handles incoming WhatsApp events
func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.PairSuccess:
		fmt.Println("Pair success:", v.ID.User)
		initialClient()
	case *events.Message:
		text := v.Message.GetConversation()
		fmt.Println("Message:", text)

		// Fetch business data based on the message text
		// Using default page and limit
		businesses, err := fetchBusinessData(text, defaultPage, defaultLimit)
		if err != nil {
			log.Printf("Failed to fetch business data: %v", err)
			return
		}

		// Print the business information
		if len(businesses) == 0 {
			fmt.Println("No business information found.")
		} else {
			for _, b := range businesses {
				fmt.Printf("Name: %s\nBusiness Name: %s\nAddress: %s\nSocial Media: %s\nWhatsApp: %s\nCategory: %s\nBusiness Photo: %s\nProduct Photo: %s\n\n",
					b.Name, b.BusinessName, b.FullAddress, b.SocialMediaUrl, b.WhatsappNumber, b.Category, b.BusinessPhoto, b.ProductPhoto)
			}
		}

	default:
		fmt.Printf("Unhandled event: %T\n", v)
	}
}
