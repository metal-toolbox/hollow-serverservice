package serverservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_parseQueryServerComponentsListParams(t *testing.T) {
	testCases := []struct {
		query    string
		expected []ServerComponentListParams
		testName string
	}{
		// TODO: add test cases for non-happy paths
		{
			"",
			[]ServerComponentListParams{},
			"zero map query params returns empty list",
		},
		{
			"sc_0[model]=blue_tang",
			[]ServerComponentListParams{
				{
					Model: "blue_tang",
				},
			},
			"map query with a single key value",
		},
		{
			"sc_0[vendor]=fish&sc_0[model]=blue_tang&sc_0[name]=dory%20fish&sc_0[serial]=1234&sc_0[type]=fins",
			[]ServerComponentListParams{
				{
					Name:                "dory fish",
					Vendor:              "fish",
					Model:               "blue_tang",
					Serial:              "1234",
					ServerComponentType: "fins",
				},
			},
			"map query with multiple key values",
		},
		{
			"sc_0[model]=blue_tang&sc_0[name]=dory&sc_1[model]=clownfish&sc_1[name]=marlin",
			[]ServerComponentListParams{
				{
					Model: "blue_tang",
					Name:  "dory",
				},
				{
					Model: "clownfish",
					Name:  "marlin",
				},
			},
			"multiple map queries with multiple key values",
		},
		{
			"sc_0[model]=blue_tang&sc_0_attr=name.space~foo~eq~bar",
			[]ServerComponentListParams{
				{
					Model: "blue_tang",
					AttributeListParams: []AttributeListParams{
						{
							Namespace: "name.space",
							Keys: []string{
								"foo",
							},
							Operator: "eq",
							Value:    "bar",
						},
					},
				},
			},
			"map query with attribute query",
		},
		{
			"sc_0[model]=blue_tang&sc_0_ver_attr=name.space~foo~eq~bar",
			[]ServerComponentListParams{
				{
					Model: "blue_tang",
					VersionedAttributeListParams: []AttributeListParams{
						{
							Namespace: "name.space",
							Keys: []string{
								"foo",
							},
							Operator: "eq",
							Value:    "bar",
						},
					},
				},
			},
			"map query with versioned attribute query",
		},
		{
			"sc_0_attr=name.space~foo~eq~bar",
			[]ServerComponentListParams{
				{
					AttributeListParams: []AttributeListParams{
						{
							Namespace: "name.space",
							Keys: []string{
								"foo",
							},
							Operator: "eq",
							Value:    "bar",
						},
					},
				},
			},
			"versioned attribute attribute query",
		},
	}

	setupGinCtx := func(queryParam string) *gin.Context {
		ctx, _ := gin.CreateTestContext(nil)
		// release mode set for less noise in tests
		gin.SetMode(gin.ReleaseMode)

		ctx.Request = httptest.NewRequest(
			http.MethodGet,
			"https://hollow.sh/components?"+queryParam,
			nil,
		)

		return ctx
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ginCtx := setupGinCtx(tc.query)
			got, _ := parseQueryServerComponentsListParams(ginCtx)
			assert.Equal(t, tc.expected, got)
		})
	}
}
