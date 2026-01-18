package mapper

import (
	"license/config"
	"license/jetbrains/code/entity"
	"license/logger"

	"gorm.io/gorm/clause"
)

type ProductMapper interface {
	Truncate() error
	SaveBatch(products []*entity.ProductEntity) error
	UpsertBatch(products []*entity.ProductEntity) error
	List() ([]entity.ProductEntity, error)
}

type GormProductMapper struct{}

func (m *GormProductMapper) Truncate() error {
	result := config.DB.Exec("DELETE FROM sys_jetbrains_product")
	if result.Error != nil {
		logger.Error("Failed to truncate product table:", result.Error)
		return result.Error
	}
	return nil
}

func (m *GormProductMapper) SaveBatch(products []*entity.ProductEntity) error {
	result := config.DB.CreateInBatches(products, len(products))
	if result.Error != nil {
		logger.Error("Failed to save products:", result.Error)
		return result.Error
	}
	return nil
}

// UpsertBatch inserts new products or updates existing ones based on product_code unique index
func (m *GormProductMapper) UpsertBatch(products []*entity.ProductEntity) error {
	result := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"product_name", "product_detail"}),
	}).CreateInBatches(products, 100)
	if result.Error != nil {
		logger.Error("Failed to upsert products:", result.Error)
		return result.Error
	}
	return nil
}

func (m *GormProductMapper) List() ([]entity.ProductEntity, error) {
	var products []entity.ProductEntity
	result := config.DB.Find(&products)
	if result.Error != nil {
		logger.Error("Failed to fetch products:", result.Error)
		return nil, result.Error
	}
	return products, nil
}

type PluginMapper interface {
	Truncate() error
	SaveBatch(products []*entity.PluginEntity) error
	UpsertBatch(plugins []*entity.PluginEntity) error
	List() ([]entity.PluginEntity, error)
}

type GormPluginMapper struct{}

func (m *GormPluginMapper) Truncate() error {
	result := config.DB.Exec("DELETE FROM sys_jetbrains_paid_plugin")
	if result.Error != nil {
		logger.Error("Failed to truncate plugin table:", result.Error)
		return result.Error
	}
	return nil
}

func (m *GormPluginMapper) SaveBatch(plugins []*entity.PluginEntity) error {
	result := config.DB.CreateInBatches(plugins, len(plugins))
	if result.Error != nil {
		logger.Error("Failed to save plugins:", result.Error)
		return result.Error
	}
	return nil
}

// UpsertBatch inserts new plugins or updates existing ones based on plugin_code unique index
func (m *GormPluginMapper) UpsertBatch(plugins []*entity.PluginEntity) error {
	result := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"plugin_id", "plugin_name", "plugin_api_detail"}),
	}).CreateInBatches(plugins, 100)
	if result.Error != nil {
		logger.Error("Failed to upsert plugins:", result.Error)
		return result.Error
	}
	return nil
}

func (m *GormPluginMapper) List() ([]entity.PluginEntity, error) {
	var plugins []entity.PluginEntity
	result := config.DB.Find(&plugins)
	if result.Error != nil {
		logger.Error("Failed to fetch plugins:", result.Error)
		return nil, result.Error
	}
	return plugins, nil
}
