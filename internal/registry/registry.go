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
	{
		Name: "python", Pkg: "python", From: SourcePacman, Binary: "python3",
		Desc: "Python programming language",
	},
	{
		Name: "python-pip", Pkg: "python-pip", From: SourcePacman, Binary: "pip",
		Desc: "Python package installer",
	},
	{
		Name: "go", Pkg: "go", From: SourcePacman, Binary: "go",
		Desc: "Go programming language",
	},
	{
		Name: "rust", Pkg: "rust", From: SourcePacman, Binary: "rustc",
		Desc: "Rust programming language",
	},
	{
		Name: "deno", Pkg: "deno", From: SourcePacman, Binary: "deno",
		Desc: "JavaScript and TypeScript runtime",
	},
	{
		Name: "bun", Pkg: "bun", From: SourcePacman, Binary: "bun",
		Desc: "Fast JavaScript runtime and toolkit",
	},
	{
		Name: "java", Pkg: "jdk-openjdk", From: SourcePacman, Binary: "java",
		Desc: "OpenJDK Java development kit",
	},
	{
		Name: "sqlite", Pkg: "sqlite", From: SourcePacman, Binary: "sqlite3",
		Desc: "SQLite embedded database",
	},
	{
		Name: "certbot", Pkg: "certbot", From: SourcePacman, Binary: "certbot",
		Desc: "Let's Encrypt SSL certificate tool",
	},
	{
		Name: "memcached", Pkg: "memcached", From: SourcePacman, Binary: "memcached",
		Service: "memcached", Desc: "Distributed memory object caching",
	},
	{
		Name: "mongodb", Pkg: "mongodb-bin", From: SourceYay, Binary: "mongod",
		Service: "mongodb", Desc: "MongoDB document database (AUR)",
	},
	{
		Name: "rabbitmq", Pkg: "rabbitmq", From: SourcePacman, Binary: "rabbitmq-server",
		Service: "rabbitmq", Desc: "AMQP message broker",
	},
	{
		Name: "fail2ban", Pkg: "fail2ban", From: SourcePacman, Binary: "fail2ban-client",
		Service: "fail2ban", Desc: "Brute force protection daemon",
	},
	{
		Name: "ffmpeg", Pkg: "ffmpeg", From: SourcePacman, Binary: "ffmpeg",
		Desc: "Audio and video processing toolkit",
	},
	{
		Name: "imagemagick", Pkg: "imagemagick", From: SourcePacman, Binary: "convert",
		Desc: "Image manipulation tool",
	},
	{
		Name: "tmux", Pkg: "tmux", From: SourcePacman, Binary: "tmux",
		Desc: "Terminal multiplexer",
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
