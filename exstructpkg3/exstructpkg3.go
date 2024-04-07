package exstructpkg3

type SingleGenericStructPkg3[T any, K any] struct {
	Data           []T          `json:"data"`
	Data2          map[string]K `json:"data2"`
	PaginationInfo string       `json:"pagination_info"`
}
