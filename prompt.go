package gopromptui

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

type Prompt struct {
}

type item struct {
	ID         string
	IsSelected bool
}

func New() *Prompt {
	return &Prompt{}
}

func (p *Prompt) ValidateFloat(input string) error {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("please enter a valid number")
	}
	return nil
}

func (p *Prompt) ValidateInt(input string) error {
	_, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errors.New("please enter a valid number")
	}
	return nil
}

func (p *Prompt) ValidateString(input string) error {
	if len(input) < 1 {
		return errors.New("please enter a valid string")
	}
	return nil
}

// AskString is a wrapper for promptui.Prompt
// values:
// 1 - label
// 2 - default value
func (p *Prompt) AskString(values ...string) (string, error) {
	// if values 1 then label is first
	// if values 2 then label is first value, and default is second value
	var label, def string
	switch len(values) {
	case 1:
		label = values[0]
	case 2:
		label = values[0]
		def = values[1]
	default:
		return "", errors.New("invalid number of arguments")
	}

	return p.AskStringWithValidator(label, def, p.ValidateString)
}

func (p *Prompt) AskStringWithValidator(label, defaultValue string, validateFunc promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateFunc,
		Default:  defaultValue,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *Prompt) AskInt(label string, def int) (int, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: p.ValidateInt,
		Default:  strconv.Itoa(def),
	}

	result, err := prompt.Run()

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}

func (p *Prompt) AskBool(label string) (bool, error) {
	res, err := p.AskFromListString(label, []string{"yes", "no"})
	if err != nil {
		return false, err
	}

	return res == "yes", nil
}

func (p *Prompt) AskFromListString(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *Prompt) AskFromListInt(label string, items []int) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}

func (p *Prompt) AskFromListStringMultiple(label string, items []string) ([]string, error) {
	var allItems []*item
	for _, i := range items {
		allItems = append(allItems, &item{
			ID: i,
		})
	}

	selectedItems, err := p.selectItems("Done ✔", 0, allItems)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, i := range selectedItems {
		result = append(result, i.ID)
	}

	return result, nil
}

// https://liza.io/implementing-multiple-choice-selection-in-go-with-promptui/
func (p *Prompt) selectItems(doneID string, selectedPos int, allItems []*item) ([]*item, error) {
	// Always prepend a "Done" item to the slice if it doesn't
	// already exist.
	if len(allItems) > 0 && allItems[0].ID != doneID {
		var items = []*item{
			{
				ID: doneID,
			},
		}
		allItems = append(items, allItems...)
	}

	// Define promptui template
	templates := &promptui.SelectTemplates{
		Label: `{{if .IsSelected}}
                    ✔
                {{end}} {{ .ID }} - label`,
		Active:   "→ {{if .IsSelected}}✔ {{end}}{{ .ID | yellow }}",
		Inactive: "{{if .IsSelected}}✔ {{end}}{{ .ID }}",
	}

	prompt := promptui.Select{
		Label:     "Item",
		Items:     allItems,
		Templates: templates,
		// Start the cursor at the currently selected index
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenItem := allItems[selectionIdx]

	if chosenItem.ID != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return p.selectItems(doneID, selectionIdx, allItems)
	}

	// If the user selected the "Done" item, return
	// all selected items.
	var selectedItems []*item
	for _, i := range allItems {
		if i.IsSelected {
			selectedItems = append(selectedItems, i)
		}
	}
	return selectedItems, nil
}
