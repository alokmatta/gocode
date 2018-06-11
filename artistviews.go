/*

chaincode using a Top K Algorithm to maintain a list of top ten Artists and notify them accordingly.
Limited hyperledger fabric PoC as the calculation is done during invoke rather a more ideal offchain Blockchain listener microservice (using the events SDK) that runs once
the View transaction is actually verified. Also ideally should have PubSub layer between the Blockchain listener and count microservice for better reliability & abstraction from any
blockchain platform.

CLI Commands to test (via https://github.com/hyperledger/fabric-samples/tree/release-1.1/first-network) on Azure Cloud VM with state stored on CouchDB

peer chaincode install -n artistViews -v 1 -p github.com/gocode/

peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n artistViews -v 1 -c '{"Args":[""]}' -P 'OR('\''Org1MSP.peer'\'','\''Org2MSP.peer'\'')'

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n artistViews -c '{"Args":["invoke","Coldplay","1"]}'

Sample Output by attaching to the Instantiated Docker instance:
root@everyvm:/home/every# docker attach 74e1270ecc8a
ArtistViewsChaincode Invoke Called
Current Top Ten [Britney         ]
Britney 2 0
ArtistViewsChaincode Invoke Called
New Artist in top 10 - Coldplay
Current Top Ten [Britney Coldplay        ]
Britney 2 0
Coldplay 1 0

*/

package main

import (
	"flag"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/dgryski/go-topk"
)

// ArtistViewsChaincode Chaincode implementation
type ArtistViewsChaincode struct {
}

var top10 [10]string
var newTop10 [10]string
var tk *topk.Stream
var k *int

// Init callback representing the invocation of a chaincode
func (t *ArtistViewsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	k := flag.Int("n", 10, "k")
	tk = topk.New(*k)
	return shim.Success(nil)
}

func (t *ArtistViewsChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("ArtistViewsChaincode Invoke Called")
        var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting Artist & Listener")
	}
        artist := args[0]
        listener := args[1]


        artistViewJSONasString := `{"docType":"ArtistView",  "Artist": "` + artist + `", "listener": "` + listener + `"}`
	artistViewJSONasBytes := []byte(artistViewJSONasString)

	// === Save ArtistView to state ===
	err = stub.PutState("artistView", artistViewJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
        }
        
	tk.Insert(artist, 1)

	KeysArray := tk.Keys()

	if len(top10) > 0 {
		if top10[0] != KeysArray[0].Key {
			fmt.Println("New number 1 - " + KeysArray[0].Key)
			fmt.Println(KeysArray[0].Key + " Notified as new Number 1 Artist")
			//Actually notify user by SMS and Email here using Twilio/SendGrid API's
		}
	}

	for i, v := range KeysArray {

		if !contains(v.Key, top10) {
			fmt.Println("New Artist in top 10 - " + v.Key)
			fmt.Println(v.Key + " Notified as Top 10 Artist")
			//Actually notify user by SMS and Email here using Twilio/SendGrid API's
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
