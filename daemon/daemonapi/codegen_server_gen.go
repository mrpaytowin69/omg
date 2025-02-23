// Package daemonapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package daemonapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth/token)
	PostAuthToken(w http.ResponseWriter, r *http.Request, params PostAuthTokenParams)

	// (GET /daemon/events)
	GetDaemonEvents(w http.ResponseWriter, r *http.Request, params GetDaemonEventsParams)

	// (POST /daemon/logs/control)
	PostDaemonLogsControl(w http.ResponseWriter, r *http.Request)

	// (GET /daemon/running)
	GetDaemonRunning(w http.ResponseWriter, r *http.Request)

	// (GET /daemon/status)
	GetDaemonStatus(w http.ResponseWriter, r *http.Request, params GetDaemonStatusParams)

	// (POST /daemon/stop)
	PostDaemonStop(w http.ResponseWriter, r *http.Request)

	// (POST /daemon/sub/action)
	PostDaemonSubAction(w http.ResponseWriter, r *http.Request)

	// (POST /node/clear)
	PostNodeClear(w http.ResponseWriter, r *http.Request)

	// (POST /node/monitor)
	PostNodeMonitor(w http.ResponseWriter, r *http.Request)

	// (GET /nodes/info)
	GetNodesInfo(w http.ResponseWriter, r *http.Request)

	// (POST /object/abort)
	PostObjectAbort(w http.ResponseWriter, r *http.Request)

	// (POST /object/clear)
	PostObjectClear(w http.ResponseWriter, r *http.Request)

	// (GET /object/config)
	GetObjectConfig(w http.ResponseWriter, r *http.Request, params GetObjectConfigParams)

	// (GET /object/file)
	GetObjectFile(w http.ResponseWriter, r *http.Request, params GetObjectFileParams)

	// (POST /object/monitor)
	PostObjectMonitor(w http.ResponseWriter, r *http.Request)

	// (POST /object/progress)
	PostObjectProgress(w http.ResponseWriter, r *http.Request)

	// (GET /object/selector)
	GetObjectSelector(w http.ResponseWriter, r *http.Request, params GetObjectSelectorParams)

	// (POST /object/status)
	PostObjectStatus(w http.ResponseWriter, r *http.Request)

	// (POST /object/switchTo)
	PostObjectSwitchTo(w http.ResponseWriter, r *http.Request)

	// (GET /public/openapi)
	GetSwagger(w http.ResponseWriter, r *http.Request)

	// (GET /relay/message)
	GetRelayMessage(w http.ResponseWriter, r *http.Request, params GetRelayMessageParams)

	// (POST /relay/message)
	PostRelayMessage(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// PostAuthToken operation middleware
func (siw *ServerInterfaceWrapper) PostAuthToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthTokenParams

	// ------------- Optional query parameter "role" -------------
	if paramValue := r.URL.Query().Get("role"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "role", r.URL.Query(), &params.Role)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "role", Err: err})
		return
	}

	// ------------- Optional query parameter "duration" -------------
	if paramValue := r.URL.Query().Get("duration"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "duration", r.URL.Query(), &params.Duration)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "duration", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostAuthToken(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetDaemonEvents operation middleware
func (siw *ServerInterfaceWrapper) GetDaemonEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDaemonEventsParams

	// ------------- Optional query parameter "duration" -------------
	if paramValue := r.URL.Query().Get("duration"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "duration", r.URL.Query(), &params.Duration)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "duration", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------
	if paramValue := r.URL.Query().Get("limit"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "filter" -------------
	if paramValue := r.URL.Query().Get("filter"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "filter", r.URL.Query(), &params.Filter)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDaemonEvents(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostDaemonLogsControl operation middleware
func (siw *ServerInterfaceWrapper) PostDaemonLogsControl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDaemonLogsControl(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetDaemonRunning operation middleware
func (siw *ServerInterfaceWrapper) GetDaemonRunning(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDaemonRunning(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetDaemonStatus operation middleware
func (siw *ServerInterfaceWrapper) GetDaemonStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDaemonStatusParams

	// ------------- Optional query parameter "namespace" -------------
	if paramValue := r.URL.Query().Get("namespace"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "namespace", r.URL.Query(), &params.Namespace)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	// ------------- Optional query parameter "relatives" -------------
	if paramValue := r.URL.Query().Get("relatives"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "relatives", r.URL.Query(), &params.Relatives)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "relatives", Err: err})
		return
	}

	// ------------- Optional query parameter "selector" -------------
	if paramValue := r.URL.Query().Get("selector"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "selector", r.URL.Query(), &params.Selector)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "selector", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDaemonStatus(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostDaemonStop operation middleware
func (siw *ServerInterfaceWrapper) PostDaemonStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDaemonStop(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostDaemonSubAction operation middleware
func (siw *ServerInterfaceWrapper) PostDaemonSubAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDaemonSubAction(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostNodeClear operation middleware
func (siw *ServerInterfaceWrapper) PostNodeClear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostNodeClear(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostNodeMonitor operation middleware
func (siw *ServerInterfaceWrapper) PostNodeMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostNodeMonitor(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetNodesInfo operation middleware
func (siw *ServerInterfaceWrapper) GetNodesInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNodesInfo(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectAbort operation middleware
func (siw *ServerInterfaceWrapper) PostObjectAbort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectAbort(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectClear operation middleware
func (siw *ServerInterfaceWrapper) PostObjectClear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectClear(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetObjectConfig operation middleware
func (siw *ServerInterfaceWrapper) GetObjectConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetObjectConfigParams

	// ------------- Required query parameter "path" -------------
	if paramValue := r.URL.Query().Get("path"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "path"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "path", r.URL.Query(), &params.Path)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "path", Err: err})
		return
	}

	// ------------- Optional query parameter "evaluate" -------------
	if paramValue := r.URL.Query().Get("evaluate"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "evaluate", r.URL.Query(), &params.Evaluate)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "evaluate", Err: err})
		return
	}

	// ------------- Optional query parameter "impersonate" -------------
	if paramValue := r.URL.Query().Get("impersonate"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "impersonate", r.URL.Query(), &params.Impersonate)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "impersonate", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetObjectConfig(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetObjectFile operation middleware
func (siw *ServerInterfaceWrapper) GetObjectFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetObjectFileParams

	// ------------- Required query parameter "path" -------------
	if paramValue := r.URL.Query().Get("path"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "path"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "path", r.URL.Query(), &params.Path)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "path", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetObjectFile(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectMonitor operation middleware
func (siw *ServerInterfaceWrapper) PostObjectMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectMonitor(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectProgress operation middleware
func (siw *ServerInterfaceWrapper) PostObjectProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectProgress(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetObjectSelector operation middleware
func (siw *ServerInterfaceWrapper) GetObjectSelector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetObjectSelectorParams

	// ------------- Required query parameter "selector" -------------
	if paramValue := r.URL.Query().Get("selector"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "selector"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "selector", r.URL.Query(), &params.Selector)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "selector", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetObjectSelector(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectStatus operation middleware
func (siw *ServerInterfaceWrapper) PostObjectStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectStatus(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostObjectSwitchTo operation middleware
func (siw *ServerInterfaceWrapper) PostObjectSwitchTo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostObjectSwitchTo(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetSwagger operation middleware
func (siw *ServerInterfaceWrapper) GetSwagger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSwagger(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetRelayMessage operation middleware
func (siw *ServerInterfaceWrapper) GetRelayMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetRelayMessageParams

	// ------------- Optional query parameter "nodename" -------------
	if paramValue := r.URL.Query().Get("nodename"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "nodename", r.URL.Query(), &params.Nodename)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "nodename", Err: err})
		return
	}

	// ------------- Optional query parameter "cluster_id" -------------
	if paramValue := r.URL.Query().Get("cluster_id"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cluster_id", r.URL.Query(), &params.ClusterId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "cluster_id", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRelayMessage(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// PostRelayMessage operation middleware
func (siw *ServerInterfaceWrapper) PostRelayMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{""})

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRelayMessage(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/auth/token", wrapper.PostAuthToken)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/daemon/events", wrapper.GetDaemonEvents)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/daemon/logs/control", wrapper.PostDaemonLogsControl)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/daemon/running", wrapper.GetDaemonRunning)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/daemon/status", wrapper.GetDaemonStatus)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/daemon/stop", wrapper.PostDaemonStop)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/daemon/sub/action", wrapper.PostDaemonSubAction)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/node/clear", wrapper.PostNodeClear)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/node/monitor", wrapper.PostNodeMonitor)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/nodes/info", wrapper.GetNodesInfo)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/abort", wrapper.PostObjectAbort)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/clear", wrapper.PostObjectClear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/object/config", wrapper.GetObjectConfig)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/object/file", wrapper.GetObjectFile)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/monitor", wrapper.PostObjectMonitor)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/progress", wrapper.PostObjectProgress)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/object/selector", wrapper.GetObjectSelector)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/status", wrapper.PostObjectStatus)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/object/switchTo", wrapper.PostObjectSwitchTo)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/public/openapi", wrapper.GetSwagger)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/relay/message", wrapper.GetRelayMessage)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/relay/message", wrapper.PostRelayMessage)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
