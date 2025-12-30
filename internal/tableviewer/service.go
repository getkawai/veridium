package tableviewer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// Service provides table viewer operations
type Service struct {
	db *sql.DB
}

// NewService creates a new TableViewer service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetAllTables returns all tables and views in the database
func (s *Service) GetAllTables() ([]TableBasicInfo, error) {
	query := `
		SELECT name, type
		FROM sqlite_master
		WHERE type IN ('table', 'view')
		  AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []TableBasicInfo
	for rows.Next() {
		var table TableBasicInfo
		if err := rows.Scan(&table.Name, &table.Type); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		// Get count for each table
		count, err := s.getTableCount(table.Name)
		if err != nil {
			// Log error but continue
			count = 0
		}
		table.Count = count

		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	return tables, nil
}

// GetTableDetails returns detailed structure information for a table
func (s *Service) GetTableDetails(tableName string) ([]TableColumnInfo, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table details: %w", err)
	}
	defer rows.Close()

	var columns []TableColumnInfo
	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notnull   int
			dfltValue sql.NullString
			pk        int
		)

		if err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		column := TableColumnInfo{
			Name:         name,
			Type:         colType,
			Nullable:     notnull == 0,
			IsPrimaryKey: pk > 0,
		}

		if dfltValue.Valid {
			column.DefaultValue = &dfltValue.String
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}

	return columns, nil
}

// GetTableData returns paginated table data with optional filtering and sorting
func (s *Service) GetTableData(tableName string, pagination PaginationParams, filters []FilterCondition) (*TableDataResult, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	// Build query
	selectClause := fmt.Sprintf("SELECT * FROM %s", tableName)
	whereClause := ""
	orderClause := ""
	var args []interface{}

	// Add filters
	if len(filters) > 0 {
		var whereConditions []string
		for _, filter := range filters {
			switch filter.Operator {
			case "equals":
				whereConditions = append(whereConditions, fmt.Sprintf("%s = ?", filter.Column))
				args = append(args, filter.Value)
			case "contains":
				whereConditions = append(whereConditions, fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", filter.Column))
				args = append(args, fmt.Sprintf("%%%v%%", filter.Value))
			case "startsWith":
				whereConditions = append(whereConditions, fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", filter.Column))
				args = append(args, fmt.Sprintf("%v%%", filter.Value))
			case "endsWith":
				whereConditions = append(whereConditions, fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", filter.Column))
				args = append(args, fmt.Sprintf("%%%v", filter.Value))
			}
		}

		if len(whereConditions) > 0 {
			whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
		}
	}

	// Add sorting
	if pagination.SortBy != nil && *pagination.SortBy != "" {
		direction := "ASC"
		if pagination.SortOrder != nil && strings.ToUpper(*pagination.SortOrder) == "DESC" {
			direction = "DESC"
		}
		orderClause = fmt.Sprintf(" ORDER BY %s %s", *pagination.SortBy, direction)
	}

	// Add pagination
	limitClause := " LIMIT ? OFFSET ?"
	queryArgs := append(args, pagination.PageSize, offset)

	query := selectClause + whereClause + orderClause + limitClause

	// Get data
	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Read all rows
	var data []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{} to hold each column value
		columns := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create map for this row
		rowMap := make(map[string]interface{})
		for i, colName := range columnNames {
			val := columns[i]
			// Convert []byte to string for better JSON serialization
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		data = append(data, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) as total FROM %s%s", tableName, whereClause)
	var total int
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &TableDataResult{
		Data: data,
		Pagination: PaginationResult{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Total:    total,
		},
	}, nil
}

// UpdateRow updates a row in the table
func (s *Service) UpdateRow(tableName string, id string, primaryKeyColumn string, data map[string]interface{}) (map[string]interface{}, error) {
	// Build UPDATE query
	var setParts []string
	var values []interface{}

	for key, value := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, strings.Join(setParts, ", "), primaryKeyColumn)

	_, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to update row: %w", err)
	}

	// Get the updated row
	return s.getRowByPrimaryKey(tableName, primaryKeyColumn, id)
}

// DeleteRow deletes a row from the table
func (s *Service) DeleteRow(tableName string, id string, primaryKeyColumn string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, primaryKeyColumn)
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete row: %w", err)
	}
	return nil
}

// InsertRow inserts a new row into the table
func (s *Service) InsertRow(tableName string, data map[string]interface{}) (map[string]interface{}, error) {
	var columns []string
	var placeholders []string
	var values []interface{}

	for key, value := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	result, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert row: %w", err)
	}

	// Get the last inserted row
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Query the inserted row
	getQuery := fmt.Sprintf("SELECT * FROM %s WHERE rowid = ?", tableName)
	return s.queryRowAsMap(getQuery, lastID)
}

// BatchDelete deletes multiple rows
func (s *Service) BatchDelete(tableName string, ids []string, primaryKeyColumn string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)",
		tableName,
		primaryKeyColumn,
		strings.Join(placeholders, ", "))

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to batch delete: %w", err)
	}
	return nil
}

// getTableCount gets the total row count for a table
func (s *Service) getTableCount(tableName string) (int64, error) {
	// Use parameterized query to prevent SQL injection
	// Note: Table names cannot be parameterized in standard SQL, so we validate the name
	if !isValidTableName(tableName) {
		return 0, fmt.Errorf("invalid table name: %s", tableName)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)
	var count int64
	if err := s.db.QueryRow(query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get table count: %w", err)
	}
	return count, nil
}

// isValidTableName validates table name to prevent SQL injection
func isValidTableName(name string) bool {
	// Allow alphanumeric, underscore, and hyphen
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 64
}

// getRowByPrimaryKey gets a single row by its primary key
func (s *Service) getRowByPrimaryKey(tableName string, primaryKeyColumn string, id string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, primaryKeyColumn)
	return s.queryRowAsMap(query, id)
}

// queryRowAsMap executes a query and returns the first row as a map
func (s *Service) queryRowAsMap(query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no row found")
	}

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	columns := make([]interface{}, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := rows.Scan(columnPointers...); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	rowMap := make(map[string]interface{})
	for i, colName := range columnNames {
		val := columns[i]
		if b, ok := val.([]byte); ok {
			rowMap[colName] = string(b)
		} else {
			rowMap[colName] = val
		}
	}

	return rowMap, nil
}

// ExecuteRawQuery executes a raw SQL query and returns results as JSON
// This is useful for the PgTable UI to execute arbitrary queries
func (s *Service) ExecuteRawQuery(query string, args []interface{}) (string, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columnNames {
			val := columns[i]
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	jsonBytes, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(jsonBytes), nil
}
