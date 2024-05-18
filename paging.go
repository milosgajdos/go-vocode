package vocode

import (
	"fmt"

	"github.com/milosgajdos/go-vocode/request"
)

type Paging struct {
	Page        int  `json:"page"`
	Size        int  `json:"size"`
	Total       int  `json:"total"`
	HasMore     bool `json:"has_more"`
	IsEstimated bool `json:"total_is_estimated"`
}

type Sort struct {
	Col  string
	Desc bool
}

type PageParams struct {
	Page int
	Size int
	Sort *Sort
}

func (l *PageParams) Encode() request.PageParams {
	params := map[string]string{}
	params["page"] = fmt.Sprintf("%d", l.Page)
	params["size"] = fmt.Sprintf("%d", l.Size)
	params["sort_column"] = l.Sort.Col
	params["sort_desc"] = fmt.Sprintf("%v", l.Sort.Desc)
	return params
}
