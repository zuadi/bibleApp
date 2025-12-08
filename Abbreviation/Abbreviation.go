package abbreviation

import "strings"

func Replace(input string) string {

	switch {
	case input == "Gn":
		return strings.Replace(input, "Gn", "Gen", 1)
	case input == "Lv":
		return strings.Replace(input, "Lv", "Lev", 1)
	case input == "Nm":
		return strings.Replace(input, "Nm", "Num", 1)
	case input == "Du":
		return strings.Replace(input, "Du", "Deu", 1)
	case input == "Dt":
		return strings.Replace(input, "Dt", "Deu", 1)
	case input == "Jdg":
		return strings.Replace(input, "Jdg", "Judg", 1)
	case input == "Jgs":
		return strings.Replace(input, "Jgs", "Judg", 1)
	case input == "Rth":
		return strings.Replace(input, "Rth", "Rut", 1)
	case input == "Sm":
		return strings.Replace(input, "Sm", "Sam", 1)
	case input == "Kgs":
		return strings.Replace(input, "Kgs", "Kin", 1)
	case input == "Kg":
		return strings.Replace(input, "Kg", "Kin", 1)
	case input == "Jb":
		return strings.Replace(input, "Jb", "Job", 1)
	case input == "Psalms":
		return strings.Replace(input, "Psalms", "Psalm", 1)
	case input == "Prv":
		return strings.Replace(input, "Prv", "Pro", 1)
	case input == "Sg":
		return strings.Replace(input, "Sg", "Song", 1)
	case input == "Hld":
		return strings.Replace(input, "Hld", "Hohe", 1)
	case input == "Klgl":
		return strings.Replace(input, "Klgl", "Klage", 1)
	case input == "Ezk":
		return strings.Replace(input, "Ezk", "Eze", 1)
	case input == "Dan":
		return strings.Replace(input, "Dn", "Dan", 1)
	case input == "Hb":
		return strings.Replace(input, "Hb", "Hab", 1)
	case input == "Obd":
		return strings.Replace(input, "Obd", "Oba", 1)
	case input == "Zp":
		return strings.Replace(input, "Zp", "Zep", 1)
	case input == "Hg":
		return strings.Replace(input, "Hg", "Hag", 1)
	case input == "Mt":
		return strings.Replace(input, "Mt", "Mat", 1)
	case input == "Mk":
		return strings.Replace(input, "Mk", "Mar", 1)
	case input == "Mrk":
		return strings.Replace(input, "Mrk", "Mar", 1)
	case input == "Lk":
		return strings.Replace(input, "Lk", "Luk", 1)
	case input == "Jn":
		return strings.Replace(input, "Jn", "Joh", 1)
	case input == "Jhn":
		return strings.Replace(input, "Jhn", "Joh", 1)
	case input == "Apg":
		return strings.Replace(input, "Apg", "Apos", 1)
	case input == "Rm":
		return strings.Replace(input, "Rm", "Rom", 1)
	case input == "Php":
		return strings.Replace(input, "Php", "Phili", 1)
	case input == "Phm":
		return strings.Replace(input, "Phm", "Phile", 1)
	case input == "Phlm":
		return strings.Replace(input, "Phlm", "Phile", 1)
	case input == "Pt":
		return strings.Replace(input, "Pt", "Pet", 1)
	case input == "Offb":
		return strings.Replace(input, "Offb", "Offen", 1)
	case input == "S. ":
		return strings.Replace(input, "S. ", "", 1)
	}

	return input
}
