package indi.yuluo.personal.config;

import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.Random;

@Configuration
public class GatewayRouterConfig {

    @Bean
    public RouteLocator customRouteLocator(RouteLocatorBuilder builder) {

        return builder.routes()
                .route("user_service", r -> r
                        .order(100)
                        .path("/gw/**")
                        .filters(f -> f.stripPrefix(1))
                        .uri("http://127.0.0.1:8081")
                )
                .build();
    }

}
