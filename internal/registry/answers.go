package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var answersCache_ map[string]map[string]string

func init() {
	answersCache_ = loadAnswers()
}

func GetCachedAnswers(toolName string) map[string]string {
	if answersCache_ == nil {
		answersCache_ = loadAnswers()
	}
	ans, ok := answersCache_[toolName]
	if !ok {
		return make(map[string]string)
	}
	return ans
}

func SetCachedAnswers(toolName string, answers map[string]string) {
	if answersCache_ == nil {
		answersCache_ = loadAnswers()
	}
	answersCache_[toolName] = answers
	saveAnswers(answersCache_)
}

func loadAnswers() map[string]map[string]string {
	data, err := os.ReadFile(answersPath())
	if err != nil {
		return make(map[string]map[string]string)
	}
	var cache map[string]map[string]string
	if err := json.Unmarshal(data, &cache); err != nil {
		return make(map[string]map[string]string)
	}
	return cache
}

func saveAnswers(cache map[string]map[string]string) {
	path := answersPath()
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(cache, "", "  ")
	os.WriteFile(path, data, 0644)
}

func answersPath() string {
	cfgDir := os.Getenv("XDG_CONFIG_HOME")
	if cfgDir == "" {
		home := os.Getenv("HOME")
		cfgDir = filepath.Join(home, ".config")
	}
	return filepath.Join(cfgDir, "archstack", "answers.json")
}
