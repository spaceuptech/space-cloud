package eventing

import (
	"context"
	"errors"
	"math/rand"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// HookDBCreateIntent handles the create intent request
func (m *Module) CreateFileIntentHook(ctx context.Context, req *model.CreateFileRequest, meta map[string]interface{}) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := ksuid.New().String()

	rules := m.getMatchingRules(utils.EventFileCreate, map[string]string{})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
			batchID, utils.EventStatusIntent, rule.Url, &model.QueueEventRequest{
				Type: utils.EventFileCreate,
				Payload: &model.FilePayload{
					Meta: meta,
					Path: req.Path,
				},
			}))
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HookDBDeleteIntent handles the delete intent requests
func (m *Module) DeleteFileIntentHook(ctx context.Context, path string, meta map[string]interface{}) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create a unique batch id and token
	batchID := ksuid.New().String()
	token := rand.Intn(utils.MaxEventTokens)

	rules := m.getMatchingRules(utils.EventFileDelete, map[string]string{})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
			batchID, utils.EventStatusIntent, rule.Url, &model.QueueEventRequest{
				Type: utils.EventFileDelete,
				Payload: &model.FilePayload{
					Meta: meta,
					Path: path,
				},
			}))
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}
