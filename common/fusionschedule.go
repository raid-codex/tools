package common

import (
	"fmt"
	"sort"
	"time"
)

type FusionSchedule struct {
	DateStart string                `json:"-"`
	DateEnd   string                `json:"-"`
	Raw       []*FusionScheduleItem `json:"raw"`
	Daily     map[string][]int      `json:"daily"`
}

type FusionScheduleItem struct {
	Index         int      `json:"index"`
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	DateStart     string   `json:"date_start"`
	DateEnd       string   `json:"date_end"`
	ChampionSlugs []string `json:"champion_slugs"`
}

func (fs *FusionSchedule) Sanitize() error {
	if fs.Raw == nil {
		return nil
	}
	fs.Daily = make(map[string][]int)
	current, err := time.Parse("2006-01-02", fs.DateStart)
	if err != nil {
		return err
	}
	end, err := time.Parse("2006-01-02", fs.DateEnd)
	if err != nil {
		return err
	}
	for current.Before(end) || current == end {
		fs.Daily[current.Format("2006-01-02")] = make([]int, 0)
		current = current.AddDate(0, 0, 1)
	}
	for idx, fsi := range fs.Raw {
		fsi.Index = idx + 1
		if err := fsi.Sanitize(); err != nil {
			return err
		}
		fsiCurrent, err := time.Parse("2006-01-02", fsi.DateStart)
		if err != nil {
			return err
		}
		fsiEnd, err := time.Parse("2006-01-02", fsi.DateEnd)
		if err != nil {
			return err
		}
		for fsiCurrent.Before(fsiEnd) || fsiCurrent == fsiEnd {
			fs.Daily[fsiCurrent.Format("2006-01-02")] = append(fs.Daily[fsiCurrent.Format("2006-01-02")], fsi.Index)
			fsiCurrent = fsiCurrent.AddDate(0, 0, 1)
		}
	}
	return nil
}

func (fsi *FusionScheduleItem) Sanitize() error {
	for _, slug := range fsi.ChampionSlugs {
		if champions, err := GetChampions(FilterChampionSlug(slug)); err != nil {
			return err
		} else if len(champions) != 1 {
			return fmt.Errorf("found %d champions for slug %s", len(champions), slug)
		}
	}
	sort.Strings(fsi.ChampionSlugs)
	return nil
}
