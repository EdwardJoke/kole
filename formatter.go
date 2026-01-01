package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	formatStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(2).
			Width(100)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#28A745")).
			Padding(1)

	keywordStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6"))

	stringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	commentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))
)

func highlightSyntax(line string) string {
	keywords := []string{"export", "alias", "function"}
	for _, kw := range keywords {
		if strings.HasPrefix(line, kw) {
			line = keywordStyle.Render(kw) + line[len(kw):]
			break
		}
	}

	if strings.HasPrefix(line, "#") {
		return commentStyle.Render(line)
	}

	quotedPattern := regexp.MustCompile(`"[^"]*"|'[^']*'`)
	line = quotedPattern.ReplaceAllStringFunc(line, func(s string) string {
		return stringStyle.Render(s)
	})

	pathPattern := regexp.MustCompile(`\$[A-Z_][A-Z0-9_]*|/[^"'\s]+`)
	line = pathPattern.ReplaceAllStringFunc(line, func(s string) string {
		return pathStyle.Render(s)
	})

	return line
}

type ShellConfig struct {
	FilePath  string
	FileType  string
	Lines     []string
	Exports   []string
	Aliases   []string
	Functions []string
	Comments  []string
	Others    []string
}

func detectFileType(path string) string {
	if strings.HasSuffix(path, ".bashrc") {
		return "bash"
	}
	if strings.HasSuffix(path, ".zshrc") {
		return "zsh"
	}
	return "unknown"
}

func parseShellConfig(path string) (*ShellConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ShellConfig{
		FilePath: path,
		FileType: detectFileType(path),
	}

	exportPattern := regexp.MustCompile(`^export\s+(\w+)=`)
	aliasPattern := regexp.MustCompile(`^alias\s+(\w+)=`)
	functionPattern := regexp.MustCompile(`^(\w+)\s*\(\)\s*\{`)
	commentPattern := regexp.MustCompile(`^#`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case exportPattern.MatchString(line):
			config.Exports = append(config.Exports, line)
		case aliasPattern.MatchString(line):
			config.Aliases = append(config.Aliases, line)
		case functionPattern.MatchString(line):
			config.Functions = append(config.Functions, line)
		case commentPattern.MatchString(line):
			config.Comments = append(config.Comments, line)
		default:
			config.Others = append(config.Others, line)
		}
	}

	return config, scanner.Err()
}

func formatShellConfig(config *ShellConfig) string {
	var lines []string

	lines = append(lines, highlightSyntax("# Comments"))
	for _, comment := range config.Comments {
		lines = append(lines, highlightSyntax(comment))
	}

	if len(config.Comments) > 0 && (len(config.Exports) > 0 || len(config.Aliases) > 0 || len(config.Functions) > 0) {
		lines = append(lines, "")
	}

	lines = append(lines, highlightSyntax("# Exports"))
	sort.Strings(config.Exports)
	for _, exp := range config.Exports {
		lines = append(lines, highlightSyntax(exp))
	}

	if len(config.Exports) > 0 && (len(config.Aliases) > 0 || len(config.Functions) > 0) {
		lines = append(lines, "")
	}

	lines = append(lines, highlightSyntax("# Aliases"))
	sort.Strings(config.Aliases)
	for _, alias := range config.Aliases {
		lines = append(lines, highlightSyntax(alias))
	}

	if len(config.Aliases) > 0 && len(config.Functions) > 0 {
		lines = append(lines, "")
	}

	lines = append(lines, highlightSyntax("# Functions"))
	sort.Strings(config.Functions)
	for _, fn := range config.Functions {
		lines = append(lines, highlightSyntax(fn))
	}

	if len(config.Others) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, highlightSyntax("# Other"))
		for _, other := range config.Others {
			lines = append(lines, highlightSyntax(other))
		}
	}

	return strings.Join(lines, "\n")
}

func writeFormattedConfig(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func showFormatterMenu() string {
	var selectedAction string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Shell Config Formatter").
				Description("Format and organize your .bashrc or .zshrc files").
				Options(
					huh.NewOption("Format .bashrc", "format_bashrc"),
					huh.NewOption("Format .zshrc", "format_zshrc"),
					huh.NewOption("Format custom path", "format_custom"),
					huh.NewOption("Back to main menu", "back"),
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

func formatFile(path string) error {
	config, err := parseShellConfig(path)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	formatted := formatShellConfig(config)

	var confirm bool
	huh.NewConfirm().
		Title(fmt.Sprintf("Preview formatted %s?", config.FileType)).
		Affirmative("Show preview").
		Negative("Skip preview").
		Value(&confirm).
		Run()

	if confirm {
		preview := formatted
		if len(preview) > 1500 {
			preview = preview[:1500] + "\n..."
		}
		fmt.Println(formatStyle.Render(fmt.Sprintf("Preview:\n\n%s", preview)))
		huh.NewConfirm().
			Title("Apply formatting?").
			Affirmative("Yes, apply").
			Negative("No, cancel").
			Value(&confirm).
			Run()

		if confirm {
			err = writeFormattedConfig(path, formatted)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Println(infoStyle.Render(fmt.Sprintf("Successfully formatted %s", path)))
		}
	}

	return nil
}

func runFormatter() {
	for {
		action := showFormatterMenu()

		switch action {
		case "format_bashrc":
			home, _ := os.UserHomeDir()
			path := fmt.Sprintf("%s/.bashrc", home)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println(infoStyle.Render(".bashrc not found in home directory"))
				huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()
				continue
			}
			if err := formatFile(path); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()

		case "format_zshrc":
			home, _ := os.UserHomeDir()
			path := fmt.Sprintf("%s/.zshrc", home)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println(infoStyle.Render(".zshrc not found in home directory"))
				huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()
				continue
			}
			if err := formatFile(path); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()

		case "format_custom":
			var path string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter path to shell config file").
						Placeholder("/path/to/.bashrc or .zshrc").
						Value(&path),
				),
			).WithTheme(huh.ThemeCharm())
			form.Run()

			if path != "" {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Println(infoStyle.Render("File not found"))
				} else {
					if err := formatFile(path); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}
				huh.NewConfirm().Title("Press Enter to continue").Affirmative("OK").Negative("").Run()
			}

		case "back":
			return
		}
	}
}
