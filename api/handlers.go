package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"address_issuer/pkg/client"
	"address_issuer/pkg/config"
	"address_issuer/pkg/log"
	"address_issuer/pkg/models"
	"address_issuer/pkg/server"
	"address_issuer/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

func health(config *config.Config) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(healthHandler(config), mdw...)
}
func healthHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func listConnections(config *config.Config, client *client.Client) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(listConnectionsHandler(config, client), mdw...)
}
func listConnectionsHandler(config *config.Config, client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
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

		log.Info.Println("Listing connections...")

		connections, err := client.ListConnections()
		if err != nil {
			log.Error.Printf("Failed to list connections: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to list connections: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		log.Info.Print("Connections listed successfully!")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(connections.Results)
	}
}

func listCredentials(config *config.Config, client *client.Client) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(listCredentialsHandler(config, client), mdw...)
}
func listCredentialsHandler(config *config.Config, client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
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

		log.Info.Println("Listing credentials records...")

		records, err := client.ListCredentialRecords()
		if err != nil {
			log.Error.Printf("Failed to list credential records: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "Failed to list credential records: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		log.Info.Print("Credential records listed successfully!")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records.Results)
	}
}

func getCredential(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(getCredentialHandler(config, client, cache), mdw...)
}
func getCredentialHandler(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
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

		// Step 1: Retrieve user information
		var userInfo models.AddressCredentialRequest
		err := json.NewDecoder(r.Body).Decode(&userInfo)
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

		// Step 2: Validate ID number
		log.Info.Println("Validating ID number...")
		validID := userInfo.IDNumber
		_, err = utils.IDValidator(validID)
		if err != nil {
			log.Error.Printf("ID validation failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "ID validation failed: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Info.Println("ID Validation passed")

		// Step 3: Create Invitation
		invitationRequest := models.CreateInvitationRequest{}

		invitation, err := client.CreateInvitation(invitationRequest)
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

		// Step 4: Cache user data for webhookEventsHandler
		err = cache.UpdateStruct(invitation.Invitation.RecipientKeys[0], userInfo)
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

		w.WriteHeader(http.StatusOK)
		res := server.Response{
			"success":    true,
			"credential": invitation.InvitationURL,
		}
		json.NewEncoder(w).Encode(res)
	}
}

func getCredentialByEmail(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(getCredentialByEmailHandler(config, client, cache), mdw...)
}
func getCredentialByEmailHandler(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
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

		// Step 1: Retrieve user information
		var userInfo models.AddressCredentialRequest
		err := json.NewDecoder(r.Body).Decode(&userInfo)
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
		err = utils.ValidEmail(userInfo.Email)
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

		// Step 3: Validate ID number
		log.Info.Println("Validating ID number...")
		validID := userInfo.IDNumber
		_, err = utils.IDValidator(validID)
		if err != nil {
			log.Error.Printf("ID validation failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := server.Response{
				"success": false,
				"msg":     "ID validation failed: " + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Info.Println("ID Validation passed")

		// Step 4: Create Invitation
		invitationRequest := models.CreateInvitationRequest{}

		invitation, err := client.CreateInvitation(invitationRequest)
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

		// Step 5: Cache user data
		err = cache.UpdateStruct(invitation.Invitation.RecipientKeys[0], userInfo)
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

		// Step 6: Generate a qr code for email
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

		// Step 7: Send email
		// prefix := string(userInfo.IDNumber[6])
		// if prefix >= "5" {
		// 	prefix = "Mr "
		// } else {
		// 	prefix = "Ms/Mrs "
		// }

		prefix := "Mr/Ms/Mrs "

		err = utils.SendCredentialByEmail(prefix+userInfo.Surname, userInfo.Email, invitation.Invitation.RecipientKeys[0], qrCodePng, config)
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

		// Step 8: Remove qr from os
		err = os.Remove("./" + invitation.Invitation.RecipientKeys[0] + ".png")
		if err != nil {
			log.Warning.Print("Failed to remove QR code: ", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func webhookEvents(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
	mdw := []server.Middleware{
		server.LogAPIRequest,
	}

	return server.ChainMiddleware(webhookEventsHandler(config, client, cache), mdw...)
}
func webhookEventsHandler(config *config.Config, client *client.Client, cache *utils.BigCache) http.HandlerFunc {
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

				_, err := client.PingConnection(request.ConnectionID, pingRequest)
				if err != nil {
					log.Error.Printf("Failed to ping holder: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			if request.State == "active" {
				userInfo, err := cache.ReadStruct(request.InvitationKey)
				if err != nil {
					log.Error.Printf("Failed to read cached user data: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// err = cache.UpdateString(userInfo.IDNumber, request.ConnectionID)
				// if err != nil {
				// 	log.Error.Printf("Failed to cache connID %s", err)
				// 	w.WriteHeader(http.StatusBadRequest)
				// 	return
				// }

				// credentialRequest := models.IssueCredentialRequest{
				// 	AutoRemove:      false,
				// 	ConnectionID:    request.ConnectionID,
				// 	Comment:         "Physical Address Credential",
				// 	CredDefID:       config.GetCredDefID(),
				// 	IssuerDid:       config.GetPublicDID(),
				// 	SchemaID:        config.GetSchemaID(),
				// 	SchemaIssuerDid: config.GetPublicDID(),
				// 	SchemaName:      config.GetSchemaName(),
				// 	SchemaVersion:   config.GetSchemaVersion(),
				// 	Trace:           false,
				// 	CredentialProposal: models.CredentialProposal{
				// 		Type: "issue-credential/1.0/credential-preview",
				// 		Attributes: []models.Attribute{
				// 			{
				// 				Name:  "ID Number",
				// 				Value: userInfo.IDNumber,
				// 			},
				// 			{
				// 				Name:  "First Names",
				// 				Value: userInfo.FirstNames,
				// 			},
				// 			{
				// 				Name:  "Surname",
				// 				Value: userInfo.Surname,
				// 			},
				// 			{
				// 				Name:  "Statement Issuer",
				// 				Value: userInfo.StatementIssuer,
				// 			},
				// 			{
				// 				Name:  "Statement Date",
				// 				Value: userInfo.StatementDate,
				// 			},
				// 			{
				// 				Name:  "Address Line 1",
				// 				Value: userInfo.AddressLine1,
				// 			},
				// 			{
				// 				Name:  "Address Line 2",
				// 				Value: userInfo.AddressLine2,
				// 			},
				// 			{
				// 				Name:  "Address Line 3",
				// 				Value: userInfo.AddressLine3,
				// 			},
				// 			{
				// 				Name:  "City",
				// 				Value: userInfo.City,
				// 			},
				// 			{
				// 				Name:  "Postal Code",
				// 				Value: userInfo.PostalCode,
				// 			},
				// 			{
				// 				Name:  "Expiry Date",
				// 				Value: userInfo.ExpiryDate,
				// 			},
				// 		},
				// 	},
				// }

				// _, err = client.IssueCredential(credentialRequest)
				// if err != nil {
				// 	log.Error.Printf("Failed to send credential offer: %s", err)
				// 	w.WriteHeader(http.StatusBadRequest)
				// 	return
				// }

				// // For email credential notification
				// err = cache.UpdateStruct(request.ConnectionID, userInfo)
				// if err != nil {
				// 	log.Error.Printf("Failed to cache user data: %s", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }

				// cache.DeleteStruct(request.InvitationKey)

				// log.Info.Println("Credential offer sent")
				// w.WriteHeader(http.StatusOK)

				credDefID := config.GetCornerstoneCredDefID()

				idProof := map[string]interface{}{
					"name": "ID Number",
					"restrictions": []map[string]interface{}{
						{
							"cred_def_id": credDefID,
						},
					},
				}

				namesProof := map[string]interface{}{
					"name": "First Names",
					"restrictions": []map[string]interface{}{
						{
							"cred_def_id": credDefID,
						},
					},
				}

				surnameProof := map[string]interface{}{
					"name": "Surname",
					"restrictions": []map[string]interface{}{
						{
							"cred_def_id": credDefID,
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
							idProof,
							namesProof,
							surnameProof,
						},
						RequestedPredicates: models.RequestedPredicates{},
					},
				}

				_, err = client.SendProofRequest(identityProofRequest)
				if err != nil {
					log.Error.Printf("Failed to send proof request: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				err = cache.UpdateStruct(request.ConnectionID, userInfo)
				if err != nil {
					log.Error.Printf("Failed to cache user data: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cache.DeleteStruct(request.InvitationKey)

				log.Info.Println("Proof request sent")
				w.WriteHeader(http.StatusOK)
			}

		case "issue_credential":
			var request models.IssueCredentialWebhookResponse
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				log.Error.Printf("Failed to decode request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userInfo, err := cache.ReadStruct(request.ConnectionID)
			if err != nil {
				log.Error.Printf("Failed to read cached user data: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if request.State == "credential_issued" && userInfo.Email != "" {
				log.Info.Println("Sending credential issued notification...")

				// prefix := string(userInfo.IDNumber[6])
				// if prefix >= "5" {
				// 	prefix = "Mr "
				// } else {
				// 	prefix = "Ms/Mrs "
				// }

				prefix := "Mr/Ms/Mrs"

				err = utils.SendNotificationEmail(prefix+userInfo.Surname, userInfo.Email, config)
				if err != nil {
					log.Error.Printf("Failed to send credential notification email: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cache.DeleteStruct(request.ConnectionID)

				log.Info.Println("Notified user successfully about issued credential!")
				w.WriteHeader(http.StatusOK)
			}

		case "present_proof":
			var request models.PresentProofWebhookResponse
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				log.Error.Printf("Failed to decode request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userInfo, err := cache.ReadStruct(request.ConnectionID)
			if err != nil {
				log.Error.Printf("Failed to read cached user data: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if request.State == "verified" {
				presExRecord, err := client.GetPresExRecord(request.PresentationExchangeID)
				if err != nil {
					log.Error.Printf("Failed to get presentation record: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if presExRecord.Verified == "true" {
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
								Name:  "Statement Issuer",
								Value: userInfo.StatementIssuer,
							},
							{
								Name:  "Statement Date",
								Value: userInfo.StatementDate,
							},
							{
								Name:  "Address Line 1",
								Value: userInfo.AddressLine1,
							},
							{
								Name:  "Address Line 2",
								Value: userInfo.AddressLine2,
							},
							{
								Name:  "Address Line 3",
								Value: userInfo.AddressLine3,
							},
							{
								Name:  "City",
								Value: userInfo.City,
							},
							{
								Name:  "Postal Code",
								Value: userInfo.PostalCode,
							},
							{
								Name:  "Expiry Date",
								Value: userInfo.ExpiryDate,
							},
						},
					},
				}

				fmt.Println("step 5")
				_, err = client.IssueCredential(credentialRequest)
				if err != nil {
					log.Error.Printf("Failed to send credential offer: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// For email credential notification
				if userInfo.Email == "" {
					cache.DeleteStruct(request.ConnectionID)
				}

				log.Info.Println("Credential offer sent")
				w.WriteHeader(http.StatusOK)
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
