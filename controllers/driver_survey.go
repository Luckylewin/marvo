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
// @Param	mac		query	string 	true	"mac"
// @Param	sexuality	query	string 	true	"性别"
// @Param	age		query	string 	true	"年龄"
// @Param	email		query	string 	true	"邮箱"
// @Param	name		query	string 	true	"名称"
// @Param	area		query	string 	true	"地区/国家"
// @Param	types		query	string 	true	"型号"
// @Param	game		query	string 	false	"游戏"
// @Param	suggest		query	string 	false	"建议"
// @Param	facebook	query	string 	false	"facebook"
// @Param	sign		query	string 	true	"签名"
// @Success 201 {int} models.DriverSurvey
// @Failure 200 :sign is empty
// @router / [get]
func (c *DriverSurveyController) Submit() {
	var f models.SurveyForm
	var v models.DriverSurvey

	var mac string
	c.Ctx.Input.Bind(&mac, "mac")  
	f.Mac = mac

	var sexuality string
	c.Ctx.Input.Bind(&sexuality, "sexuality")  
	f.Sexuality = sexuality

	var age string
	c.Ctx.Input.Bind(&age, "age")  
	f.Age = age

	var area string
	c.Ctx.Input.Bind(&area, "area")
	f.Area = area

	var email string
	c.Ctx.Input.Bind(&email, "email")  
	f.Email = email

	var name string
	c.Ctx.Input.Bind(&name, "name")  
	f.Name = name

	var game string
	c.Ctx.Input.Bind(&game, "game") 
	f.Game = game

	var suggest string
	c.Ctx.Input.Bind(&suggest, "suggest")  
	f.Suggest = suggest

	var types string
	c.Ctx.Input.Bind(&types, "types")  
	f.Types = types
	
	var facebook string
	c.Ctx.Input.Bind(&facebook, "facebook")  
	f.Facebook = facebook

	var sign string
	c.Ctx.Input.Bind(&sign, "sign")  
	f.Sign = sign

	// 检查签名
 	validateSignResut := validateSign(f)
	if validateSignResut == false {
			c.Data["json"] = map[string]interface{}{"status":"1001","message":"invalid sign"}
			c.ServeJSON()
			return
	}

	// 校验数据数据
	validResult, Message := validateParams(f)	
	if validResult == false {
		c.Data["json"] = map[string]interface{}{"status":"1002","message":Message}
		c.ServeJSON()
		return
	} 

	v.Mac 		= f.Mac
	v.Sexuality = f.Sexuality 
	v.Area 		= f.Area 	
	v.Email 	= f.Email	
	v.Name 		= f.Name	
	v.Suggest 	= f.Suggest 	
	v.Game 		= f.Game 		
	v.Age 		= f.Age 		
	v.Types     = f.Types
	v.Facebook  = f.Facebook

	// 先判断一下这条数据是否存在
	if item, err := models.GetDriverSurveyByMac(v.Mac); err == nil {
			// 判断是否有更新数据
			if isNeedUpdate(f,item) {
				v.Id = item.Id
				if err := models.UpdateDriverSurveyById(&v); err == nil {
					c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
				} else {
					c.Data["json"] = map[string]interface{}{"status":"1003","message":"update fail"}
				}
			} else {
				c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
			}
	} else {
		// 不存在则新增
		if _, err := models.AddDriverSurvey(&v); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
		} else {
				c.Data["json"] = err.Error()
		} 			
	} 	

	c.ServeJSON()
}

// Post ...
// @Title Post提交数据
// @Description 提交用户填写的调查数据
// @Success 201 {int} models.DriverSurvey
// @Param	body body models.SurveyForm true	"json格式，包含项Mac地址:Mac,年龄:Age,邮箱:Email,游戏:Game,名称:Name,地区/国家:Area,性别:Sexuality,签名:Sign,建议:Suggest,型号:Types,脸书:facebook"
// @Failure 200 body is empty
// @router / [post]
func (c *DriverSurveyController) Post() {
	
	var v models.DriverSurvey
	var f models.SurveyForm

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		json.Unmarshal(c.Ctx.Input.RequestBody, &f)
		
	// 检查签名
    	validateSignResut := validateSign(f)
		if validateSignResut == false {
			c.Data["json"] = map[string]interface{}{"status":"1001","message":"invalid sign"}
			c.ServeJSON()
			return
		}
		
		// 校验数据数据
		validResult, Message := validateParams(f)	
		if validResult == false {
			c.Data["json"] = map[string]interface{}{"status":"1002","message":Message}
			c.ServeJSON()
			return
		} 

		// 先判断一下这条数据是否存在
		if item, err := models.GetDriverSurveyByMac(v.Mac); err == nil {
				// 判断是否有更新数据
				if isNeedUpdate(f,item) {
					v.Id = item.Id
					if err := models.UpdateDriverSurveyById(&v); err == nil {
						c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
					} else {
						c.Data["json"] = map[string]interface{}{"status":"1003","message":"update fail"}
					}
				} else {
					c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
				}
		} else {
				// 数据验证通过 写入
			if _, err := models.AddDriverSurvey(&v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = map[string]interface{}{"status":"1000","message":"success"}
			} else {
					c.Data["json"] = err.Error()
			}
		}

	} else {
		c.Data["json"] = err.Error()
	}

	c.ServeJSON()
}

// 判断是否有更新数据
func isNeedUpdate(form models.SurveyForm, object *models.DriverSurvey) bool {
	changeFlag := false
	if object.Area != form.Area  ||
	   object.Sexuality != form.Sexuality || 
	   object.Email != form.Email || 
	   object.Name != form.Name || 
	   object.Suggest != form.Suggest || 
	   object.Game != form.Game ||
	   object.Facebook != form.Facebook ||
	   object.Age != form.Age {
		return true
	}
	
	return changeFlag
}

// 校验Sign
func validateSign(form models.SurveyForm) bool {
	raw := make(map[string]interface{})
	raw ["mac"] = form.Mac
	raw ["sexuality"] = form.Sexuality
	raw ["area"] = form.Area
	raw ["email"] = form.Email
	raw ["name"] = form.Name
	raw ["age"] = form.Age
	raw ["suggest"] = form.Suggest
	raw ["game"] = form.Game
	raw ["types"] = form.Types
	raw ["facebook"] = form.Facebook
	
	sign := libs.MakeSignature(raw,beego.AppConfig.String("apiKey"))
	
	if sign != form.Sign {
		return false
	}

	return true
}

// 校验json里面的字段
func validateParams(form models.SurveyForm) (result bool, message string) {
	valid := validation.Validation{}
	
	valid.Required(form.Mac, "mac")
	valid.Required(form.Sexuality, "sexuality")
	valid.Required(form.Age, "age")
	valid.Required(form.Area, "area")
	valid.Required(form.Name, "name")
	valid.Required(form.Email, "email")
	valid.Required(form.Types, "types")

	valid.MaxSize(form.Mac,64,"mac")
	valid.MaxSize(form.Area,255,"area")
	valid.MaxSize(form.Name,255,"name")
	valid.MaxSize(form.Email,255,"email")
	valid.MaxSize(form.Game,255,"game")
	valid.MaxSize(form.Types,255,"types")
	valid.MaxSize(form.Facebook,255,"facebook")

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
