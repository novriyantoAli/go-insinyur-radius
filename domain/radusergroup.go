package domain

import "context"

// Radusergroup ...
type Radusergroup struct {
	Username  *string `json:"username"`
	Groupname *string `json:"groupname"`
	Priority  *int64  `json:"priority"`
}

//Profile ...
type Profile struct {
	ProfileName          string `json:"profile_name"`
	PrefixName           string `json:"prefix_name"`
	Priority             int64  `json:"priority"`
	PoolName             string `json:"pool_name"`
	UseLimitSession      bool   `json:"use_limit_session"`
	LimitSession         int    `json:"limit_session"`
	UseLimitSpeed        bool   `json:"use_limit_speed"`
	LimitCIRUpload       int64  `json:"limit_cir_upload"`
	LimitCIRUploadUnit   string `json:"limit_cir_upload_unit"`
	LimitMIRUpload       int64  `json:"limit_mir_upload"`
	LimitMIRUploadUnit   string `json:"limit_mir_upload_unit"`
	LimitCIRDownload     int64  `json:"limit_cir_download"`
	LimitCIRDownloadUnit string `json:"limit_cir_download_unit"`
	LimitMIRDownload     int64  `json:"limit_mir_download"`
	LimitMIRDownloadUnit string `json:"limit_mir_download_unit"`
}

// RadusergroupRepository ...
type RadusergroupRepository interface {
	Fetch(ctx context.Context) (res []Radusergroup, err error)
	Get(ctx context.Context, username string) (res Radusergroup, err error)
	Insert(ctx context.Context, rug Radusergroup) (err error)
	SaveProfile(ctx context.Context, profile Profile) (err error)
	Delete(ctx context.Context, username string) (err error)
}

// RadusergroupUsecase ...
type RadusergroupUsecase interface {
	Fetch(ctx context.Context) (res []Radusergroup, err error)
	Get(ctx context.Context, username string) (res Radusergroup, err error)
	LoadProfile(ctx context.Context, profile string) (res Profile, err error)
	Save(ctx context.Context, m *Radusergroup) (err error)
	SaveProfile(ctx context.Context, m *Profile) (result Radusergroup, err error)
	Delete(ctx context.Context, username string) (err error)
}
