package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"golang-xds-impl/middleware"
	"golang-xds-impl/pkg/constant"
	"golang-xds-impl/pkg/log"
	"golang-xds-impl/pkg/resources"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// control plane server 服务端代理，envoy 会监听此 server 端口，来完成配置的动态更新。

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

	// 创建 gRPC 服务
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
		// envoy 需要监听的 xDS server 端口。
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

	// 启动 xDS 资源动态更新测试 服务, 可以通过请求 /test 进行版本更替 [并不是 GRPC 服务端口]
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
