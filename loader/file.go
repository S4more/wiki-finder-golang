package file

import (
	"os"
	"slices"

	timer "github.com/s4more/go-wiki/utils"
	"github.com/sugawarayuuta/sonnet"
)

func ReadFile() []byte {
	data, err := os.ReadFile("/mnt/c/Users/Samore/web-dev-fall-2022-wiki/search_index/test_links.json")
	if err != nil {
		panic(err)
	}
	return data
}

func UnmarshalMessage(message []byte) [][]uint32 {
	defer timer.Timer("Unmarshall")()
	var links [][]uint32
	err := sonnet.Unmarshal(message, &links)

	for _, v := range links {
		slices.Sort(v)
		slices.Compact(v)
	}

	if err != nil {
		panic(err)
	}

	return links
}
