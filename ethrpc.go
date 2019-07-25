package ethrpc

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// EthError - ethereum error
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EthError) Error() string {
	return fmt.Sprintf("EthError %d (%s)", err.Code, err.Message)
}

type EthResponse struct {
	ID      int             `json:"id,omitempty"`
	Version string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type EthRequest struct {
	ID      int           `json:"id,omitempty"`
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewEthRequest(method string, params ...interface{}) EthRequest {
	return EthRequest{
		ID:      1,
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}

var (
	_ EthRPC = (*NodeAPI)(nil)
	_ EthRPC = (*InfuraAPI)(nil)
	_ EthRPC = (*EtherscanAPI)(nil)
)

type EthRPC interface {
	String() string
	Debug(bool)
	BatchCall(reqs []EthRequest) ([]EthResponse, error)
	Call(method string, params ...interface{}) (json.RawMessage, error)
	Web3ClientVersion() (string, error)
	NetVersion() (string, error)
	NetListening() (bool, error)
	NetPeerCount() (int, error)
	EthProtocolVersion() (string, error)
	EthSyncing() (*Syncing, error)
	EthGasPrice() (big.Int, error)
	EthBlockNumber() (int, error)
	EthGetBalance(address, block string) (big.Int, error)
	EthGetStorageAt(data string, position int, tag string) (string, error)
	EthGetTransactionCount(address, block string) (int, error)
	EthGetBlockTransactionCountByNumber(number int) (int, error)
	EthGetBlockByNumber(number int, withTransactions bool) (*Block, error)
	EthGetTransactionByHash(hash string) (*Transaction, error)
	EthGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error)
	EthGetTransactionReceipt(hash string) (*TransactionReceipt, error)
	EthGetLogs(params FilterParams) ([]Log, error)
	EthGetUncleByBlockNumberAndIndex(number int, pos int) (*Block, error)
}