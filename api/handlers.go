package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"address-issuer/pkg/acapy"
	"address-issuer/pkg/config"
	"address-issuer/pkg/log"
	"address-issuer/pkg/models"
	"address-issuer/pkg/server"
	"address-issuer/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

func getCredential(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(getCredentialHandler(config, acapy, cache), mdw...)
}
func getCredentialHandler(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			log.Warning.Print("Incorrect request method!")
			w.WriteHeader(http.StatusMethodNotAllowed)
			res := server.Response{
				"success": false,
				"msg":     "Warning: Incorrect request method!",
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		defer r.Body.Close()

		log.Info.Println("Creating credential request...")

		// Step 1: Retrieve address information
		var addrInfo models.AddressCredentialRequest
		err := json.NewDecoder(r.Body).Decode(&addrInfo)
		if err != nil {
			log.Error.Printf("Failed to decode credential data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to decode credential data: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 2: Create Invitation
		invitationRequest := models.CreateInvitationRequest{}

		invitation, err := acapy.CreateInvitation(invitationRequest)
		if err != nil {
			log.Error.Printf("Failed to create invitation: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to create invitation: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 3: Cache address data for webhookEventsHandler
		err = cache.Address(invitation.Invitation.RecipientKeys[0], addrInfo)
		if err != nil {
			log.Error.Printf("Failed to cache address data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to cache address data: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		res := server.Response{
			"success":    true,
			"credential": invitation.InvitationURL,
		}
		json.NewEncoder(w).Encode(res)
	}
}

func getCredentialByEmail(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(getCredentialByEmailHandler(config, acapy, cache), mdw...)
}
func getCredentialByEmailHandler(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			log.Warning.Print("Incorrect request method!")
			w.WriteHeader(http.StatusMethodNotAllowed)
			res := server.Response{
				"success": false,
				"msg":     "Warning: Incorrect request method!",
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		defer r.Body.Close()

		log.Info.Println("Creating credential request...")

		// Step 1: Retrieve address information
		var addrInfo models.AddressCredentialRequest
		err := json.NewDecoder(r.Body).Decode(&addrInfo)
		if err != nil {
			log.Error.Printf("Failed to decode credential data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to decode credential data: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 2: Validate email address
		err = utils.ValidEmail(addrInfo.Email)
		if err != nil {
			log.Error.Printf("Failed %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 3: Create Invitation
		invitationRequest := models.CreateInvitationRequest{}

		invitation, err := acapy.CreateInvitation(invitationRequest)
		if err != nil {
			log.Error.Printf("Failed to create invitation: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to create invitation: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 4: Cache address data for webhookEventsHandler
		err = cache.Address(invitation.Invitation.RecipientKeys[0], addrInfo)
		if err != nil {
			log.Error.Printf("Failed to cache user data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to cache user data: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 5: Generate a qr code for email
		qrCodePng, err := qrcode.Encode(invitation.InvitationURL, qrcode.Medium, 256)
		if err != nil {
			log.Warning.Print("Failed to create QR code: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to create QR code: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 6: Send email
		err = utils.SendCredentialByEmail(addrInfo.Email, invitation.Invitation.RecipientKeys[0], qrCodePng, config)
		if err != nil {
			log.Warning.Print("Failed to send credential by email: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to send credential by email",
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		// Step 7: Remove qr from os once email is sent
		err = os.Remove("./" + invitation.Invitation.RecipientKeys[0] + ".png")
		if err != nil {
			log.Warning.Print("Failed to remove QR code: ", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func webhookEvents(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(webhookEventsHandler(config, acapy, cache), mdw...)
}
func webhookEventsHandler(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			log.Warning.Print("Incorrect request method!")
			w.WriteHeader(http.StatusMethodNotAllowed)
			res := server.Response{
				"success": false,
				"msg":     "Warning: Incorrect request method!",
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		defer r.Body.Close()

		topic := mux.Vars(r)["topic"]

		switch topic {
		case "connections":
			var request models.Connection
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				log.Error.Printf("Failed to decode request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if request.State == "response" {
				pingRequest := models.PingConnectionRequest{
					Comment: "Ping",
				}

				_, err := acapy.PingConnection(request.ConnectionID, pingRequest)
				if err != nil {
					log.Error.Printf("Failed to ping holder: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			if request.State == "active" {
				addrInfo, err := cache.ReadAddress(request.InvitationKey)
				if err != nil {
					log.Error.Printf("Failed to read cached user data: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				schemaID := config.GetCornerstoneSchemaID()

				namesProof := map[string]interface{}{
					"name": "names",
					"restrictions": []map[string]interface{}{
						{
							"schema_id": schemaID,
						},
					},
				}

				surnameProof := map[string]interface{}{
					"name": "surname",
					"restrictions": []map[string]interface{}{
						{
							"schema_id": schemaID,
						},
					},
				}

				identityProofRequest := models.ProofRequest{
					Comment:      "Cornerstone Proof Request",
					ConnectionID: request.ConnectionID,
					PresentationRequest: models.PresentationRequest{
						Name:    "Proof of Identity",
						Version: "1.0",
						RequestedAttributes: models.RequestedAttributes{
							namesProof,
							surnameProof,
						},
						RequestedPredicates: models.RequestedPredicates{},
					},
				}

				_, err = acapy.SendProofRequest(identityProofRequest)
				if err != nil {
					log.Error.Printf("Failed to send proof request: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				err = cache.Address(request.ConnectionID, addrInfo)
				if err != nil {
					log.Error.Printf("Failed to cache user data: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cache.DeleteAddress(request.InvitationKey)

				log.Info.Println("Proof request sent")
			}

		case "issue_credential":
			var request models.IssueCredentialWebhookResponse
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				log.Error.Printf("Failed to decode request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userInfo, err := cache.ReadAddress(request.ConnectionID)
			if err != nil {
				log.Error.Printf("Failed to read cached user data: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if request.State == "credential_issued" && userInfo.Email != "" {
				log.Info.Println("Sending credential issued notification...")

				err = utils.SendNotificationEmail(userInfo.Email, config)
				if err != nil {
					log.Error.Printf("Failed to send credential notification email: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cache.DeleteAddress(request.ConnectionID)

				log.Info.Println("Notified user successfully about issued credential!")
			}

		case "present_proof":
			var request models.PresentProofWebhookResponse
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				log.Error.Printf("Failed to decode request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			addrInfo, err := cache.ReadAddress(request.ConnectionID)
			if err != nil {
				log.Error.Printf("Failed to read cached user data: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if request.State == "verified" {
				presExRecord, err := acapy.GetPresExRecord(request.PresentationExchangeID)
				if err != nil {
					log.Error.Printf("Failed to get presentation record: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if presExRecord.Verified == "true" {
					log.Info.Println("Cornerstone Credential Verification: Successful!")

					sdJSON, _ := time.Parse("2006-01-02T15:04:05.999Z", addrInfo.StatementDate)
					sd := sdJSON.Format("20060102")

					now := time.Now()
					threeMonths := time.Hour * 24 * 90
					expiryPeriod := now.Add(threeMonths)
					ed := expiryPeriod.Format("20060102")

					si64 := utils.ImageBase64()
					si := ""

					if addrInfo.StatementImage == "" {
						si = si64
					}

					credentialRequest := models.IssueCredentialRequest{
						AutoRemove:      false,
						ConnectionID:    request.ConnectionID,
						Comment:         "Physical Address Credential",
						CredDefID:       config.GetCredDefID(),
						IssuerDid:       config.GetPublicDID(),
						SchemaID:        config.GetSchemaID(),
						SchemaIssuerDid: config.GetPublicDID(),
						SchemaName:      config.GetSchemaName(),
						SchemaVersion:   config.GetSchemaVersion(),
						Trace:           false,
						CredentialProposal: models.CredentialProposal{
							Type: "issue-credential/1.0/credential-preview",
							Attributes: []models.Attribute{
								{
									Name:  "address_line",
									Value: addrInfo.AddressLine,
								},
								{
									Name:  "city",
									Value: addrInfo.City,
								},
								{
									Name:  "province",
									Value: addrInfo.Province,
								},
								{
									Name:  "postal_code",
									Value: addrInfo.PostalCode,
								},
								{
									Name:  "country_code",
									Value: addrInfo.CountryCode,
								},
								{
									Name:  "statement_date",
									Value: sd,
								},
								{
									Name:  "expiry_date",
									Value: ed,
								},
								{
									Name:  "statement_image",
									Value: si,
								},
								{
									Name:  "self_attested",
									Value: addrInfo.SelfAttested,
								},
							},
						},
					}

					_, err = acapy.IssueCredential(credentialRequest)
					if err != nil {
						log.Error.Printf("Failed to send credential offer: %s", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					if addrInfo.Email == "" {
						cache.DeleteAddress(request.ConnectionID)
					}

					log.Info.Println("Credential offer sent")
				} else {
					log.Info.Println("Cornerstone Credential Verification: Unsuccessful!")
				}
			}

		case "basicmessages":
		case "revocation_registry":
		case "problem_report":
		case "issuer_cred_rev":

		default:
			log.Warning.Printf("Unexpected topic: %s", topic)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
