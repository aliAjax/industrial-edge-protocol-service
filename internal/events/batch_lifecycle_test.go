package events_test

import (
	"sync"
	"testing"
	"time"

	"industrial-edge-protocol/internal/alerts"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/events"
	"industrial-edge-protocol/internal/rules"
)

func TestAlarmDispatchWaitsForAllHandlers(t *testing.T) {
	bus := events.NewBus()
	bus.Subscribe("alarm", func(events.Event) {})
	start := make(chan struct{})
	counts := make(chan int, 2)
	var callers sync.WaitGroup
	callers.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer callers.Done()
			<-start
			counts <- bus.PublishBatch([]events.Event{{Type: "alarm"}, {Type: "alarm"}})
		}()
	}
	close(start)
	callers.Wait()
	close(counts)
	for count := range counts {
		if count != 2 {
			t.Fatalf("batch returned before handlers completed: %d", count)
		}
	}
}

func TestAlarmDispatchFailureDoesNotHang(t *testing.T) {
	manager := alerts.New()
	start := make(chan struct{})
	counts := make(chan int, 2)
	var callers sync.WaitGroup
	callers.Add(2)
	for worker := 0; worker < 2; worker++ {
		go func(offset int) {
			defer callers.Done()
			<-start
			counts <- manager.OpenBatch([]domain.Alarm{{ID: domain.ID("a" + string(rune('0'+offset)))}, {ID: domain.ID("b" + string(rune('0'+offset)))}})
		}(worker)
	}
	close(start)
	callers.Wait()
	close(counts)
	for count := range counts {
		if count != 2 {
			t.Fatalf("alarm batch completed early: %d", count)
		}
	}
}

func TestRuleEvaluationPublishesEveryAlarm(t *testing.T) {
	engine := rules.NewEngine()
	rule := domain.Rule{ID: "rule-1", Expression: "value > 5", Window: 1}
	start := make(chan struct{})
	counts := make(chan int, 2)
	var callers sync.WaitGroup
	callers.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer callers.Done()
			<-start
			counts <- engine.EvaluateBatch(rule, []domain.Reading{{PointID: "p1", Value: 8}, {PointID: "p1", Value: 9}})
		}()
	}
	close(start)
	callers.Wait()
	close(counts)
	for count := range counts {
		if count != 2 {
			t.Fatalf("rule batch lost results: %d", count)
		}
	}
}

func TestWindowSnapshotSurvivesAsyncDispatch(t *testing.T) {
	window := rules.NewWindow(time.Minute)
	start := make(chan struct{})
	counts := make(chan int, 2)
	var callers sync.WaitGroup
	callers.Add(2)
	for i := 0; i < 2; i++ {
		go func(offset int) {
			defer callers.Done()
			<-start
			at := time.Unix(int64(100+offset), 0)
			counts <- window.AddBatch([]domain.Reading{{PointID: "p1", ObservedAt: at}, {PointID: "p1", ObservedAt: at.Add(time.Second)}})
		}(i)
	}
	close(start)
	callers.Wait()
	close(counts)
	for count := range counts {
		if count != 2 {
			t.Fatalf("window batch returned before writes: %d", count)
		}
	}
}
