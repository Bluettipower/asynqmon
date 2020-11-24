package main

import (
	"time"

	"github.com/hibiken/asynq"
)

type QueueStateSnapshot struct {
	// Name of the queue.
	Queue string `json:"queue"`
	// Total number of tasks in the queue.
	Size int `json:"size"`
	// Number of tasks in each state.
	Active    int `json:"active"`
	Pending   int `json:"pending"`
	Scheduled int `json:"scheduled"`
	Retry     int `json:"retry"`
	Dead      int `json:"dead"`

	// Total number of tasks processed during the given date.
	// The number includes both succeeded and failed tasks.
	Processed int `json:"processed"`
	// Breakdown of processed tasks.
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	// Paused indicates whether the queue is paused.
	// If true, tasks in the queue will not be processed.
	Paused bool `json:"paused"`
	// Time when this snapshot was taken.
	Timestamp time.Time `json:"timestamp"`
}

func toQueueStateSnapshot(s *asynq.QueueStats) *QueueStateSnapshot {
	return &QueueStateSnapshot{
		Queue:     s.Queue,
		Size:      s.Size,
		Active:    s.Active,
		Pending:   s.Pending,
		Scheduled: s.Scheduled,
		Retry:     s.Retry,
		Dead:      s.Dead,
		Processed: s.Processed,
		Succeeded: s.Processed - s.Failed,
		Failed:    s.Failed,
		Paused:    s.Paused,
		Timestamp: s.Timestamp,
	}
}

type DailyStats struct {
	Queue     string    `json:"queue"`
	Processed int       `json:"processed"`
	Succeeded int       `json:"succeeded"`
	Failed    int       `json:"failed"`
	Date      time.Time `json:"date"`
}

func toDailyStats(s *asynq.DailyStats) *DailyStats {
	return &DailyStats{
		Queue:     s.Queue,
		Processed: s.Processed,
		Succeeded: s.Processed - s.Failed,
		Failed:    s.Failed,
		Date:      s.Date,
	}
}

type BaseTask struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Payload asynq.Payload `json:"payload"`
	Queue   string        `json:"queue"`
}

type ActiveTask struct {
	*BaseTask
}

func toActiveTask(t *asynq.ActiveTask) *ActiveTask {
	base := &BaseTask{
		ID:      t.ID,
		Type:    t.Type,
		Payload: t.Payload,
		Queue:   t.Queue,
	}
	return &ActiveTask{base}
}

func toActiveTasks(in []*asynq.ActiveTask) []*ActiveTask {
	out := make([]*ActiveTask, len(in))
	for i, t := range in {
		out[i] = toActiveTask(t)
	}
	return out
}

type PendingTask struct {
	*BaseTask
}

func toPendingTask(t *asynq.PendingTask) *PendingTask {
	base := &BaseTask{
		ID:      t.ID,
		Type:    t.Type,
		Payload: t.Payload,
		Queue:   t.Queue,
	}
	return &PendingTask{base}
}

func toPendingTasks(in []*asynq.PendingTask) []*PendingTask {
	out := make([]*PendingTask, len(in))
	for i, t := range in {
		out[i] = toPendingTask(t)
	}
	return out
}

type ScheduledTask struct {
	*BaseTask
	NextProcessAt time.Time `json:"next_process_at"`
}

func toScheduledTask(t *asynq.ScheduledTask) *ScheduledTask {
	base := &BaseTask{
		ID:      t.ID,
		Type:    t.Type,
		Payload: t.Payload,
		Queue:   t.Queue,
	}
	return &ScheduledTask{
		BaseTask:      base,
		NextProcessAt: t.NextProcessAt,
	}
}

func toScheduledTasks(in []*asynq.ScheduledTask) []*ScheduledTask {
	out := make([]*ScheduledTask, len(in))
	for i, t := range in {
		out[i] = toScheduledTask(t)
	}
	return out
}

type RetryTask struct {
	*BaseTask
	NextProcessAt time.Time `json:"next_process_at"`
	MaxRetry      int       `json:"max_retry"`
	Retried       int       `json:"retried"`
	ErrorMsg      string    `json:"error_message"`
}

func toRetryTask(t *asynq.RetryTask) *RetryTask {
	base := &BaseTask{
		ID:      t.ID,
		Type:    t.Type,
		Payload: t.Payload,
		Queue:   t.Queue,
	}
	return &RetryTask{
		BaseTask:      base,
		NextProcessAt: t.NextProcessAt,
		MaxRetry:      t.MaxRetry,
		Retried:       t.Retried,
		ErrorMsg:      t.ErrorMsg,
	}
}

func toRetryTasks(in []*asynq.RetryTask) []*RetryTask {
	out := make([]*RetryTask, len(in))
	for i, t := range in {
		out[i] = toRetryTask(t)
	}
	return out
}

type DeadTask struct {
	*BaseTask
	MaxRetry     int       `json:"max_retry"`
	Retried      int       `json:"retried"`
	ErrorMsg     string    `json:"error_message"`
	LastFailedAt time.Time `json:"last_failed_at"`
}

func toDeadTask(t *asynq.DeadTask) *DeadTask {
	base := &BaseTask{
		ID:      t.ID,
		Type:    t.Type,
		Payload: t.Payload,
		Queue:   t.Queue,
	}
	return &DeadTask{
		BaseTask:     base,
		MaxRetry:     t.MaxRetry,
		Retried:      t.Retried,
		ErrorMsg:     t.ErrorMsg,
		LastFailedAt: t.LastFailedAt,
	}
}

func toDeadTasks(in []*asynq.DeadTask) []*DeadTask {
	out := make([]*DeadTask, len(in))
	for i, t := range in {
		out[i] = toDeadTask(t)
	}
	return out
}
