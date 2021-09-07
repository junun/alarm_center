package persistence

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type sonarRepo struct {
	client  *repo.SonarClient
}

func NewSonarRepository(c *repo.SonarClient) repo.SonarRepository {
	return &sonarRepo{
		client: c,
	}
}

func (s *sonarRepo) NewSonarClient() error {
	c		:= &http.Client{}
	data 	:= url.Values{}
	data.Set("login", config.SonarConfig.Login)
	data.Set("password", config.SonarConfig.Password)
	req, err := http.NewRequest("POST", config.SonarConfig.Url+db.SonarApiLogin, strings.NewReader(data.Encode()))
	if err != nil {
		return  err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.Do(req)
	if err != nil {
		return  err
	}
	defer resp.Body.Close()

	// 默认过期时间1个小时
	maxAge	:= 3600
	jwt 	:= ""
	for s, i := range resp.Header {
		if s == "Set-Cookie" {
			for _, s2 := range i {
				r := regexp.MustCompile(`(?i)XSRF`)
				if r.MatchString(s2) {
					maxAge,_ = strconv.Atoi(strings.Split(strings.Split(s2, ";")[1], "=")[1])
				}
				j := regexp.MustCompile(`(?i)JWT-SESSION`)
				if j.MatchString(s2) {
					jwt = strings.Split(strings.Split(s2, "=")[1], ";")[0]
				}
			}
		}
	}

	expire := time.Now().Add(time.Duration(maxAge) * time.Second)
	client := &repo.SonarClient{
		Name: "JWT-SESSION",
		JwtSession: jwt,
		ExpireTime: expire,
		Path: "/",
	}
	s.client = client
	return nil
}

func (s *sonarRepo) SonarReq(method, uri string) (*http.Request, error) {
	cookie := &http.Cookie{Name: s.client.Name, Value: s.client.JwtSession, Path: s.client.Path, Expires: s.client.ExpireTime}
	req, _ := http.NewRequest("GET", config.SonarConfig.Url +  db.SonarUserValidate,nil)
	req.Header.Add("cache-control", "no-cache")
	req.AddCookie(cookie)
	if time.Now().After(s.client.ExpireTime) {
		if s.ValidateUser(req) {
			req, _ = http.NewRequest(method, config.SonarConfig.Url + uri,nil)
			return req, nil
		}
	}
	return s.NewSonarReq(method, uri)
}

func (s *sonarRepo) NewSonarReq(method, uri string) (*http.Request, error) {
	err := s.NewSonarClient()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{Name: s.client.Name, Value: s.client.JwtSession, Path: s.client.Path, Expires: s.client.ExpireTime}
	req, _ := http.NewRequest(method, config.SonarConfig.Url + uri,nil)
	req.Header.Add("cache-control", "no-cache")
	req.AddCookie(cookie)
	return req, nil
}

func (s *sonarRepo) ValidateUser(req *http.Request) bool {
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if res.Status != "200" {
		return false
	}

	return true
}

func(s *sonarRepo) GetMeasuresComponent(uri, component string) (repo.SonarComponentResponse, error) {
	url :=  uri + "?component=" + component +"&metricKeys=alert_status,bugs,new_bugs," +
		"vulnerabilities,new_vulnerabilities,code_smells,new_code_smells,security_hotspots,new_security_hotspots"
	var scp repo.SonarComponentResponse

	req, err := s.SonarReq("GET", url)
	if err != nil {
		return scp, err
	}

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//if res.Status != "200" {
	//	return scp, errors.New(string(body))
	//}

	var target repo.SonarResult
	err = json.Unmarshal(body, &target)
	if err != nil {
		return scp, err
	}

	for _, measure := range target.Component.Measures {
		switch measure.Metric {
		case "alert_status":
			scp.AlertStatus = measure.Value
		case "bugs":
			scp.Bugs = measure.Value
		case "new_bugs":
			scp.NewBugs = measure.Period.Value
		case "vulnerabilities":
			scp.Vulnerabilities =measure.Value
		case "new_vulnerabilities":
			scp.NewVulnerabilities = measure.Period.Value
		case "code_smells":
			scp.CodeSmells = measure.Value
		case "new_code_smells":
			scp.NewCodeSmells = measure.Period.Value
		case "security_hotspots":
			scp.SecurityHotspots = measure.Value
		case "new_security_hotspots":
			scp.NewSecurityHotspots = measure.Period.Value
		}
	}

	return scp, nil
}

