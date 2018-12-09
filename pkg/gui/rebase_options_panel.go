package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

type option struct {
	value string
}

// GetDisplayStrings is a function.
func (r *option) GetDisplayStrings() []string {
	return []string{r.value}
}

func (gui *Gui) handleCreateRebaseOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	options := []*option{
		{value: "continue"},
		{value: "abort"},
	}

	if gui.State.WorkingTreeState == "rebasing" {
		options = append(options, &option{value: "skip"})
	}

	handleMenuPress := func(index int) error {
		// need to get status again in case something has changed since menu appeared
		status := gui.State.WorkingTreeState

		if status != "merging" && status != "rebasing" {
			return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NotMergingOrRebasing"))
		}

		commandType := strings.Replace(status, "ing", "e", 1)
		// we should end up with a command like 'git merge --continue'
		command := fmt.Sprintf("git %s --%s", commandType, options[index].value)

		gui.Log.Info("going in")
		// gui.OSCommand.RunCommand(command)
		gui.OSCommand.RunCommand(command)
		// go func() { // TODO: WTF
		// }()

		if err := gui.refreshSidePanels(gui.g); err != nil {
			return err
		}
		return nil
	}

	var title string
	if gui.State.WorkingTreeState == "merging" {
		title = gui.Tr.SLocalize("MergeOptionsTitle")
	} else {
		title = gui.Tr.SLocalize("RebaseOptionsTitle")
	}

	return gui.createMenu(title, options, handleMenuPress)
}
