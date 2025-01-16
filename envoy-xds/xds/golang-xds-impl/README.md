# xDS Golang 实现

本示例中演示，如何使用 golang 搭建一个 xDS server，用来给 envoy 动态下发规则。

## 环境准备

首先需要启动一个 envoy 和两个 nginx 应用的 docker 容器。

1. 修改 envoy-example/docker-compose 中的配置挂载为 envoy-dynamic.yaml

    ```yaml
    volumes:
      # 单个代理配置
      # - ./envoy/conf/envoy-standalone.yaml:/etc/envoy/envoy.yaml
      # 代理百度
      # - ./envoy/conf/envoy-proxy-baidu.yaml:/etc/envoy/envoy.yaml
      # 代理多个应用，集群代理
      # - ./envoy/conf/envoy-cluster.yaml:/etc/envoy/envoy.yaml
      # xds 动态代理
      - ./envoy/conf/envoy-dynamic.yaml:/etc/envoy/envoy.yaml
      # 通过文件的形式设置代理
      # - ./envoy/files/envoy.yaml:/etc/envoy/envoy.yaml
      # - ./envoy/files/cds.yaml:/etc/envoy/cds.yaml
      # - ./envoy/files/lds.yaml:/etc/envoy/lds.yaml
    ```

2. 启动容器

    ```shell
    cd envoy-example
    docker compose up -d
    ```

## envoy-dynamic.yaml 配置分析

```yaml
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

node:
  # 节点名称
  cluster: xds_test_node
  # 唯一标识， 服务发现时通过唯一标识去识别节点
  id: test_1

dynamic_resources:
  cds_config:
    # 使用的资源版本号
    resource_api_version: V3
    api_config_source:
      # 连接类型
      api_type: GRPC
      transport_api_version: V3
      refresh_delay: 5s
      # grpc 地址
      grpc_services:
        - envoy_grpc:
            cluster_name: grpc_xds_cluster

  lds_config:
    resource_api_version: V3
    api_config_source:
      api_type: GRPC
      transport_api_version: V3
      grpc_services:
        - envoy_grpc:
            cluster_name: grpc_xds_cluster

# 通过静态配置一个 xds_cluster 服务用于做 xDS 控制面，数据面为 envoy 
static_resources:
  clusters:
    - connect_timeout: 1s
      load_assignment:
        cluster_name: grpc_xds_cluster
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      # control-plane xds 控制面服务地址
                      # 需要修改为本机 ip 地址，和 docker 容器通信
                      address: 192.168.2.17
                      port_value: 9090
      http2_protocol_options: {}
      name: grpc_xds_cluster
```

## goalng-xds-imple 代码分析

### 代码目录

```shell
├── README.md
├── cmd
│   ├── client.go           # xDS client ，获取 envoy 中的 xDS 配置
│   └── server.go           # xDS server，envoy 监听此服务端口，来完成 xDS 更新
├── go.mod
├── go.sum
├── middleware
│   └── call_back.go        # xDS 回调
└── pkg
    ├── constant
    │   └── constants.go    # example 中用到的常量定义
    ├── log
    │   └── logger.go       # 日志组件
    └── resources
        └── resource.go     # 具体 xDS resource 的组装
```

### 代码分析

这里分析下 server.go 和 resource.go 代码文件

```go
func main() {

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,

		// 一条 GRPC 连接允许并发的发送和接收多个 Stream
		grpc.MaxConcurrentStreams(1000),

		// 连接超过多少时间 不活跃，则会去探测 是否依然 alive
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Second * 30,
			Timeout: time.Second * 5,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{

			// 发送ping之前最少要等待 xx 时间
			MinTime: time.Second * 30,

			// 连接空闲时仍然发送 PING 帧进行监测
			PermitWithoutStream: true,
		}),

		// 设置 GRPC 帧大小, mac 运行没有问题，windows 运行会出现 http2 帧太大的问题。
		// grpc.MaxSendMsgSize(math.MaxInt64),
		// grpc.MaxRecvMsgSize(math.MaxInt64),
	)

	// 创建 GRPC 服务
	grpcServer := grpc.NewServer(grpcOptions...)

	// 开启 debug 日志模式，在 callback 中打印日志。
	xLog := log.XLogger{Debug: true}

	// 创建缓存系统，实质是内部维护了一个缓存来存储 xds 配置信息。
	c := cache.NewSnapshotCache(false, cache.IDHash{}, xLog)

	// envoy 配置的缓存快照, 1 是版本号, 通过版本号的变更进行配置更新。
	snapshot := resources.GenerateSnapshot("1")
	if err := snapshot.Consistent(); err != nil {
		xLog.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	// Add the snapshot to the cache. nodeID 必须要设置，对应 envoy/envoy-dynamic.yaml 中设置的 nodeId
	nodeID := constant.NODE_ID
	if err := c.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		os.Exit(1)
	}

	// 请求回调, 类似于中间件，在 envoy 接受到配置更新时打印日志信息。
	cb := middleware.CallBacksMiddleWares{XLog: xLog}

	// 官方提供的控制面server
	srv := server.NewServer(context.Background(), c, &cb)
	// 注册 集群服务
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	// 注册 listener
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)
	// 由于在 listener 下需要创建路由, 所以需要加入
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)

	errCh := make(chan error)

	go func() {
		// envoy 需要监听的 xds server 端口。
		fmt.Printf("GRPC server started in %v ... \n", constant.SERVER_PORT)
		lis, err := net.Listen(constant.NETWORK_MODEL, fmt.Sprintf(":%d", constant.SERVER_PORT))
		if err != nil {
			errCh <- err
			return
		}
		if err = grpcServer.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	// 启动动态测试服务, 可以通过请求 /test 进行版本更替 [并不是 GRPC 服务端口]
	go func() {
		r := gin.New()
		r.GET(constant.TEST_URI, func(ctx *gin.Context) {
			// 如果部署 2个 nginx 容器, 可以通过这个 IP 的调整测试出是否成功从代理 v1-nginx 转成代理  v2-nginx
			resources.UpstreamHost = constant.UPSTREAM_SERVER_2
			// 通过版本控制snapshot 的更新
			ss := resources.GenerateSnapshot("2")
			if err := c.SetSnapshot(ctx, nodeID, ss); err != nil {
				ctx.String(400, err.Error())
				return
			}
			ctx.String(200, "OK")

			log.XLogger{}.Infof("update snapshot success")
		})

		if err := r.Run(constant.WILDCARD_HOST + ":" + strconv.Itoa(constant.DYNAMIC_TEST_SERVER_PORT)); err != nil {
			errCh <- err
		}

	}()

	err := <-errCh

	xLog.Errorf("err is %s", err.Error())
}
```

```go
var (
	ClusterName  = constant.CLUSTER_NAME
	RouteName    = constant.ROUTER_NAME
	ListerName   = constant.LISTENER_NAME
	ListenerPort = constant.LISTENER_PORT
	UpstreamHost = constant.UPSTREAM_SERVER_1
	UpstreamPort = constant.UPSTREAM_SERVER_PORT
)

// CDS 配置
func makeCluster(clusterName string) *clusterv3.Cluster {

	return &clusterv3.Cluster{
		Name:                 ClusterName,
		ConnectTimeout:       durationpb.New(3 * time.Second),
		// 集群类型
		ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_LOGICAL_DNS},
		// 负载均衡策略
		LbPolicy:             clusterv3.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily:      clusterv3.Cluster_V4_ONLY,
	}
}

// makeEndpoint 设置 endpoint
func makeEndpoint(clusterName string) *endpointv3.ClusterLoadAssignment {
	return &endpointv3.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpointv3.LocalityLbEndpoints{
			{
				LbEndpoints: []*endpointv3.LbEndpoint{
					{
						HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
							Endpoint: &endpointv3.Endpoint{
								Address: &core.Address{
									Address: &core.Address_SocketAddress{
										SocketAddress: &core.SocketAddress{
											// 需要代理的服务信息
											Protocol: core.SocketAddress_TCP,
											// host 地址
											Address:  UpstreamHost,
											PortSpecifier: &core.SocketAddress_PortValue{
												// 需要代理的服务端口，nginx 应用运行在 nginx-app 容器的 80 端口
												PortValue: uint32(UpstreamPort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeRoute(routeName string, clusterName string) *routev3.RouteConfiguration {
	return &routev3.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*routev3.VirtualHost{
			{
				Name:    "xds_test_hosts",
				Domains: []string{"*"},
				Routes: []*routev3.Route{
					{
						Match: &routev3.RouteMatch{
							PathSpecifier: &routev3.RouteMatch_Prefix{
								// 路由名字，不设置重写的时候，会全部转发给 nginx-app 容器
								Prefix: "/",
							},
						},
						Action: &routev3.Route_Route{
							Route: &routev3.RouteAction{
								ClusterSpecifier: &routev3.RouteAction_Cluster{
									Cluster: clusterName,
								},
								HostRewriteSpecifier: &routev3.RouteAction_HostRewriteLiteral{
									HostRewriteLiteral: UpstreamHost,
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeHTTPListener(listenerName string, route string) *listenerv3.Listener {

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
			// fix: gRPC config for type.googleapis.com/envoy.config.listener.v3.Listener rejected: Error adding/updating listener(s) listener_0: Didn't find a registered implementation for 'envoy.filters.http.router' with type URL: ''
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: &anypb.Any{
					TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router",
				},
			},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listenerv3.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  constant.WILDCARD_HOST,
					PortSpecifier: &core.SocketAddress_PortValue{
						// envoy 监听的端口，可以通过 ip://ListenerPort 访问代理的服务
						PortValue: uint32(ListenerPort),
					},
				},
			},
		},
		// 过滤器链扩展
		FilterChains: []*listenerv3.FilterChain{{
			Filters: []*listenerv3.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listenerv3.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					// 这里的 ClusterName 需要和 上面的 ClusterName 区分开，这里是配置文件中配置的 xds server 的 ClusterName. 
					// 如果写错 envoy 不会报错，但是 envoy 代理服务访问失败。
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: constant.GRPC_STATIC_CLUSTER_NAME},
				},
			}},
		},
	}
	return source
}

// GenerateSnapshot 创建缓存快照
// 真正存放不同xds配置的文件, 通过snapshot构建出不同资源类型的服务发现资源(其实就是构造不同的配置结构体)
func GenerateSnapshot(version string) *cache.Snapshot {

	snap, _ := cache.NewSnapshot(version,
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster(ClusterName)},
			resource.RouteType:    {makeRoute(RouteName, ClusterName)},
			resource.ListenerType: {makeHTTPListener(ListerName, RouteName)},
		},
	)

	return snap
}
```

### 示例演示

1. 构建和启动 xds-server

    ```shell
    # build
    go build -o xds-server cmd/server.go

    # 启动
    ./xds-server
    ```

2. 访问演示

    启动

    ```shell
    $ curl 127.0.0.1:8082

    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Envoy-Nginx-App</title>
    </head>
    <body>
        <p>This is nginx app pages for 1.</p>
    </body>
    </html> 
    ```

3. 切换 xDS 资源

    ```shell
    curl 127.0.0.1:19090/test

    # 访问测试

    $ curl 127.0.0.1:8082      
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Envoy-Nginx-App</title>
    </head>
    <body>
        <p>This is nginx app pages for 2.</p>
    </body>
    </html>
    ```
