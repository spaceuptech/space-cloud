package schema

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

func TestSchema_generateCreationQueries(t *testing.T) {
	type fields struct {
		SchemaDoc          schemaType
		crud               *crud.Module
		project            string
		config             config.Crud
		removeProjectScope bool
	}
	type args struct {
		ctx           context.Context
		dbAlias       string
		tableName     string
		project       string
		parsedSchema  schemaType
		currentSchema schemaCollection
	}

	crudPostgres := crud.Init(false)
	crudPostgres.SetConfig("test", config.Crud{"postgres": {Type: "sql-postgres", Enabled: false}})

	crudMySql := crud.Init(false)
	crudMySql.SetConfig("test", config.Crud{"mysql": {Type: "sql-mysql", Enabled: false}})

	crudSqlServer := crud.Init(false)
	crudSqlServer.SetConfig("test", config.Crud{"sqlserver": {Type: "sql-sqlserver", Enabled: false}})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "adding two columns",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)", "ALTER TABLE test.table1 ADD col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 datetime"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 datetime"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 bigint NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP INDEX c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NOT NULL", "ALTER TABLE test.table1 ADD PRIMARY KEY (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NULL", "ALTER TABLE test.table1 DROP PRIMARY KEY"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP FOREIGN KEY c_table1_col1", "ALTER TABLE test.table1 DROP INDEX c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},

		// //sql-server

		{
			name: "adding two columns",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)", "ALTER TABLE test.table1 ADD col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 timestamp NULL"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 timestamp NULL"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 bigint NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 bigint NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 PRIMARY KEY CLUSTERED (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},

		// //postgres
		{
			name: "adding two columns",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)", "ALTER TABLE test.table1 ADD COLUMN col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 timestamp"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 timestamp"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 PRIMARY KEY (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc:          tt.fields.SchemaDoc,
				crud:               tt.fields.crud,
				project:            tt.fields.project,
				config:             tt.fields.config,
				removeProjectScope: tt.fields.removeProjectScope,
			}
			got, err := s.generateCreationQueries(tt.args.ctx, tt.args.dbAlias, tt.args.tableName, tt.args.project, tt.args.parsedSchema, tt.args.currentSchema)
			if (err != nil) != tt.wantErr {
				t.Errorf("name = %v, Schema.generateCreationQueries() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("name = %v, Schema.generateCreationQueries() = %v, want %v", tt.name, got, tt.want)
				return
			}

			for i, v := range got {
				if tt.want[i] != v {
					t.Errorf("name = %v, Schema.generateCreationQueries() = %v, want %v", tt.name, got, tt.want)
					break
				}
			}
		})
	}
}
