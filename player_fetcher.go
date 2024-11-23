package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func NewPlayerFetcher(cache io.ReadWriteSeeker) (PlayerFetcher, error) {
	pf := PlayerFetcher{
		Cache:     make(map[string]QuestCacheItem),
		cacheFile: cache,
		cacheTime: time.Minute * 5,
		cacheLock: &sync.RWMutex{},
	}

	err := pf.readCache()
	if err != nil {
		return PlayerFetcher{}, err
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

func (p *PlayerFetcher) readCache() error {
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()
	_, err := p.cacheFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = json.NewDecoder(p.cacheFile).Decode(&p.Cache)

	return err
}

func (p *PlayerFetcher) updateUserCache(username string, quests []Quest) error {
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()
	p.Cache[username] = QuestCacheItem{
		Quests: quests,
		Time:   time.Now(),
	}

	_, err := p.cacheFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = json.NewEncoder(p.cacheFile).Encode(p.Cache)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayerFetcher) fetchUserQuests(username string) ([]Quest, error) {
	cachedQuests, err := p.fetchUserQuestsFromCache(username)
	if err == nil {
		log.Printf("user %s was cached", username)
		return cachedQuests, nil
	}
	log.Printf("user %s NOT cached, fetching (%v)", username, err)
	u := "https://apps.runescape.com/runemetrics/quests?user=" + username
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var decoded RuneMetricsResponseBody
	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		return nil, err
	}

	err = p.updateUserCache(username, decoded.Quests)
	if err != nil {
		return nil, err
	}

	return decoded.Quests, nil
}

type QuestCache map[string]QuestCacheItem
type QuestCacheItem struct {
	Quests []Quest
	Time   time.Time
}

type RuneMetricsResponseBody struct {
	Quests []Quest `json:"quests"`
	//LoggedIn string  `json:"loggedIn"`
}

type Quest struct {
	Title        string `json:"title"`
	Status       string `json:"status"`
	Difficulty   int    `json:"difficulty"`
	Members      bool   `json:"members"`
	QuestPoints  int    `json:"questPoints"`
	UserEligible bool   `json:"userEligible"`
}
