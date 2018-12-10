package fortesting

import (
	"bufio"
	"os"

	t "gopkg.in/check.v1"
)

var testData []string = []string{
	"1014", "1021", "1040", "1055", "1081", "1087", "1099", "10990", "1103", "1118", "1123", "1124", "1125", "1134", "11341", "1140", "1141", "1151",
	"1175", "1179", "1180", "1186", "1190", "1201", "1259", "1260", "1267", "12720", "1280", "1291", "1296", "1313", "1318", "1329", "1345", "1350",
	"1363", "1371", "1372", "1376", "1385", "1390", "1392", "1398", "140", "1404", "1412", "1414", "1419", "1420", "1422", "1426", "1431", "1434",
	"1439", "1449", "1466", "147", "1479", "1483", "1490", "1499", "1506", "1514", "1515", "1527", "1570", "158", "1591", "1619", "1625", "1647",
	"1650", "1655", "1671", "168", "1687", "1698", "1705", "1724", "1729", "1734", "1739", "1750", "1754", "1755", "1758", "176", "1766", "1772",
	"1781", "179", "17963", "180", "1807", "1821", "1826", "1831", "1840", "1842", "1866", "1876", "1890", "1904", "1906", "19061", "1910", "1913",
	"1919", "1920", "1923", "193", "1947", "1958", "1971", "1973", "1974", "1976", "1991", "1992", "2018", "2044", "2057", "2069", "2071", "2108",
	"2112", "2120", "2139", "2147", "2178", "2187", "2196", "2208", "2210", "2242", "2262", "2266", "2275", "2306", "2316", "2324", "2336", "304",
	"2344", "2354", "2374", "2378", "2399", "24", "2417", "2435", "2457", "2458", "2468", "2479", "2481", "2493", "24958", "2503", "2506", "2509",
	"2518", "253", "25466", "2558", "2571", "2594", "2601", "2609", "2628", "2630", "2633", "2636", "264", "2655", "2662", "2687", "2703", "272",
	"2725", "27258", "2727", "2733", "2735", "27561", "2761", "2768", "2769", "2773", "2787", "28", "280", "281", "2811", "2815", "2817", "2821",
	"2822", "284", "2845", "2868", "2871", "2877", "288", "2883", "2886", "292", "2925", "2948", "2983", "2988", "299", "2999", "301", "3032", "3043",
	"3069", "3090", "3102", "3121", "3138", "3144", "3149", "3177", "3178", "3179", "3181", "3214", "3215", "3224", "323", "3240", "3250", "3262",
	"3282", "3290", "3295", "3309", "331", "3319", "3323", "3327", "3329", "3338", "33398", "3361", "3363", "3364", "337", "3373", "3374", "339",
	"3421", "3423", "3439", "3442", "3444", "346", "3469", "3473", "3475", "3485", "350", "3504", "3506", "3507", "3513", "3526", "3528", "354",
	"3540", "355", "3562", "3582", "3583", "35993", "36", "3605", "3606", "3622", "3654", "3658", "3661", "3676", "369", "3722", "3738", "3749",
	"3754", "3770", "3779", "3790", "3830", "3831", "3844", "3850", "3853", "3859", "3867", "3880", "3883", "3902", "3909", "3911", "3920", "3939",
	"3941", "396", "3968", "3991", "4006", "4022", "4026", "4033", "4051", "406", "4066", "4068", "4080", "4081", "4083", "4084", "409", "4092",
	"4109", "4112", "4119", "4120", "4124", "4126", "4128", "4153", "4171", "4197", "42", "420", "4217", "4234", "427", "4277", "428", "4286", "48289",
	"4296", "43", "4315", "434", "43465", "4356", "4365", "4374", "4383", "4399", "4411", "4425", "4440", "4448", "4449", "445", "4451", "4464",
	"4469", "4476", "4497", "4518", "4524", "4581", "4592", "4603", "4619", "4620", "4624", "4630", "4631", "4632", "4643", "4658", "4693", "46931",
	"4697", "4701", "4709", "4714", "4717", "472", "4733", "4734", "4735", "4742", "47436", "4747", "4751", "4767", "4768", "4769", "4783", "4786",
	"4816", "4824", "4827", "483", "4856", "4886", "4891", "4894", "4897", "4916", "4928", "4941", "4948", "4957", "4958", "499", "4993", "5001",
	"5002", "5004", "5008", "5026", "50447", "5047", "5050", "50648", "5070", "510", "5117", "5136", "5150", "5161", "517", "5184", "5194", "5201",
	"5203", "5211", "5223", "5233", "5239", "5264", "5286", "5289", "5290", "5295", "5296", "5299", "5305", "5313", "532", "5325", "5326", "5349",
	"535", "5381", "5384", "5385", "5390", "5397", "5400", "5409", "5413", "5421", "5436", "5439", "5451", "5459", "5465", "5470", "5471", "5482",
	"5487", "5488", "5512", "5524", "5551", "5552", "5565", "5574", "5575", "5577", "558", "5595", "5612", "5637", "5653", "5659", "5670", "5676",
	"5681", "56895", "5695", "5728", "5740", "5749", "5755", "5762", "5763", "577", "5770", "5780", "579", "5814", "5829", "5839", "5843", "5844",
	"585", "5871", "5874", "5876", "5879", "5884", "5912", "5913", "5935", "5936", "5937", "5939", "5979", "5980", "5984", "5995", "6011", "6023",
	"6031", "6036", "604", "6044", "6063", "6065", "6068", "6075", "6084", "6085", "6114", "6115", "6138", "6150", "6163", "6164", "6175", "620",
	"6201", "6203", "6207", "6213", "622", "6235", "6247", "6253", "6257", "6267", "62803", "6281", "6285", "6286", "6301", "6305", "6309", "63228",
	"6335", "6345", "6358", "6363", "6368", "6372", "6391", "6392", "6393", "6398", "6408", "6410", "6420", "6426", "6438", "6441", "6444", "647",
	"6475", "6480", "6489", "6493", "6505", "6509", "6510", "6511", "6518", "6532", "6568", "6575", "6588", "6589", "6591", "6602", "6616", "6617",
	"6624", "6625", "6627", "6633", "6646", "6657", "6665", "6672", "6693", "6713", "6726", "6744", "6745", "6749", "6752", "6781", "6785", "68",
	"6805", "6808", "6828", "6829", "6836", "6900", "6911", "6932", "6935", "6944", "6951", "6960", "6968", "6971", "6992", "6996", "6999", "7003",
	"7013", "7014", "7018", "70237", "7028", "7045", "7048", "7054", "7064", "707", "71", "710", "7118", "7124", "713", "7140", "7170", "7181",
	"7182", "7218", "7226", "7228", "723", "7239", "7243", "7244", "7263", "7275", "7282", "7289", "7293", "7299", "730", "7304", "7307", "7321",
	"7329", "7335", "7338", "7352", "7364", "7366", "7399", "741", "7423", "74303", "744", "7440", "7442", "7444", "7464", "7471", "748", "7481",
	"7486", "7487", "7491", "75", "751", "7522", "7530", "7537", "7539", "7544", "7571", "7574", "7585", "7587", "75953", "76", "7634", "764",
	"7652", "7656", "7659", "7672", "7677", "7678", "7695", "770", "7713", "7765", "778", "7780", "7784", "7789", "7799", "7814", "7826", "7829",
	"7831", "7832", "7834", "7851", "7867", "7868", "7878", "7888", "7894", "79", "7905", "7921", "7934", "7945", "7950", "7951", "7966", "7969",
	"7973", "7979", "8", "802", "8024", "8033", "8042", "8046", "8050", "8060", "8063", "8069", "8071", "808", "8080", "8110", "813", "8137",
	"81436", "8149", "8166", "8173", "8175", "8187", "8188", "82", "8203", "8209", "8222", "8228", "8232", "8238", "8243", "8256", "8257", "8258",
	"8262", "8273", "8277", "8282", "8287", "8290", "8313", "8315", "8333", "8342", "8362", "8368", "8379", "8386", "8405", "8406", "8415", "843",
	"8439", "8449", "8496", "85", "8505", "8532", "8537", "8564", "8587", "859", "8595", "8599", "8609", "8619", "862", "8640", "8646", "8685",
	"8697", "8698", "8699", "8711", "8715", "872", "8730", "8732", "8735", "8739", "874", "876", "8768", "8787", "879", "8794", "8799", "8800",
	"8835", "8841", "8850", "8856", "8867", "8870", "890", "8900", "8906", "8930", "8938", "8941", "8948", "8952", "8960", "89722", "8976", "9002",
	"9013", "9015", "902", "9025", "9031", "9050", "9069", "9081", "9084", "9085", "90933", "9119", "9149", "9157", "916", "9164", "9169", "9171",
	"9174", "9180", "9190", "9193", "9202", "9210", "9220", "9223", "9234", "9249", "9261", "9265", "9268", "927", "9273", "9276", "9301", "9306",
	"9320", "9325", "9327", "9331", "9342", "9345", "9349", "9350", "9351", "9352", "9377", "9395", "9415", "9420", "9421", "9430", "9435", "9440",
	"9450", "9453", "946", "9464", "9470", "9478", "9487", "94970", "9500", "9503", "9510", "9528", "953", "9541", "9543", "9556", "9562", "9570",
	"9586", "9599", "9602", "9611", "9618", "9623", "9624", "9628", "9632", "966", "9664", "969", "9712", "9714", "9726", "9729", "9739", "9761",
	"9778", "9781", "9785", "979", "9799", "9814", "9841", "101489", "102189", "104089", "105589", "108189", "108789", "109989", "1099089",
	"110389", "111889", "112389", "112489", "112589", "113489", "1134189", "114089", "114189", "115189", "117589", "117989", "118089", "118689",
	"119089", "120189", "125989", "126089", "126789", "1272089", "128089", "129189", "129689", "131389", "131889", "132989", "134589", "135089",
	"136389", "137189", "137289", "137689", "138589", "139089", "139289", "139889", "14089", "140489", "141289", "141489", "141989", "142089",
	"142289", "142689", "143189", "143489", "143989", "144989", "146689", "14789", "147989", "148389", "149089", "149989", "150689", "151489",
	"151589", "152789", "157089", "15889", "159189", "161989", "162589", "164789", "165089", "165589", "167189", "16889", "168789", "169889",
	"170589", "172489", "172989", "173489", "173989", "175089", "175489", "175589", "175889", "17689", "176689", "177289", "178189", "17989",
	"1796389", "18089", "180789", "182189", "182689", "183189", "184089", "184289", "186689", "187689", "189089", "190489", "190689", "1906189",
	"191089", "191389", "191989", "192089", "192389", "19389", "194789", "195889", "197189", "197389", "197489", "197689", "199189", "199289",
	"201889", "204489", "205789", "206989", "207189", "210889", "211289", "212089", "213989", "214789", "217889", "218789", "219689", "220889",
	"221089", "224289", "226289", "226689", "227589", "230689", "231689", "232489", "233689", "234489", "235489", "237489", "237889", "239989",
	"2489", "241789", "243589", "245789", "245889", "246889", "247989", "248189", "249389", "2495889", "250389", "250689", "250989", "251889",
	"25389", "2546689", "255889", "257189", "259489", "260189", "260989", "262889", "263089", "263389", "263689", "26489", "265589", "266289",
	"268789", "270389", "27289", "272589", "2725889", "272789", "273389", "273589", "2756189", "276189", "276889", "276989", "277389", "278789",
	"2889", "28089", "28189", "281189", "281589", "281789", "282189", "282289", "28489", "284589", "286889", "287189", "287789", "28889",
	"288389", "288689", "29289", "292589", "294889", "298389", "298889", "29989", "299989", "30189", "303289", "30489", "304389", "306989",
	"309089", "310289", "312189", "313889", "314489", "314989", "317789", "317889", "317989", "318189", "321489", "321589", "322489", "32389",
	"324089", "325089", "326289", "328289", "329089", "329589", "330989", "33189", "331989", "332389", "332789", "332989", "333889", "3339889",
	"336189", "336389", "336489", "33789", "337389", "337489", "33989", "342189", "342389", "343989", "344289", "344489", "34689", "346989",
	"347389", "347589", "348589", "35089", "350489", "350689", "350789", "351389", "352689", "352889", "35489", "354089", "35589", "356289",
	"358289", "358389", "3599389", "3689", "360589", "360689", "362289", "365489", "365889", "366189", "367689", "36989", "372289", "373889",
	"374989", "375489", "377089", "377989", "379089", "383089", "383189", "384489", "385089", "385389", "385989", "386789", "388089", "388389",
	"390289", "390989", "391189", "392089", "393989", "394189", "39689", "396889", "399189", "400689", "402289", "402689", "403389", "405189",
	"40689", "406689", "406889", "408089", "408189", "408389", "408489", "40989", "409289", "410989", "411289", "411989", "412089", "412489",
	"412689", "412889", "415389", "417189", "419789", "4289", "42089", "421789", "423489", "42789", "427789", "42889", "428689", "428989",
	"429689", "4389", "431589", "43489", "4346589", "435689", "436589", "437489", "438389", "439989", "441189", "442589", "444089", "444889",
	"444989", "44589", "445189", "446489", "446989", "447689", "449789", "451889", "452489", "458189", "459289", "460389", "461989", "462089",
	"462489", "463089", "463189", "463289", "464389", "465889", "469389", "4693189", "469789", "470189", "470989", "471489", "471789", "47289",
	"473389", "473489", "473589", "474289", "4743689", "474789", "475189", "476789", "476889", "476989", "478389", "478689", "481689", "482489",
	"482789", "48389", "485689", "488689", "489189", "489489", "489789", "491689", "492889", "494189", "494889", "495789", "495889", "49989",
	"499389", "500189", "500289", "500489", "500889", "502689", "5044789", "504789", "505089", "5064889", "507089", "51089", "511789", "513689",
	"515089", "516189", "51789", "518489", "519489", "520189", "520389", "521189", "522389", "523389", "523989", "526489", "528689", "528989",
	"529089", "529589", "529689", "529989", "530589", "531389", "53289", "532589", "532689", "534989", "53589", "538189", "538489", "538589",
	"539089", "539789", "540089", "540989", "541389", "542189", "543689", "543989", "545189", "545989", "546589", "547089", "547189", "548289",
	"548789", "548889", "551289", "552489", "555189", "555289", "556589", "557489", "557589", "557789", "55889", "559589", "561289", "563789",
	"565389", "565989", "567089", "567689", "568189", "5689589", "569589", "572889", "574089", "574989", "575589", "576289", "576389", "57789",
	"577089", "578089", "57989", "581489", "582989", "583989", "584389", "584489", "58589", "587189", "587489", "587689", "587989", "588489",
	"591289", "591389", "593589", "593689", "593789", "593989", "597989", "598089", "598489", "599589", "601189", "602389", "603189", "603689",
	"60489", "604489", "606389", "606589", "606889", "607589", "608489", "608589", "611489", "611589", "613889", "615089", "616389", "616489",
	"617589", "62089", "620189", "620389", "620789", "621389", "62289", "623589", "624789", "625389", "625789", "626789", "6280389", "628189",
	"628589", "628689", "630189", "630589", "630989", "6322889", "633589", "634589", "635889", "636389", "636889", "637289", "639189", "639289",
	"639389", "639889", "640889", "641089", "642089", "642689", "643889", "644189", "644489", "64789", "647589", "648089", "648989", "649389",
	"650589", "650989", "651089", "651189", "651889", "653289", "656889", "657589", "658889", "658989", "659189", "660289", "661689", "661789",
	"662489", "662589", "662789", "663389", "664689", "665789", "666589", "667289", "669389", "671389", "672689", "674489", "674589", "674989",
	"675289", "678189", "678589", "6889", "680589", "680889", "682889", "682989", "683689", "690089", "691189", "693289", "693589", "694489",
	"695189", "696089", "696889", "697189", "699289", "699689", "699989", "700389", "701389", "701489", "701889", "7023789", "702889", "704589",
	"704889", "705489", "706489", "70789", "7189", "71089", "711889", "712489", "71389", "714089", "717089", "718189", "718289", "721889",
	"722689", "722889", "72389", "723989", "724389", "724489", "726389", "727589", "728289", "728989", "729389", "729989", "73089", "730489",
	"730789", "732189", "732989", "733589", "733889", "735289", "736489", "736689", "739989", "74189", "742389", "7430389", "74489", "744089",
	"744289", "744489", "746489", "747189", "74889", "748189", "748689", "748789", "749189", "7589", "75189", "752289", "753089", "753789",
	"753989", "754489", "757189", "757489", "758589", "758789", "7595389", "7689", "763489", "76489", "765289", "765689", "765989", "767289",
	"767789", "767889", "769589", "77089", "771389", "776589", "77889", "778089", "778489", "778989", "779989", "781489", "782689", "782989",
	"783189", "783289", "783489", "785189", "786789", "786889", "787889", "788889", "789489", "7989", "790589", "792189", "793489", "794589",
	"795089", "795189", "796689", "796989", "797389", "797989", "889", "80289", "802489", "803389", "804289", "804689", "805089", "806089",
	"806389", "806989", "807189", "80889", "808089", "811089", "81389", "813789", "8143689", "814989", "816689", "817389", "817589", "818789",
	"818889", "8289", "820389", "820989", "822289", "822889", "823289", "823889", "824389", "825689", "825789", "825889", "826289", "827389",
	"827789", "828289", "828789", "829089", "831389", "831589", "833389", "834289", "836289", "836889", "837989", "838689", "840589", "840689",
	"841589", "84389", "843989", "844989", "849689", "8589", "850589", "853289", "853789", "856489", "858789", "85989", "859589", "859989",
	"860989", "861989", "86289", "864089", "864689", "868589", "869789", "869889", "869989", "871189", "871589", "87289", "873089", "873289",
	"873589", "873989", "87489", "87689", "876889", "878789", "87989", "879489", "879989", "880089", "883589", "884189", "885089", "885689",
	"886789", "887089", "89089", "890089", "890689", "893089", "893889", "894189", "894889", "895289", "896089", "8972289", "897689", "900289",
	"901389", "901589", "90289", "902589", "903189", "905089", "906989", "908189", "908489", "908589", "9093389", "911989", "914989", "915789",
	"91689", "916489", "916989", "917189", "917489", "918089", "919089", "919389", "920289", "921089", "922089", "922389", "923489", "924989",
	"926189", "926589", "926889", "92789", "927389", "927689", "930189", "930689", "932089", "932589", "932789", "933189", "934289", "934589",
	"934989", "935089", "935189", "935289", "937789", "939589", "941589", "942089", "942189", "943089", "943589", "944089", "945089", "945389",
	"94689", "946489", "947089", "947889", "948789", "9497089", "950089", "950389", "951089", "952889", "95389", "954189", "954389", "955689",
	"956289", "957089", "958689", "959989", "960289", "961189", "961889", "962389", "962489", "962889", "963289", "96689", "966489", "96989",
	"971289", "971489", "972689", "972989", "973989", "976189", "977889", "978189", "978589", "97989", "979989", "981489", "984189", "98589",
	"985289", "986189", "987589", "988289", "988489", "989389", "990389", "990689", "991389", "991889", "991989", "994589", "994789", "995589",
	"995789", "99789", "997989", "999589", "999989", "985", "9852", "9861", "9875", "9882", "9884", "9893", "9903", "9906", "9913", "9918", "9919",
	"9945", "9947", "9955", "9957", "997", "9979", "9995", "9999",
}

func ArrayForTesting() []string {
	return testData
}

func Dir() string {
	return os.Getenv("FILE_DIR")
}

// func ArrayForTesting() []string {

// 	testData := []string{}

// 	st := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
// 	for _, a := range st {
// 		for _, b := range st {
// 			for _, c := range st {
// 				testData = append(testData, a+"-"+b+"-"+c)
// 			}
// 		}
// 	}

// 	return testData
// }

func CheckFiles(c *t.C, file1, file2 string) {
	b1, err := ReadFile(file1)
	c.Assert(err, t.IsNil)

	b2, err := ReadFile(file2)
	c.Assert(err, t.IsNil)

	c.Assert(len(b1), t.Equals, len(b2))
	for i := range b1 {
		c.Assert(b1[i], t.DeepEquals, b2[i])
	}
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
