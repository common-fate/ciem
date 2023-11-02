// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/common_fate/v1alpha1/common_fate.proto

package common_fatev1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/common-fate/ciem/gen/proto/common_fate/v1alpha1"
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
	UserManagementServiceName = "common_fate.v1alpha1.UserManagementService"
	// ConnectionsServiceName is the fully-qualified name of the ConnectionsService service.
	ConnectionsServiceName = "common_fate.v1alpha1.ConnectionsService"
	// UsageMetricsServiceName is the fully-qualified name of the UsageMetricsService service.
	UsageMetricsServiceName = "common_fate.v1alpha1.UsageMetricsService"
	// AccessServiceName is the fully-qualified name of the AccessService service.
	AccessServiceName = "common_fate.v1alpha1.AccessService"
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
	UserManagementServiceListUsersProcedure = "/common_fate.v1alpha1.UserManagementService/ListUsers"
	// ConnectionsServiceListConnectionsProcedure is the fully-qualified name of the
	// ConnectionsService's ListConnections RPC.
	ConnectionsServiceListConnectionsProcedure = "/common_fate.v1alpha1.ConnectionsService/ListConnections"
	// ConnectionsServiceAddConnectionProcedure is the fully-qualified name of the ConnectionsService's
	// AddConnection RPC.
	ConnectionsServiceAddConnectionProcedure = "/common_fate.v1alpha1.ConnectionsService/AddConnection"
	// ConnectionsServiceRemoveConnectionProcedure is the fully-qualified name of the
	// ConnectionsService's RemoveConnection RPC.
	ConnectionsServiceRemoveConnectionProcedure = "/common_fate.v1alpha1.ConnectionsService/RemoveConnection"
	// ConnectionsServiceListAWSAccountsProcedure is the fully-qualified name of the
	// ConnectionsService's ListAWSAccounts RPC.
	ConnectionsServiceListAWSAccountsProcedure = "/common_fate.v1alpha1.ConnectionsService/ListAWSAccounts"
	// ConnectionsServiceGetAWSRolesForAccountProcedure is the fully-qualified name of the
	// ConnectionsService's GetAWSRolesForAccount RPC.
	ConnectionsServiceGetAWSRolesForAccountProcedure = "/common_fate.v1alpha1.ConnectionsService/GetAWSRolesForAccount"
	// UsageMetricsServiceGetAWSRoleMetricsProcedure is the fully-qualified name of the
	// UsageMetricsService's GetAWSRoleMetrics RPC.
	UsageMetricsServiceGetAWSRoleMetricsProcedure = "/common_fate.v1alpha1.UsageMetricsService/GetAWSRoleMetrics"
	// UsageMetricsServiceGetUsageForRoleProcedure is the fully-qualified name of the
	// UsageMetricsService's GetUsageForRole RPC.
	UsageMetricsServiceGetUsageForRoleProcedure = "/common_fate.v1alpha1.UsageMetricsService/GetUsageForRole"
	// AccessServiceListEntitlementsForProviderProcedure is the fully-qualified name of the
	// AccessService's ListEntitlementsForProvider RPC.
	AccessServiceListEntitlementsForProviderProcedure = "/common_fate.v1alpha1.AccessService/ListEntitlementsForProvider"
	// AccessServiceCreateAccessRequestProcedure is the fully-qualified name of the AccessService's
	// CreateAccessRequest RPC.
	AccessServiceCreateAccessRequestProcedure = "/common_fate.v1alpha1.AccessService/CreateAccessRequest"
)

// UserManagementServiceClient is a client for the common_fate.v1alpha1.UserManagementService
// service.
type UserManagementServiceClient interface {
	ListUsers(context.Context, *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error)
}

// NewUserManagementServiceClient constructs a client for the
// common_fate.v1alpha1.UserManagementService service. By default, it uses the Connect protocol with
// the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use
// the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
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

// ListUsers calls common_fate.v1alpha1.UserManagementService.ListUsers.
func (c *userManagementServiceClient) ListUsers(ctx context.Context, req *connect_go.Request[v1alpha1.ListUsersRequest]) (*connect_go.Response[v1alpha1.ListUsersResponse], error) {
	return c.listUsers.CallUnary(ctx, req)
}

// UserManagementServiceHandler is an implementation of the
// common_fate.v1alpha1.UserManagementService service.
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
	return "/common_fate.v1alpha1.UserManagementService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.UserManagementService.ListUsers is not implemented"))
}

// ConnectionsServiceClient is a client for the common_fate.v1alpha1.ConnectionsService service.
type ConnectionsServiceClient interface {
	ListConnections(context.Context, *connect_go.Request[v1alpha1.ListConnectionsRequest]) (*connect_go.Response[v1alpha1.ListConnectionsResponse], error)
	AddConnection(context.Context, *connect_go.Request[v1alpha1.AddConnectionRequest]) (*connect_go.Response[v1alpha1.AddConnectionsResponse], error)
	RemoveConnection(context.Context, *connect_go.Request[v1alpha1.RemoveConnectionRequest]) (*connect_go.Response[v1alpha1.RemoveConnectionResponse], error)
	ListAWSAccounts(context.Context, *connect_go.Request[v1alpha1.ListAWSAccountsRequest]) (*connect_go.Response[v1alpha1.ListAWSAccountsResponse], error)
	GetAWSRolesForAccount(context.Context, *connect_go.Request[v1alpha1.ListAWSRolesForAccountRequest]) (*connect_go.Response[v1alpha1.ListAWSRolesForAccountResponse], error)
}

// NewConnectionsServiceClient constructs a client for the common_fate.v1alpha1.ConnectionsService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewConnectionsServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ConnectionsServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &connectionsServiceClient{
		listConnections: connect_go.NewClient[v1alpha1.ListConnectionsRequest, v1alpha1.ListConnectionsResponse](
			httpClient,
			baseURL+ConnectionsServiceListConnectionsProcedure,
			opts...,
		),
		addConnection: connect_go.NewClient[v1alpha1.AddConnectionRequest, v1alpha1.AddConnectionsResponse](
			httpClient,
			baseURL+ConnectionsServiceAddConnectionProcedure,
			opts...,
		),
		removeConnection: connect_go.NewClient[v1alpha1.RemoveConnectionRequest, v1alpha1.RemoveConnectionResponse](
			httpClient,
			baseURL+ConnectionsServiceRemoveConnectionProcedure,
			opts...,
		),
		listAWSAccounts: connect_go.NewClient[v1alpha1.ListAWSAccountsRequest, v1alpha1.ListAWSAccountsResponse](
			httpClient,
			baseURL+ConnectionsServiceListAWSAccountsProcedure,
			opts...,
		),
		getAWSRolesForAccount: connect_go.NewClient[v1alpha1.ListAWSRolesForAccountRequest, v1alpha1.ListAWSRolesForAccountResponse](
			httpClient,
			baseURL+ConnectionsServiceGetAWSRolesForAccountProcedure,
			opts...,
		),
	}
}

// connectionsServiceClient implements ConnectionsServiceClient.
type connectionsServiceClient struct {
	listConnections       *connect_go.Client[v1alpha1.ListConnectionsRequest, v1alpha1.ListConnectionsResponse]
	addConnection         *connect_go.Client[v1alpha1.AddConnectionRequest, v1alpha1.AddConnectionsResponse]
	removeConnection      *connect_go.Client[v1alpha1.RemoveConnectionRequest, v1alpha1.RemoveConnectionResponse]
	listAWSAccounts       *connect_go.Client[v1alpha1.ListAWSAccountsRequest, v1alpha1.ListAWSAccountsResponse]
	getAWSRolesForAccount *connect_go.Client[v1alpha1.ListAWSRolesForAccountRequest, v1alpha1.ListAWSRolesForAccountResponse]
}

// ListConnections calls common_fate.v1alpha1.ConnectionsService.ListConnections.
func (c *connectionsServiceClient) ListConnections(ctx context.Context, req *connect_go.Request[v1alpha1.ListConnectionsRequest]) (*connect_go.Response[v1alpha1.ListConnectionsResponse], error) {
	return c.listConnections.CallUnary(ctx, req)
}

// AddConnection calls common_fate.v1alpha1.ConnectionsService.AddConnection.
func (c *connectionsServiceClient) AddConnection(ctx context.Context, req *connect_go.Request[v1alpha1.AddConnectionRequest]) (*connect_go.Response[v1alpha1.AddConnectionsResponse], error) {
	return c.addConnection.CallUnary(ctx, req)
}

// RemoveConnection calls common_fate.v1alpha1.ConnectionsService.RemoveConnection.
func (c *connectionsServiceClient) RemoveConnection(ctx context.Context, req *connect_go.Request[v1alpha1.RemoveConnectionRequest]) (*connect_go.Response[v1alpha1.RemoveConnectionResponse], error) {
	return c.removeConnection.CallUnary(ctx, req)
}

// ListAWSAccounts calls common_fate.v1alpha1.ConnectionsService.ListAWSAccounts.
func (c *connectionsServiceClient) ListAWSAccounts(ctx context.Context, req *connect_go.Request[v1alpha1.ListAWSAccountsRequest]) (*connect_go.Response[v1alpha1.ListAWSAccountsResponse], error) {
	return c.listAWSAccounts.CallUnary(ctx, req)
}

// GetAWSRolesForAccount calls common_fate.v1alpha1.ConnectionsService.GetAWSRolesForAccount.
func (c *connectionsServiceClient) GetAWSRolesForAccount(ctx context.Context, req *connect_go.Request[v1alpha1.ListAWSRolesForAccountRequest]) (*connect_go.Response[v1alpha1.ListAWSRolesForAccountResponse], error) {
	return c.getAWSRolesForAccount.CallUnary(ctx, req)
}

// ConnectionsServiceHandler is an implementation of the common_fate.v1alpha1.ConnectionsService
// service.
type ConnectionsServiceHandler interface {
	ListConnections(context.Context, *connect_go.Request[v1alpha1.ListConnectionsRequest]) (*connect_go.Response[v1alpha1.ListConnectionsResponse], error)
	AddConnection(context.Context, *connect_go.Request[v1alpha1.AddConnectionRequest]) (*connect_go.Response[v1alpha1.AddConnectionsResponse], error)
	RemoveConnection(context.Context, *connect_go.Request[v1alpha1.RemoveConnectionRequest]) (*connect_go.Response[v1alpha1.RemoveConnectionResponse], error)
	ListAWSAccounts(context.Context, *connect_go.Request[v1alpha1.ListAWSAccountsRequest]) (*connect_go.Response[v1alpha1.ListAWSAccountsResponse], error)
	GetAWSRolesForAccount(context.Context, *connect_go.Request[v1alpha1.ListAWSRolesForAccountRequest]) (*connect_go.Response[v1alpha1.ListAWSRolesForAccountResponse], error)
}

// NewConnectionsServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewConnectionsServiceHandler(svc ConnectionsServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	connectionsServiceListConnectionsHandler := connect_go.NewUnaryHandler(
		ConnectionsServiceListConnectionsProcedure,
		svc.ListConnections,
		opts...,
	)
	connectionsServiceAddConnectionHandler := connect_go.NewUnaryHandler(
		ConnectionsServiceAddConnectionProcedure,
		svc.AddConnection,
		opts...,
	)
	connectionsServiceRemoveConnectionHandler := connect_go.NewUnaryHandler(
		ConnectionsServiceRemoveConnectionProcedure,
		svc.RemoveConnection,
		opts...,
	)
	connectionsServiceListAWSAccountsHandler := connect_go.NewUnaryHandler(
		ConnectionsServiceListAWSAccountsProcedure,
		svc.ListAWSAccounts,
		opts...,
	)
	connectionsServiceGetAWSRolesForAccountHandler := connect_go.NewUnaryHandler(
		ConnectionsServiceGetAWSRolesForAccountProcedure,
		svc.GetAWSRolesForAccount,
		opts...,
	)
	return "/common_fate.v1alpha1.ConnectionsService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ConnectionsServiceListConnectionsProcedure:
			connectionsServiceListConnectionsHandler.ServeHTTP(w, r)
		case ConnectionsServiceAddConnectionProcedure:
			connectionsServiceAddConnectionHandler.ServeHTTP(w, r)
		case ConnectionsServiceRemoveConnectionProcedure:
			connectionsServiceRemoveConnectionHandler.ServeHTTP(w, r)
		case ConnectionsServiceListAWSAccountsProcedure:
			connectionsServiceListAWSAccountsHandler.ServeHTTP(w, r)
		case ConnectionsServiceGetAWSRolesForAccountProcedure:
			connectionsServiceGetAWSRolesForAccountHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedConnectionsServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedConnectionsServiceHandler struct{}

func (UnimplementedConnectionsServiceHandler) ListConnections(context.Context, *connect_go.Request[v1alpha1.ListConnectionsRequest]) (*connect_go.Response[v1alpha1.ListConnectionsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.ConnectionsService.ListConnections is not implemented"))
}

func (UnimplementedConnectionsServiceHandler) AddConnection(context.Context, *connect_go.Request[v1alpha1.AddConnectionRequest]) (*connect_go.Response[v1alpha1.AddConnectionsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.ConnectionsService.AddConnection is not implemented"))
}

func (UnimplementedConnectionsServiceHandler) RemoveConnection(context.Context, *connect_go.Request[v1alpha1.RemoveConnectionRequest]) (*connect_go.Response[v1alpha1.RemoveConnectionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.ConnectionsService.RemoveConnection is not implemented"))
}

func (UnimplementedConnectionsServiceHandler) ListAWSAccounts(context.Context, *connect_go.Request[v1alpha1.ListAWSAccountsRequest]) (*connect_go.Response[v1alpha1.ListAWSAccountsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.ConnectionsService.ListAWSAccounts is not implemented"))
}

func (UnimplementedConnectionsServiceHandler) GetAWSRolesForAccount(context.Context, *connect_go.Request[v1alpha1.ListAWSRolesForAccountRequest]) (*connect_go.Response[v1alpha1.ListAWSRolesForAccountResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.ConnectionsService.GetAWSRolesForAccount is not implemented"))
}

// UsageMetricsServiceClient is a client for the common_fate.v1alpha1.UsageMetricsService service.
type UsageMetricsServiceClient interface {
	GetAWSRoleMetrics(context.Context, *connect_go.Request[v1alpha1.GetAWSRoleMetricsRequest]) (*connect_go.Response[v1alpha1.GetAWSRoleMetricsResponse], error)
	GetUsageForRole(context.Context, *connect_go.Request[v1alpha1.GetUsageForRoleRequest]) (*connect_go.Response[v1alpha1.GetUsageForRoleResponse], error)
}

// NewUsageMetricsServiceClient constructs a client for the common_fate.v1alpha1.UsageMetricsService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewUsageMetricsServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) UsageMetricsServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &usageMetricsServiceClient{
		getAWSRoleMetrics: connect_go.NewClient[v1alpha1.GetAWSRoleMetricsRequest, v1alpha1.GetAWSRoleMetricsResponse](
			httpClient,
			baseURL+UsageMetricsServiceGetAWSRoleMetricsProcedure,
			opts...,
		),
		getUsageForRole: connect_go.NewClient[v1alpha1.GetUsageForRoleRequest, v1alpha1.GetUsageForRoleResponse](
			httpClient,
			baseURL+UsageMetricsServiceGetUsageForRoleProcedure,
			opts...,
		),
	}
}

// usageMetricsServiceClient implements UsageMetricsServiceClient.
type usageMetricsServiceClient struct {
	getAWSRoleMetrics *connect_go.Client[v1alpha1.GetAWSRoleMetricsRequest, v1alpha1.GetAWSRoleMetricsResponse]
	getUsageForRole   *connect_go.Client[v1alpha1.GetUsageForRoleRequest, v1alpha1.GetUsageForRoleResponse]
}

// GetAWSRoleMetrics calls common_fate.v1alpha1.UsageMetricsService.GetAWSRoleMetrics.
func (c *usageMetricsServiceClient) GetAWSRoleMetrics(ctx context.Context, req *connect_go.Request[v1alpha1.GetAWSRoleMetricsRequest]) (*connect_go.Response[v1alpha1.GetAWSRoleMetricsResponse], error) {
	return c.getAWSRoleMetrics.CallUnary(ctx, req)
}

// GetUsageForRole calls common_fate.v1alpha1.UsageMetricsService.GetUsageForRole.
func (c *usageMetricsServiceClient) GetUsageForRole(ctx context.Context, req *connect_go.Request[v1alpha1.GetUsageForRoleRequest]) (*connect_go.Response[v1alpha1.GetUsageForRoleResponse], error) {
	return c.getUsageForRole.CallUnary(ctx, req)
}

// UsageMetricsServiceHandler is an implementation of the common_fate.v1alpha1.UsageMetricsService
// service.
type UsageMetricsServiceHandler interface {
	GetAWSRoleMetrics(context.Context, *connect_go.Request[v1alpha1.GetAWSRoleMetricsRequest]) (*connect_go.Response[v1alpha1.GetAWSRoleMetricsResponse], error)
	GetUsageForRole(context.Context, *connect_go.Request[v1alpha1.GetUsageForRoleRequest]) (*connect_go.Response[v1alpha1.GetUsageForRoleResponse], error)
}

// NewUsageMetricsServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewUsageMetricsServiceHandler(svc UsageMetricsServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	usageMetricsServiceGetAWSRoleMetricsHandler := connect_go.NewUnaryHandler(
		UsageMetricsServiceGetAWSRoleMetricsProcedure,
		svc.GetAWSRoleMetrics,
		opts...,
	)
	usageMetricsServiceGetUsageForRoleHandler := connect_go.NewUnaryHandler(
		UsageMetricsServiceGetUsageForRoleProcedure,
		svc.GetUsageForRole,
		opts...,
	)
	return "/common_fate.v1alpha1.UsageMetricsService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case UsageMetricsServiceGetAWSRoleMetricsProcedure:
			usageMetricsServiceGetAWSRoleMetricsHandler.ServeHTTP(w, r)
		case UsageMetricsServiceGetUsageForRoleProcedure:
			usageMetricsServiceGetUsageForRoleHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedUsageMetricsServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedUsageMetricsServiceHandler struct{}

func (UnimplementedUsageMetricsServiceHandler) GetAWSRoleMetrics(context.Context, *connect_go.Request[v1alpha1.GetAWSRoleMetricsRequest]) (*connect_go.Response[v1alpha1.GetAWSRoleMetricsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.UsageMetricsService.GetAWSRoleMetrics is not implemented"))
}

func (UnimplementedUsageMetricsServiceHandler) GetUsageForRole(context.Context, *connect_go.Request[v1alpha1.GetUsageForRoleRequest]) (*connect_go.Response[v1alpha1.GetUsageForRoleResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.UsageMetricsService.GetUsageForRole is not implemented"))
}

// AccessServiceClient is a client for the common_fate.v1alpha1.AccessService service.
type AccessServiceClient interface {
	ListEntitlementsForProvider(context.Context, *connect_go.Request[v1alpha1.ListEntitlementsForProviderRequest]) (*connect_go.Response[v1alpha1.ListEntitlementsForProviderResponse], error)
	CreateAccessRequest(context.Context, *connect_go.Request[v1alpha1.CreateAccessRequestRequest]) (*connect_go.Response[v1alpha1.CreateAccessRequestResponse], error)
}

// NewAccessServiceClient constructs a client for the common_fate.v1alpha1.AccessService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccessServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AccessServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accessServiceClient{
		listEntitlementsForProvider: connect_go.NewClient[v1alpha1.ListEntitlementsForProviderRequest, v1alpha1.ListEntitlementsForProviderResponse](
			httpClient,
			baseURL+AccessServiceListEntitlementsForProviderProcedure,
			opts...,
		),
		createAccessRequest: connect_go.NewClient[v1alpha1.CreateAccessRequestRequest, v1alpha1.CreateAccessRequestResponse](
			httpClient,
			baseURL+AccessServiceCreateAccessRequestProcedure,
			opts...,
		),
	}
}

// accessServiceClient implements AccessServiceClient.
type accessServiceClient struct {
	listEntitlementsForProvider *connect_go.Client[v1alpha1.ListEntitlementsForProviderRequest, v1alpha1.ListEntitlementsForProviderResponse]
	createAccessRequest         *connect_go.Client[v1alpha1.CreateAccessRequestRequest, v1alpha1.CreateAccessRequestResponse]
}

// ListEntitlementsForProvider calls common_fate.v1alpha1.AccessService.ListEntitlementsForProvider.
func (c *accessServiceClient) ListEntitlementsForProvider(ctx context.Context, req *connect_go.Request[v1alpha1.ListEntitlementsForProviderRequest]) (*connect_go.Response[v1alpha1.ListEntitlementsForProviderResponse], error) {
	return c.listEntitlementsForProvider.CallUnary(ctx, req)
}

// CreateAccessRequest calls common_fate.v1alpha1.AccessService.CreateAccessRequest.
func (c *accessServiceClient) CreateAccessRequest(ctx context.Context, req *connect_go.Request[v1alpha1.CreateAccessRequestRequest]) (*connect_go.Response[v1alpha1.CreateAccessRequestResponse], error) {
	return c.createAccessRequest.CallUnary(ctx, req)
}

// AccessServiceHandler is an implementation of the common_fate.v1alpha1.AccessService service.
type AccessServiceHandler interface {
	ListEntitlementsForProvider(context.Context, *connect_go.Request[v1alpha1.ListEntitlementsForProviderRequest]) (*connect_go.Response[v1alpha1.ListEntitlementsForProviderResponse], error)
	CreateAccessRequest(context.Context, *connect_go.Request[v1alpha1.CreateAccessRequestRequest]) (*connect_go.Response[v1alpha1.CreateAccessRequestResponse], error)
}

// NewAccessServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccessServiceHandler(svc AccessServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	accessServiceListEntitlementsForProviderHandler := connect_go.NewUnaryHandler(
		AccessServiceListEntitlementsForProviderProcedure,
		svc.ListEntitlementsForProvider,
		opts...,
	)
	accessServiceCreateAccessRequestHandler := connect_go.NewUnaryHandler(
		AccessServiceCreateAccessRequestProcedure,
		svc.CreateAccessRequest,
		opts...,
	)
	return "/common_fate.v1alpha1.AccessService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AccessServiceListEntitlementsForProviderProcedure:
			accessServiceListEntitlementsForProviderHandler.ServeHTTP(w, r)
		case AccessServiceCreateAccessRequestProcedure:
			accessServiceCreateAccessRequestHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAccessServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccessServiceHandler struct{}

func (UnimplementedAccessServiceHandler) ListEntitlementsForProvider(context.Context, *connect_go.Request[v1alpha1.ListEntitlementsForProviderRequest]) (*connect_go.Response[v1alpha1.ListEntitlementsForProviderResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.AccessService.ListEntitlementsForProvider is not implemented"))
}

func (UnimplementedAccessServiceHandler) CreateAccessRequest(context.Context, *connect_go.Request[v1alpha1.CreateAccessRequestRequest]) (*connect_go.Response[v1alpha1.CreateAccessRequestResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("common_fate.v1alpha1.AccessService.CreateAccessRequest is not implemented"))
}
