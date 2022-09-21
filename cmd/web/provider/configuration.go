package provider

import (
	"fmt"
	"net/http"
	"terraform-plan-parser/internals"

	"github.com/gin-gonic/gin"
)

func (siteSvc *SiteService) generateProviderIndex(data internals.ProviderConfigs) {

	providersRoot := &TreeItemAsync{
		Id:          "providers",
		DisplayName: "Providers",
		Parent:      "#",
		Type:        "Folder",
		Children:    true,
	}
	siteSvc.configIndex.AddItem(providersRoot, nil)
	for provId, prov := range data {
		item := &TreeItemAsync{
			Id:          hash(provId),
			DisplayName: provId,
			Parent:      "providers",
			Type:        "Provider",
			Children:    false,
		}
		siteSvc.configIndex.AddItem(item, prov)
	}

}

func (siteSvc *SiteService) generateModuleConfigIndex(data internals.ModuleConfig) {

	rootRoot := &TreeItemAsync{
		Id:          "root",
		DisplayName: "Root",
		Parent:      "#",
		Type:        "Folder",
		Children:    true,
	}
	siteSvc.configIndex.AddItem(rootRoot, nil)

	siteSvc.recursiveModuleConfig("root", data)

}

func (siteSvc *SiteService) Config(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "Configuration",
		"page":     "config.tmpl",
		"IsConfig": true,
		"data":     siteSvc.plan.PlannedValue.RootModule,
	})
}
func (siteSvc *SiteService) GetConfigAsync(c *gin.Context) {
	nodeId := c.Query("id")
	resp := siteSvc.configIndex.GetItemsByParent(nodeId)
	c.IndentedJSON(http.StatusOK, resp)
}

func (siteSvc *SiteService) GetItem(c *gin.Context) {
	addr := c.Param("address")

	item := siteSvc.configIndex.GetResource(addr)
	if item == nil {
		item = map[string]string{
			"message": "no data",
		}
	}

	c.IndentedJSON(http.StatusOK, item)
}

func (siteSvc *SiteService) recursiveModuleConfig(parentId string, parent internals.ModuleConfig) {

	// Outputs
	rootOutputs := &TreeItemAsync{
		Id:          hash(fmt.Sprintf("%soutputs", parentId)),
		DisplayName: "Outputs",
		Parent:      parentId,
		Type:        "Folder",
		Children:    len(parent.Outputs) > 0,
	}
	siteSvc.configIndex.AddItem(rootOutputs, nil)

	for outputId, output := range parent.Outputs {
		item := &TreeItemAsync{
			Id:          hash(fmt.Sprintf("%s%s", parentId, outputId)),
			DisplayName: outputId,
			Parent:      rootOutputs.Id,
			Type:        "Output",
			Children:    false,
		}
		siteSvc.configIndex.AddItem(item, output)
	}

	// Resource Configuration
	rootResources := &TreeItemAsync{
		Id:          hash(fmt.Sprintf("%sresources", parentId)),
		DisplayName: "Resources",
		Parent:      parentId,
		Type:        "Folder",
		Children:    len(parent.Resources) > 0,
	}
	siteSvc.configIndex.AddItem(rootResources, nil)

	for _, res := range parent.Resources {
		item := &TreeItemAsync{
			Id:          hash(fmt.Sprintf("%s%s", parentId, res.Address)),
			DisplayName: res.Name,
			Parent:      rootResources.Id,
			Type:        "Resource",
			Children:    false,
		}
		siteSvc.configIndex.AddItem(item, res)
	}

	// SubModules
	rootModules := &TreeItemAsync{
		Id:          hash(fmt.Sprintf("%smodules", parentId)),
		DisplayName: "Modules",
		Parent:      parentId,
		Type:        "Folder",
		Children:    len(parent.Modules) > 0,
	}
	siteSvc.configIndex.AddItem(rootModules, nil)

	for moduleId, mod := range parent.Modules {
		item := &TreeItemAsync{
			Id:          hash(fmt.Sprintf("%s%s", parentId, moduleId)),
			DisplayName: moduleId,
			Parent:      rootModules.Id,
			Type:        "Module",
			Children:    true,
		}
		siteSvc.configIndex.AddItem(item, mod)

		siteSvc.recursiveModuleConfig(item.Id, mod.Module)

	}
}
