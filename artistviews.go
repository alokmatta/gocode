package main

import (
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	
	"github.com/dgryski/go-topk"
)

// SimpleChaincode example simple Chaincode implementation
type ArtistViewsChaincode struct {
}

//var A, B string
//var Aval, Bval, X int
//var Count int = 0
k := flag.Int("n", 10, "k")
tk := topk.New(*k)
var top10 [10]string
var newTop10 [10]string

// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *ArtistViewsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	
	return shim.Success(nil)
}

func (t *ArtistViewsChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ArtistViewsChaincode Invoke Called")
	//fmt.Print("Artist Count ")
	//Count = Count + 1
	//fmt.Print(Count)
	item := args[0]
	tk.Insert(item, 1)
	
	KeysArray := tk.Keys()
	
	if len(top10) > 0 {
			if top10[0] != KeysArray[0].Key {
				fmt.Println("New number 1 - " + KeysArray[0].Key)
			}
		}

		for i, v := range KeysArray {

			if !contains(v.Key, top10) {
				fmt.Println("New Artist in top 10 - " + v.Key)
			}

			newTop10[i] = KeysArray[i].Key

		}
		top10 = newTop10
		fmt.Print("Current Top Ten ")
		fmt.Println(top10)
	
		for _, v := range tk.Keys() {
		fmt.Println(v.Key, v.Count, v.Error)
	}
	
	return shim.Success(nil)
}

func contains(e string, s [10]string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
