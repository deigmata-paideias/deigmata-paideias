package indi.yuluo.entity;

import com.baomidou.mybatisplus.annotation.TableName;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@TableName("data")
public class Data {

	private Long id;

	private String name;

	private Integer age;

	private String email;

	public Long getId() {
		return id;
	}

	public void setId(Long id) {
		this.id = id;
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	public Integer getAge() {
		return age;
	}

	public void setAge(Integer age) {
		this.age = age;
	}

	public String getEmail() {
		return email;
	}

	public void setEmail(String email) {
		this.email = email;
	}

	@Override
	public String toString() {

		return "Data{" + "id=" + id
				+ ", name='" + name + '\''
				+ ", age=" + age
				+ ", email='" + email + '\''
				+ '}';
	}

}
