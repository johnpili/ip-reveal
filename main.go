package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-zoo/bone"
	"github.com/johnpili/ip-echo/controllers"
	"github.com/johnpili/ip-echo/models"
	"gopkg.in/yaml.v2"
)

// Configurations / Settings
var (
	configuration models.Config
	BuildVersion  = ""
)

func main() {

	pid := os.Getpid()
	err := ioutil.WriteFile("application.pid", []byte(strconv.Itoa(pid)), 0666) // Used to kill this program
	if err != nil {
		log.Print(err)
	}

	var configLocation string
	flag.StringVar(&configLocation, "config", "config.yml", "Set the location of configuration file")
	flag.Parse()

	printProductBanner()

	configuration = loadConfiguration(configLocation)

	port := strconv.Itoa(configuration.HTTP.Port)
	if os.Getenv("ASPNETCORE_PORT") != "" {
		port = os.Getenv("ASPNETCORE_PORT") // Override port if deployed in IIS
	}

	viewBox := rice.MustFindBox("views")
	staticBox := rice.MustFindBox("static")

	countries := make(map[string]string)
	countries["AF"] = "Afghanistan"
	countries["AL"] = "Albania"
	countries["DZ"] = "Algeria"
	countries["AS"] = "American Samoa"
	countries["AD"] = "Andorra"
	countries["AO"] = "Angola"
	countries["AI"] = "Anguilla"
	countries["AQ"] = "Antarctica"
	countries["AG"] = "Antigua and Barbuda"
	countries["AR"] = "Argentina"
	countries["AM"] = "Armenia"
	countries["AW"] = "Aruba"
	countries["AU"] = "Australia"
	countries["AT"] = "Austria"
	countries["AZ"] = "Azerbaijan"
	countries["BS"] = "Bahamas (the)"
	countries["BH"] = "Bahrain"
	countries["BD"] = "Bangladesh"
	countries["BB"] = "Barbados"
	countries["BY"] = "Belarus"
	countries["BE"] = "Belgium"
	countries["BZ"] = "Belize"
	countries["BJ"] = "Benin"
	countries["BM"] = "Bermuda"
	countries["BT"] = "Bhutan"
	countries["BO"] = "Bolivia (Plurinational State of)"
	countries["BQ"] = "Bonaire, Sint Eustatius and Saba"
	countries["BA"] = "Bosnia and Herzegovina"
	countries["BW"] = "Botswana"
	countries["BV"] = "Bouvet Island"
	countries["BR"] = "Brazil"
	countries["IO"] = "British Indian Ocean Territory (the)"
	countries["BN"] = "Brunei Darussalam"
	countries["BG"] = "Bulgaria"
	countries["BF"] = "Burkina Faso"
	countries["BI"] = "Burundi"
	countries["CV"] = "Cabo Verde"
	countries["KH"] = "Cambodia"
	countries["CM"] = "Cameroon"
	countries["CA"] = "Canada"
	countries["KY"] = "Cayman Islands (the)"
	countries["CF"] = "Central African Republic (the)"
	countries["TD"] = "Chad"
	countries["CL"] = "Chile"
	countries["CN"] = "China"
	countries["CX"] = "Christmas Island"
	countries["CC"] = "Cocos (Keeling) Islands (the)"
	countries["CO"] = "Colombia"
	countries["KM"] = "Comoros (the)"
	countries["CD"] = "Congo (the Democratic Republic of the)"
	countries["CG"] = "Congo (the)"
	countries["CK"] = "Cook Islands (the)"
	countries["CR"] = "Costa Rica"
	countries["HR"] = "Croatia"
	countries["CU"] = "Cuba"
	countries["CW"] = "Curaçao"
	countries["CY"] = "Cyprus"
	countries["CZ"] = "Czechia"
	countries["CI"] = "Côte d'Ivoire"
	countries["DK"] = "Denmark"
	countries["DJ"] = "Djibouti"
	countries["DM"] = "Dominica"
	countries["DO"] = "Dominican Republic (the)"
	countries["EC"] = "Ecuador"
	countries["EG"] = "Egypt"
	countries["SV"] = "El Salvador"
	countries["GQ"] = "Equatorial Guinea"
	countries["ER"] = "Eritrea"
	countries["EE"] = "Estonia"
	countries["SZ"] = "Eswatini"
	countries["ET"] = "Ethiopia"
	countries["FK"] = "Falkland Islands (the) [Malvinas]"
	countries["FO"] = "Faroe Islands (the)"
	countries["FJ"] = "Fiji"
	countries["FI"] = "Finland"
	countries["FR"] = "France"
	countries["GF"] = "French Guiana"
	countries["PF"] = "French Polynesia"
	countries["TF"] = "French Southern Territories (the)"
	countries["GA"] = "Gabon"
	countries["GM"] = "Gambia (the)"
	countries["GE"] = "Georgia"
	countries["DE"] = "Germany"
	countries["GH"] = "Ghana"
	countries["GI"] = "Gibraltar"
	countries["GR"] = "Greece"
	countries["GL"] = "Greenland"
	countries["GD"] = "Grenada"
	countries["GP"] = "Guadeloupe"
	countries["GU"] = "Guam"
	countries["GT"] = "Guatemala"
	countries["GG"] = "Guernsey"
	countries["GN"] = "Guinea"
	countries["GW"] = "Guinea-Bissau"
	countries["GY"] = "Guyana"
	countries["HT"] = "Haiti"
	countries["HM"] = "Heard Island and McDonald Islands"
	countries["VA"] = "Holy See (the)"
	countries["HN"] = "Honduras"
	countries["HK"] = "Hong Kong"
	countries["HU"] = "Hungary"
	countries["IS"] = "Iceland"
	countries["IN"] = "India"
	countries["ID"] = "Indonesia"
	countries["IR"] = "Iran (Islamic Republic of)"
	countries["IQ"] = "Iraq"
	countries["IE"] = "Ireland"
	countries["IM"] = "Isle of Man"
	countries["IL"] = "Israel"
	countries["IT"] = "Italy"
	countries["JM"] = "Jamaica"
	countries["JP"] = "Japan"
	countries["JE"] = "Jersey"
	countries["JO"] = "Jordan"
	countries["KZ"] = "Kazakhstan"
	countries["KE"] = "Kenya"
	countries["KI"] = "Kiribati"
	countries["KP"] = "Korea (the Democratic People's Republic of)"
	countries["KR"] = "Korea (the Republic of)"
	countries["KW"] = "Kuwait"
	countries["KG"] = "Kyrgyzstan"
	countries["LA"] = "Lao People's Democratic Republic (the)"
	countries["LV"] = "Latvia"
	countries["LB"] = "Lebanon"
	countries["LS"] = "Lesotho"
	countries["LR"] = "Liberia"
	countries["LY"] = "Libya"
	countries["LI"] = "Liechtenstein"
	countries["LT"] = "Lithuania"
	countries["LU"] = "Luxembourg"
	countries["MO"] = "Macao"
	countries["MG"] = "Madagascar"
	countries["MW"] = "Malawi"
	countries["MY"] = "Malaysia"
	countries["MV"] = "Maldives"
	countries["ML"] = "Mali"
	countries["MT"] = "Malta"
	countries["MH"] = "Marshall Islands (the)"
	countries["MQ"] = "Martinique"
	countries["MR"] = "Mauritania"
	countries["MU"] = "Mauritius"
	countries["YT"] = "Mayotte"
	countries["MX"] = "Mexico"
	countries["FM"] = "Micronesia (Federated States of)"
	countries["MD"] = "Moldova (the Republic of)"
	countries["MC"] = "Monaco"
	countries["MN"] = "Mongolia"
	countries["ME"] = "Montenegro"
	countries["MS"] = "Montserrat"
	countries["MA"] = "Morocco"
	countries["MZ"] = "Mozambique"
	countries["MM"] = "Myanmar"
	countries["NA"] = "Namibia"
	countries["NR"] = "Nauru"
	countries["NP"] = "Nepal"
	countries["NL"] = "Netherlands (the)"
	countries["NC"] = "New Caledonia"
	countries["NZ"] = "New Zealand"
	countries["NI"] = "Nicaragua"
	countries["NE"] = "Niger (the)"
	countries["NG"] = "Nigeria"
	countries["NU"] = "Niue"
	countries["NF"] = "Norfolk Island"
	countries["MP"] = "Northern Mariana Islands (the)"
	countries["NO"] = "Norway"
	countries["OM"] = "Oman"
	countries["PK"] = "Pakistan"
	countries["PW"] = "Palau"
	countries["PS"] = "Palestine, State of"
	countries["PA"] = "Panama"
	countries["PG"] = "Papua New Guinea"
	countries["PY"] = "Paraguay"
	countries["PE"] = "Peru"
	countries["PH"] = "Philippines (the)"
	countries["PN"] = "Pitcairn"
	countries["PL"] = "Poland"
	countries["PT"] = "Portugal"
	countries["PR"] = "Puerto Rico"
	countries["QA"] = "Qatar"
	countries["MK"] = "Republic of North Macedonia"
	countries["RO"] = "Romania"
	countries["RU"] = "Russian Federation (the)"
	countries["RW"] = "Rwanda"
	countries["RE"] = "Réunion"
	countries["BL"] = "Saint Barthélemy"
	countries["SH"] = "Saint Helena, Ascension and Tristan da Cunha"
	countries["KN"] = "Saint Kitts and Nevis"
	countries["LC"] = "Saint Lucia"
	countries["MF"] = "Saint Martin (French part)"
	countries["PM"] = "Saint Pierre and Miquelon"
	countries["VC"] = "Saint Vincent and the Grenadines"
	countries["WS"] = "Samoa"
	countries["SM"] = "San Marino"
	countries["ST"] = "Sao Tome and Principe"
	countries["SA"] = "Saudi Arabia"
	countries["SN"] = "Senegal"
	countries["RS"] = "Serbia"
	countries["SC"] = "Seychelles"
	countries["SL"] = "Sierra Leone"
	countries["SG"] = "Singapore"
	countries["SX"] = "Sint Maarten (Dutch part)"
	countries["SK"] = "Slovakia"
	countries["SI"] = "Slovenia"
	countries["SB"] = "Solomon Islands"
	countries["SO"] = "Somalia"
	countries["ZA"] = "South Africa"
	countries["GS"] = "South Georgia and the South Sandwich Islands"
	countries["SS"] = "South Sudan"
	countries["ES"] = "Spain"
	countries["LK"] = "Sri Lanka"
	countries["SD"] = "Sudan (the)"
	countries["SR"] = "Suriname"
	countries["SJ"] = "Svalbard and Jan Mayen"
	countries["SE"] = "Sweden"
	countries["CH"] = "Switzerland"
	countries["SY"] = "Syrian Arab Republic"
	countries["TW"] = "Taiwan (Province of China)"
	countries["TJ"] = "Tajikistan"
	countries["TZ"] = "Tanzania, United Republic of"
	countries["TH"] = "Thailand"
	countries["TL"] = "Timor-Leste"
	countries["TG"] = "Togo"
	countries["TK"] = "Tokelau"
	countries["TO"] = "Tonga"
	countries["TT"] = "Trinidad and Tobago"
	countries["TN"] = "Tunisia"
	countries["TR"] = "Turkey"
	countries["TM"] = "Turkmenistan"
	countries["TC"] = "Turks and Caicos Islands (the)"
	countries["TV"] = "Tuvalu"
	countries["UG"] = "Uganda"
	countries["UA"] = "Ukraine"
	countries["AE"] = "United Arab Emirates (the)"
	countries["GB"] = "United Kingdom of Great Britain and Northern Ireland (the)"
	countries["UM"] = "United States Minor Outlying Islands (the)"
	countries["US"] = "United States of America (the)"
	countries["UY"] = "Uruguay"
	countries["UZ"] = "Uzbekistan"
	countries["VU"] = "Vanuatu"
	countries["VE"] = "Venezuela (Bolivarian Republic of)"
	countries["VN"] = "Viet Nam"
	countries["VG"] = "Virgin Islands (British)"
	countries["VI"] = "Virgin Islands (U.S.)"
	countries["WF"] = "Wallis and Futuna"
	countries["EH"] = "Western Sahara"
	countries["YE"] = "Yemen"
	countries["ZM"] = "Zambia"
	countries["ZW"] = "Zimbabwe"
	countries["AX"] = "Åland Islands"

	pageController := controllers.PageController{
		ViewBox:       viewBox,
		Configuration: &configuration,
		Countries:     countries,
	}

	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))

	router := bone.New()
	router.Handle("/static/", staticFileServer)
	router.HandleFunc("/", pageController.IndexHandler)
	router.HandleFunc("/json", pageController.JSONHandler)
	router.HandleFunc("/text", pageController.TextHandler)
	router.HandleFunc("/txt", pageController.TextHandler)
	router.HandleFunc("/ip", pageController.TextHandler)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	if configuration.HTTP.IsTLS {
		log.Printf("Server running at https://localhost:%s/\n", port)
		log.Fatal(httpServer.ListenAndServeTLS(configuration.HTTP.ServerCert, configuration.HTTP.ServerKey))
		return
	}
	log.Printf("Server running at http://localhost:%s/\n", port)
	log.Fatal(httpServer.ListenAndServe())
}

func loadConfiguration(c string) models.Config {
	f, err := os.Open(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	var configuration models.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err.Error())
	}
	return configuration
}

func printProductBanner() {
	fmt.Println("")
	fmt.Println("Product: ip.johnpili.com")
	fmt.Println("Version: 2021.1.0")
	fmt.Println("Operating System: ", runtime.GOOS)
	fmt.Println("Architecture: ", runtime.GOARCH)
	fmt.Println("Build Version: ", BuildVersion)
	fmt.Println("System Date/Time: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("")
	fmt.Println("")
}
