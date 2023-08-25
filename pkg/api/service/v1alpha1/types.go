package v1alpha1

import (
	"github.com/nduyphuong/gorya/internal/models"
	"time"
)

type VersionInfo struct {
	Version      string    ` json:"version,omitempty"`
	GitCommit    string    `json:"git_commit,omitempty"`
	GitTreeDirty bool      `json:"git_tree_dirty,omitempty"`
	BuildTime    time.Time `json:"build_time,omitempty"`
	GoVersion    string    `json:"go_version,omitempty"`
	Compiler     string    `json:"compiler,omitempty"`
	Platform     string    `json:"platform,omitempty"`
}

type OkResponse struct {
	Message string `json:"message"`
}

type GetVersionInfoResponse struct {
	VersionInfo *VersionInfo `json:"version_info,omitempty"`
}

type GetTimeZoneResponse struct {
	TimeZones []string `json:"Timezones,omitempty"`
}

type AddScheduleRequest struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"displayname,omitempty"`
	Dtype       string  `json:"dtype"`
	Corder      bool    `json:"corder"`
	Shape       []int   `json:"shape"`
	NdArray     [][]int `json:"__ndarray__"`
	TimeZone    string  `json:"timezone"`
}

type AddPolicyRequest struct {
	Name         string              `json:"name"`
	DisplayName  string              `json:"displayname,omitempty"`
	Tags         []map[string]string `json:"tags"`
	Projects     []models.Project    `json:"projects"`
	ScheduleName string              `json:"schedulename"`
}

type GetScheduleResponse struct {
	Name        string  ` json:"name" gorm:"size:191;unique"`
	DisplayName string  `json:"displayname"`
	TimeZone    string  `json:"timezone"`
	Dtype       string  `json:"dtype"`
	Corder      bool    `json:"Corder"`
	Shape       []int   `json:"Shape"`
	NdArray     [][]int `json:"__ndarray__"`
}

type ListResponsesVerbose []ListResponseVerbose

type ListResponseVerbose struct {
	Name string `json:"name"`
	//compatible with UI
	DisplayName string `json:"displayName"`
}

type ListResponse []string

type GetPolicyResponse struct {
	models.Policy
}

type ChangeStateRequest struct {
	Action        int    `json:"action"`
	Project       string `json:"project"`
	CredentialRef string `json:"credentialRef"`
	TagKey        string `json:"tagkey"`
	TagValue      string `json:"tagvalue"`
}
