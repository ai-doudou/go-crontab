package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gorhill/cronexpr"
	"time"
)

const (
	EtcdJobPrefix     = "/cron/job/"
	EtcdKillJobPrefix = "/cron/kill/"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type JobSchedulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type JobEvent struct {
	Job       *Job
	EventType mvccpb.Event_EventType
}

func BuildJobName(job *Job) string {
	return fmt.Sprintf(EtcdJobPrefix+"%s", job.Name)
}
func BuildKillJobName(job *Job) string {
	return fmt.Sprintf(EtcdKillJobPrefix+"%s", job.Name)
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResponse(code int, msg string, data interface{}) []byte {
	resp := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	bs, _ := json.Marshal(resp)
	return bs
}

func UnPackJob(bs []byte) (*Job, error) {
	var job = &Job{}
	err := json.Unmarshal(bs, &job)
	return job, err
}

func BuildJobPlan(event *JobEvent) (*JobSchedulePlan, error) {
	if event.Job == nil {
		return nil, errors.New("job is nil")
	}
	expr, err := cronexpr.Parse(event.Job.CronExpr)
	if err != nil {
		return nil, err
	}

	jobPlan := &JobSchedulePlan{
		Job:      event.Job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return jobPlan, nil
}
