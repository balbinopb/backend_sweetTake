// internal/analytics/charts.go

package controllers

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"sweetake/models"
)

// =========================
// Helpers
// =========================

// sameDay returns true if both times are on the same calendar day (local).
func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// floorToStartOfDay returns t at 00:00:00.000.
func floorToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// ceilToEndOfDay returns t at 23:59:59.999999999.
func ceilToEndOfDay(t time.Time) time.Time {
	start := floorToStartOfDay(t)
	return start.Add(24*time.Hour - time.Nanosecond)
}

// isoWeekKey builds "YYYY-Www" from a time.Time using ISOWeek.
func isoWeekKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%04d-W%02d", year, week)
}

// monthKey builds "YYYY-MM" from a time.Time.
func monthKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month()))
}

// =========================
// DAILY CHART
// =========================

// BuildDailyConsumptionChart aggregates sugar grams per day in [start, end] for a user,
// fills missing calendar days with zero, sorts ascending, and returns JSON:
// [
//
//	{"day":"YYYY-MM-DD","sugar":36.5},
//	{"day":"YYYY-MM-DD","sugar":28.0},
//	...
//
// ]
func BuildDailyConsumptionChart(cons []models.Consumption, userID uint, start, end time.Time) ([]byte, error) {
	start = floorToStartOfDay(start)
	end = ceilToEndOfDay(end)

	// Pre-create all day buckets to fill gaps with zeros.
	daily := make(map[string]float64)
	orderedDays := make([]string, 0)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		daily[key] = 0
		orderedDays = append(orderedDays, key)
	}

	// Aggregate consumptions in the range.
	for _, c := range cons {
		if c.UserID != userID {
			continue
		}
		if c.DateTime.Before(start) || c.DateTime.After(end) {
			continue
		}
		key := c.DateTime.Format("2006-01-02")
		daily[key] += safeSugar(c)
	}

	// Build sorted JSON array.
	out := make([]map[string]interface{}, 0, len(orderedDays))
	for _, key := range orderedDays {
		out = append(out, map[string]interface{}{
			"day":   key,
			"sugar": daily[key],
		})
	}

	return json.Marshal(out)
}

// =========================
// WEEKLY CHART
// =========================

// BuildWeeklyConsumptionChart aggregates sugar grams per ISO week in [start, end] for a user,
// includes all calendar weeks that intersect the period (gap-fill zeros), sorts ascending by (year, week),
// and returns JSON:
// [
//
//	{"week":"YYYY-Www","total_sugar":210.0,"avg_per_day":30.0},
//	...
//
// ]
func BuildWeeklyConsumptionChart(cons []models.Consumption, userID uint, start, end time.Time) ([]byte, error) {
	start = floorToStartOfDay(start)
	end = ceilToEndOfDay(end)

	// Build calendar week keys covering the range.
	// Find the Monday of the first week (ISO week starts Monday).
	// Go doesn't have direct Monday-of-week; compute by stepping back to Monday.
	weekday := int(start.Weekday())
	if weekday == 0 { // Sunday -> treat as 7
		weekday = 7
	}
	firstMonday := floorToStartOfDay(start.AddDate(0, 0, -(weekday - 1)))

	weeks := make(map[string]float64)
	orderedWeeks := make([]string, 0)

	for wStart := firstMonday; !wStart.After(end); wStart = wStart.AddDate(0, 0, 7) {
		key := isoWeekKey(wStart)
		if _, exists := weeks[key]; !exists {
			weeks[key] = 0
			orderedWeeks = append(orderedWeeks, key)
		}
	}

	// Aggregate by ISO week.
	for _, c := range cons {
		if c.UserID != userID {
			continue
		}
		if c.DateTime.Before(start) || c.DateTime.After(end) {
			continue
		}
		key := isoWeekKey(c.DateTime)
		weeks[key] += safeSugar(c)
	}

	// Sort orderedWeeks by (year, week).
	sort.Slice(orderedWeeks, func(i, j int) bool {
		var yi, wi, yj, wj int
		fmt.Sscanf(orderedWeeks[i], "%d-W%d", &yi, &wi)
		fmt.Sscanf(orderedWeeks[j], "%d-W%d", &yj, &wj)
		if yi == yj {
			return wi < wj
		}
		return yi < yj
	})

	// Build JSON with avg per day (7-day weeks).
	out := make([]map[string]interface{}, 0, len(orderedWeeks))
	for _, wk := range orderedWeeks {
		total := weeks[wk]
		out = append(out, map[string]interface{}{
			"week":        wk,
			"total_sugar": total,
			"avg_per_day": total / 7.0,
		})
	}

	return json.Marshal(out)
}

// =========================
// MONTHLY CHART
// =========================

// BuildMonthlyConsumptionChart aggregates sugar grams per calendar month in [start, end] for a user,
// includes all months that intersect the period (gap-fill zeros), sorts ascending (YYYY-MM),
// and returns JSON:
// [
//
//	{"month":"YYYY-MM","total_sugar":900.0,"avg_per_day":30.0},
//	...
//
// ]
func BuildMonthlyConsumptionChart(cons []models.Consumption, userID uint, start, end time.Time) ([]byte, error) {
	start = floorToStartOfDay(start)
	end = ceilToEndOfDay(end)

	// Build month keys from start's first day of month to end's last day of month.
	firstMonth := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	lastMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())

	months := make(map[string]float64)
	orderedMonths := make([]string, 0)

	for m := firstMonth; !m.After(lastMonth); m = m.AddDate(0, 1, 0) {
		key := monthKey(m)
		months[key] = 0
		orderedMonths = append(orderedMonths, key)
	}

	// Aggregate by month.
	for _, c := range cons {
		if c.UserID != userID {
			continue
		}
		if c.DateTime.Before(start) || c.DateTime.After(end) {
			continue
		}
		key := monthKey(c.DateTime)
		months[key] += safeSugar(c)
	}

	// Sort month keys "YYYY-MM" lexicographically (works because zero-padded).
	sort.Strings(orderedMonths)

	// avg_per_day uses the actual number of calendar days in the month.
	out := make([]map[string]interface{}, 0, len(orderedMonths))
	for _, mk := range orderedMonths {
		total := months[mk]

		// Parse month key to determine days in month.
		var year, month int
		fmt.Sscanf(mk, "%d-%d", &year, &month)
		first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, start.Location())
		next := first.AddDate(0, 1, 0)
		daysInMonth := int(next.Sub(first).Hours()/24 + 0.5) // handle DST oddities

		out = append(out, map[string]interface{}{
			"month":       mk,
			"total_sugar": total,
			"avg_per_day": total / float64(daysInMonth),
		})
	}

	return json.Marshal(out)
}
