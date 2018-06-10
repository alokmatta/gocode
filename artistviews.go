package main

import (
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type ArtistViewsChaincode struct {
}

//var A, B string
//var Aval, Bval, X int
var Count int = 0

// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *ArtistViewsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	
	return shim.Success(nil)
}

func (t *ArtistViewsChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ArtistViewsChaincode Invoke Called")
	fmt.Print("Artist Count ")
	count = count + 1
	fmt.Print(count)
	return shim.Success(nil)
}

func (t *ArtistViewsChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\"")
}

func main() {
	err := shim.Start(new(ArtistViewsChaincode))
	if err != nil {
		fmt.Printf("Error starting ArtistViewsChaincode: %s", err)
	}
}
