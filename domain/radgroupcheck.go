package domain

import "context"

// Radgroupcheck ...
type Radgroupcheck struct {
	ID        *int64  `json:"id"`
	Groupname *string `json:"groupname"`
	Attribute *string `json:"attribute"`
	OP        *string `json:"op"`
	Value     *string `json:"value"`
}

// RadgroupcheckRepository ...
type RadgroupcheckRepository interface {
	Find(ctx context.Context, spec Radgroupcheck)(res Radgroupcheck, err error)
	DeleteWithGroupname(ctx context.Context, name string) (err error)
}

// RadgroupcheckUsecase ...
type RadgroupcheckUsecase interface {
	Find(ctx context.Context, spec Radgroupcheck)(res Radgroupcheck, err error)
}
