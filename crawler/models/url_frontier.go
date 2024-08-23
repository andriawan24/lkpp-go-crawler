package models

const (
	URL_STATUS_READY   = 0
	URL_STATUS_CRAWLED = 1
	URL_STATUS_ERROR   = 2
)

type UrlFrontier struct {
	ID        int
	Url       string
	Crawler   string
	Status    int64
	CreatedAt int64
	UpdatedAt int64
}
