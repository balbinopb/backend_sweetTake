//internal/analytics/metrics.go

package analytics

import (
	"math"
	"time"

	"sweetake/models"
)

// -------------------- Configurable thresholds --------------------

var (
	DefaultLowSugarGoalGrams         = 25.0  // WHO recommended free sugar limit (grams/day)
	DefaultUpperSugarLimit           = 50.0  // Upper limit (grams/day)
	HighSugarPerEntryThreshold       = 15.0  // grams -> considered high per entry
	DiabeticSpikeFactor              = 0.5   // expected mg/dL spike per gram sugar (simple model)
	SpikeHighThreshold               = 20.0  // mg/dL expected spike considered high
	BloodSugarSDWarning              = 20.0  // mg/dL SD threshold for unstable
	WeightLossSugarCaloriesThreshold = 120.0 // kcal from sugar (e.g., 30g => 30*4=120)
)

// -------------------- Helpers --------------------

// Returns sugar grams; 0 if SugarData is nil.
func safeSugar(c models.Consumption) float64 {
	if c.SugarData == nil {
		return 0
	}
	return *c.SugarData
}

// Combine MeasureDate + MeasureTime into a single timestamp.
// If MeasureTime is zero, fall back to MeasureDate.
func measurementTimestamp(m models.BloodSugarMetric) time.Time {
	// If MeasureTime wasn't set, use MeasureDate directly.
	if m.MeasureTime.IsZero() {
		return m.MeasureDate
	}
	return time.Date(
		m.MeasureDate.Year(), m.MeasureDate.Month(), m.MeasureDate.Day(),
		m.MeasureTime.Hour(), m.MeasureTime.Minute(), m.MeasureTime.Second(),
		m.MeasureTime.Nanosecond(), m.MeasureDate.Location(),
	)
}

// -------------------- Utility Calculation Functions --------------------

// Sum sugar grams in given consumptions (optionally filter by day)
func SumSugar(cons []models.Consumption) float64 {
	var s float64
	for _, c := range cons {
		s += safeSugar(c)
	}
	return s
}

// Get consumptions for a user within a date range inclusive of start and end
func FilterConsumptionByUserAndRange(cons []models.Consumption, userID uint, start, end time.Time) []models.Consumption {
	var out []models.Consumption
	for _, c := range cons {
		if c.UserID != userID {
			continue
		}
		if !c.DateTime.Before(start) && !c.DateTime.After(end) {
			out = append(out, c)
		}
	}
	return out
}

// Daily total for a specific date (local)
func DailyTotal(cons []models.Consumption, userID uint, date time.Time) float64 {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return SumSugar(FilterConsumptionByUserAndRange(cons, userID, start, end))
}

// Weekly total for the last 7 days ending at 'endDate' (inclusive)
func WeeklyTotal(cons []models.Consumption, userID uint, endDate time.Time) float64 {
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	start := end.AddDate(0, 0, -6) // 7-day window
	return SumSugar(FilterConsumptionByUserAndRange(cons, userID, start, end))
}

// Monthly total for last 30 days ending at endDate
func MonthlyTotal(cons []models.Consumption, userID uint, endDate time.Time) float64 {
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	start := end.AddDate(0, 0, -29)
	return SumSugar(FilterConsumptionByUserAndRange(cons, userID, start, end))
}

// Weekly average (daily average over 7 days)
func WeeklyAverage(cons []models.Consumption, userID uint, endDate time.Time) float64 {
	return WeeklyTotal(cons, userID, endDate) / 7.0
}

// Standard deviation for blood sugar values in a date range
func StdDevBloodSugar(metrics []models.BloodSugarMetric, userID uint, start, end time.Time) float64 {
	var vals []float64
	for _, m := range metrics {
		if m.UserID != userID {
			continue
		}
		t := measurementTimestamp(m)
		if !t.Before(start) && !t.After(end) {
			vals = append(vals, m.Value)
		}
	}
	if len(vals) == 0 {
		return 0
	}
	// mean
	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	// variance
	var vsum float64
	for _, v := range vals {
		diff := v - mean
		vsum += diff * diff
	}
	variance := vsum / float64(len(vals))
	return math.Sqrt(variance)
}

// Expected spike estimate from a consumption entry (simple model)
func ExpectedSpike(sugarGrams float64) float64 {
	// simplistic linear model: mg/dL increase per gram sugar
	return sugarGrams * DiabeticSpikeFactor
}

// Maximum expected spike in a window
func MaxExpectedSpike(cons []models.Consumption, userID uint, start, end time.Time) float64 {
	s := FilterConsumptionByUserAndRange(cons, userID, start, end)
	var max float64
	for _, c := range s {
		spike := ExpectedSpike(safeSugar(c))
		if spike > max {
			max = spike
		}
	}
	return max
}

// Sugar calories estimate
func SugarCalories(grams float64) float64 {
	return grams * 4.0
}

// Group daily sums for a date range (useful for charting)
func DailySums(cons []models.Consumption, userID uint, start, end time.Time) map[string]float64 {
	out := map[string]float64{}
	// pre-create day keys
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		out[key] = 0
	}
	for _, c := range cons {
		if c.UserID != userID {
			continue
		}
		if c.DateTime.Before(start) || c.DateTime.After(end) {
			continue
		}
		key := c.DateTime.Format("2006-01-02")
		out[key] += safeSugar(c)
	}
	return out
}
