package models

// Business represents the structure of the business data returned by the API
type Business struct {
	Name           string `json:"name"`
	BusinessName   string `json:"businessName"`
	FullAddress    string `json:"fullAddress"`
	SocialMediaUrl string `json:"socialMediaUrl"`
	WhatsappNumber string `json:"whatsappNumber"`
	Category       string `json:"category"`
	BusinessPhoto  string `json:"businessPhoto"`
	ProductPhoto   string `json:"productPhoto"`
}
