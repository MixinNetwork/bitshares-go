package types

import (
	"github.com/MixinNetwork/bitshares-go/encoding/transaction"
	"github.com/pkg/errors"
)

type Transaction struct {
	RefBlockNum    int64      `json:"ref_block_num"`
	RefBlockPrefix int64      `json:"ref_block_prefix"`
	Expiration     Time       `json:"expiration"`
	Operations     Operations `json:"operations"`
	Signatures     []string   `json:"signatures"`
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoder *transaction.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := transaction.NewRollingEncoder(encoder)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op)
	}

	// Extensions are not supported yet.
	enc.EncodeUVarint(0)
	return enc.Err()
}

// PushOperation can be used to add an operation into the transaction.
func (tx *Transaction) PushOperation(op Operation) {
	tx.Operations = append(tx.Operations, op)
}
