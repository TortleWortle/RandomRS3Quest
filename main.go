package main

import (
	"RandomRS3Quest/runemetrics"
	"RandomRS3Quest/templates"
	"RandomRS3Quest/wheeler"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"log"
	"net/http"
	"strings"
)

func main() {
	//fetcher, err := runemetrics.NewPlayerFetcher()
	//if err != nil {
	//	log.Println("could create player fetcher, exiting.")
	//	return
	//}

	m := http.NewServeMux()

	fetcher, err := runemetrics.NewPlayerFetcher()
	if err != nil {
		log.Println("could create player fetcher, exiting.")
		return
	}

	m.Handle("GET /", templ.Handler(templates.Welcome()))
	m.HandleFunc("GET /generate", func(w http.ResponseWriter, request *http.Request) {
		queryNames := request.URL.Query().Get("names")
		usernames := strings.Split(queryNames, "\r\n")
		if len(usernames) < 1 {
			log.Printf("denying request due to no names given")

			http.Redirect(w, request, request.Referer(), 301)
			return
		}

		if len(usernames) > 5 {
			log.Printf("denying request due to too many names given")

			http.Redirect(w, request, request.Referer(), 301)
			return
		}
		for _, name := range usernames {
			if len(name) > 12 || len(name) < 1 {
				// They know what they did
				log.Printf("denying request due to bad name: %v", name)
				http.Redirect(w, request, request.Referer(), 301)
				return
			}
		}

		availableQuestPairs := make(map[string][]runemetrics.Quest)
		questDigits := make(map[string]int)

		for _, name := range usernames {
			log.Printf("fetching metrics for user %s", name)
			quests, err := fetcher.FetchUserQuests(name)
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

		wheel := wheeler.GenerateWheel(availableQuestTitles)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", strings.Join(usernames, "_")+".wheel"))
		err := json.NewEncoder(w).Encode(wheel)
		if err != nil {
			log.Printf("could not encode wheel for %v", usernames)
			return
		}
	})
	err = http.ListenAndServe(":8080", m)
	if err != nil {
		log.Fatalf("listenAndServe: %v\n", err)
		return
	}
}

func prettyTitle(quest runemetrics.Quest) string {
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
