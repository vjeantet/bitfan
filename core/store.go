package core

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vjeantet/bitfan/core/models"
)

type Store struct {
	db *bolthold.Store
}

type StorePipeline struct {
	Uuid string `json:"uuid" boltholdKey:"Uuid"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StartedAt time.Time `json:"started_at"`

	Label       string `json:"label"`
	Description string `json:"description"`
	AutoStart   bool   `json:"auto_start" boltholdIndex:"AutoStart" mapstructure:"auto_start"`

	// Assets
	Assets []StoreAssetRef `json:"assets"`
}

type StoreAssetRef struct {
	Uuid  string `json:"uuid" boltholdIndex:"PipelineAssetUUID"`
	Label string
	Type  string
}

type StoreAsset struct {
	Uuid      string    `json:"uuid" boltholdKey:"Uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PipelineUUID string `boltholdIndex:"AssetPipelineUUID"`

	Label       string
	Type        string `boltholdIndex:"AssetType"`
	ContentType string
	Value       []byte
	Size        int
}

func NewStore(location string) (*Store, error) {
	database, err := bolthold.Open(filepath.Join(location, "bitfan.bolt.db"), 0666, nil)
	return &Store{db: database}, err
}

func (s *Store) close() {
	s.db.Close()
}

func (s *Store) FindPipelinesWithAutoStart() []models.Pipeline {
	pps := []models.Pipeline{}

	var sps []StorePipeline
	err := s.db.Find(&sps, bolthold.Where("AutoStart").Eq(true))
	if err != nil {
		Log().Errorf("Store : FindPipelinesWithAutoStart - %s", err.Error())
		return pps
	}

	for _, p := range sps {
		var tPipeline models.Pipeline
		tPipeline.Uuid = p.Uuid
		tPipeline.CreatedAt = p.CreatedAt
		tPipeline.UpdatedAt = p.UpdatedAt
		tPipeline.StartedAt = p.StartedAt
		tPipeline.Label = p.Label
		tPipeline.Description = p.Description
		tPipeline.AutoStart = p.AutoStart

		for _, a := range p.Assets {
			tPipeline.Assets = append(tPipeline.Assets, models.Asset{
				Uuid: a.Uuid,
				Name: a.Label,
				Type: a.Type,
			})
		}
		pps = append(pps, tPipeline)
	}

	return pps
}

func (s *Store) FindOnePipelineByUUID(UUID string, withAssetValues bool) (models.Pipeline, error) {
	tPipeline := models.Pipeline{Uuid: UUID}

	var sps []StorePipeline
	err := s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(UUID))
	if err != nil {
		return tPipeline, err
	}
	if len(sps) == 0 {
		return tPipeline, fmt.Errorf("Pipeline %s not found", UUID)
	}

	tPipeline.CreatedAt = sps[0].CreatedAt
	tPipeline.UpdatedAt = sps[0].UpdatedAt
	tPipeline.StartedAt = sps[0].StartedAt
	tPipeline.Label = sps[0].Label
	tPipeline.Description = sps[0].Description
	tPipeline.AutoStart = sps[0].AutoStart

	for _, a := range sps[0].Assets {
		asset := models.Asset{
			Uuid: a.Uuid,
			Name: a.Label,
			Type: a.Type,
		}

		if withAssetValues {
			var sas []StoreAsset
			err := s.db.Find(&sas, bolthold.Where(bolthold.Key).Eq(a.Uuid))
			if err != nil {
				return tPipeline, err
			}
			if len(sas) == 0 {
				return tPipeline, fmt.Errorf("Asset %s not found", a.Uuid)
			}
			asset.Value = sas[0].Value
			asset.Size = sas[0].Size
			asset.ContentType = sas[0].ContentType
		}
		tPipeline.Assets = append(tPipeline.Assets, asset)
	}

	return tPipeline, nil
}

func (s *Store) CreatePipeline(p *models.Pipeline) {
	sp := &StorePipeline{
		Uuid:        p.Uuid,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Label:       p.Label,
		Description: p.Description,
		AutoStart:   p.AutoStart,
	}

	for _, a := range p.Assets {
		sp.Assets = append(sp.Assets,
			StoreAssetRef{
				Uuid:  a.Uuid,
				Label: a.Name,
				Type:  a.Type,
			})
		sav := &StoreAsset{
			Uuid:         a.Uuid,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			PipelineUUID: sp.Uuid,
			Label:        a.Name,
			Type:         a.Type,
			ContentType:  a.ContentType,
			Value:        a.Value,
			Size:         a.Size,
		}
		s.db.Upsert(sav.Uuid, sav)
	}

	s.db.Upsert(sp.Uuid, sp)
}

func (s *Store) SavePipeline(p *models.Pipeline) {
	sp := &StorePipeline{
		Uuid:        p.Uuid,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   time.Now(),
		Label:       p.Label,
		Description: p.Description,
		AutoStart:   p.AutoStart,
	}

	for _, a := range p.Assets {
		sp.Assets = append(sp.Assets,
			StoreAssetRef{
				Uuid:  a.Uuid,
				Label: a.Name,
				Type:  a.Type,
			})

		if a.PipelineUUID != "" { // a.PipelineUUID = "" means its a RefAsset
			sav := &StoreAsset{
				Uuid:         a.Uuid,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				PipelineUUID: sp.Uuid,
				Label:        a.Name,
				Type:         a.Type,
				ContentType:  a.ContentType,
				Value:        a.Value,
				Size:         a.Size,
			}
			s.db.Upsert(sav.Uuid, sav)
		}
	}

	s.db.Upsert(sp.Uuid, sp)
}

func (s *Store) DeletePipeline(p *models.Pipeline) {
	err := s.db.DeleteMatching(&StoreAsset{}, bolthold.Where("PipelineUUID").Eq(p.Uuid))
	if err != nil {
		Log().Errorf("Store : DeletePipeline - %s", err.Error())
		return
	}

	err = s.db.Delete(p.Uuid, &StorePipeline{})
	if err != nil {
		Log().Errorf("Store : DeletePipeline - %s", err.Error())
		return
	}

}

func (s *Store) FindPipelines(withAssetValues bool) []models.Pipeline {
	pps := []models.Pipeline{}

	var sps []StorePipeline
	err := s.db.Find(&sps, &bolthold.Query{})
	if err != nil {
		Log().Errorf("Store : FindPipelines - %s", err.Error())
		return pps
	}

	for _, p := range sps {
		var tPipeline models.Pipeline
		tPipeline.Uuid = p.Uuid
		tPipeline.CreatedAt = p.CreatedAt
		tPipeline.UpdatedAt = p.UpdatedAt
		tPipeline.StartedAt = p.StartedAt
		tPipeline.Label = p.Label
		tPipeline.Description = p.Description
		tPipeline.AutoStart = p.AutoStart

		// for _, a := range p.Assets {
		// 	tPipeline.Assets = append(tPipeline.Assets, models.Asset{
		// 		Uuid: a.Uuid,
		// 		Name: a.Label,
		// 		Type: a.Type,
		// 	})
		// }

		for _, a := range p.Assets {
			asset := models.Asset{
				Uuid: a.Uuid,
				Name: a.Label,
				Type: a.Type,
			}

			if withAssetValues {
				var sas []StoreAsset
				err := s.db.Find(&sas, bolthold.Where(bolthold.Key).Eq(a.Uuid))
				if err != nil {
					return pps
				}
				if len(sas) == 0 {
					return pps
				}
				asset.Value = sas[0].Value
				asset.Size = sas[0].Size
				asset.ContentType = sas[0].ContentType
			}
			tPipeline.Assets = append(tPipeline.Assets, asset)
		}
		pps = append(pps, tPipeline)
	}

	return pps
}

func (s *Store) CreateAsset(a *models.Asset) {
	sav := &StoreAsset{
		Uuid:         a.Uuid,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PipelineUUID: a.PipelineUUID,
		Label:        a.Name,
		Type:         a.Type,
		ContentType:  a.ContentType,
		Value:        a.Value,
		Size:         a.Size,
	}
	s.db.Upsert(sav.Uuid, sav)

	var sps []StorePipeline
	err := s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(a.PipelineUUID))
	if err != nil {
		Log().Errorf("Store : createAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		Log().Errorf("Store : createAsset - can not find pipeline (%s)", a.PipelineUUID)
		return
	}

	sps[0].Assets = append(sps[0].Assets, StoreAssetRef{
		Uuid:  a.Uuid,
		Label: a.Name,
		Type:  a.Type,
	})
	s.db.Upsert(sps[0].Uuid, sps[0])

}

func (s *Store) SaveAsset(a *models.Asset) {
	sav := &StoreAsset{
		Uuid:         a.Uuid,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    time.Now(),
		PipelineUUID: a.PipelineUUID,
		Label:        a.Name,
		Type:         a.Type,
		ContentType:  a.ContentType,
		Value:        a.Value,
		Size:         a.Size,
	}
	s.db.Upsert(sav.Uuid, sav)

	var sps []StorePipeline
	err := s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(a.PipelineUUID))
	if err != nil {
		Log().Errorf("Store : SaveAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		Log().Errorf("Store : SaveAsset - can not find pipeline (%s)", a.PipelineUUID)
		return
	}

	for i, _ := range sps[0].Assets {
		if sps[0].Assets[i].Uuid == a.Uuid {
			sps[0].Assets[i].Label = a.Name
			sps[0].Assets[i].Type = a.Type
		}
	}
	s.db.Upsert(sps[0].Uuid, sps[0])

}

func (s *Store) DeleteAsset(a *models.Asset) {
	err := s.db.Delete(a.Uuid, &StoreAsset{})
	if err != nil {
		Log().Errorf("Store : DeleteAsset - %s", err.Error())
		return
	}

	var sps []StorePipeline
	err = s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(a.PipelineUUID))
	if err != nil {
		Log().Errorf("Store : DeleteAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		Log().Errorf("Store : DeleteAsset - can not find pipeline (%s)", a.PipelineUUID)
		return
	}

	sars := []StoreAssetRef{}
	for _, ar := range sps[0].Assets {
		if ar.Uuid != a.Uuid {
			sars = append(sars, ar)
		}
	}
	sps[0].Assets = sars
	s.db.Upsert(sps[0].Uuid, sps[0])

}

func (s *Store) FindOneAssetByUUID(uuid string) (models.Asset, error) {
	asset := models.Asset{Uuid: uuid}

	var sas []StoreAsset
	err := s.db.Find(&sas, bolthold.Where(bolthold.Key).Eq(uuid))
	if err != nil {
		return asset, err
	}
	if len(sas) == 0 {
		return asset, fmt.Errorf("Asset %s not found", uuid)
	}

	asset.ContentType = sas[0].ContentType
	asset.CreatedAt = sas[0].CreatedAt
	asset.Name = sas[0].Label
	asset.PipelineUUID = sas[0].PipelineUUID
	asset.Size = sas[0].Size
	asset.Type = sas[0].Type
	asset.UpdatedAt = sas[0].UpdatedAt
	asset.Uuid = sas[0].Uuid
	asset.Value = sas[0].Value

	return asset, nil
}
