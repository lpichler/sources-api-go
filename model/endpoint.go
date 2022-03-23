package model

import (
	"strconv"
	"time"

	"github.com/RedHatInsights/sources-api-go/util"
)

type Endpoint struct {
	AvailabilityStatus
	Pause

	ID        int64     `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Role                    *string `json:"role,omitempty"`
	Port                    *int    `json:"port,omitempty"`
	Default                 *bool   `json:"default,omitempty"`
	Scheme                  *string `json:"scheme,omitempty"`
	Host                    *string `json:"host,omitempty"`
	Path                    *string `json:"path,omitempty"`
	VerifySsl               *bool   `json:"verify_ssl,omitempty"`
	CertificateAuthority    *string `json:"certificate_authority,omitempty"`
	ReceptorNode            *string `json:"receptor_node,omitempty"`
	AvailabilityStatusError *string `json:"availability_status_error,omitempty"`

	SourceID int64 `json:"source_id"`
	Source   Source

	TenantID int64
	Tenant   Tenant
}

func (endpoint *Endpoint) ToEvent() interface{} {
	asEvent := AvailabilityStatusEvent{AvailabilityStatus: util.StringValueOrNil(endpoint.AvailabilityStatus.AvailabilityStatus),
		LastAvailableAt: util.DateTimeToRecordFormat(endpoint.LastAvailableAt),
		LastCheckedAt:   util.DateTimeToRecordFormat(endpoint.LastCheckedAt)}

	endpointEvent := &EndpointEvent{
		AvailabilityStatusEvent: asEvent,
		PauseEvent:              PauseEvent{PausedAt: util.DateTimeToRecordFormat(endpoint.PausedAt)},
		ID:                      endpoint.ID,
		CertificateAuthority:    endpoint.CertificateAuthority,
		Host:                    endpoint.Host,
		Port:                    endpoint.Port,
		ReceptorNode:            endpoint.ReceptorNode,
		Role:                    endpoint.Role,
		Scheme:                  endpoint.Scheme,
		SourceID:                endpoint.SourceID,
		VerifySsl:               endpoint.VerifySsl,
		Default:                 endpoint.Default,
		Path:                    endpoint.Path,
		CreatedAt:               util.DateTimeToRecordFormat(endpoint.CreatedAt),
		UpdatedAt:               util.DateTimeToRecordFormat(endpoint.UpdatedAt),
		AvailabilityStatusError: util.StringValueOrNil(endpoint.AvailabilityStatusError),
		Tenant:                  &endpoint.Tenant.ExternalTenant,
	}

	return endpointEvent
}

func (endpoint *Endpoint) ToResponse() *EndpointResponse {
	id := strconv.FormatInt(endpoint.ID, 10)
	sourceId := strconv.FormatInt(endpoint.SourceID, 10)
	asResponse := AvailabilityStatusResponse{
		AvailabilityStatus: util.StringValueOrNil(endpoint.AvailabilityStatus.AvailabilityStatus),
		LastCheckedAt:      util.DateTimeToRFC3339(endpoint.LastCheckedAt),
		LastAvailableAt:    util.DateTimeToRFC3339(endpoint.LastAvailableAt),
	}

	return &EndpointResponse{
		AvailabilityStatusResponse: asResponse,
		PauseResponse:              PauseResponse{PausedAt: util.DateTimeToRFC3339(endpoint.PausedAt)},
		ID:                         id,
		CreatedAt:                  util.DateTimeToRFC3339(endpoint.CreatedAt),
		UpdatedAt:                  util.DateTimeToRFC3339(endpoint.UpdatedAt),
		Role:                       endpoint.Role,
		Port:                       endpoint.Port,
		Default:                    endpoint.Default,
		Scheme:                     endpoint.Scheme,
		Host:                       endpoint.Host,
		Path:                       endpoint.Path,
		VerifySsl:                  endpoint.VerifySsl,
		CertificateAuthority:       endpoint.CertificateAuthority,
		ReceptorNode:               endpoint.ReceptorNode,
		AvailabilityStatusError:    endpoint.AvailabilityStatusError,
		SourceID:                   sourceId,
	}
}

func (endpoint *Endpoint) ToEmailNotificationInfo(previousStatus string) *EmailNotificationInfo {
	return &EmailNotificationInfo{
		SourceID:                   strconv.FormatInt(endpoint.SourceID, 10),
		SourceName:                 endpoint.Source.Name,
		ResourceDisplayName:        "Endpoint",
		CurrentAvailabilityStatus:  endpoint.AvailabilityStatus.AvailabilityStatus,
		PreviousAvailabilityStatus: previousStatus,
	}
}
