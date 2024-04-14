package chaincode

import (
  "fmt"
  "time"
  "encoding/json"
  "log"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type VirusChaincode struct {
  contractapi.Contract
}
// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type VirusSignature struct {
    IPFSHash      string  `json:"IPFSHash"`
    SignatureID   string  `json:"SignatureID"`
    Timestamp     int64   `json:"Timestamp"`
    Uploader      string  `json:"Uploader"`
    VirusName     string  `json:"VirusName"`
}

func (t *VirusChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
    virusSignatures := []VirusSignature{
        {SignatureID: "1", VirusName: "SampleVirus1", IPFSHash: "QmZQHmuXvF1AifghrGnNH4uey5iF1hzeRZvfevF2kg19nV", Uploader: "Org1", Timestamp: time.Now().Unix()},
        {SignatureID: "2", VirusName: "SampleVirus2", IPFSHash: "QmZQHmuXvF1AifghrGnNH4uey5iF1hzeRZvfevF2kg19nW", Uploader: "Org2", Timestamp: time.Now().Unix()},
        // Add more sample virus signatures as needed
    }

    for _, signature := range virusSignatures {
        virusJSON, err := json.Marshal(signature)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(signature.SignatureID, virusJSON)
        if err != nil {
			      return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}
func (t *VirusChaincode) UploadSignature(ctx contractapi.TransactionContextInterface,ipfsHash string, signatureID string, uploader string, virusName string) error {
    signature := VirusSignature{
        IPFSHash:    ipfsHash,
        SignatureID: signatureID,
        Timestamp:   time.Now().Unix(),
        Uploader:    uploader,
        VirusName:   virusName,
    }

    virusJSON, err := json.Marshal(signature)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(signatureID, virusJSON)
}

func (t *VirusChaincode) GetSignature(ctx contractapi.TransactionContextInterface, signatureID string) (*VirusSignature, error) {
    virusJSON, err := ctx.GetStub().GetState(signatureID)
    if err != nil {
        return nil, err
    }
    if virusJSON == nil {
        return nil, fmt.Errorf("Virus signature with ID %s not found", signatureID)
    }

    var signature VirusSignature
    err = json.Unmarshal(virusJSON, &signature)
    if err != nil {
        return nil, err
    }

    return &signature, nil
}
func (t *VirusChaincode) UpdateSignature(ctx contractapi.TransactionContextInterface, signatureID string, newVirusName string, newIPFSHash string) error {
    // Retrieve the existing virus signature from the ledger
    virusJSON, err := ctx.GetStub().GetState(signatureID)
    if err != nil {
        return fmt.Errorf("failed to read virus signature from ledger: %v", err)
    }
    if virusJSON == nil {
        return fmt.Errorf("virus signature with ID %s does not exist", signatureID)
    }

    var existingSignature VirusSignature
    err = json.Unmarshal(virusJSON, &existingSignature)
    if err != nil {
        return err
    }

    // Update the fields of the existing virus signature
    if newVirusName != "" {
        existingSignature.VirusName = newVirusName
    }
    if newIPFSHash != "" {
        existingSignature.IPFSHash = newIPFSHash
    }

    // Marshal the updated virus signature back to JSON
    updatedVirusJSON, err := json.Marshal(existingSignature)
    if err != nil {
        return err
    }

    // Write the updated virus signature back to the ledger
    err = ctx.GetStub().PutState(signatureID, updatedVirusJSON)
    if err != nil {
        return fmt.Errorf("failed to update virus signature on ledger: %v", err)
    }

    return nil
}
func (t *VirusChaincode) DeleteSignature(ctx contractapi.TransactionContextInterface, signatureID string) error {
    // Check if the virus signature exists
    virusJSON, err := ctx.GetStub().GetState(signatureID)
    if err != nil {
        return fmt.Errorf("failed to read virus signature from ledger: %v", err)
    }
    if virusJSON == nil {
        return fmt.Errorf("virus signature with ID %s does not exist", signatureID)
    }

    // Delete the virus signature from the ledger
    err = ctx.GetStub().DelState(signatureID)
    if err != nil {
        return fmt.Errorf("failed to delete virus signature from ledger: %v", err)
    }

    return nil
}

func (t *VirusChaincode) SignatureExists(ctx contractapi.TransactionContextInterface, signatureID string) (bool, error) {
    // Check if the virus signature exists
    virusJSON, err := ctx.GetStub().GetState(signatureID)
    if err != nil {
        return false, fmt.Errorf("failed to read virus signature from ledger: %v", err)
    }
    if virusJSON == nil {
        return false, nil
    }
    return true, nil
}

func (t *VirusChaincode) GetAllSignatures(ctx contractapi.TransactionContextInterface) ([]*VirusSignature, error) {
    // Retrieve all virus signatures from the ledger
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, fmt.Errorf("failed to read virus signatures from ledger: %v", err)
    }
    defer resultsIterator.Close()

    var virusSignatures []*VirusSignature
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var signature VirusSignature
        err = json.Unmarshal(queryResponse.Value, &signature)
        if err != nil {
            return nil, err
        }
        virusSignatures = append(virusSignatures, &signature)
    }

    return virusSignatures, nil
}

func main() {
    virusChaincode := new(VirusChaincode)
    contractAPI, err := contractapi.NewChaincode(virusChaincode)
    if err != nil {
        log.Fatal("Error creating virus chaincode: ", err)
    }

    if err := contractAPI.Start(); err != nil {
        log.Fatal("Error starting virus chaincode: ", err)
    }
}
