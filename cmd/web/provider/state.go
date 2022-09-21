package provider

import (
	"net/http"
	"terraform-plan-parser/internals"

	"github.com/gin-gonic/gin"
)

// Serve the state page
func (siteSvc *SiteService) State(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Previous State",
		"page":    "state.tmpl",
		"IsState": true,
	})
}

// Api call to list all the state objects
func (siteSvc *SiteService) GetStateList(c *gin.Context) {
	modules := make([]TreeItem, 0)
	resources := make([]TreeItem, 0)

	for _, item := range siteSvc.stateIndex {
		respItem := TreeItem{
			Id:          item.Id,
			Parent:      item.Parent,
			DisplayName: item.Name,
			Type:        item.Type,
			LiAttr: map[string]string{
				"data-addr": item.Address,
			},
		}
		if respItem.Type == "Module" {
			modules = append(modules, respItem)
		} else {
			resources = append(resources, respItem)
		}

	}

	result := append(modules, resources...)
	c.IndentedJSON(http.StatusOK, result)
}

// Api call to get a single state resource
func (siteSvc *SiteService) GetStateResource(c *gin.Context) {
	addr := c.Param("address")

	res, exists := siteSvc.stateIndex[addr]
	if !exists {
		resp := map[string]string{
			"message": "state resource not found",
		}

		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	c.IndentedJSON(http.StatusOK, res.Resource)
}

func generateStateIndex(root internals.RootModule) Tree {
	result := Tree{}

	for _, child := range root.ChildModules {
		childRes := RecurseStateResources("#", &child)
		result = MergeMaps(result, childRes)
	}

	for _, res := range root.Resources {
		node := &TreeNode{
			Id:       hash(res.Address),
			Type:     "Resource",
			Address:  res.Address,
			Parent:   "#",
			Name:     res.Name,
			Resource: res,
		}
		result[node.Id] = node
	}

	return result
}

func RecurseStateResources(parentId string, parent *internals.ChildModule) Tree {
	result := Tree{}

	rootNode := &TreeNode{
		Id:      hash(parent.Address),
		Type:    "Module",
		Address: parent.Address,
		Parent:  parentId,
		Name:    parent.Address,
	}
	result[rootNode.Id] = rootNode

	for _, child := range parent.ChildModules {
		childRes := RecurseStateResources(rootNode.Id, &child)
		result = MergeMaps(result, childRes)
	}

	for _, res := range parent.Resources {
		node := &TreeNode{
			Id:       hash(res.Address),
			Address:  res.Address,
			Type:     "Resource",
			Parent:   rootNode.Id,
			Name:     res.Name,
			Resource: res,
		}
		result[node.Id] = node
	}

	return result
}
