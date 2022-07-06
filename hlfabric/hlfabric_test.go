package hlfabric_test

import (
	"testing"

	"github.com/ewangplay/gokv/hlfabric"
	"github.com/philippgille/gokv/test"
)

func TestClient(t *testing.T) {
	if !checkTestNetwork() {
		t.Skip("No fabirc test network running")
	}

	// Test NewClient with hyperledger fabric test network
	t.Run("hlfabric-test-network", func(t *testing.T) {
		client := createClient(t)
		defer client.Close()
		test.TestStore(client, t)
	})
}

func createClient(t *testing.T) *hlfabric.Client {
	opts := hlfabric.Options{
		ChannelName: "mychannel",
		ContractID:  "kvstore",
		MspID:       "Org1MSP",
		WalletPath:  "sampleconfig/wallet",
		CcpPath:     "sampleconfig/connection-org1.yaml",
		AppUser: hlfabric.AppUser{
			Name:    "appUser",
			MspPath: "sampleconfig/appUser/msp",
		},
		EndorsingPeers: []string{"peer0.org1.example.com:7051"},
	}
	client, err := hlfabric.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// checkTestNetwork returns true if a hyperledger fabric test network is running
func checkTestNetwork() bool {
	return false
}
