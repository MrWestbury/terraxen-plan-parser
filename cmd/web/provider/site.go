package provider

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"terraform-plan-parser/internals"

	"github.com/gin-gonic/gin"
)

type SiteService struct {
	plan            *internals.Plan
	stateIndex      Tree
	changeIndex     ChangeIndex
	changeAttrIndex ChangeAttrIndex
	configIndex     *TreeItemIndexer
}

func NewSiteService(plan *internals.Plan) *SiteService {
	siteSvc := &SiteService{
		plan:            plan,
		stateIndex:      nil,
		changeIndex:     nil,
		changeAttrIndex: nil,
		configIndex:     NewTreeItemIndexer(),
	}

	siteSvc.stateIndex = generateStateIndex(plan.PriorState.Values.RootModule)
	siteSvc.changeAttrIndex = generateChangeAttrIndex(plan.RelevantAttributes)
	siteSvc.changeIndex = generateChangeIndex(plan.ResourceChanges)
	siteSvc.generateProviderIndex(plan.Configuration.ProviderConfig)
	siteSvc.generateModuleConfigIndex(plan.Configuration.Root)
	return siteSvc
}

func (siteSvc *SiteService) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "Home",
		"page":      "home.tmpl",
		"version":   siteSvc.plan.FormatVersion,
		"terraform": siteSvc.plan.TerraformVersion,
	})
}

func MergeMaps(m1 Tree, m2 Tree) Tree {
	for ia, va := range m1 {
		m2[ia] = va
	}
	return m2
}

func hash(s string) string {
	hash := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}
