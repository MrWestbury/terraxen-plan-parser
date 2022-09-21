package provider

import (
	"log"
	"net/http"
	"strings"
	"terraform-plan-parser/internals"

	"github.com/gin-gonic/gin"
)

func (siteSvc *SiteService) Changes(c *gin.Context) {

	items := make([]internals.ResourceChange, 0)

	for _, item := range siteSvc.plan.ResourceChanges {
		actionStr := strings.Join(item.Changes.Actions, "/")
		if actionStr == "no-op" {
			continue
		}

		items = append(items, item)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "Changes",
		"page":      "changes.tmpl",
		"IsChanges": true,
		"data":      items,
	})
}

func (siteSvc *SiteService) GetChangeList(c *gin.Context) {
	modules := make([]TreeItem, 0)
	moduleIds := map[string]bool{}
	result := make([]TreeItem, 0)

	for _, item := range siteSvc.changeIndex {
		moduleId := hash(item.ModuleAddress)
		mod, exists := siteSvc.stateIndex[moduleId]
		parentId := "#"
		if exists {

			parentId = moduleId
			if _, exists2 := moduleIds[moduleId]; !exists2 {
				modItem := TreeItem{
					Id:          moduleId,
					Parent:      mod.Parent,
					DisplayName: mod.Name,
					Type:        "Module",
				}
				modules = append(modules, modItem)
				moduleIds[moduleId] = true
			}
		}

		respItem := TreeItem{
			Id:          item.Id,
			Parent:      parentId,
			DisplayName: item.Name,
			Type:        item.Change,
			LiAttr: map[string]string{
				"data-addr": item.Address,
			},
		}
		result = append(result, respItem)
	}

	result = append(modules, result...)

	c.IndentedJSON(http.StatusOK, result)
}

func (siteSvc *SiteService) GetChangeResource(c *gin.Context) {
	addr := c.Param("address")

	res, exists := siteSvc.changeIndex[addr]
	if !exists {
		resp := map[string]string{
			"message": "change resource not found",
		}

		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := map[string]interface{}{
		"resource":   res.Resource,
		"attributes": siteSvc.changeAttrIndex[res.Address],
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func generateChangeAttrIndex(changes []internals.Attribute) ChangeAttrIndex {
	index := ChangeAttrIndex{}

	for _, chg := range changes {
		_, ok := index[chg.Resource]
		if !ok {
			index[chg.Resource] = make([]string, 0)
		}
		index[chg.Resource] = append(index[chg.Resource], chg.Attribute)
	}
	log.Printf("Have %d atrributes", len(index))
	return index
}

func generateChangeIndex(changes []internals.ResourceChange) ChangeIndex {
	index := ChangeIndex{}

	for _, chg := range changes {
		item := ChangeItem{
			Id:            hash(chg.Address),
			Name:          chg.Name,
			Address:       chg.Address,
			ModuleAddress: chg.ModuleAddress,
			Change:        strings.Join(chg.Changes.Actions, "/"),
			Resource:      chg,
		}
		index[item.Id] = item
	}

	return index
}
