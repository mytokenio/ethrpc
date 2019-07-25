package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

type NodeAPI struct {
	url    string
	client *http.Client
	debug  bool
}

// New create new rpc client with given url
func NewNodeAPI(url string, options ...func(x *NodeAPI)) *NodeAPI {
	rpc := &NodeAPI{
		url:    url,
		client: http.DefaultClient,
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

func (x *NodeAPI) String() string {
	return "node-" + x.url
}

func (x *NodeAPI) Debug(debug bool) {
	x.debug = debug
}

func (x *NodeAPI) call(method string, target interface{}, params ...interface{}) error {
	result, err := x.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	if bytes.Equal(result, []byte("null")) {
		return fmt.Errorf("result null")
	}
	return json.Unmarshal(result, target)
}

//batch request
func (x *NodeAPI) BatchCall(reqs []EthRequest) ([]EthResponse, error) {
	body, err := json.Marshal(reqs)
	if err != nil {
		return nil, err
	}

	response, err := x.client.Post(x.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code error %d", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if x.debug {
		fmt.Printf("requests %s responses %s\n", body, data)
	}

	resps := &[]EthResponse{}
	if err := json.Unmarshal(data, resps); err != nil {
		return nil, err
	}
	return *resps, nil
}

// Call returns raw response of method call
func (x *NodeAPI) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := EthRequest{
		ID:      1,
		Version: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := x.client.Post(x.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code error %d", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if x.debug {
		fmt.Printf("request %s %s response %s\n", method, body, data)
	}

	resp := new(EthResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}

// Web3ClientVersion returns the current client version.
func (x *NodeAPI) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := x.call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// NetVersion returns the current network protocol version.
func (x *NodeAPI) NetVersion() (string, error) {
	var version string

	err := x.call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (x *NodeAPI) NetListening() (bool, error) {
	var listening bool

	err := x.call("net_listening", &listening)
	return listening, err
}

// NetPeerCount returns number of peers currently connected to the client.
func (x *NodeAPI) NetPeerCount() (int, error) {
	var response string
	if err := x.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthProtocolVersion returns the current ethereum protocol version.
func (x *NodeAPI) EthProtocolVersion() (string, error) {
	var protocolVersion string

	err := x.call("eth_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

// EthSyncing returns an object with data about the sync status or false.
func (x *NodeAPI) EthSyncing() (*Syncing, error) {
	result, err := x.Call("eth_syncing")
	if err != nil {
		return nil, err
	}
	syncing := new(Syncing)
	if bytes.Equal(result, []byte("false")) {
		return syncing, nil
	}
	err = json.Unmarshal(result, syncing)
	return syncing, err
}

// EthGasPrice returns the current price per gas in wei.
func (x *NodeAPI) EthGasPrice() (big.Int, error) {
	var response string
	if err := x.call("eth_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthBlockNumber returns the number of most recent block.
func (x *NodeAPI) EthBlockNumber() (int, error) {
	var response string
	if err := x.call("eth_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBalance returns the balance of the account of given address in wei.
func (x *NodeAPI) EthGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := x.call("eth_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthGetStorageAt returns the value from a storage position at a given address.
func (x *NodeAPI) EthGetStorageAt(data string, position int, tag string) (string, error) {
	var result string

	err := x.call("eth_getStorageAt", &result, data, IntToHex(position), tag)
	return result, err
}

// EthGetTransactionCount returns the number of transactions sent from an address.
func (x *NodeAPI) EthGetTransactionCount(address, block string) (int, error) {
	var response string

	if err := x.call("eth_getTransactionCount", &response, address, block); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (x *NodeAPI) EthGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	if err := x.call("eth_getBlockTransactionCountByNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (x *NodeAPI) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	var response ProxyBlock
	if withTransactions {
		response = new(ProxyBlockWithTransactions)
	} else {
		response = new(ProxyBlockWithoutTransactions)
	}

	err := x.call(method, response, params...)
	if err != nil {
		return nil, err
	}
	block := response.ToBlock()
	if len(block.Hash) == 0 {
		return nil, fmt.Errorf("block not found")
	}

	return &block, nil
}

// EthGetBlockByNumber returns information about a block by block number.
func (x *NodeAPI) EthGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return x.getBlock("eth_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
}

func (x *NodeAPI) EthGetUncleByBlockNumberAndIndex(number int, pos int) (*Block, error) {
	uncleBlock := new(UncleBlock)

	err := x.call("eth_getUncleByBlockNumberAndIndex", uncleBlock, IntToHex(number), IntToHex(pos))
	if err != nil {
		return nil, err
	}

	block := uncleBlock.ToBlock()
	if len(block.Hash) == 0 {
		return nil, fmt.Errorf("block not found")
	}

	return &block, nil
}

func (x *NodeAPI) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := x.call(method, transaction, params...)
	if len(transaction.Hash) == 0 {
		return nil, fmt.Errorf("tx not found")
	}
	return transaction, err
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (x *NodeAPI) EthGetTransactionByHash(hash string) (*Transaction, error) {
	return x.getTransaction("eth_getTransactionByHash", hash)
}

// EthGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (x *NodeAPI) EthGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	return x.getTransaction("eth_getTransactionByBlockNumberAndIndex", IntToHex(blockNumber), IntToHex(transactionIndex))
}

// EthGetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note That the receipt is not available for pending transactions.
func (x *NodeAPI) EthGetTransactionReceipt(hash string) (*TransactionReceipt, error) {
	transactionReceipt := new(TransactionReceipt)

	err := x.call("eth_getTransactionReceipt", transactionReceipt, hash)
	if err != nil {
		return nil, err
	}

	if len(transactionReceipt.TransactionHash) == 0 {
		return nil, fmt.Errorf("receipt not found")
	}

	return transactionReceipt, nil
}

// EthGetLogs returns an array of all logs matching a given filter object.
func (x *NodeAPI) EthGetLogs(params FilterParams) ([]Log, error) {
	var logs []Log
	err := x.call("eth_getLogs", &logs, params)
	return logs, err
}
