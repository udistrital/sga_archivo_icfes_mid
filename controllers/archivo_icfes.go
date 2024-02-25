package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_archivo_icfes_mid/services"
	"github.com/udistrital/utils_oas/requestresponse"
)

// ArchivoIcfesController ...
type ArchivoIcfesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArchivoIcfesController) URLMapping() {
	c.Mapping("PostArchivoIcfes", c.PostArchivoIcfes)
}

// PostArchivoIcfes ...
// @Title PostArchivoIcfes
// @Description Agregar ArchivoIcfes
// @Param id path int	true "el id del periodo"
// @Param   archivo_icfes	formData  file	true   "body Agregar ArchivoIcfes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /:id [post]
func (c *ArchivoIcfesController) PostArchivoIcfes() {
	periodo_id := c.Ctx.Input.Param(":id")
	archivoIcfes, _, err := c.GetFile("archivo_icfes")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	} else {
		respuesta := services.ArchivosIcfes(periodo_id, archivoIcfes)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
	}
	c.ServeJSON()
}
