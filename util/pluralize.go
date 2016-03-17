package util

import "regexp"

// Source fom https://github.com/acsellers/inflections
//
// NounPluralize is based on the algorithm described at
// www.csse.monash.edu.au/~damian/papers/HTML/Plurals.html
// It doesn't deal with classical pluralizations at the moment.
// It could have been based on Rails pluralize
// implementation, but that code is not clear due to the
// way it was implemented. It is only implemented for English.
//
// Copyright (c) 2013, Andrew Sellers <andrew@andrewcsellers.com>
// All rights reserved.
//
// License: BSD 3-clause
func Pluralize(str string) string {
	switch {
	case _UserDefined(str):
		return userDefinedPluralize(str)
	case _NonInflecting(str):
		return str
	case _Pronoun(str):
		return pronounPluralize(str)
	case _StandardIrregular(str):
		return stdIrregularPluralize(str)
	case _CommonSuffixes(str):
		return commonSuffixPluralize(str)
	case _ExAssimilatedClassic(str):
		return exAssimilatedClassPluralize(str)
	case _UmAssimilatedClassic(str):
		return umAssimilatedClassPluralize(str)
	case _AAssimilatedClassic(str):
		return aAssimilatedClassPluralize(str)
	case _ChSh(str):
		return chShPluralize(str)
	case _Ss(str):
		return ssPluralize(str)
	case _F(str):
		return fPluralize(str)
	case _Fe(str):
		return fePluralize(str)
	case _Y(str):
		return yPluralize(str)
	case _O(str):
		return oPluralize(str)
	//case _CompoundWord(str):
	//return compoundPluralize(str)
	default:
		return defaultPluralize(str)
	}
}

var UserDefinedInflections map[string]string

func init() {
	UserDefinedInflections = make(map[string]string)
}
func _UserDefined(str string) bool {
	_, found := UserDefinedInflections[str]
	return found
}
func userDefinedPluralize(str string) string {
	return UserDefinedInflections[str]
}

var uninflected = []string{"bison", "bream", "breeches", "britches", "carp", "chassis",
	"clippers", "cod", "contretemps", "corps", "debris", "diabetes", "djinn", "eland", "elk",
	"flounder", "gallows", "graffiti", "headquarters", "herpes", "high-jinks", "homework",
	"innings", "jackanapes", "mackarels", "measles", "mews", "mumps", "news", "pincers",
	"pliers", "proceedings", "rabies", "salmon", "scissors", "sea-bass", "series", "shears",
	"species", "swine", "trout", "tuna", "whiting", "wildebeest"}

func _NonInflecting(str string) bool {
	for _, candidate := range uninflected {
		if candidate == str {
			return true
		}
	}
	return false
}

var pronouns = map[string]string{"i": "we", "me": "us", "myself": "ourselves",
	"you": "you", "thou": "you", "thee": "you", "yourself": "yourself", "thyself": "yourself",
	"she": "they", "he": "they", "it": "they", "they": "they", "her": "them", "him": "them",
	"them": "them", "herself": "themselves", "himself": "themselves", "itself": "themselves",
	"themself": "themselves", "oneself": "oneselves", "mine": "ours", "yours": "yours",
	"thine": "yours", "hers": "theirs", "his": "theirs", "its": "theirs", "theirs": "theirs",
	"my": "our", "your": "your", "thy": "your", "their": "their"}

func _Pronoun(str string) bool {
	_, found := pronouns[str]
	return found
}
func pronounPluralize(str string) string {
	return pronouns[str]
}

var standardIrregular = map[string]string{"beef": "beefs", "child": "children",
	"ephemeris": "ephemerides", "money": "monies", "mongoose": "mongooses",
	"mythos": "mythoi", "octopus": "octopuses", "ox": "oxen", "soliloquy": "soliloquies",
	"trilby": "trilbys"}

func _StandardIrregular(str string) bool {
	_, found := standardIrregular[str]
	return found
}
func stdIrregularPluralize(str string) string {
	return standardIrregular[str]
}

var commonSuffixes = regexp.MustCompile("[man|[lm]ouse|tooth|goose|foot|zoon|[csx]is]$")

func _CommonSuffixes(str string) bool {
	return commonSuffixes.MatchString(str)
}
func commonSuffixPluralize(str string) string {
	switch {
	case str[len(str)-3:] == "man":
		return str[:len(str)-3] + "men"
	case str[len(str)-4:] == "ouse":
		return str[:len(str)-4] + "ice"
	case str[len(str)-5:] == "tooth":
		return str[:len(str)-4] + "eeth"
	case str[len(str)-4:] == "oose":
		return str[:len(str)-4] + "eese"
	case str[:len(str)-3] == "oot":
		return str[:len(str)-3] + "eet"
	case str[:len(str)-3] == "oon":
		return str[:len(str)-3] + "oa"
	case str[:len(str)-2] == "is":
		return str[:len(str)-2] + "es"
	}
	return str + "s"
}

func _ExAssimilatedClassic(str string) bool {
	return str == "codex" || str == "murex" || str == "silex"
}
func exAssimilatedClassPluralize(str string) string {
	return str[:len(str)-2] + "ices"
}

func _AAssimilatedClassic(str string) bool {
	return str == "alumna" || str == "alga" || str == "vertebra"
}
func aAssimilatedClassPluralize(str string) string {
	return str + "e"
}

var uma = []string{"agendum", "bacterium", "candelabrum", "datum",
	"desideratum", "erratum", "extremm", "stratum", "ovum"}

func _UmAssimilatedClassic(str string) bool {
	for _, candidate := range uma {
		if str == candidate {
			return true
		}
	}
	return false
}
func umAssimilatedClassPluralize(str string) string {
	return str[:len(str)-2] + "a"
}

func _ChSh(str string) bool {
	return str[len(str)-2:] == "ch" || str[len(str)-2:] == "sh"
}
func chShPluralize(str string) string {
	return str[:len(str)-1] + "hes"
}
func _Ss(str string) bool {
	return str[len(str)-2:] == "ss"
}
func ssPluralize(str string) string {
	return str + "es"
}

var (
	aeolf = regexp.MustCompile("[aeo]lf$")
	eaf   = regexp.MustCompile("[^d]eaf$")
)

func _F(str string) bool {
	if aeolf.MatchString(str) || eaf.MatchString(str) || str[len(str)-3:] == "arf" {
		return true
	}
	return false
}
func fPluralize(str string) string {
	return str[:len(str)-1] + "ves"
}

var nlwife = regexp.MustCompile("[nlw]ife$")

func _Fe(str string) bool {
	return nlwife.MatchString(str)
}
func fePluralize(str string) string {
	return str[:len(str)-2] + "ves"
}

var (
	vowelY  = regexp.MustCompile("[aeiou]y$")
	properY = regexp.MustCompile("^[A-Z].*y$")
)

func _Y(str string) bool {
	return str[len(str)-1] == 'y'
}
func yPluralize(str string) string {
	if vowelY.MatchString(str) || properY.MatchString(str) {
		return str + "s"
	} else {
		return str[:len(str)-2] + "ies"
	}
}

var (
	oToOs = []string{"albino", "archipelago", "armadillo", "commando",
		"ditto", "dynamo", "embryo", "fiasco", "generalissimo", "ghetto",
		"guano", "inferno", "jumbo", "lingo", "lumbago", "magneto",
		"manifesto", "medico", "octavo", "photo", "pro", "quarto", "rhino",
		"stylo"}
	anglOToOs = []string{"alto", "basso", "canto", "contralto",
		"crescendo", "solo", "soprano", "tempo"}
)

func _O(str string) bool {
	if str[len(str)-1] != 'o' {
		return false
	}

	for _, candidate := range oToOs {
		if candidate == str {
			return true
		}
	}
	for _, candidate := range anglOToOs {
		if candidate == str {
			return true
		}
	}
	return false
}
func oPluralize(str string) string {
	return str + "s"
}

func compoundPluralize(str string) string {
	return str
}

func defaultPluralize(str string) string {
	return str + "s"
}
