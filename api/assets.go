package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/core/models"
)

type AssetApiController struct {
	database *gorm.DB
	path     string
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

	a.database.Create(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	c.JSON(200, asset)
}

func (a *AssetApiController) DownloadOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+asset.Name+"\"")
	c.Data(200, "application/octet-stream", asset.Value)
}

func (a *AssetApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	err := c.BindJSON(&asset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	a.database.Save(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	a.database.Delete(&asset)

	c.JSON(204, "")
}

func (a *AssetApiController) ReplaceByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}
	tmpasset := models.Asset{}
	err := c.BindJSON(&tmpasset)
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

	a.database.Save(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}
