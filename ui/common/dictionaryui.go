package common

import (
	"fmt"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddDictEntryItem = &Item{Label: "Add an Entry", Details: "Add a new dictionary entry"}
var RemoveDictEntryItem = &Item{Label: "Remove an Entry", Details: "Remove a dictionary entry"}
var RenameDictEntryItem = &Item{Label: "Rename an Entry", Details: "Rename a dictionary entry"}
var FromAnotherDictItem = &Item{Label: "From Another Dictionary", Details: "Select a dictionary to choose from"}

func StringFromDictionaryUI(dictionary string) (dictname, name, value string, success bool) {
	selectToRemove := false
	selectToRename := false

	for {
		items := []*Item{}
		dict := state.GetDictionary(dictionary)
		//sort the keys
		keys := make([]string, 0, len(dict))
		for k := range dict {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			items = append(items, &Item{Label: k, Value: dict[k]})
		}
		if !selectToRemove {

			items = append(items, AddDictEntryItem, RemoveDictEntryItem, RenameDictEntryItem, FromAnotherDictItem)
		}
		items = append(items, Back)
		spr := promptui.Select{Label: "Dictionary: " + dictionary, Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			return "", "", "", false
		}
		switch sel {
		case Back.Label:
			if selectToRemove {
				selectToRemove = false
				continue
			}
			return "", "", "", false
		case AddDictEntryItem.Label:
			return AddDictEntryUI(dictionary)
		case RemoveDictEntryItem.Label:
			selectToRemove = true
		case RenameDictEntryItem.Label:
			selectToRename = true
		case FromAnotherDictItem.Label:
			dict_name, entry_name, entry_value, ok := StringFromAllDictionariesUI(dictionary)
			if !ok {
				continue
			}
			new_entry_name := dict_name + entry_name
			dict[new_entry_name] = entry_value
			state.Save()
			return dictionary, new_entry_name, entry_value, true
		default:
			if selectToRemove {
				delete(dict, sel)
				state.Save()
				selectToRemove = false
				continue
			}
			if selectToRename {
				RenameItem := &Item{Label: fmt.Sprintf("New name for >>%s<<", sel)}
				err := InputNewStringUI(RenameItem)
				if err != nil {
					fmt.Println(err)
				} else {
					nname := RenameItem.Value.(string)
					dict[nname] = dict[sel]
					delete(dict, sel)
					state.Save()
				}
				selectToRename = false
				continue
			}
			return dictionary, sel, dict[sel], true
		}
	}

}

func StringFromAllDictionariesUI(dictionary string) (dictname, name, value string, success bool) {
	items := []*Item{}
	for _, dict := range state.GetDictionaries() {
		items = append(items, &Item{Label: dict})
	}
	items = append(items, Back)
	spr := promptui.Select{Label: "Select Dictionary", Items: items, Templates: ItemTemplate, Size: 10}
	_, sel, err := spr.Run()
	if err != nil {
		return "", "", "", false
	}
	switch sel {
	case Back.Label:
		return "", "", "", false
	default:
		return StringFromDictionaryUI(sel)
	}

}

func AddDictEntryUI(dictionary string) (dictname, name, value string, success bool) {
	dictname = dictionary
	it := &Item{Label: "Name of the entry"}
	err := InputNewStringUI(it)
	if err == nil {
		name = it.Value.(string)
	} else {
		return
	}
	it = &Item{Label: "Value of the entry"}
	err = InputNewStringUI(it)
	if err == nil {
		value = it.Value.(string)
	} else {
		return
	}
	state.GetDictionary(dictionary)[name] = value
	state.Save()
	success = true
	return
}
