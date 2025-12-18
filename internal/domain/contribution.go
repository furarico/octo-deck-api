package domain

import "time"

type Contribution struct {
	Date  time.Time
	Count int
}

func NewContribution(date time.Time, count int) *Contribution {
	return &Contribution{
		Date:  date,
		Count: count,
	}
}
