package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDomainJSONSerialization 测试 Domain 模型的 JSON 序列化
// 运行命令: go test -v -run TestDomainJSONSerialization ./pkg/model/
func TestDomainJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		domain   Domain
		expected string
	}{
		{
			name: "basic domain",
			domain: Domain{
				Name:        "example.com",
				EnableHTTPS: EnableHTTPSOn,
			},
			expected: `{"name":"example.com","enableHttps":"on"}`,
		},
		{
			name: "domain with certificate",
			domain: Domain{
				Name:           "secure.example.com",
				EnableHTTPS:    EnableHTTPSForce,
				CertIdentifier: "my-cert",
			},
			expected: `{"name":"secure.example.com","enableHttps":"force","certIdentifier":"my-cert"}`,
		},
		{
			name: "domain with version",
			domain: Domain{
				Name:        "versioned.example.com",
				Version:     "12345",
				EnableHTTPS: EnableHTTPSOff,
			},
			expected: `{"name":"versioned.example.com","version":"12345","enableHttps":"off"}`,
		},
		{
			name:     "empty domain",
			domain:   Domain{},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.domain)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

// TestDomainJSONDeserialization 测试 Domain 模型的 JSON 反序列化
// 运行命令: go test -v -run TestDomainJSONDeserialization ./pkg/model/
func TestDomainJSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected Domain
	}{
		{
			name:    "basic domain",
			jsonStr: `{"name":"example.com","enableHttps":"on"}`,
			expected: Domain{
				Name:        "example.com",
				EnableHTTPS: EnableHTTPSOn,
			},
		},
		{
			name:    "domain with certificate",
			jsonStr: `{"name":"secure.example.com","enableHttps":"force","certIdentifier":"my-cert"}`,
			expected: Domain{
				Name:           "secure.example.com",
				EnableHTTPS:    EnableHTTPSForce,
				CertIdentifier: "my-cert",
			},
		},
		{
			name:     "empty json",
			jsonStr:  `{}`,
			expected: Domain{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var domain Domain
			err := json.Unmarshal([]byte(tt.jsonStr), &domain)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, domain)
		})
	}
}

// TestDomainValidate 测试 Domain 模型验证
// 运行命令: go test -v -run TestDomainValidate ./pkg/model/
func TestDomainValidate(t *testing.T) {
	tests := []struct {
		name    string
		domain  Domain
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty name should fail",
			domain:  Domain{},
			wantErr: true,
			errMsg:  "domain name is required",
		},
		{
			name: "valid domain",
			domain: Domain{
				Name: "example.com",
			},
			wantErr: false,
		},
		{
			name: "valid domain with https",
			domain: Domain{
				Name:        "secure.example.com",
				EnableHTTPS: EnableHTTPSOn,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.domain.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDomainIsHTTPS 测试 Domain 的 IsHTTPS 方法
// 运行命令: go test -v -run TestDomainIsHTTPS ./pkg/model/
func TestDomainIsHTTPS(t *testing.T) {
	tests := []struct {
		name    string
		domain  Domain
		isHTTPS bool
		isForce bool
	}{
		{
			name:    "https off",
			domain:  Domain{EnableHTTPS: EnableHTTPSOff},
			isHTTPS: false,
			isForce: false,
		},
		{
			name:    "https on",
			domain:  Domain{EnableHTTPS: EnableHTTPSOn},
			isHTTPS: true,
			isForce: false,
		},
		{
			name:    "https force",
			domain:  Domain{EnableHTTPS: EnableHTTPSForce},
			isHTTPS: true,
			isForce: true,
		},
		{
			name:    "empty https",
			domain:  Domain{},
			isHTTPS: false,
			isForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isHTTPS, tt.domain.IsHTTPS())
			assert.Equal(t, tt.isForce, tt.domain.IsForceHTTPS())
		})
	}
}

// TestEnableHTTPSConstants 测试 HTTPS 常量
// 运行命令: go test -v -run TestEnableHTTPSConstants ./pkg/model/
func TestEnableHTTPSConstants(t *testing.T) {
	assert.Equal(t, "off", EnableHTTPSOff)
	assert.Equal(t, "on", EnableHTTPSOn)
	assert.Equal(t, "force", EnableHTTPSForce)
}
