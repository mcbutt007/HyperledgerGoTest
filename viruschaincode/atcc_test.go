package main

import (
    "testing"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestVirusChaincode(t *testing.T) {
    // Create a new instance of the chaincode
    chaincode := new(VirusChaincode)

    // Create a new ChaincodeMockStub
    stub := shimtest.NewMockStub("mockstub", chaincode)

    // Test InitLedger function
    response := stub.MockInit("1", nil)
    if response.Status != shim.OK {
        t.Errorf("InitLedger failed: %s", response.Message)
    }

    // Test UploadSignature function
    response = stub.MockInvoke("1", [][]byte{[]byte("UploadSignature"), []byte("1"), []byte("Virus1"), []byte("QmHash")})
    if response.Status != shim.OK {
        t.Errorf("UploadSignature failed: %s", response.Message)
    }

    // Test GetSignature function
    response = stub.MockInvoke("1", [][]byte{[]byte("GetSignature"), []byte("1")})
    if response.Status != shim.OK {
        t.Errorf("GetSignature failed: %s", response.Message)
    }

    // Test UpdateSignature function
    response = stub.MockInvoke("1", [][]byte{[]byte("UpdateSignature"), []byte("1"), []byte("NewVirusName"), []byte("NewQmHash")})
    if response.Status != shim.OK {
        t.Errorf("UpdateSignature failed: %s", response.Message)
    }

    // Test DeleteSignature function
    response = stub.MockInvoke("1", [][]byte{[]byte("DeleteSignature"), []byte("1")})
    if response.Status != shim.OK {
        t.Errorf("DeleteSignature failed: %s", response.Message)
    }

    // Test SignatureExists function
    response = stub.MockInvoke("1", [][]byte{[]byte("SignatureExists"), []byte("1")})
    if response.Status != shim.OK {
        t.Errorf("SignatureExists failed: %s", response.Message)
    }

    // Test GetAllSignatures function
    response = stub.MockInvoke("1", [][]byte{[]byte("GetAllSignatures")})
    if response.Status != shim.OK {
        t.Errorf("GetAllSignatures failed: %s", response.Message)
    }
}

