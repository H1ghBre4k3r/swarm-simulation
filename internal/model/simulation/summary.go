package simulation

import (
	"fmt"
	"time"
)

type Summary struct {
	Participants uint64           `json:"participants"`
	Collisions   CollisionSummary `json:"collisions"`
	Runtime      RuntimeSummary   `json:"runtime"`
}

func (s Summary) String() string {
	summary := `
-----------------------------
		SUMMARY
-----------------------------
N° of Participants:		%v
%v
%v
`
	return fmt.Sprintf(summary, s.Participants, s.Collisions, s.Runtime)
}

type CollisionSummary struct {
	Total   uint64 `json:"total"`
	Most    uint64 `json:"most"`
	Least   uint64 `json:"least"`
	Average uint64 `json:"avg"`
}

func (c CollisionSummary) String() string {
	return fmt.Sprintf(`
Total Collisions: 		%v
Most Collisions: 		%v
Least Collisions: 		%v
Avg Collisions: 		%v`, c.Total, c.Most, c.Least, c.Average)
}

type RuntimeSummary struct {
	Ticks             uint64        `json:"ticks"`
	Duration          time.Duration `json:"duration"`
	AverageTickLength time.Duration `json:"avg"`
}

func (r RuntimeSummary) String() string {
	return fmt.Sprintf(`
Total ticks: 			%v
Duration:			%v
Avg Tick Length: 		%v`, r.Ticks, r.Duration, r.AverageTickLength)
}

func GenerateSummary(s *Simulation) Summary {
	most := uint64(0)
	least := uint64(0)
	total := uint64(0)
	participants := s.entities.Get()
	for i, e := range participants {
		collisions := e.GetCollisions()
		total += collisions
		if i == 0 {
			least = collisions
			most = collisions
		} else {
			if collisions > most {
				most = collisions
			} else if collisions < least {
				least = collisions
			}
		}
	}

	return Summary{
		Participants: uint64(len(participants)),
		Collisions: CollisionSummary{
			Total:   total,
			Most:    most,
			Least:   least,
			Average: total / uint64(len(participants)),
		},
		Runtime: RuntimeSummary{
			Ticks:             s.ticks,
			Duration:          s.duration,
			AverageTickLength: time.Duration(s.duration / time.Duration(s.ticks)),
		},
	}
}
