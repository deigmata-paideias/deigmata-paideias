package indi.yuluo.personal.result;


public enum ResultEnum implements IResult {


    SUCCESS(8291, "接口调用成功"),
    VALIDATE_FAILED(8292, "参数校验失败"),
    COMMON_FAILED(8293, "接口调用失败"),
    FORBIDDEN(8294, "没有权限访问资源");

    private final Integer code;
    private final String message;

    ResultEnum(Integer code, String message) {
        this.code = code;
        this.message = message;
    }

    @Override
    public Integer getCode() {
        return code;
    }

    @Override
    public String getMessage() {
        return message;
    }

}
