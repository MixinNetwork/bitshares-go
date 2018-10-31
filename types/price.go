package types

import (
	"encoding/json"

	"github.com/MixinNetwork/bitshares-go/encoding/transaction"
)

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

type AssetAmount struct {
	Amount  string   `json:"amount"`
	AssetID ObjectID `json:"asset_id"`
}

func (aa AssetAmount) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(aa.Amount)
	enc.Encode(aa.AssetID)
	return enc.Err()
}

// RPC client might return asset amount as uint64 or string,
// therefore a custom unmarshaller is used
func (aa *AssetAmount) UnmarshalJSON(b []byte) (err error) {
	var body AssetAmount
	err = json.Unmarshal(b, &body)
	if err == nil {
		aa.AssetID = body.AssetID
		aa.Amount = body.Amount
	}

	return err
}
