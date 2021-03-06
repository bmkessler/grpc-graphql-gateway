package runtime

import (
	"errors"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

type GraphqlRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

// ParseRequest parses graphql query and variables from each request methods
func parseRequest(r *http.Request) (*GraphqlRequest, error) {
	var body []byte

	// Get request body
	switch r.Method {
	case http.MethodPost:
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.New("malformed request body, " + err.Error())
		}
		body = buf
	case http.MethodGet:
		body = []byte(r.URL.Query().Get("query"))
	default:
		return nil, errors.New("invalid request method: '" + r.Method + "'")
	}

	// And try to parse
	var req GraphqlRequest
	if err := json.Unmarshal(body, &req); err != nil {
		// If error, the request body may come with single query line
		req.Query = string(body)
	}
	return &req, nil
}

func MarshalRequest(args map[string]interface{}, v interface{}) error {
	buf, err := json.Marshal(args)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &v)
}
