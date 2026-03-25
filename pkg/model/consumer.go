// Package model provides data models for the SDK
package model

import "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"

// Credential 凭证接口
type Credential interface {
	// GetType 返回凭证类型
	GetType() string
	// Validate 验证凭证
	Validate(forUpdate bool) error
}

// Consumer 服务消费者
type Consumer struct {
	// Name 消费者名称
	Name string `json:"name,omitempty"`

	// Credentials 消费者凭证列表
	Credentials []Credential `json:"credentials,omitempty"`
}

// CredentialType 凭证类型常量
const (
	CredentialTypeKeyAuth = "key-auth"
)

// KeyAuthCredentialSource Key Auth凭证来源
type KeyAuthCredentialSource string

const (
	// KeyAuthCredentialSourceBearer 使用 Authorization: Bearer token 头
	KeyAuthCredentialSourceBearer KeyAuthCredentialSource = "BEARER"
	// KeyAuthCredentialSourceHeader 使用 HTTP 头
	KeyAuthCredentialSourceHeader KeyAuthCredentialSource = "HEADER"
	// KeyAuthCredentialSourceQuery 使用查询参数
	KeyAuthCredentialSourceQuery KeyAuthCredentialSource = "QUERY"
)

// IsKeyRequired 返回是否需要Key
func (s KeyAuthCredentialSource) IsKeyRequired() bool {
	return s != KeyAuthCredentialSourceBearer
}

// ParseKeyAuthCredentialSource 解析凭证来源
func ParseKeyAuthCredentialSource(str string) KeyAuthCredentialSource {
	switch str {
	case "BEARER":
		return KeyAuthCredentialSourceBearer
	case "HEADER":
		return KeyAuthCredentialSourceHeader
	case "QUERY":
		return KeyAuthCredentialSourceQuery
	default:
		return ""
	}
}

// KeyAuthCredential Key Auth凭证
type KeyAuthCredential struct {
	// Source 凭证来源 (BEARER, HEADER, QUERY)
	Source string `json:"source,omitempty"`
	// Key API Key名称 (HEADER/QUERY时需要)
	Key string `json:"key,omitempty"`
	// Values API Key值列表
	Values []string `json:"values,omitempty"`
}

// NewKeyAuthCredential 创建KeyAuth凭证
func NewKeyAuthCredential(source, key string, values []string) *KeyAuthCredential {
	return &KeyAuthCredential{
		Source: source,
		Key:    key,
		Values: values,
	}
}

// GetType 返回凭证类型
func (c *KeyAuthCredential) GetType() string {
	return CredentialTypeKeyAuth
}

// Validate 验证KeyAuthCredential
func (c *KeyAuthCredential) Validate(forUpdate bool) error {
	if c.Source == "" {
		return errors.NewValidationError("source cannot be blank")
	}

	source := ParseKeyAuthCredentialSource(c.Source)
	if source == "" {
		return errors.NewValidationError("unknown source value: " + c.Source)
	}

	if source.IsKeyRequired() && c.Key == "" {
		return errors.NewValidationError("key cannot be blank")
	}

	if !forUpdate && len(c.Values) == 0 {
		return errors.NewValidationError("values cannot be empty")
	}

	return nil
}

// Validate 验证Consumer
func (c *Consumer) Validate(forUpdate bool) error {
	if c.Name == "" {
		return errors.NewValidationError("name cannot be blank")
	}
	if len(c.Credentials) == 0 {
		return errors.NewValidationError("credentials cannot be empty")
	}
	for _, cred := range c.Credentials {
		if err := cred.Validate(forUpdate); err != nil {
			return err
		}
	}
	return nil
}
