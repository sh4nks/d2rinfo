package controller

import (
	"d2rinfo/config"
	"d2rinfo/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	otter "github.com/maypok86/otter/v2"
)

const (
	D2EMU_TZ_API     = "https://d2emu.com/api/v1/tz"
	D2EMU_DCLONE_API = "https://d2emu.com/api/v1/dclone"
)

func NewD2RInfoController(config *config.Config, cache *otter.Cache[string, any]) *D2RInfoController {
	return &D2RInfoController{
		config: config,
		cache:  cache,
	}
}

type D2RInfoController struct {
	config *config.Config
	cache  *otter.Cache[string, any]
}

func (ctrl *D2RInfoController) GetD2RInfoData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	cacheKey := "d2rinfo"

	// Try to get combined data from cache first.
	cachedValue, found := ctrl.cache.GetIfPresent(cacheKey)
	if found {
		// If found, attempt to assert type and return.
		if data, ok := cachedValue.(map[string]any); ok {
			data["cache"] = "hit"
			json.NewEncoder(w).Encode(data)
			log.Printf("Cache hit for '%s'.", cacheKey)
			return
		}
	}

	log.Println("Fetching and generating new D2R info data (cache miss).")
	combinedData := make(map[string]any)

	if ctrl.config.Username == "" {
		log.Println("Warning: Variable D2EMU_USERNAME is not set. API calls might fail.")
	}
	if ctrl.config.Token == "" {
		log.Println("Warning: variable D2EMU_TOKEN is not set. API calls might fail.")
	}

	// Fetch TZ data
	tzData, err := utils.FetchJSON(ctx, D2EMU_TZ_API, ctrl.config.Username, ctrl.config.Token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch TZ data: %v", err), http.StatusInternalServerError)
		log.Printf("Error fetching TZ data: %v", err)
		return
	}
	combinedData["tz"] = tzData

	// Fetch Dclone data
	dcloneData, err := utils.FetchJSON(ctx, D2EMU_DCLONE_API, ctrl.config.Username, ctrl.config.Token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch Dclone data: %v", err), http.StatusInternalServerError)
		log.Printf("Error fetching Dclone data: %v", err)
		return
	}
	combinedData["dclone"] = dcloneData
	combinedData["generated_at"] = time.Now().Format(time.RFC3339)

	// Store the combined data in the cache
	ctrl.cache.Set(cacheKey, combinedData)
	combinedData["cache"] = "set"

	json.NewEncoder(w).Encode(combinedData)
}
