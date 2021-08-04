package hollow

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
)

// PaginationParams allow you to paginate the results
type PaginationParams struct {
	Limit  int    `json:"limit,omitempty"`
	Page   int    `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type paginationData struct {
	pageCount  int
	totalCount int64
	nextCursor string
	pager      db.Pagination
}

func encodeCursor(t time.Time) string {
	key := t.Format(time.RFC3339Nano)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func decodeCursor(encodedCursor string) (res *time.Time, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	t, err := time.Parse(time.RFC3339Nano, string(byt))
	if err != nil {
		return
	}

	res = &t

	return
}

func parsePagination(c *gin.Context) (db.Pagination, error) {
	var err error
	// Initializing default
	limit := db.DefaultPaginationSize
	page := 1
	query := c.Request.URL.Query()

	var cursor *time.Time

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "cursor":
			cursor, err = decodeCursor(queryValue)
			if err != nil {
				return db.Pagination{}, err
			}
		}
	}

	return db.Pagination{
		Limit:  limit,
		Page:   page,
		Cursor: cursor,
	}, nil
}

func (p *PaginationParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}

	if p.Page != 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}

	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
}
