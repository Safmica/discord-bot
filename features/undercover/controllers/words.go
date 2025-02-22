package controllers

import (
	"math/rand"

	models "github.com/Safmica/discord-bot/features/undercover/models"
)

func AssignRandomWords() bool {
	var unusedIndexes []int
	for i, entry := range models.ActiveGame.Words {
		if !entry.Used {
			unusedIndexes = append(unusedIndexes, i)
		}
	}

	if len(unusedIndexes) == 0 {
		return false
	}

	randomIndex := unusedIndexes[rand.Intn(len(unusedIndexes))]

	words := models.ActiveGame.Words[randomIndex].Word
	models.ActiveGame.Words[randomIndex].Used = true

	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	models.ActiveGame.CivilianWords = words[0]
	models.ActiveGame.UndercoverWords = words[1]

	return true
}
