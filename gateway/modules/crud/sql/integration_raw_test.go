// +build integration

package sql

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_CreateDatabaseIfNotExist(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		query   string
		args    args
		wantErr bool
		want    []interface{}
	}{
		{
			name:  "Db Creation check",
			query: "SELECT schema_name FROM information_schema.schemata where SCHEMA_NAME = 'myproject';",
			args: args{
				ctx:  context.Background(),
				name: *dbType,
			},
			wantErr: false,
			want:    []interface{}{map[string]interface{}{"SCHEMA_NAME": "myproject"}},
		},
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("CreateDatabaseIfNotExist() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.CreateDatabaseIfNotExist(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabaseIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			rows, err := db.client.Queryx(tt.query)
			if err != nil {
				t.Error("CreateDatabaseIfNotExist() query error", err)
				return
			}
			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("CreateDatabaseIfNotExist() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if !reflect.DeepEqual(tt.want, readResult) {
				t.Errorf("CreateDatabaseIfNotExist() got (%v) want (%v)", readResult, tt.want)
			}
		})
	}
}

func TestSQL_GetConnectionState(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Connection State",
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("GetConnectionState() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetConnectionState(tt.args.ctx); got != tt.want {
				t.Errorf("GetConnectionState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQL_RawBatch(t *testing.T) {
	type args struct {
		ctx     context.Context
		queries []string
	}
	tests := []struct {
		name    string
		query   string
		args    args
		wantErr bool
		want    []map[string]interface{}
	}{
		{
			name:  "Raw Batch List of Queries",
			query: "SELECT * FROM raw_batch",
			args: args{
				ctx: context.Background(),
				queries: []string{
					"INSERT INTO raw_batch (id,score) VALUE ('11',20)",
					"INSERT INTO raw_batch (id,score) VALUE ('22',30)",
				},
			},
			want:    []map[string]interface{}{{"id": "11", "score": int64(20)}, {"id": "22", "score": int64(30)}},
			wantErr: false,
		},
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("RawBatch() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.RawBatch(tt.args.ctx, tt.args.queries); (err != nil) != tt.wantErr {
				t.Errorf("RawBatch() error = %v, wantErr %v", err, tt.wantErr)
			}

			rows, err := db.client.Queryx(tt.query)
			if err != nil {
				t.Error("RawBatch() query error", err)
				return
			}
			readResult := []map[string]interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("RawBatch() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}

			if !reflect.DeepEqual(tt.want, readResult) {
				t.Errorf("RawBatch() got (%v) want (%v)", readResult, tt.want)
			}
		})
	}
	if _, err := db.client.Exec("TRUNCATE TABLE raw_batch"); err != nil {
		t.Log("RawBatch() Couldn't truncate table", err)
	}
}

func TestSQL_RawQuery(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name       string
		query      string
		args       args
		want       int64
		want1      interface{}
		wantErr    bool
		wantResult []interface{}
	}{
		{
			name:  "Raw Prepared Query",
			query: "SELECT * FROM raw_query",
			args: args{
				ctx:   context.Background(),
				query: "INSERT INTO raw_query (id,score) VALUE (?,?)",
				args:  []interface{}{"1", 20},
			},
			want:       0,
			want1:      []interface{}{},
			wantErr:    false,
			wantResult: []interface{}{map[string]interface{}{"id": "1", "score": int64(20)}},
		},
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("RawQuery() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if utils.DBType(*dbType) == utils.SQLServer {
				tt.args.query = db.generateQuerySQLServer(tt.args.query)
			}
			got, got1, err := db.RawQuery(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RawQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RawQuery() got1 = %v, want %v", got1, tt.want1)
			}

			rows, err := db.client.Queryx(tt.query)
			if err != nil {
				t.Error("RawQuery() query error", err)
				return
			}
			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("RawQuery() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if !reflect.DeepEqual(tt.wantResult, readResult) {
				t.Errorf("RawQuery() got (%v) want (%v)", readResult, tt.want)
			}
		})
	}
	if _, err := db.client.Exec("TRUNCATE TABLE raw_query"); err != nil {
		t.Log("RawQuery() Couldn't truncate table", err)
	}
}
