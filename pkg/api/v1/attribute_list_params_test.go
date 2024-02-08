package fleetdbapi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAttributesListParams(t *testing.T) {
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
			"query with namespace",
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
		{
			[]AttributeListParams{
				{
					Namespace:         "hollow.versioned",
					Keys:              []string{"a", "b"},
					Operator:          "lt",
					Value:             "5",
					AttributeOperator: AttributeLogicalOR,
				},
			},
			"key",
			make(url.Values),
			url.Values{
				"key": []string{
					"hollow.versioned~a.b~lt~5~or",
				},
			},
			"query with namespace, keys and operator, value, OR Attribute Operator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			encodeAttributesListParams(tc.alp, tc.key, tc.queryValues)
			assert.Equal(t, tc.expected, tc.queryValues)
		})
	}
}

func TestParseQueryAttributesListParams(t *testing.T) {
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
		{
			"attr",
			"attr=hollow.versioned~a.name~like~nemo&attr=hollow.versioned~a.name~like~bluefin~or",
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
				{
					Namespace: "hollow.versioned",
					Keys: []string{
						"a",
						"name",
					},
					Operator:          "like",
					Value:             "bluefin%",
					AttributeOperator: AttributeLogicalOR,
				},
			},
			"query with Attribute Operator - OR",
		},
		{
			"attr",
			"attr=hollow.versioned~a.name~like~nemo&attr=hollow.versioned~a.name~like~bluefin~wtfbbq",
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
				{
					Namespace: "hollow.versioned",
					Keys: []string{
						"a",
						"name",
					},
					Operator: "like",
					Value:    "bluefin%",
				},
			},
			"query with invalid attribute operator defaults to Attribute operator - AND",
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

func TestSetJSONBWhereClause(t *testing.T) {
	tblName := "foo"

	testCases := []struct {
		jsonPath         string
		param            AttributeListParams
		expectedValues   []interface{}
		expectedWhereStr string
		testName         string
	}{
		{
			"?",
			AttributeListParams{
				Namespace: "hollow.versioned",
				Keys: []string{
					"a",
					"b",
				},
				Operator: "lt",
				Value:    "5",
			},
			[]interface{}{"5"},
			"json_extract_path_text(foo.data::JSONB, ?)::int < ?",
			"where less than",
		},
		{
			"?",
			AttributeListParams{
				Namespace: "hollow.versioned",
				Keys: []string{
					"a",
					"b",
				},
				Operator: "gt",
				Value:    "5",
			},
			[]interface{}{"5"},
			"json_extract_path_text(foo.data::JSONB, ?)::int > ?",
			"where greater than",
		},
		{
			"?",
			AttributeListParams{
				Namespace: "hollow.versioned",
				Keys: []string{
					"a",
					"b",
				},
				Operator: "like",
				Value:    "foobar",
			},
			[]interface{}{"foobar"},
			"json_extract_path_text(foo.data::JSONB, ?) LIKE ?",
			"like",
		},
		{
			"?",
			AttributeListParams{
				Namespace: "hollow.versioned",
				Keys: []string{
					"a",
					"b",
				},
				Operator: "eq",
				Value:    "10",
			},
			[]interface{}{"10"},
			"json_extract_path_text(foo.data::JSONB, ?) = ?",
			"equal",
		},
		{
			"",
			AttributeListParams{
				Namespace: "hollow.versioned",
				Keys: []string{
					"a",
					"b",
				},
			},
			[]interface{}{},
			`foo.data::JSONB -> ? \? ?`,
			"default - keys specified with no operator",
		},
		{
			"?",
			AttributeListParams{
				Namespace: "hollow.versioned",
			},
			[]interface{}{},
			"foo.data::JSONB",
			"default - no keys",
		},
	}

	values := []interface{}{}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gotWhereStr, gotValues := tc.param.setJSONBWhereClause(tblName, tc.jsonPath, values)
			assert.Equal(t, tc.expectedValues, gotValues, "values")
			assert.Equal(t, tc.expectedWhereStr, gotWhereStr, "where stmt")
		})
	}
}
