package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============ Mock ConfigProvider ============

type mockConfig struct {
	accessExpireMin int64
	accessSecret    string
	refreshExpireHr int64
	refreshSecret   string
}

func (m *mockConfig) AccessTokenExpiredTime() int64  { return m.accessExpireMin }
func (m *mockConfig) AccessTokenSecret() string      { return m.accessSecret }
func (m *mockConfig) RefreshTokenExpiredTime() int64 { return m.refreshExpireHr }
func (m *mockConfig) RefreshTokenSecret() string     { return m.refreshSecret }

// ============ 辅助函数 ============

func defaultConfig() *mockConfig {
	return &mockConfig{
		accessExpireMin: 15,
		accessSecret:    "access-secret-key",
		refreshExpireHr: 24,
		refreshSecret:   "refresh-secret-key",
	}
}

func fixedTime() time.Time {
	return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
}

func fixedTokenID() string {
	return "test-token-id-12345"
}

func sameSecretConfig() *mockConfig {
	return &mockConfig{
		accessExpireMin: 15,
		accessSecret:    "shared-secret-key",
		refreshExpireHr: 24,
		refreshSecret:   "shared-secret-key",
	}
}

// ============ NewJWTService 测试 ============

func TestNewJWTService_DefaultOptions(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	require.NotNil(t, svc)
	assert.Equal(t, cfg, svc.config)
	assert.NotNil(t, svc.nowFunc)
	assert.NotNil(t, svc.tokenIDFn)
}

func TestNewJWTService_WithNowFunc(t *testing.T) {
	cfg := defaultConfig()
	fn := fixedTime
	svc := NewJWTService(cfg, WithNowFunc(fn))

	require.NotNil(t, svc.nowFunc)
	assert.Equal(t, fn(), svc.nowFunc())
}

func TestNewJWTService_WithTokenIDFunc(t *testing.T) {
	cfg := defaultConfig()
	fn := fixedTokenID
	svc := NewJWTService(cfg, WithTokenIDFunc(fn))

	require.NotNil(t, svc.tokenIDFn)
	assert.Equal(t, fn(), svc.tokenIDFn())
}

// ============ GenAccessToken 测试 ============

func TestGenAccessToken_Success(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, err := svc.GenAccessToken("user-001")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// 手动解析验证 claims
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(cfg.AccessTokenSecret()), nil
		},
		jwt.WithTimeFunc(fixedTime),
	)
	require.NoError(t, err)

	claims, ok := token.Claims.(*CustomClaims)
	require.True(t, ok)
	assert.True(t, token.Valid)

	// 验证各字段
	assert.Equal(t, TokenTypeAccess, claims.TokenType)
	assert.Equal(t, "user-001", claims.Subject)
	assert.Equal(t, fixedTokenID(), claims.ID)

	expectedExp := fixedTime().Add(time.Duration(cfg.AccessTokenExpiredTime()) * time.Minute)
	assert.True(t, claims.ExpiresAt.Time.Equal(expectedExp), "expires at mismatch")
	assert.True(t, claims.IssuedAt.Time.Equal(fixedTime()), "issued at mismatch")
	assert.True(t, claims.NotBefore.Time.Equal(fixedTime()), "not before mismatch")
}

// ============ VerifyAccessToken 测试 ============

func TestVerifyAccessToken_Success(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenAccessToken("user-001")
	uid, err := svc.VerifyAccessToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, "user-001", uid)
}

func TestVerifyAccessToken_Expired(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenAccessToken("user-001")

	futureTime := fixedTime().Add(16 * time.Minute)
	verifySvc := NewJWTService(cfg, WithNowFunc(func() time.Time { return futureTime }))

	_, err := verifySvc.VerifyAccessToken(tokenStr)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestVerifyAccessToken_WrongSecret(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenAccessToken("user-001")

	wrongCfg := &mockConfig{
		accessExpireMin: 15,
		accessSecret:    "wrong-secret",
		refreshExpireHr: 24,
		refreshSecret:   "refresh-secret-key",
	}
	wrongSvc := NewJWTService(wrongCfg, WithNowFunc(fixedTime))

	_, err := wrongSvc.VerifyAccessToken(tokenStr)
	assert.Error(t, err)
}

func TestVerifyAccessToken_RefreshTokenUsedAsAccess(t *testing.T) {
	cfg := sameSecretConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	refreshToken, _ := svc.GenRefreshToken("user-001")

	_, err := svc.VerifyAccessToken(refreshToken)
	assert.EqualError(t, err, "token type mismatch: expected access")
}

func TestVerifyAccessToken_InvalidTokenString(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	_, err := svc.VerifyAccessToken("invalid.token.string")
	assert.Error(t, err)
}

func TestVerifyAccessToken_EmptyToken(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	_, err := svc.VerifyAccessToken("")
	assert.Error(t, err)
}

// ============ GenRefreshToken 测试 ============

func TestGenRefreshToken_Success(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, err := svc.GenRefreshToken("user-001")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(cfg.RefreshTokenSecret()), nil
		},
		jwt.WithTimeFunc(fixedTime),
	)
	require.NoError(t, err)

	claims, ok := token.Claims.(*CustomClaims)
	require.True(t, ok)
	assert.True(t, token.Valid)

	assert.Equal(t, TokenTypeRefresh, claims.TokenType)
	assert.Equal(t, "user-001", claims.Subject)
	assert.Equal(t, fixedTokenID(), claims.ID)

	expectedExp := fixedTime().Add(time.Duration(cfg.RefreshTokenExpiredTime()) * time.Hour)
	assert.True(t, claims.ExpiresAt.Time.Equal(expectedExp), "expires at mismatch")
}

// ============ VerifyRefreshToken 测试 ============

func TestVerifyRefreshToken_Success(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenRefreshToken("user-001")
	uid, err := svc.VerifyRefreshToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, "user-001", uid)
}

func TestVerifyRefreshToken_Expired(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenRefreshToken("user-001")

	futureTime := fixedTime().Add(25 * time.Hour)
	verifySvc := NewJWTService(cfg, WithNowFunc(func() time.Time { return futureTime }))

	_, err := verifySvc.VerifyRefreshToken(tokenStr)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestVerifyRefreshToken_WrongSecret(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenRefreshToken("user-001")

	wrongCfg := &mockConfig{
		accessExpireMin: 15,
		accessSecret:    "access-secret-key",
		refreshExpireHr: 24,
		refreshSecret:   "wrong-refresh-secret",
	}
	wrongSvc := NewJWTService(wrongCfg, WithNowFunc(fixedTime))

	_, err := wrongSvc.VerifyRefreshToken(tokenStr)
	assert.Error(t, err)
}

func TestVerifyRefreshToken_AccessTokenUsedAsRefresh(t *testing.T) {
	cfg := sameSecretConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	accessToken, _ := svc.GenAccessToken("user-001")

	_, err := svc.VerifyRefreshToken(accessToken)
	assert.EqualError(t, err, "token type mismatch: expected refresh")
}

func TestVerifyRefreshToken_InvalidTokenString(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	_, err := svc.VerifyRefreshToken("invalid.token.string")
	assert.Error(t, err)
}

func TestVerifyRefreshToken_EmptyToken(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	_, err := svc.VerifyRefreshToken("")
	assert.Error(t, err)
}

// ============ RenewAccessToken 测试 ============

func TestRenewAccessToken_Success(t *testing.T) {
	cfg := defaultConfig()
	callCount := 0
	svc := NewJWTService(cfg,
		WithNowFunc(fixedTime),
		WithTokenIDFunc(func() string {
			callCount++
			return fmt.Sprintf("token-id-%d", callCount)
		}),
	)

	refreshToken, _ := svc.GenRefreshToken("user-001")

	newAccess, newRefresh, err := svc.RenewAccessToken(refreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)

	uid, err := svc.VerifyAccessToken(newAccess)
	require.NoError(t, err)
	assert.Equal(t, "user-001", uid)

	uid, err = svc.VerifyRefreshToken(newRefresh)
	require.NoError(t, err)
	assert.Equal(t, "user-001", uid)
}

func TestRenewAccessToken_WithInvalidRefreshToken(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg)

	_, _, err := svc.RenewAccessToken("invalid-refresh-token")
	assert.Error(t, err)
}

func TestRenewAccessToken_WithAccessToken(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	accessToken, _ := svc.GenAccessToken("user-001")

	_, _, err := svc.RenewAccessToken(accessToken)
	assert.Error(t, err)
}

func TestRenewAccessToken_WithExpiredRefreshToken(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	refreshToken, _ := svc.GenRefreshToken("user-001")

	futureTime := fixedTime().Add(25 * time.Hour)
	renewSvc := NewJWTService(cfg, WithNowFunc(func() time.Time { return futureTime }))

	_, _, err := renewSvc.RenewAccessToken(refreshToken)
	assert.Error(t, err)
}

// ============ Token 类型常量测试 ============

func TestTokenTypeConstants(t *testing.T) {
	assert.Equal(t, TokenType("access"), TokenTypeAccess)
	assert.Equal(t, TokenType("refresh"), TokenTypeRefresh)
}

// ============ 签名算法篡改测试 ============

func TestVerifyAccessToken_TamperedSigningMethod(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	claims := CustomClaims{
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(fixedTime().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(fixedTime()),
			NotBefore: jwt.NewNumericDate(fixedTime()),
			Subject:   "hacker",
			ID:        "fake-id",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Skipf("could not create none-signed token: %v", err)
	}

	_, err = svc.VerifyAccessToken(tokenStr)
	assert.Error(t, err)
}

// ============ 不同 uid 边界测试 ============

func TestGenAccessToken_DifferentUIDs(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	uids := []string{"", "a", "user-123", "非常长的用户ID_🔥🔥🔥"}
	for _, uid := range uids {
		tokenStr, err := svc.GenAccessToken(uid)
		require.NoError(t, err, "GenAccessToken(%q) failed", uid)

		gotUID, err := svc.VerifyAccessToken(tokenStr)
		require.NoError(t, err, "VerifyAccessToken for uid=%q failed", uid)
		assert.Equal(t, uid, gotUID, "uid mismatch")
	}
}

func TestGenRefreshToken_DifferentUIDs(t *testing.T) {
	cfg := defaultConfig()
	svc := NewJWTService(cfg, WithNowFunc(fixedTime), WithTokenIDFunc(fixedTokenID))

	uids := []string{"", "a", "user-456", "非常长的用户ID_🔥🔥🔥"}
	for _, uid := range uids {
		tokenStr, err := svc.GenRefreshToken(uid)
		require.NoError(t, err, "GenRefreshToken(%q) failed", uid)

		gotUID, err := svc.VerifyRefreshToken(tokenStr)
		require.NoError(t, err, "VerifyRefreshToken for uid=%q failed", uid)
		assert.Equal(t, uid, gotUID, "uid mismatch")
	}
}

// ============ NotBefore 测试 ============

func TestVerifyAccessToken_NotBeforeViolation(t *testing.T) {
	cfg := defaultConfig()
	pastTime := fixedTime()
	svc := NewJWTService(cfg, WithNowFunc(func() time.Time { return pastTime }), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenAccessToken("user-001")

	earlyTime := pastTime.Add(-1 * time.Minute)
	earlySvc := NewJWTService(cfg, WithNowFunc(func() time.Time { return earlyTime }))

	_, err := earlySvc.VerifyAccessToken(tokenStr)
	assert.Error(t, err)
}

func TestVerifyRefreshToken_NotBeforeViolation(t *testing.T) {
	cfg := defaultConfig()
	pastTime := fixedTime()
	svc := NewJWTService(cfg, WithNowFunc(func() time.Time { return pastTime }), WithTokenIDFunc(fixedTokenID))

	tokenStr, _ := svc.GenRefreshToken("user-001")

	earlyTime := pastTime.Add(-1 * time.Minute)
	earlySvc := NewJWTService(cfg, WithNowFunc(func() time.Time { return earlyTime }))

	_, err := earlySvc.VerifyRefreshToken(tokenStr)
	assert.Error(t, err)
}
