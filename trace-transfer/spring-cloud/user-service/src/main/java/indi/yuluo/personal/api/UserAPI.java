package indi.yuluo.personal.api;

import indi.yuluo.personal.entity.User;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

@FeignClient(name = "kratos-user-svc", url = "http://localhost:8000")
public interface UserAPI {

    @GetMapping(value = "/api/v1/user/{id}")
    User getUser(@PathVariable("id") String id);

}
