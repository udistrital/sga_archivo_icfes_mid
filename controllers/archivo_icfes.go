package controllers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_archivo_icfes/models"
	"github.com/udistrital/sga_mid_archivo_icfes/services"
	"github.com/udistrital/utils_oas/errorhandler"
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
	defer errorhandler.HandlePanic(&c.Controller)

	periodo_id := c.Ctx.Input.Param(":id")
	ArchivoIcfes := "Archivo procesado"
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	fmt.Println("name", c.GetString("name"))
	fmt.Println("periodo", periodo_id)
	multipartFile, _, err := c.GetFile("archivo_icfes")
	if err != nil {
		services.ManejoError(&alerta, &alertas, "err reading multipartFile", err)
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	file, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		services.ManejoError(&alerta, &alertas, "err reading file", err)
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	lines := strings.Split(strings.Replace(string(file), "\r\n", "\n", -1), "\n")
	//Probando que el archivo tenga el contenido necesario
	if len(lines) < 2 {
		services.ManejoError(&alerta, &alertas, "err in file content")
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	testHeaderFile := strings.Split(lines[0], ",")[0]
	if testHeaderFile != "CODREGSNP" {
		services.ManejoError(&alerta, &alertas, "err in file content")
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	lines = lines[1:] // remove first element
	evaluacionesInscripcion := make([]map[string]interface{}, 0)
	detallesEvaluacion := make([]map[string]interface{}, 0)
	for _, line := range lines {
		// 0 código ICFEs del estudianate
		// 1 para nombre del estudiante
		// 11 PLC Español
		// 12 PMA Matematicas
		// 13 PSC Sociela
		// 14 PCN Ciencia Naturales
		// 15 PIN Ingles
		recordFields := strings.Split(line, ",")
		if len(recordFields) > 1 {
			if !services.PeticionInscripciones(recordFields, periodo_id, &alerta, &alertas, &detallesEvaluacion, &evaluacionesInscripcion) {
				c.Ctx.Output.SetStatus(404)
				c.ServeJSON()
			}
		}
	}

	services.ManejoPeticiones(evaluacionesInscripcion, &alerta, &alertas, detallesEvaluacion, ArchivoIcfes)

	c.Ctx.Output.SetStatus(200)
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
