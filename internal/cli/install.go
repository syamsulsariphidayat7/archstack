package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func listCmd() error {
	return printStatusTable(registry.AllTools())
}

func installCmd(names []string) error {
	hasError := false
	for _, name := range names {
		tool, ok := registry.GetTool(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "[error] Tool tidak dikenal: %s\n", name)
			hasError = true
			continue
		}
		if system.IsInstalled(tool.Binary, tool.Pkg) {
			fmt.Printf("[ok] %s sudah terinstal\n", name)
			continue
		}
		if len(tool.Prompts) > 0 {
			cached := registry.GetCachedAnswers(tool.Name)
			answers, err := collectAnswers(tool, cached)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] %s: %v\n", name, err)
				hasError = true
				continue
			}
			registry.SetCachedAnswers(tool.Name, answers)
			printAnswers(tool, answers)
		}
		if err := installOne(tool); err != nil {
			fmt.Fprintf(os.Stderr, "[error] %s: %v\n", name, err)
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("beberapa tool gagal diinstal")
	}
	return nil
}

func installOne(tool *registry.Tool) error {
	fmt.Printf("[install] Menginstal %s...\n", tool.Name)
	switch tool.From {
	case registry.SourcePacman:
		return system.InstallPacman(tool.Pkg)
	case registry.SourceYay:
		return system.InstallYay(tool.Pkg)
	}
	return nil
}

func collectAnswers(tool *registry.Tool, cached map[string]string) (map[string]string, error) {
	answers := make(map[string]string)
	fmt.Printf("\n=== Menginstal %s ===\n", tool.Name)
	for _, p := range tool.Prompts {
		fmt.Printf("%s ", p.Question)
		defaultVal := p.Recommended
		if cv, ok := cached[p.Key]; ok {
			defaultVal = cv
		}
		switch p.Type {
		case registry.PromptSelect:
			fmt.Printf("[%s]: ", strings.Join(p.Options, "/"))
			var input string
			fmt.Scanln(&input)
			if input == "" {
				input = defaultVal
			}
			answers[p.Key] = input
		case registry.PromptText:
			fmt.Printf("(terakhir: %s): ", defaultVal)
			var input string
			fmt.Scanln(&input)
			if input == "" {
				input = defaultVal
			}
			answers[p.Key] = input
		}
	}
	fmt.Println()
	return answers, nil
}

func printAnswers(tool *registry.Tool, answers map[string]string) {
	for _, p := range tool.Prompts {
		fmt.Printf("  %s: %s\n", p.Question, answers[p.Key])
	}
}
