// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: commonfate/authz/v1alpha1/authz.proto

package authzv1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// AuthzServiceName is the fully-qualified name of the AuthzService service.
	AuthzServiceName = "commonfate.authz.v1alpha1.AuthzService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AuthzServiceBatchPutEntityProcedure is the fully-qualified name of the AuthzService's
	// BatchPutEntity RPC.
	AuthzServiceBatchPutEntityProcedure = "/commonfate.authz.v1alpha1.AuthzService/BatchPutEntity"
	// AuthzServiceBatchDeleteEntityProcedure is the fully-qualified name of the AuthzService's
	// BatchDeleteEntity RPC.
	AuthzServiceBatchDeleteEntityProcedure = "/commonfate.authz.v1alpha1.AuthzService/BatchDeleteEntity"
	// AuthzServiceBatchPutPolicyProcedure is the fully-qualified name of the AuthzService's
	// BatchPutPolicy RPC.
	AuthzServiceBatchPutPolicyProcedure = "/commonfate.authz.v1alpha1.AuthzService/BatchPutPolicy"
	// AuthzServiceBatchAuthorizeProcedure is the fully-qualified name of the AuthzService's
	// BatchAuthorize RPC.
	AuthzServiceBatchAuthorizeProcedure = "/commonfate.authz.v1alpha1.AuthzService/BatchAuthorize"
	// AuthzServiceLookupResourcesProcedure is the fully-qualified name of the AuthzService's
	// LookupResources RPC.
	AuthzServiceLookupResourcesProcedure = "/commonfate.authz.v1alpha1.AuthzService/LookupResources"
	// AuthzServiceListPoliciesProcedure is the fully-qualified name of the AuthzService's ListPolicies
	// RPC.
	AuthzServiceListPoliciesProcedure = "/commonfate.authz.v1alpha1.AuthzService/ListPolicies"
	// AuthzServiceFilterEntitiesProcedure is the fully-qualified name of the AuthzService's
	// FilterEntities RPC.
	AuthzServiceFilterEntitiesProcedure = "/commonfate.authz.v1alpha1.AuthzService/FilterEntities"
	// AuthzServiceGetEntityProcedure is the fully-qualified name of the AuthzService's GetEntity RPC.
	AuthzServiceGetEntityProcedure = "/commonfate.authz.v1alpha1.AuthzService/GetEntity"
)

// AuthzServiceClient is a client for the commonfate.authz.v1alpha1.AuthzService service.
type AuthzServiceClient interface {
	// creates or updates entities for a particular policy store in the authorization service.
	BatchPutEntity(context.Context, *connect_go.Request[v1alpha1.BatchPutEntityRequest]) (*connect_go.Response[v1alpha1.BatchPutEntityResponse], error)
	// removes entities from the authorization service.
	BatchDeleteEntity(context.Context, *connect_go.Request[v1alpha1.BatchDeleteEntityRequest]) (*connect_go.Response[v1alpha1.BatchDeleteEntityResponse], error)
	// adds Cedar policies for a particular policy store
	BatchPutPolicy(context.Context, *connect_go.Request[v1alpha1.BatchPutPolicyRequest]) (*connect_go.Response[v1alpha1.BatchPutPolicyResponse], error)
	// run multiple authorization decisions and returns allow + deny for each.
	BatchAuthorize(context.Context, *connect_go.Request[v1alpha1.BatchAuthorizeRequest]) (*connect_go.Response[v1alpha1.BatchAuthorizeResponse], error)
	// look up which resources a particular principal can access
	LookupResources(context.Context, *connect_go.Request[v1alpha1.LookupResourcesRequest]) (*connect_go.Response[v1alpha1.LookupResourcesResponse], error)
	ListPolicies(context.Context, *connect_go.Request[v1alpha1.ListPoliciesRequest]) (*connect_go.Response[v1alpha1.ListPoliciesResponse], error)
	// Query for entities matching filter conditions.
	FilterEntities(context.Context, *connect_go.Request[v1alpha1.FilterEntitiesRequest]) (*connect_go.Response[v1alpha1.FilterEntitiesResponse], error)
	// Query for entity by UID.
	GetEntity(context.Context, *connect_go.Request[v1alpha1.GetEntityRequest]) (*connect_go.Response[v1alpha1.GetEntityResponse], error)
}

// NewAuthzServiceClient constructs a client for the commonfate.authz.v1alpha1.AuthzService service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAuthzServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AuthzServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &authzServiceClient{
		batchPutEntity: connect_go.NewClient[v1alpha1.BatchPutEntityRequest, v1alpha1.BatchPutEntityResponse](
			httpClient,
			baseURL+AuthzServiceBatchPutEntityProcedure,
			opts...,
		),
		batchDeleteEntity: connect_go.NewClient[v1alpha1.BatchDeleteEntityRequest, v1alpha1.BatchDeleteEntityResponse](
			httpClient,
			baseURL+AuthzServiceBatchDeleteEntityProcedure,
			opts...,
		),
		batchPutPolicy: connect_go.NewClient[v1alpha1.BatchPutPolicyRequest, v1alpha1.BatchPutPolicyResponse](
			httpClient,
			baseURL+AuthzServiceBatchPutPolicyProcedure,
			opts...,
		),
		batchAuthorize: connect_go.NewClient[v1alpha1.BatchAuthorizeRequest, v1alpha1.BatchAuthorizeResponse](
			httpClient,
			baseURL+AuthzServiceBatchAuthorizeProcedure,
			opts...,
		),
		lookupResources: connect_go.NewClient[v1alpha1.LookupResourcesRequest, v1alpha1.LookupResourcesResponse](
			httpClient,
			baseURL+AuthzServiceLookupResourcesProcedure,
			opts...,
		),
		listPolicies: connect_go.NewClient[v1alpha1.ListPoliciesRequest, v1alpha1.ListPoliciesResponse](
			httpClient,
			baseURL+AuthzServiceListPoliciesProcedure,
			opts...,
		),
		filterEntities: connect_go.NewClient[v1alpha1.FilterEntitiesRequest, v1alpha1.FilterEntitiesResponse](
			httpClient,
			baseURL+AuthzServiceFilterEntitiesProcedure,
			opts...,
		),
		getEntity: connect_go.NewClient[v1alpha1.GetEntityRequest, v1alpha1.GetEntityResponse](
			httpClient,
			baseURL+AuthzServiceGetEntityProcedure,
			opts...,
		),
	}
}

// authzServiceClient implements AuthzServiceClient.
type authzServiceClient struct {
	batchPutEntity    *connect_go.Client[v1alpha1.BatchPutEntityRequest, v1alpha1.BatchPutEntityResponse]
	batchDeleteEntity *connect_go.Client[v1alpha1.BatchDeleteEntityRequest, v1alpha1.BatchDeleteEntityResponse]
	batchPutPolicy    *connect_go.Client[v1alpha1.BatchPutPolicyRequest, v1alpha1.BatchPutPolicyResponse]
	batchAuthorize    *connect_go.Client[v1alpha1.BatchAuthorizeRequest, v1alpha1.BatchAuthorizeResponse]
	lookupResources   *connect_go.Client[v1alpha1.LookupResourcesRequest, v1alpha1.LookupResourcesResponse]
	listPolicies      *connect_go.Client[v1alpha1.ListPoliciesRequest, v1alpha1.ListPoliciesResponse]
	filterEntities    *connect_go.Client[v1alpha1.FilterEntitiesRequest, v1alpha1.FilterEntitiesResponse]
	getEntity         *connect_go.Client[v1alpha1.GetEntityRequest, v1alpha1.GetEntityResponse]
}

// BatchPutEntity calls commonfate.authz.v1alpha1.AuthzService.BatchPutEntity.
func (c *authzServiceClient) BatchPutEntity(ctx context.Context, req *connect_go.Request[v1alpha1.BatchPutEntityRequest]) (*connect_go.Response[v1alpha1.BatchPutEntityResponse], error) {
	return c.batchPutEntity.CallUnary(ctx, req)
}

// BatchDeleteEntity calls commonfate.authz.v1alpha1.AuthzService.BatchDeleteEntity.
func (c *authzServiceClient) BatchDeleteEntity(ctx context.Context, req *connect_go.Request[v1alpha1.BatchDeleteEntityRequest]) (*connect_go.Response[v1alpha1.BatchDeleteEntityResponse], error) {
	return c.batchDeleteEntity.CallUnary(ctx, req)
}

// BatchPutPolicy calls commonfate.authz.v1alpha1.AuthzService.BatchPutPolicy.
func (c *authzServiceClient) BatchPutPolicy(ctx context.Context, req *connect_go.Request[v1alpha1.BatchPutPolicyRequest]) (*connect_go.Response[v1alpha1.BatchPutPolicyResponse], error) {
	return c.batchPutPolicy.CallUnary(ctx, req)
}

// BatchAuthorize calls commonfate.authz.v1alpha1.AuthzService.BatchAuthorize.
func (c *authzServiceClient) BatchAuthorize(ctx context.Context, req *connect_go.Request[v1alpha1.BatchAuthorizeRequest]) (*connect_go.Response[v1alpha1.BatchAuthorizeResponse], error) {
	return c.batchAuthorize.CallUnary(ctx, req)
}

// LookupResources calls commonfate.authz.v1alpha1.AuthzService.LookupResources.
func (c *authzServiceClient) LookupResources(ctx context.Context, req *connect_go.Request[v1alpha1.LookupResourcesRequest]) (*connect_go.Response[v1alpha1.LookupResourcesResponse], error) {
	return c.lookupResources.CallUnary(ctx, req)
}

// ListPolicies calls commonfate.authz.v1alpha1.AuthzService.ListPolicies.
func (c *authzServiceClient) ListPolicies(ctx context.Context, req *connect_go.Request[v1alpha1.ListPoliciesRequest]) (*connect_go.Response[v1alpha1.ListPoliciesResponse], error) {
	return c.listPolicies.CallUnary(ctx, req)
}

// FilterEntities calls commonfate.authz.v1alpha1.AuthzService.FilterEntities.
func (c *authzServiceClient) FilterEntities(ctx context.Context, req *connect_go.Request[v1alpha1.FilterEntitiesRequest]) (*connect_go.Response[v1alpha1.FilterEntitiesResponse], error) {
	return c.filterEntities.CallUnary(ctx, req)
}

// GetEntity calls commonfate.authz.v1alpha1.AuthzService.GetEntity.
func (c *authzServiceClient) GetEntity(ctx context.Context, req *connect_go.Request[v1alpha1.GetEntityRequest]) (*connect_go.Response[v1alpha1.GetEntityResponse], error) {
	return c.getEntity.CallUnary(ctx, req)
}

// AuthzServiceHandler is an implementation of the commonfate.authz.v1alpha1.AuthzService service.
type AuthzServiceHandler interface {
	// creates or updates entities for a particular policy store in the authorization service.
	BatchPutEntity(context.Context, *connect_go.Request[v1alpha1.BatchPutEntityRequest]) (*connect_go.Response[v1alpha1.BatchPutEntityResponse], error)
	// removes entities from the authorization service.
	BatchDeleteEntity(context.Context, *connect_go.Request[v1alpha1.BatchDeleteEntityRequest]) (*connect_go.Response[v1alpha1.BatchDeleteEntityResponse], error)
	// adds Cedar policies for a particular policy store
	BatchPutPolicy(context.Context, *connect_go.Request[v1alpha1.BatchPutPolicyRequest]) (*connect_go.Response[v1alpha1.BatchPutPolicyResponse], error)
	// run multiple authorization decisions and returns allow + deny for each.
	BatchAuthorize(context.Context, *connect_go.Request[v1alpha1.BatchAuthorizeRequest]) (*connect_go.Response[v1alpha1.BatchAuthorizeResponse], error)
	// look up which resources a particular principal can access
	LookupResources(context.Context, *connect_go.Request[v1alpha1.LookupResourcesRequest]) (*connect_go.Response[v1alpha1.LookupResourcesResponse], error)
	ListPolicies(context.Context, *connect_go.Request[v1alpha1.ListPoliciesRequest]) (*connect_go.Response[v1alpha1.ListPoliciesResponse], error)
	// Query for entities matching filter conditions.
	FilterEntities(context.Context, *connect_go.Request[v1alpha1.FilterEntitiesRequest]) (*connect_go.Response[v1alpha1.FilterEntitiesResponse], error)
	// Query for entity by UID.
	GetEntity(context.Context, *connect_go.Request[v1alpha1.GetEntityRequest]) (*connect_go.Response[v1alpha1.GetEntityResponse], error)
}

// NewAuthzServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAuthzServiceHandler(svc AuthzServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	authzServiceBatchPutEntityHandler := connect_go.NewUnaryHandler(
		AuthzServiceBatchPutEntityProcedure,
		svc.BatchPutEntity,
		opts...,
	)
	authzServiceBatchDeleteEntityHandler := connect_go.NewUnaryHandler(
		AuthzServiceBatchDeleteEntityProcedure,
		svc.BatchDeleteEntity,
		opts...,
	)
	authzServiceBatchPutPolicyHandler := connect_go.NewUnaryHandler(
		AuthzServiceBatchPutPolicyProcedure,
		svc.BatchPutPolicy,
		opts...,
	)
	authzServiceBatchAuthorizeHandler := connect_go.NewUnaryHandler(
		AuthzServiceBatchAuthorizeProcedure,
		svc.BatchAuthorize,
		opts...,
	)
	authzServiceLookupResourcesHandler := connect_go.NewUnaryHandler(
		AuthzServiceLookupResourcesProcedure,
		svc.LookupResources,
		opts...,
	)
	authzServiceListPoliciesHandler := connect_go.NewUnaryHandler(
		AuthzServiceListPoliciesProcedure,
		svc.ListPolicies,
		opts...,
	)
	authzServiceFilterEntitiesHandler := connect_go.NewUnaryHandler(
		AuthzServiceFilterEntitiesProcedure,
		svc.FilterEntities,
		opts...,
	)
	authzServiceGetEntityHandler := connect_go.NewUnaryHandler(
		AuthzServiceGetEntityProcedure,
		svc.GetEntity,
		opts...,
	)
	return "/commonfate.authz.v1alpha1.AuthzService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AuthzServiceBatchPutEntityProcedure:
			authzServiceBatchPutEntityHandler.ServeHTTP(w, r)
		case AuthzServiceBatchDeleteEntityProcedure:
			authzServiceBatchDeleteEntityHandler.ServeHTTP(w, r)
		case AuthzServiceBatchPutPolicyProcedure:
			authzServiceBatchPutPolicyHandler.ServeHTTP(w, r)
		case AuthzServiceBatchAuthorizeProcedure:
			authzServiceBatchAuthorizeHandler.ServeHTTP(w, r)
		case AuthzServiceLookupResourcesProcedure:
			authzServiceLookupResourcesHandler.ServeHTTP(w, r)
		case AuthzServiceListPoliciesProcedure:
			authzServiceListPoliciesHandler.ServeHTTP(w, r)
		case AuthzServiceFilterEntitiesProcedure:
			authzServiceFilterEntitiesHandler.ServeHTTP(w, r)
		case AuthzServiceGetEntityProcedure:
			authzServiceGetEntityHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAuthzServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAuthzServiceHandler struct{}

func (UnimplementedAuthzServiceHandler) BatchPutEntity(context.Context, *connect_go.Request[v1alpha1.BatchPutEntityRequest]) (*connect_go.Response[v1alpha1.BatchPutEntityResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.BatchPutEntity is not implemented"))
}

func (UnimplementedAuthzServiceHandler) BatchDeleteEntity(context.Context, *connect_go.Request[v1alpha1.BatchDeleteEntityRequest]) (*connect_go.Response[v1alpha1.BatchDeleteEntityResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.BatchDeleteEntity is not implemented"))
}

func (UnimplementedAuthzServiceHandler) BatchPutPolicy(context.Context, *connect_go.Request[v1alpha1.BatchPutPolicyRequest]) (*connect_go.Response[v1alpha1.BatchPutPolicyResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.BatchPutPolicy is not implemented"))
}

func (UnimplementedAuthzServiceHandler) BatchAuthorize(context.Context, *connect_go.Request[v1alpha1.BatchAuthorizeRequest]) (*connect_go.Response[v1alpha1.BatchAuthorizeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.BatchAuthorize is not implemented"))
}

func (UnimplementedAuthzServiceHandler) LookupResources(context.Context, *connect_go.Request[v1alpha1.LookupResourcesRequest]) (*connect_go.Response[v1alpha1.LookupResourcesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.LookupResources is not implemented"))
}

func (UnimplementedAuthzServiceHandler) ListPolicies(context.Context, *connect_go.Request[v1alpha1.ListPoliciesRequest]) (*connect_go.Response[v1alpha1.ListPoliciesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.ListPolicies is not implemented"))
}

func (UnimplementedAuthzServiceHandler) FilterEntities(context.Context, *connect_go.Request[v1alpha1.FilterEntitiesRequest]) (*connect_go.Response[v1alpha1.FilterEntitiesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.FilterEntities is not implemented"))
}

func (UnimplementedAuthzServiceHandler) GetEntity(context.Context, *connect_go.Request[v1alpha1.GetEntityRequest]) (*connect_go.Response[v1alpha1.GetEntityResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.authz.v1alpha1.AuthzService.GetEntity is not implemented"))
}