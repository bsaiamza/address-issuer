package models

type AddressCredentialRequest struct {
	Email           string `json:"email"`
	IDNumber        string `json:"id_number"`
	Surname         string `json:"surname"`
	Forenames       string `json:"forenames"`
	StatementIssuer string `json:"statement_issuer"`
	StatementDate   string `json:"statement_date"`
	Street1         string `json:"street_1"`
	Street2         string `json:"street_2"`
	Street3         string `json:"street_3"`
	City            string `json:"city"`
	PostalCode      string `json:"postal_code"`
	ExpiryDate      string `json:"expiry_date"`
}
