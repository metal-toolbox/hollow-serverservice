package dbtools

import (
	"context"
	"encoding/base64"

	"gocloud.dev/secrets"
)

// Encrypt provides a wrapper to handle encrypting a string with the secrets keeper
// and returns it already base64 encoded
func Encrypt(ctx context.Context, keeper *secrets.Keeper, str string) (string, error) {
	cipher, err := keeper.Encrypt(ctx, []byte(str))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipher), nil
}

// Decrypt provides a wrapper to handle decrypting a base64 encoded string with
// the secrets keeper
func Decrypt(ctx context.Context, keeper *secrets.Keeper, base64str string) (string, error) {
	plain, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return "", err
	}

	decrypted, err := keeper.Decrypt(ctx, plain)

	return string(decrypted), err
}
