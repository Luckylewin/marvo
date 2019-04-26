package controllers

import (
	"encoding/json"
	"marvoAPI/models"
	"strconv"
	
	//"errors"
	//"strings"
	"github.com/astaxie/beego"
)

// DriverSurveyController operations for DriverSurvey
type DriverSurveyController struct {
	beego.Controller
}

// URLMapping ...
func (c *DriverSurveyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
}

// Post ...
// @Title Post
// @Description 提交用户填写的调查数据
// @Param	sexuality	sexuality 	models.sexuality true "调查表单:性别 取值[male|female]"
// @Param	region	region 	models.region	true "调查表单:地区/国家"
// @Param	email	email 	models.email	true "调查表单:邮箱"
// @Param	suggest	suggest models.suggest	false "调查表单:建议"
// @Param	age		age 	models.age		false "调查表单:年龄"
// @Param	game	game 	models.game		false "调查表单:游戏名称"
// @Success 201 成功提交数据到服务
// @Failure 403 sexuality is empty
// @Failure 403 region is empty
// @Failure 403 email is empty
// @router / [post]
func (c *DriverSurveyController) Post() {
	var v models.DriverSurvey
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddDriverSurvey(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
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