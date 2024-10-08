package monday

import "testing"

var (
	testingLayoutsData = []struct {
		layout    string
		matches   []string
		unmatches []string
		locales   []Locale
	}{
		{
			layout: "Mon, 2 Jan 2006 15:4:5 -0700",
			matches: []string{
				"Нд, 22 Гру 2019 00:35:44 +0200",
				"Вт, 2 Янв 1997 12:01:44 -0330",
				"Tue, 12 Feb 2014 22:03:40 -0300",
				"Fri, 12 Feb 2016 16:01:00 +0000",
			},
			unmatches: []string{
				"Mon, 2 Jan 2006 15:4:5 -0700 ",
				"Mon, 2  Jan 2006 15:4:5 -0700",
				"Mon, 2 Jan 2006 15:4:5",
				"Mon, 2 Jan 2006 15:4:5 4 -0700",
				"Mon, 2 Jan 2006 15:4:5 4 +0200",
			},
			locales: []Locale{
				LocaleUkUA,
				LocaleRuRU,
				LocaleEnUS,
				LocaleEnUS,
			},
		},
		{
			layout: "2006-01-02T15:04:05-07:00",
			matches: []string{
				"2019-12-22T00:35:44+02:00",
				"2016-02-02T00:00:00+03:00",
				"2016-01-26T00:10:21-03:00",
				"2016-01-26T22:15:30+03:00",
			},
			locales: []Locale{
				LocaleEnUS,
				LocaleEnUS,
				LocaleEnUS,
				LocaleEnUS,
			},
		},
		{
			layout: "Пон, 2 Янв 2006 15:4:5 -0700",
		},
	}
)

func TestLayoutValidator(t *testing.T) {
	ld := NewLocaleDetector()
	for _, ltd := range testingLayoutsData {
		ld.prepareLayout(ltd.layout)
		for i, m := range ltd.matches {
			if !ld.validateValue(ltd.layout, m) {
				t.Errorf("'%s' not matches to '%s' last error position = %d\n", m, ltd.layout, ld.errorPosition())
			}
			locale := ld.detectLocale(m)
			if !compareLocales(locale, ltd.locales[i]) {
				t.Errorf("locales detect error, expected '%s', result '%s'\n", ltd.locales[i], locale)
			}
		}
		for _, u := range ltd.unmatches {
			if ld.validateValue(ltd.layout, u) {
				t.Errorf("'%s' matches to '%s'\n", u, ltd.layout)
			}
		}
	}
}

// TODO: locale groups.
var englishLocales = newSet(LocaleEnUS, LocaleEnGB)

func compareLocales(a, b Locale) bool {
	if a == b {
		return true
	}
	if englishLocales.Has(a) && englishLocales.Has(b) {
		return true
	}
	return false
}

func TestCompareLocales(t *testing.T) {
	if !compareLocales(LocaleEnGB, LocaleEnUS) {
		t.Errorf("compareLocales not works as expected")
	}
}

func TestParsing(t *testing.T) {
	ld := NewLocaleDetector()
	for _, ltd := range testingLayoutsData {
		for i, formatted := range ltd.matches {
			dt, err := ld.Parse(ltd.layout, formatted)
			if err != nil {
				t.Errorf("error parsing '%s' with layout: '%s' [error:%s]\n", formatted, ltd.layout, err.Error())
			} else {
				restored := Format(dt, ltd.layout, ltd.locales[i])
				if restored != formatted {
					dt2, err2 := ld.Parse(ltd.layout, restored)
					if err2 != nil {
						t.Errorf("parsed time '%s' (%s) does not match restored time '%s'\n", formatted, dt, restored)
					}
					if dt2.Unix() != dt.Unix() {
						t.Errorf("restored time '%v' != parsed time '%v' (restored='%s')", dt2, dt, restored)
					}
				}
			}
		}
	}
}
