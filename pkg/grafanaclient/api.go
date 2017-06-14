// Copyright © Mirantis.Inc 2017.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package grafanaclient provide a simple API to manage Grafana 4.0 DataSources and Dashboards in Go.
// It's using Grafana 4.0 REST API.
// Credits https://github.com/adejoux/grafanaclient Alain Dejoux <adejoux@djouxtech.net>.
package grafanaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const timeout = 5

// SessionInterface Interface with all method definations
type SessionInterface interface {
	DoLogon() error
	CreateDataSource(DataSource) error
	GetDataSourceName(string) (DataSource, error)
	DeleteDataSource(DataSource) error
	GetDataSourceList() ([]DataSource, error)
	GetDataSourceListID(int) (DataSource, error)
	GetUsers() ([]User, error)
	GetUserID(int) (UserID, error)
	CreateUser(AdminCreateUser) error
	DeleteUser(User) error
	GetOrgs() ([]OrgList, error)
	CreateOrg(Org) error
	GetOrgID(int) (OrgID, error)
	DeleteOrg(OrgList) error
	GetOrgUsers(OrgList) ([]OrgUserList, error)
	CreateOrgUser(AdminCreateUser, CreateUserInOrg, OrgList) error
	DeleteOrgUser(User) error
}

// GrafanaError is a error structure to handle error messages in this library
type GrafanaError struct {
	Code        int
	Description string
}

// A GrafanaMessage contains the json error message received when http request failed
type GrafanaMessage struct {
	Message string `json:"message"`
}

// Error generate a text error message.
// If Code is zero, we know it's not a http error.
func (h GrafanaError) Error() string {
	if h.Code != 0 {
		return fmt.Sprintf("HTTP %d: %s", h.Code, h.Description)
	}
	return fmt.Sprintf("ERROR: %s", h.Description)
}

// Session contains user credentials, url and a pointer to http client session.
type Session struct {
	client   *http.Client
	User     string
	Password string
	url      string
}

// A Login contains the json structure of Grafana authentication request
type Login struct {
	User     string `json:"user"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// A DataSource contains the json structure of Grafana DataSource
type DataSource struct {
	ID                int    `json:"ID"`
	OrgID             int    `json:"orgID"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Access            string `json:"access"`
	URL               string `json:"url"`
	Password          string `json:"password"`
	User              string `json:"user"`
	Database          string `json:"database"`
	BasicAuthUser     string `json:"basicAuthUser"`
	BasicAuthPassword string `json:"basicAuthPassword"`
	BasicAuth         bool   `json:"basicAuth"`
	IsDefault         bool   `json:"isDefault"`
}

// OrgUserList Get Users in organization
type OrgUserList struct {
	OrgID  int    `json:"orgID"`
	UserID int    `json:"userID"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

// CreateUserInOrg Creates User in organization
type CreateUserInOrg struct {
	LoginOrEmail string `json:"loginOrEmail"`
	Role         string `json:"role"`
}

// OrgList Get organization list
type OrgList struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

// Org Create organization
type Org struct {
	Name string `json:"name"`
}

// OrgID Org by ID List
type OrgID struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Address AddressJSON
}

// AddressJSON has details of organization
type AddressJSON struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

// OrgUsers List of all users in organization
type OrgUsers struct {
	OrgID  string `json:"orgID"`
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

// AdminCreateUser Creates User
type AdminCreateUser struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"Required"`
}

// User gets Users List
type User struct {
	ID      int    `json:"ID"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	IsAdmin bool   `json:"isAdmin"`
}

// UserID gets Users by ID List
type UserID struct {
	OrgID          int    `json:"orgID"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Login          string `json:"login"`
	Theme          string `json:"theme"`
	IsGrafanaAdmin bool   `json:"isGrafanaAdmin"`
}

// NewSession It returns a Session struct pointer.
func NewSession(user string, password string, url string) *Session {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Session{client: &http.Client{Jar: jar, Timeout: time.Second * timeout}, User: user, Password: password, url: url}
}

// httpRequest handle the request to Grafana server.
//It returns the response body and a error if something went wrong
func (s *Session) httpRequest(method string, url string, body io.Reader) (result io.Reader, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return result, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return result, err
	}

	//	defer response.Body.Close()
	if response.StatusCode != 200 {
		dec := json.NewDecoder(response.Body)
		var gMess GrafanaMessage
		err := dec.Decode(&gMess)
		if err != nil {
         	       return result, err
	        }

		return result, GrafanaError{response.StatusCode, gMess.Message}
	}
	result = response.Body
	return
}

// DoLogon uses  a new http connection using the credentials stored in the Session struct.
// It returns a error if it cannot perform the login.
func (s *Session) DoLogon() (err error) {
	reqURL := s.url + "/login"

	login := Login{User: s.User, Password: s.Password}
	jsonStr, err := json.Marshal(login)
	if err != nil {
                return
        }

	_, err = s.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// CreateDataSource creates a Grafana DataSource.
// It take a DataSource struct in parameter.
// It returns a error if it cannot perform the creation.
func (s *Session) CreateDataSource(ds DataSource) (err error) {
	reqURL := s.url + "/api/datasources"

	jsonStr, err := json.Marshal(ds)
	if err != nil {
                return
        }

	_, err = s.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// GetDataSourceName get a existing DataSource by name.
// It return a DataSource struct.
// It returns a error if a problem occurs when trying to retrieve the DataSource.
func (s *Session) GetDataSourceName(name string) (ds DataSource, err error) {
	dslist, err := s.GetDataSourceList()
	if err != nil {
		return
	}

	for _, elem := range dslist {
		if elem.Name == name {
			ds = elem
		}
	}
	return
}

// DeleteDataSource deletes a Grafana DataSource.
// It take a existing DataSource struct in parameter.
// It returns a error if it cannot perform the deletion.
func (s *Session) DeleteDataSource(ds DataSource) (err error) {

	reqURL := fmt.Sprintf("%s/api/datasources/%d", s.url, ds.ID)

	jsonStr, err := json.Marshal(ds)
	if err != nil {
                return
        }

	_, err = s.httpRequest("DELETE", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// GetDataSourceList return a list of existing Grafana DataSources.
// It return a array of DataSource struct.
// It returns a error if it cannot get the DataSource list.
func (s *Session) GetDataSourceList() (ds []DataSource, err error) {
	reqURL := s.url + "/api/datasources"

	body, err := s.httpRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&ds)
	return
}

// GetDataSourceListID by ID returns single Grafana DataSources.
// It return a array of DataSource struct.
// It returns a error if it cannot get the DataSource list.
func (s *Session) GetDataSourceListID(ID int) (ds DataSource, err error) {
	reqURL := fmt.Sprintf("%s/api/datasources/%d", s.url, ID)

	body, err := s.httpRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&ds)
	return
}

// GetUsers returns list of users
func (s *Session) GetUsers() (user []User, err error) {
	reqURL := s.url + "/api/users"
	body, err := s.httpRequest("GET", reqURL, nil)

	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&user)
	return
}

// GetUserID Get User by ID
func (s *Session) GetUserID(ID int) (userID UserID, err error) {
	reqURL := fmt.Sprintf("%s/api/users/%d", s.url, ID)
	body, err := s.httpRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&userID)
	return
}

// CreateUser creates a user
func (s *Session) CreateUser(user AdminCreateUser) (err error) {
	reqURL := s.url + "/api/admin/users"
	jsonStr, err := json.Marshal(user)
	if err != nil {
                return
        }
        
	_, err = s.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// DeleteUser Delete the user with given id
func (s *Session) DeleteUser(user User) (err error) {
	reqURL := fmt.Sprintf("%s/api/admin/users/%d", s.url, user.ID)
	jsonStr, err := json.Marshal(user)
	if err != nil {
                return
        }

	_, err = s.httpRequest("DELETE", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// GetOrgs returns list of organizations
func (s *Session) GetOrgs() (org []OrgList, err error) {
	reqURL := s.url + "/api/orgs"
	body, err := s.httpRequest("GET", reqURL, nil)

	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&org)
	return
}

// CreateOrg creates a organization
func (s *Session) CreateOrg(org Org) (err error) {
	reqURL := s.url + "/api/orgs"
	jsonStr, err := json.Marshal(org)
	if err != nil {
                return
        }

	_, err = s.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// GetOrgID Get Org by ID
func (s *Session) GetOrgID(ID int) (orgID OrgID, err error) {
	reqURL := fmt.Sprintf("%s/api/orgs/%d", s.url, ID)
	body, err := s.httpRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&orgID)
	return
}

// DeleteOrg Delete the organization with given id
func (s *Session) DeleteOrg(org OrgList) (err error) {
	reqURL := fmt.Sprintf("%s/api/orgs/%d", s.url, org.ID)
	jsonStr, err := json.Marshal(org)
	if err != nil {
                return
        }

	_, err = s.httpRequest("DELETE", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// GetOrgUsers gets Users in Organisation
func (s *Session) GetOrgUsers(orgID OrgList) (org []OrgUserList, err error) {
	reqURL := fmt.Sprintf("%s/api/orgs/%d/users", s.url, orgID.ID)
	body, err := s.httpRequest("GET", reqURL, nil)

	if err != nil {
		return
	}
	dec := json.NewDecoder(body)
	err = dec.Decode(&org)
	return
}

// CreateOrgUser Add User in Organisation
func (s *Session) CreateOrgUser(user AdminCreateUser, createOrguser CreateUserInOrg, org OrgList) (err error) {
	// Create a user using Admin api and then add that user to organization
	reqUserURL := s.url + "/api/admin/users"
	jsonUsrStr, err := json.Marshal(user)
	if err != nil {
                return
        }

	_, err = s.httpRequest("POST", reqUserURL, bytes.NewBuffer(jsonUsrStr))
	if err != nil {
		return
	}

	reqURL := fmt.Sprintf("%s/api/orgs/%d/users", s.url, org.ID)
	jsonStr, err := json.Marshal(createOrguser)
	if err != nil {
                return
        }

	_, err = s.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr))

	return
}

// DeleteOrgUser Delete User in Organisation
func (s *Session) DeleteOrgUser(user User) (err error) {
	// Deleting the user through admin api deletes that user from organization
	reqURL := fmt.Sprintf("%s/api/admin/users/%d", s.url, user.ID)
	jsonStr, err := json.Marshal(user)
	if err != nil {
                return
        }

	_, err = s.httpRequest("DELETE", reqURL, bytes.NewBuffer(jsonStr))

	return
}
