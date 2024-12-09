package helpers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func UpdateCarCounts(result []int, r *http.Request, w http.ResponseWriter) {
	cookie, err := r.Cookie("car_counts")
	carCounts := make(map[int]int)

	if err == nil {
		carCounts = ParseCarCounts(cookie.Value)
	}

	for _, id := range result {
		carCounts[id]++
	}

	cookieValue := formatCarCounts(carCounts)
	http.SetCookie(w, &http.Cookie{
		Name:  "car_counts",
		Value: cookieValue,
		Path:  "/",
	})
}

func ParseCarCounts(cookieValue string) map[int]int {
	carCounts := make(map[int]int)
	if cookieValue == "" {
		return carCounts
	}
	pairs := strings.Split(cookieValue, "|")

	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			continue
		}
		id, err1 := strconv.Atoi(kv[0])
		count, err2 := strconv.Atoi(kv[1])
		if err1 != nil || err2 != nil {
			continue
		}
		carCounts[id] = count
	}

	return carCounts
}

func formatCarCounts(carCounts map[int]int) string {
	var parts []string
	for id, count := range carCounts {
		parts = append(parts, fmt.Sprintf("%d:%d", id, count))
	}
	return strings.Join(parts, "|")
}
