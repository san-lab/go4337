package ui

import (
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddDictEntryItem = &Item{Label: "Add an Entry", Details: "Add a new dictionary entry"}
var RemoveDictEntryItem = &Item{Label: "Remove an Entry", Details: "Remove a dictionary entry"}
var FromAnotherDictItem = &Item{Label: "From Another Dictionary", Details: "Select a dictionary to choose from"}

func StringFromDictionaryUI(dictionary string) (dictname, name, value string, success bool) {
	selectToRemove := false

	for {
		items := []*Item{}
		dict := state.GetDictionary(dictionary)
		for k, v := range dict {
			items = append(items, &Item{Label: k, Value: v})
		}
		if !selectToRemove {

			items = append(items, AddDictEntryItem, RemoveDictEntryItem, FromAnotherDictItem)
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
