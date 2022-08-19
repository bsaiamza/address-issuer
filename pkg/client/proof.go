package client

import (
	"address_issuer/pkg/log"
	"address_issuer/pkg/models"
	"fmt"
)

func (c *Client) ListVerificationRecords() (models.ListVerificationRecordsResponse, error) {
	var verificationRecords models.ListVerificationRecordsResponse

	arg := models.AcapyGetRequest{
		Endpoint: "/present-proof/records",
		Response: &verificationRecords,
	}

	err := c.get(arg)
	if err != nil {
		log.Error.Printf("Failed on ACA-py /present-proof/records: %s", err.Error())
		return models.ListVerificationRecordsResponse{}, err
	}
	return verificationRecords, nil
}

func (c *Client) SendProofRequest(request models.ProofRequest) (models.SendProofRequestResponse, error) {
	var proof models.SendProofRequestResponse

	arg := models.AcapyPostRequest{
		Endpoint: "/present-proof/send-request",
		Body:     request,
		Response: &proof,
	}

	err := c.post(arg)
	if err != nil {
		log.Error.Printf("Failed on ACA-py /present-proof/send-request: %s", err.Error())
		return models.SendProofRequestResponse{}, err
	}
	return proof, nil
}

func (c *Client) GetPresExRecord(presExID string) (models.ProofRecords, error) {
	endpoint := fmt.Sprintf("/present-proof/records/%s", presExID)
	var presExRecord models.ProofRecords

	arg := models.AcapyGetRequest{
		Endpoint: endpoint,
		Response: &presExRecord,
	}

	err := c.get(arg)
	if err != nil {
		log.Error.Printf("Failed on ACA-py /present-proof/records/{pres_ex_id}: %s", err.Error())
		return models.ProofRecords{}, err
	}
	return presExRecord, nil
}
