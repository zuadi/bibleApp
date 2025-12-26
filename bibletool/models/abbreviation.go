package models

import "strings"

func ReplaceAbbreviation(input string) string {
	switch input {
	case "Gn":
		return strings.Replace(input, "Gn", "Gen", 1)
	case "Lv":
		return strings.Replace(input, "Lv", "Lev", 1)
	case "Nm":
		return strings.Replace(input, "Nm", "Num", 1)
	case "Du":
		return strings.Replace(input, "Du", "Deu", 1)
	case "Dt":
		return strings.Replace(input, "Dt", "Deu", 1)
	case "Jdg":
		return strings.Replace(input, "Jdg", "Judg", 1)
	case "Jgs":
		return strings.Replace(input, "Jgs", "Judg", 1)
	case "Rth":
		return strings.Replace(input, "Rth", "Rut", 1)
	case "Sm":
		return strings.Replace(input, "Sm", "Sam", 1)
	case "Kgs":
		return strings.Replace(input, "Kgs", "Kin", 1)
	case "Kg":
		return strings.Replace(input, "Kg", "Kin", 1)
	case "Jb":
		return strings.Replace(input, "Jb", "Job", 1)
	case "Psalms":
		return strings.Replace(input, "Psalms", "Psalm", 1)
	case "Prv":
		return strings.Replace(input, "Prv", "Pro", 1)
	case "Sg":
		return strings.Replace(input, "Sg", "Song", 1)
	case "Hld":
		return strings.Replace(input, "Hld", "Hohe", 1)
	case "Klgl":
		return strings.Replace(input, "Klgl", "Klage", 1)
	case "Ezk":
		return strings.Replace(input, "Ezk", "Eze", 1)
	case "Dan":
		return strings.Replace(input, "Dn", "Dan", 1)
	case "Hb":
		return strings.Replace(input, "Hb", "Hab", 1)
	case "Obd":
		return strings.Replace(input, "Obd", "Oba", 1)
	case "Zp":
		return strings.Replace(input, "Zp", "Zep", 1)
	case "Hg":
		return strings.Replace(input, "Hg", "Hag", 1)
	case "Mt":
		return strings.Replace(input, "Mt", "Mat", 1)
	case "Mk":
		return strings.Replace(input, "Mk", "Mar", 1)
	case "Mrk":
		return strings.Replace(input, "Mrk", "Mar", 1)
	case "Lk":
		return strings.Replace(input, "Lk", "Luk", 1)
	case "Jn":
		return strings.Replace(input, "Jn", "Joh", 1)
	case "Jhn":
		return strings.Replace(input, "Jhn", "Joh", 1)
	case "Apg":
		return strings.Replace(input, "Apg", "Apos", 1)
	case "Rm":
		return strings.Replace(input, "Rm", "Rom", 1)
	case "Php":
		return strings.Replace(input, "Php", "Phili", 1)
	case "Phm":
		return strings.Replace(input, "Phm", "Phile", 1)
	case "Phlm":
		return strings.Replace(input, "Phlm", "Phile", 1)
	case "Pt":
		return strings.Replace(input, "Pt", "Pet", 1)
	case "Offb":
		return strings.Replace(input, "Offb", "Offen", 1)
	case "S. ":
		return strings.Replace(input, "S. ", "", 1)
	}
	return input
}
