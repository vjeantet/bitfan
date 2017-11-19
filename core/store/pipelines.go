package store

import (
	"fmt"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vjeantet/bitfan/core/models"
)

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

func (s *Store) FindPipelinesWithAutoStart() []models.Pipeline {
	pps := []models.Pipeline{}

	var sps []StorePipeline
	err := s.db.Find(&sps, bolthold.Where("AutoStart").Eq(true))
	if err != nil {
		s.log.Error("Store : FindPipelinesWithAutoStart " + err.Error())
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
		s.log.Error("Store : DeletePipeline -" + err.Error())
		return
	}

	err = s.db.Delete(p.Uuid, &StorePipeline{})
	if err != nil {
		s.log.Error("Store : DeletePipeline - " + err.Error())
		return
	}

}

func (s *Store) FindPipelines(withAssetValues bool) []models.Pipeline {
	pps := []models.Pipeline{}

	var sps []StorePipeline
	err := s.db.Find(&sps, &bolthold.Query{})
	if err != nil {
		s.log.Error("Store : FindPipelines - " + err.Error())
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
