package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (r *Routing) modifyRequest(ctx context.Context, modules modulesInterface, route *config.Route, req *http.Request) (string, interface{}, error) {
	// Return if the rule is allow
	if route.Rule == nil || route.Rule.Rule == "allow" {
		return "", nil, nil
	}

	// Extract the token
	token := utils.GetTokenFromHeader(req)

	// Extract the params only if content-type is `application/json`
	var params interface{}
	var data []byte
	var err error
	if req.Header.Get("Content-Type") == "application/json" {
		data, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return "", nil, err
		}

		if err := json.Unmarshal(data, &params); err != nil {
			return "", nil, utils.LogError("Unable to unmarshal body to JSON", module, handleRequest, err)
		}
	}

	// Finally we authorize the request
	a := modules.Auth()
	auth, err := a.AuthorizeRequest(ctx, route.Rule, route.Project, token, params)
	if err != nil {
		return "", nil, err
	}

	// Set the headers
	state := map[string]interface{}{"args": params, "auth": auth}
	for k, v := range route.Modify.Headers {
		// Load the string if it exists
		value, err := utils.LoadValue(v, state)
		if err == nil {
			if temp, ok := value.(string); ok {
				v = temp
			}
		}

		// Set the header
		req.Header.Add(k, v)
	}

	// Don't forget to reset the body
	if params != nil {
		// Generate new request body if template was provided
		newParams, err := r.adjustBody(route.Project, token, route, auth, params)
		if err != nil {
			return "", nil, err
		}

		// Marshal it then set it
		data, _ = json.Marshal(newParams)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	return token, auth, err
}

func (r *Routing) modifyResponse(res *http.Response, route *config.Route, token string, auth interface{}) error {
	if res.Header.Get("Content-Type") == "application/json" {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var params interface{}
		if err := json.Unmarshal(data, &params); err != nil {
			return utils.LogError("Unable to unmarshal body to JSON", module, handleResponse, err)
		}

		newParams, err := r.adjustBody(route.Project, token, route, auth, params)
		if err != nil {
			return err
		}

		// Marshal it then set it
		data, _ = json.Marshal(newParams)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	return nil
}
