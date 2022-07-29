package models

type AddressCredentialRequest struct {
	Email           string `json:"email"`
	IDNumber        string `json:"id_number"`
	Surname         string `json:"surname"`
	FirstNames      string `json:"first_names"`
	StatementIssuer string `json:"statement_issuer"`
	StatementDate   string `json:"statement_date"`
	AddressLine1    string `json:"address_line_1"`
	AddressLine2    string `json:"address_line_2"`
	AddressLine3    string `json:"address_line_3"`
	City            string `json:"city"`
	PostalCode      string `json:"postal_code"`
	ExpiryDate      string `json:"expiry_date"`
}
