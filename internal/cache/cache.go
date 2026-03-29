package cache

import (
	"encoding/gob"
	"os"
	"strings"
	"time"

	"github.com/maypress/CurEx/internal/api"
)

const CacheFileName = "./cache/rates.cache"

var cacheTime time.Time
var cachedRates *api.DailyRates

type CacheData struct {
	Rates api.DailyRates
	Time  time.Time
}

func init() {
	gob.Register(api.DailyRates{})
	loadCache()
}

func loadCache() {
	data, err := os.ReadFile(CacheFileName)
	if err != nil {
		return
	}

	var cd CacheData
	if err := gob.NewDecoder(strings.NewReader(string(data))).Decode(&cd); err != nil {
		return
	}

	if time.Since(cd.Time) < time.Hour {
		cachedRates = &cd.Rates
		cacheTime = cd.Time
	}
}

func SaveRates(rates *api.DailyRates) {
	if rates == nil {
		return
	}
	
	cd := CacheData{Rates: *rates, Time: time.Now()}
	file, err := os.Create(CacheFileName)
	if err != nil {
		return
	}
	defer file.Close()
	
	gob.NewEncoder(file).Encode(cd)
	cachedRates = rates
	cacheTime = time.Now()
}

func GetRates() *api.DailyRates {
	return cachedRates
}

func GetCacheTime() time.Time {
	return cacheTime
}