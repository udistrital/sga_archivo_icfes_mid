// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/sga_mid_archivo_icfes/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/astaxie/beego"
)

func init() {

	beego.ErrorController(&errorhandler.ErrorHandlerController{})
	
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/icfes",
			beego.NSInclude(
				&controllers.ArchivoIcfesController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
