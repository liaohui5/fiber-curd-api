package passwd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PasswdEncrypt_Success(t *testing.T) {
	origin := "mySecretPassword123"
	hash, err := PasswdEncrypt(origin)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	// bcrypt 哈希总是以 $2a$ 开头
	assert.Contains(t, hash, "$2a$")
	// 哈希值不应该等于原始密码
	assert.NotEqual(t, origin, hash)
}

func Test_PasswdEncrypt_DifferentPasswordsProduceDifferentHashes(t *testing.T) {
	hash1, err1 := PasswdEncrypt("passwordA")
	hash2, err2 := PasswdEncrypt("passwordB")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// 不同的密码应产生不同的哈希
	assert.NotEqual(t, hash1, hash2)
}

func Test_PasswdEncrypt_SamePasswordProducesDifferentHashes(t *testing.T) {
	hash1, err1 := PasswdEncrypt("samePassword")
	hash2, err2 := PasswdEncrypt("samePassword")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// 相同密码因盐值不同也应产生不同的哈希
	assert.NotEqual(t, hash1, hash2)
}

func Test_PasswdEncrypt_EmptyString(t *testing.T) {
	hash, err := PasswdEncrypt("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$2a$")
}

func Test_PasswdEncrypt_LongPassword(t *testing.T) {
	// bcrypt 有 72 字节的长度限制，超过该限制会返回错误
	longPassword := "this_is_a_very_very_long_password_that_exceeds_the_72_bytes_limit_of_bcrypt_and_will_error"
	hash, err := PasswdEncrypt(longPassword)

	assert.Error(t, err)
	assert.Empty(t, hash)
}

func Test_PasswdVerify_CorrectPassword(t *testing.T) {
	origin := "correctPassword123"
	hash, err := PasswdEncrypt(origin)
	assert.NoError(t, err)

	result := PasswdVerify(origin, hash)
	assert.True(t, result)
}

func Test_PasswdVerify_WrongPassword(t *testing.T) {
	origin := "correctPassword123"
	hash, err := PasswdEncrypt(origin)
	assert.NoError(t, err)

	result := PasswdVerify("wrongPassword456", hash)
	assert.False(t, result)
}

func Test_PasswdVerify_EmptyPassword(t *testing.T) {
	hash, err := PasswdEncrypt("somePassword")
	assert.NoError(t, err)

	result := PasswdVerify("", hash)
	assert.False(t, result)
}

func Test_PasswdVerify_EmptyHash(t *testing.T) {
	result := PasswdVerify("somePassword", "")
	assert.False(t, result)
}

func Test_PasswdVerify_InvalidHash(t *testing.T) {
	result := PasswdVerify("somePassword", "not-a-valid-hash")
	assert.False(t, result)
}

func Test_PasswdVerify_DifferentHashForSamePassword(t *testing.T) {
	origin := "sharedPassword"
	hash1, err1 := PasswdEncrypt(origin)
	hash2, err2 := PasswdEncrypt(origin)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// 同一个密码的不同哈希都能通过验证
	assert.True(t, PasswdVerify(origin, hash1))
	assert.True(t, PasswdVerify(origin, hash2))
}

func Test_PasswdEncryptAndVerify_RoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"P@ssw0rd!",
		"中文密码测试",
		"emoji🔐password",
		"a",
		"exactly_72_bytes_padding_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	for _, pw := range passwords {
		t.Run(pw, func(t *testing.T) {
			hash, err := PasswdEncrypt(pw)
			assert.NoError(t, err)
			assert.True(t, PasswdVerify(pw, hash))
			assert.False(t, PasswdVerify(pw+"_wrong", hash))
		})
	}
}
