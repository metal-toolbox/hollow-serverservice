package hollow

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
)

// ServerResponse represents the data that the server will return on any given call
type ServerResponse struct {
	PageSize         int                 `json:"page_size,omitempty"`
	Page             int                 `json:"page,omitempty"`
	PageCount        int                 `json:"page_count,omitempty"`
	TotalPages       int                 `json:"total_pages,omitempty"`
	TotalRecordCount int64               `json:"total_record_count,omitempty"`
	NextCursor       string              `json:"next_cursor,omitempty"`
	Links            ServerResponseLinks `json:"_links,omitempty"`
	Message          string              `json:"message,omitempty"`
	Error            string              `json:"error,omitempty"`
	Slug             string              `json:"slug,omitempty"`
	Record           interface{}         `json:"record,omitempty"`
	Records          interface{}         `json:"records,omitempty"`
}

// ServerResponseLinks represent links that could be returned on a page
type ServerResponseLinks struct {
	Self       *Link `json:"self,omitempty"`
	NextCursor *Link `json:"next_cursor,omitempty"`
	First      *Link `json:"first,omitempty"`
	Previous   *Link `json:"previous,omitempty"`
	Next       *Link `json:"next,omitempty"`
	Last       *Link `json:"last,omitempty"`
}

// Link represents an address to a page
type Link struct {
	Href string `json:"href,omitempty"`
}

func newErrorResponse(m string, err error) *ServerResponse {
	return &ServerResponse{
		Message: m,
		Error:   err.Error(),
	}
}

func badRequestResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, newErrorResponse(message, err))
}

func notFoundResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, newErrorResponse("resource not found", err))
}

func createdResponse(c *gin.Context, slug string) {
	r := &ServerResponse{
		Message: "resource created",
		Slug:    slug,
		Links: ServerResponseLinks{
			Self: &Link{Href: fmt.Sprintf("%s/%s", uriWithoutQueryParams(c), slug)},
		},
	}

	c.JSON(http.StatusOK, r)
}

func deletedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, &ServerResponse{Message: "resource deleted"})
}

func updatedResponse(c *gin.Context, slug string) {
	r := &ServerResponse{
		Message: "resource updated",
		Slug:    slug,
		Links: ServerResponseLinks{
			Self: &Link{Href: uriWithoutQueryParams(c)},
		},
	}

	c.JSON(http.StatusOK, r)
}

func dbFailureResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("datastore error", err))
}

func dbErrorResponse(c *gin.Context, err error) {
	if errors.Is(err, db.ErrNotFound) {
		notFoundResponse(c, err)
	} else {
		dbFailureResponse(c, err)
	}
}

func failedConvertingToVersioned(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("failed parsing the datastore results", err))
}

func listResponse(c *gin.Context, i interface{}, p paginationData) {
	uri := c.Request.URL

	r := &ServerResponse{
		PageSize:  p.pager.LimitUsed(),
		PageCount: p.pageCount,
		Records:   i,
		Links: ServerResponseLinks{
			Self: &Link{Href: uri.String()},
		},
	}

	// Only include total counts if we are not using a cursor, otherwise counts will be wrong
	if p.pager.Cursor == nil {
		d := float64(p.totalCount) / float64(p.pager.LimitUsed())
		r.TotalPages = int(math.Ceil(d))
		r.Page = p.pager.Page
		r.TotalRecordCount = p.totalCount

		r.Links.First = &Link{Href: getURIWithQuerySet(*uri, "page", "1")}
		r.Links.Last = &Link{Href: getURIWithQuerySet(*uri, "page", strconv.Itoa(r.TotalPages))}

		if r.Page < r.TotalPages {
			r.Links.Next = &Link{Href: getURIWithQuerySet(*uri, "page", strconv.Itoa(r.Page+1))}
		}

		if r.Page != 1 {
			r.Links.Previous = &Link{Href: getURIWithQuerySet(*uri, "page", strconv.Itoa(r.Page-1))}
		}
	}

	if p.nextCursor != "" && p.pageCount == p.pager.LimitUsed() {
		r.NextCursor = p.nextCursor
		r.Links.NextCursor = &Link{Href: getURIWithQuerySet(*uri, "cursor", p.nextCursor)}
	}

	c.JSON(http.StatusOK, r)
}

func itemResponse(c *gin.Context, i interface{}) {
	r := &ServerResponse{
		Record: i,
		Links: ServerResponseLinks{
			Self: &Link{Href: c.Request.URL.String()},
		},
	}
	c.JSON(http.StatusOK, r)
}

func getURIWithQuerySet(uri url.URL, key, value string) string {
	q := uri.Query()
	q.Del(key)
	q.Add(key, value)
	uri.RawQuery = q.Encode()

	return uri.String()
}

func uriWithoutQueryParams(c *gin.Context) string {
	uri := c.Request.URL
	uri.RawQuery = ""

	return uri.String()
}
