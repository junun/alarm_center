package repo

import (
	"net/http"
	"time"
)

type SonarClient struct{
	Name		string
	JwtSession 	string
	ExpireTime 	time.Time
	Path 		string
}

type SonarComponentResponse struct {
	AlertStatus 		string
	Bugs				string
	NewBugs				string
	Vulnerabilities		string
	NewVulnerabilities	string
	CodeSmells			string
	NewCodeSmells		string
	SecurityHotspots	string
	NewSecurityHotspots	string
}

type Period struct {
	Index int
	Value string
	BestValue 	bool
}
type Measures struct {
	Metric		string
	Value		string
	BestValue 	bool
	Periods 	[]Period
	Period		Period
}

type Component struct {
	Key 		string
	Name 		string
	Qualifier	string
	Measures	[]Measures
}

type SonarResult struct {
	Component 	Component
}

type SonarRepository interface {
	NewSonarClient() error
	SonarReq(method, uri string) (*http.Request, error)
	NewSonarReq(method, uri string) (*http.Request, error)
	ValidateUser(req *http.Request) bool
	GetMeasuresComponent(uri, component string) (SonarComponentResponse, error)
}