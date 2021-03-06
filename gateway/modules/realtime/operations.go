package realtime

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	authHelpers "github.com/spaceuptech/space-cloud/gateway/modules/auth/helpers"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Subscribe performs the realtime subscribe operation.
func (m *Module) Subscribe(clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error) {
	// Create a 20 second context to process request
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if data.Group == "" || data.DBType == "" || data.Where == nil {
		return nil, errors.New("invalid request parameters provided")
	}
	readReq := model.ReadRequest{Find: data.Where, Operation: utils.All}

	// Check if the user is authorised to make the request
	actions, reqParams, err := m.auth.IsReadOpAuthorised(ctx, data.Project, data.DBType, data.Group, data.Token, &readReq, model.ReturnWhereStub{})
	if err != nil {
		return nil, err
	}

	return m.DoRealtimeSubscribe(ctx, clientID, data, actions, reqParams, sendFeed)
}

// DoRealtimeSubscribe makes the realtime query
func (m *Module) DoRealtimeSubscribe(ctx context.Context, clientID string, data *model.RealtimeRequest, actions *model.PostProcess, reqParams model.RequestParams, sendFeed model.SendFeed) ([]*model.FeedData, error) {
	readReq := &model.ReadRequest{Find: data.Where, Operation: utils.All}
	if data.Options.SkipInitial {
		m.AddLiveQuery(data.ID, data.Project, data.DBType, data.Group, clientID, data.Where, actions, sendFeed)
		return []*model.FeedData{}, nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, _, err := m.crud.Read(ctx2, data.DBType, data.Group, readReq, reqParams)
	if err != nil {
		return nil, err
	}

	_ = authHelpers.PostProcessMethod(ctx, m.aesKey, actions, result)

	feedData := make([]*model.FeedData, 0)
	array, ok := result.([]interface{})
	if ok {
		for _, row := range array {
			// Get the appropriate find object
			find := m.prepareFindObject(data.DBType, data.Group, row.(map[string]interface{}))

			// Create the feed data object
			feedData = append(feedData, &model.FeedData{
				Group:     data.Group,
				Type:      utils.RealtimeInitial,
				TimeStamp: 1,
				Find:      find,
				DBType:    data.DBType,
				Payload:   row,
				QueryID:   data.ID,
			})
		}
	}

	// Add the live query
	m.AddLiveQuery(data.ID, data.Project, data.DBType, data.Group, clientID, data.Where, actions, sendFeed)

	return feedData, nil
}

// Unsubscribe performs the realtime unsubscribe operation.
func (m *Module) Unsubscribe(ctx context.Context, data *model.RealtimeRequest, clientID string) error {
	return m.RemoveLiveQuery(ctx, data.DBType, data.Group, clientID, data.ID)
}

// HandleRealtimeEvent handles an incoming realtime event from the eventing module
func (m *Module) HandleRealtimeEvent(ctxRoot context.Context, eventDoc *model.CloudEventPayload) error {

	ids := m.syncMan.GetSpaceCloudNodeIDs(m.project)

	// Create wait group
	var wg sync.WaitGroup
	wg.Add(len(ids))

	// Create success & error channels
	successCh := make(chan struct{}, 1)
	errCh := make(chan error, len(ids))

	ctx, cancel := context.WithTimeout(ctxRoot, 5*time.Second)
	defer cancel()

	for _, i := range ids {
		go func(id string) {
			defer wg.Done()

			if err := m.pubsubClient.Send(ctx, getSendTopic(id), eventDoc); err != nil {
				errCh <- err
				return
			}
		}(i)
	}

	go func() {
		wg.Wait()
		successCh <- struct{}{}
	}()

	select {
	case err := <-errCh:
		cancel()
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Realtime Module: Event handler error", map[string]interface{}{"error": err})
		return err

	case <-successCh:
		return nil
	}
}

// ProcessRealtimeRequests handles an incoming realtime process event
func (m *Module) ProcessRealtimeRequests(ctx context.Context, eventDoc *model.CloudEventPayload) error {

	dbEvent := new(model.DatabaseEventMessage)
	if err := mapstructure.Decode(eventDoc.Data, dbEvent); err != nil {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Realtime Module Request Handler Error", map[string]interface{}{"error": err})
		return err
	}

	t, _ := time.Parse(time.RFC3339, eventDoc.Time)
	feedData := &model.FeedData{
		Type:      eventingToRealtimeEvent(eventDoc.Type),
		Payload:   dbEvent.Doc,
		TimeStamp: t.UnixNano() / int64(time.Millisecond),
		Group:     dbEvent.Col,
		DBType:    dbEvent.DBType,
		Find:      dbEvent.Find,
	}

	m.helperSendFeed(ctx, feedData)

	return nil
}
