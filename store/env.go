package store

import (
	"fmt"
	"time"

	"github.com/awillis/bitfan/api/models"
	"github.com/timshannon/bolthold"
)

type StoreEnv struct {
	Uuid string `json:"uuid" boltholdKey:"Uuid"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string `json:"name"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

func (s *Store) CreateEnv(xp *models.Env) {
	sxp := &StoreEnv{
		Uuid:      xp.Uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Name:   xp.Name,
		Value:  xp.Value,
		Secret: xp.Secret,
	}

	s.db.Upsert(sxp.Uuid, sxp)
}

func (s *Store) DeleteEnv(p *models.Env) {
	err := s.db.DeleteMatching(&StoreAsset{}, bolthold.Where(bolthold.Key).Eq(p.Uuid))
	if err != nil {
		s.log.Error("Store : DeleteEnv -" + err.Error())
		return
	}

	err = s.db.Delete(p.Uuid, &StoreEnv{})
	if err != nil {
		s.log.Error("Store : DeleteEnv - " + err.Error())
		return
	}
}

func (s *Store) FindEnvs() []models.Env {
	xps := []models.Env{}

	var sxps []StoreEnv
	query := &bolthold.Query{}
	err := s.db.Find(&sxps, query)
	if err != nil {
		s.log.Error("Store : FindEnvs " + err.Error())
		return xps
	}
	for _, xp := range sxps {
		var tEnv models.Env
		tEnv.Uuid = xp.Uuid
		tEnv.Name = xp.Name
		tEnv.Value = xp.Value
		tEnv.Secret = xp.Secret

		xps = append(xps, tEnv)
	}

	return xps
}

func (s *Store) FindOneEnvByUUID(UUID string) (models.Env, error) {
	tEnv := models.Env{Uuid: UUID}

	var sps []StoreEnv
	err := s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(UUID))
	if err != nil {
		return tEnv, err
	}
	if len(sps) == 0 {
		return tEnv, fmt.Errorf("Env %s not found", UUID)
	}

	tEnv.CreatedAt = sps[0].CreatedAt
	tEnv.UpdatedAt = sps[0].UpdatedAt

	tEnv.Name = sps[0].Name
	tEnv.Value = sps[0].Value
	tEnv.Secret = sps[0].Secret

	return tEnv, nil
}
