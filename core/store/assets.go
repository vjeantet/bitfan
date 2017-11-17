package store

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/timshannon/bolthold"
	"github.com/vjeantet/bitfan/core/models"
)

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
		log.Errorf("Store : createAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		log.Errorf("Store : createAsset - can not find pipeline (%s)", a.PipelineUUID)
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
		log.Errorf("Store : SaveAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		log.Errorf("Store : SaveAsset - can not find pipeline (%s)", a.PipelineUUID)
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
		log.Errorf("Store : DeleteAsset - %s", err.Error())
		return
	}

	var sps []StorePipeline
	err = s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(a.PipelineUUID))
	if err != nil {
		log.Errorf("Store : DeleteAsset - %s", err.Error())
		return
	}
	if len(sps) == 0 {
		log.Errorf("Store : DeleteAsset - can not find pipeline (%s)", a.PipelineUUID)
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
