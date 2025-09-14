package indi.yuluo.personal.user;

import indi.yuluo.personal.api.UserAPI;
import indi.yuluo.personal.entity.User;
import jakarta.servlet.http.HttpServletRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/users")
public class UserController {

    private static final Logger logger = LoggerFactory.getLogger(UserController.class);

    private final UserAPI userAPI;

    public UserController(UserAPI userAPI) {
        this.userAPI = userAPI;
    }

    @GetMapping("/hi")
    public String hi() {

        return "Hello, User!";
    }

    @GetMapping("/rpc/{id}")
    public String rpc(@PathVariable("id") String id, HttpServletRequest request) {

        logger.info("User 服务发起 RPC 调用, 请求路径: {}", request.getRequestURI());
        // 打印所有 header
        java.util.Enumeration<String> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
            String headerName = headerNames.nextElement();
            logger.info("[user-service] header: {} = {}", headerName, request.getHeader(headerName));
        }
        User user = userAPI.getUser(id);
        return user.getName();
    }

}
