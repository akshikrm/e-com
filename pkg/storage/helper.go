package storage

import (
	"fmt"
	"net/url"
)

func buildFilterQuery(query string, filter url.Values) string {
	startDate := filter.Get("start_date")
	endDate := filter.Get("end_date")
	categoryId := filter.Get("category_id")

	if categoryId != "" {
		query = fmt.Sprintf("%s AND p.category_id = '%s'", query, categoryId)
	}
	if endDate != "" {
		query = fmt.Sprintf("%s AND p.created_at < '%s'", query, endDate)
	}
	if startDate != "" {
		if endDate != "" {
			query = fmt.Sprintf("%s AND p.created_at > '%s'", query, startDate)
		} else {
			query = fmt.Sprintf("%s AND p.created_at < '%s'", query, startDate)
		}
	}
	return query
}
