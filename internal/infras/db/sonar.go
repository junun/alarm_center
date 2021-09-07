package db

import (
	"alarm_center/internal/domain/repo"
)

const (
	SonarApiLogin = "/api/authentication/login"
	SonarUserValidate = "/api/authentication/validate"
	SonarMeasuresComponent	= "/api/measures/component"
)

func  InitSonarClient() *repo.SonarClient {
	return &repo.SonarClient{}
}

