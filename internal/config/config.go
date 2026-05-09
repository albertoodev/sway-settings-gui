package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	ModKey   string
	Terminal string
	Menu     string
	Opacity  float64

	GapsInner    int
	GapsOuter    int
	CornerRadius int
	BorderSize   int
	BlurEnabled  bool
	BlurPasses   int
	BlurRadius   int
	DimInactive  float64

	Outputs []Output
}

type Output struct {
	Name       string
	Resolution string
	Position   string
	Transform  string
}

func SwayConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sway")
}

func MainConfigPath() string  { return filepath.Join(SwayConfigDir(), "config") }
func AppearancePath() string  { return filepath.Join(SwayConfigDir(), "config.d", "appearance.conf") }
func OutputsPath() string     { return filepath.Join(SwayConfigDir(), "config.d", "outputs.conf") }

func Load() (*Config, error) {
	cfg := &Config{
		ModKey:       "Mod4",
		Terminal:     "kitty",
		Opacity:      1.0,
		GapsInner:    4,
		GapsOuter:    2,
		CornerRadius: 15,
		BorderSize:   1,
		BlurEnabled:  true,
		BlurPasses:   2,
		BlurRadius:   8,
		DimInactive:  0.1,
	}
	if err := cfg.loadMain(); err != nil {
		return nil, fmt.Errorf("main config: %w", err)
	}
	if err := cfg.loadAppearance(); err != nil {
		return nil, fmt.Errorf("appearance: %w", err)
	}
	if err := cfg.loadOutputs(); err != nil {
		return nil, fmt.Errorf("outputs: %w", err)
	}
	return cfg, nil
}

func (c *Config) loadMain() error {
	data, err := os.ReadFile(MainConfigPath())
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if m := regexp.MustCompile(`^set \$mod (\S+)`).FindStringSubmatch(line); m != nil {
			c.ModKey = m[1]
		}
		if m := regexp.MustCompile(`^set \$term (.+)`).FindStringSubmatch(line); m != nil {
			c.Terminal = strings.TrimSpace(m[1])
		}
		if m := regexp.MustCompile(`^set \$menu (.+)`).FindStringSubmatch(line); m != nil {
			c.Menu = strings.TrimSpace(m[1])
		}
		if m := regexp.MustCompile(`^set \$opacity ([0-9.]+)`).FindStringSubmatch(line); m != nil {
			c.Opacity, _ = strconv.ParseFloat(m[1], 64)
		}
	}
	return nil
}

func (c *Config) loadAppearance() error {
	data, err := os.ReadFile(AppearancePath())
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if m := regexp.MustCompile(`^gaps inner (\d+)`).FindStringSubmatch(line); m != nil {
			c.GapsInner, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^gaps outer (\d+)`).FindStringSubmatch(line); m != nil {
			c.GapsOuter, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^corner_radius (\d+)`).FindStringSubmatch(line); m != nil {
			c.CornerRadius, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^default_border pixel (\d+)`).FindStringSubmatch(line); m != nil {
			c.BorderSize, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^blur (enable|disable)`).FindStringSubmatch(line); m != nil {
			c.BlurEnabled = m[1] == "enable"
		}
		if m := regexp.MustCompile(`^blur_passes (\d+)`).FindStringSubmatch(line); m != nil {
			c.BlurPasses, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^blur_radius (\d+)`).FindStringSubmatch(line); m != nil {
			c.BlurRadius, _ = strconv.Atoi(m[1])
		}
		if m := regexp.MustCompile(`^default_dim_inactive ([0-9.]+)`).FindStringSubmatch(line); m != nil {
			c.DimInactive, _ = strconv.ParseFloat(m[1], 64)
		}
	}
	return nil
}

func (c *Config) loadOutputs() error {
	data, err := os.ReadFile(OutputsPath())
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`^output (\S+) resolution (\S+) position (\S+)(?:\s+transform\s+(\S+))?`)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if m := re.FindStringSubmatch(line); m != nil {
			c.Outputs = append(c.Outputs, Output{
				Name:       m[1],
				Resolution: m[2],
				Position:   m[3],
				Transform:  m[4],
			})
		}
	}
	return nil
}

type replacement struct {
	re  *regexp.Regexp
	val string
}

func replaceInFile(path string, replacements []replacement) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, r := range replacements {
			if r.re.MatchString(trimmed) {
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + r.val
				break
			}
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func (c *Config) Save() error {
	if err := c.saveMain(); err != nil {
		return fmt.Errorf("main config: %w", err)
	}
	if err := c.saveAppearance(); err != nil {
		return fmt.Errorf("appearance: %w", err)
	}
	if err := c.saveOutputs(); err != nil {
		return fmt.Errorf("outputs: %w", err)
	}
	return nil
}

func (c *Config) saveMain() error {
	return replaceInFile(MainConfigPath(), []replacement{
		{regexp.MustCompile(`^set \$mod `), fmt.Sprintf("set $mod %s", c.ModKey)},
		{regexp.MustCompile(`^set \$term `), fmt.Sprintf("set $term %s", c.Terminal)},
		{regexp.MustCompile(`^set \$menu `), fmt.Sprintf("set $menu %s", c.Menu)},
		{regexp.MustCompile(`^set \$opacity `), fmt.Sprintf("set $opacity %g", c.Opacity)},
	})
}

func (c *Config) saveAppearance() error {
	blurVal := "enable"
	if !c.BlurEnabled {
		blurVal = "disable"
	}
	return replaceInFile(AppearancePath(), []replacement{
		{regexp.MustCompile(`^gaps inner \d+`), fmt.Sprintf("gaps inner %d", c.GapsInner)},
		{regexp.MustCompile(`^gaps outer \d+`), fmt.Sprintf("gaps outer %d", c.GapsOuter)},
		{regexp.MustCompile(`^corner_radius \d+`), fmt.Sprintf("corner_radius %d", c.CornerRadius)},
		{regexp.MustCompile(`^default_border pixel \d+`), fmt.Sprintf("default_border pixel %d", c.BorderSize)},
		{regexp.MustCompile(`^default_floating_border pixel \d+`), fmt.Sprintf("default_floating_border pixel %d", c.BorderSize)},
		{regexp.MustCompile(`^blur (enable|disable)$`), fmt.Sprintf("blur %s", blurVal)},
		{regexp.MustCompile(`^blur_passes \d+`), fmt.Sprintf("blur_passes %d", c.BlurPasses)},
		{regexp.MustCompile(`^blur_radius \d+`), fmt.Sprintf("blur_radius %d", c.BlurRadius)},
		{regexp.MustCompile(`^default_dim_inactive [0-9.]+`), fmt.Sprintf("default_dim_inactive %.2f", c.DimInactive)},
	})
}

func (c *Config) saveOutputs() error {
	data, err := os.ReadFile(OutputsPath())
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	nameRe := regexp.MustCompile(`^output (\S+) resolution`)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		m := nameRe.FindStringSubmatch(trimmed)
		if m == nil {
			continue
		}
		for _, out := range c.Outputs {
			if out.Name != m[1] {
				continue
			}
			newLine := fmt.Sprintf("output %s resolution %s position %s", out.Name, out.Resolution, out.Position)
			if out.Transform != "" {
				newLine += " transform " + out.Transform
			}
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = indent + newLine
		}
	}
	return os.WriteFile(OutputsPath(), []byte(strings.Join(lines, "\n")), 0644)
}
