package indi.yuluo.service;

import com.baomidou.mybatisplus.extension.service.IService;
import indi.yuluo.entity.Data;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

public interface MySQLDataService extends IService<Data> {

	String getContextData(String name);

}
