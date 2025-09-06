package postgresql_sqlboiler

import (
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/hdget/common/types"
)

type psqlProvider struct {
	defaultDb types.DbClient
	masterDb  types.DbClient
	slaveDbs  []types.DbClient
	extraDbs  map[string]types.DbClient
}

func New(configProvider types.ConfigProvider, logger types.LoggerProvider) (types.DbProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	p := &psqlProvider{
		slaveDbs: make([]types.DbClient, len(config.Slaves)),
		extraDbs: make(map[string]types.DbClient),
	}

	if config.Default != nil {
		p.defaultDb, err = newClient(config.Default)
		if err != nil {
			logger.Fatal("init postgresql default connection", "err", err)
		}

		// 设置boil的缺省db
		boil.SetDB(p.defaultDb)
		logger.Debug("init postgresql default", "host", config.Default.Host)
	}

	if config.Master != nil {
		p.masterDb, err = newClient(config.Master)
		if err != nil {
			logger.Fatal("init postgresql master connection", "err", err)
		}
		logger.Debug("init postgresql master", "host", config.Master.Host)
	}

	for i, slaveConf := range config.Slaves {
		p.slaveDbs[i], err = newClient(slaveConf)
		if err != nil {
			logger.Fatal("init postgresql slave connection", "slave", i, "err", err)
		}

		logger.Debug("init postgresql slave", "index", i, "host", slaveConf.Host)
	}

	for _, extraConf := range config.Items {
		p.extraDbs[extraConf.Name], err = newClient(extraConf)
		if err != nil {
			logger.Fatal("new postgresql extra connection", "name", extraConf.Name, "err", err)
		}

		logger.Debug("init postgresql extra", "name", extraConf.Name, "host", extraConf.Host)
	}

	return p, nil
}

func (p *psqlProvider) GetCapability() types.Capability {
	return Capability
}

func (p *psqlProvider) My() types.DbClient {
	return p.defaultDb
}

func (p *psqlProvider) Master() types.DbClient {
	return p.masterDb
}

func (p *psqlProvider) Slave(i int) types.DbClient {
	return p.slaveDbs[i]
}

func (p *psqlProvider) By(name string) types.DbClient {
	return p.extraDbs[name]
}
