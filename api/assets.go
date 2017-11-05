package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/core/models"
)

type AssetApiController struct {
	path string
}

func (a *AssetApiController) Create(c *gin.Context) {
	var asset models.Asset
	err := c.BindJSON(&asset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	uid, _ := uuid.NewV4()
	asset.Uuid = uid.String()
	asset.Size = len(asset.Value)

	n := 512
	if len(asset.Value) < 512 {
		n = len(asset.Value)
	}
	asset.ContentType = http.DetectContentType(asset.Value[:n])

	core.Storage().CreateAsset(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset, err := core.Storage().FindOneAssetByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	c.JSON(200, asset)
}

func (a *AssetApiController) DownloadOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset, err := core.Storage().FindOneAssetByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+asset.Name+"\"")
	c.Data(200, "application/octet-stream", asset.Value)
}

func (a *AssetApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset, err := core.Storage().FindOneAssetByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	err = c.BindJSON(&asset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	core.Storage().SaveAsset(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset, err := core.Storage().FindOneAssetByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	core.Storage().DeleteAsset(&asset)

	c.JSON(204, "")
}

func (a *AssetApiController) ReplaceByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset, err := core.Storage().FindOneAssetByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	tmpasset := models.Asset{}
	err = c.BindJSON(&tmpasset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	asset.Name = tmpasset.Name
	asset.Value = tmpasset.Value
	asset.Size = len(asset.Value)
	n := 512
	if len(asset.Value) < 512 {
		n = len(asset.Value)
	}
	asset.ContentType = http.DetectContentType(asset.Value[:n])

	core.Storage().SaveAsset(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}
