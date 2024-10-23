package model

type HttpResult struct {
	Msg        string
	StatusCode int
	Object     interface{}
}
type CountResponse struct {
	Total      int64
	Validated  int64
	TotalHttp  int64
	ValidHttp  int64
	TotalSocks int64
	ValidSocks int64
	TotalVmess int64
	ValidVmess int64
}
