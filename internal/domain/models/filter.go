package models

type Filter struct {
	SearchQuery string
	Limit       int64
	Page        int64
	OrderBy     string
	Desc        bool
}
