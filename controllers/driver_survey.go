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

type reponse struct {
	status int
	message string
}

var SurveyParams map[string]interface{}

// URLMapping ...
func (c *DriverSurveyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Submit", c.Submit)
}

// Submit ...
// @Title GET 方式提交数据
// @Description GET 方式提交数据
// @Param	sexuality	query	string 	true	"性别"
// @Param	age			query	string 	true	"年龄"
// @Param	email		query	string 	true	"邮箱"
// @Param	name		query	string 	true	"名称"
// @Param	region		query	string 	true	"地区/国家"
// @Param	game		query	string 	false	"游戏"
// @Param	suggest		query	string 	false	"建议"
// @Param	sign		query	string 	true	"签名"
// @Success 201 {int} models.DriverSurvey
// @Failure 200 :sign is empty
// @router / [get]
func (controller *DriverSurveyController) Submit() {
	var f models.SurveyForm
	var v models.DriverSurvey

	var sexuality string
	controller.Ctx.Input.Bind(&sexuality, "sexuality")  
	f.Sexuality = sexuality

	var age string
	controller.Ctx.Input.Bind(&age, "age")  
	f.Age = age

	var region string
	controller.Ctx.Input.Bind(&region, "region")
	f.Region = region

	var email string
	controller.Ctx.Input.Bind(&email, "email")  
	f.Email = email

	var name string
	controller.Ctx.Input.Bind(&name, "name")  
	f.Name = name

	var game string
	controller.Ctx.Input.Bind(&game, "game") 
	f.Game = game

	var suggest string
	controller.Ctx.Input.Bind(&suggest, "suggest")  
	f.Suggest = suggest

	var sign string
	controller.Ctx.Input.Bind(&sign, "sign")  
	f.Sign = sign

	// 检查签名
    validateSignResut := validateSign(f)
	if validateSignResut == false {
			controller.Data["json"] = map[string]interface{}{"status":"1001","message":"invalid sign"}
			controller.ServeJSON()
			return
	}

	// 校验数据数据
	validResult, Message := validateParams(f)	
	if validResult == false {
		controller.Data["json"] = map[string]interface{}{"status":"1002","message":Message}
		controller.ServeJSON()
		return
	} 

	v.Sexuality = f.Sexuality 
	v.Region 	= f.Region 	
	v.Email 	= f.Email	
	v.Name 		= f.Name	
	v.Suggest 	= f.Suggest 	
	v.Game 		= f.Game 		
	v.Age 		= f.Age 		
	
	// 数据验证通过 写入
	if _, err := models.AddDriverSurvey(&v); err == nil {
			controller.Ctx.Output.SetStatus(201)
			controller.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
	} else {
			controller.Data["json"] = err.Error()
	} 		

	controller.ServeJSON()
}

// Post ...
// @Title Post提交数据
// @Description 提交用户填写的调查数据
// @Success 201 {int} models.DriverSurvey
// @Param	body body models.SurveyForm true	"json格式，包含项年龄:Age,邮箱:Email,游戏:Game,名称:Name,地区/国家:Region,性别:Sexuality,签名:Sign,建议:Suggest"
// @Failure 200 body is empty
// @router / [post]
func (controller *DriverSurveyController) Post() {
	
	var v models.DriverSurvey
	var f models.SurveyForm

	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, &v); err == nil {
		json.Unmarshal(controller.Ctx.Input.RequestBody, &f)
		
		// 检查签名
    	validateSignResut := validateSign(f)
		if validateSignResut == false {
			controller.Data["json"] = map[string]interface{}{"status":"1001","message":"invalid sign"}
			controller.ServeJSON()
			return
		}
		
		// 校验数据数据
		validResult, Message := validateParams(f)	
		if validResult == false {
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
	raw ["sexuality"] = form.Sexuality
	raw ["region"] = form.Region
	raw ["email"] = form.Email
	raw ["name"] = form.Name
	raw ["age"] = form.Age
	raw ["suggest"] = form.Suggest
	raw ["game"] = form.Game

	sign := libs.MakeSignature(raw,beego.AppConfig.String("apiKey"))
	
	if sign != form.Sign {
		return false
	}

	return true
}

// 校验json里面的字段
func validateParams(form models.SurveyForm) (result bool, message string) {
	valid := validation.Validation{}
	
	valid.Required(form.Sexuality, "sexuality")
	valid.Required(form.Age, "age")
	valid.Required(form.Region, "region")
	valid.Required(form.Name, "name")
	valid.Required(form.Email, "email")

	valid.MaxSize(form.Region,255,"region")
	valid.MaxSize(form.Name,255,"name")
	valid.MaxSize(form.Email,255,"email")
	valid.MaxSize(form.Game,255,"game")

	valid.Email(form.Email,"email")
	
	log := logs.GetLogger()
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return false,err.Key + " " +err.Message
		}
	}

	return true,"pass"
}
