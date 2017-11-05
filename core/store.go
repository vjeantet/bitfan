package core

import (
	"fmt"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vjeantet/bitfan/core/models"
)

type Store struct {
	db *gorm.DB
}

func NewStore(location string) (*Store, error) {
	database, err := gorm.Open("sqlite3", filepath.Join(location, "bitfan.db"))
	// database.LogMode(true)

	// Migrate the schema
	database.AutoMigrate(&models.Pipeline{}, &models.Asset{})

	return &Store{db: database}, err
}

func (s *Store) close() {
	s.db.Close()
}

func (s *Store) FindPipelinesWithAutoStart() []models.Pipeline {
	pipelinesToStart := []models.Pipeline{}
	s.db.Where(&models.Pipeline{AutoStart: true}).Find(&pipelinesToStart).RecordNotFound()
	return pipelinesToStart
}

func (s *Store) FindOnePipelineByUUID(UUID string) (models.Pipeline, error) {
	tPipeline := models.Pipeline{Uuid: UUID}
	if s.db.Preload("Assets").Where(&tPipeline).First(&tPipeline).RecordNotFound() {
		return tPipeline, fmt.Errorf("Pipeline %s not found", UUID)
	}
	return tPipeline, nil
}

func (s *Store) CreatePipeline(p *models.Pipeline) {
	s.db.Create(p)
}

func (s *Store) SavePipeline(p *models.Pipeline) {
	s.db.Save(p)
}

func (s *Store) DeletePipeline(p *models.Pipeline) {
	s.db.Delete(models.Asset{}, "pipeline_uuid = ?", p.Uuid)
	s.db.Delete(&p)
}

func (s *Store) FindPipelines() []models.Pipeline {
	pps := []models.Pipeline{}
	s.db.Find(&pps)
	return pps
}

func (s *Store) CreateAsset(a *models.Asset) {
	s.db.Create(a)
}

func (s *Store) SaveAsset(a *models.Asset) {
	s.db.Save(a)
}

func (s *Store) DeleteAsset(a *models.Asset) {
	s.db.Delete(&a)
}

func (s *Store) FindOneAssetByUUID(uuid string) (models.Asset, error) {
	asset := models.Asset{Uuid: uuid}
	if s.db.Where(&asset).First(&asset).RecordNotFound() {
		return asset, fmt.Errorf("Asset %s not found", uuid)
	}
	return asset, nil
}
