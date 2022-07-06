package hlfabric

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/philippgille/gokv/util"
)

// AppUser defines hlfabric app user
type AppUser struct {
	Name    string
	MspPath string
}

// Options are the options for the HLFabric client
type Options struct {
	ChannelName    string
	ContractID     string
	MspID          string
	WalletPath     string
	CcpPath        string
	AppUser        AppUser
	EndorsingPeers []string
}

// Client represents hyperledger fabric blockchain client
type Client struct {
	opts     *Options
	wallet   *gateway.Wallet
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
}

// You must call the Close() method on the client when you're done working with it.
func NewClient(opts Options) (*Client, error) {

	wallet, err := gateway.NewFileSystemWallet(opts.WalletPath)
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		return nil, err
	}

	if !wallet.Exists(opts.AppUser.Name) {
		err = populateWallet(wallet, opts.MspID, &opts.AppUser)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			return nil, err
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(opts.CcpPath))),
		gateway.WithIdentity(wallet, opts.AppUser.Name),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		return nil, err
	}

	network, err := gw.GetNetwork(opts.ChannelName)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		return nil, err
	}

	contract := network.GetContract(opts.ContractID)

	c := &Client{
		opts:     &opts,
		wallet:   wallet,
		gateway:  gw,
		network:  network,
		contract: contract,
	}

	return c, nil
}

func populateWallet(wallet *gateway.Wallet, mspID string, appUser *AppUser) error {
	// read the certificate pem
	certPath := filepath.Join(appUser.MspPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(appUser.MspPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	err = wallet.Put(appUser.Name, identity)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.Submit("Set", k, string(data))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	data, err := c.Evaluate("Get", k)
	if err != nil {
		return false, err
	}
	if len(data) == 0 {
		return false, nil
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) Delete(k string) error {
	var err error
	if err = util.CheckKey(k); err != nil {
		return err
	}

	_, err = c.Submit("Delete", k)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return nil
}

// Submit will submit a transaction to the ledger
func (c *Client) Submit(fn string, args ...string) ([]byte, error) {
	var err error
	var txn *gateway.Transaction

	if len(c.opts.EndorsingPeers) > 0 {
		txn, err = c.contract.CreateTransaction(fn,
			gateway.WithEndorsingPeers(c.opts.EndorsingPeers...))
	} else {
		txn, err = c.contract.CreateTransaction(fn)
	}

	if err != nil {
		return nil, err
	}

	return txn.Submit(args...)
}

// Evaluate will evaluate a transaction function and return its results
func (c *Client) Evaluate(fn string, args ...string) ([]byte, error) {
	return c.contract.EvaluateTransaction(fn, args...)
}
