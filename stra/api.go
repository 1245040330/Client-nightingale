package stra

import (
	"fmt"
	"github.com/didi/nightingale/src/model"
	"github.com/didi/nightingale/src/toolkits/str"
	"reflect"
	"time"
)


func NewApiCollect( step int, name,api string, modTime time.Time) *model.ApiCollect {
	return &model.ApiCollect{
		CollectType:   "api",
		Name:        name,
		Step:          step,
		Api:          api,
		LastUpdated:   modTime,
	}
}
type Res struct {
	Dat []model.ApiCollect `json:"dat"`
	Err string        `json:"err"`
}
func copyPoint(m *model.ApiCollect) *model.ApiCollect{
	vt := reflect.TypeOf(m).Elem()
	fmt.Println(vt)
	newoby := reflect.New(vt)
	newoby.Elem().Set(reflect.ValueOf(m).Elem())
	return newoby.Interface().(*model.ApiCollect)
}
func GetApiCollects() map[string]*model.ApiCollect {
	apis := make(map[string]*model.ApiCollect)
	if StraConfig.Enable {
		apis =Collect.GetApi()

		for _, p := range apis {
			tagsMap := str.DictedTagstring("api")
			tagsMap["api"] = p.Api

			p.Comment = str.SortedTags(tagsMap)
		}
	}
	return apis
}

