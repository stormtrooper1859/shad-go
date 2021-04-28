// +build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

type entrance struct {
	entranceType int
	date         int
}

func ComputeLoad(guests []Guest) []Load {
	n := len(guests) * 2
	ar := make([]entrance, 0, n)

	for _, v := range guests {
		ar = append(ar,
			entrance{entranceType: 1, date: v.CheckInDate},
			entrance{entranceType: -1, date: v.CheckOutDate})
	}

	sort.Slice(ar, func(i, j int) bool {
		if ar[i].date == ar[j].date {
			return ar[i].entranceType < ar[j].date
		}
		return ar[i].date < ar[j].date
	})

	result := make([]Load, 0)
	currentVisitors := 0

	for i := 0; i < n; i++ {
		previousVisitors := currentVisitors
		currentDate := ar[i].date

		for i+1 < n && ar[i+1].date == currentDate {
			currentVisitors += ar[i].entranceType
			i++
		}
		currentVisitors += ar[i].entranceType

		if currentVisitors != previousVisitors {
			result = append(result, Load{currentDate, currentVisitors})
		}
	}

	return result
}
