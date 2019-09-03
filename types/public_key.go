package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/MixinNetwork/bitshares-go/encoding/transaction"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ripemd160"
)

var (
	ErrInvalidPublicKey             = fmt.Errorf("invalid PublicKey")
	ErrPublicKeyChainPrefixMismatch = fmt.Errorf("PublicKey chain prefix mismatch")
)

type PublicKey struct {
	key      *btcec.PublicKey
	prefix   string
	checksum []byte
}

func (p PublicKey) String() string {
	b := append(p.Bytes(), p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var key string

	if err := json.Unmarshal(data, &key); err != nil {
		return errors.Wrap(err, "Unmarshal")
	}

	pub, err := NewPublicKeyFromString(key)
	if err != nil {
		return errors.Wrap(err, "NewPublicKeyFromString")
	}

	p.key = pub.key
	p.prefix = pub.prefix
	p.checksum = pub.checksum
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p PublicKey) MarshalTransaction(enc *transaction.Encoder) error {
	return enc.Encode(p.Bytes())
}

func (p PublicKey) Bytes() []byte {
	return p.key.SerializeCompressed()
}

func (p PublicKey) Equal(pub *PublicKey) bool {
	return p.key.IsEqual(pub.key)
}

func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	return p.key.ToECDSA()
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func (p PublicKey) MaxSharedKeyLength() int {
	return (p.key.ToECDSA().Curve.Params().BitSize + 7) / 8
}

func NewPublicKeyFromString(key string) (*PublicKey, error) {
	return NewPublicKeyWithPrefixFromString(key, "BTS")
}

//NewPublicKey creates a new PublicKey from string
func NewPublicKeyWithPrefixFromString(key, prefixChain string) (*PublicKey, error) {
	prefix := key[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, ErrPublicKeyChainPrefixMismatch
	}

	b58 := base58.Decode(key[len(prefixChain):])
	if len(b58) < 5 {
		return nil, ErrInvalidPublicKey
	}

	chk1 := b58[len(b58)-4:]

	keyBytes := b58[:len(b58)-4]
	chk2, err := Ripemd160Checksum(keyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Ripemd160Checksum")
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, ErrInvalidPublicKey
	}

	pub, err := btcec.ParsePubKey(keyBytes, btcec.S256())
	if err != nil {
		return nil, errors.Wrap(err, "ParsePubKey")
	}

	k := PublicKey{
		key:      pub,
		prefix:   prefix,
		checksum: chk1,
	}

	return &k, nil
}

func Ripemd160(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, errors.Wrap(err, "Write")
	}

	sum := h.Sum(nil)
	return sum, nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	buf, err := Ripemd160(in)
	if err != nil {
		return nil, errors.Wrap(err, "Ripemd160")
	}

	return buf[:4], nil
}
