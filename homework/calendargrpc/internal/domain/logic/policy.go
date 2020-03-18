package logic

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
)

type RetentionPolicy struct {
	d time.Duration
}

func NewRetentionPolicy(d time.Duration) interfaces.RetentionPolicy {
	return &RetentionPolicy{d: d}
}

func (r RetentionPolicy) GetDuration() time.Duration {
	return r.d
}
