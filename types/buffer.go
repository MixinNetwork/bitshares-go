package types

import (
	"encoding/hex"
	"encoding/json"

	"github.com/MixinNetwork/bitshares-go/encoding/transaction"
	"github.com/pkg/errors"
)

type Buffer []byte

func (p *Buffer) UnmarshalJSON(data []byte) error {
	var b string
	if err := json.Unmarshal(data, &b); err != nil {
		return errors.Wrap(err, "Unmarshal")
	}

	return p.FromString(b)
}

func (p Buffer) Bytes() []byte {
	return p
}

func (p Buffer) Length() int {
	return len(p)
}

func (p Buffer) String() string {
	return hex.EncodeToString(p)
}

func (p *Buffer) FromString(data string) error {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return errors.Wrap(err, "DecodeString")
	}

	*p = buf
	return nil
}

func (p Buffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Buffer) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Wrap(err, "encode length")
	}

	if err := enc.Encode(p.Bytes()); err != nil {
		return errors.Wrap(err, "encode bytes")
	}

	return nil
}
