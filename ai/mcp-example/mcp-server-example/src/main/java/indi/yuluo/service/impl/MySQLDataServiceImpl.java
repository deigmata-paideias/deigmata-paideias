package indi.yuluo.service.impl;

import java.util.Objects;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import indi.yuluo.entity.Data;
import indi.yuluo.repository.DataMapper;
import indi.yuluo.service.MySQLDataService;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class MySQLDataServiceImpl extends ServiceImpl<DataMapper, Data> implements MySQLDataService {

	public String getContextData(String name) {

		Data user = this.getOne(
				new QueryWrapper<Data>().eq("name", name));

		if (Objects.nonNull(user)) {
			return user.toString();
		}

		return "未查询到用户名为 " + name + " 的用户信息";
	}

}
