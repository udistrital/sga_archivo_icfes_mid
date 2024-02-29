package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_archivo_icfes/controllers:ArchivoIcfesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_archivo_icfes/controllers:ArchivoIcfesController"],
        beego.ControllerComments{
            Method: "PostArchivoIcfes",
            Router: "/archivos/:id_periodo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
