package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (game *Game) loadScores() {
	data, err := os.ReadFile(GetScoresPath())
	if err != nil {
		if os.IsNotExist(err) {
			game.scores = []int{}
		}
		fmt.Printf("Error loading scores: %s\n", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	scores := make([]int, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		score, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Error parsing scores: %s\n", err)
		}
		scores = append(scores, score)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	game.scores = scores
}

func (game *Game) saveScores() {
	os.MkdirAll(GetGameCacheDir(), os.ModePerm)
	file, err := os.OpenFile(GetScoresPath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	for _, score := range game.scores {
		scoreStr := fmt.Sprintf("%d\n", score)
		_, err = file.WriteString(scoreStr)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			return
		}
	}
}

func GetScoresPath() string {
	return fmt.Sprintf("%s/scores.txt", GetGameCacheDir())
}

func (game *Game) appendCurrentScore() {
	if game.score == 0 {
		return
	}

	game.scores = append(game.scores, game.score)
	sort.Sort(sort.Reverse(sort.IntSlice(game.scores)))
	go game.saveScores()
}
