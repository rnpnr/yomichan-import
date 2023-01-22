package yomichan

import (
	"regexp"
	"strings"

	zig "foosoft.net/projects/zero-epwing-go"
)

type daijisenExtractor struct {
	partsExp     *regexp.Regexp
	expShapesExp *regexp.Regexp
	expMultiExp  *regexp.Regexp
	expVarExp    *regexp.Regexp
	readGroupExp *regexp.Regexp
	metaExp      *regexp.Regexp
	v5Exp        *regexp.Regexp
	v1Exp        *regexp.Regexp
}

func makeDaijisenExtractor() epwingExtractor {
	return &daijisenExtractor{
		partsExp:     regexp.MustCompile(`([^【]+)(?:【(.*)】)?`),
		expShapesExp: regexp.MustCompile(`[×△＝‐]+`),
		expMultiExp:  regexp.MustCompile(`】[^【】]*【`),
		expVarExp:    regexp.MustCompile(`（([^）]*)）`),
		readGroupExp: regexp.MustCompile(`[‐・]+`),
		metaExp:      regexp.MustCompile(`［([^］]*)］`),
		v5Exp:        regexp.MustCompile(`(動.[四五](［[^］]+］)?)|(動..二)`),
		v1Exp:        regexp.MustCompile(`(動..一)`),
	}
}

func (e *daijisenExtractor) extractTerms(entry zig.BookEntry, sequence int) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	var expressions []string
	if expression := matches[2]; len(expression) > 0 {
		expression = e.expMultiExp.ReplaceAllString(expression, "・")
		expression = e.expShapesExp.ReplaceAllString(expression, "")
		for _, split := range strings.Split(expression, "・") {
			splitInc := e.expVarExp.ReplaceAllString(split, "$1")
			expressions = append(expressions, splitInc)
			if split != splitInc {
				splitExc := e.expVarExp.ReplaceAllLiteralString(split, "")
				expressions = append(expressions, splitExc)
			}
		}
	}

	var reading string
	if reading = matches[1]; len(reading) > 0 {
		reading = e.readGroupExp.ReplaceAllLiteralString(reading, "")
		reading = e.expVarExp.ReplaceAllLiteralString(reading, "")
	}

	var tags []string
	for _, split := range strings.Split(entry.Text, "\n") {
		if matches := e.metaExp.FindStringSubmatch(split); matches != nil {
			for _, tag := range strings.Split(matches[1], "・") {
				tags = append(tags, tag)
			}
		}
	}

	var terms []dbTerm
	if len(expressions) == 0 {
		term := dbTerm{
			Expression: reading,
			Glossary:   []any{entry.Text},
			Sequence:   sequence,
		}

		e.exportRules(&term, tags)
		terms = append(terms, term)

	} else {
		for _, expression := range expressions {
			term := dbTerm{
				Expression: expression,
				Reading:    reading,
				Glossary:   []any{entry.Text},
				Sequence:   sequence,
			}

			e.exportRules(&term, tags)
			terms = append(terms, term)
		}
	}

	return terms
}

func (*daijisenExtractor) extractKanji(entry zig.BookEntry) []dbKanji {
	return nil
}

func (e *daijisenExtractor) exportRules(term *dbTerm, tags []string) {
	for _, tag := range tags {
		if tag == "形" {
			term.addRules("adj-i")
		} else if tag == "動サ変" && (strings.HasSuffix(term.Expression, "する") || strings.HasSuffix(term.Expression, "為る")) {
			term.addRules("vs")
		} else if term.Expression == "来る" {
			term.addRules("vk")
		} else if e.v5Exp.MatchString(tag) {
			term.addRules("v5")
		} else if e.v1Exp.MatchString(tag) {
			term.addRules("v1")
		}
	}
}

func (*daijisenExtractor) getRevision() string {
	return "daijisen2"
}

func (*daijisenExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		0xa121: " ",
		0xa122: "¡",
		0xa123: "¢",
		0xa124: "£",
		0xa125: "¤",
		0xa126: "¥",
		0xa127: "¦",
		0xa128: "§",
		0xa129: "¨",
		0xa12a: "©",
		0xa12b: "ª",
		0xa12c: "«",
		0xa12d: "¬",
		0xa12e: "­",
		0xa12f: "®",
		0xa130: "¯",
		0xa131: "°",
		0xa132: "±",
		0xa133: "²",
		0xa134: "³",
		0xa135: "´",
		0xa136: "µ",
		0xa137: "¶",
		0xa138: "·",
		0xa139: "¸",
		0xa13a: "¹",
		0xa13b: "º",
		0xa13c: "»",
		0xa13d: "¼",
		0xa13e: "½",
		0xa13f: "¾",
		0xa140: "¿",
		0xa141: "À",
		0xa142: "Á",
		0xa143: "Â",
		0xa144: "Ã",
		0xa145: "Ä",
		0xa146: "Å",
		0xa147: "Æ",
		0xa148: "Ç",
		0xa149: "È",
		0xa14a: "É",
		0xa14b: "Ê",
		0xa14c: "Ë",
		0xa14d: "Ì",
		0xa14e: "Í",
		0xa14f: "Î",
		0xa150: "Ï",
		0xa151: "Ð",
		0xa152: "Ñ",
		0xa153: "Ò",
		0xa154: "Ó",
		0xa155: "Ô",
		0xa156: "Õ",
		0xa157: "Ö",
		0xa158: "×",
		0xa159: "Ø",
		0xa15a: "Ù",
		0xa15b: "Ú",
		0xa15c: "Û",
		0xa15d: "Ü",
		0xa15e: "Ý",
		0xa15f: "Þ",
		0xa160: "ß",
		0xa161: "à",
		0xa162: "á",
		0xa163: "â",
		0xa164: "ã",
		0xa165: "ä",
		0xa166: "å",
		0xa167: "æ",
		0xa168: "ç",
		0xa169: "è",
		0xa16a: "é",
		0xa16b: "ê",
		0xa16c: "ë",
		0xa16d: "ì",
		0xa16e: "í",
		0xa16f: "î",
		0xa170: "ï",
		0xa171: "ð",
		0xa172: "ñ",
		0xa173: "ò",
		0xa174: "ó",
		0xa175: "ô",
		0xa176: "õ",
		0xa177: "ö",
		0xa178: "÷",
		0xa179: "ø",
		0xa17a: "ù",
		0xa17b: "ú",
		0xa17c: "û",
		0xa17d: "ü",
		0xa17e: "ý",
		0xa221: "þ",
		0xa222: "ÿ",
		0xa223: "Ā",
		0xa224: "ā",
		0xa225: "Ă",
		0xa226: "ă",
		0xa227: "Ą",
		0xa228: "ą",
		0xa229: "Ć",
		0xa22a: "ć",
		0xa22b: "Ĉ",
		0xa22c: "ĉ",
		0xa22d: "Ċ",
		0xa22e: "ċ",
		0xa22f: "Č",
		0xa230: "č",
		0xa231: "Ď",
		0xa232: "ď",
		0xa233: "Đ",
		0xa234: "đ",
		0xa235: "Ē",
		0xa236: "ē",
		0xa237: "Ĕ",
		0xa238: "ĕ",
		0xa239: "Ė",
		0xa23a: "ė",
		0xa23b: "Ę",
		0xa23c: "ę",
		0xa23d: "Ě",
		0xa23e: "ě",
		0xa23f: "Ĝ",
		0xa240: "ĝ",
		0xa241: "Ğ",
		0xa242: "ğ",
		0xa243: "Ġ",
		0xa244: "ġ",
		0xa245: "Ģ",
		0xa246: "ģ",
		0xa247: "Ĥ",
		0xa248: "ĥ",
		0xa249: "Ħ",
		0xa24a: "ħ",
		0xa24b: "Ĩ",
		0xa24c: "ĩ",
		0xa24d: "Ī",
		0xa24e: "ī",
		0xa24f: "Ĭ",
		0xa250: "ĭ",
		0xa251: "Į",
		0xa252: "į",
		0xa253: "İ",
		0xa254: "ı",
		0xa255: "Ĳ",
		0xa256: "ĳ",
		0xa257: "Ĵ",
		0xa258: "ĵ",
		0xa259: "Ķ",
		0xa25a: "ķ",
		0xa25b: "ĸ",
		0xa25c: "Ĺ",
		0xa25d: "ĺ",
		0xa25e: "Ļ",
		0xa25f: "ļ",
		0xa260: "Ľ",
		0xa261: "ľ",
		0xa262: "Ŀ",
		0xa263: "ŀ",
		0xa264: "Ł",
		0xa265: "ł",
		0xa266: "Ń",
		0xa267: "ń",
		0xa268: "Ņ",
		0xa269: "ņ",
		0xa26a: "Ň",
		0xa26b: "ň",
		0xa26c: "ŉ",
		0xa26d: "Ŋ",
		0xa26e: "ŋ",
		0xa26f: "Ō",
		0xa270: "ō",
		0xa271: "Ŏ",
		0xa272: "ŏ",
		0xa273: "Ő",
		0xa274: "ő",
		0xa275: "Œ",
		0xa276: "œ",
		0xa277: "Ŕ",
		0xa278: "ŕ",
		0xa279: "Ŗ",
		0xa27a: "ŗ",
		0xa27b: "Ř",
		0xa27c: "ř",
		0xa27d: "Ś",
		0xa27e: "ś",
		0xa321: "Ŝ",
		0xa322: "ŝ",
		0xa323: "Ş",
		0xa324: "ş",
		0xa325: "Š",
		0xa326: "š",
		0xa327: "Ţ",
		0xa328: "ţ",
		0xa329: "Ť",
		0xa32a: "ť",
		0xa32b: "Ŧ",
		0xa32c: "ŧ",
		0xa32d: "Ũ",
		0xa32e: "ũ",
		0xa32f: "Ū",
		0xa330: "ū",
		0xa331: "Ŭ",
		0xa332: "ŭ",
		0xa333: "Ů",
		0xa334: "ů",
		0xa335: "Ű",
		0xa336: "ű",
		0xa337: "Ų",
		0xa338: "ų",
		0xa339: "Ŵ",
		0xa33a: "ŵ",
		0xa33b: "Ŷ",
		0xa33c: "ŷ",
		0xa33d: "Ÿ",
		0xa33e: "Ź",
		0xa33f: "ź",
		0xa340: "Ż",
		0xa341: "ż",
		0xa342: "Ž",
		0xa343: "ž",
		0xa344: "ſ",
		0xa34d: "ƒ",
		0xa34e: "ˆ",
		0xa34f: "˜",
		0xa362: "Ḍ",
		0xa363: "Ḥ",
		0xa364: "Ṛ",
		0xa365: "Ṣ",
		0xa366: "Ẓ",
		0xa367: "ạ́",
		0xa368: "ḅ",
		0xa369: "ī",
		0xa36a: "ḍ",
		0xa36b: "ḥ",
		0xa36c: "i",
		0xa36d: "ị̄",
		0xa36e: "ị́",
		0xa36f: "ị̈",
		0xa370: "î",
		0xa371: "ḳ",
		0xa372: "ṁ",
		0xa373: "ṃ",
		0xa374: "ṅ",
		0xa375: "ṇ",
		0xa376: "ṛ",
		0xa377: "ṣ",
		0xa378: "ṭ",
		0xa379: "ẓ",
	}
}

func (*daijisenExtractor) getFontWide() map[int]string {
	return map[int]string{
		0xb021: "嗩",
		0xb022: "盎",
		0xb023: "盔",
		0xb024: "荽",
		0xb025: "芡",
		0xb026: "蕹",
		0xb027: "螋",
		0xb028: "蛺",
		0xb029: "蚨",
		0xb02a: "蠊",
		0xb02b: "闈",
		0xb02c: "獒",
		0xb02d: "犰",
		0xb02e: "鑊",
		0xb02f: "眶",
		0xb030: "睽",
		0xb031: "熒",
		0xb032: "莆",
		0xb033: "芮",
		0xb034: "苕",
		0xb035: "蝤",
		0xb036: "獫",
		0xb037: "狳",
		0xb038: "猻",
		0xb039: "晡",
		0xb03a: "曛",
		0xb03b: "洄",
		0xb03c: "洹",
		0xb03d: "硇",
		0xb03e: "擻",
		0xb03f: "拄",
		0xb040: "瞟",
		0xb041: "眙",
		0xb042: "眚",
		0xb043: "芰",
		0xb044: "萏",
		0xb045: "蘅",
		0xb046: "螵",
		0xb047: "蛑",
		0xb048: "狁",
		0xb049: "狻",
		0xb04a: "猢",
		0xb04b: "肫",
		0xb04c: "臃",
		0xb04d: "刖",
		0xb04e: "脞",
		0xb04f: "鏟",
		0xb050: "坷",
		0xb051: "畎",
		0xb052: "譙",
		0xb053: "蘼",
		0xb054: "菇",
		0xb055: "螬",
		0xb056: "虺",
		0xb057: "膘",
		0xb058: "澌",
		0xb059: "涿",
		0xb05a: "垸",
		0xb05b: "詿",
		0xb05c: "謭",
		0xb05d: "訕",
		0xb05e: "詘",
		0xb05f: "撿",
		0xb061: "賾",
		0xb062: "臬",
		0xb063: "葒",
		0xb064: "萁",
		0xb065: "蕤",
		0xb066: "翬",
		0xb067: "翥",
		0xb068: "炫",
		0xb069: "榷",
		0xb06a: "棖",
		0xb06b: "鑣",
		0xb06c: "坨",
		0xb06e: "儈",
		0xb06f: "綈",
		0xb070: "踹",
		0xb071: "橛",
		0xb072: "椐",
		0xb073: "憨",
		0xb074: "緦",
		0xb075: "繒",
		0xb076: "黧",
		0xb077: "輞",
		0xb078: "軔",
		0xb07a: "沔",
		0xb07b: "洱",
		0xb07c: "浠",
		0xb07d: "欏",
		0xb07e: "桕",
		0xb121: "桫",
		0xb122: "怍",
		0xb123: "悱",
		0xb124: "戕",
		0xb125: "緗",
		0xb126: "蘞",
		0xb127: "蚱",
		0xb128: "蚍",
		0xb129: "螭",
		0xb12a: "蚜",
		0xb12b: "轔",
		0xb12c: "鼹",
		0xb12d: "闋",
		0xb12e: "駙",
		0xb12f: "涪",
		0xb130: "渲",
		0xb131: "棼",
		0xb132: "鐲",
		0xb133: "卬",
		0xb134: "厓",
		0xb135: "唵",
		0xb136: "啡",
		0xb137: "墝",
		0xb138: "墩",
		0xb139: "壔",
		0xb13b: "kg",
		0xb13c: "cc",
		0xb13d: "畺",
		0xb13e: "仿",
		0xb13f: "厲",
		0xb140: "饜",
		0xb141: "嘎",
		0xb142: "壠",
		0xb143: "氐",
		0xb144: "你",
		0xb145: "佉",
		0xb146: "淸",
		0xb148: "颺",
		0xb149: "嘻",
		0xb14a: "嘰",
		0xb14b: "噉",
		0xb14c: "噲",
		0xb14d: "奝",
		0xb14f: "丰",
		0xb150: "繇",
		0xb151: "燄",
		0xb152: "囉",
		0xb153: "、",
		0xb159: "俏",
		0xb15a: "剗",
		0xb15b: "剡",
		0xb15c: "吒",
		0xb15d: "吧",
		0xb15e: "媳",
		0xb15f: "Ⅰ",
		0xb160: "Ⅱ",
		0xb161: "弇",
		0xb162: "傖",
		0xb163: "埿",
		0xb164: "嫩",
		0xb16b: "Ⅲ",
		0xb16c: "Ⅳ",
		0xb16d: "Ⅴ",
		0xb16e: "Ⅵ",
		0xb16f: "Ⅶ",
		0xb170: "Ⅷ",
		0xb171: "Ⅹ",
		0xb172: "垜",
		0xb173: "漪",
		0xb174: "莧",
		0xb175: "陘",
		0xb176: "寵",
		0xb177: "濞",
		0xb178: "邙",
		0xb17a: "∫",
		0xb17c: "沅",
		0xb17d: "濹",
		0xb17e: "鄧",
		0xb268: "扑",
		0xb269: "灤",
		0xb26a: "蔞",
		0xb26b: "蓯",
		0xb26c: "蓰",
		0xb270: "拖",
		0xb271: "蔯",
		0xb272: "邈",
		0xb274: "邛",
		0xb275: "挘",
		0xb276: "挹",
		0xb277: "芎",
		0xb278: "芩",
		0xb279: "薏",
		0xb27b: "帒",
		0xb27c: "帮",
		0xb27d: "幫",
		0xb27e: "毟",
		0xb321: "苆",
		0xb322: "\n㋘",
		0xb323: "\n㋙",
		0xb324: "\n㋚",
		0xb325: "\n㋛",
		0xb326: "\n㋜",
		0xb327: "\n㋝",
		0xb328: "漚",
		0xb329: "荃",
		0xb32a: "莒",
		0xb32b: "惲",
		0xb32c: "愒",
		0xb332: "胳",
		0xb333: "燋",
		0xb334: "毗",
		0xb335: "畯",
		0xb336: "礱",
		0xb338: "璜",
		0xb339: "琨",
		0xb33a: "砰",
		0xb33b: "惋",
		0xb33c: "娌",
		0xb33d: "″",
		0xb342: "腊",
		0xb343: "楣",
		0xb344: "刁",
		0xb345: "邢",
		0xb347: "賖",
		0xb348: "砭",
		0xb349: "㊙",
		0xb34b: "牓",
		0xb34c: "痎",
		0xb34d: "瘀",
		0xb34e: "惕",
		0xb34f: "忡",
		0xb350: "鑲",
		0xb351: "閫",
		0xb352: "閽",
		0xb353: "髡",
		0xb354: "划",
		0xb355: "檞",
		0xb356: "瘙",
		0xb357: "贛",
		0xb358: "圳",
		0xb359: "塌",
		0xb35a: "夤",
		0xb35b: "晷",
		0xb35c: "榨",
		0xb35d: "礴",
		0xb35e: "枘",
		0xb35f: "珉",
		0xb360: "琮",
		0xb361: "癭",
		0xb362: "婺",
		0xb363: "宓",
		0xb364: "柒",
		0xb365: "殂",
		0xb366: "縈",
		0xb367: "愜",
		0xb368: "祆",
		0xb369: "祜",
		0xb36a: "櫧",
		0xb36b: "，",
		0xb36e: "徉",
		0xb36f: "徜",
		0xb371: "靛",
		0xb372: "籮",
		0xb373: "縐",
		0xb374: "鸝",
		0xb375: "鸇",
		0xb376: "鷉",
		0xb377: "鷚",
		0xb378: "鸊",
		0xb379: "鷴",
		0xb37a: "栬",
		0xb37b: "桲",
		0xb37c: "裊",
		0xb37d: "釃",
		0xb37e: "醅",
		0xb421: "鵒",
		0xb422: "鴞",
		0xb423: "虢",
		0xb424: "↔",
		0xb425: "烑",
		0xb426: "煆",
		0xb427: "睜",
		0xb428: "睢",
		0xb429: "筎",
		0xb42a: "汴",
		0xb42b: "糙",
		0xb42c: "繳",
		0xb42d: "珧",
		0xb42e: "咖",
		0xb42f: "筠",
		0xb430: "閒",
		0xb431: "帔",
		0xb432: "幘",
		0xb433: "鱘",
		0xb434: "鵂",
		0xb435: "飥",
		0xb436: "∘",
		0xb437: "翎",
		0xb438: "骶",
		0xb439: "邡",
		0xb43a: "裰",
		0xb43b: "鰳",
		0xb43c: "鰣",
		0xb43d: "巋",
		0xb43e: "阼",
		0xb43f: "ħ",
		0xb440: "醃",
		0xb441: "雒",
		0xb442: "雞",
		0xb443: "魦",
		0xb444: "褚",
		0xb445: "鯧",
		0xb446: "鯪",
		0xb447: "厴",
		0xb448: "陔",
		0xb449: "邳",
		0xb44a: "邶",
		0xb44b: "〻",
		0xb44c: "ノ",
		0xb44d: "鈸",
		0xb44e: "逭",
		0xb44f: "荇",
		0xb450: "菀",
		0xb451: "孽",
		0xb452: "麇",
		0xb453: "瘵",
		0xb454: "痱",
		0xb455: "、",
		0xb456: "イ̇",
		0xb457: "℧",
		0xb458: "跑",
		0xb459: "剕",
		0xb45a: "鰶",
		0xb45b: "褰",
		0xb45c: "窳",
		0xb45d: "郿",
		0xb45e: "郅",
		0xb45f: "龑",
		0xb460: "紓",
		0xb461: "絁",
		0xb462: "豳",
		0xb463: "劂",
		0xb464: "嚕",
		0xb465: "哆",
		0xb54b: "錘",
		0xb54c: "緌",
		0xb54d: "蟖",
		0xb54e: "顬",
		0xb54f: "劓",
		0xb550: "蒯",
		0xb551: "勖",
		0xb552: "蜓",
		0xb553: "殮",
		0xb554: "屣",
		0xb555: "嬀",
		0xb556: "婕",
		0xb557: "娓",
		0xb558: "嬙",
		0xb559: "喈",
		0xb55a: "カ゚",
		0xb55b: "ケ゚",
		0xb55c: "蠼",
		0xb55d: "靚",
		0xb55e: "鏁",
		0xb55f: "鯁",
		0xb560: "鱺",
		0xb561: "鱭",
		0xb562: "鰱",
		0xb563: "儋",
		0xb564: "佾",
		0xb565: "嫠",
		0xb566: "唼",
		0xb567: "©",
		0xb568: "\nⓐ",
		0xb569: "\nⓑ",
		0xb56a: "\nⓒ",
		0xb56b: "桒",
		0xb56c: "咩",
		0xb56d: "鮏",
		0xb56e: "靏",
		0xb56f: "簱",
		0xb570: "罇",
		0xb571: "沆",
		0xb572: "忞",
		0xb573: "昱",
		0xb574: "荊",
		0xb575: "勛",
		0xb576: "棱",
		0xb577: "涇",
		0xb578: "銈",
		0xb579: "嘈",
		0xb57a: "誾",
		0xb57b: "鉸",
		0xb57c: "摠",
		0xb57d: "鈼",
		0xb57e: "嶸",
		0xb621: "昉",
		0xb622: "昺",
		0xb623: "兗",
		0xb624: "泫",
		0xb625: "昕",
		0xb626: "珙",
		0xb627: "珖",
		0xb628: "琦",
		0xb629: "徧",
		0xb62a: "煜",
		0xb62b: "跆",
		0xb62d: "楨",
		0xb62e: "愷",
		0xb62f: "熲",
		0xb630: "鍰",
		0xb631: "棘",
		0xb632: "稹",
		0xb633: "蕙",
		0xb634: "𩊠",
		0xb635: "奭",
		0xb636: "鋧",
		0xb637: "璟",
		0xb638: "儛",
		0xb639: "鐖",
		0xb63a: "埵",
		0xb63b: "桄",
		0xb63c: "澧",
		0xb63d: "瘖",
		0xb63e: "玫",
		0xb63f: "妤",
		0xb640: "炻",
		0xb641: "釗",
		0xb642: "紝",
		0xb643: "盌",
		0xb644: "羗",
		0xb645: "倜",
		0xb646: "\n㋐",
		0xb647: "\n㋑",
		0xb648: "\n㋒",
		0xb649: "\n㋓",
		0xb64a: "\n㋔",
		0xb64b: "\n㋕",
		0xb64c: "\n㋖",
		0xb64d: "\n㋗",
		0xb64e: "㈠",
		0xb64f: "㈡",
		0xb650: "㈢",
		0xb651: "㈣",
		0xb652: "虗",
		0xb653: "啐",
		0xb654: "跎",
		0xb655: "滇",
		0xb656: "潢",
		0xb657: "燁",
		0xb658: "嶠",
		0xb659: "髹",
		0xb65a: "錡",
		0xb65b: "盦",
		0xb65c: "舢",
		0xb65d: "♨",
		0xb65e: "摹",
		0xb65f: "彀",
		0xb660: "騭",
		0xb661: "惝",
		0xb662: "腭",
		0xb663: "呍",
		0xb664: "擤",
		0xb665: "捥",
		0xb666: "梲",
		0xb667: "踠",
		0xb668: "窠",
		0xb669: "桛",
		0xb66a: "魞",
		0xb66b: "噦",
		0xb66c: "圊",
		0xb66d: "睺",
		0xb66e: "驎",
		0xb66f: "袪",
		0xb670: "彔",
		0xb671: "鱲",
		0xb672: "糈",
		0xb673: "湑",
		0xb674: "楉",
		0xb675: "縠",
		0xb676: "絓",
		0xb677: "簁",
		0xb678: "棰",
		0xb679: "糝",
		0xb67a: "搢",
		0xb67b: "炷",
		0xb67c: "縕",
		0xb67d: "鎺",
		0xb67e: "袘",
		0xb721: "𧘱",
		0xb722: "韛",
		0xb723: "跗",
		0xb724: "糗",
		0xb725: "辦",
		0xb726: "犛",
		0xb727: "獐",
		0xb728: "獱",
		0xb729: "玕",
		0xb72a: "瑇",
		0xb72b: "稭",
		0xb72c: "籙",
		0xb72d: "虬",
		0xb72e: "螈",
		0xb72f: "裑",
		0xb730: "貛",
		0xb731: "鶍",
		0xb732: "鵼",
		0xb733: "麞",
		0xb734: "鼯",
		0xb735: "梣",
		0xb736: "楤",
		0xb737: "槵",
		0xb738: "橅",
		0xb739: "瘭",
		0xb73a: "戶",
		0xb73b: "硨",
		0xb73c: "磲",
		0xb73d: "篊",
		0xb73e: "聱",
		0xb73f: "蘩",
		0xb740: "蜾",
		0xb741: "蜱",
		0xb742: "蠃",
		0xb743: "豇",
		0xb744: "魳",
		0xb745: "魬",
		0xb746: "鮄",
		0xb747: "鯇",
		0xb748: "𩸽",
		0xb749: "鯥",
		0xb74a: "鯷",
		0xb74b: "鰧",
		0xb831: "鱓",
		0xb832: "鱩",
		0xb833: "鱝",
		0xb834: "孒",
		0xb835: "偓",
		0xb836: "汶",
		0xb837: "柷",
		0xb839: "⿐", /* correct radical; renders incorrectly in japanese fonts */
		0xb83a: "瑁",
		0xb83b: "閩",
		0xb83c: "猽",
		0xb83d: "茅",
		0xb83e: "觔",
		0xb83f: "紈",
		0xb840: "醞",
		0xb841: "猨",
		0xb842: "莩",
		0xb843: "橉",
		0xb844: "隄",
		0xb845: "產",
		0xb846: "黑",
		0xb847: "佪",
		0xb848: "枻",
		0xb849: "柀",
		0xb84a: "玦",
		0xb84b: "詡",
		0xb84c: "朓",
		0xb84d: "絺",
		0xb84e: "庾",
		0xb84f: "龐",
		0xb851: "〻",
		0xb852: "⇒",
		0xb853: "璐",
		0xb854: "踔",
		0xb855: "棭",
		0xb856: "燾",
		0xb857: "菝",
		0xb858: "葜",
		0xb859: "獦",
		0xb85a: "氅",
		0xb85b: "簎",
		0xb85c: "芷",
		0xb85d: "淼",
		0xb85e: "丨",
		0xb85f: "乚",
		0xb860: "𠆢",
		0xb861: "亻",
		0xb862: "刂",
		0xb863: "㔾",
		0xb866: "彐",
		0xb867: "⺕",
		0xb868: "⻌",
		0xb869: "辶",
		0xb86a: "辵",
		0xb86b: "阝",
		0xb86d: "忄",
		0xb86e: "⺗",
		0xb86f: "扌",
		0xb871: "氵",
		0xb872: "氺",
		0xb873: "灬",
		0xb874: "狀",
		0xb876: "爫",
		0xb877: "⺤",
		0xb878: "牜",
		0xb879: "犭",
		0xb87a: "耂",
		0xb87c: "疒",
		0xb87d: "礻",
		0xb87e: "禸",
		0xb921: "罒",
		0xb922: "艹",
		0xb923: "⺿",
		0xb924: "衤",
		0xb927: "飠",
		0xb928: "𧾷",
		0xb929: "㓁",
		0xb92b: "杮",
		0xb92c: "茛",
		0xb92d: "痀",
		0xb92e: "噯",
		0xb92f: "芾",
		0xb930: "焄",
		0xb931: "倘",
		0xb932: "暠",
		0xb933: "璘",
		0xb934: "甪",
		0xb935: "觖",
		0xb936: "觫",
		0xb937: "觳",
		0xb938: "觽",
		0xb939: "侗",
		0xb93a: "黃",
		0xb93b: "樾",
		0xb93c: "擎",
		0xb93d: "翟",
		0xb93e: "驌",
		0xb93f: "邕",
		0xb940: "澍",
		0xb941: "灝",
		0xb942: "皥",
		0xb943: "梘",
		0xb944: "嗉",
		0xb945: "痤",
		0xb946: "釻",
		0xb947: "龔",
		0xb948: "勰",
		0xb949: "倻",
		0xb94a: "壒",
		0xb94b: "芿",
		0xb94c: "薷",
		0xb94d: "藿",
		0xb94e: "蒁",
		0xb94f: "豉",
		0xb950: "緣",
		0xb951: "臗",
		0xb952: "臏",
		0xb953: "滎",
		0xb954: "牖",
		0xb955: "瘂",
		0xb956: "蠲",
		0xb957: "鈹",
		0xb958: "顖",
		0xb959: "𩕄",
		0xb95a: "鼷",
		0xb95b: "齗",
		0xb95c: "嚳",
		0xb95d: "姧",
		0xb95e: "嬥",
		0xb95f: "彽",
		0xb960: "挵",
		0xb961: "洎",
		0xb962: "蕞",
		0xb963: "蕺",
		0xb964: "旰",
		0xb965: "梻",
		0xb966: "焠",
		0xb967: "禖",
		0xb968: "皂",
		0xb969: "皪",
		0xb96a: "眴",
		0xb96b: "裱",
		0xb96c: "簶",
		0xb96d: "蛽",
		0xb96e: "蜹",
		0xb96f: "蝲",
		0xb970: "蹰",
		0xb971: "頊",
		0xb972: "顓",
		0xb973: "顥",
		0xb974: "餼",
		0xb975: "鮞",
		0xb976: "鮬",
		0xb977: "鯎",
		0xb978: "鮸",
		0xb979: "鯘",
		0xb97a: "鰙",
		0xb97b: "鱮",
		0xb97c: "鄴",
		0xb97d: "璡",
		0xb97e: "磈",
		0xba21: "鄱",
		0xba22: "琚",
		0xba23: "艜",
		0xba24: "佺",
		0xba25: "偁",
		0xba26: "劻",
		0xba27: "噶",
		0xba28: "墪",
		0xba29: "埦",
		0xba2a: "嵆",
		0xba2b: "耼",
		0xba2c: "裒",
		0xba2d: "塤",
		0xba2e: "壎",
		0xba2f: "嚞",
		0xba30: "姒",
		0xba31: "姮",
		0xba75: "媧",
		0xba76: "幞",
		0xba77: "廆",
		0xba78: "弽",
		0xba79: "弴",
		0xba7a: "皙",
		0xba7b: "泔",
		0xba7c: "淛",
		0xba7d: "淝",
		0xba7e: "淄",
		0xbb21: "潙",
		0xbb22: "澶",
		0xbb23: "濊",
		0xbb24: "菡",
		0xbb25: "菪",
		0xbb26: "蒴",
		0xbb27: "蠆",
		0xbb28: "蘐",
		0xbb29: "鄯",
		0xbb2a: "适",
		0xbb2b: "忉",
		0xbb2c: "敔",
		0xbb2d: "鼂",
		0xbb2e: "昀",
		0xbb2f: "枲",
		0xbb30: "栱",
		0xbb31: "栝",
		0xbb32: "棅",
		0xbb33: "櫆",
		0xbb34: "烤",
		0xbb35: "犍",
		0xbb36: "珅",
		0xbb37: "玢",
		0xbb38: "珣",
		0xbb39: "琰",
		0xbb3a: "琫",
		0xbb3b: "瑀",
		0xbb3c: "瑄",
		0xbb3d: "瑒",
		0xbb3e: "瑭",
		0xbb3f: "瑫",
		0xbb40: "璵",
		0xbb41: "璩",
		0xbb42: "璿",
		0xbb43: "瓚",
		0xbb44: "癋",
		0xbb45: "磤",
		0xbb46: "竽",
		0xbb47: "筲",
		0xbb48: "糕",
		0xbb49: "紇",
		0xbb4a: "縑",
		0xbb4b: "羿",
		0xbb4c: "翺",
		0xbb4d: "詵",
		0xbb4e: "𨏍",
		0xbb4f: "銙",
		0xbb50: "錑",
		0xbb51: "錕",
		0xbb52: "鍱",
		0xbb53: "𨫤",
		0xbb54: "閦",
		0xbb55: "闐",
		0xbb57: "靮",
		0xbb58: "韴",
		0xbb59: "歆",
		0xbb5a: "頫",
		0xbb5b: "顒",
		0xbb5c: "顗",
		0xbb5d: "餛",
		0xbb5e: "餺",
		0xbb5f: "魹",
		0xbb60: "鷟",
		0xbb61: "毱",
		0xbb62: "﨟",
		0xbb63: "㐂",
		0xbb64: "鶡",
		0xbb65: "鸕",
		0xbb66: "鼐",
		0xbb67: "酈",
		0xbb68: "睟",
		0xbb69: "鹿子", /* FIXME: ⿰鹿子 (waseda just writes the characters separated */
		0xbb6a: "媞",
		0xbb6b: "彤",
		0xbb6c: "淩", /* goo uses 凌 but the image font has 氵 */
		0xbb6d: "葳",
		0xbb6e: "昫",
		0xbb6f: "簏",
		0xbb70: "騃",
		0xbb71: "輶",
		0xbb72: "莘",
		0xbb73: "摭",
		0xbb74: "茲",
		0xbb75: "咜",
		0xbb76: "晫",
		0xbb77: "昪",
		0xbb78: "枓",
		0xbb79: "翃",
		0xbb7a: "艠",
		0xbb7b: "酛",
		0xbb7c: "禛",
		0xbb7d: "¯",
		0xbc21: "▤",
		0xbc23: "♣",
		0xbc24: "♥",
		0xbc25: "♠",
		0xbc26: "♦",
		0xbc2a: "♮",
		0xbc2b: "゛",
		0xbc2c: "∘",
		0xbc32: "〳",
		0xbc33: "〵",
		0xbc36: "ℊ",
		0xbc37: "ς",
		0xbc3a: "〽",
		0xbc3b: "敱",
		0xbc3c: "哯",
		0xbc3d: "𠺕",
		0xbc3e: "擌",
		0xbc3f: "枛",
		0xbc40: "熮",
		0xbc41: "瓼",
		0xbc42: "𤸎",
		0xbc43: "癁",
		0xbc44: "瞙",
		0xbc45: "矪",
		0xbc46: "窼",
		0xbc47: "𥶡",
		0xbc48: "𥻨",
		0xbc49: "𦀌",
		0xbc4a: "縨",
		0xbc4b: "縬",
		0xbc4c: "繀",
		0xbc4d: "膲",
		0xbc4e: "𦨖",
		0xbc4f: "𦨞",
		0xbc50: "𦬇",
		0xbc51: "蚇",
		0xbc52: "蜏",
		0xbc53: "䗳",
		0xbc54: "䘣",
		0xbc55: "蹳",
		0xbc56: "鐁",
		0xbc57: "鞖",
		0xbc58: "䪊",
		0xbc59: "鯳",
		0xbc5a: "鳹",
		0xbc5b: "𪀚",
		0xbc5c: "䳑",
		0xbc5d: "鸍",
		0xbc5e: "璙",
		0xbc5f: "秇",
		0xbc60: "羕",
		0xbc62: "皶",
		0xbc64: "𥇒",
		0xbc65: "嚧",
		0xbc66: "坅",
		0xbc67: "蘸",
		0xbc68: "𤭯",
		0xbc69: "牫",
		0xbc6a: "矠",
		0xbc6b: "硾",
		0xbc6c: "𥿠",
		0xbc6d: "𥄢",
		0xbc6e: "𪃹",
		0xbc6f: "韡",
		0xbc70: "异",
		0xbc72: "麀",
		0xbc73: "𥫱",
		0xbc74: "朮",
		0xbc75: "𧐐",
		0xbd5c: "慒",
		0xbd5d: "𪙉",
		0xbd5e: "毈",
		0xbd5f: "薼",
		0xbd60: "𡣞",
		0xbd61: "𤣥",
		0xbd62: "𨉷",
		0xbd63: "𡻕",
		0xbd64: "荗",
		0xbd65: "麬",
		0xbd66: "㸅",
		0xbd67: "𫛉",
		0xbd68: "磂",
		0xbd69: "坶",
		0xbd6a: "䗈",
		0xbd6b: "檫",
		0xbd6c: "欛",
		0xbd6d: "欙",
		0xbd6e: "殭",
		0xbd6f: "甗",
		0xbd70: "軑",
		0xbd71: "輀",
		0xbd72: "輭",
		0xbd73: "轘",
		0xbd74: "菑",
		0xbd75: "葖",
		0xbd76: "蓇",
		0xbd77: "蘘",
		0xbd78: "杈",
		0xbd79: "焮",
		0xbd7a: "昰",
		0xbd7b: "尟",
		0xbd7c: "賙",
		0xbd7d: "璫",
		0xbd7e: "璠",
		0xbe21: "疢",
		0xbe22: "瞤",
		0xbe23: "矬",
		0xbe24: "矻",
		0xbe25: "磠",
		0xbe26: "穭",
		0xbe27: "窅",
		0xbe28: "笭",
		0xbe29: "簋",
		0xbe2a: "簠",
		0xbe2b: "耦",
		0xbe2c: "蝘",
		0xbe2d: "豨",
		0xbe2e: "飣",
		0xbe2f: "餖",
		0xbe30: "膆",
		0xbe31: "臛",
		0xbe32: "欬",
		0xbe33: "羖",
		0xbe34: "疿",
		0xbe35: "蝱",
		0xbe36: "嚲",
		0xbe37: "匜",
		0xbe38: "刵",
		0xbe39: "剉",
		0xbe3a: "箚",
		0xbe3b: "𤴡",
		0xbe3d: "伋",
		0xbe3e: "睠",
		0xbe3f: "僄",
		0xbe40: "儵",
		0xbe41: "煠",
		0xbe42: "熅",
		0xbe43: "熛",
		0xbe44: "僶",
		0xbe45: "隤",
		0xbe46: "扆",
		0xbe47: "璆",
		0xbe48: "攩",
		0xbe49: "洿",
		0xbe4a: "涑",
		0xbe4b: "攙",
		0xbe4c: "瓈",
		0xbe4d: "罝",
		0xbe4e: "盬",
		0xbe4f: "鈇",
		0xbe50: "鉧",
		0xbe51: "銍",
		0xbe52: "淯",
		0xbe53: "湉",
		0xbe54: "滃",
		0xbe55: "噠",
		0xbe56: "嘽",
		0xbe57: "嚬",
		0xbe58: "鋙",
		0xbe59: "鎛",
		0xbe5a: "鏽",
		0xbe5b: "矰",
		0xbe5c: "灃",
		0xbe5d: "忼",
		0xbe5e: "怵",
		0xbe5f: "怳",
		0xbe60: "惛",
		0xbe61: "愐",
		0xbe62: "嚚",
		0xbe63: "篅",
		0xbe64: "慠",
		0xbe65: "籑",
		0xbe66: "籰",
		0xbe67: "臲",
		0xbe68: "稃",
		0xbe69: "惸",
		0xbe6a: "慠",
		0xbe6b: "庤",
		0xbe6c: "閟",
		0xbe6d: "玁",
		0xbe6e: "餗",
		0xbe6f: "餧",
		0xbe70: "䭔",
		0xbe71: "餻",
		0xbe72: "饆",
		0xbe73: "皝",
		0xbe74: "鵟",
		0xbe75: "皻",
		0xbe76: "堄",
		0xbe77: "埤",
		0xbe78: "塼",
		0xbe79: "饞",
		0xbe7a: "饠",
		0xbe7b: "姁",
		0xbe7c: "姞",
		0xbe7d: "媢",
		0xbe7e: "媿",
		0xbf21: "孼",
		0xbf22: "籹",
		0xbf23: "粔",
		0xbf24: "顇",
		0xbf25: "顦",
		0xbf26: "蜋",
		0xbf27: "蜐",
		0xbf28: "蜺",
		0xbf29: "苾",
		0xbf2a: "莕",
		0xbf2b: "葅",
		0xbf2c: "蓏",
		0xbf2d: "薁",
		0xbf2e: "薟",
		0xbf2f: "藊",
		0xbf30: "綷",
		0xbf31: "纑",
		0xbf32: "駰",
		0xbf33: "朳",
		0xbf34: "杇",
		0xbf35: "蝥",
		0xbf36: "螠",
		0xbf37: "蠔",
		0xbf38: "蠭",
		0xbf39: "醨",
		0xbf3a: "醼",
		0xbf3b: "耷",
		0xbf3c: "掽",
		0xbf3d: "枒",
		0xbf3e: "柹",
		0xbf3f: "杴",
		0xbf40: "杻",
		0xbf41: "棬",
		0xbf42: "躃",
		0xbf43: "搩",
		0xbf44: "摽",
		0xbf45: "艴",
		0xbf46: "穀",
		0xbf47: "榰", /* FIXME: 「榰」different image character but same kanji */
		0xbf48: "樝",
		0xbf49: "檛",
		0xbf4a: "檉",
		0xbf4b: "靿",
		0xbf4c: "鞚",
		0xbf4d: "韉",
		0xbf4e: "趯",
		0xbf4f: "虁",
		0xbf50: "蓀",
		0xbf51: "确",
		0xbf52: "拼",
		0xbf53: "蘡",
		0xbf54: "霅",
		0xbf55: "鮊",
		0xbf56: "誮",
		0xbf57: "土",
		0xbf58: "棻",
		0xbf59: "碭",
		0xbf5a: "獼",
		0xc041: "鯝",
		0xc042: "鱟",
		0xc043: "鱵",
		0xc044: "鬅",
		0xc045: "鬌",
		0xc046: "扃",
		0xc047: "橐",
		0xc048: "鱏",
		0xc049: "氳",
		0xc04a: "罾",
		0xc04b: "攲",
		0xc04c: "巹",
		0xc04d: "齁",
		0xc04e: "呫",
		0xc04f: "麯",
		0xc050: "魣",
		0xc051: "𦯶",
		0xc052: "蠐",
		0xc053: "蓽",
		0xc054: "柃",
		0xc055: "𧏛",
		0xc056: "髐",
		0xc057: "𨗈",
		0xc058: "笇",
		0xc059: "匾",
		0xc05a: "蒾",
		0xc05b: "鴗",
		0xc05c: "偟",
		0xc05d: "藋",
		0xc05e: "甆",
		0xc05f: "穇",
		0xc060: "蜟",
		0xc061: "壚",
		0xc063: "牁",
		0xc064: "胘",
		0xc065: "黮",
		0xc066: "婥",
		0xc068: "止",
		0xc069: "嗢",
		0xc06a: "鳭",
		0xc06b: "麥",
		0xc06c: "鶬",
		0xc06d: "虯",
		0xc06e: "庪",
		0xc06f: "秭",
		0xc070: "岏",
		0xc071: "⻞",
		0xc072: "阝",
		0xc073: "槩",
		0xc074: "毿",
		0xc075: "灔",
		0xc076: "么",
		0xc077: "鼺",
		0xc078: "蠁",
		0xc079: "麨",
		0xc07a: "碰",
		0xc07b: "俰",
		0xc07c: "筟",
		0xc07d: "鱪",
		0xc07e: "仐",
		0xc121: "牱",
		0xc123: "犎",
		0xc124: "猲",
		0xc125: "袽",
		0xc126: "蕽",
		0xc127: "桵",
		0xc128: "椂",
		0xc129: "傔",
		0xc12a: "儃",
		0xc12b: "扡",
		0xc12c: "挃",
		0xc12d: "詤",
		0xc12e: "誷",
		0xc12f: "鮾",
		0xc130: "鱐",
		0xc131: "箯",
		0xc132: "荇",
		0xc133: "蓎",
		0xc134: "茝",
		0xc135: "檨",
		0xc136: "蕡",
		0xc137: "醶",
		0xc138: "簄",
		0xc139: "觘",
		0xc13a: "鑯",
		0xc13b: "𣑥",
		0xc13c: "猍",
		0xc13d: "葼",
		0xc13e: "箶",
		0xc13f: "粶",
		0xc140: "迱",
		0xc141: "髩",
		0xc142: "橒",
		0xc143: "龗",
		0xc144: "籡",
		0xc145: "粏",
		0xc146: "蚸",
		0xc147: "螇",
		0xc148: "鞺",
		0xc149: "鰖",
		0xc14a: "鱰",
		0xc14b: "鴲",
		0xc14c: "鷀",
		0xc14d: "彇",
		0xc14e: "鋐",
		0xc14f: "𡱖",
		0xc150: "笧",
		0xc151: "篗",
		0xc152: "糄",
		0xc153: "𫒒",
		0xc154: "鐴",
		0xc155: "篔",
		0xc156: "舃",
		0xc157: "忩",
		0xc158: "𩺊",
		0xc159: "芸", /* FIXME: Should have split radical */
		0xc15a: "簳",
		0xc15b: "𤭖",
		0xc15c: "蔲",
		0xc15d: "竈",
		0xc15e: "鉏",
		0xc15f: "尩",
		0xc160: "邌",
		0xc161: "鮧",
		0xc162: "鱁",
		0xc163: "鱛",
		0xc164: "鬂",
		0xc165: "酤",
		0xc166: "樏",
		0xc167: "襅",
		0xc168: "蒅",
		0xc169: "躮",
		0xc16a: "鮲",
		0xc16b: "鰘",
		0xc16c: "鵇",
		0xc16d: "嚈",
		0xc16e: "憍",
		0xc16f: "𣪘",
		0xc170: "璱",
		0xc171: "褹",
		0xc172: "緂",
		0xc173: "鬠",
		0xc174: "鐧",
		0xc175: "㝢",
		0xc176: "洀",
		0xc177: "襀",
		0xc178: "嚩",
		0xc179: "挍",
		0xc17a: "𩊱",
		0xc17b: "妋",
		0xc17c: "熇",
		0xc17d: "戭",
		0xc17e: "煑",
		0xc221: "顊",
		0xc222: "斲",
		0xc223: "鄽",
		0xc224: "柲",
		0xc225: "齝",
		0xc226: "鯯",
		0xc227: "㮶",
		0xc228: "檝",
		0xc229: "蚉",
		0xc22a: "蛁",
		0xc22b: "蟟",
		0xc22c: "洦",
		0xc22d: "孁",
		0xc22e: "𡑮",
		0xc22f: "鏱",
		0xc230: "裛",
		0xc231: "礜",
		0xc232: "𣑊",
		0xc233: "籭",
		0xc234: "儞",
		0xc235: "頞",
		0xc236: "㒵",
		0xc237: "𩅧",
		0xc238: "魶",
		0xc239: "鷧",
		0xc23a: "瞔",
		0xc23b: "橖",
		0xc23c: "紞",
		0xc23d: "韝",
		0xc23e: "弣",
		0xc23f: "芺",
		0xc240: "惔",
		0xc241: "唽",
		0xc327: "榺",
		0xc328: "笯",
		0xc329: "砑",
		0xc32a: "畾",
		0xc32b: "灩",
		0xc32c: "埸",
		0xc32d: "釱",
		0xc32e: "炗",
		0xc32f: "鬜",
		0xc330: "鯽",
		0xc331: "癤",
		0xc332: "梂",
		0xc333: "蔤",
		0xc334: "鋂",
		0xc335: "壍",
		0xc336: "痟",
		0xc337: "齵",
		0xc338: "鸜",
		0xc339: "泬",
		0xc33a: "釽",
		0xc33b: "籗",
		0xc33c: "楲",
		0xc33d: "窬",
		0xc33e: "貒",
		0xc33f: "悞",
		0xc340: "尰",
		0xc341: "巑",
		0xc342: "葈",
		0xc343: "藦",
		0xc344: "腅",
		0xc345: "糫",
		0xc346: "蟭",
		0xc347: "曌",
		0xc348: "孙",
		0xc349: "乐",
		0xc34a: "车",
		0xc34b: "产",
		0xc34c: "电",
		0xc34d: "榰", /* FIXME: 「榰」different image character but same kanji */
		0xc34e: "蜐",
		0xc350: "⼝",
		0xc351: "⼟",
		0xc352: "⼥",
		0xc353: "⼭",
		0xc355: "⼱",
		0xc356: "⼸",
		0xc357: "⽇",
		0xc358: "⺝",
		0xc359: "⽊",
		0xc35a: "⽕",
		0xc35b: "𤣩",
		0xc35c: "⽬",
		0xc35d: "⽯",
		0xc35e: "⽲",
		0xc35f: "⽶",
		0xc360: "糹",
		0xc362: "⽿",
		0xc364: "⾍",
		0xc365: "訁",
		0xc367: "⾙",
		0xc368: "⾞",
		0xc369: "⾣",
		0xc36a: "釒",
		0xc36b: "⻗",
		0xc36c: "⾰",
		0xc36d: "⿂",
		0xc36e: "❶",
		0xc36f: "❷",
		0xc370: "❸",
		0xc371: "❹",
		0xc372: "❺",
		0xc373: "\n①",
		0xc374: "\n②",
		0xc375: "\n③",
		0xc376: "\n④",
		0xc377: "\n⑤",
		0xc378: "\n⑥",
		0xc379: "\n⑦",
		0xc37a: "\n⑧",
		0xc37b: "\n⑨",
		0xc37c: "\n⑩",
		0xc37d: "\n⑪",
		0xc37e: "\n⑫",
		0xc421: "\n⑬",
		0xc422: "\n⑭",
		0xc423: "\n⑮",
		0xc424: "\n⑯",
		0xc425: "\n⑰",
		0xc426: "\n⑱",
		0xc427: "\n⑲",
		0xc428: "\n⑳",
		0xc429: "\n㉑",
		0xc42a: "\n㉒",
		0xc42b: "\n㉓",
		0xc42c: "\n㉔",
		0xc42d: "\n㉕",
		0xc430: "𝄉",
		0xc431: "\n㊀",
		0xc432: "\n㊁",
		0xc433: "\n㊂",
		0xc434: "\n㊃",
		0xc435: "\n㊄",
		0xc437: "\n㊀",
		0xc438: "\n㊁",
		0xc439: "\n㊂",
		0xc43a: "\n㊃",
		0xc43b: "\n㊄",
		0xc43c: "\n㊅",
		0xc43d: "\n㊆",
		0xc43e: "\n㊇",
		0xc43f: "\n㊈",
		0xc440: "\n㉖",
		0xc441: "\n㉗",
		0xc442: "\n㉘",
		0xc443: "\n㉙",
		0xc444: "\n㉚",
		0xc445: "\n㉛",
		0xc446: "\n㉜",
		0xc447: "\n㉜",
		0xc448: "\n㉝",
		0xc449: "\n㉞",
		0xc44a: "\n㉟",
		0xc44b: "䷝",
		0xc44c: "䷲",
		0xc44d: "䷸",
		0xc44e: "䷜",
		0xc44f: "䷳",
		0xc450: "䷁",
		0xc451: "乀",
		0xc453: "∨",
		0xc454: "𝄐",
		0xc455: "㉆",
		0xc457: "𝒜",
		0xc458: "▧",
		0xc459: "⿴",
		0xc45a: "भर", /* FIXME: 勃嚕唵/bhrūṃ  in sanskrit */
		0xc45c: "ϒ", /* ♈︎ is also similar (for 雁点)*/
		0xc461: "\n㈥",
		0xc462: "\n㊅",
		0xc463: "", /* FIXME: Camera icon. indicates an image. not needed. */
		0xc464: "", /* FIXME: Camera icon. indicates an image. not needed. */
		0xc465: "", /* FIXME: ♪. indicates audio. not needed */
		0xc466: "", /* FIXME: Paintbrush icon. indicates stroke diagram. not needed. */
		0xc468: "《色》", /* FIXME: Ink pot icon. indicates ex. color. */
		0xc46d: "あ･",
		0xc46e: "き･",
		0xc46f: "こ|",
		0xc470: "じ|",
		0xc471: "ず|",
		0xc472: "せ･",
		0xc473: "た･",
		0xc474: "の|",
		0xc475: "は･",
		0xc476: "も|",
		0xc479: "ß",
		0xc47a: "Æ",
		0xc47b: "æ",
		0xc47c: "œ",
		0xc47d: "Ⅰ",
		0xc47e: "Ⅱ",
		0xc526: "Ⅴ",
		0xc527: "ɔ",
		0xc56b: "ʃ",
		0xc56c: "ɑ",
		0xc56d: "ː",
		0xc56e: "ɡ",
		0xc56f: "ŋ",
		0xc570: "ʒ",
		0xc571: "ø",
		0xc572: "ɲ",
		0xc573: "ɟ",
		0xc579: "Φ",
		0xc57a: "⇄",
		0xc57b: "ø",
		0xc57c: "∛",
		0xc57d: "ⁿ√",
		0xc57e: "√",
		0xc621: "㏋",
		0xc623: "∓",
		0xc624: "√2",
		0xc625: "√𝑎",
		0xc626: "℥",
		0xc627: "⫅",
		0xc628: "Δ",
		0xc629: "ʌ",
		0xc62a: "Π",
		0xc62b: "Σ",
		0xc62c: "Φ",
		0xc62d: "Ω",
		0xc62e: "α",
		0xc62f: "β",
		0xc630: "γ",
		0xc631: "θ",
		0xc632: "λ",
		0xc633: "μ",
		0xc634: "π",
		0xc635: "ϕ",
		0xc636: "ɯ̈",
	}
}
