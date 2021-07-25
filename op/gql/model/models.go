package model

import (
	"fmt"
	"io"
	"strconv"
)

type ObjectRecord struct {
	Id         string                 `json:"id,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type QueryResult struct {
	Edges      []*Edge   `json:"edges,omitempty"`
	PageInfo   *PageInfo `json:"pageInfo,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

type AuthInfoInput struct {
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type OnUpdatedObjectInput struct {
	Name       string                 `json:"name,omitempty"`
	Id         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type NewObjectRecordInput struct {
	Name       string                 `json:"name,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type HealthError struct {
	Test    string `json:"test,omitempty"`
	Message string `json:"message,omitempty"`
}

type UpdateObjectRecordInput struct {
	Name       string                 `json:"name,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type PageInput struct {
	First int    `json:"first,omitempty"`
	After string `json:"after,omitempty"`
}

type ConnectorHealth struct {
	Status      *HealthStatus  `json:"status,omitempty"`
	LastChecked int            `json:"lastChecked,omitempty"`
	Errors      []*HealthError `json:"errors,omitempty"`
}

type OnCreatedObjectInput struct {
	Name       string                 `json:"name,omitempty"`
	Id         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor,omitempty"`
	HasNextPage bool   `json:"hasNextPage,omitempty"`
}

type Edge struct {
	Cursor string        `json:"cursor,omitempty"`
	Node   *ObjectRecord `json:"node,omitempty"`
}

type Object struct {
	Name       string                 `json:"name,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "HEALTHY"
	HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
)

var AllHealthStatus = []HealthStatus{HealthStatusHealthy, HealthStatusUnhealthy}

func (h HealthStatus) IsValid() bool {
	switch h {
	case HealthStatusHealthy, HealthStatusUnhealthy:
		return true
	}
	return false
}
func (h HealthStatus) String() string {
	return string(h)
}
func (h *HealthStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}
	*h = HealthStatus(str)
	if !h.IsValid() {
		return fmt.Errorf("enums must be strings")
	}
	return nil
}
func (h HealthStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(h.String()))
}
