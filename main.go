package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	usernames := []string{
		"TinierTortle",
		"Grimburd",
		"Aza Saindu",
	}

	cacheFile, err := os.OpenFile("questcache.dat", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("could not open/create cachefile, exiting.")
		return
	}

	fetcher, err := NewPlayerFetcher(cacheFile)
	if err != nil {
		log.Println("could create player fetcher, exiting.")
		return
	}

	availableQuestPairs := make(map[string][]Quest)
	questDigits := make(map[string]int)

	for _, name := range usernames {
		log.Printf("fetching metrics for user %s", name)
		quests, err := fetcher.fetchUserQuests(name)
		if err != nil {
			log.Printf("failed fetching metrics for user %s", name)
			continue
		}
		for _, quest := range quests {
			if quest.UserEligible && quest.Status != "COMPLETED" {
				// hacky but here we insert the difficulty indicator
				quest.Title = prettyTitle(quest)
				availableQuestPairs[name] = append(availableQuestPairs[name], quest)
				questDigits[quest.Title] += 1
			}
		}
	}

	var availableQuestTitles []string

	for title, digit := range questDigits {
		if digit == len(usernames) {
			availableQuestTitles = append(availableQuestTitles, title)
		}
	}

	wheel := GenerateWheel(availableQuestTitles)

	wheelFile, err := os.OpenFile("./quests.wheel", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("could not open/create wheelFile, exiting.")
		return
	}

	err = wheelFile.Truncate(0)
	if err != nil {
		log.Println("could not truncate wheelFile, exiting.")
		return
	}

	err = json.NewEncoder(wheelFile).Encode(wheel)
	if err != nil {
		log.Println("Could not encode/save wheel file.")
		return
	}
	log.Printf("Saved %d quests in quests.wheel have fun!", len(availableQuestTitles))
}

func prettyTitle(quest Quest) string {
	title := quest.Title

	if quest.Difficulty > 5 {
		title += "🥲"
	} else {
		for i := 0; i < quest.Difficulty; i++ {
			title += "⭐"
		}
	}
	return title
}
