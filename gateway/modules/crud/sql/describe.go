package sql

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// DescribeTable return a description of sql table & foreign keys in table
// NOTE: not to be exposed externally
func (s *SQL) DescribeTable(ctx context.Context, col string) ([]model.InspectorFieldType, []model.IndexType, error) {
	fields, err := s.getDescribeDetails(ctx, s.name, col)
	if err != nil {
		return nil, nil, err
	}
	index, err := s.getIndexDetails(ctx, s.name, col)
	if err != nil {
		return nil, nil, err
	}

	return fields, index, nil
}

func (s *SQL) getDescribeDetails(ctx context.Context, project, col string) ([]model.InspectorFieldType, error) {
	queryString := ""
	args := make([]interface{}, 0)
	switch model.DBType(s.dbType) {
	case model.MySQL:
		queryString = `
select a.table_schema  AS 'TABLE_SCHEMA',
       a.table_name AS 'TABLE_NAME',

       a.column_name AS 'COLUMN_NAME',
       a.data_type 'DATA_TYPE',
       a.is_nullable AS 'IS_NULLABLE',
       a.ordinal_position AS 'ORDINAL_POSITION',
       CASE
           WHEN a.column_default = '1' THEN 'true'
           WHEN a.column_default = '0' THEN 'false'
           WHEN a.column_default = "b\'1\'" THEN 'true'
           WHEN a.column_default = "b\'0\'" THEN 'false'
           ELSE coalesce(a.column_default,'')
       END AS 'DEFAULT',
       IF(upper(a.extra) = 'AUTO_INCREMENT', 'true', 'false') AS 'AUTO_INCREMENT',
       coalesce(a.character_maximum_length,0) AS 'CHARACTER_MAXIMUM_LENGTH',
       coalesce(a.numeric_precision,0) AS 'NUMERIC_PRECISION',
       coalesce(a.numeric_scale,0) AS 'NUMERIC_SCALE',
        
       coalesce(d.constraint_name,'') AS 'CONSTRAINT_NAME',
       coalesce(d.delete_rule,'') AS 'DELETE_RULE',
       coalesce(d.referenced_table_schema,'') AS 'REFERENCED_TABLE_SCHEMA',
       coalesce(d.referenced_table_name,'') AS 'REFERENCED_TABLE_NAME',
       coalesce(d.referenced_column_name,'') AS 'REFERENCED_COLUMN_NAME'
from information_schema.columns a
         left join (select x.constraint_schema , x.table_name, x.constraint_name, y.delete_rule,
                           z.referenced_table_schema, z.referenced_table_name, z.referenced_column_name, z.column_name
                    from INFORMATION_SCHEMA.TABLE_CONSTRAINTS x
                             left join INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS y on x.constraint_schema=y.constraint_schema and
                                                                                       x.constraint_name=y.constraint_name and x.table_name=y.table_name
                             left join  INFORMATION_SCHEMA.KEY_COLUMN_USAGE z on x.constraint_schema=z.constraint_schema and
                                                                                 x.constraint_name=z.constraint_name and x.table_name=z.table_name
                    where x.CONSTRAINT_TYPE in ('FOREIGN KEY')) d on a.table_schema = d.constraint_schema
    and a.table_name=d.table_name and a.column_name=d.column_name
where a.table_name= ? and a.table_schema= ? ;
`
		args = append(args, col, project)

	case model.Postgres:
		queryString = `select c.table_schema  AS "TABLE_SCHEMA",
       c.table_name AS "TABLE_NAME",

       c.column_name AS "COLUMN_NAME",
       c.data_type "DATA_TYPE",
       c.is_nullable AS "IS_NULLABLE",
       c.ordinal_position AS "ORDINAL_POSITION",
       SPLIT_PART(REPLACE(coalesce(c.column_default,''),'''',''), '::', 1) AS "DEFAULT",
       case when upper(c.column_default) like 'NEXTVAL%' then 'true' else 'false' end AS "AUTO_INCREMENT",
       coalesce(c.character_maximum_length,0) AS "CHARACTER_MAXIMUM_LENGTH",
       coalesce(c.numeric_precision,0) AS "NUMERIC_PRECISION",
       coalesce(c.numeric_scale,0) AS "NUMERIC_SCALE",

       coalesce(fk.constraint_name,'') AS "CONSTRAINT_NAME",
       coalesce(fk.delete_rule,'') AS "DELETE_RULE",
       coalesce(fk.foreign_table_schema,'') AS "REFERENCED_TABLE_SCHEMA",
       coalesce(fk.foreign_table_name,'') AS "REFERENCED_TABLE_NAME",
       coalesce(fk.foreign_column_name,'') AS "REFERENCED_COLUMN_NAME"
from information_schema.columns c
         left join (SELECT tc.table_schema, tc.constraint_name, tc.table_name, kcu.column_name,
                           ccu.table_schema AS foreign_table_schema, ccu.table_name AS foreign_table_name,
                           ccu.column_name AS foreign_column_name, rc.delete_rule as delete_rule
                    FROM information_schema.table_constraints AS tc
                             INNER JOIN information_schema.key_column_usage AS kcu
                                        ON tc.constraint_name = kcu.constraint_name
                                            AND tc.table_schema = kcu.table_schema
                             INNER JOIN information_schema.constraint_column_usage AS ccu
                                        ON ccu.constraint_name = tc.constraint_name
                                            AND ccu.table_schema = tc.table_schema
                             INNER JOIN  information_schema.referential_constraints rc
                                         on rc.constraint_name = tc.constraint_name and rc.constraint_schema = tc.table_schema
                    WHERE tc.constraint_type = 'FOREIGN KEY' ) fk
                   on fk.table_name = c.table_name and fk.column_name = c.column_name and fk.table_schema = c.table_schema
where c.table_name = $1 and c.table_schema = $2
order by c.ordinal_position;`

		args = append(args, col, project)
	case model.SQLServer:

		queryString = `SELECT DISTINCT C.COLUMN_NAME as 'Field', C.IS_NULLABLE as 'Null' ,
                case when C.DATA_TYPE = 'varchar' then concat(C.DATA_TYPE,'(',REPLACE(c.CHARACTER_MAXIMUM_LENGTH,'-1','max'),')') else C.DATA_TYPE end as 'Type',
                REPLACE(REPLACE(REPLACE(coalesce(C.COLUMN_DEFAULT,''),'''',''),'(',''),')','') as 'Default',
                CASE
                    WHEN TC.CONSTRAINT_TYPE = 'PRIMARY KEY' THEN 'PRI'
                    WHEN TC.CONSTRAINT_TYPE = 'UNIQUE' THEN 'UNI'
                    WHEN TC.CONSTRAINT_TYPE = 'FOREIGN KEY' THEN 'MUL'
                    ELSE isnull(TC.CONSTRAINT_TYPE,'')
                    END AS 'Key',
                coalesce(c.CHARACTER_MAXIMUM_LENGTH,50) AS 'VarcharSize',
                CASE WHEN I.NAME IS NOT NULL THEN 'true' ELSE 'false' END AS 'AutoIncrement'
FROM INFORMATION_SCHEMA.COLUMNS AS C
         LEFT JOIN SYS.IDENTITY_COLUMNS I ON C.table_name = OBJECT_NAME(I.OBJECT_ID) AND C.COLUMN_NAME = I.NAME
         FULL JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE AS CC
                   ON C.COLUMN_NAME = CC.COLUMN_NAME
         FULL JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
                   ON CC.CONSTRAINT_NAME = TC.CONSTRAINT_NAME
WHERE C.TABLE_SCHEMA=@p2 AND C.table_name = @p1`

		args = append(args, col, project)
	}
	rows, err := s.getClient().QueryxContext(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := make([]model.InspectorFieldType, 0)
	count := 0
	for rows.Next() {
		count++
		fieldType := new(model.InspectorFieldType)

		if err := rows.StructScan(fieldType); err != nil {
			return nil, err
		}

		result = append(result, *fieldType)
	}
	if count == 0 {
		return result, errors.New(s.dbType + ":" + col + " not found during inspection")
	}
	return result, nil
}

func (s *SQL) getIndexDetails(ctx context.Context, project, col string) ([]model.IndexType, error) {
	queryString := ""
	switch model.DBType(s.dbType) {

	case model.MySQL:
		queryString = `
select b.table_schema AS 'TABLE_SCHEMA',
       b.table_name AS 'TABLE_NAME',
       b.column_name AS 'COLUMN_NAME',
       b.index_name AS 'INDEX_NAME',
       b.seq_in_index AS 'SEQ_IN_INDEX',
	   case when b.collation = "A" then "asc" else "desc" end as SORT,
       case when b.non_unique=0 then true else false end 'IS_UNIQUE',
       case when upper(b.index_name)='PRIMARY' then 1 else 0 end 'IS_PRIMARY'
from INFORMATION_SCHEMA.STATISTICS  b
where b.table_schema= ? and b.table_name= ?;`

	case model.Postgres:
		queryString = `select
    n.nspname AS "TABLE_SCHEMA",
    t.relname AS "TABLE_NAME" ,
    b.attname AS "COLUMN_NAME",
    a.relname AS "INDEX_NAME",
    array_position(i.indkey, b.attnum)+1 "SEQ_IN_INDEX",
	case when i.indoption[array_position(i.indkey, b.attnum)] = 0 then 'asc' else 'desc' END AS "SORT",
    i.indisunique AS "IS_UNIQUE",
    i.indisprimary "IS_PRIMARY"
from pg_class a
         left join pg_namespace n on n.oid = a.relnamespace
         left join pg_index i on a.oid = i.indexrelid and a.relkind='i' and i.indisvalid = true
         left join pg_class t on t.oid = i.indrelid
         left join pg_attribute b on b.attrelid = t.oid and b.attnum = ANY(i.indkey)
where n.nspname= $1 and t.relname= $2;`
	case model.SQLServer:
		queryString = `SELECT 
    	TABLE_NAME = t.name,
    	COLUMN_NAME = col.name,
    	INDEX_NAME = ind.name,
    	SEQ_IN_INDEX = ic.index_column_id,
    	case when ind.is_unique = 0 then 'no' else 'yes' end as IS_UNIQUE,
    	case when ic.is_descending_key = 0 then 'asc' else 'desc' end as SORT 
			FROM 
     			sys.indexes ind 
				INNER JOIN 
     			sys.index_columns ic ON  ind.object_id = ic.object_id and ind.index_id = ic.index_id 
				INNER JOIN 
     			sys.columns col ON ic.object_id = col.object_id and ic.column_id = col.column_id 
				INNER JOIN 
     			sys.tables t ON ind.object_id = t.object_id 
				INNER JOIN 
        	sys.schemas s ON t.schema_id = s.schema_id
			WHERE 
     			ind.is_primary_key = 0  and s.name = @p1 and t.name = @p2 `
	}
	rows, err := s.getClient().QueryxContext(ctx, queryString, []interface{}{project, col}...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := make([]model.IndexType, 0)
	for rows.Next() {
		indexKey := new(model.IndexType)

		if err := rows.StructScan(indexKey); err != nil {
			return nil, err
		}

		result = append(result, *indexKey)
	}
	return result, nil
}
