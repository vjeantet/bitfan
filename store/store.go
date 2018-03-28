package store

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/timshannon/bolthold"
	"github.com/vjeantet/bitfan/api/models"
	"github.com/vjeantet/bitfan/commons"
)

type Store struct {
	db              *bolthold.Store
	pipelineTmpPath string
	log             commons.Logger
}

func New(location string, log commons.Logger) (*Store, error) {
	database, err := bolthold.Open(filepath.Join(location, "bitfan.bolt.db"), 0666, nil)
	pipelineTmpPath := filepath.Join(location, "_pipelines")

	return &Store{db: database, log: log, pipelineTmpPath: pipelineTmpPath}, err
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) PipelineTmpPath(tPipeline *models.Pipeline) string {
	uidString := fmt.Sprintf("%s_%d", tPipeline.Uuid, time.Now().Unix())
	pipelinePath := filepath.Join(s.pipelineTmpPath, uidString)

	if tPipeline.Playground == true {
		pipelinePath = filepath.Join(os.TempDir(), "bitfan-playground", tPipeline.Uuid)
	}

	os.MkdirAll(pipelinePath, os.ModePerm)
	return pipelinePath
}

func (s *Store) preparePipelineAssetExecutionStage(cwd string, tAssets *[]models.Asset) (string, error) {
	var entrypointLocation string
	for _, asset := range *tAssets {
		dest := filepath.Join(cwd, asset.Name)
		dir := filepath.Dir(dest)
		os.MkdirAll(dir, os.ModePerm)
		s.log.Debugf("configuration %s stored to %s", asset.Name, cwd)
		if err := ioutil.WriteFile(dest, asset.Value, 07770); err != nil {
			return "", err
		}

		if asset.Type == models.ASSET_TYPE_ENTRYPOINT {
			entrypointLocation = filepath.Join(cwd, asset.Name)
		}
	}

	return entrypointLocation, nil
}

func (s *Store) PreparePipelineExecutionStage(tPipeline *models.Pipeline) (string, error) {
	//Save assets to cwd
	cwd := s.PipelineTmpPath(tPipeline)

	s.log.Debugf("configuration %s storage to %s", tPipeline.Uuid, cwd)

	// If Playground With Base
	if tPipeline.PlaygroundBaseUUID != "" {
		// Find all assets from Base Pipeline
		assets, err := s.FindAssetsByPipelineUUID(tPipeline.PlaygroundBaseUUID)
		if err != nil {
			return "", err
		}

		// Save them with playground pipeline
		_, err = s.preparePipelineAssetExecutionStage(cwd, &assets)
		if err != nil {
			return "", err
		}
	}

	entrypointLocation, err := s.preparePipelineAssetExecutionStage(cwd, &tPipeline.Assets)
	if err != nil {
		return "", err
	}
	tPipeline.ConfigLocation = entrypointLocation

	if tPipeline.ConfigLocation == "" {
		return "", fmt.Errorf("missing entrypoint for pipeline %s", tPipeline.Uuid)
	}

	s.log.Debugf("configuration %s pipeline %s ready to be loaded", tPipeline.Uuid, tPipeline.ConfigLocation)
	return tPipeline.ConfigLocation, nil
}

// CopyTo writes the raw database's content to given io.Writer
func (s *Store) CopyTo(w io.Writer) (int, error) {
	size := 0
	err := s.db.Bolt().View(func(tx *bolt.Tx) error {
		size = int(tx.Size())
		_, err := tx.WriteTo(w)
		return err
	})
	return size, err
}
