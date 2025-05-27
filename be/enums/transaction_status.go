package enums

type TransactionStatus string

const (
	TransactionPENDING TransactionStatus = `PENDING`
	TransactionSUCCESS TransactionStatus = `SUCCESS`
	TransactionFAILED  TransactionStatus = `FAILED`
)

func (s TransactionStatus) String() string {
	return string(s)
}
