package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"sort"
	"time"
)

const (
	// MMRetries defines how many times the player keep retrying on match making
	MMRetries = 3
	// MMWaitTime defines how long the player will wait if there is no other player in player pool
	MMWaitTime = 3000 * time.Millisecond
	// MMTimeout determines how long the player should be in the player pool
	MMTimeout = 120 * time.Second
)

// MatchMakingFunc calculates the score based on the given profile target and candidate
type MatchMakingFunc func(target *zb_data.PlayerProfile, candidate *zb_data.PlayerProfile) float64

// mmf is the gobal match making function that calculates the match making score
var mmf MatchMakingFunc = func(target *zb_data.PlayerProfile, candidate *zb_data.PlayerProfile) float64 {
	targetData := target.RegistrationData
	candidateData := candidate.RegistrationData
	// map tag to the same tag group
	if len(targetData.Tags) > 0 || len(candidateData.Tags) > 0 {
		if compareTags(targetData.Tags, candidateData.Tags) {
			return 1
		} else {
			return 0
		}
	}

	// version
	if targetData.Version != candidateData.Version {
		return 0
	}

	// backend side game logic
	if targetData.UseBackendGameLogic != candidateData.UseBackendGameLogic {
		return 0
	}

	// debug cheats
	if targetData.DebugCheats.Enabled != candidateData.DebugCheats.Enabled {
		return 0
	}

	return 1
}

// compareTags compares string slices so order of string matters
func compareTags(tag1, tag2 []string) bool {
	if len(tag1) != len(tag2) {
		return false
	}
	for i, v := range tag1 {
		if v != tag2[i] {
			return false
		}
	}
	return true
}

// PlayerScore simply maintains the player id and score tuple
type PlayerScore struct {
	score float64
	id    string
}

type byPlayersScore []*PlayerScore

func (p byPlayersScore) Len() int { return len(p) }

func (p byPlayersScore) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p byPlayersScore) Less(i, j int) bool {
	return p[i].score > p[j].score
}

func sortByPlayerScore(ps []*PlayerScore) []*PlayerScore {
	sort.Sort(byPlayersScore(ps))
	return ps
}

func findPlayerProfileByID(pool *zb_data.PlayerPool, id string) *zb_data.PlayerProfile {
	for _, pp := range pool.PlayerProfiles {
		if pp.RegistrationData.UserId == id {
			return pp
		}
	}
	return nil
}

func removePlayerFromPool(pool *zb_data.PlayerPool, id string) *zb_data.PlayerPool {
	var newpool zb_data.PlayerPool
	for _, pp := range pool.PlayerProfiles {
		if pp.RegistrationData.UserId != id {
			newpool.PlayerProfiles = append(newpool.PlayerProfiles, pp)
		}
	}
	return &newpool
}
