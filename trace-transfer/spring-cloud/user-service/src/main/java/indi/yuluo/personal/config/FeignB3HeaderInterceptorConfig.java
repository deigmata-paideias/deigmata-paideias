package indi.yuluo.personal.config;

import feign.RequestInterceptor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

@Configuration
public class FeignB3HeaderInterceptorConfig {

    @Bean
    public RequestInterceptor b3HeaderInterceptor() {
        return template -> {
            ServletRequestAttributes attrs = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
            if (attrs != null) {
                HttpServletRequest request = attrs.getRequest();
                String b3 = request.getHeader("traceparent");
                if (b3 != null) {
                    template.header("traceparent", b3);
                }
            }
        };
    }
}
