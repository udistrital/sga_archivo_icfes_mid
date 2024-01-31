package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/udistrital/sga_mid_archivo_icfes/models"
)

func ManejoError(alerta *models.Alert, alertas *[]interface{}, mensaje string, err ...error) {
	var msj string
	if len(err) > 0 && err[0] != nil {
		msj = mensaje + err[0].Error()
	} else {
		msj = mensaje
	}
	*alertas = append(*alertas, msj)
	(*alerta).Body = *alertas
	(*alerta).Type = "error"
	(*alerta).Code = "400"
}

func AsignacionNotas(criteriosRes []map[string]interface{}, porcentajes map[string]interface{}, aspirante_puntajes map[string]interface{}, detallesEvaluacion *[]map[string]interface{}, evaluacionesInscripcion *[]map[string]interface{}, proyecto_inscripcion interface{}, aspirante_codigo_icfes *string, inscripcion map[string]interface{}) {
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
			*evaluacionesInscripcion = append(*evaluacionesInscripcion, map[string]interface{}{
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

			*detallesEvaluacion = append(*detallesEvaluacion, map[string]interface{}{
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
