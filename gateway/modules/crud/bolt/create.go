package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (b *Bolt) Create(ctx context.Context, project, col string, req *model.CreateRequest) (int64, error) {
	objs := []interface{}{}
	switch req.Operation {
	case utils.All, utils.One:
		if req.Operation == utils.One {
			doc, ok := req.Document.(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("error inserting into bboltdb cannot assert document to map")
			}
			objs = append(objs, doc)
		} else {
			docs, ok := req.Document.([]interface{})
			if !ok {
				return 0, fmt.Errorf("error inserting into bboltdb cannot assert document to slice of interface")
			}
			objs = docs
		}

		if err := b.client.Update(func(tx *bbolt.Tx) error {

			for _, objToSet := range objs {
				// get _id from create request
				id, ok := objToSet.(map[string]interface{})["_id"]
				if !ok {
					return fmt.Errorf("error creating _id not found in create request")
				}
				// check if specified already exists in database
				count, _, err := b.Read(ctx, project, col, &model.ReadRequest{
					Find: map[string]interface{}{
						"_id": id,
					},
					Operation: utils.Count,
				})
				if count > 0 || err != nil {
					logrus.Errorf("error inserting into bboltdb data already exists - %v", err)
				}

				b, err := tx.CreateBucketIfNotExists([]byte(project))
				if err != nil {
					logrus.Errorf("error creating bucket in bboltdb while inserting- %v", err)
					return err
				}

				// store value as json string
				value, err := json.Marshal(&objToSet)
				if err != nil {
					logrus.Errorf("error marshalling while inserting in bboltdb - %v", err)
					return err
				}

				// insert document in bucket
				if err = b.Put([]byte(fmt.Sprintf("%s/%s", col, id)), value); err != nil {
					return fmt.Errorf("error inserting in bbolt db - %v", err)
				}
			}
			return nil
		}); err != nil {
			return 0, err
		}
		return int64(len(objs)), nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
