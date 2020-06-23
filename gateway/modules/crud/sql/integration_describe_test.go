// +build integration

package sql

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_DescribeTable(t *testing.T) {
	var firstColumn = "column1"
	// var secondColumn = "column2"
	type args struct {
		ctx context.Context
		col string
	}
	tests := []struct {
		name        string
		createQuery []string
		scQueries   []string
		args        args
		fields      []utils.FieldType
		foreignKeys []utils.ForeignKeysType
		indexKeys   []utils.IndexType
		wantErr     bool
	}{
		{
			name:        "MySQL field col1 with type ID",
			createQuery: []string{"CREATE TABLE table1 (column1 varchar(50))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []utils.FieldType{{FieldName: "column1", FieldType: "varchar(50)", FieldNull: "YES"}},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type String",
			createQuery: []string{"CREATE TABLE table1 (column1 text)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []utils.FieldType{{FieldName: "column1", FieldType: "text", FieldNull: "YES"}},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean",
			createQuery: []string{"CREATE TABLE table1 (column1 boolean)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "tinyint", FieldNull: "YES"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer",
			createQuery: []string{"CREATE TABLE table1 (column1 bigint)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "YES"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float",
			createQuery: []string{"CREATE TABLE table1 (column1 float)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "YES"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type JSON",
			createQuery: []string{"CREATE TABLE table1 (column1 json)"},

			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "YES"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type DateTime",
			createQuery: []string{"CREATE TABLE table1 (column1 datetime)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "YES"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type ID ",
			createQuery: []string{"CREATE TABLE table1 (column1 varchar(50) NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type String ",
			createQuery: []string{"CREATE TABLE table1 (column1 text NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Boolean ",
			createQuery: []string{"CREATE TABLE table1 (column1 boolean NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "tinyint", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Integer ",
			createQuery: []string{"CREATE TABLE table1 (column1 bigint NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Float ",
			createQuery: []string{"CREATE TABLE table1 (column1 float NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type DateTime ",
			createQuery: []string{"CREATE TABLE table1 (column1 datetime NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type JSON ",
			createQuery: []string{"CREATE TABLE table1 (column1 json NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []utils.IndexType{},
			foreignKeys: []utils.ForeignKeysType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}},
			wantErr:     false,
		},
		// There is a bug in code, inspection cannot detect @createdAt,@updatedAt directives
		// TODO: What other special directives do we have ?
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @createdAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @updatedAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// NOTE: JSON & text type cannot have default value
		{
			name:        "MySQL field col1 which is not null with type ID having default value INDIA",
			createQuery: []string{"CREATE TABLE table1 (column1 varchar(50) NOT NULL DEFAULT 'INDIA')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldDefault: "INDIA", FieldExtra: "INDIA"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE table1 (column1 boolean NOT NULL DEFAULT true)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "tinyint", FieldNull: "NO", FieldDefault: "1", FieldExtra: "1"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Integer having default value 100",
			createQuery: []string{"CREATE TABLE table1 (column1 bigint NOT NULL DEFAULT 100)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", FieldExtra: "100"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Float having default value 9.8",
			createQuery: []string{"CREATE TABLE table1 (column1 float NOT NULL DEFAULT 9.8)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO", FieldDefault: "9.8", FieldExtra: "9.8"}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE table1 (column1 datetime NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []utils.ForeignKeysType{},
			indexKeys:   []utils.IndexType{},
			fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO", FieldDefault: "2020-05-30 00:42:05", FieldExtra: "2020-05-30 00:42:05"}},
			wantErr:     false,
		},
		// {
		// 	name:      "MySQL field col1 with type ID which is not null having primary key constraint",
		// 	scQueries: []string{`type table1 { id : ID! @primary, name : String!}`},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "PRI"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name:        "MySQL field col1 with type ID which is not null having foreign key constraint created through or not from space cloud",
		// 	createQuery: []string{`type table2 { id : ID! @primary, name : String! }`, `type table1 { id : ID! @primary, column1 : ID! @foreign(table:"table2",to:"id") }`},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String which is not null having foreign key constraint created through or not from space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer which is not null having foreign key constraint created through or not from space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float which is not null having foreign key constraint created through or not from space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime which is not null having foreign key constraint created through or not from space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON which is not null having foreign key constraint created through or not from space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{{TableName: "table1", ColumnName: firstColumn, RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
		// 	indexKeys:   []utils.IndexType{},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO", FieldKey: "MUL"}},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "tinyint", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON which is not null having single unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON which is not null having single unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON which is not null having single index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "json", FieldNull: "NO"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON which is not null having single index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys:   []utils.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "text", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "text", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "float", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "float", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint not created through space cloud",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "json", FieldNull: "NO"}, {FieldName: secondColumn, FieldType: "json", FieldNull: "NO"}},
		// 	indexKeys: []utils.IndexType{
		// 		{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 		{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{},
		// 	wantErr:     false,
		// },
		// {
		// 	name: "identify varchar with any size",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "table1",
		// 	},
		// 	foreignKeys: []utils.ForeignKeysType{}, fields: []utils.FieldType{{FieldName: firstColumn, FieldType: "varchar(5550)", FieldNull: "NO", FieldKey: "PRI"}},
		// 	wantErr: false,
		// },
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("DescribeTable() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create table in db
			if err := db.RawBatch(context.Background(), tt.createQuery); err != nil {
				t.Errorf("DescribeTable() couldn't insert rows got error - (%v)", err)
			}

			got, got1, got2, err := db.DescribeTable(tt.args.ctx, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.fields) {
				t.Errorf("DescribeTable() got = %v, want %v", got, tt.fields)
			}
			if !reflect.DeepEqual(got1, tt.foreignKeys) {
				t.Errorf("DescribeTable() got1 = %v, want %v", got1, tt.foreignKeys)
			}
			if !reflect.DeepEqual(got2, tt.indexKeys) {
				t.Errorf("DescribeTable() got2 = %v, want %v", got2, tt.indexKeys)
			}
			if _, err := db.client.Exec("DROP TABLE IF EXISTS table1"); err != nil {
				t.Log("DescribeTable() Couldn't truncate table", err)
			}
			if _, err := db.client.Exec("DROP TABLE IF EXISTS table2"); err != nil {
				t.Log("DescribeTable() Couldn't truncate table", err)
			}
		})
	}
}

func getIndexName(tableName, indexName string) string {
	return fmt.Sprintf("index__%s__%s", tableName, indexName)
}

func getConstraintName(tableName, columnName string) string {
	return fmt.Sprintf("c_%s_%s", tableName, columnName)
}
