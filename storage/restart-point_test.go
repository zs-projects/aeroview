package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromSortedSlice(t *testing.T) {
	cities := []string{
		"Adelanto",
		"Agoura",
		"Alameda",
		"Albany",
		"Alhambra",
		"Aliso",
		"AlturasCounty",
		"Amador",
		"American",
		"Anaheim",
		"Anderson",
		"Angels",
		"Antioch",
		"Apple",
		"Arcadia",
		"Arcata",
		"Arroyo",
		"Artesia",
		"Arvin",
		"Atascadero",
		"Atherton",
		"Atwater",
		"AuburnCounty",
		"Avalon",
		"Avenal",
		"Azusa",
		"BakersfieldCounty",
		"Baldwin",
		"Banning",
		"Barstow",
		"Beaumont",
		"Bell",
		"Bell",
		"Bellflower",
		"Belmont",
		"Belvedere",
		"Benicia",
		"Berkeley",
		"Beverly",
		"Big",
		"Biggs",
		"Bishop",
		"Blue",
		"Blythe",
		"Bradbury",
		"Brawley",
		"Brea",
		"Brentwood",
		"Brisbane",
		"Buellton",
		"Buena",
		"Burbank",
		"Burlingame",
	}
	r := FromSortedSlice(cities)
	raw := r.ToRawStrings()
	assert.Equal(t, cities, raw)
}
