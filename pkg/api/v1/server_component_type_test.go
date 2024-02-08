package fleetdbapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerComponentTypeSliceSearch(t *testing.T) {
	id0 := "2fc27650-c2d1-013a-6e20-2cde48001122"
	id1 := "404cc820-c2d1-013a-6e21-2cde48001122"

	fixture := ServerComponentTypeSlice{
		&ServerComponentType{
			ID:   id0,
			Name: "Fins",
			Slug: "fins",
		},
		&ServerComponentType{
			ID:   id1,
			Name: "Tails",
			Slug: "tails",
		},
	}

	cases := []struct {
		id       string
		name     string
		slug     string
		expected *ServerComponentType
		testName string
	}{
		// match by ID
		{
			id0,
			"",
			"",
			fixture[0],
			"find by ID",
		},
		// match by Name
		{
			"",
			"Tails",
			"",
			fixture[1],
			"find by Name",
		},
		// match by Slug
		{
			"",
			"",
			"fins",
			fixture[0],
			"find by Slug",
		},
		// no match
		{
			"",
			"Foo",
			"",
			nil,
			"no match",
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			var got *ServerComponentType
			if tc.id != "" {
				got = fixture.ByID(tc.id)
			}

			if tc.name != "" {
				got = fixture.ByName(tc.name)
			}

			if tc.slug != "" {
				got = fixture.BySlug(tc.slug)
			}

			assert.Equal(t, tc.expected, got)
		})
	}
}
