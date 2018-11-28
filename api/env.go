package api

import (
	"fmt"
	"os"
	"strings"

	"bitfan/api/models"
	"bitfan/core"
	"github.com/gin-gonic/gin"
	gouuid "github.com/nu7hatch/gouuid"
)

type EnvApiController struct {
	path string
}

func (p *EnvApiController) Find(c *gin.Context) {
	envs := core.Storage().FindEnvs()

	for i := range envs {
		if envs[i].Secret == true {
			envs[i].Value = "********"
		}
	}

	envs = append([]models.Env{
		{Name: "PATH", Value: os.Getenv("PATH")},
		{Name: "BF_COMMONS_PATH", Value: os.Getenv("BF_COMMONS_PATH")},
		{Name: "BF_DATA_PATH", Value: os.Getenv("BF_DATA_PATH")},
		{Name: "BF_HTTPD_ADDR", Value: os.Getenv("BF_HTTPD_ADDR")},
	}, envs...)

	c.JSON(200, envs)
}

func (p *EnvApiController) Create(c *gin.Context) {
	var varenv models.Env
	err := c.BindJSON(&varenv)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	varenv.Name = strings.TrimSpace(varenv.Name)

	if varenv.Name == "" {
		c.JSON(400, models.Error{Message: "env name can not be empty"})
		return
	}

	//TODO: find one by name, if err --> return err (duplicate)

	if varenv.Uuid == "" {
		uid, _ := gouuid.NewV4()
		varenv.Uuid = uid.String()
	}

	core.Storage().CreateEnv(&varenv)
	if err := os.Setenv(varenv.Name, varenv.Value); err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}
	c.Redirect(302, fmt.Sprintf("/%s/env/%s", p.path, varenv.Uuid))
}

func (p *EnvApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	varenv, err := core.Storage().FindOneEnvByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}
	if varenv.Secret == true {
		varenv.Value = "********"
	}
	c.JSON(200, varenv)
}

func (p *EnvApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	varenv, err := core.Storage().FindOneEnvByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	core.Storage().DeleteEnv(&varenv)
	if err := os.Unsetenv(varenv.Name); err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}
	c.JSON(204, "")
}
