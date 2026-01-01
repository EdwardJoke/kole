package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(22)

	secondaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#4A90E2")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(22)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#28A745")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(22)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			PaddingBottom(1)
)

type EnvVar struct {
	Name  string
	Value string
}

func getEnvVars() []EnvVar {
	var envVars []EnvVar
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars = append(envVars, EnvVar{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}
	return envVars
}

func getPathEntries() []string {
	path := os.Getenv("PATH")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, ":")
}

func showMainMenu() string {
	_ = getEnvVars() // Load environment variables

	var selectedAction string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Environment Variable Manager").
				Options(
					huh.NewOption("View All Environment Variables", "view_all"),
					huh.NewOption("Manage PATH Variable", "manage_path"),
					huh.NewOption("Add New Environment Variable", "add_var"),
					huh.NewOption("Edit Environment Variable", "edit_var"),
					huh.NewOption("Delete Environment Variable", "delete_var"),
					huh.NewOption("Search Environment Variables", "search_var"),
					huh.NewOption("Format Shell Config (.bashrc/.zshrc)", "format_shell"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&selectedAction),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		os.Exit(1)
	}

	return selectedAction
}

func viewAllEnvVars() {
	envVars := getEnvVars()

	var content strings.Builder
	content.WriteString("All Environment Variables\n")
	content.WriteString(strings.Repeat("=", 50) + "\n\n")

	for _, env := range envVars {
		displayValue := env.Value
		if len(displayValue) > 80 {
			displayValue = displayValue[:80] + "..."
		}
		displayValue = strings.ReplaceAll(displayValue, "\n", "\\n")
		content.WriteString(fmt.Sprintf("%-30s -> %s\n", env.Name, displayValue))
	}

	fmt.Println(borderStyle.Render(content.String()))

	var confirm bool
	huh.NewConfirm().
		Title("Press Enter to continue").
		Affirmative("OK").
		Negative("").
		Value(&confirm).
		Run()
}

func managePath() {
	pathEntries := getPathEntries()

	var options []huh.Option[string]
	for i, entry := range pathEntries {
		display := fmt.Sprintf("%d. %s", i+1, entry)
		if len(display) > 70 {
			display = display[:70] + "..."
		}
		options = append(options, huh.NewOption(display, entry))
	}

	var selectedPath string
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("PATH Entries Management").
				Description(fmt.Sprintf("Total entries: %d", len(pathEntries))).
				Options(options...).
				Value(&selectedPath),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do with this PATH entry?").
				Options(
					huh.NewOption("View full path", "view"),
					huh.NewOption("Delete this entry", "delete"),
					huh.NewOption("Move up", "move_up"),
					huh.NewOption("Move down", "move_down"),
					huh.NewOption("Back to main menu", "back"),
				).
				Value(&action),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	switch action {
	case "view":
		fmt.Println(borderStyle.Render(fmt.Sprintf("Full PATH Entry:\n\n%s", selectedPath)))
		huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()
	case "delete":
		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Are you sure you want to delete:\n%s", selectedPath)).
			Affirmative("Yes, delete it").
			Negative("No, keep it").
			Value(&confirm).
			Run()
		if confirm {
			fmt.Println(successStyle.Render("PATH entry deleted successfully!"))
			fmt.Println("Note: To permanently remove this entry, add it to your shell profile with the entry removed from PATH")
		}
	case "back":
		return
	}
}

func addNewEnvVar() {
	var name, value string
	var varConfirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the environment variable name").
				Placeholder("MY_VARIABLE").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Name cannot be empty")
					}
					if strings.Contains(s, "=") {
						return fmt.Errorf("Name cannot contain '='")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Enter the value").
				Placeholder("Enter the value here...").
				Value(&value),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Set %s=%s?", name, value)).
				Affirmative("Yes, set it").
				Negative("No, cancel").
				Value(&varConfirmed),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	if varConfirmed {
		os.Setenv(name, value)
		fmt.Println(successStyle.Render(fmt.Sprintf("✅ Successfully set %s=%s", name, value)))
		fmt.Println("\n💡 Note: This change is temporary for this session.")
		fmt.Println("To make it permanent, add this line to your shell profile:")
		fmt.Printf("   export %s=\"%s\"\n", name, value)
	}
}

func editEnvVar() {
	envVars := getEnvVars()

	var options []huh.Option[string]
	for _, env := range envVars {
		displayValue := env.Value
		if len(displayValue) > 40 {
			displayValue = displayValue[:40] + "..."
		}
		displayValue = strings.ReplaceAll(displayValue, "\n", "\\n")
		options = append(options, huh.NewOption(fmt.Sprintf("%s = %s", env.Name, displayValue), env.Name))
	}

	var selectedVar string
	var newValue string
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the environment variable to edit").
				Options(options...).
				Value(&selectedVar),
		),
		huh.NewGroup(
			huh.NewText().
				Title(fmt.Sprintf("Edit value for %s", selectedVar)).
				Value(&newValue),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Update %s?", selectedVar)).
				Affirmative("Yes, update").
				Negative("No, cancel").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	if confirmed {
		os.Setenv(selectedVar, newValue)
		fmt.Println(successStyle.Render(fmt.Sprintf("Updated %s successfully", selectedVar)))
	}
}

func deleteEnvVar() {
	envVars := getEnvVars()

	var options []huh.Option[string]
	for _, env := range envVars {
		options = append(options, huh.NewOption(env.Name, env.Name))
	}

	var selectedVar string
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the environment variable to delete").
				Options(options...).
				Value(&selectedVar),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete %s?", selectedVar)).
				Affirmative("Yes, delete").
				Negative("No, keep").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	if confirmed {
		os.Unsetenv(selectedVar)
		fmt.Println(successStyle.Render(fmt.Sprintf("Deleted %s successfully", selectedVar)))
		fmt.Println("\nNote: This change is temporary for this session.")
	}
}

func searchEnvVars() {
	var searchTerm string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Search environment variables").
				Placeholder("Enter search term...").
				Value(&searchTerm),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	if searchTerm == "" {
		return
	}

	envVars := getEnvVars()
	var matches []EnvVar
	searchLower := strings.ToLower(searchTerm)

	for _, env := range envVars {
		if strings.Contains(strings.ToLower(env.Name), searchLower) ||
			strings.Contains(strings.ToLower(env.Value), searchLower) {
			matches = append(matches, env)
		}
	}

	var content strings.Builder
	if len(matches) == 0 {
		content.WriteString(fmt.Sprintf("No matches found for '%s'\n", searchTerm))
	} else {
		content.WriteString(fmt.Sprintf("Found %d matches for '%s'\n", len(matches), searchTerm))
		content.WriteString(strings.Repeat("=", 50) + "\n\n")

		for _, env := range matches {
			displayValue := env.Value
			if len(displayValue) > 80 {
				displayValue = displayValue[:80] + "..."
			}
			displayValue = strings.ReplaceAll(displayValue, "\n", "\\n")
			content.WriteString(fmt.Sprintf("%-30s -> %s\n", env.Name, displayValue))
		}
	}

	fmt.Println(borderStyle.Render(content.String()))

	huh.NewConfirm().
		Title("Press Enter to continue").
		Affirmative("OK").
		Negative("").
		Run()
}

func main() {
	for {
		action := showMainMenu()

		switch action {
		case "view_all":
			viewAllEnvVars()
		case "manage_path":
			managePath()
		case "add_var":
			addNewEnvVar()
		case "edit_var":
			editEnvVar()
		case "delete_var":
			deleteEnvVar()
		case "search_var":
			searchEnvVars()
		case "format_shell":
			runFormatter()
		case "exit":
			fmt.Println("\nGoodbye! Thanks for using Environment Variable Manager!")
			return
		}
	}
}
