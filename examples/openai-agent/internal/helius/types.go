package helius

// EnhancedTransaction represents a single parsed transaction from the Helius API.
type EnhancedTransaction struct {
	Description      string           `json:"description"`
	Type             string           `json:"type"`
	Source           string           `json:"source"`
	Fee              int64            `json:"fee"`
	FeePayer         string           `json:"feePayer"`
	Signature        string           `json:"signature"`
	Slot             int64            `json:"slot"`
	Timestamp        int64            `json:"timestamp"`
	NativeTransfers  []NativeTransfer `json:"nativeTransfers"`
	TokenTransfers   []TokenTransfer  `json:"tokenTransfers"`
	AccountData      []AccountData    `json:"accountData"`
	TransactionError *TxError         `json:"transactionError"`
	Events           Events           `json:"events"`
}

// NativeTransfer represents a SOL transfer between accounts.
type NativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
	Amount          int64  `json:"amount"` // lamports
}

// TokenTransfer represents a token transfer between accounts.
type TokenTransfer struct {
	FromUserAccount  string  `json:"fromUserAccount"`
	ToUserAccount    string  `json:"toUserAccount"`
	FromTokenAccount string  `json:"fromTokenAccount"`
	ToTokenAccount   string  `json:"toTokenAccount"`
	TokenAmount      float64 `json:"tokenAmount"`
	Mint             string  `json:"mint"`
}

// AccountData holds account-level balance change data.
type AccountData struct {
	Account              string               `json:"account"`
	NativeBalanceChange  int64                `json:"nativeBalanceChange"`
	TokenBalanceChanges  []TokenBalanceChange `json:"tokenBalanceChanges"`
}

// TokenBalanceChange represents a token balance change for one account.
type TokenBalanceChange struct {
	UserAccount    string         `json:"userAccount"`
	TokenAccount   string         `json:"tokenAccount"`
	Mint           string         `json:"mint"`
	RawTokenAmount RawTokenAmount `json:"rawTokenAmount"`
}

// RawTokenAmount holds a raw token amount with its decimals.
type RawTokenAmount struct {
	TokenAmount string `json:"tokenAmount"`
	Decimals    int    `json:"decimals"`
}

// TxError represents a transaction error.
type TxError struct {
	Error string `json:"error"`
}

// Events holds the structured event data parsed by Helius.
type Events struct {
	Swap *SwapEvent `json:"swap"`
}

// SwapEvent represents a parsed swap event from the transaction.
type SwapEvent struct {
	NativeInput  *NativeAmount  `json:"nativeInput"`
	NativeOutput *NativeAmount  `json:"nativeOutput"`
	TokenInputs  []SwapToken    `json:"tokenInputs"`
	TokenOutputs []SwapToken    `json:"tokenOutputs"`
	TokenFees    []SwapToken    `json:"tokenFees"`
	NativeFees   []NativeAmount `json:"nativeFees"`
	InnerSwaps   []InnerSwap    `json:"innerSwaps"`
}

// NativeAmount represents a native SOL amount tied to an account.
type NativeAmount struct {
	Account string `json:"account"`
	Amount  string `json:"amount"` // lamports as string
}

// SwapToken represents a token involved in a swap event.
type SwapToken struct {
	UserAccount    string         `json:"userAccount"`
	TokenAccount   string         `json:"tokenAccount"`
	Mint           string         `json:"mint"`
	RawTokenAmount RawTokenAmount `json:"rawTokenAmount"`
}

// InnerSwap represents a single hop within a multi-hop swap.
type InnerSwap struct {
	TokenInputs  []TokenTransfer `json:"tokenInputs"`
	TokenOutputs []TokenTransfer `json:"tokenOutputs"`
	TokenFees    []TokenTransfer `json:"tokenFees"`
	NativeFees   []NativeTransfer `json:"nativeFees"`
	ProgramInfo  *ProgramInfo    `json:"programInfo"`
}

// ProgramInfo identifies the DEX program used in an inner swap.
type ProgramInfo struct {
	Source          string `json:"source"`
	Account         string `json:"account"`
	ProgramName     string `json:"programName"`
	InstructionName string `json:"instructionName"`
}
