package syncman

import "github.com/spaceuptech/space-cloud/gateway/config"

// AdminSyncmanInterface is an interface consisting of functions of admin module used by eventing module
type AdminSyncmanInterface interface {
	GetInternalAccessToken() (string, error)
	IsTokenValid(token, resource, op string, attr map[string]string) error
	ValidateProjectSyncOperation(c *config.Config, project *config.Project) bool
	SetConfig(admin *config.Admin) error
	GetConfig() *config.Admin
}

type preparedQueryResponse struct {
	ID        string       `json:"id"`
	DBAlias   string       `json:"db"`
	SQL       string       `json:"sql"`
	Arguments []string     `json:"arguments" yaml:"arguments"`
	Rule      *config.Rule `json:"rule"`
}

type dbRulesResponse struct {
	IsRealTimeEnabled bool                    `json:"isRealtimeEnabled"`
	Rules             map[string]*config.Rule `json:"rules"`
}

type dbSchemaResponse struct {
	Schema string `json:"schema"`
}
