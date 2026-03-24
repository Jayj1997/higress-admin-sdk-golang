// Package model provides data models for the SDK
package model

import "github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"

// Consumer 服务消费者
type Consumer struct {
	// Name 消费者名称
	Name string `json:"name,omitempty"`

	// Credentials 消费者凭证列表
	Credentials []Credential `json:"credentials,omitempty"`
}

// Credential 凭证基类
type Credential struct {
	// Type 凭证类型
	Type string `json:"type,omitempty"`

	// Properties 凭证属性
	Properties map[string]interface{} `json:"-"`
}

// CredentialType 凭证类型常量
const (
	CredentialTypeKeyAuth = "key-auth"
)

// KeyAuthCredential Key Auth凭证
type KeyAuthCredential struct {
	Credential
	// Key API Key名称
	Key string `json:"key,omitempty"`
	// Value API Key值
	Value string `json:"value,omitempty"`
	// In 凭证位置 (header, query)
	In string `json:"in,omitempty"`
}

// Validate 验证Consumer
func (c *Consumer) Validate(forUpdate bool) error {
	if c.Name == "" {
		return errors.NewValidationError("name cannot be blank")
	}
	if len(c.Credentials) == 0 {
		return errors.NewValidationError("credentials cannot be empty")
	}
	for i := range c.Credentials {
		if err := c.Credentials[i].Validate(forUpdate); err != nil {
			return err
		}
	}
	return nil
}

// Validate 验证Credential
func (c *Credential) Validate(forUpdate bool) error {
	// 基类验证为空，子类可以覆盖
	return nil
}

// Validate 验证KeyAuthCredential
func (c *KeyAuthCredential) Validate(forUpdate bool) error {
	if c.Key == "" {
		return errors.NewValidationError("key cannot be blank")
	}
	if c.Value == "" {
		return errors.NewValidationError("value cannot be blank")
	}
	return nil
}
