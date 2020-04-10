package metrics

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/spaceuptech/space-api-go/db"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

type mockAdminInterface struct {
	mock.Mock
}

func (i *mockAdminInterface) GetClusterID() string {
	return i.Called().String(0)
}

type mockSyncManInterface struct {
	mock.Mock
}

func (i *mockSyncManInterface) GetNodesInCluster() int {
	return i.Called().Int(0)
}

func TestModule_generateMetricsRequest(t *testing.T) {
	type args struct {
		project *config.Project
		ssl     *config.SSL
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 map[string]interface{}
		want3 map[string]interface{}
	}{
		{
			name: "valid config",
			args: args{
				project: &config.Project{
					ID: "project",
					Modules: &config.Modules{
						Crud: map[string]*config.CrudStub{"db": {
							Type: "postgres",
							Collections: map[string]*config.TableRule{
								"table1": {},
							},
							Enabled: true,
						}},
						Auth: config.Auth{"auth": {
							ID:      "email",
							Enabled: true,
						}},
						Services: &config.ServicesModule{Services: map[string]*config.Service{"service": {}}},
						FileStore: &config.FileStore{
							Enabled:   true,
							StoreType: "local",
							Rules:     []*config.FileRule{{ID: "file"}},
						},
						Eventing: config.Eventing{
							Enabled: true,
							DBAlias: "db",
							Rules: map[string]config.EventingRule{
								"type": {
									Type: "type",
								},
							},
						},
						LetsEncrypt: config.LetsEncrypt{
							ID:                 "letsEncrypt",
							WhitelistedDomains: []string{"1"},
						},
						Routes: config.Routes{{
							ID: "route",
						}},
					},
				},
				ssl: &config.SSL{
					Enabled: true,
				},
			},
			want:  "clusterID",
			want1: "project",
			want2: map[string]interface{}{
				"nodes":        1,
				"os":           runtime.GOOS,
				"is_prod":      false,
				"version":      utils.BuildVersion,
				"distribution": "ce",
				"last_updated": time.Now().UnixNano() / int64(time.Millisecond),
				"ssl_enabled":  true,
				"project":      "project",
				"crud": map[string]interface{}{
					"db": map[string]interface{}{
						"tables": -2,
					},
				},
				"databases": map[string][]string{
					"databases": {"postgres"},
				},
				"file_store_store_type": "local",
				"file_store_rules":      1,
				"auth": map[string]interface{}{
					"providers": 1,
				},
				"services":     1,
				"lets_encrypt": 1,
				"routes":       1,
				"total_events": 1,
			},
			want3: map[string]interface{}{"start_time": ""},
		},
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1, got2, got3 := m.generateMetricsRequest(tt.args.project, tt.args.ssl)
			if got != tt.want {
				t.Errorf("generateMetricsRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("generateMetricsRequest() got1 = %v, want %v", got1, tt.want1)
			}
			for key, value := range tt.want2 {
				if key == "last_updated" {
					continue
				}
				gotValue, ok := got2[key]
				if !ok {
					t.Errorf("createCrudDocuments() key = %s doesn't exist in result", key)
					continue
				}
				if !reflect.DeepEqual(gotValue, value) {
					t.Errorf("createCrudDocuments() got value = %v %T want = %v %T", gotValue, gotValue, value, value)
				}
			}
			if _, ok := got3["start_time"]; !ok {
				t.Errorf("generateMetricsRequest() got3 = %v, want %v", got3, tt.want3)
			}

		})
	}
}

func TestModule_updateSCMetrics(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		clusterID string
		projectID string
		set       map[string]interface{}
		min       map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			m.updateSCMetrics(tt.args.clusterID, tt.args.projectID, tt.args.set, tt.args.min)
		})
	}
}