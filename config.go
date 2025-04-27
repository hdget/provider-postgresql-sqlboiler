package postgresql_sqlboiler

import (
	"github.com/hdget/common/intf"
	"github.com/pkg/errors"
)

type psqlProviderConfig struct {
	Default *psqlConfig   `mapstructure:"default"`
	Master  *psqlConfig   `mapstructure:"master"`
	Slaves  []*psqlConfig `mapstructure:"slaves"`
	Items   []*psqlConfig `mapstructure:"items"`
}

type psqlConfig struct {
	Name            string `mapstructure:"name"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Schema          string `mapstructure:"schema"`
	UsePgBouncer    bool   `mapstructure:"use_pg_bouncer"`
	MaxOpenConn     int    `mapstructure:"max_client_conn"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // unit: seconds
}

const (
	configSection = "sdk.postgresql"
)

var (
	errInvalidConfig = errors.New("invalid postgresql provider config")
	errEmptyConfig   = errors.New("empty postgresql provider config")
)

func newConfig(configProvider intf.ConfigProvider) (*psqlProviderConfig, error) {
	if configProvider == nil {
		return nil, errInvalidConfig
	}

	var c *psqlProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate postgresql provider config")
	}

	return c, nil
}

func (c *psqlProviderConfig) validate() error {
	if c.Default != nil {
		err := c.validateInstance(c.Default)
		if err != nil {
			return err
		}
	}

	if c.Master != nil {
		err := c.validateInstance(c.Master)
		if err != nil {
			return err
		}
	}

	for _, slave := range c.Slaves {
		err := c.validateInstance(slave)
		if err != nil {
			return err
		}
	}

	for _, item := range c.Items {
		err := c.validateExtraInstance(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *psqlProviderConfig) validateInstance(ic *psqlConfig) error {
	if ic == nil || ic.Host == "" || ic.User == "" {
		return errEmptyConfig
	}

	if ic.Schema == "" {
		ic.Schema = "public"
	}

	// setup default config value
	if ic.Port == 0 {
		if ic.UsePgBouncer {
			ic.Port = 6432
		} else {
			ic.Port = 5432
		}
	}

	if ic.MaxOpenConn == 0 {
		ic.MaxOpenConn = 100
	}

	if ic.ConnMaxLifetime == 0 {
		ic.ConnMaxLifetime = 50 * 60
	}

	return nil
}

func (c *psqlProviderConfig) validateExtraInstance(ic *psqlConfig) error {
	if ic == nil || ic.Host == "" || ic.Name == "" {
		return errEmptyConfig
	}

	// setup default config value
	if ic.Port == 0 {
		ic.Port = 3306
	}
	return nil
}
