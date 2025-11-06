package tableviewer

// TableBasicInfo represents basic information about a database table
type TableBasicInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"` // "BASE TABLE" or "VIEW"
	Count int64  `json:"count"`
}

// TableColumnInfo represents detailed information about a table column
type TableColumnInfo struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Nullable     bool        `json:"nullable"`
	IsPrimaryKey bool        `json:"isPrimaryKey"`
	DefaultValue *string     `json:"defaultValue,omitempty"`
	ForeignKey   *ForeignKey `json:"foreignKey,omitempty"`
}

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// PaginationParams represents pagination and sorting parameters
type PaginationParams struct {
	Page      int     `json:"page"`
	PageSize  int     `json:"pageSize"`
	SortBy    *string `json:"sortBy,omitempty"`
	SortOrder *string `json:"sortOrder,omitempty"` // "asc" or "desc"
}

// FilterCondition represents a filter condition for queries
type FilterCondition struct {
	Column   string      `json:"column"`
	Operator string      `json:"operator"` // "equals", "contains", "startsWith", "endsWith"
	Value    interface{} `json:"value"`
}

// PaginationResult wraps paginated data with metadata
type PaginationResult struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// TableDataResult represents the result of a paginated table query
type TableDataResult struct {
	Data       []map[string]interface{} `json:"data"`
	Pagination PaginationResult         `json:"pagination"`
}
