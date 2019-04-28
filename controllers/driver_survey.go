package controllers

import (
	"encoding/json"
	"marvoAPI/models"
	"marvoAPI/libs"
	"github.com/astaxie/beego"
 	"github.com/astaxie/beego/logs"
 	"github.com/astaxie/beego/validation"
)

// 数据收集相关API
type DriverSurveyController struct {
	beego.Controller
}

var SurveyParams map[string]interface{}

// URLMapping ...
func (c *DriverSurveyController) URLMapping() {
	c.Mapping("Post", c.Post)
}


// Post ...
// @Title Post提交数据
// @Description 提交用户填写的调查数据
// @Success 201 {int} models.DriverSurvey
// @Param	body body models.SurveyForm true	"json格式，包含项年龄:Age,邮箱:Email,游戏:Game,名称:Name,地区/国家:Region,性别:Sexuality,签名:Sign,建议:Suggest"
// @Failure 400 body is empty
// @router / [post]
func (controller *DriverSurveyController) Post() {
	
	var v models.DriverSurvey
	var f models.SurveyForm

	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, &v); err == nil {
		json.Unmarshal(controller.Ctx.Input.RequestBody, &f)
		
		// 检查签名
    		validateSignResut := validateSign(f)
		if validateSignResut == false {
			controller.Ctx.Output.SetStatus(400)
			controller.Data["json"] = map[string]interface{}{"status":"1001","message":"invalid sign"}
			controller.ServeJSON()
			return
		}

		// 校验数据数据
		validResult, Message := validateParams(f)	
		if validResult == false {
			controller.Ctx.Output.SetStatus(400)
			controller.Data["json"] = map[string]interface{}{"status":"1002","message":Message}
			controller.ServeJSON()
			return
		} 

		// 数据验证通过 写入
		if _, err := models.AddDriverSurvey(&v); err == nil {
				controller.Ctx.Output.SetStatus(201)
				controller.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
		} else {
				controller.Data["json"] = err.Error()
		}

	} else {
		controller.Data["json"] = err.Error()
	}

	controller.ServeJSON()
}

// 校验Sign
func validateSign(form models.SurveyForm) bool {
	raw := make(map[string]interface{})
	raw ["Sexuality"] = form.Sexuality
	raw ["Region"] = form.Region
	raw ["Email"] = form.Email
	raw ["Name"] = form.Name
	raw ["Age"] = form.Age
	raw ["Suggest"] = form.Suggest
	raw ["Game"] = form.Game

	sign := libs.MakeSignature(raw,beego.AppConfig.String("apiKey"))

	if sign != form.Sign {
		return false
	}

	return true
}

// 校验json里面的字段
func validateParams(form models.SurveyForm) (result bool, message string) {
	valid := validation.Validation{}
	
	valid.Required(form.Sexuality, "Sexuality")
	valid.Required(form.Age, "Age")
	valid.Required(form.Region, "Region")
	valid.Required(form.Name, "Name")
	valid.Required(form.Email, "Email")

	valid.MaxSize(form.Region,255,"Region")
	valid.MaxSize(form.Name,255,"Name")
	valid.MaxSize(form.Email,255,"Email")
	valid.MaxSize(form.Game,255,"Game")

	valid.Email(form.Email,"Email")
	
	log := logs.GetLogger()
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return false,err.Key + " " +err.Message
		}
	}

	return true,"pass"
}
