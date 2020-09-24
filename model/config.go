package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Configs is the main configuration object
type Configs struct {
	TodoHost    string `json:"todo_host"    yaml:"todo_host"    mapstructure:"todo_host"`
	AdminHeader string `json:"admin_header" yaml:"admin_header" mapstructure:"admin_header"`
	AdminSecret string `json:"admin_secret" yaml:"admin_secret" mapstructure:"admin_secret"`
}

// Validate the Configurations
func (confs Configs) Validate() error {
	return validation.ValidateStruct(&confs,
		validation.Field(&confs.TodoHost, validation.Required, is.URL),
		validation.Field(&confs.AdminHeader, validation.Required, validation.Length(6, 64)),
		validation.Field(&confs.AdminSecret, validation.Required, validation.Length(6, 64)),
	)
}

// DefaultConfigs returns the sensible default configs
func DefaultConfigs() *Configs {
	return &Configs{
		TodoHost:    "http://127.0.0.1:6969",
		AdminHeader: "TodoAdminHeader",
		AdminSecret: "TodoAdminSecret",
	}
}
