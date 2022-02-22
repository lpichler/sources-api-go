package dao

import (
	"fmt"

	m "github.com/RedHatInsights/sources-api-go/model"
	"github.com/RedHatInsights/sources-api-go/util"
)

type ApplicationAuthenticationDaoImpl struct {
	TenantID *int64
}

func (a *ApplicationAuthenticationDaoImpl) ApplicationAuthenticationsByApplications(applications []m.Application) ([]m.ApplicationAuthentication, error) {
	var applicationAuthentications []m.ApplicationAuthentication

	applicationIDs := make([]int64, 0)
	for _, value := range applications {
		applicationIDs = append(applicationIDs, value.ID)
	}

	err := DB.Preload("Tenant").Where("application_id IN ?", applicationIDs).Find(&applicationAuthentications).Error
	if err != nil {
		return nil, err
	}

	return applicationAuthentications, nil
}

func (a *ApplicationAuthenticationDaoImpl) ApplicationAuthenticationsByAuthentications(authentications []m.Authentication) ([]m.ApplicationAuthentication, error) {
	var applicationAuthentications []m.ApplicationAuthentication

	authenticationUIDs := make([]string, 0)
	for _, value := range authentications {
		authenticationUIDs = append(authenticationUIDs, value.ID)
	}

	result := DB.Preload("Tenant").Where("authentication_uid IN ?", authenticationUIDs).Find(&applicationAuthentications)
	if result.Error != nil {
		return nil, result.Error
	}

	return applicationAuthentications, nil
}

func (a *ApplicationAuthenticationDaoImpl) ApplicationAuthenticationsByResource(resourceType string, applications []m.Application, authentications []m.Authentication) ([]m.ApplicationAuthentication, error) {
	if resourceType == "Source" {
		return a.ApplicationAuthenticationsByApplications(applications)
	}

	return a.ApplicationAuthenticationsByAuthentications(authentications)
}

func (a *ApplicationAuthenticationDaoImpl) List(limit int, offset int, filters []util.Filter) ([]m.ApplicationAuthentication, int64, error) {
	appAuths := make([]m.ApplicationAuthentication, 0, limit)
	query := DB.Debug().Model(&m.ApplicationAuthentication{}).
		Offset(offset).
		Where("tenant_id = ?", a.TenantID)

	query, err := applyFilters(query, filters)
	if err != nil {
		return nil, 0, err
	}

	count := int64(0)
	query.Count(&count)

	result := query.Limit(limit).Find(&appAuths)
	return appAuths, count, result.Error
}

func (a *ApplicationAuthenticationDaoImpl) GetById(id *int64) (*m.ApplicationAuthentication, error) {
	appAuth := &m.ApplicationAuthentication{ID: *id}
	result := DB.First(&appAuth)

	return appAuth, result.Error
}

func (a *ApplicationAuthenticationDaoImpl) Create(appAuth *m.ApplicationAuthentication) error {
	result := DB.Create(appAuth)
	return result.Error
}

func (a *ApplicationAuthenticationDaoImpl) Update(appAuth *m.ApplicationAuthentication) error {
	result := DB.Updates(appAuth)
	return result.Error
}

func (a *ApplicationAuthenticationDaoImpl) Delete(id *int64) error {
	appAuth := &m.ApplicationAuthentication{ID: *id}
	if result := DB.Delete(appAuth); result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete application id %v", *id)
	}

	return nil
}

func (a *ApplicationAuthenticationDaoImpl) Tenant() *int64 {
	return a.TenantID
}
