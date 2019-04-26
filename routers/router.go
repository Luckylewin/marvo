// @APIVersion 1.0.0
// @Title Marvo Game API
// @Description Marvo Game 相关的API
// @Contact 876505905@qq.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"marvoAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/driver_survey",
			beego.NSInclude(
				&controllers.DriverSurveyController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
