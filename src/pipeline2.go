package model

import (
	"time"

	"github.com/zhu/qvm/server/lib/params"
	"gopkg.in/mgo.v2/bson"
)

type IMeter interface {
	QueryMeter(m MeterQuery) (MeterResult, error)
	QueryMeterItems(m MeterQuery) ([]*MeterItems, error)
	SupportSpec(spec string) bool
}

type MeterResult interface {
	ToMeterResult() *params.Meter
}

type MeterItems struct {
	Name string `bson:"_id"`
}

type MeterUser struct {
	Uid uint32 `bson:"_id"`
}

type CountMeter struct {
	Id        interface{} `bson:"_id"`
	StartTime time.Time   `bson:"-"`
	Count     int64       `bson:"count"`
}

func (cm *CountMeter) ToMeterResult() (res *params.Meter) {
	res = &params.Meter{
		Time: cm.StartTime.Format(params.MeterTimeFormat),
		Vals: map[string]int64{
			"count": cm.Count,
		},
	}
	return
}

type MeterQuery struct {
	StartTime time.Time
	EndTime   time.Time
	Query     bson.M
}

func NewMeterQuery() MeterQuery {
	return MeterQuery{
		Query: make(bson.M),
	}
}

func (m MeterQuery) Start(t time.Time) MeterQuery {
	m.StartTime = t
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["start_time"] = bson.M{
		"$gte": t,
	}
	return m
}

func (m MeterQuery) End(t time.Time) MeterQuery {
	m.EndTime = t
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["end_time"] = bson.M{
		"$lte": t,
	}
	return m
}

func (m MeterQuery) Uid(uid uint32) MeterQuery {
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["uid"] = uid
	return m
}

func (m MeterQuery) Region(region string) MeterQuery {
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["region_id"] = region
	return m
}

func (m MeterQuery) Spec(spec string) MeterQuery {
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["spec"] = spec
	return m
}
