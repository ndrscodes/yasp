package db

import (
	"fmt"
	"time"
)

const (
	OPERATIONAL = iota
	DEGRADED    = iota
	DOWN        = iota
)

type System struct {
	Id     uint
	Name   string
	Status int
}

type Incident struct {
	Id        uint
	Title     string
	Details   string
	CreatedAt time.Time
}

type Update struct {
	Id        uint
	Title     string
	Details   string
	Status    int
	CreatedAt time.Time
}

type SystemWithIncidents struct {
	System
	Incidents []Incident
}

func NewSystem(name string) System {
	return System{
		Name: name,
	}
}

func NewIncident(title, details string) Incident {
	return Incident{
		Title:     title,
		Details:   details,
		CreatedAt: time.Now(),
	}
}

func NewUpdate(title, details string, status int) Update {
	return Update{
		Title:     title,
		Details:   details,
		Status:    status,
		CreatedAt: time.Now(),
	}
}

func ToStatus(dbStatus string) (int, error) {
	switch dbStatus {
	case "operational":
		return OPERATIONAL, nil
	case "degraded":
		return DEGRADED, nil
	case "down":
		return DOWN, nil
	default:
		return OPERATIONAL, fmt.Errorf("unknown status: %s", dbStatus)
	}
}

func FromStatus(status int) (string, error) {
	switch status {
	case OPERATIONAL:
		return "operational", nil
	case DEGRADED:
		return "degraded", nil
	case DOWN:
		return "down", nil
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}
}
