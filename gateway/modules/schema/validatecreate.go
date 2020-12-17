package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SchemaValidator function validates the schema which it gets from module
func (s *Schema) SchemaValidator(ctx context.Context, dbAlias, col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	for schemaKey := range doc {
		if _, p := collectionFields[schemaKey]; !p {
			return nil, errors.New("The field " + schemaKey + " is not present in schema of " + col)
		}
	}

	mutatedDoc := map[string]interface{}{}
	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]

		if fieldValue.IsLinked {
			if ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("cannot insert value for a linked field %s", fieldKey), nil, nil)
			}
			continue
		}

		if fieldValue.IsAutoIncrement {
			continue
		}

		if !ok && fieldValue.IsDefault {
			value = fieldValue.Default
			ok = true
		}

		if fieldValue.Kind == model.TypeID && !ok {
			value = ksuid.New().String()
			ok = true
		}

		if fieldValue.IsCreatedAt || fieldValue.IsUpdatedAt {
			mutatedDoc[fieldKey] = time.Now().UTC()
			continue
		}

		if fieldValue.IsFieldTypeRequired {
			if !ok {
				return nil, errors.New("required field " + fieldKey + " from " + col + " not present in request")
			}
		}

		// check type
		val, err := s.checkType(ctx, dbAlias, col, value, fieldValue)
		if err != nil {
			return nil, err
		}

		mutatedDoc[fieldKey] = val
	}
	return mutatedDoc, nil
}

// ValidateCreateOperation validates schema on create operation
func (s *Schema) ValidateCreateOperation(ctx context.Context, dbAlias, col string, req *model.CreateRequest) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.SchemaDoc == nil {
		return errors.New("schema not initialized")
	}

	v := make([]interface{}, 0)

	switch t := req.Document.(type) {
	case []interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	collection, ok := s.SchemaDoc[dbAlias]
	if !ok {
		return errors.New("No db was found named " + dbAlias)
	}
	collectionFields, ok := collection[col]
	if !ok {
		return nil
	}

	for index, docTemp := range v {
		doc, ok := docTemp.(map[string]interface{})
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("document provided for collection (%s:%s)", dbAlias, col), nil, nil)
		}
		newDoc, err := s.SchemaValidator(ctx, dbAlias, col, collectionFields, doc)
		if err != nil {
			return err
		}

		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}
func (s *Schema) checkType(ctx context.Context, dbAlias, col string, value interface{}, fieldValue *model.FieldType) (interface{}, error) {
	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v)/1000, 0), nil
		case model.TypeInteger:
			return value, nil
		case model.TypeFloat:
			return float64(v), nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case string:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			unitTimeInRFC3339Nano, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid datetime format recieved for field %s in collection %s - use RFC3339 fromat", fieldValue.FieldName, col), nil, nil)
			}
			return unitTimeInRFC3339Nano, nil
		case model.TypeID, model.TypeString, model.TypeTime, model.TypeDate, model.TypeUUID:
			return value, nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got String", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case float32, float64:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v.(float64))/1000, 0), nil
		case model.TypeFloat:
			return value, nil
		case model.TypeInteger:
			return int64(value.(float64)), nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Float", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}
	case bool:
		switch fieldValue.Kind {
		case model.TypeBoolean:
			return value, nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Bool", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case time.Time, *time.Time:
		return v, nil

	case map[string]interface{}:
		if fieldValue.Kind == model.TypeJSON {
			dbType, ok := s.dbAliasDBTypeMapping[dbAlias]
			if !ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown db alias provided (%s)", dbAlias), nil, nil)
			}
			if model.DBType(dbType) == model.Mongo {
				return value, nil
			}
			data, err := json.Marshal(value)
			if err != nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error checking type in schema module unable to marshal data for field having type json", err, nil)
			}
			return string(data), nil
		}
		if fieldValue.Kind != model.TypeObject {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
		}

		return s.SchemaValidator(ctx, dbAlias, col, fieldValue.NestedObject, v)

	case []interface{}:
		if !fieldValue.IsList {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type (array) received for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
		}

		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := s.checkType(ctx, dbAlias, col, value, fieldValue)
			if err != nil {
				return nil, err
			}
			arr[index] = val
		}
		return arr, nil
	default:
		if !fieldValue.IsFieldTypeRequired {
			return nil, nil
		}

		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("no matching type found for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
	}
}
