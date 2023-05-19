//nolint:wsl
package serverservice

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"

	"go.hollow.sh/serverservice/internal/models"
)

func TestSerialization(t *testing.T) {
	srv := &models.Server{
		Name:         null.StringFrom("server-name"),
		FacilityCode: null.StringFrom("fc13"),
		ID:           "some-uuid-str",
	}

	_, err := NewCreateServerMessage((*models.Server)(nil))
	require.ErrorIs(t, err, ErrNilServer, "nil input")

	byt, err := NewCreateServerMessage(srv)
	require.NoError(t, err, "good server obj")
	require.NotNil(t, byt, "good server obj")

	bogus := []byte("bogus")
	_, err = DeserializeCreateServer(bogus)
	require.ErrorIs(t, err, ErrBadJSONIn, "bogus deserialize")

	exp := &CreateServer{
		Name:         null.StringFrom("server-name"),
		FacilityCode: null.StringFrom("fc13"),
		ID:           "some-uuid-str",
	}

	cs, err := DeserializeCreateServer(byt)
	require.NoError(t, err, "good deserialize")
	require.Equal(t, exp.Name, cs.Name, "good deserialize name")
	require.Equal(t, exp.FacilityCode, cs.FacilityCode, "good deserialize facility")
	require.Equal(t, exp.ID, cs.ID, "good deserialize id")
}
