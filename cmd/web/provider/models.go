package provider

type TreeItem struct {
	Id          string            `json:"id"`
	Parent      string            `json:"parent"`
	DisplayName string            `json:"text"`
	LiAttr      map[string]string `json:"li_attr,omitempty"`
	Type        string            `json:"type"`
}

type TreeItemAsync struct {
	Id          string            `json:"id"`
	Parent      string            `json:"parent"`
	DisplayName string            `json:"text"`
	LiAttr      map[string]string `json:"li_attr,omitempty"`
	Type        string            `json:"type"`
	Children    bool              `json:"children"`
}

type TreeNode struct {
	Id       string      `json:"id"`
	Type     string      `json:"type"`
	Address  string      `json:"address"`
	Parent   string      `json:"parent"`
	Name     string      `json:"name"`
	Resource interface{} `json:"resource"`
}

type Tree map[string]*TreeNode

type ChangeItem struct {
	Id            string
	ModuleAddress string
	Name          string
	Address       string
	Change        string
	Resource      interface{}
}

type ChangeIndex map[string]ChangeItem

type ChangeAttrIndex map[string][]string

type GeneralIndex map[string]interface{}

type ParentIndex map[string][]string
