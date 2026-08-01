package registry

type SourceType string

const (
	SourcePacman SourceType = "pacman"
	SourceYay    SourceType = "yay"
)

type PromptType string

const (
	PromptSelect PromptType = "select"
	PromptText   PromptType = "text"
)

type Prompt struct {
	Key         string
	Question    string
	Type        PromptType
	Options     []string
	Recommended string
}

type Tool struct {
	Name    string
	Pkg     string
	From    SourceType
	Binary  string
	Service string
	Desc    string
	Prompts []Prompt
}

var defaultTools = []Tool{
	{
		Name: "php", Pkg: "php", From: SourcePacman, Binary: "php",
		Desc: "PHP scripting language",
	},
	{
		Name: "php-fpm", Pkg: "php-fpm", From: SourcePacman,
		Service: "php-fpm", Desc: "PHP FastCGI Process Manager",
	},
	{
		Name: "node", Pkg: "nodejs", From: SourcePacman, Binary: "node",
		Desc: "Node.js JavaScript runtime",
	},
	{
		Name: "npm", Pkg: "npm", From: SourcePacman, Binary: "npm",
		Desc: "Node.js package manager",
	},
	{
		Name: "pnpm", Pkg: "pnpm", From: SourceYay, Binary: "pnpm",
		Desc: "Fast, disk space efficient package manager",
	},
	{
		Name: "postgres", Pkg: "postgresql", From: SourcePacman, Binary: "psql",
		Service: "postgresql", Desc: "PostgreSQL database",
	},
	{
		Name: "redis", Pkg: "redis", From: SourcePacman, Binary: "redis-cli",
		Service: "redis", Desc: "Redis in-memory store",
	},
	{
		Name: "mariadb", Pkg: "mariadb", From: SourcePacman, Binary: "mariadb",
		Service: "mariadb", Desc: "MariaDB database server",
	},
	{
		Name: "docker", Pkg: "docker", From: SourcePacman, Binary: "docker",
		Service: "docker", Desc: "Docker container platform",
	},
	{
		Name: "git", Pkg: "git", From: SourcePacman, Binary: "git",
		Desc: "Distributed version control system",
	},
	{
		Name: "composer", Pkg: "composer", From: SourcePacman, Binary: "composer",
		Desc: "PHP dependency manager",
	},
	{
		Name: "nginx", Pkg: "nginx", From: SourcePacman, Binary: "nginx",
		Service: "nginx", Desc: "Nginx web server",
		Prompts: []Prompt{
			{
				Key:      "version",
				Question: "Pilih versi nginx",
				Type:     PromptSelect,
				Options:  []string{"stable (repo resmi pacman)", "mainline (AUR nginx-mainline)"},
				Recommended: "stable (repo resmi pacman)",
			},
			{
				Key:         "project_root",
				Question:    "Root project (kosongkan kalau gak perlu)",
				Type:        PromptText,
				Recommended: "/srv/http/<nama-project>",
			},
		},
	},
}

func AllTools() []Tool {
	return defaultTools
}

func GetTool(name string) (*Tool, bool) {
	for _, t := range defaultTools {
		if t.Name == name {
			return &t, true
		}
	}
	return nil, false
}
