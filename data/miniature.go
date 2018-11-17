package data

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	// nameIndex is the index of the name field from a record in the CVS file
	nameIndex            = 0
	lineageIndex         = 1
	aspectIndex          = 2
	spawnCostIndex       = 3
	aspectCostIndex      = 4
	powerIndex           = 5
	defenseIndex         = 6
	lifeIndex            = 7
	abilitiesIndex       = 8
	flavorTextIndex      = 9
	collectorNumberIndex = 10
	setIndex             = 11
	rarityIndex          = 12
)

// Miniature is a miniature in Dreamblade. This struct is immutable.
type Miniature struct {
	id              string
	name            string
	lineage         string
	aspect          string
	spawnCost       string
	aspectCost      string
	power           string
	defense         string
	life            string
	abilities       string
	flavorText      string
	collectorNumber string
	set             string
	rarity          string
	nextMiniID      string
	prevMiniID      string
}

func newMiniature(r []string) Miniature {
	return Miniature{
		id:              createIDFromName(r[nameIndex]),
		name:            r[nameIndex],
		lineage:         r[lineageIndex],
		aspect:          r[aspectIndex],
		spawnCost:       r[spawnCostIndex],
		aspectCost:      r[aspectCostIndex],
		power:           r[powerIndex],
		defense:         r[defenseIndex],
		life:            r[lifeIndex],
		abilities:       r[abilitiesIndex],
		flavorText:      r[flavorTextIndex],
		collectorNumber: r[collectorNumberIndex],
		set:             r[setIndex],
		rarity:          r[rarityIndex],
	}
}

func (mini Miniature) ID() string {
	return mini.id
}
func (mini Miniature) Name() string {
	return mini.name
}
func (mini Miniature) Lineage() string {
	return mini.lineage
}
func (mini Miniature) Aspect() string {
	return mini.aspect
}
func (mini Miniature) SpawnCost() string {
	return mini.spawnCost
}
func (mini Miniature) AspectCost() string {
	return mini.aspectCost
}
func (mini Miniature) Power() string {
	return mini.power
}
func (mini Miniature) Defense() string {
	return mini.defense
}
func (mini Miniature) Life() string {
	return mini.life
}
func (mini Miniature) Abilities() string {
	return mini.abilities
}
func (mini Miniature) FlavorText() string {
	return mini.flavorText
}

// CollectorNumberAsInt returns the collector number of the miniature
// as an integer. If the value cannot be converted to an integer, a
// -1 is returned.
func (mini Miniature) CollectorNumberAsInt() int {
	n1, err := strconv.ParseInt(mini.collectorNumber, 10, 0)
	if err != nil {
		return -1
	}
	return int(n1)
}
func (mini Miniature) CollectorNumber() string {
	return mini.collectorNumber
}
func (mini Miniature) SetCode() string {
	return mini.set
}
func (mini Miniature) Set() (miniSet *MiniatureSet, exists bool) {
	var err error
	miniSet, err = GetMiniatureSetByID(mini.set)
	exists = err == nil
	return
}
func (mini Miniature) Rarity() string {
	return mini.rarity
}
func (mini Miniature) NextMiniID() string {
	return mini.nextMiniID
}
func (mini Miniature) PrevMiniID() string {
	return mini.prevMiniID
}
func (mini Miniature) String() string {
	return fmt.Sprintf("[%s] %s (%s, %s)", mini.id, mini.name, mini.set, mini.collectorNumber)
}

func miniComparator(m1 *Miniature, m2 *Miniature) bool {
	return m1.CollectorNumberAsInt() < m2.CollectorNumberAsInt()
}

// GetMiniatureByID retrieves a miniature by it's ID
func GetMiniatureByID(id string) (*Miniature, error) {
	recordIndex, exists := idToIndex[strings.ToLower(id)]
	if !exists {
		return nil, fmt.Errorf("Unable to find miniature with ID '%s'", id)
	}
	d := data[recordIndex]
	return &d, nil
}

// GetMiniaturesBySet retrieves all of the miniatures that correspond
// to the given set code.
func GetMiniaturesBySet(setCode string) ([]*Miniature, error) {
	set, exists := setToMinis[strings.ToLower(setCode)]
	if !exists {
		return nil, fmt.Errorf("Unable to find set with code '%s'", setCode)
	}
	return set, nil
}

func loadDataFromFile() [][]string {
	filepath := os.Getenv("DATA")
	if len(filepath) == 0 {
		log.Fatalf("Environment variable DATA not found. Please specify the location of the Dreamblade data as a CSV file.")
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Unable to open data file: %s\n%s", filepath, err)
	}
	r := csv.NewReader(file)
	d, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read data file: %s\n%s", filepath, err)
	}
	return d[1:]
}

func buildIDToIndex(miniatures []Miniature) {
	idToIndex = make(map[string]int)
	setToMinis = make(map[string][]*Miniature)

	for i := range miniatures {
		miniature := &miniatures[i]
		idToIndex[miniature.id] = i
		setKey := strings.ToLower(miniature.set)

		if _, exists := setToMinis[setKey]; !exists {
			setToMinis[setKey] = []*Miniature{}
		}
		setToMinis[setKey] = append(setToMinis[setKey], &miniatures[i])
	}

	for _, miniArray := range setToMinis {

		// sort all of the sets by the collector number
		sort.Slice(miniArray, func(i, j int) bool {
			return miniComparator(miniArray[i], miniArray[j])
		})

		// set up the next, prev minis
		for i, miniPtr := range miniArray {
			if i > 0 {
				miniArray[i-1].nextMiniID = miniPtr.id
				miniPtr.prevMiniID = miniArray[i-1].id
			}
		}
	}
}

func createIDFromName(name string) string {
	s := strings.Replace(name, " ", "_", -1)
	return strings.ToLower(s)
}

func convertRecordToMiniature(records [][]string) []Miniature {
	var data []Miniature
	for _, r := range records {
		data = append(data, newMiniature(r))
	}
	return data
}

func init() {
	rawData := loadDataFromFile()
	data = convertRecordToMiniature(rawData)
	buildIDToIndex(data)
}

var data []Miniature
var idToIndex map[string]int
var setToMinis map[string][]*Miniature
