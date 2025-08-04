package common

type LanguageCode string

const (
	// Afrikaans
	LanguageCodeAF LanguageCode = "af"
	// Akan
	LanguageCodeAK LanguageCode = "ak"
	// Amharic
	LanguageCodeAM LanguageCode = "am"
	// Arabic
	LanguageCodeAR LanguageCode = "ar"
	// Assamese
	LanguageCodeAS LanguageCode = "as"
	// Azerbaijani
	LanguageCodeAZ LanguageCode = "az"
	// Belarusian
	LanguageCodeBE LanguageCode = "be"
	// Bulgarian
	LanguageCodeBG LanguageCode = "bg"
	// Bambara
	LanguageCodeBM LanguageCode = "bm"
	// Bangla
	LanguageCodeBN LanguageCode = "bn"
	// Tibetan
	LanguageCodeBO LanguageCode = "bo"
	// Breton
	LanguageCodeBR LanguageCode = "br"
	// Bosnian
	LanguageCodeBS LanguageCode = "bs"
	// Catalan
	LanguageCodeCA LanguageCode = "ca"
	// Chechen
	LanguageCodeCE LanguageCode = "ce"
	// Corsican
	LanguageCodeCO LanguageCode = "co"
	// Czech
	LanguageCodeCS LanguageCode = "cs"
	// Church Slavic
	LanguageCodeCU LanguageCode = "cu"
	// Welsh
	LanguageCodeCY LanguageCode = "cy"
	// Danish
	LanguageCodeDA LanguageCode = "da"
	// German
	LanguageCodeDE LanguageCode = "de"
	// Austrian German
	LanguageCodeDE_AT LanguageCode = "de_AT"
	// Swiss High German
	LanguageCodeDE_CH LanguageCode = "de_CH"
	// Dzongkha
	LanguageCodeDZ LanguageCode = "dz"
	// Ewe
	LanguageCodeEE LanguageCode = "ee"
	// Greek
	LanguageCodeEL LanguageCode = "el"
	// English
	LanguageCodeEN LanguageCode = "en"
	// Australian English
	LanguageCodeEN_AU LanguageCode = "en_AU"
	// Canadian English
	LanguageCodeEN_CA LanguageCode = "en_CA"
	// British English
	LanguageCodeEN_GB LanguageCode = "en_GB"
	// American English
	LanguageCodeEN_US LanguageCode = "en_US"
	// Esperanto
	LanguageCodeEO LanguageCode = "eo"
	// Spanish
	LanguageCodeES LanguageCode = "es"
	// European Spanish
	LanguageCodeES_ES LanguageCode = "es_ES"
	// Mexican Spanish
	LanguageCodeES_MX LanguageCode = "es_MX"
	// Estonian
	LanguageCodeET LanguageCode = "et"
	// Basque
	LanguageCodeEU LanguageCode = "eu"
	// Persian
	LanguageCodeFA LanguageCode = "fa"
	// Dari
	LanguageCodeFA_AF LanguageCode = "fa_AF"
	// Fulah
	LanguageCodeFF LanguageCode = "ff"
	// Finnish
	LanguageCodeFI LanguageCode = "fi"
	// Faroese
	LanguageCodeFO LanguageCode = "fo"
	// French
	LanguageCodeFR LanguageCode = "fr"
	// Canadian French
	LanguageCodeFR_CA LanguageCode = "fr_CA"
	// Swiss French
	LanguageCodeFR_CH LanguageCode = "fr_CH"
	// Western Frisian
	LanguageCodeFY LanguageCode = "fy"
	// Irish
	LanguageCodeGA LanguageCode = "ga"
	// Scottish Gaelic
	LanguageCodeGD LanguageCode = "gd"
	// Galician
	LanguageCodeGL LanguageCode = "gl"
	// Gujarati
	LanguageCodeGU LanguageCode = "gu"
	// Manx
	LanguageCodeGV LanguageCode = "gv"
	// Hausa
	LanguageCodeHA LanguageCode = "ha"
	// Hebrew
	LanguageCodeHE LanguageCode = "he"
	// Hindi
	LanguageCodeHI LanguageCode = "hi"
	// Croatian
	LanguageCodeHR LanguageCode = "hr"
	// Haitian Creole
	LanguageCodeHT LanguageCode = "ht"
	// Hungarian
	LanguageCodeHU LanguageCode = "hu"
	// Armenian
	LanguageCodeHY LanguageCode = "hy"
	// Interlingua
	LanguageCodeIA LanguageCode = "ia"
	// Indonesian
	LanguageCodeID LanguageCode = "id"
	// Igbo
	LanguageCodeIG LanguageCode = "ig"
	// Sichuan Yi
	LanguageCodeII LanguageCode = "ii"
	// Icelandic
	LanguageCodeIS LanguageCode = "is"
	// Italian
	LanguageCodeIT LanguageCode = "it"
	// Japanese
	LanguageCodeJA LanguageCode = "ja"
	// Javanese
	LanguageCodeJV LanguageCode = "jv"
	// Georgian
	LanguageCodeKA LanguageCode = "ka"
	// Kikuyu
	LanguageCodeKI LanguageCode = "ki"
	// Kazakh
	LanguageCodeKK LanguageCode = "kk"
	// Kalaallisut
	LanguageCodeKL LanguageCode = "kl"
	// Khmer
	LanguageCodeKM LanguageCode = "km"
	// Kannada
	LanguageCodeKN LanguageCode = "kn"
	// Korean
	LanguageCodeKO LanguageCode = "ko"
	// Kashmiri
	LanguageCodeKS LanguageCode = "ks"
	// Kurdish
	LanguageCodeKU LanguageCode = "ku"
	// Cornish
	LanguageCodeKW LanguageCode = "kw"
	// Kyrgyz
	LanguageCodeKY LanguageCode = "ky"
	// Latin
	LanguageCodeLA LanguageCode = "la"
	// Luxembourgish
	LanguageCodeLB LanguageCode = "lb"
	// Ganda
	LanguageCodeLG LanguageCode = "lg"
	// Lingala
	LanguageCodeLN LanguageCode = "ln"
	// Lao
	LanguageCodeLO LanguageCode = "lo"
	// Lithuanian
	LanguageCodeLT LanguageCode = "lt"
	// Luba-Katanga
	LanguageCodeLU LanguageCode = "lu"
	// Latvian
	LanguageCodeLV LanguageCode = "lv"
	// Malagasy
	LanguageCodeMG LanguageCode = "mg"
	// Maori
	LanguageCodeMI LanguageCode = "mi"
	// Macedonian
	LanguageCodeMK LanguageCode = "mk"
	// Malayalam
	LanguageCodeML LanguageCode = "ml"
	// Mongolian
	LanguageCodeMN LanguageCode = "mn"
	// Marathi
	LanguageCodeMR LanguageCode = "mr"
	// Malay
	LanguageCodeMS LanguageCode = "ms"
	// Maltese
	LanguageCodeMT LanguageCode = "mt"
	// Burmese
	LanguageCodeMY LanguageCode = "my"
	// Norwegian Bokmål
	LanguageCodeNB LanguageCode = "nb"
	// North Ndebele
	LanguageCodeND LanguageCode = "nd"
	// Nepali
	LanguageCodeNE LanguageCode = "ne"
	// Dutch
	LanguageCodeNL LanguageCode = "nl"
	// Flemish
	LanguageCodeNL_BE LanguageCode = "nl_BE"
	// Norwegian Nynorsk
	LanguageCodeNN LanguageCode = "nn"
	// Nyanja
	LanguageCodeNY LanguageCode = "ny"
	// Oromo
	LanguageCodeOM LanguageCode = "om"
	// Odia
	LanguageCodeOR LanguageCode = "or"
	// Ossetic
	LanguageCodeOS LanguageCode = "os"
	// Punjabi
	LanguageCodePA LanguageCode = "pa"
	// Polish
	LanguageCodePL LanguageCode = "pl"
	// Pashto
	LanguageCodePS LanguageCode = "ps"
	// Portuguese
	LanguageCodePT LanguageCode = "pt"
	// Brazilian Portuguese
	LanguageCodePT_BR LanguageCode = "pt_BR"
	// European Portuguese
	LanguageCodePT_PT LanguageCode = "pt_PT"
	// Quechua
	LanguageCodeQU LanguageCode = "qu"
	// Romansh
	LanguageCodeRM LanguageCode = "rm"
	// Rundi
	LanguageCodeRN LanguageCode = "rn"
	// Romanian
	LanguageCodeRO LanguageCode = "ro"
	// Moldavian
	LanguageCodeRO_MD LanguageCode = "ro_MD"
	// Russian
	LanguageCodeRU LanguageCode = "ru"
	// Kinyarwanda
	LanguageCodeRW LanguageCode = "rw"
	// Sanskrit
	LanguageCodeSA LanguageCode = "sa"
	// Sindhi
	LanguageCodeSD LanguageCode = "sd"
	// Northern Sami
	LanguageCodeSE LanguageCode = "se"
	// Sango
	LanguageCodeSG LanguageCode = "sg"
	// Sinhala
	LanguageCodeSI LanguageCode = "si"
	// Slovak
	LanguageCodeSK LanguageCode = "sk"
	// Slovenian
	LanguageCodeSL LanguageCode = "sl"
	// Samoan
	LanguageCodeSM LanguageCode = "sm"
	// Shona
	LanguageCodeSN LanguageCode = "sn"
	// Somali
	LanguageCodeSO LanguageCode = "so"
	// Albanian
	LanguageCodeSQ LanguageCode = "sq"
	// Serbian
	LanguageCodeSR LanguageCode = "sr"
	// Southern Sotho
	LanguageCodeST LanguageCode = "st"
	// Sundanese
	LanguageCodeSU LanguageCode = "su"
	// Swedish
	LanguageCodeSV LanguageCode = "sv"
	// Swahili
	LanguageCodeSW LanguageCode = "sw"
	// Congo Swahili
	LanguageCodeSW_CD LanguageCode = "sw_CD"
	// Tamil
	LanguageCodeTA LanguageCode = "ta"
	// Telugu
	LanguageCodeTE LanguageCode = "te"
	// Tajik
	LanguageCodeTG LanguageCode = "tg"
	// Thai
	LanguageCodeTH LanguageCode = "th"
	// Tigrinya
	LanguageCodeTI LanguageCode = "ti"
	// Turkmen
	LanguageCodeTK LanguageCode = "tk"
	// Tongan
	LanguageCodeTO LanguageCode = "to"
	// Turkish
	LanguageCodeTR LanguageCode = "tr"
	// Tatar
	LanguageCodeTT LanguageCode = "tt"
	// Uyghur
	LanguageCodeUG LanguageCode = "ug"
	// Ukrainian
	LanguageCodeUK LanguageCode = "uk"
	// Urdu
	LanguageCodeUR LanguageCode = "ur"
	// Uzbek
	LanguageCodeUZ LanguageCode = "uz"
	// Vietnamese
	LanguageCodeVI LanguageCode = "vi"
	// Volapük
	LanguageCodeVO LanguageCode = "vo"
	// Wolof
	LanguageCodeWO LanguageCode = "wo"
	// Xhosa
	LanguageCodeXH LanguageCode = "xh"
	// Yiddish
	LanguageCodeYI LanguageCode = "yi"
	// Yoruba
	LanguageCodeYO LanguageCode = "yo"
	// Chinese
	LanguageCodeZH LanguageCode = "zh"
	// Simplified Chinese
	LanguageCodeZH_HANS LanguageCode = "zh_Hans"
	// Traditional Chinese
	LanguageCodeZH_HANT LanguageCode = "zh_Hant"
	// Zulu
	LanguageCodeZU LanguageCode = "zu"
)