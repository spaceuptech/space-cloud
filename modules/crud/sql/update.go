package sql

import (
	"fmt"
	"context"
	"strings"
	"strconv"

	goqu "github.com/doug-martin/goqu/v8"

	_ "github.com/go-sql-driver/mysql"                 // Import for MySQL
	_ "github.com/lib/pq"                              // Import for postgres
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"  // Dialect for postgres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Update updates the document(s) which match the condition provided.
func (s *SQL) Update(ctx context.Context, project, col string, req *model.UpdateRequest) error {
	return s.update(ctx, project, col, req, s.client)
}

func (s *SQL) update(ctx context.Context, project, col string, req *model.UpdateRequest, executor interface{}) error {
	if req == nil {
		return utils.ErrInvalidParams
	}
	if req.Update == nil {
		return utils.ErrInvalidParams
	}
	switch req.Operation {
	case utils.All:
		for k := range req.Update {
			switch k {
			case "$set", "$inc", "$mul", "$max", "$min", "$currentDate":
				sqlQuery, args, err := s.generateUpdateQuery(ctx, project, col, req, k)
				if err != nil {
					return err
				}
				err = doExecContext(ctx, sqlQuery, args, executor)
				if err != nil {
					return err
				}
			default: // (case "$push", "$unset", "$rename")
				return utils.ErrInvalidParams
			}
		}
	case utils.Upsert:
		for k := range req.Update {
			switch k {
			case "$set", "$inc", "$mul", "$max", "$min", "$currentDate":
				err := s.upsert(ctx, project, col, req, k)
				if err != nil {
					return err
				}
			default: // (case "$push", "$unset", "$rename")
				return utils.ErrInvalidParams
			}
		}
	default: // (case utils.One)
		return utils.ErrInvalidParams
	}
	return nil
}

//generateUpdateQuery makes query for update operations
func (s *SQL) generateUpdateQuery(ctx context.Context, project, col string, req *model.UpdateRequest, op string) (string, []interface{}, error) {
	// Generate a prepared query builder
	dialect := goqu.Dialect(s.dbType)
	query := dialect.From(col)
	if op == "$set" {
		query = query.Prepared(true)
	}

	if req.Find != nil {
		// Get the where clause from query object
		var err error
		query, err = generateWhereClause(query, req.Find)
		if err != nil {
			return "", nil, err
		}
	}

	if req.Update == nil {
		return "", nil, utils.ErrInvalidParams
	}
	m, ok := req.Update[op].(map[string]interface{})
	if !ok {
		return "", nil, utils.ErrInvalidParams
	}

	if op == "$currentDate" {
		for k,v := range m {
			mm, ok := v.(map[string]interface{})
			if !ok {
				return "", nil, utils.ErrInvalidParams
			}
			for _, val := range mm {
				val, ok := val.(string)
				if !ok {
					return "", nil, utils.ErrInvalidParams
				}
				switch val {
				case "date", "timestamp":
					m[k] = val
				default:
					return "", nil, utils.ErrInvalidParams
				}
			}
		}
	}
	
	record, err := generateRecord(req.Update[op])
	if err != nil {
		return "", nil, err
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Update().Set(record).ToSQL()
	if err != nil {
		return "", nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	switch op {
	case "$set":
	case "$inc":
		for k, v := range m {
			val, err := numToString(v)
			if err != nil {
				return "", nil, err
			}
			sqlString = strings.Replace(sqlString, k+"="+val, k+"="+k+"+"+val, -1)
		}
	case "$mul":
		for k, v := range m {
			val, err := numToString(v)
			if err != nil {
				return "", nil, err
			}
			sqlString = strings.Replace(sqlString, k+"="+val, k+"="+k+"*"+val, -1)
		}
	case "$max":
		for k, v := range m {
			val, err := numToString(v)
			if err != nil {
				return "", nil, err
			}
			sqlString = strings.Replace(sqlString, k+"="+val, k+"=GREATEST("+k+","+val+")", -1)
		}
	case "$min":
		for k, v := range m {
			val, err := numToString(v)
			if err != nil {
				return "", nil, err
			}
			sqlString = strings.Replace(sqlString, k+"="+val, k+"=LEAST("+k+","+val+")", -1)
		}
	case "$currentDate":
		for k, v := range m {
			val, ok := v.(string)
			if !ok {
				return "", nil, utils.ErrInvalidParams
			}
			if val == "timestamp" {
				sqlString = strings.Replace(sqlString, k+"='"+val+"'", k+"=CURRENT_TIMESTAMP()", -1)
			} else {
				sqlString = strings.Replace(sqlString, k+"='"+val+"'", k+"=CURRENT_DATE()", -1)
			}
		}
	default:
		return "", nil, utils.ErrInvalidParams
	}
	return sqlString, args, nil
}

// helper function to upsert data
func (s *SQL) upsert(ctx context.Context, project, col string, req *model.UpdateRequest, op string) error {
	tx, err := s.client.BeginTxx(ctx, nil) //TODO - Write *sqlx.TxOption instead of nil
	if err != nil {
		fmt.Println("Error in initiating Transaction for Upsert")
		return err
	}

	sqlString, args, err := s.generateUpdateQuery(ctx, project, col, req, op)
	if err != nil {
		return err
	}
	
	stmt, err := tx.PreparexContext(ctx, sqlString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		if op == "$currentDate" {
			return utils.ErrInvalidParams
		}
		sqlQuery, args, err := s.generateCreateQuery(ctx, project, col, &model.CreateRequest{Document: req.Update[op], Operation: utils.One})
		if err != nil {
			return err
		}
		err = doExecContext(ctx, sqlQuery, args, tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func numToString(v interface{}) (string, error) {
	val, ok := v.(float64)
	if !ok {
		return "", utils.ErrInvalidParams
	}
	return strconv.FormatFloat(val, 'f', -1, 64), nil
}
