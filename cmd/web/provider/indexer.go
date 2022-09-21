package provider

import "log"

type ItemIndex map[string]interface{}

type TreeItemAsyncIndex map[string]*TreeItemAsync

type TreeItemIndexer struct {
	items   ItemIndex
	index   TreeItemAsyncIndex
	parents ParentIndex
}

func NewTreeItemIndexer() *TreeItemIndexer {
	tii := &TreeItemIndexer{
		items:   ItemIndex{},
		index:   TreeItemAsyncIndex{},
		parents: ParentIndex{},
	}

	return tii
}

func (tii *TreeItemIndexer) AddItem(item *TreeItemAsync, resource interface{}) {
	_, exists := tii.items[item.Id]
	if exists {
		log.Print("Warning, item exists:")
		log.Printf("   Indexed: %s %v", tii.index[item.Id].Id, tii.index[item.Id])
		log.Printf("  New item: %s %v", item.Id, item)
	}

	tii.items[item.Id] = resource
	tii.index[item.Id] = item

	lst, ok := tii.parents[item.Parent]
	if !ok {
		lst = []string{}
		tii.parents[item.Parent] = lst
	}
	tii.parents[item.Parent] = append(lst, item.Id)
}

func (tii *TreeItemIndexer) GetItemsByParent(parentId string) []*TreeItemAsync {
	resp := []*TreeItemAsync{}
	items, ok := tii.parents[parentId]
	if !ok {
		return resp
	}

	for _, itemId := range items {
		item := tii.index[itemId]
		if parentId != item.Parent {
			log.Printf("parent ID mismatch")
		}
		resp = append(resp, item)
	}
	return resp
}

func (tii *TreeItemIndexer) GetResource(itemId string) interface{} {
	return tii.items[itemId]
}
