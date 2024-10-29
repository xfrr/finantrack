package assetdomain

import "errors"

const (
	// AssetTypeCash represents the cash asset type.
	AssetTypeCash AssetType = "cash"

	// AssetTypeBank represents the bank asset type.
	AssetTypeBank AssetType = "bank"

	// AssetTypeInvestment represents the investment asset type.
	AssetTypeInvestment AssetType = "investment"

	// AssetTypeOther represents the other asset type.
	AssetTypeOther AssetType = "other"
)

var (
	// ErrInvalidAssetType represents the error when the asset type is invalid.
	ErrInvalidAssetType = errors.New("invalid asset type")
)

// AssetType represents the asset type.
type AssetType string

// String returns the string representation of the asset type.
func (at AssetType) String() string {
	return string(at)
}

// Validate validates the asset type.
func (at AssetType) Validate() error {
	switch at {
	case AssetTypeCash, AssetTypeBank, AssetTypeInvestment, AssetTypeOther:
		return nil
	}

	return ErrInvalidAssetType
}
