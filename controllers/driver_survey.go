package controllers

import (
	"encoding/json"
	"marvoAPI/models"
	"strconv"
	
	//"errors"
	//"strings"
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
	c.Mapping("GetOne", c.GetOne)
}

// Post ...
// @Title Post
// @Description 提交用户填写的调查数据
// @Success 201 {int} models.DriverSurvey
// @Param	body body models.SurveyForm true	"json数据"
// @Failure 403 body is empty
// @router / [post]
func (controller *DriverSurveyController) Post() {
	log := logs.GetLogger()
    
	var v models.DriverSurvey
	var f models.SurveyForm

	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, &v); err == nil {
		// 校验数据数据
		json.Unmarshal(controller.Ctx.Input.RequestBody, &f)
		validResult, Message := checkFormParams(f)

		log.Println("收到的年龄",f.Age)
		log.Println("收到的名称",f.Name)
		log.Println("收到的地区",f.Region)
		log.Println("收到的邮箱",f.Email)
		log.Println("收到的性别",f.Sexuality)

		log.Println(f)
		log.Println(validResult,Message)

		if validResult == false {
			controller.Ctx.Output.SetStatus(400)
			controller.Data["json"] = map[string]interface{}{"status":"1001","test":Message}
		} else {
			if _, err := models.AddDriverSurvey(&v); err == nil {
				controller.Ctx.Output.SetStatus(201)
				controller.Data["json"] = v
			} else {
				controller.Data["json"] = err.Error()
			}
		}
	} else {
		controller.Data["json"] = err.Error()
	}

	controller.ServeJSON()
}

func checkFormParams(form models.SurveyForm) (result bool, message string) {
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

// GetOne ...
// @Title Get One
// @Description get DriverSurvey by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DriverSurvey
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DriverSurveyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetDriverSurveyById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}