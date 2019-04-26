package controllers

import (
	"encoding/json"
	"marvoAPI/models"
	"strconv"
	
	//"errors"
	//"strings"
	"github.com/astaxie/beego"
)

// 数据收集相关API
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
// @Success 201 {int} models.DriverSurvey
// @Param	body body models.SurveyForm true	"json数据"
// @Failure 403 body is empty
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