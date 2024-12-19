package models

type Application struct {
	Name    string
	Version string
}

type ApplicationStatus struct {
	Sync   bool
	Health bool
}
