package helpers

import (
	"fmt"
	"regexp"
	"strconv"
)

func EngineDelta(eng1Str, eng2Str string) string {
	eng1, eng2 := engineFloat(eng1Str), engineFloat(eng2Str)

	if eng1 == -1 || eng2 == -1 {
		return "err"
	}

	delta := eng2 - eng1

	arrow1, arrow2 := arrows(delta)

	return fmt.Sprintf("%s %.1f L %s", arrow1, delta, arrow2)
}

func PowerDelta(hp1, hp2 int) string {
	delta := hp2 - hp1

	arrow1, arrow2 := arrows(float64(delta))

	return fmt.Sprintf("%s %d Hp %s", arrow1, delta, arrow2)
}

func arrows(delta float64) (string, string) {
	if delta < 0 {
		return "⬆", "⬇"
	} else if delta == 0 {
		return "", ""
	} else {
		return "⬇", "⬆"
	}
}

func engineFloat(engineStr string) float64 {
	if engineStr == "Electric Motor" {
		return 0
	}

	// Extract the numeric part before "L"
	re := regexp.MustCompile(`(\d+\.\d+)L`)
	matches := re.FindStringSubmatch(engineStr)
	if len(matches) < 2 {
		fmt.Println("Error extracting engine size:", engineStr)
		return -1
	}

	engine, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		fmt.Println("Error converting to float:", err)
		return -1
	}

	return engine
}
