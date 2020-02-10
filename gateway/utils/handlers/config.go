package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"isProd": adminMan.LoadEnv(), "version": utils.BuildVersion})
	}
}

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		status, token, err := adminMan.Login(req.User, req.Pass)
		if err != nil {
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		c := syncMan.GetGlobalConfig()
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "projects": c.Projects})
	}
}

// HandleRefreshToken creates the refresh-token endpoint
func HandleRefreshToken(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()
		
		newToken, err := adminMan.RefreshToken(token)
		if err != nil {
			logrus.Errorf("Error while refreshing token handleRefreshToken - %s ",err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return 
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": newToken})
	}
}



// HandleLoadProjects returns the handler to load the projects via a REST endpoint
func HandleLoadProjects(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load config from file
		c := syncMan.GetGlobalConfig()

		// Create a projects array
		projects := []*config.Project{}

		// Iterate over all projects
		for _, p := range c.Projects {
			// Add the project to the array if user has read access
			_, err := adminMan.IsAdminOpAuthorised(token, p.ID)
			if err == nil {
				projects = append(projects, p)
			}

			// Add an empty collections object is not present already
			for k, v := range p.Modules.Crud {
				if v.Collections == nil {
					p.Modules.Crud[k].Collections = map[string]*config.TableRule{}
				}
			}
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"projects": projects})
	}
}

// HandleGlobalConfig returns the handler to store the global config of a project via a REST endpoint
func HandleGlobalConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the body of the request
		c := new(config.Project)
		err := json.NewDecoder(r.Body).Decode(c)
		defer utils.CloseTheCloser(r.Body)

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type - " + err.Error()})
			return
		}

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, c.ID)
		if err != nil {
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		if err := syncMan.SetProjectGlobalConfig(ctx, c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

//HandleGetGlobalConfig returns handler to get config of the project
func HandleGetGlobalConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"aesKey": "", "secret": project.Secret, "name": project.Name})
	}
}

// HandleStoreProjectConfig returns the handler to store the config of a project via a REST endpoint
func HandleStoreProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the body of the request
		c := new(config.Project)
		err := json.NewDecoder(r.Body).Decode(c)
		defer utils.CloseTheCloser(r.Body)

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type - " + err.Error()})
			return
		}

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, c.ID)
		if err != nil {
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		if err := syncMan.SetProjectConfig(ctx, c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleDeleteProjectConfig returns the handler to delete the config of a project via a REST endpoint
func HandleDeleteProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Give negative acknowledgement
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "Operation not supported"})
	}
}
