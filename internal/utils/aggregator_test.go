package utils

import (
    "testing"
    "time"
)

// helper to make APIUsageEntry with created at given time
func mkEntry(t time.Time) APIUsageEntry {
    return APIUsageEntry{
        ID:      "id",
        Model:   "gpt-4o",
        Created: t.Unix(),
        Usage:   APIUsage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
        Cost:    0.01,
    }
}

func TestFilterEntriesByDateRange_Inclusive(t *testing.T) {
    start := time.Date(2025, 9, 1, 0, 0, 0, 0, time.Local)
    end := time.Date(2025, 9, 1, 23, 59, 59, 0, time.Local)

    inside := mkEntry(start.Add(12 * time.Hour))
    atStart := mkEntry(start)
    atEnd := mkEntry(end)
    before := mkEntry(start.Add(-time.Second))
    after := mkEntry(end.Add(time.Second))

    entries := []APIUsageEntry{inside, atStart, atEnd, before, after}

    filtered := filterEntriesByDateRange(entries, start, end)
    if len(filtered) != 3 {
        t.Fatalf("expected 3 entries, got %d", len(filtered))
    }
}

func TestAggregateDailyUsage_LocalDateGrouping(t *testing.T) {
    // Choose a fixed local date to avoid flakiness
    base := time.Date(2025, 9, 2, 10, 0, 0, 0, time.Local)
    e1 := mkEntry(base)
    e2 := mkEntry(base.Add(2 * time.Hour))

    daily := AggregateDailyUsage([]APIUsageEntry{e1, e2})
    if len(daily) != 1 {
        t.Fatalf("expected 1 day, got %d", len(daily))
    }
    wantDate := base.Format("2006-01-02")
    if daily[0].Date != wantDate {
        t.Fatalf("expected date %s, got %s", wantDate, daily[0].Date)
    }
}

func TestAggregateMonthlyUsage_IncludesBoundaryDays(t *testing.T) {
    // Build entries at first and last second of a month in local time
    monthStart := time.Date(2025, 8, 1, 0, 0, 0, 0, time.Local)
    monthEnd := time.Date(2025, 8, 31, 23, 59, 59, 0, time.Local)

    eStart := mkEntry(monthStart)
    eEnd := mkEntry(monthEnd)
    eMid := mkEntry(monthStart.Add(15 * 24 * time.Hour))

    monthly := AggregateMonthlyUsage([]APIUsageEntry{eStart, eMid, eEnd})
    if len(monthly) == 0 {
        t.Fatalf("expected at least one month aggregate")
    }
    found := false
    for _, m := range monthly {
        if m.Month == "2025-08" {
            // Boundaries should be included; total requests = 3
            if m.RequestCount != 3 {
                t.Fatalf("expected 3 requests, got %d", m.RequestCount)
            }
            // Daily breakdown should include entries on boundaries
            if len(m.DailyBreakdown) == 0 {
                t.Fatalf("expected non-empty daily breakdown")
            }
            found = true
        }
    }
    if !found {
        t.Fatalf("expected month 2025-08 in results")
    }
}

