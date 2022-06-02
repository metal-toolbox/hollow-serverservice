package serverservice

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_encodeAttributesListParams(t *testing.T) {
	testCases := []struct {
		alp         []AttributeListParams
		key         string
		queryValues url.Values
		expected    url.Values
		testName    string
	}{
		// TODO: add test cases for non-happy paths
		{
			[]AttributeListParams{
				{
					Namespace: "hollow.versioned",
				},
			},
			"key",
			make(url.Values),
			url.Values{
				"key": []string{
					"hollow.versioned",
				},
			},
			"query with",
		},
		{
			[]AttributeListParams{
				{
					Namespace: "hollow.versioned",
					Keys:      []string{"a", "b"},
				},
			},
			"key",
			make(url.Values),
			url.Values{
				"key": []string{
					"hollow.versioned~a.b",
				},
			},
			"query with namespace, keys",
		},
		{
			[]AttributeListParams{
				{
					Namespace: "hollow.versioned",
					Keys:      []string{"a", "b"},
					Operator:  "lt",
					Value:     "5",
				},
			},
			"key",
			make(url.Values),
			url.Values{
				"key": []string{
					"hollow.versioned~a.b~lt~5",
				},
			},
			"query with namespace, keys and operator, value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			encodeAttributesListParams(tc.alp, tc.key, tc.queryValues)
			assert.Equal(t, tc.expected, tc.queryValues)
		})
	}
}

func Test_parseQueryAttributesListParams(t *testing.T) {
	testCases := []struct {
		key      string
		query    string
		expected []AttributeListParams
		testName string
	}{
		{
			// TODO: add test cases for non-happy paths
			"attr",
			"attr=hollow.versioned~a.b~lt~5",
			[]AttributeListParams{
				{
					Namespace: "hollow.versioned",
					Keys: []string{
						"a",
						"b",
					},
					Operator: "lt",
					Value:    "5",
				},
			},
			"query with namespace, attribute list params",
		},
		{
			"attr",
			"attr=hollow.versioned~a.name~like~nemo",
			[]AttributeListParams{
				{
					Namespace: "hollow.versioned",
					Keys: []string{
						"a",
						"name",
					},
					Operator: "like",
					Value:    "nemo%",
				},
			},
			"query with namespace, attribute list params and 'like' operator",
		},
	}

	setupGinCtx := func(queryURL string) *gin.Context {
		ctx, _ := gin.CreateTestContext(nil)
		// release mode set for less noise in tests
		gin.SetMode(gin.ReleaseMode)

		ctx.Request = httptest.NewRequest(
			http.MethodGet,
			"https://hollow.sh/servers?"+queryURL,
			nil,
		)

		return ctx
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ginCtx := setupGinCtx(tc.query)
			got := parseQueryAttributesListParams(ginCtx, tc.key)
			assert.Equal(t, tc.expected, got)
		})
	}
}
