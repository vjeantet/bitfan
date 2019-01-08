package store

import (
	"fmt"
	"time"

	"github.com/awillis/bitfan/api/models"
	"github.com/timshannon/bolthold"
)

type StoreXProcessor struct {
	Uuid string `json:"uuid" boltholdKey:"Uuid"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Label                 string   `json:"label" boltholdIndex:"Label" mapstructure:"label"`
	Description           string   `json:"description"`
	Behavior              string   `json:"behavior" boltholdIndex:"Behavior" mapstructure:"behavior"`
	Stream                bool     `json:"stream"`
	Kind                  string   `json:"kind"`
	Args                  []string `json:"args"`
	Code                  string   `json:"code"`
	Command               string   `json:"command"`
	StdinAs               string   `json:"stdin_as"`
	StdoutAs              string   `json:"stdout_as"`
	OptionsCompositionTpl string   `json:"options_composition_tpl"`
	HasDoc                bool     `json:"has_doc"`
}

func (s *Store) CreateXProcessor(xp *models.XProcessor) {
	sxp := &StoreXProcessor{
		Uuid:      xp.Uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Label:                 xp.Label,
		Behavior:              xp.Behavior,
		Stream:                xp.Stream,
		Command:               xp.Command,
		Code:                  xp.Code,
		Args:                  xp.Args,
		Kind:                  xp.Kind,
		StdinAs:               xp.StdinAs,
		StdoutAs:              xp.StdoutAs,
		Description:           xp.Description,
		OptionsCompositionTpl: xp.OptionsCompositionTpl,
		HasDoc:                xp.HasDoc,
	}

	s.db.Upsert(sxp.Uuid, sxp)
}

func (s *Store) SaveXProcessor(xp *models.XProcessor) {
	sp := &StoreXProcessor{
		Uuid:      xp.Uuid,
		CreatedAt: xp.CreatedAt,
		UpdatedAt: time.Now(),

		Label:                 xp.Label,
		Behavior:              xp.Behavior,
		Stream:                xp.Stream,
		Command:               xp.Command,
		Code:                  xp.Code,
		Args:                  xp.Args,
		Kind:                  xp.Kind,
		StdinAs:               xp.StdinAs,
		StdoutAs:              xp.StdoutAs,
		Description:           xp.Description,
		OptionsCompositionTpl: xp.OptionsCompositionTpl,
		HasDoc:                xp.HasDoc,
	}

	s.db.Upsert(sp.Uuid, sp)
}

func (s *Store) DeleteXProcessor(p *models.XProcessor) {
	err := s.db.Delete(p.Uuid, &StoreXProcessor{})
	if err != nil {
		s.log.Error("Store : DeleteXProcessor - " + err.Error())
		return
	}
}

func (s *Store) FindXProcessors(behavior string) []models.XProcessor {
	xps := []models.XProcessor{}

	var sxps []StoreXProcessor
	query := &bolthold.Query{}
	if behavior != "" {
		query = bolthold.Where("Behavior").Eq(behavior)
	}
	err := s.db.Find(&sxps, query)
	if err != nil {
		s.log.Error("Store : FindXProcessors " + err.Error())
		return xps
	}
	for _, xp := range sxps {
		var tXProcessor models.XProcessor
		tXProcessor.Uuid = xp.Uuid
		tXProcessor.Label = xp.Label
		tXProcessor.Behavior = xp.Behavior
		tXProcessor.Stream = xp.Stream
		tXProcessor.Args = xp.Args
		tXProcessor.Kind = xp.Kind
		tXProcessor.Command = xp.Command
		tXProcessor.Code = xp.Code
		tXProcessor.StdinAs = xp.StdinAs
		tXProcessor.StdoutAs = xp.StdoutAs
		tXProcessor.Description = xp.Description
		tXProcessor.OptionsCompositionTpl = xp.OptionsCompositionTpl
		tXProcessor.HasDoc = xp.HasDoc

		xps = append(xps, tXProcessor)
	}

	return xps
}

func (s *Store) FindOneXProcessorByUUID(UUID string, withAssetValues bool) (models.XProcessor, error) {
	tXProcessor := models.XProcessor{Uuid: UUID}

	var sps []StoreXProcessor
	err := s.db.Find(&sps, bolthold.Where(bolthold.Key).Eq(UUID))
	if err != nil {
		return tXProcessor, err
	}
	if len(sps) == 0 {
		return tXProcessor, fmt.Errorf("XProcessor %s not found", UUID)
	}

	tXProcessor.CreatedAt = sps[0].CreatedAt
	tXProcessor.UpdatedAt = sps[0].UpdatedAt

	tXProcessor.Label = sps[0].Label
	tXProcessor.Behavior = sps[0].Behavior
	tXProcessor.Stream = sps[0].Stream
	tXProcessor.Kind = sps[0].Kind
	tXProcessor.Args = sps[0].Args
	tXProcessor.Command = sps[0].Command
	tXProcessor.Code = sps[0].Code
	tXProcessor.StdinAs = sps[0].StdinAs
	tXProcessor.StdoutAs = sps[0].StdoutAs
	tXProcessor.Description = sps[0].Description
	tXProcessor.OptionsCompositionTpl = sps[0].OptionsCompositionTpl
	tXProcessor.HasDoc = sps[0].HasDoc

	return tXProcessor, nil
}

func (s *Store) FindOneXProcessorByName(name string) (models.XProcessor, error) {
	tXProcessor := models.XProcessor{Label: name}

	var sps []StoreXProcessor
	err := s.db.Find(&sps, bolthold.Where("Label").Eq(name))
	if err != nil {
		return tXProcessor, err
	}
	if len(sps) == 0 {
		return tXProcessor, fmt.Errorf("XProcessor %s not found", name)
	}

	tXProcessor.CreatedAt = sps[0].CreatedAt
	tXProcessor.UpdatedAt = sps[0].UpdatedAt

	tXProcessor.Uuid = sps[0].Uuid
	tXProcessor.Label = sps[0].Label
	tXProcessor.Behavior = sps[0].Behavior
	tXProcessor.Stream = sps[0].Stream
	tXProcessor.Kind = sps[0].Kind
	tXProcessor.Args = sps[0].Args
	tXProcessor.Command = sps[0].Command
	tXProcessor.Code = sps[0].Code
	tXProcessor.StdinAs = sps[0].StdinAs
	tXProcessor.StdoutAs = sps[0].StdoutAs
	tXProcessor.Description = sps[0].Description
	tXProcessor.OptionsCompositionTpl = sps[0].OptionsCompositionTpl
	tXProcessor.HasDoc = sps[0].HasDoc

	return tXProcessor, nil
}
