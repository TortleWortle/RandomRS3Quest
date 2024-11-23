package main

import (
	"RandomRS3Quest/runemetrics"
	"RandomRS3Quest/templates"
	"RandomRS3Quest/wheeler"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

func main() {
	m := http.NewServeMux()

	fetcher, err := runemetrics.NewPlayerFetcher()
	if err != nil {
		log.Println("could create player fetcher, exiting.")
		return
	}

	m.Handle("GET /", templ.Handler(templates.Welcome()))
	m.HandleFunc("GET /generate", func(w http.ResponseWriter, request *http.Request) {
		reqID := uuid.New()
		logger := slog.With(slog.String("requestID", reqID.String()))
		queryNames := request.URL.Query().Get("names")
		usernames := strings.Split(queryNames, "\r\n")

		if len(usernames) < 1 {
			logger.Info("denying request due to no names given")

			http.Redirect(w, request, request.Referer(), 301)
			return
		}

		if len(usernames) > 5 {
			logger.Info("denying request due to too many names given", slog.Int("amount", len(usernames)))

			http.Redirect(w, request, request.Referer(), 301)
			return
		}
		for _, name := range usernames {
			if len(name) > 12 || len(name) < 1 {
				// They know what they did
				logger.Info("denying request due to bad name", slog.String("rsn", name))
				http.Redirect(w, request, request.Referer(), 301)
				return
			}
		}

		availableQuestPairs := make(map[string][]runemetrics.Quest)
		questDigits := make(map[string]int)

		logger.Info("starting fetcher", slog.String("rsns", strings.Join(usernames, ",")))
		for _, name := range usernames {
			logger.Info("fetching metrics", slog.String("rsn", name))
			quests, err := fetcher.FetchUserQuests(name)
			if err != nil {
				slog.Info("failed fetching metrics", slog.String("rsn", name))
				continue
			}
			logger.Info("fetched quests", slog.String("rsn", name), slog.Int("amount", len(quests)))
			for _, quest := range quests {
				if quest.UserEligible && quest.Status != "COMPLETED" {
					// hacky but here we insert the difficulty indicator
					quest.Title = prettyTitle(quest)
					availableQuestPairs[name] = append(availableQuestPairs[name], quest)
					questDigits[quest.Title] += 1
				}
			}
			logger.Info("eligible quests", slog.String("rsn", name), slog.Int("amount", len(availableQuestPairs[name])))
		}

		var availableQuestTitles []string

		for title, digit := range questDigits {
			if digit == len(usernames) {
				availableQuestTitles = append(availableQuestTitles, title)
			}
		}

		logger.Info("combined eligible quests", slog.String("rsns", strings.Join(usernames, ",")), slog.Int("amount", len(availableQuestTitles)))

		logger.Info("generating wheel")
		wheel := wheeler.GenerateWheel(availableQuestTitles)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", strings.Join(usernames, "_")+".wheel"))
		logger.Info("encoding wheel")
		err := json.NewEncoder(w).Encode(wheel)
		if err != nil {
			logger.Error("could not encode wheel", slog.String("rsns", strings.Join(usernames, ",")))
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
