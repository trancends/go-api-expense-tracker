package sharedmodel

type Paging struct {
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPage"`
	TotalRows   int `json:"totalRows"`
	TotalPages  int `json:"totalPages"`
}

// Rumus pagination = limit & offset
// 1 2 3 4 5 6 7 8 9 10
// offset = (page - 1) * limit
// ? = (1 - 1) * 5 ===> 0
// ? = (2 - 1) * 5 ===> 5
// ? = (3 - 1) * 5 ===> 10
// Data ada 10, yang ditampilkan 5, page=1,limit=5 offset 0
// page = 2 limit=5 offset=5
