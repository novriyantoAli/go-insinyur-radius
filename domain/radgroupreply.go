package domain

import "context"

// Radgroupreply ...
type Radgroupreply struct {
	ID        *int64  `json:"id"`
	Groupname *string `json:"groupname"`
	Attribute *string `json:"attribute"`
	OP        *string `json:"op"`
	Value     *string `json:"value"`
}

// RadgroupreplyRepository ...
type RadgroupreplyRepository interface {
	Find(ctx context.Context, spec Radgroupreply) (res Radgroupreply, err error)
	DeleteWithGroupname(ctx context.Context, name string) (err error)
}

// RadgroupreplyUsecase ...
type RadgroupreplyUsecase interface {
	Find(ctx context.Context, spec Radgroupreply) (res Radgroupreply, err error)
}
