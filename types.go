package ethrpc

import (
	"bytes"
	"encoding/json"
	"math/big"
	"unsafe"
)

// Syncing - object with syncing data info
type Syncing struct {
	IsSyncing     bool
	StartingBlock int
	CurrentBlock  int
	HighestBlock  int
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Syncing) UnmarshalJSON(data []byte) error {
	proxy := new(proxySyncing)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	proxy.IsSyncing = true
	*s = *(*Syncing)(unsafe.Pointer(proxy))

	return nil
}

// Transaction - transaction object
type Transaction struct {
	Hash             string
	Nonce            int
	BlockHash        string
	BlockNumber      *int
	TransactionIndex *int
	From             string
	To               string
	Value            big.Int
	Gas              int
	GasPrice         big.Int
	Input            string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Transaction) UnmarshalJSON(data []byte) error {
	proxy := new(proxyTransaction)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	*t = *(*Transaction)(unsafe.Pointer(proxy))

	return nil
}

// Log - log object
type Log struct {
	Removed          bool
	LogIndex         int
	TransactionIndex int
	TransactionHash  string
	BlockNumber      int
	BlockHash        string
	Address          string
	Data             string
	Topics           []string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (log *Log) UnmarshalJSON(data []byte) error {
	proxy := new(proxyLog)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	*log = *(*Log)(unsafe.Pointer(proxy))

	return nil
}

// FilterParams - Filter parameters object
type FilterParams struct {
	FromBlock string     `json:"fromBlock,omitempty"`
	ToBlock   string     `json:"toBlock,omitempty"`
	Address   []string   `json:"address,omitempty"`
	Topics    [][]string `json:"topics,omitempty"`
}

// TransactionReceipt - transaction receipt object
type TransactionReceipt struct {
	TransactionHash   string
	TransactionIndex  int
	BlockHash         string
	BlockNumber       int
	CumulativeGasUsed int
	GasUsed           int
	ContractAddress   string
	Logs              []Log
	LogsBloom         string
	Root              string
	Status            int
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionReceipt) UnmarshalJSON(data []byte) error {
	proxy := new(proxyTransactionReceipt)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	*t = *(*TransactionReceipt)(unsafe.Pointer(proxy))

	return nil
}

// Block - block object
type Block struct {
	Number           int
	Hash             string
	ParentHash       string
	Nonce            string
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	Miner            string
	Difficulty       big.Int
	TotalDifficulty  big.Int
	ExtraData        string
	Size             int
	GasLimit         int
	GasUsed          int
	Timestamp        int
	Uncles           []string
	Transactions     []Transaction
}

// {"difficulty":"0xcb5d1dadda318",
// "extraData":"0x737061726b706f6f6c2d636e2d6e6f64652d3032",
// "gasLimit":"0x79b6ea","gasUsed":"0x5bff7f",
// "hash":"0xef7fa50f455e5c40f3435f2e1ede71fabf670f265ad623fb804ae2eb3c1d0db3",
// "logsBloom":"0xc10234002148000000898a09010090001000d00c08001041402285024029034930020098803100080090804027430682a2
// 000008a80100100820a02400382402849050c0801120184080510a002b508062a00004106009c08000450440400220092020220420040026198040000a0100595
// 20401084505001000001084802102006208124218000000020c408c200400281111000000042018018304010022005828100042c208189bc600001442
// 2010200805000040c000082200c006208910a081101388020109402808b10009a444088290a4c638030a000a88270002100568002101000200000020008480
// 0200044004001503102829040009a000101080","miner":"0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c",
// "mixHash":"0xbe59def406d63402e08374dea3bc803aaa1f341f3fc50b66f38a1cea85401ab5",
// "nonce":"0xe764982ec0e9a94e","number":"0x5dd091",
// "parentHash":"0x0bc5f310f4017a7add06b1b0419beeb0cf2bd667be1f25e7aa348b4dd51ab4a5",
// "receiptsRoot":"0xa2cef1828213fa7c2fd568e10fe42cef90194889167f6faf12dc1efb4524ac28",
// "sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
// "size":"0x21d","stateRoot":"0xad23b36dbaf20fe8387307b733e1279e7aa41e4ffc14cff8a08da88df65885b0",
// "timestamp":"0x5b734e23","totalDifficulty":null,"transactionsRoot":"0xf4ce86641f301d2ee4be3397d22dec3dd101edb9dce0cb8baed213b083bd9e55","uncles":[]}}

//request eth_getUncleByBlockNumberAndIndex {"id":1,"jsonrpc":"2.0","method":"eth_getUncleByBlockNumberAndIndex","params":["0x65b652","0x0"]}
// response {"jsonrpc":"2.0","result":{"author":"0xf35074bbd0a9aee46f4ea137971feec024ab704e","difficulty":"0xae765ab687703",
// "extraData":"0x736f6c6f706f6f6c2e6f7267","gasLimit":"0x7a121d","gasUsed":"0x79e444","hash":"0x52fa126b7bbad9a6f9235bc57163881e17eece47bd743f6021ea57eddc6de5d2",
// "logsBloom":"0x0a2c304d02e20f1a80420484684d1242a3a8248800406008a741c700c02e218182c088a12084cd21040183021343089c8215740108a0007108c02001002460104008a64708c1801900c03108018
// 64ea2024e00600254312902008c870081408218315220a5c031b601a9600890400494041012080054a0601cc232182080000000012d10004f014a10406498c08044380013400480400081280260212081821502844a00064d10
// 8821c6842d0523850024200601b021980918006408182001021520400390c8010551092a300482bc06396a010b82490c018864060d6910300224193244d2002116040400cf404084150192c826414a031402e0620400000082",
// "miner":"0xf35074bbd0a9aee46f4ea137971feec024ab704e","mixHash":"0xf9e3975ac162d3ce75b6d8abc6899dc6858f5da2c1fa950e4ea1ae77eba106e3","nonce":"0x7cd838000635bb78",
// "number":"0x65b64f","parentHash":"0xfe76b816fc985a131a8e4a9bf77be997e89a3b7433324f9b0d783950a6a1f836","receiptsRoot":"0xde022e159186410c56ce2ff6dfe42d6b35880e999c99ff789335368ca99359f6",
// "sealFields":["0xa0f9e3975ac162d3ce75b6d8abc6899dc6858f5da2c1fa950e4ea1ae77eba106e3","0x887cd838000635bb78"],"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
// "size":null,"stateRoot":"0xe9340e24e817e41353e94145a8c807d7b96284f761900c94da88cb73f1e65549","timestamp":"0x5be412f4","totalDifficulty":"0x1a0e5a2ee25b0583f63",
// "transactions":[],"transactionsRoot":"0x74824e98fe52018ed3269bb4483197ed1155186a4ad36b923632a76d0b257480","uncles":[]},"id":1}

type UncleBlock struct {
	Number           hexInt `json:"number"`
	Hash             string `json:"hash"`
	ParentHash       string `json:"parentHash"`
	Nonce            string `json:"nonce"`
	Sha3Uncles       string `json:"sha3Uncles"`
	LogsBloom        string `json:"logsBloom"`
	TransactionsRoot string `json:"transactionsRoot"`
	StateRoot        string `json:"stateRoot"`
	Miner            string `json:"miner"`
	Difficulty       hexBig `json:"difficulty"`
}

func (proxy *UncleBlock) ToBlock() Block {
	return Block{
		Number:           int(proxy.Number),
		Hash:             proxy.Hash,
		ParentHash:       proxy.ParentHash,
		Nonce:            proxy.Nonce,
		Sha3Uncles:       proxy.Sha3Uncles,
		LogsBloom:        proxy.LogsBloom,
		TransactionsRoot: proxy.TransactionsRoot,
		StateRoot:        proxy.StateRoot,
		Miner:            proxy.Miner,
		Difficulty:       big.Int(proxy.Difficulty),
	}
}

type proxySyncing struct {
	IsSyncing     bool   `json:"-"`
	StartingBlock hexInt `json:"startingBlock"`
	CurrentBlock  hexInt `json:"currentBlock"`
	HighestBlock  hexInt `json:"highestBlock"`
}

type proxyTransaction struct {
	Hash             string  `json:"hash"`
	Nonce            hexInt  `json:"nonce"`
	BlockHash        string  `json:"blockHash"`
	BlockNumber      *hexInt `json:"blockNumber"`
	TransactionIndex *hexInt `json:"transactionIndex"`
	From             string  `json:"from"`
	To               string  `json:"to"`
	Value            hexBig  `json:"value"`
	Gas              hexInt  `json:"gas"`
	GasPrice         hexBig  `json:"gasPrice"`
	Input            string  `json:"input"`
}

type proxyLog struct {
	Removed          bool     `json:"removed"`
	LogIndex         hexInt   `json:"logIndex"`
	TransactionIndex hexInt   `json:"transactionIndex"`
	TransactionHash  string   `json:"transactionHash"`
	BlockNumber      hexInt   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	Address          string   `json:"address"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
}

type proxyTransactionReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  hexInt `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       hexInt `json:"blockNumber"`
	CumulativeGasUsed hexInt `json:"cumulativeGasUsed"`
	GasUsed           hexInt `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	Logs              []Log  `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Root              string `json:"root"`
	Status            hexInt `json:"status"`
}

type hexInt int

func (i *hexInt) UnmarshalJSON(data []byte) error {
	result, err := ParseInt(string(bytes.Trim(data, `"`)))
	*i = hexInt(result)

	return err
}

type hexBig big.Int

func (i *hexBig) UnmarshalJSON(data []byte) error {
	result, err := ParseBigInt(string(bytes.Trim(data, `"`)))
	*i = hexBig(result)

	return err
}

type ProxyBlockWithTransactions struct {
	Number           hexInt             `json:"number"`
	Hash             string             `json:"hash"`
	ParentHash       string             `json:"parentHash"`
	Nonce            string             `json:"nonce"`
	Sha3Uncles       string             `json:"sha3Uncles"`
	LogsBloom        string             `json:"logsBloom"`
	TransactionsRoot string             `json:"transactionsRoot"`
	StateRoot        string             `json:"stateRoot"`
	Miner            string             `json:"miner"`
	Difficulty       hexBig             `json:"difficulty"`
	TotalDifficulty  hexBig             `json:"totalDifficulty"`
	ExtraData        string             `json:"extraData"`
	Size             hexInt             `json:"size"`
	GasLimit         hexInt             `json:"gasLimit"`
	GasUsed          hexInt             `json:"gasUsed"`
	Timestamp        hexInt             `json:"timestamp"`
	Uncles           []string           `json:"uncles"`
	Transactions     []proxyTransaction `json:"transactions"`
}

func (proxy *ProxyBlockWithTransactions) ToBlock() Block {
	return *(*Block)(unsafe.Pointer(proxy))
}

type ProxyBlock interface {
	ToBlock() Block
}

type ProxyBlockWithoutTransactions struct {
	Number           hexInt   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       hexBig   `json:"difficulty"`
	TotalDifficulty  hexBig   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             hexInt   `json:"size"`
	GasLimit         hexInt   `json:"gasLimit"`
	GasUsed          hexInt   `json:"gasUsed"`
	Timestamp        hexInt   `json:"timestamp"`
	Uncles           []string `json:"uncles"`
	Transactions     []string `json:"transactions"`
}

func (proxy *ProxyBlockWithoutTransactions) ToBlock() Block {
	block := Block{
		Number:           int(proxy.Number),
		Hash:             proxy.Hash,
		ParentHash:       proxy.ParentHash,
		Nonce:            proxy.Nonce,
		Sha3Uncles:       proxy.Sha3Uncles,
		LogsBloom:        proxy.LogsBloom,
		TransactionsRoot: proxy.TransactionsRoot,
		StateRoot:        proxy.StateRoot,
		Miner:            proxy.Miner,
		Difficulty:       big.Int(proxy.Difficulty),
		TotalDifficulty:  big.Int(proxy.TotalDifficulty),
		ExtraData:        proxy.ExtraData,
		Size:             int(proxy.Size),
		GasLimit:         int(proxy.GasLimit),
		GasUsed:          int(proxy.GasUsed),
		Timestamp:        int(proxy.Timestamp),
		Uncles:           proxy.Uncles,
	}

	block.Transactions = make([]Transaction, len(proxy.Transactions))
	for i := range proxy.Transactions {
		block.Transactions[i] = Transaction{
			Hash: proxy.Transactions[i],
		}
	}

	return block
}
