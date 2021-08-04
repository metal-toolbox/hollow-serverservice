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

// HasNextPage will return if there are additional resources to load on additional pages
func (r *ServerResponse) HasNextPage() bool {
	return r.Records != nil && (r.Links.NextCursor != nil || r.Links.Next != nil)
}

func badRequestResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, &ServerResponse{Message: message, Error: err.Error()})
}

func createdResponse(c *gin.Context, slug string) {
	uri := fmt.Sprintf("%s/%s", uriWithoutQueryParams(c), slug)
	r := &ServerResponse{
		Message: "resource created",
		Slug:    slug,
		Links: ServerResponseLinks{
			Self: &Link{Href: uri},
		},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
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

func dbErrorResponse(c *gin.Context, err error) {
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, &ServerResponse{Message: "resource not found", Error: err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "datastore error", Error: err.Error()})
	}
}

func failedConvertingToVersioned(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "failed parsing the datastore results", Error: err.Error()})
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
		r.Links.NextCursor = &Link{Href: getURIWithCursor(*uri, p.nextCursor)}
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

func getURIWithCursor(uri url.URL, value string) string {
	q := uri.Query()
	q.Del("cursor")
	q.Del("page")
	q.Add("cursor", value)
	uri.RawQuery = q.Encode()

	return uri.String()
}

func uriWithoutQueryParams(c *gin.Context) string {
	uri := c.Request.URL
	uri.RawQuery = ""

	return uri.String()
}
