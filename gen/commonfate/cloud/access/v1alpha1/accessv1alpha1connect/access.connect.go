// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: commonfate/cloud/access/v1alpha1/access.proto

package accessv1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/common-fate/ciem/gen/commonfate/cloud/access/v1alpha1"
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
	// UserManagementServiceName is the fully-qualified name of the UserManagementService service.
	UserManagementServiceName = "commonfate.cloud.access.v1alpha1.UserManagementService"
	// ResourceServiceName is the fully-qualified name of the ResourceService service.
	ResourceServiceName = "commonfate.cloud.access.v1alpha1.ResourceService"
	// AccessServiceName is the fully-qualified name of the AccessService service.
	AccessServiceName = "commonfate.cloud.access.v1alpha1.AccessService"
	// ControlPlaneServiceName is the fully-qualified name of the ControlPlaneService service.
	ControlPlaneServiceName = "commonfate.cloud.access.v1alpha1.ControlPlaneService"
	// AccessRequestServiceName is the fully-qualified name of the AccessRequestService service.
	AccessRequestServiceName = "commonfate.cloud.access.v1alpha1.AccessRequestService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// UserManagementServiceListUsersProcedure is the fully-qualified name of the
	// UserManagementService's ListUsers RPC.
	UserManagementServiceListUsersProcedure = "/commonfate.cloud.access.v1alpha1.UserManagementService/ListUsers"
	// ResourceServiceListResourcesForProviderProcedure is the fully-qualified name of the
	// ResourceService's ListResourcesForProvider RPC.
	ResourceServiceListResourcesForProviderProcedure = "/commonfate.cloud.access.v1alpha1.ResourceService/ListResourcesForProvider"
	// AccessServiceGrantProcedure is the fully-qualified name of the AccessService's Grant RPC.
	AccessServiceGrantProcedure = "/commonfate.cloud.access.v1alpha1.AccessService/Grant"
	// ControlPlaneServiceGetExistingAccessRequestProcedure is the fully-qualified name of the
	// ControlPlaneService's GetExistingAccessRequest RPC.
	ControlPlaneServiceGetExistingAccessRequestProcedure = "/commonfate.cloud.access.v1alpha1.ControlPlaneService/GetExistingAccessRequest"
	// AccessRequestServiceListAccessRequestsProcedure is the fully-qualified name of the
	// AccessRequestService's ListAccessRequests RPC.
	AccessRequestServiceListAccessRequestsProcedure = "/commonfate.cloud.access.v1alpha1.AccessRequestService/ListAccessRequests"
	// AccessRequestServiceGetAccessRequestProcedure is the fully-qualified name of the
	// AccessRequestService's GetAccessRequest RPC.
	AccessRequestServiceGetAccessRequestProcedure = "/commonfate.cloud.access.v1alpha1.AccessRequestService/GetAccessRequest"
	// AccessRequestServiceRevokeAccessRequestProcedure is the fully-qualified name of the
	// AccessRequestService's RevokeAccessRequest RPC.
	AccessRequestServiceRevokeAccessRequestProcedure = "/commonfate.cloud.access.v1alpha1.AccessRequestService/RevokeAccessRequest"
	// AccessRequestServiceCancelAccessRequestProcedure is the fully-qualified name of the
	// AccessRequestService's CancelAccessRequest RPC.
	AccessRequestServiceCancelAccessRequestProcedure = "/commonfate.cloud.access.v1alpha1.AccessRequestService/CancelAccessRequest"
	// AccessRequestServiceReviewAccessRequestProcedure is the fully-qualified name of the
	// AccessRequestService's ReviewAccessRequest RPC.
	AccessRequestServiceReviewAccessRequestProcedure = "/commonfate.cloud.access.v1alpha1.AccessRequestService/ReviewAccessRequest"
)

// UserManagementServiceClient is a client for the
// commonfate.cloud.access.v1alpha1.UserManagementService service.
type UserManagementServiceClient interface {
	ListUsers(context.Context, *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error)
}

// NewUserManagementServiceClient constructs a client for the
// commonfate.cloud.access.v1alpha1.UserManagementService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewUserManagementServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) UserManagementServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &userManagementServiceClient{
		listUsers: connect_go.NewClient[v1alpha1.ListUsersRequest, v1alpha1.ListUsersResponse](
			httpClient,
			baseURL+UserManagementServiceListUsersProcedure,
			opts...,
		),
	}
}

// userManagementServiceClient implements UserManagementServiceClient.
type userManagementServiceClient struct {
	listUsers *connect_go.Client[v1alpha1.ListUsersRequest, v1alpha1.ListUsersResponse]
}

// ListUsers calls commonfate.cloud.access.v1alpha1.UserManagementService.ListUsers.
func (c *userManagementServiceClient) ListUsers(ctx context.Context, req *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error) {
	return c.listUsers.CallUnary(ctx, req)
}

// UserManagementServiceHandler is an implementation of the
// commonfate.cloud.access.v1alpha1.UserManagementService service.
type UserManagementServiceHandler interface {
	ListUsers(context.Context, *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error)
}

// NewUserManagementServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewUserManagementServiceHandler(svc UserManagementServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	userManagementServiceListUsersHandler := connect_go.NewUnaryHandler(
		UserManagementServiceListUsersProcedure,
		svc.ListUsers,
		opts...,
	)
	return "/commonfate.cloud.access.v1alpha1.UserManagementService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case UserManagementServiceListUsersProcedure:
			userManagementServiceListUsersHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedUserManagementServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedUserManagementServiceHandler struct{}

func (UnimplementedUserManagementServiceHandler) ListUsers(context.Context, *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.UserManagementService.ListUsers is not implemented"))
}

// ResourceServiceClient is a client for the commonfate.cloud.access.v1alpha1.ResourceService
// service.
type ResourceServiceClient interface {
	ListResourcesForProvider(context.Context, *connect_go.Request[v1alpha1.ListResourcesForProviderRequest]) (*connect_go.Response[v1alpha1.ListResourcesForProviderResponse], error)
}

// NewResourceServiceClient constructs a client for the
// commonfate.cloud.access.v1alpha1.ResourceService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewResourceServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ResourceServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &resourceServiceClient{
		listResourcesForProvider: connect_go.NewClient[v1alpha1.ListResourcesForProviderRequest, v1alpha1.ListResourcesForProviderResponse](
			httpClient,
			baseURL+ResourceServiceListResourcesForProviderProcedure,
			opts...,
		),
	}
}

// resourceServiceClient implements ResourceServiceClient.
type resourceServiceClient struct {
	listResourcesForProvider *connect_go.Client[v1alpha1.ListResourcesForProviderRequest, v1alpha1.ListResourcesForProviderResponse]
}

// ListResourcesForProvider calls
// commonfate.cloud.access.v1alpha1.ResourceService.ListResourcesForProvider.
func (c *resourceServiceClient) ListResourcesForProvider(ctx context.Context, req *connect_go.Request[v1alpha1.ListResourcesForProviderRequest]) (*connect_go.Response[v1alpha1.ListResourcesForProviderResponse], error) {
	return c.listResourcesForProvider.CallUnary(ctx, req)
}

// ResourceServiceHandler is an implementation of the
// commonfate.cloud.access.v1alpha1.ResourceService service.
type ResourceServiceHandler interface {
	ListResourcesForProvider(context.Context, *connect_go.Request[v1alpha1.ListResourcesForProviderRequest]) (*connect_go.Response[v1alpha1.ListResourcesForProviderResponse], error)
}

// NewResourceServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewResourceServiceHandler(svc ResourceServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	resourceServiceListResourcesForProviderHandler := connect_go.NewUnaryHandler(
		ResourceServiceListResourcesForProviderProcedure,
		svc.ListResourcesForProvider,
		opts...,
	)
	return "/commonfate.cloud.access.v1alpha1.ResourceService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ResourceServiceListResourcesForProviderProcedure:
			resourceServiceListResourcesForProviderHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedResourceServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedResourceServiceHandler struct{}

func (UnimplementedResourceServiceHandler) ListResourcesForProvider(context.Context, *connect_go.Request[v1alpha1.ListResourcesForProviderRequest]) (*connect_go.Response[v1alpha1.ListResourcesForProviderResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.ResourceService.ListResourcesForProvider is not implemented"))
}

// AccessServiceClient is a client for the commonfate.cloud.access.v1alpha1.AccessService service.
type AccessServiceClient interface {
	// Grant is a high-level declarative API which can be called to ensure access has been provisioned to an entitlement.
	//
	// The method checks whether the entitlement has been provisioned to the user.
	// If the entitlement has not been provisioned, an Access Request will be created for the entitlement.
	// If a pending Access Request exists for the entitlement, this request is returned.
	//
	// In future, this method may trigger an extension to any Access Requests which are due to expire.
	//
	// This method is used by the Common Fate CLI in commands like 'cf exec gcp -- <command>' to ensure access
	// is provisioned prior to running a command.
	Grant(context.Context, *connect_go.Request[v1alpha1.GrantRequest]) (*connect_go.Response[v1alpha1.GrantResponse], error)
}

// NewAccessServiceClient constructs a client for the commonfate.cloud.access.v1alpha1.AccessService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccessServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AccessServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accessServiceClient{
		grant: connect_go.NewClient[v1alpha1.GrantRequest, v1alpha1.GrantResponse](
			httpClient,
			baseURL+AccessServiceGrantProcedure,
			opts...,
		),
	}
}

// accessServiceClient implements AccessServiceClient.
type accessServiceClient struct {
	grant *connect_go.Client[v1alpha1.GrantRequest, v1alpha1.GrantResponse]
}

// Grant calls commonfate.cloud.access.v1alpha1.AccessService.Grant.
func (c *accessServiceClient) Grant(ctx context.Context, req *connect_go.Request[v1alpha1.GrantRequest]) (*connect_go.Response[v1alpha1.GrantResponse], error) {
	return c.grant.CallUnary(ctx, req)
}

// AccessServiceHandler is an implementation of the commonfate.cloud.access.v1alpha1.AccessService
// service.
type AccessServiceHandler interface {
	// Grant is a high-level declarative API which can be called to ensure access has been provisioned to an entitlement.
	//
	// The method checks whether the entitlement has been provisioned to the user.
	// If the entitlement has not been provisioned, an Access Request will be created for the entitlement.
	// If a pending Access Request exists for the entitlement, this request is returned.
	//
	// In future, this method may trigger an extension to any Access Requests which are due to expire.
	//
	// This method is used by the Common Fate CLI in commands like 'cf exec gcp -- <command>' to ensure access
	// is provisioned prior to running a command.
	Grant(context.Context, *connect_go.Request[v1alpha1.GrantRequest]) (*connect_go.Response[v1alpha1.GrantResponse], error)
}

// NewAccessServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccessServiceHandler(svc AccessServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	accessServiceGrantHandler := connect_go.NewUnaryHandler(
		AccessServiceGrantProcedure,
		svc.Grant,
		opts...,
	)
	return "/commonfate.cloud.access.v1alpha1.AccessService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AccessServiceGrantProcedure:
			accessServiceGrantHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAccessServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccessServiceHandler struct{}

func (UnimplementedAccessServiceHandler) Grant(context.Context, *connect_go.Request[v1alpha1.GrantRequest]) (*connect_go.Response[v1alpha1.GrantResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessService.Grant is not implemented"))
}

// ControlPlaneServiceClient is a client for the
// commonfate.cloud.access.v1alpha1.ControlPlaneService service.
type ControlPlaneServiceClient interface {
	// GetExistingAccessRequest checks if there is an existing access request for the particular entitlement.
	// It returns a nil response if no access request exists.
	GetExistingAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetExistingAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetExistingAccessRequestResponse], error)
}

// NewControlPlaneServiceClient constructs a client for the
// commonfate.cloud.access.v1alpha1.ControlPlaneService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewControlPlaneServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ControlPlaneServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &controlPlaneServiceClient{
		getExistingAccessRequest: connect_go.NewClient[v1alpha1.GetExistingAccessRequestRequest, v1alpha1.GetExistingAccessRequestResponse](
			httpClient,
			baseURL+ControlPlaneServiceGetExistingAccessRequestProcedure,
			opts...,
		),
	}
}

// controlPlaneServiceClient implements ControlPlaneServiceClient.
type controlPlaneServiceClient struct {
	getExistingAccessRequest *connect_go.Client[v1alpha1.GetExistingAccessRequestRequest, v1alpha1.GetExistingAccessRequestResponse]
}

// GetExistingAccessRequest calls
// commonfate.cloud.access.v1alpha1.ControlPlaneService.GetExistingAccessRequest.
func (c *controlPlaneServiceClient) GetExistingAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.GetExistingAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetExistingAccessRequestResponse], error) {
	return c.getExistingAccessRequest.CallUnary(ctx, req)
}

// ControlPlaneServiceHandler is an implementation of the
// commonfate.cloud.access.v1alpha1.ControlPlaneService service.
type ControlPlaneServiceHandler interface {
	// GetExistingAccessRequest checks if there is an existing access request for the particular entitlement.
	// It returns a nil response if no access request exists.
	GetExistingAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetExistingAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetExistingAccessRequestResponse], error)
}

// NewControlPlaneServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewControlPlaneServiceHandler(svc ControlPlaneServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	controlPlaneServiceGetExistingAccessRequestHandler := connect_go.NewUnaryHandler(
		ControlPlaneServiceGetExistingAccessRequestProcedure,
		svc.GetExistingAccessRequest,
		opts...,
	)
	return "/commonfate.cloud.access.v1alpha1.ControlPlaneService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ControlPlaneServiceGetExistingAccessRequestProcedure:
			controlPlaneServiceGetExistingAccessRequestHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedControlPlaneServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedControlPlaneServiceHandler struct{}

func (UnimplementedControlPlaneServiceHandler) GetExistingAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetExistingAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetExistingAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.ControlPlaneService.GetExistingAccessRequest is not implemented"))
}

// AccessRequestServiceClient is a client for the
// commonfate.cloud.access.v1alpha1.AccessRequestService service.
type AccessRequestServiceClient interface {
	ListAccessRequests(context.Context, *connect_go.Request[v1alpha1.ListAccessRequestsRequest]) (*connect_go.Response[v1alpha1.ListAccessRequestsResponse], error)
	GetAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetAccessRequestResponse], error)
	RevokeAccessRequest(context.Context, *connect_go.Request[v1alpha1.RevokeAccessRequestRequest]) (*connect_go.Response[v1alpha1.RevokeAccessRequestResponse], error)
	CancelAccessRequest(context.Context, *connect_go.Request[v1alpha1.CancelAccessRequestRequest]) (*connect_go.Response[v1alpha1.CancelAccessRequestResponse], error)
	ReviewAccessRequest(context.Context, *connect_go.Request[v1alpha1.ReviewAccessRequestRequest]) (*connect_go.Response[v1alpha1.ReviewAccessRequestResponse], error)
}

// NewAccessRequestServiceClient constructs a client for the
// commonfate.cloud.access.v1alpha1.AccessRequestService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccessRequestServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AccessRequestServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accessRequestServiceClient{
		listAccessRequests: connect_go.NewClient[v1alpha1.ListAccessRequestsRequest, v1alpha1.ListAccessRequestsResponse](
			httpClient,
			baseURL+AccessRequestServiceListAccessRequestsProcedure,
			opts...,
		),
		getAccessRequest: connect_go.NewClient[v1alpha1.GetAccessRequestRequest, v1alpha1.GetAccessRequestResponse](
			httpClient,
			baseURL+AccessRequestServiceGetAccessRequestProcedure,
			opts...,
		),
		revokeAccessRequest: connect_go.NewClient[v1alpha1.RevokeAccessRequestRequest, v1alpha1.RevokeAccessRequestResponse](
			httpClient,
			baseURL+AccessRequestServiceRevokeAccessRequestProcedure,
			opts...,
		),
		cancelAccessRequest: connect_go.NewClient[v1alpha1.CancelAccessRequestRequest, v1alpha1.CancelAccessRequestResponse](
			httpClient,
			baseURL+AccessRequestServiceCancelAccessRequestProcedure,
			opts...,
		),
		reviewAccessRequest: connect_go.NewClient[v1alpha1.ReviewAccessRequestRequest, v1alpha1.ReviewAccessRequestResponse](
			httpClient,
			baseURL+AccessRequestServiceReviewAccessRequestProcedure,
			opts...,
		),
	}
}

// accessRequestServiceClient implements AccessRequestServiceClient.
type accessRequestServiceClient struct {
	listAccessRequests  *connect_go.Client[v1alpha1.ListAccessRequestsRequest, v1alpha1.ListAccessRequestsResponse]
	getAccessRequest    *connect_go.Client[v1alpha1.GetAccessRequestRequest, v1alpha1.GetAccessRequestResponse]
	revokeAccessRequest *connect_go.Client[v1alpha1.RevokeAccessRequestRequest, v1alpha1.RevokeAccessRequestResponse]
	cancelAccessRequest *connect_go.Client[v1alpha1.CancelAccessRequestRequest, v1alpha1.CancelAccessRequestResponse]
	reviewAccessRequest *connect_go.Client[v1alpha1.ReviewAccessRequestRequest, v1alpha1.ReviewAccessRequestResponse]
}

// ListAccessRequests calls
// commonfate.cloud.access.v1alpha1.AccessRequestService.ListAccessRequests.
func (c *accessRequestServiceClient) ListAccessRequests(ctx context.Context, req *connect_go.Request[v1alpha1.ListAccessRequestsRequest]) (*connect_go.Response[v1alpha1.ListAccessRequestsResponse], error) {
	return c.listAccessRequests.CallUnary(ctx, req)
}

// GetAccessRequest calls commonfate.cloud.access.v1alpha1.AccessRequestService.GetAccessRequest.
func (c *accessRequestServiceClient) GetAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.GetAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetAccessRequestResponse], error) {
	return c.getAccessRequest.CallUnary(ctx, req)
}

// RevokeAccessRequest calls
// commonfate.cloud.access.v1alpha1.AccessRequestService.RevokeAccessRequest.
func (c *accessRequestServiceClient) RevokeAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.RevokeAccessRequestRequest]) (*connect_go.Response[v1alpha1.RevokeAccessRequestResponse], error) {
	return c.revokeAccessRequest.CallUnary(ctx, req)
}

// CancelAccessRequest calls
// commonfate.cloud.access.v1alpha1.AccessRequestService.CancelAccessRequest.
func (c *accessRequestServiceClient) CancelAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.CancelAccessRequestRequest]) (*connect_go.Response[v1alpha1.CancelAccessRequestResponse], error) {
	return c.cancelAccessRequest.CallUnary(ctx, req)
}

// ReviewAccessRequest calls
// commonfate.cloud.access.v1alpha1.AccessRequestService.ReviewAccessRequest.
func (c *accessRequestServiceClient) ReviewAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.ReviewAccessRequestRequest]) (*connect_go.Response[v1alpha1.ReviewAccessRequestResponse], error) {
	return c.reviewAccessRequest.CallUnary(ctx, req)
}

// AccessRequestServiceHandler is an implementation of the
// commonfate.cloud.access.v1alpha1.AccessRequestService service.
type AccessRequestServiceHandler interface {
	ListAccessRequests(context.Context, *connect_go.Request[v1alpha1.ListAccessRequestsRequest]) (*connect_go.Response[v1alpha1.ListAccessRequestsResponse], error)
	GetAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetAccessRequestResponse], error)
	RevokeAccessRequest(context.Context, *connect_go.Request[v1alpha1.RevokeAccessRequestRequest]) (*connect_go.Response[v1alpha1.RevokeAccessRequestResponse], error)
	CancelAccessRequest(context.Context, *connect_go.Request[v1alpha1.CancelAccessRequestRequest]) (*connect_go.Response[v1alpha1.CancelAccessRequestResponse], error)
	ReviewAccessRequest(context.Context, *connect_go.Request[v1alpha1.ReviewAccessRequestRequest]) (*connect_go.Response[v1alpha1.ReviewAccessRequestResponse], error)
}

// NewAccessRequestServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccessRequestServiceHandler(svc AccessRequestServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	accessRequestServiceListAccessRequestsHandler := connect_go.NewUnaryHandler(
		AccessRequestServiceListAccessRequestsProcedure,
		svc.ListAccessRequests,
		opts...,
	)
	accessRequestServiceGetAccessRequestHandler := connect_go.NewUnaryHandler(
		AccessRequestServiceGetAccessRequestProcedure,
		svc.GetAccessRequest,
		opts...,
	)
	accessRequestServiceRevokeAccessRequestHandler := connect_go.NewUnaryHandler(
		AccessRequestServiceRevokeAccessRequestProcedure,
		svc.RevokeAccessRequest,
		opts...,
	)
	accessRequestServiceCancelAccessRequestHandler := connect_go.NewUnaryHandler(
		AccessRequestServiceCancelAccessRequestProcedure,
		svc.CancelAccessRequest,
		opts...,
	)
	accessRequestServiceReviewAccessRequestHandler := connect_go.NewUnaryHandler(
		AccessRequestServiceReviewAccessRequestProcedure,
		svc.ReviewAccessRequest,
		opts...,
	)
	return "/commonfate.cloud.access.v1alpha1.AccessRequestService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AccessRequestServiceListAccessRequestsProcedure:
			accessRequestServiceListAccessRequestsHandler.ServeHTTP(w, r)
		case AccessRequestServiceGetAccessRequestProcedure:
			accessRequestServiceGetAccessRequestHandler.ServeHTTP(w, r)
		case AccessRequestServiceRevokeAccessRequestProcedure:
			accessRequestServiceRevokeAccessRequestHandler.ServeHTTP(w, r)
		case AccessRequestServiceCancelAccessRequestProcedure:
			accessRequestServiceCancelAccessRequestHandler.ServeHTTP(w, r)
		case AccessRequestServiceReviewAccessRequestProcedure:
			accessRequestServiceReviewAccessRequestHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAccessRequestServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccessRequestServiceHandler struct{}

func (UnimplementedAccessRequestServiceHandler) ListAccessRequests(context.Context, *connect_go.Request[v1alpha1.ListAccessRequestsRequest]) (*connect_go.Response[v1alpha1.ListAccessRequestsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessRequestService.ListAccessRequests is not implemented"))
}

func (UnimplementedAccessRequestServiceHandler) GetAccessRequest(context.Context, *connect_go.Request[v1alpha1.GetAccessRequestRequest]) (*connect_go.Response[v1alpha1.GetAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessRequestService.GetAccessRequest is not implemented"))
}

func (UnimplementedAccessRequestServiceHandler) RevokeAccessRequest(context.Context, *connect_go.Request[v1alpha1.RevokeAccessRequestRequest]) (*connect_go.Response[v1alpha1.RevokeAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessRequestService.RevokeAccessRequest is not implemented"))
}

func (UnimplementedAccessRequestServiceHandler) CancelAccessRequest(context.Context, *connect_go.Request[v1alpha1.CancelAccessRequestRequest]) (*connect_go.Response[v1alpha1.CancelAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessRequestService.CancelAccessRequest is not implemented"))
}

func (UnimplementedAccessRequestServiceHandler) ReviewAccessRequest(context.Context, *connect_go.Request[v1alpha1.ReviewAccessRequestRequest]) (*connect_go.Response[v1alpha1.ReviewAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("commonfate.cloud.access.v1alpha1.AccessRequestService.ReviewAccessRequest is not implemented"))
}
