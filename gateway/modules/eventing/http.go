package eventing

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *Module) logInvocation(ctx context.Context, eventID string, payload []byte, responseStatusCode int, responseBody, errorMsg string) error {
	invocationDoc := &model.InvocationDocument{
		EventID:            eventID,
		InvocationTime:     time.Now().Format(time.RFC3339),
		RequestPayload:     string(payload),
		ResponseStatusCode: responseStatusCode,
		ResponseBody:       responseBody,
		ErrorMessage:       errorMsg,
	}
	createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.One, IsBatch: true}
	if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
		return errors.New("eventing module couldn't log the request - " + err.Error())
	}
	return nil
}

// MakeInvocationHTTPRequest fires an http request and returns a response
func (s *Module) MakeInvocationHTTPRequest(ctx context.Context, method string, eventDoc *model.EventDocument, token, scToken string, params, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, method, eventDoc.URL, bytes.NewBuffer(data))
	if err != nil {
		if err := s.logInvocation(ctx, eventDoc.ID, data, 0, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	// Add the headers
	if token != "" {
		// Add the token only if its provided
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-sc-token", "Bearer "+scToken)

	// Create a http client and fire the request
	client := &http.Client{}

	// if s.storeType && s.isConsulConnectEnabled && strings.Contains(url, "https") && strings.Contains(url, ".consul") {
	// 	 client = s.consulService.HTTPClient()
	// }

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		if err := s.logInvocation(ctx, eventDoc.ID, data, resp.StatusCode, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}
	defer utils.CloseTheCloser(resp.Body)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err := s.logInvocation(ctx, eventDoc.ID, data, 0, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	if err := json.Unmarshal(responseBody, vPtr); err != nil {
		err := s.logInvocation(ctx, eventDoc.ID, data, resp.StatusCode, string(responseBody), err.Error())
		if err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := s.logInvocation(ctx, eventDoc.ID, data, resp.StatusCode, string(responseBody), err.Error())
		if err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	if err := s.logInvocation(ctx, eventDoc.ID, data, resp.StatusCode, string(responseBody), ""); err != nil {
		logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
		return err
	}

	return nil
}
