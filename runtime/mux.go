package runtime

import (
	"errors"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ysugimoto/grpc-graphql-gateway/middleware"
)

type GraphqlHandler func(w http.ResponseWriter, r *http.Request) *graphql.Result
type GraphqlErrorHandler func(errs gqlerrors.FormattedErrors)

func Respond(w http.ResponseWriter, status int, message string) {
	m := []byte(message)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(m)))
	w.WriteHeader(status)
	if len(m) > 0 {
		w.Write(m)
	}
}

type ServeMux struct {
	middlewares  []middleware.MiddlewareFunc
	ErrorHandler func(errs []gqlerrors.FormattedError) error
	schema       graphql.Schema
	Handler      GraphqlHandler
}

func NewServeMux(ms ...middleware.MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,
	}
}

func (s *ServeMux) Use(ms ...middleware.MiddlewareFunc) *ServeMux {
	s.middlewares = append(s.middlewares, ms...)
	return s
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.middlewares {
		if err := m(w, r); err != nil {
			Respond(w, http.StatusBadRequest, "middleware error occured: "+err.Error())
			return
		}
	}

	if s.Handler == nil {
		Respond(w, http.StatusBadRequest, "graphql handler is not registered")
		return
	}

	result := s.Handler(w, r)
	if result == nil {
		return
	}

	if len(result.Errors) > 0 {
		if s.ErrorHandler != nil {
			if err := s.ErrorHandler(result.Errors); err != nil {
				Respond(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	out, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func ParseRequest(r *http.Request) (
	query string,
	variables map[string]interface{},
	parseError error,
) {
	switch r.Method {
	case http.MethodPost:
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			parseError = errors.New("malformed request body, " + err.Error())
			return
		}
		query = string(buf)
	case http.MethodGet:
		query = r.URL.Query().Get("query")
	default:
		parseError = errors.New("invalid request method: '" + r.Method + "'")
	}
	return
}