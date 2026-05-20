package passwd

import "golang.org/x/crypto/bcrypt"

//////////////////////////////////////////////
// NOTE: 密码处理工具方法
//////////////////////////////////////////////

// PasswdEncrypt 对明文密码进行 bcrypt 哈希加密
func PasswdEncrypt(origin string) (string, error) {
	// bcrypt.GenerateFromPassword 内部会自动生成盐值并将其包含在最终的哈希字符串中
	// bcrypt.DefaultCost 是默认的开销因子(默认为 10), 数值越大加密越安全但越耗时
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(origin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwdHash), nil
}

// PasswdVerify 校验明文密码和目标哈希值是否匹配
func PasswdVerify(origin string, target string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(target), []byte(origin))
	return err == nil
}
