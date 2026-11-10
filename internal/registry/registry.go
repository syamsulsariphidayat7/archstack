package registry

type SourceType string

const (
	SourcePacman SourceType = "pacman"
	SourceYay    SourceType = "yay"
	// SourceServiceOnly: archstack tidak memasang/menghapus binary, hanya
	// mengelola service-nya (binary diasumsikan sudah dipasang manual).
	SourceServiceOnly SourceType = "service"
)

// ServiceScope menentukan domain systemd yang dipakai sebuah service.
type ServiceScope int

const (
	ServiceSystem ServiceScope = iota // systemctl + sudo (default)
	ServiceUser                       // systemctl --user (tanpa sudo)
)

type Tool struct {
	Name            string
	Pkg             string
	From            SourceType
	Binary          string
	Service         string
	ServiceScope    ServiceScope
	Linger          bool     // aktifkan loginctl enable-linger (user service)
	PidFile         string   // opsional, "~/.app/server/.pid" untuk deteksi instance manual
	HealthURL       string   // opsional, dicek setelah service dinyalakan
	UnitArgs        []string // opsional, argumen ExecStart (mis. serve --no-open)
	UnitDescription string   // opsional, Description di unit; default: Desc
	Desc            string
}

// IsServiceOnly menandai tool yang binary-nya tidak dikelola archstack.
func (t Tool) IsServiceOnly() bool {
	return t.From == SourceServiceOnly
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
	{
		Name: "omniroute", From: SourceServiceOnly, Binary: "omniroute",
		Service: "omniroute", ServiceScope: ServiceUser, Linger: true,
		PidFile:         "~/.omniroute/server/.pid",
		HealthURL:       "http://127.0.0.1:20128/api/health",
		UnitArgs:        []string{"serve", "--no-open"},
		UnitDescription: "OmniRoute AI proxy router",
		Desc:            "OmniRoute AI proxy router — autostart service (binary manual)",
	},
	{
		Name: "archstack-web", From: SourceServiceOnly, Binary: "archstack",
		Service: "archstack-web", ServiceScope: ServiceUser, Linger: true,
		HealthURL:       "http://127.0.0.1:8090/api/tools",
		UnitArgs:        []string{"web"},
		UnitDescription: "ArchStack Web UI",
		Desc:            "ArchStack Web UI — dashboard lokal autostart (binary archstack sendiri)",
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
