package g

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
	"sync"
	"log"
	"errors"
)

type EventDto struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Counter  string `json:"counter"`

	Func       string `json:"func"`
	LeftValue  string `json:"leftValue"`
	Operator   string `json:"operator"`
	RightValue string `json:"rightValue"`

	Note string `json:"note"`

	MaxStep     int `json:"maxStep"`
	CurrentStep int `json:"currentStep"`
	Priority    int `json:"priority"`

	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`

	ExpressionId int `json:"expressionId"`
	StrategyId   int `json:"strategyId"`
	TemplateId   int `json:"templateId"`

	// abandoned in domeos
	Link string `json:"link"`
}

type SafeEvents struct {
	sync.RWMutex
	// id -> EventDto
	M map[string]*EventDto
}

type OrderedEvents []*EventDto

func (this OrderedEvents) Len() int {
	return len(this)
}
func (this OrderedEvents) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this OrderedEvents) Less(i, j int) bool {
	return this[i].Timestamp < this[j].Timestamp
}

var Events = &SafeEvents{M: make(map[string]*EventDto)}

func (this *SafeEvents) Delete(id string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, id)

	// add mysql
	sql := fmt.Sprintf("delete from alarm_event_info_draft where id = '%s'", id)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail: ", err)
	}
}

func (this *SafeEvents) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *SafeEvents) Clone() map[string]*EventDto {
	m := make(map[string]*EventDto)
	this.RLock()
	defer this.RUnlock()
	for key, val := range this.M {
		m[key] = val
	}
	return m
}

func (this *SafeEvents) Put(event *model.Event) error {
	if event.Status == "OK" {
		this.Delete(event.Id)
		return nil
	}

	// verify the strategy this event binds to still exists
	var strategyCount int
	check := fmt.Sprintf(
		"select count(*) from portal.strategy where id = %d",
		event.StrategyId(),
	)
	err := DB.QueryRow(check).Scan(&strategyCount)
	if err != nil {
		return errors.New("MySQL error in select strategy")
	}
	if strategyCount == 0 {
		return errors.New("non-existed strategy, drop")
	}

	dto := &EventDto{}
	dto.Id = event.Id
	dto.Endpoint = event.Endpoint
	dto.Metric = event.Metric()
	dto.Counter = event.Counter()
	dto.Func = event.Func()
	dto.LeftValue = utils.ReadableFloat(event.LeftValue)
	dto.Operator = event.Operator()
	dto.RightValue = utils.ReadableFloat(event.RightValue())
	dto.Note = event.Note()

	dto.MaxStep = event.MaxStep()
	dto.CurrentStep = event.CurrentStep
	dto.Priority = event.Priority()

	dto.Status = event.Status
	dto.Timestamp = event.EventTime

	dto.ExpressionId = event.ExpressionId()
	dto.StrategyId = event.StrategyId()
	dto.TemplateId = event.TplId()

	// abandoned in domeos
	dto.Link = ""

	this.Lock()
	defer this.Unlock()
	this.M[dto.Id] = dto

	// add mysql
	sql := fmt.Sprintf(
		"insert into alarm_event_info_draft(id, endpoint, metric, counter, func, left_value, operator, right_value, note, max_step, current_step, priority, " +
		"status, timestamp, expression_id, strategy_id, template_id) values ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', " +
		"%d, %d, %d, '%s', %d, %d, %d, %d) on duplicate key update endpoint='%s', metric='%s', counter='%s', func='%s', " +
		"left_value='%s', operator='%s', right_value='%s', note='%s', max_step=%d, current_step=%d, priority=%d, status='%s', timestamp=%d, " +
		"expression_id=%d, strategy_id=%d, template_id=%d",
		dto.Id,
		dto.Endpoint,
		dto.Metric,
		dto.Counter,
		dto.Func,
		dto.LeftValue,
		dto.Operator,
		dto.RightValue,
		dto.Note,
		dto.MaxStep,
		dto.CurrentStep,
		dto.Priority,
		dto.Status,
		dto.Timestamp,
		dto.ExpressionId,
		dto.StrategyId,
		dto.TemplateId,
		dto.Endpoint,
		dto.Metric,
		dto.Counter,
		dto.Func,
		dto.LeftValue,
		dto.Operator,
		dto.RightValue,
		dto.Note,
		dto.MaxStep,
		dto.CurrentStep,
		dto.Priority,
		dto.Status,
		dto.Timestamp,
		dto.ExpressionId,
		dto.StrategyId,
		dto.TemplateId,
	)

	_, err = DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail: ", err)
	}

	return nil

}

/*
func Link(event *model.Event) string {
	tplId := event.TplId()
	if tplId != 0 {
		// a link for template viewing
		return fmt.Sprintf("%s/api/alarm/template/view/%d", Config().Api.DomeOS, tplId)
	}

	eid := event.ExpressionId()
	if eid != 0 {
		// a link for expression viewing
		return fmt.Sprintf("%s/api/alarm/expression/view/%d", Config().Api.DomeOS, eid)
	}

	return ""
}
*/
