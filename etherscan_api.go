package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type EtherscanAPI struct {
	tokens []string
	url string
	client *http.Client
	debug  bool
}

// New create new rpc client with given url
func NewEtherscanAPI(token string, options ...func(x *EtherscanAPI)) *EtherscanAPI {
	rpc := &EtherscanAPI{
		tokens: strings.Split(token, ","),
		url:    "https://api.etherscan.io/api",
		client: http.DefaultClient,
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

func (x *EtherscanAPI) Debug(debug bool) {
	x.debug = debug
}

func (x *EtherscanAPI) String() string {
	return "etherscan"
}

func (x *EtherscanAPI) call(method string, target interface{}, params ...interface{}) error {
	result, err := x.Call(method, params)
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
func (x *EtherscanAPI) BatchCall(reqs []EthRequest) ([]EthResponse, error) {
	return nil, fmt.Errorf("TODO")
}

// Call returns raw response of method call
func (x *EtherscanAPI) Call(method string, params ...interface{}) (json.RawMessage, error) {
	retry := 0

retry:
	tokenIndex := retry % len(x.tokens)
	token := x.tokens[tokenIndex]

	if retry > 5 {
		return nil, fmt.Errorf("etherscan retry failed")
	}
	retry++

	req, err := http.NewRequest("GET", x.url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("module", "proxy")
	q.Add("action", method)
	q.Add("apikey", token)

	if params[0] != nil {
		m, ok := params[0].(map[string]string)
		if ok {
			for k, v := range m {
				fmt.Println(k + v)
			}
		}
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	response, err := x.client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	//check 403 to change apikey and retry
	if response.StatusCode == 403 {
		fmt.Printf("etherscan response 403, retrying %d...\n", retry)
		time.Sleep(3e8)
		goto retry
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code error %d", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if x.debug {
		fmt.Printf("request %s %s response %s\n", method, q.Encode(), data)
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
func (x *EtherscanAPI) Web3ClientVersion() (string, error) {
	return "", fmt.Errorf("TODO")
}

// NetVersion returns the current network protocol version.
func (x *EtherscanAPI) NetVersion() (string, error) {
	return "", fmt.Errorf("TODO")
}

// NetListening returns true if client is actively listening for network connections.
func (x *EtherscanAPI) NetListening() (bool, error) {
	return false, fmt.Errorf("TODO")
}

// NetPeerCount returns number of peers currently connected to the client.
func (x *EtherscanAPI) NetPeerCount() (int, error) {
	return 0, fmt.Errorf("TODO")
}

// EthProtocolVersion returns the current ethereum protocol version.
func (x *EtherscanAPI) EthProtocolVersion() (string, error) {
	return "", fmt.Errorf("TODO")
}

// EthSyncing returns an object with data about the sync status or false.
func (x *EtherscanAPI) EthSyncing() (*Syncing, error) {
	return nil, fmt.Errorf("TODO")
}

// EthGasPrice returns the current price per gas in wei.
func (x *EtherscanAPI) EthGasPrice() (big.Int, error) {
	var response string
	if err := x.call("eth_gasPrice", &response, nil); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthBlockNumber returns the number of most recent block.
func (x *EtherscanAPI) EthBlockNumber() (int, error) {
	var response string
	if err := x.call("eth_blockNumber", &response, nil); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBalance returns the balance of the account of given address in wei.
func (x *EtherscanAPI) EthGetBalance(address, block string) (big.Int, error) {
	return big.Int{}, fmt.Errorf("TODO")
}

// EthGetStorageAt returns the value from a storage position at a given address.
func (x *EtherscanAPI) EthGetStorageAt(data string, position int, tag string) (string, error) {
	var result string

	params := map[string]string{
		"address": data,
		"position": IntToHex(position),
		"tag": tag,
	}
	err := x.call("eth_getStorageAt", &result, params)
	return result, err
}

// EthGetTransactionCount returns the number of transactions sent from an address.
func (x *EtherscanAPI) EthGetTransactionCount(address, block string) (int, error) {
	var response string

	params := map[string]string{
		"address": address,
		"tag": "latest",
	}
	if err := x.call("eth_getTransactionCount", &response, params); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (x *EtherscanAPI) EthGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	params := map[string]string{
		"tag": IntToHex(number),
	}

	if err := x.call("eth_getBlockTransactionCountByNumber", &response, params); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (x *EtherscanAPI) EthGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	var response ProxyBlock
	if withTransactions {
		response = new(ProxyBlockWithTransactions)
	} else {
		response = new(ProxyBlockWithoutTransactions)
	}

	params := map[string]string{
		"tag": IntToHex(number),
		"boolean": fmt.Sprintf("%v", withTransactions),
	}

	err := x.call("eth_getBlockByNumber", response, params)
	if err != nil {
		return nil, err
	}
	block := response.ToBlock()
	if len(block.Hash) == 0 {
		return nil, fmt.Errorf("block not found")
	}

	return &block, nil
}

func (x *EtherscanAPI) EthGetUncleByBlockNumberAndIndex(number int, pos int) (*Block, error) {
	return nil, fmt.Errorf("TODO")
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (x *EtherscanAPI) EthGetTransactionByHash(hash string) (*Transaction, error) {
	transaction := new(Transaction)

	params := map[string]string{
		"txhash": hash,
	}

	err := x.call("eth_getTransactionByHash", transaction, params)
	if len(transaction.Hash) == 0 {
		return nil, fmt.Errorf("tx not found")
	}
	return transaction, err
}

// EthGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (x *EtherscanAPI) EthGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	transaction := new(Transaction)

	params := map[string]string{
		"tag": IntToHex(blockNumber),
		"index": IntToHex(transactionIndex),
	}

	err := x.call("eth_getTransactionByBlockNumberAndIndex", transaction, params)
	if len(transaction.Hash) == 0 {
		return nil, fmt.Errorf("tx not found")
	}
	return transaction, err
}

// EthGetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note That the receipt is not available for pending transactions.
func (x *EtherscanAPI) EthGetTransactionReceipt(hash string) (*TransactionReceipt, error) {
	transactionReceipt := new(TransactionReceipt)

	params := map[string]string{
		"txhash": hash,
	}

	err := x.call("eth_getTransactionReceipt", transactionReceipt, params)
	if err != nil {
		return nil, err
	}

	if len(transactionReceipt.TransactionHash) == 0 {
		return nil, fmt.Errorf("receipt not found")
	}

	return transactionReceipt, nil
}

// EthGetLogs returns an array of all logs matching a given filter object.
func (x *EtherscanAPI) EthGetLogs(params FilterParams) ([]Log, error) {
	var logs []Log
	return logs, fmt.Errorf("TODO")
}