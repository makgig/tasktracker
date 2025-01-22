package task

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxDaysInterval = 400
	daysInWeek      = 7
	monthsInYear    = 12
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if err := ValidateDate(date); err != nil {
		return "", err
	}

	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	baseDate, err := ParseDate(date)
	if err != nil {
		return "", err
	}

	rule, err := ParseRepeatRule(repeat)
	if err != nil {
		return "", err
	}

	nextDate, err := rule.calculateNextDate(now, baseDate)
	if err != nil {
		return "", err
	}

	return FormatDate(nextDate), nil
}

type RepeatRule struct {
	Type      string
	Days      int
	DaysWeek  []int
	MonthDays []int
	Months    []int
}

func ParseRepeatRule(rule string) (*RepeatRule, error) {
	parts := strings.Fields(rule)
	if len(parts) < 1 {
		return nil, fmt.Errorf("некорректный формат правила")
	}

	r := &RepeatRule{Type: parts[0]}

	switch r.Type {
	case "d":
		if err := r.parseDayRule(parts); err != nil {
			return nil, err
		}
	case "y":
		if len(parts) != 1 {
			return nil, fmt.Errorf("правило y не требует дополнительных параметров")
		}
	case "w":
		if err := r.parseWeekRule(parts); err != nil {
			return nil, err
		}
	case "m":
		if err := r.parseMonthRule(parts); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("неподдерживаемый тип правила: %s", r.Type)
	}

	return r, nil
}

func (r *RepeatRule) parseDayRule(parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("правило d требует указания количества дней")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("некорректное количество дней: %w", err)
	}

	if days <= 0 || days > maxDaysInterval {
		return fmt.Errorf("количество дней должно быть в диапазоне от 1 до %d", maxDaysInterval)
	}

	r.Days = days
	return nil
}

func (r *RepeatRule) parseWeekRule(parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("правило w требует указания дней недели")
	}

	daysStr := strings.Split(parts[1], ",")
	r.DaysWeek = make([]int, 0, len(daysStr))

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return fmt.Errorf("некорректный день недели: %s", dayStr)
		}
		if day < 1 || day > daysInWeek {
			return fmt.Errorf("день недели должен быть от 1 до 7: %d", day)
		}
		r.DaysWeek = append(r.DaysWeek, day)
	}

	return nil
}

func (r *RepeatRule) parseMonthRule(parts []string) error {
	if len(parts) < 2 || len(parts) > 3 {
		return fmt.Errorf("правило m требует указания дней и опционально месяцев")
	}

	daysStr := strings.Split(parts[1], ",")
	r.MonthDays = make([]int, 0, len(daysStr))

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return fmt.Errorf("некорректный день месяца: %s", dayStr)
		}

		// Проверяем допустимые значения
		if day > 0 && day > 31 {
			return fmt.Errorf("день месяца должен быть от 1 до 31 или -1, -2: %d", day)
		}
		if day < 0 && day < -2 {
			return fmt.Errorf("отрицательный день месяца может быть только -1 или -2: %d", day)
		}

		r.MonthDays = append(r.MonthDays, day)
	}

	if len(parts) == 3 {
		monthsStr := strings.Split(parts[2], ",")
		r.Months = make([]int, 0, len(monthsStr))

		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				return fmt.Errorf("некорректный месяц: %s", monthStr)
			}
			if month < 1 || month > monthsInYear {
				return fmt.Errorf("месяц должен быть от 1 до 12: %d", month)
			}
			r.Months = append(r.Months, month)
		}
	}

	return nil
}

func (r *RepeatRule) calculateNextDate(now, base time.Time) (time.Time, error) {
	switch r.Type {
	case "d":
		return r.calculateDayRule(now, base)
	case "y":
		return r.calculateYearRule(now, base)
	case "w":
		return r.calculateWeekRule(now, base)
	case "m":
		return r.calculateMonthRule(now, base)
	default:
		return time.Time{}, fmt.Errorf("неподдерживаемый тип правила: %s", r.Type)
	}
}

func (r *RepeatRule) calculateDayRule(now, base time.Time) (time.Time, error) {
	days := r.Days
	next := base

	// Если базовая дата раньше текущей
	if next.Before(now) {
		// Точный расчет следующей даты с учетом интервала
		for next.Before(now) {
			next = next.AddDate(0, 0, days)
		}
	} else {
		// Если базовая дата еще не наступила
		next = next.AddDate(0, 0, days)
	}

	return next, nil
}

func (r *RepeatRule) calculateYearRule(now, base time.Time) (time.Time, error) {
	next := base

	if next.After(now) || next.Equal(now) {
		next = next.AddDate(1, 0, 0)
	} else {
		for !next.After(now) {
			// Сохраним изначальные месяц и день до добавления года
			month := next.Month()
			day := next.Day()

			next = next.AddDate(1, 0, 0)

			// Если была дата 29 февраля и после добавления года оказались в марте,
			// значит следующий год не високосный - установим дату на 1 марта
			if month == time.February && day == 29 && next.Month() == time.March {
				next = time.Date(next.Year(), time.March, 1, 0, 0, 0, 0, next.Location())
			}
		}
	}
	return next, nil
}

func (r *RepeatRule) calculateWeekRule(now, base time.Time) (time.Time, error) {
	next := base
	if next.Before(now) {
		next = now
	}

	for i := 0; i <= daysInWeek; i++ {
		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		for _, day := range r.DaysWeek {
			if weekday == day && next.After(now) {
				return next, nil
			}
		}
		next = next.AddDate(0, 0, 1)
	}

	return next, nil
}

func (r *RepeatRule) calculateMonthRule(now, base time.Time) (time.Time, error) {
	months := make(map[int]bool)
	if len(r.Months) > 0 {
		for _, month := range r.Months {
			months[month] = true
		}
	}

	next := base
	if next.Before(now) {
		next = now
	}

	// Поиск в течение 2 лет
	for i := 0; i < 24*31; i++ {
		currentMonth := int(next.Month())

		// Проверка на допустимые месяцы
		if len(months) > 0 && !months[currentMonth] {
			next = time.Date(next.Year(), next.Month()+1, 1, 0, 0, 0, 0, next.Location())
			continue
		}

		// Определяем последний день текущего месяца
		lastDay := time.Date(next.Year(), next.Month()+1, 0, 0, 0, 0, 0, next.Location()).Day()

		var candidates []time.Time

		for _, day := range r.MonthDays {
			var targetDay int

			if day < 0 {
				targetDay = lastDay + day + 1 // -1 -> последний день, -2 -> предпоследний
			} else {
				targetDay = day
			}

			if targetDay < 1 || targetDay > lastDay {
				continue
			}

			candidateDate := time.Date(next.Year(), next.Month(), targetDay, 0, 0, 0, 0, next.Location())
			if candidateDate.After(now) {
				candidates = append(candidates, candidateDate)
			}
		}

		if len(candidates) > 0 {
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].Before(candidates[j])
			})
			return candidates[0], nil
		}

		next = next.AddDate(0, 0, 1)
	}

	return time.Time{}, fmt.Errorf("не удалось найти подходящую дату")
}
