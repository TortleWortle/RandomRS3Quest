package runemetrics

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

func NewPlayerFetcher() (PlayerFetcher, error) {
	pf := PlayerFetcher{
		Cache:     make(map[string]QuestCacheItem),
		cacheTime: time.Minute * 5,
		cacheLock: &sync.RWMutex{},
	}
	return pf, nil
}

type PlayerFetcher struct {
	Cache     QuestCache
	cacheFile io.ReadWriteSeeker
	cacheTime time.Duration
	cacheLock *sync.RWMutex
}

func (p *PlayerFetcher) fetchUserQuestsFromCache(username string) ([]Quest, error) {
	p.cacheLock.RLock()
	defer p.cacheLock.RUnlock()
	cachedQuests, ok := p.Cache[username]
	if !ok {
		return nil, errors.New("user not cached")
	}
	if cachedQuests.Time.Add(p.cacheTime).Before(time.Now()) {
		// expired
		return nil, errors.New("user Cache expired")
	}

	return cachedQuests.Quests, nil
}

func (p *PlayerFetcher) updateUserCache(username string, quests []Quest) error {
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()
	p.Cache[username] = QuestCacheItem{
		Quests: quests,
		Time:   time.Now(),
	}
	return nil
}

func (p *PlayerFetcher) FetchUserQuests(username string) ([]Quest, error) {
	cachedQuests, err := p.fetchUserQuestsFromCache(username)
	logger := slog.With(slog.String("rsn", username))
	if err == nil {
		logger.Info("user was cached")
		return cachedQuests, nil
	}
	logger.Info("user was NOT cached", slog.String("reason", err.Error()))
	u := "https://apps.runescape.com/runemetrics/quests?user=" + username
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var decoded ResponseBody
	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		logger.Error("could not decode body", slog.Int("status", res.StatusCode))
		return nil, err
	}

	err = p.updateUserCache(username, decoded.Quests)
	if err != nil {
		logger.Error("could not update usercache")
		return nil, err
	}

	return decoded.Quests, nil
}

type QuestCache map[string]QuestCacheItem
type QuestCacheItem struct {
	Quests []Quest
	Time   time.Time
}

type ResponseBody struct {
	Quests   []Quest `json:"quests"`
	LoggedIn string  `json:"loggedIn"`
}

type Quest struct {
	Title        string `json:"title"`
	Status       string `json:"status"`
	Difficulty   int    `json:"difficulty"`
	Members      bool   `json:"members"`
	QuestPoints  int    `json:"questPoints"`
	UserEligible bool   `json:"userEligible"`
}
