package common

type CurrencyCode string

const (
	// United Arab Emirates dirham
	CurrencyCodeAED CurrencyCode = "AED"
	// Afghan afghani
	CurrencyCodeAFN CurrencyCode = "AFN"
	// Albanian lek
	CurrencyCodeALL CurrencyCode = "ALL"
	// Armenian dram
	CurrencyCodeAMD CurrencyCode = "AMD"
	// Netherlands Antillean guilder
	CurrencyCodeANG CurrencyCode = "ANG"
	// Angolan kwanza
	CurrencyCodeAOA CurrencyCode = "AOA"
	// Argentine peso
	CurrencyCodeARS CurrencyCode = "ARS"
	// Australian dollar
	CurrencyCodeAUD CurrencyCode = "AUD"
	// Aruban florin
	CurrencyCodeAWG CurrencyCode = "AWG"
	// Azerbaijani manat
	CurrencyCodeAZN CurrencyCode = "AZN"
	// Bosnia and Herzegovina convertible mark
	CurrencyCodeBAM CurrencyCode = "BAM"
	// Barbados dollar
	CurrencyCodeBBD CurrencyCode = "BBD"
	// Bangladeshi taka
	CurrencyCodeBDT CurrencyCode = "BDT"
	// Bulgarian lev
	CurrencyCodeBGN CurrencyCode = "BGN"
	// Bahraini dinar
	CurrencyCodeBHD CurrencyCode = "BHD"
	// Burundian franc
	CurrencyCodeBIF CurrencyCode = "BIF"
	// Bermudian dollar
	CurrencyCodeBMD CurrencyCode = "BMD"
	// Brunei dollar
	CurrencyCodeBND CurrencyCode = "BND"
	// Boliviano
	CurrencyCodeBOB CurrencyCode = "BOB"
	// Brazilian real
	CurrencyCodeBRL CurrencyCode = "BRL"
	// Bahamian dollar
	CurrencyCodeBSD CurrencyCode = "BSD"
	// Bhutanese ngultrum
	CurrencyCodeBTN CurrencyCode = "BTN"
	// Botswana pula
	CurrencyCodeBWP CurrencyCode = "BWP"
	// Belarusian ruble
	CurrencyCodeBYN CurrencyCode = "BYN"
	// Belize dollar
	CurrencyCodeBZD CurrencyCode = "BZD"
	// Canadian dollar
	CurrencyCodeCAD CurrencyCode = "CAD"
	// Congolese franc
	CurrencyCodeCDF CurrencyCode = "CDF"
	// Swiss franc
	CurrencyCodeCHF CurrencyCode = "CHF"
	// Chilean peso
	CurrencyCodeCLP CurrencyCode = "CLP"
	// Renminbi (Chinese) yuan
	CurrencyCodeCNY CurrencyCode = "CNY"
	// Colombian peso
	CurrencyCodeCOP CurrencyCode = "COP"
	// Costa Rican colon
	CurrencyCodeCRC CurrencyCode = "CRC"
	// Cuban convertible peso
	CurrencyCodeCUC CurrencyCode = "CUC"
	// Cuban peso
	CurrencyCodeCUP CurrencyCode = "CUP"
	// Cape Verde escudo
	CurrencyCodeCVE CurrencyCode = "CVE"
	// Czech koruna
	CurrencyCodeCZK CurrencyCode = "CZK"
	// Djiboutian franc
	CurrencyCodeDJF CurrencyCode = "DJF"
	// Danish krone
	CurrencyCodeDKK CurrencyCode = "DKK"
	// Dominican peso
	CurrencyCodeDOP CurrencyCode = "DOP"
	// Algerian dinar
	CurrencyCodeDZD CurrencyCode = "DZD"
	// Egyptian pound
	CurrencyCodeEGP CurrencyCode = "EGP"
	// Eritrean nakfa
	CurrencyCodeERN CurrencyCode = "ERN"
	// Ethiopian birr
	CurrencyCodeETB CurrencyCode = "ETB"
	// Euro
	CurrencyCodeEUR CurrencyCode = "EUR"
	// Fiji dollar
	CurrencyCodeFJD CurrencyCode = "FJD"
	// Falkland Islands pound
	CurrencyCodeFKP CurrencyCode = "FKP"
	// Pound sterling
	CurrencyCodeGBP CurrencyCode = "GBP"
	// Georgian lari
	CurrencyCodeGEL CurrencyCode = "GEL"
	// Ghanaian cedi
	CurrencyCodeGHS CurrencyCode = "GHS"
	// Gibraltar pound
	CurrencyCodeGIP CurrencyCode = "GIP"
	// Gambian dalasi
	CurrencyCodeGMD CurrencyCode = "GMD"
	// Guinean franc
	CurrencyCodeGNF CurrencyCode = "GNF"
	// Guatemalan quetzal
	CurrencyCodeGTQ CurrencyCode = "GTQ"
	// Guyanese dollar
	CurrencyCodeGYD CurrencyCode = "GYD"
	// Hong Kong dollar
	CurrencyCodeHKD CurrencyCode = "HKD"
	// Honduran lempira
	CurrencyCodeHNL CurrencyCode = "HNL"
	// Croatian kuna
	CurrencyCodeHRK CurrencyCode = "HRK"
	// Haitian gourde
	CurrencyCodeHTG CurrencyCode = "HTG"
	// Hungarian forint
	CurrencyCodeHUF CurrencyCode = "HUF"
	// Indonesian rupiah
	CurrencyCodeIDR CurrencyCode = "IDR"
	// Israeli new shekel
	CurrencyCodeILS CurrencyCode = "ILS"
	// Indian rupee
	CurrencyCodeINR CurrencyCode = "INR"
	// Iraqi dinar
	CurrencyCodeIQD CurrencyCode = "IQD"
	// Iranian rial
	CurrencyCodeIRR CurrencyCode = "IRR"
	// Icelandic króna
	CurrencyCodeISK CurrencyCode = "ISK"
	// Jamaican dollar
	CurrencyCodeJMD CurrencyCode = "JMD"
	// Jordanian dinar
	CurrencyCodeJOD CurrencyCode = "JOD"
	// Japanese yen
	CurrencyCodeJPY CurrencyCode = "JPY"
	// Kenyan shilling
	CurrencyCodeKES CurrencyCode = "KES"
	// Kyrgyzstani som
	CurrencyCodeKGS CurrencyCode = "KGS"
	// Cambodian riel
	CurrencyCodeKHR CurrencyCode = "KHR"
	// Comoro franc
	CurrencyCodeKMF CurrencyCode = "KMF"
	// North Korean won
	CurrencyCodeKPW CurrencyCode = "KPW"
	// South Korean won
	CurrencyCodeKRW CurrencyCode = "KRW"
	// Kuwaiti dinar
	CurrencyCodeKWD CurrencyCode = "KWD"
	// Cayman Islands dollar
	CurrencyCodeKYD CurrencyCode = "KYD"
	// Kazakhstani tenge
	CurrencyCodeKZT CurrencyCode = "KZT"
	// Lao kip
	CurrencyCodeLAK CurrencyCode = "LAK"
	// Lebanese pound
	CurrencyCodeLBP CurrencyCode = "LBP"
	// Sri Lankan rupee
	CurrencyCodeLKR CurrencyCode = "LKR"
	// Liberian dollar
	CurrencyCodeLRD CurrencyCode = "LRD"
	// Lesotho loti
	CurrencyCodeLSL CurrencyCode = "LSL"
	// Libyan dinar
	CurrencyCodeLYD CurrencyCode = "LYD"
	// Moroccan dirham
	CurrencyCodeMAD CurrencyCode = "MAD"
	// Moldovan leu
	CurrencyCodeMDL CurrencyCode = "MDL"
	// Malagasy ariary
	CurrencyCodeMGA CurrencyCode = "MGA"
	// Macedonian denar
	CurrencyCodeMKD CurrencyCode = "MKD"
	// Myanmar kyat
	CurrencyCodeMMK CurrencyCode = "MMK"
	// Mongolian tögrög
	CurrencyCodeMNT CurrencyCode = "MNT"
	// Macanese pataca
	CurrencyCodeMOP CurrencyCode = "MOP"
	// Mauritanian ouguiya
	CurrencyCodeMRU CurrencyCode = "MRU"
	// Mauritian rupee
	CurrencyCodeMUR CurrencyCode = "MUR"
	// Maldivian rufiyaa
	CurrencyCodeMVR CurrencyCode = "MVR"
	// Malawian kwacha
	CurrencyCodeMWK CurrencyCode = "MWK"
	// Mexican peso
	CurrencyCodeMXN CurrencyCode = "MXN"
	// Malaysian ringgit
	CurrencyCodeMYR CurrencyCode = "MYR"
	// Mozambican metical
	CurrencyCodeMZN CurrencyCode = "MZN"
	// Namibian dollar
	CurrencyCodeNAD CurrencyCode = "NAD"
	// Nigerian naira
	CurrencyCodeNGN CurrencyCode = "NGN"
	// Nicaraguan córdoba
	CurrencyCodeNIO CurrencyCode = "NIO"
	// Norwegian krone
	CurrencyCodeNOK CurrencyCode = "NOK"
	// Nepalese rupee
	CurrencyCodeNPR CurrencyCode = "NPR"
	// New Zealand dollar
	CurrencyCodeNZD CurrencyCode = "NZD"
	// Omani rial
	CurrencyCodeOMR CurrencyCode = "OMR"
	// Panamanian balboa
	CurrencyCodePAB CurrencyCode = "PAB"
	// Peruvian sol
	CurrencyCodePEN CurrencyCode = "PEN"
	// Papua New Guinean kina
	CurrencyCodePGK CurrencyCode = "PGK"
	// Philippine peso
	CurrencyCodePHP CurrencyCode = "PHP"
	// Pakistani rupee
	CurrencyCodePKR CurrencyCode = "PKR"
	// Polish złoty
	CurrencyCodePLN CurrencyCode = "PLN"
	// Paraguayan guaraní
	CurrencyCodePYG CurrencyCode = "PYG"
	// Qatari riyal
	CurrencyCodeQAR CurrencyCode = "QAR"
	// Romanian leu
	CurrencyCodeRON CurrencyCode = "RON"
	// Serbian dinar
	CurrencyCodeRSD CurrencyCode = "RSD"
	// Russian ruble
	CurrencyCodeRUB CurrencyCode = "RUB"
	// Rwandan franc
	CurrencyCodeRWF CurrencyCode = "RWF"
	// Saudi riyal
	CurrencyCodeSAR CurrencyCode = "SAR"
	// Solomon Islands dollar
	CurrencyCodeSBD CurrencyCode = "SBD"
	// Seychelles rupee
	CurrencyCodeSCR CurrencyCode = "SCR"
	// Sudanese pound
	CurrencyCodeSDG CurrencyCode = "SDG"
	// Swedish krona/kronor
	CurrencyCodeSEK CurrencyCode = "SEK"
	// Singapore dollar
	CurrencyCodeSGD CurrencyCode = "SGD"
	// Saint Helena pound
	CurrencyCodeSHP CurrencyCode = "SHP"
	// Sierra Leonean leone
	CurrencyCodeSLL CurrencyCode = "SLL"
	// Somali shilling
	CurrencyCodeSOS CurrencyCode = "SOS"
	// Surinamese dollar
	CurrencyCodeSRD CurrencyCode = "SRD"
	// South Sudanese pound
	CurrencyCodeSSP CurrencyCode = "SSP"
	// São Tomé and Príncipe dobra
	CurrencyCodeSTN CurrencyCode = "STN"
	// Salvadoran colón
	CurrencyCodeSVC CurrencyCode = "SVC"
	// Syrian pound
	CurrencyCodeSYP CurrencyCode = "SYP"
	// Swazi lilangeni
	CurrencyCodeSZL CurrencyCode = "SZL"
	// Thai baht
	CurrencyCodeTHB CurrencyCode = "THB"
	// Tajikistani somoni
	CurrencyCodeTJS CurrencyCode = "TJS"
	// Turkmenistan manat
	CurrencyCodeTMT CurrencyCode = "TMT"
	// Tunisian dinar
	CurrencyCodeTND CurrencyCode = "TND"
	// Tongan paʻanga
	CurrencyCodeTOP CurrencyCode = "TOP"
	// Turkish lira
	CurrencyCodeTRY CurrencyCode = "TRY"
	// Trinidad and Tobago dollar
	CurrencyCodeTTD CurrencyCode = "TTD"
	// New Taiwan dollar
	CurrencyCodeTWD CurrencyCode = "TWD"
	// Tanzanian shilling
	CurrencyCodeTZS CurrencyCode = "TZS"
	// Ukrainian hryvnia
	CurrencyCodeUAH CurrencyCode = "UAH"
	// Ugandan shilling
	CurrencyCodeUGX CurrencyCode = "UGX"
	// United States dollar
	CurrencyCodeUSD CurrencyCode = "USD"
	// Uruguayan peso
	CurrencyCodeUYU CurrencyCode = "UYU"
	// Uzbekistan som
	CurrencyCodeUZS CurrencyCode = "UZS"
	// Venezuelan bolívar soberano
	CurrencyCodeVES CurrencyCode = "VES"
	// Vietnamese đồng
	CurrencyCodeVND CurrencyCode = "VND"
	// Vanuatu vatu
	CurrencyCodeVUV CurrencyCode = "VUV"
	// Samoan tala
	CurrencyCodeWST CurrencyCode = "WST"
	// CFA franc BEAC
	CurrencyCodeXAF CurrencyCode = "XAF"
	// East Caribbean dollar
	CurrencyCodeXCD CurrencyCode = "XCD"
	// CFA franc BCEAO
	CurrencyCodeXOF CurrencyCode = "XOF"
	// CFP franc (franc Pacifique)
	CurrencyCodeXPF CurrencyCode = "XPF"
	// Yemeni rial
	CurrencyCodeYER CurrencyCode = "YER"
	// South African rand
	CurrencyCodeZAR CurrencyCode = "ZAR"
	// Zambian kwacha
	CurrencyCodeZMW CurrencyCode = "ZMW"
	// Zimbabwean dollar
	CurrencyCodeZWL CurrencyCode = "ZWL"
)