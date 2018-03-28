package store

import (
	"fmt"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vjeantet/bitfan/api/models"
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
		s.log.Error("Store : createAsset - " + err.Error())
		return
	}
	if len(sps) == 0 {
		s.log.Error("Store : createAsset - can not find pipeline " + a.PipelineUUID)
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
		s.log.Error("Store : SaveAsset -" + err.Error())
		return
	}
	if len(sps) == 0 {
		s.log.Error("Store : SaveAsset - can not find pipeline : " + a.PipelineUUID)
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
		s.log.Error("Store : DeleteAsset : " + err.Error())
		return
	}

	var sps []StorePipeline
	err = s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(a.PipelineUUID))
	if err != nil {
		s.log.Error("Store : DeleteAsset - " + err.Error())
		return
	}
	if len(sps) == 0 {
		s.log.Error("Store : DeleteAsset - can not find pipeline " + a.PipelineUUID)
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

func (s *Store) FindAssetsByPipelineUUID(uuid string) ([]models.Asset, error) {
	assets := []models.Asset{}

	var sas []StoreAsset
	err := s.db.Find(&sas, bolthold.Where("PipelineUUID").Eq(uuid))
	if err != nil {
		return assets, err
	}
	if len(sas) == 0 {
		return assets, fmt.Errorf("Assets for pipeline UUID=%s not found", uuid)
	}

	for _, sa := range sas {

		var asset models.Asset
		asset.ContentType = sa.ContentType
		asset.CreatedAt = sa.CreatedAt
		asset.Name = sa.Label
		asset.PipelineUUID = sa.PipelineUUID
		asset.Size = sa.Size
		asset.Type = sa.Type
		asset.UpdatedAt = sa.UpdatedAt
		asset.Uuid = sa.Uuid
		asset.Value = sa.Value

		assets = append(assets, asset)
	}
	return assets, nil
}
