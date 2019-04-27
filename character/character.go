package character

import (
	"io/ioutil"
	"math"
	"sort"
	"strings"

	"github.com/zwzn/dnd/blade"

	"gopkg.in/yaml.v2"

	"golang.org/x/xerrors"
)

var (
	ErrInvalidFormat = xerrors.New("invalid markdown format")
)

var Skills = map[string]string{
	"Acrobatics":      "dex",
	"Animal Handling": "wis",
	"Arcana":          "int",
	"Athletics":       "str",
	"Deception":       "cha",
	"History":         "int",
	"Insight":         "wis",
	"Intimidation":    "cha",
	"Investigation":   "int",
	"Medicine":        "wis",
	"Nature":          "int",
	"Perception":      "wis",
	"Performance":     "cha",
	"Persuasion":      "cha",
	"Religion":        "int",
	"Sleight of Hand": "dex",
	"Stealth":         "dex",
	"Survival":        "wis",
}

type Character struct {
	character

	CurrentHP int
	Status    []string
}
type character struct {
	file          string
	rawMD         string
	Name          string         `yaml:"name"`
	Level         int            `yaml:"level"`
	MaxHP         int            `yaml:"max_hp"`
	Speed         int            `yaml:"speed"`
	Initiative    int            `yaml:"initiative"`
	AbilityScores map[string]int `yaml:"ability_scores"`
	Proficiencies []string       `yaml:"proficiencies"`
	Expertise     []string       `yaml:"expertise"`
}

type Skill struct {
	Prof  string
	Mod   string
	Skill string
	Bonus int
}

func NewFile(file string) (*Character, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("error opening character file: %w", err)
	}
	return New(file, string(b))
}
func New(file, md string) (*Character, error) {

	ch := &Character{}
	ch.file = file
	parts := strings.SplitN(md, "\n----", 2)

	if !strings.HasPrefix(parts[0], "----") {
		return nil, ErrInvalidFormat
	}
	if len(parts) > 1 {
		ch.rawMD = parts[1]
	}
	err := yaml.Unmarshal([]byte(parts[0][4:]), &ch.character)
	if err != nil {
		return nil, err
	}
	ch.CurrentHP = ch.MaxHP
	ch.Status = []string{}

	return ch, nil
}
func (c *Character) Blade(b *blade.Parser) {
	c.rawMD = b.Parse(c.rawMD)
}
func (c *Character) Proficiency() int {
	if c.Level < 5 {
		return 2
	}
	if c.Level < 9 {
		return 3
	}
	if c.Level < 13 {
		return 4
	}
	if c.Level < 17 {
		return 5
	}
	return 4
}

func (c *Character) AbilityScoreMods() map[string]int {
	saves := map[string]int{}
	for ability, score := range c.AbilityScores {
		saves[ability] = mod(score)
	}
	return saves
}

func (c *Character) SavingThrows() map[string]int {
	saves := map[string]int{}
	for ability, score := range c.AbilityScoreMods() {
		bonus := 0
		if inList(c.Proficiencies, ability) {
			bonus = c.Proficiency()
		}
		if inList(c.Expertise, ability) {
			bonus = c.Proficiency() * 2
		}
		saves[ability] = score + bonus
	}
	return saves
}

func (c *Character) Skills() []*Skill {
	skills := []*Skill{}
	for skill, mod := range Skills {
		prof := ""
		bonus := c.AbilityScoreMods()[mod]
		if inList(c.Expertise, skill) {
			prof = "E"
			bonus += c.Proficiency() * 2
		} else if inList(c.Proficiencies, skill) {
			prof = "P"
			bonus += c.Proficiency()
		}
		skills = append(skills, &Skill{
			Prof:  prof,
			Mod:   mod,
			Skill: skill,
			Bonus: bonus,
		})
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Skill < skills[j].Skill
	})
	return skills
}

func mod(score int) int {
	return int(math.Floor((float64(score) - 10) / 2))
}

func inList(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
