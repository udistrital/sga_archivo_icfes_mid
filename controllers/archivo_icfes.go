package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"io/ioutil"
	"sga_mid_archivo_icfes/models"
	"strconv"
	"strings"
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
	ArchivoIcfes := "Archivo procesado"
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	fmt.Println("name", c.GetString("name"))
	fmt.Println("periodo", periodo_id)
	multipartFile, _, err := c.GetFile("archivo_icfes")
	if err != nil {
		fmt.Println("err reading multipartFile", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	file, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		fmt.Println("err reading file", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	lines := strings.Split(strings.Replace(string(file), "\r\n", "\n", -1), "\n")
	// fmt.Println(lines)
	//Probando que el archivo tenga el contenido necesario
	if len(lines) < 2 {
		fmt.Println("err in file content lentg")
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err in file content")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	testHeaderFile := strings.Split(lines[0], ",")[0]
	if testHeaderFile != "CODREGSNP" {
		fmt.Println("err in file content")
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err in file content")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	// lines = lines[1:] // remove first element
	lines = lines[1:] // remove first element
	evaluacionesInscripcion := make([]map[string]interface{}, 0)
	detallesEvaluacion := make([]map[string]interface{}, 0)
	// fmt.Println(lines, len(lines))
	for _, line := range lines {
		// 0 código ICFEs del estudianate
		// 1 para nombre del estudiante
		// 11 PLC Español
		// 12 PMA Matematicas
		// 13 PSC Sociela
		// 14 PCN Ciencia Naturales
		// 15 PIN Ingles
		recordFields := strings.Split(line, ",")
		// fmt.Println("Separa")
		// fmt.Println(recordFields)
		if len(recordFields) > 1 {
			aspirante_codigo_icfes := recordFields[0]
			aspirante_nombre := recordFields[1]
			aspirante_puntajes := map[string]interface{}{
				"PLC": recordFields[11],
				"PMA": recordFields[12],
				"PSC": recordFields[13],
				"PCN": recordFields[14],
				"PIN": recordFields[15],
			}
			fmt.Println("line", aspirante_codigo_icfes, aspirante_nombre, aspirante_puntajes)
			// traer data de la inscripcion o inscripciones
			// fmt.Println("url","http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId__Activo:true,InscripcionId__EstadoInscripcionId__Id:1,InscripcionId__PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes)
			var inscripcionesRes []map[string]interface{}
			errInscripciones := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId.Activo:true,InscripcionId.EstadoInscripcionId.Id:1,InscripcionId.PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes, &inscripcionesRes)
			if errInscripciones != nil {
				alertas = append(alertas, errInscripciones)
				alerta.Body = alertas
				alerta.Type = "error"
				alerta.Code = "400"
				c.ServeJSON()
			} else {
				// fmt.Println("inscripciones", len(inscripcionesRes), inscripcionesRes)
				// fmt.Println("inscripciones", len(inscripcionesRes))
				for _, inscripcionTemp := range inscripcionesRes {
					/// fmt.Println("inscripcionTemp", inscripcionTemp)
					if inscripcionTemp["InscripcionId"] != nil {
						// Extrae info de la inscripcion para saber el proyecto y la persona
						inscripcion := inscripcionTemp["InscripcionId"].(map[string]interface{})
						proyecto_inscripcion := inscripcion["ProgramaAcademicoId"]
						// fmt.Println("ProgramaAcademicoId", proyecto_inscripcion)
						// cargar criterios de admisión con el proyecto dependiendo de la inscripcion
						var criteriosRes []map[string]interface{}
						// fmt.Println("url criterios", "http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/requisito_programa_academico?limit=0&query=Activo:true,RequisitoId__Activo:true,PeriodoId:"+periodo_id+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", proyecto_inscripcion))
						errCriterios := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/requisito_programa_academico?limit=0&query=Activo:true,RequisitoId__Activo:true,PeriodoId:"+periodo_id+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", proyecto_inscripcion.(float64)), &criteriosRes)
						if errCriterios != nil {
							alertas = append(alertas, errCriterios)
							alerta.Body = alertas
							alerta.Type = "error"
							alerta.Code = "400"
							c.ServeJSON()
						} else {
							// fmt.Println("criterios", criteriosRes)
							// formatdata.JsonPrint(criteriosRes)
							// si existe criterios para el proyecto del aspirante revisar desde aqui
							//Revisar el for no es necesarios ps ya se maneja un solo criterio para los subcriterios
							var porcentajes map[string]interface{}
							for _, criterioTemp := range criteriosRes {
								if criterioTemp["RequisitoId"] != nil {
									// fmt.Println("criterio")
									// formatdata.JsonPrint(criterioTemp["PorcentajeEspecifico"])

									if err := json.Unmarshal([]byte(criterioTemp["PorcentajeEspecifico"].(string)), &porcentajes); err != nil {
										panic(err)
									}

									// fmt.Println("salee")
									// formatdata.JsonPrint(aspirante_puntajes)
									// Calculo de notas por su respectivo area y procentaje definido por carrera
									//Matematicas
									NotaMatematicas, _ := strconv.ParseFloat(aspirante_puntajes["PMA"].(string), 64)
									PorcentajeMatematicas, _ := strconv.ParseFloat(porcentajes["Area1"].(string), 64)
									TotalMatematicas := NotaMatematicas * (PorcentajeMatematicas / 100)
									//Ciencias Naturales
									NotaNaturales, _ := strconv.ParseFloat(aspirante_puntajes["PCN"].(string), 64)
									PorcentajeNaturales, _ := strconv.ParseFloat(porcentajes["Area2"].(string), 64)
									TotalNaturales := NotaNaturales * (PorcentajeNaturales / 100)
									//Español
									NotaEspañol, _ := strconv.ParseFloat(aspirante_puntajes["PLC"].(string), 64)
									PorcentajeEspañol, _ := strconv.ParseFloat(porcentajes["Area3"].(string), 64)
									TotalEspañol := NotaEspañol * (PorcentajeEspañol / 100)

									//Sociales
									NotaSociales, _ := strconv.ParseFloat(aspirante_puntajes["PSC"].(string), 64)
									PorcentajeSociales, _ := strconv.ParseFloat(porcentajes["Area4"].(string), 64)
									TotalSociales := NotaSociales * (PorcentajeSociales / 100)

									//Ingles
									NotaIngles, _ := strconv.ParseFloat(aspirante_puntajes["PIN"].(string), 64)
									PorcentajeIngles, _ := strconv.ParseFloat(porcentajes["Area5"].(string), 64)
									TotalIngles := NotaIngles * (PorcentajeIngles / 100)
									fmt.Println("Nota Matematematicas:", TotalMatematicas, "Nota Ciencias Naturales:", TotalNaturales, "Nota Español:", TotalEspañol, "Nota Sociales:", TotalSociales, "Nota Ingles:", TotalIngles)
									// formatdata.JsonPrint(Porcentaje)
									// fmt.Println("inscripcion", aspirante_codigo_icfes, aspirante_puntajes[criterio["CodigoAbreviacion"].(string)])
									notaFinal := TotalMatematicas + TotalNaturales + TotalEspañol + TotalSociales + TotalIngles
									// notaFinal, _ := strconv.ParseFloat(aspirante_puntajes[criterio["CodigoAbreviacion"].(string)].(string), 64)
									evaluacionesInscripcion = append(evaluacionesInscripcion, map[string]interface{}{
										"InscripcionId": inscripcion["Id"],
										"NotaFinal":     notaFinal,
										"Activo":        true,
									})

									area1 := fmt.Sprintf("%v", TotalMatematicas)
									area2 := fmt.Sprintf("%v", TotalNaturales)
									area3 := fmt.Sprintf("%v", TotalEspañol)
									area4 := fmt.Sprintf("%v", TotalSociales)
									area5 := fmt.Sprintf("%v", TotalIngles)
									pma := fmt.Sprintf("%v", aspirante_puntajes["PMA"])
									pcn := fmt.Sprintf("%v", aspirante_puntajes["PCN"])
									plc := fmt.Sprintf("%v", aspirante_puntajes["PLC"])
									pcs := fmt.Sprintf("%v", aspirante_puntajes["PSC"])
									pin := fmt.Sprintf("%v", aspirante_puntajes["PIN"])

									requestBod := "{\"Puntajes\":{\"PMA\": \"" + pma + "\", \"PCN\": \"" + pcn + "\", \"PLC\":\"" + plc + "\",\"PSC\": \"" + pcs + "\", \"PIN\": \"" + pin + "\" },\"Notas\":{\"Area1\": \"" + area1 + "\",\"Area2\": \"" + area2 + "\",\"Area3\": \"" + area3 + "\",\"Area4\": \"" + area4 + "\",\"Area5\": \"" + area5 + "\"}}"

									detallesEvaluacion = append(detallesEvaluacion, map[string]interface{}{
										"EvaluacionInscripcionId":      "viene del anterior",
										"RequisitoProgramaAcademicoId": map[string]interface{}{"Id": criteriosRes[0]["Id"].(float64)},
										"NotaRequisito":                notaFinal,
										"Activo":                       true,
										"DetalleCalificacion":          requestBod,
									})
								} else {
									fmt.Println("no hay criterios para proyecto", proyecto_inscripcion, "para inscripcion", aspirante_codigo_icfes)
								}

							}

						}

					} else {
						fmt.Println("no hay inscripciones para ", aspirante_codigo_icfes)

					}

				}

			}
		}

	}

	formatdata.JsonPrint(evaluacionesInscripcion)

	for i, postevaluacion := range evaluacionesInscripcion {
		var resultadoevaluacion map[string]interface{}
		errPostevaluacion := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/evaluacion_inscripcion", "POST", &resultadoevaluacion, postevaluacion)
		if resultadoevaluacion["Type"] == "error" || errPostevaluacion != nil || resultadoevaluacion["Status"] == "404" || resultadoevaluacion["Message"] != nil {
			alertas = append(alertas, resultadoevaluacion)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			detallesEvaluacion[i]["EvaluacionInscripcionId"] = map[string]interface{}{"Id": resultadoevaluacion["Id"].(float64)}

			// alertas = append(alertas, resultadoevaluacion)
		}
	}
	formatdata.JsonPrint(detallesEvaluacion)
	for _, postdetalle := range detallesEvaluacion {
		var resultadodetalle map[string]interface{}
		errPostedetalle := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/detalle_evaluacion", "POST", &resultadodetalle, postdetalle)
		if resultadodetalle["Type"] == "error" || errPostedetalle != nil || resultadodetalle["Status"] == "404" || resultadodetalle["Message"] != nil {
			alertas = append(alertas, resultadodetalle)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {

			alertas = append(alertas, ArchivoIcfes)
		}

	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
