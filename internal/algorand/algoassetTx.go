package algorand

import (
	"fmt"

	"exitor-dapp/internal/createasset"
)



// Makes a transaction on Algorand, and then reads it into the database at the same time
func (repo *Repository) CreateOnAlgorand(ctx context.Context, claims auth.Claims, req CreatedAssetCreate, now time.Time) (* CreatedAsset, error) {
	mAlgorand := CreatedAsset{
		ID:        uuid.NewRandom().String(),
		AccountID: req.AccountID,
		WalletAddress: req.WalletAddress,
		Total:		req.Total,
		AssetName:      req.AssetName,
		Decimals: 	req.Decimals,
		DefaultFrozen: req.DefaultFrozen,
		URL:			req.URL,
		Status:    CreatedAssetStatus_Active,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// We need to derive the mnemonic of the
	// account in order to sign a transaction


	assetTotalIssuance := uint64(mAlgorand.Total)
	assetDecimalsForDisplay := uint32(mAlgorand.Decimals)
	accountsAreDefaultFrozen := mAlgorand.DefaultFrozen
	managerAddress := mAlgorand.WalletAddress
	assetReserveAddress := ""
	addressWithFreezingPrivileges := mAlgorand.WalletAddress
	addressWithClawbackPrivileges := mAlgorand.WalletAddress
	assetName := mAlgorand.AssetName
	assetUrl := mAlgorand.URL
	assetMetadataHash := ""

	tx, err := transaction.MakeAssetCreateTxn(mAlgorand.WalletAddress, assetTotalIssuance, assetDecimalsForDisplay,
			accountsAreDefaultFrozen, managerAddress, assetReserveAddress, addressWithFreezingPrivileges, addressWithClawbackPrivileges,
			assetName, assetUrl, assetMetadataHash)
	if err != nil {
			fmt.Printf("Error creating transaction: %\n", err)
			return
	}

	// Let's sign the transaction
	_, bytes, err := crypto.SignTransaction()
}
