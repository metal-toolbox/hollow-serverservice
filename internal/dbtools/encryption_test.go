package dbtools_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
)

func TestEncryptandDecrypt(t *testing.T) {
	ctx := context.TODO()
	keeper := dbtools.TestSecretKeeper(t)

	secretKey := "NotARealPassword"

	encrypted, err := dbtools.Encrypt(ctx, keeper, secretKey)
	assert.NoError(t, err)
	assert.NotEqual(t, secretKey, encrypted)

	decrypted, err := dbtools.Decrypt(ctx, keeper, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, secretKey, decrypted)
}
