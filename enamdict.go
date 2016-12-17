/*
 * Copyright (c) 2016 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"io"

	"github.com/FooSoft/jmdict"
)

func computeJmnedictTagMeta(entities map[string]string) map[string]dbTagMeta {
	tags := make(map[string]dbTagMeta)

	for name, value := range entities {
		tag := dbTagMeta{Notes: value}

		switch name {
		case "company":
		case "fem":
		case "given":
		case "masc":
		case "organization":
		case "person":
		case "place":
		case "product":
		case "station":
		case "surname":
		case "unclass":
		case "work":
			tag.Class = "name"
			tag.Order = 4
		}

		tags[name] = tag
	}

	return tags
}

func extractJmnedictTerms(enamdictEntry jmdict.JmnedictEntry) []dbTerm {
	var terms []dbTerm

	convert := func(reading jmdict.JmnedictReading, kanji *jmdict.JmnedictKanji) {
		if kanji != nil && hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var term dbTerm
		term.addTags(reading.Information...)

		if kanji == nil {
			term.Expression = reading.Reading
			term.addTags(reading.Information...)
		} else {
			term.Expression = kanji.Expression
			term.Reading = reading.Reading
			term.addTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					term.addTags(priority)
				}
			}
		}

		for _, trans := range enamdictEntry.Translations {
			term.Glossary = append(term.Glossary, trans.Translations...)
			term.addTags(trans.NameTypes...)
		}

		terms = append(terms, term)
	}

	if len(enamdictEntry.Kanji) > 0 {
		for _, kanji := range enamdictEntry.Kanji {
			for _, reading := range enamdictEntry.Readings {
				convert(reading, &kanji)
			}
		}
	} else {
		for _, reading := range enamdictEntry.Readings {
			convert(reading, nil)
		}
	}

	return terms
}

func exportJmnedictDb(outputDir, title string, reader io.Reader, flags int) error {
	dict, entities, err := jmdict.LoadJmnedictNoTransform(reader)
	if err != nil {
		return err
	}

	var terms dbTermList
	for _, e := range dict.Entries {
		terms = append(terms, extractJmnedictTerms(e)...)
	}

	return writeDb(
		outputDir,
		title,
		terms.crush(),
		nil,
		computeJmnedictTagMeta(entities),
		flags&flagPretty == flagPretty,
	)
}
