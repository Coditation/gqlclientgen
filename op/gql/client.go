package client

import (
	"context"
	"encoding/json"
	graphql "github.com/hasura/go-graphql-client"
	model "gqlclientgen/op/gql/model"
	"net/http"
)

type Client struct {
	Client *graphql.Client
}

func NewClient(url string, httpClient *http.Client) *graphql.Client {
	return graphql.NewClient(url, httpClient)
}
func (c *Client) RegisterConnection(ctx context.Context, connectionId string, input model.AuthInfoInput) (interface{}, error) {
	var mutate struct {
		RegisterConnection interface{} `graphql:"(connectionId: &connectionId,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"input":        input,
	}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) RevokeConnection(ctx context.Context, connectionId string) (interface{}, error) {
	var mutate struct {
		RevokeConnection interface{} `graphql:"(connectionId: &connectionId)"`
	}
	variables := map[string]interface{}{"connectionId": connectionId}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) CreateObjectRecord(ctx context.Context, connectionId string, input model.NewObjectRecordInput) (*model.ObjectRecord, error) {
	var mutate struct {
		CreateObjectRecord struct {
			ObjectRecord model.ObjectRecord
		} `graphql:"(connectionId: &connectionId,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"input":        input,
	}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res model.ObjectRecord
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) UpdateObjectRecord(ctx context.Context, connectionId string, id string, input model.UpdateObjectRecordInput) (*model.ObjectRecord, error) {
	var mutate struct {
		UpdateObjectRecord struct {
			ObjectRecord model.ObjectRecord
		} `graphql:"(connectionId: &connectionId,id: &id,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"id":           id,
		"input":        input,
	}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res model.ObjectRecord
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) OnObjectRecordsCreated(ctx context.Context, connectionId string, input []*model.OnCreatedObjectInput) (interface{}, error) {
	var mutate struct {
		OnObjectRecordsCreated interface{} `graphql:"(connectionId: &connectionId,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"input":        input,
	}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) OnObjectRecordsUpdated(ctx context.Context, connectionId string, input []*model.OnUpdatedObjectInput) (interface{}, error) {
	var mutate struct {
		OnObjectRecordsUpdated interface{} `graphql:"(connectionId: &connectionId,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"input":        input,
	}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) PollForChange(ctx context.Context, connectionId string) (interface{}, error) {
	var mutate struct {
		PollForChange interface{} `graphql:"(connectionId: &connectionId)"`
	}
	variables := map[string]interface{}{"connectionId": connectionId}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) KeepTokenAlive(ctx context.Context, connectionId string) (interface{}, error) {
	var mutate struct {
		KeepTokenAlive interface{} `graphql:"(connectionId: &connectionId)"`
	}
	variables := map[string]interface{}{"connectionId": connectionId}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) RunIntegrationTests(ctx context.Context) (interface{}, error) {
	var mutate struct {
		RunIntegrationTests interface{} `graphql:"runIntegrationTests"`
	}
	variables := map[string]interface{}{}

	resp, err := c.Client.QueryRaw(ctx, &mutate, variables)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) Objects(ctx context.Context, connectionId string) ([]*model.Object, error) {
	var query struct {
		Objects struct {
			Object []model.Object
		} `graphql:"(connectionId: &connectionId)"`
	}
	variables := map[string]interface{}{"connectionId": connectionId}

	resp, err := c.Client.QueryRaw(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	var res []*model.Object
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return res, nil
}

func (c *Client) ObjectRecord(ctx context.Context, connectionId string, name string, id string) (*model.ObjectRecord, error) {
	var query struct {
		ObjectRecord struct {
			ObjectRecord model.ObjectRecord
		} `graphql:"(connectionId: &connectionId,name: &name,id: &id)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"id":           id,
		"name":         name,
	}

	resp, err := c.Client.QueryRaw(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	var res model.ObjectRecord
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) Soql(ctx context.Context, connectionId string, q string, input model.PageInput) (*model.QueryResult, error) {
	var query struct {
		Soql struct {
			QueryResult model.QueryResult
		} `graphql:"(connectionId: &connectionId,q: &q,input: &input)"`
	}
	variables := map[string]interface{}{
		"connectionId": connectionId,
		"input":        input,
		"q":            q,
	}

	resp, err := c.Client.QueryRaw(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	var res model.QueryResult
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}

func (c *Client) Health(ctx context.Context) (*model.ConnectorHealth, error) {
	var query struct {
		Health struct {
			ConnectorHealth model.ConnectorHealth
		} `graphql:"health"`
	}
	variables := map[string]interface{}{}

	resp, err := c.Client.QueryRaw(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	var res model.ConnectorHealth
	if resp != nil {
		byteData, err := resp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if unMarshalErr := json.Unmarshal(byteData, &res); unMarshalErr != nil {
			return nil, unMarshalErr
		}
	}
	return &res, nil
}
