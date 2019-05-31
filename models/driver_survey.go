package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"github.com/astaxie/beego/orm"
)

type SurveyForm struct {
	Mac 	  string `orm:"column(mac)" description:"mac地址"`
	Sexuality string `orm:"column(sexuality)" description:"性别(female|male)"`
	Area      string `orm:"column(area);size(255)" description:"地区/国家"`
	Email     string `orm:"column(email);size(255)" description:"邮箱"`
	Name      string `orm:"column(name);size(255);" description:"年龄"`
	Age       string `orm:"column(age);size(50);" description:"年龄"`
	Suggest   string `orm:"column(suggest);size(1000);null" description:"建议"`
	Game      string `orm:"column(game);size(255);null" description:"游戏"`
	Types     string `orm:"column(types);size(255);null" description:"型号"`
	Facebook  string `orm:"column(facebook);size(255);null" description:"facebook"`
	Sign      string `orm:"column(Sign);size(255);null" description:"签名"`
}

type DriverSurvey struct {
	Id        int    `orm:"column(id);auto"`
	Mac 	  string `orm:"column(mac)" description:"mac地址"`
	Sexuality string `orm:"column(sexuality)" description:"性别(female|male)"`
	Area      string `orm:"column(area);size(255)" description:"地区/国家"`
	Email     string `orm:"column(email);size(255)" description:"邮箱"`
	Name      string `orm:"column(name);size(255);" description:"年龄"`
	Age       string `orm:"column(age);size(50);null" description:"年龄"`
	Suggest   string `orm:"column(suggest);size(1000);null" description:"建议"`
	Game      string `orm:"column(game);size(255);null" description:"游戏"`
	Types     string `orm:"column(types);size(255);null" description:"型号"`
	Facebook  string `orm:"column(facebook);size(255);null" description:"facebook"`
	CreatedAt int64  `orm:"column(created_at);null" description:"创建时间"`
	UpdatedAt int64  `orm:"column(updated_at);null" description:"更新时间"`
}

func (t *DriverSurvey) TableName() string {
	return "driver_surveys"
}

func init() {
	orm.RegisterModel(new(DriverSurvey))
}

// AddDriverSurvey insert a new DriverSurvey into database and returns
// last inserted Id on success.
func AddDriverSurvey(m *DriverSurvey) (id int64, err error) {
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDriverSurveyById retrieves DriverSurvey by Id. Returns error if
// Id doesn't exist
func GetDriverSurveyById(id int) (v *DriverSurvey, err error) {
	o := orm.NewOrm()
	v = &DriverSurvey{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetDriverSurveyByMac retrieves DriverSurvey by Mac. Returns true if Mac doesn't exist
func GetDriverSurveyByMac(mac string) (*DriverSurvey,error) {
	survey := new(DriverSurvey)
	o := orm.NewOrm()
    err := o.QueryTable("driver_surveys").Filter("mac", mac).One(survey)
    if err != nil {
        return nil, err
    }
    return survey, nil
}

// GetAllDriverSurvey retrieves all DriverSurvey matches certain condition. Returns empty list if
// no records exist
func GetAllDriverSurvey(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DriverSurvey))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []DriverSurvey
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateDriverSurvey updates DriverSurvey by Id and returns error if
// the record to be updated doesn't exist
func UpdateDriverSurveyById(m *DriverSurvey) (err error) {
	o := orm.NewOrm()
	v := DriverSurvey{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		m.UpdatedAt = time.Now().Unix()
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDriverSurvey deletes DriverSurvey by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDriverSurvey(id int) (err error) {
	o := orm.NewOrm()
	v := DriverSurvey{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DriverSurvey{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
