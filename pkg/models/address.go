package models

type AddressCredentialRequest struct {
	AddressLine    string `json:"address_line"`
	City           string `json:"city"`
	Province       string `json:"province"`
	PostalCode     string `json:"postal_code"`
	CountryCode    string `json:"country_code"`
	StatementDate  string `json:"statement_date"`
	ExpiryDate     string `json:"expiry_date"`
	StatementImage string `json:"statement_image"`
	SelfAttested   string `json:"self_attested"`
	Email          string `json:"email"`
}
