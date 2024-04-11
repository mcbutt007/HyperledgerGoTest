package main

import (
    "testing"
    "fmt"
    "strings"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestInstancesCreation(test *testing.T) {
    stub := InitChaincode(test)

    assetExternalId := "ID01"
    ownerId := "o1"
    ownerId2 := "o2"

    Invoke(test, stub, "create_owner", ownerId, "Username_1", "Company_1")
    Invoke(test, stub, "read_owner", ownerId)
    Invoke(test, stub, "create_asset", assetExternalId, "Sernr1234", "Matnr1234", "ObjDesc", ownerId)
    Invoke(test, stub, "read_asset", assetExternalId)
    Invoke(test, stub, "create_owner", ownerId2, "Username_2", "Company_2")
    Invoke(test, stub, "set_owner", assetExternalId, ownerId2)
    Invoke(test, stub, "read_asset", assetExternalId)
}

func InitChaincode(test *testing.T) *shimtest.MockStub {
    stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))
    result := stub.MockInit("000", nil)

    if result.Status != shim.OK {
       test.FailNow()
    }
    return stub
}

func Invoke(test *testing.T, stub *shimtest.MockStub, function string, args ...string) {

    cc_args := make([][]byte, 1+len(args))
    cc_args[0] = []byte(function)
    for i, arg := range args {
        cc_args[i + 1] = []byte(arg)
    }
    result := stub.MockInvoke("000", cc_args)
    fmt.Println("Call:    ", function, "(", strings.Join(args,","), ")")
    fmt.Println("RetCode: ", result.Status)
    fmt.Println("RetMsg:  ", result.Message)
    fmt.Println("Payload: ", string(result.Payload))
    fmt.Println()

    if result.Status != shim.OK {
        test.FailNow()
    }
}
