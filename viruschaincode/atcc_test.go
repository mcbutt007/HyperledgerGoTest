package chaincode_test

import (
	"encoding/json"
	"chaincode"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.VirusChaincode{}
	err := assetTransfer.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestCreateSignature(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.VirusChaincode{}
	err := assetTransfer.UploadSignature(transactionContext, "", "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.UploadSignature(transactionContext, "", "virus1", "", "")
	require.EqualError(t, err, "the signature virus1 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve signature"))
	err = assetTransfer.UploadSignature(transactionContext, "", "virus1", "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve signature")
}

func TestReadAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.VirusSignature{SignatureID: "virus1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.VirusChaincode{}
	asset, err := assetTransfer.GetSignature(transactionContext, "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve signature"))
	_, err = assetTransfer.GetSignature(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve signature")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = assetTransfer.GetSignature(transactionContext, "virus1")
	require.EqualError(t, err, "Virus signature with ID virus1 not found")
	require.Nil(t, asset)
}

func TestUpdateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.VirusSignature{SignatureID: "virus1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.VirusChaincode{}
	err = assetTransfer.UpdateSignature(transactionContext, "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.UpdateSignature(transactionContext, "", "virus1", "")
	require.EqualError(t, err, "virus signature with ID virus1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve signature"))
	err = assetTransfer.UpdateSignature(transactionContext, "virus1", "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve signature")
}

func TestDeleteAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	asset := &chaincode.VirusSignature{SignatureID: "virus1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	assetTransfer := chaincode.VirusChaincode{}
	err = assetTransfer.DeleteSignature(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeleteSignature(transactionContext, "virus1")
	require.EqualError(t, err, "the asset virus1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.DeleteSignature(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestGetAllAssets(t *testing.T) {
	asset := &chaincode.VirusSignature{SignatureID: "virus1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	assetTransfer := &chaincode.VirusChaincode{}
	assets, err := assetTransfer.GetAllAssets(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*chaincode.Asset{asset}, assets)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	assets, err = assetTransfer.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, assets)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	assets, err = assetTransfer.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, assets)
}
