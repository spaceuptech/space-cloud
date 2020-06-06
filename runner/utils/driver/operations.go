package driver

import (
	"context"
	"net/http"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateProject creates project
func (m *Module) CreateProject(ctx context.Context, project *model.Project) error {
	return m.driver.CreateProject(ctx, project)
}

// DeleteProject deletes project
func (m *Module) DeleteProject(ctx context.Context, projectID string) error {
	return m.driver.DeleteProject(ctx, projectID)
}

// ApplyService applies service
func (m *Module) ApplyService(ctx context.Context, service *model.Service) error {
	err := m.driver.ApplyService(ctx, service)
	if err == nil {
		m.metricHook(service.ProjectID)
	}
	return err
}

// GetLogs get logs of specified service
func (m *Module) GetLogs(ctx context.Context, projectID, serviceID, taskID, replica string, w http.ResponseWriter, r *http.Request) error {
	return m.driver.GetLogs(ctx, projectID, serviceID, taskID, replica, w, r)
}

// GetServices gets services
func (m *Module) GetServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	return m.driver.GetServices(ctx, projectID)
}

// DeleteService delete's service
func (m *Module) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	return m.driver.DeleteService(ctx, projectID, serviceID, version)
}

// AdjustScale adjust's scale
func (m *Module) AdjustScale(service *model.Service, activeReqs int32) error {
	return m.driver.AdjustScale(service, activeReqs)
}

// WaitForService waits for service
func (m *Module) WaitForService(service *model.Service) error {
	return m.driver.WaitForService(service)
}

// Type gets driver type
func (m *Module) Type() model.DriverType {
	return m.driver.Type()
}

// ApplyServiceRoutes applies service routes
func (m *Module) ApplyServiceRoutes(ctx context.Context, projectID, serviceID string, routes model.Routes) error {
	return m.driver.ApplyServiceRoutes(ctx, projectID, serviceID, routes)
}

// GetServiceRoutes get's service routes
func (m *Module) GetServiceRoutes(ctx context.Context, projectID string) (map[string]model.Routes, error) {
	return m.driver.GetServiceRoutes(ctx, projectID)
}

// CreateSecret create's secret
func (m *Module) CreateSecret(projectID string, secretObj *model.Secret) error {
	return m.driver.CreateSecret(projectID, secretObj)
}

// ListSecrets list's secrets
func (m *Module) ListSecrets(projectID string) ([]*model.Secret, error) {
	return m.driver.ListSecrets(projectID)
}

// DeleteSecret delete's secret
func (m *Module) DeleteSecret(projectID, secretName string) error {
	return m.driver.DeleteSecret(projectID, secretName)
}

// SetKey set's key for secret
func (m *Module) SetKey(projectID, secretName, secretKey string, secretObj *model.SecretValue) error {
	return m.driver.SetKey(projectID, secretName, secretKey, secretObj)
}

// DeleteKey delete's key of secret
func (m *Module) DeleteKey(projectID, secretName, secretKey string) error {
	return m.driver.DeleteKey(projectID, secretName, secretKey)
}

// SetFileSecretRootPath set's file secret root path
func (m *Module) SetFileSecretRootPath(projectID string, secretName, rootPath string) error {
	return m.driver.SetFileSecretRootPath(projectID, secretName, rootPath)
}
