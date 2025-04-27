package postgresql_sqlboiler

import (
	"github.com/hdget/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "postgresql-sqlboiler"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryDb,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
