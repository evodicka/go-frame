package persistence

import (
	"time"
)

type Status struct {
	CurrentImageId int
	LastSwitch     time.Time
	ImageDuration  int
}

var currentStatus = prepopulateStatus()

func prepopulateStatus() Status {
	return Status{
		CurrentImageId: -1,
		LastSwitch:     time.Unix(0, 0),
		ImageDuration:  300,
	}
}

func GetCurrentStatus() Status {
	return currentStatus
}

func UpdateImageStatus(newId int) {
	currentStatus.CurrentImageId = newId
	currentStatus.LastSwitch = time.Now()
}
