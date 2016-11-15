package omniboost

import (
	"strconv"

	"github.com/tim-online/go-omniboost/utils"
)

type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

func newMeta() *Meta {
	return &Meta{
		Pagination: newPagination(),
	}
}

type Pagination struct {
	Total       int    `json:"total,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	CurrentPage int    `json:"current_page,omitempty"`
	Count       int    `json:"count,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	Links       *Links `json:"links,omitempty"`
}

func newPagination() *Pagination {
	return &Pagination{
		Total:       0,
		PerPage:     0,
		CurrentPage: 1,
		Count:       0,
		TotalPages:  1,
		Links:       newLinks(),
	}
}

// IsLastPage returns true if the current page is the last
func (p *Pagination) IsLastPage() bool {
	return p.CurrentPage >= p.TotalPages
}

// Links manages links that are returned along with a List
type Links struct {
	Current *utils.URL `json:"current,omitempty"`
	First   *utils.URL `json:"first,omitempty"`
	Last    *utils.URL `json:"last,omitempty"`
	Next    *utils.URL `json:"next,omitempty"`
	Prev    *utils.URL `json:"prev,omitempty"`
}

func newLinks() *Links {
	return &Links{
		Current: nil,
		First:   nil,
		Last:    nil,
		Next:    nil,
		Prev:    nil,
	}
}

func pageForURL(urlText string) (int, error) {
	u, err := utils.ParseRequestURI(urlText)
	if err != nil {
		return 0, err
	}

	pageStr := u.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}

	return page, nil
}
