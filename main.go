package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type RegexConfig struct {
	Key    string
	Regex  string
	Ignore string
}

var regexConfigurations = []RegexConfig{
	{
		Key:    "Amazon_AWS_Access_Key_ID",
		Regex:  `([^A-Z0-9]|^)(AKIA|A3T|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{12,}`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `//s3-[a-z0-9-]+\.amazonaws\.com/[a-z0-9._-]+`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `//s3\.amazonaws\.com/[a-z0-9._-]+`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `[a-z0-9.-]+\.s3-[a-z0-9-]\.amazonaws\.com`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `[a-z0-9.-]+\.s3-website[.-](eu|ap|us|ca|sa|cn)`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `[a-z0-9.-]+\.s3\.amazonaws\.com`,
		Ignore: "",
	},
	{
		Key:    "Amazon_AWS_S3_Bucket",
		Regex:  `amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
		Ignore: "",
	},
	{
		Key:    "Microsoft_Blobs",
		Regex:  `https://[^/]+\.blob\.core\.windows\.net/[^/]+/(.+)`,
		Ignore: "",
	},
	{
		Key:    "Authorization_Basic",
		Regex:  `(b|B)asic [a-zA-Z0-9=:_\+\/-]{5,100}`,
		Ignore: "authentication|JavaScript|easing|announcement|template|states|class|subclass|editor|debounce|material|sanity|event|styles|prototype|enough|screen|logging|query|props|O-Form|cases|sorting|varfcast|steps|placement|check|understanding|transclusion|object|index|title|structure|coercion|mouse|numbers|initialisation|bindings|comparasion|comparator|catch-all|control|instance|focus|logic|browser|stroke|canvas|install|tutorial|initialization|cleanup|backdrop|label-container|copy-paste|search|numeric|sameness|access|security|tests|operations|parts|scroll|options|accounts|setup|google|dashboard|segmentation|features|button-text|behaviors|knowledge|directive|implementation|Expressions|functions|html5|debug|custom|startup|building|native|rules|properties|functionalit|alphabetical|border|alert|panel|markup|support|latin|config|wrapper|field|gauge|settings|constructor|wrapping|anchor|example|content|validation|dropdown|format|angular|arrow|layout|string|chart|javascript|change|server|usage|filter|Multilingual|realm|image|color|configuration|syntax|authorization|authentication|realm=|setting|chips|compatibility|menus|company|container|widget|in-memory|comparison|mapping|datastore|recommended|aspect|spinner|advisory|questions|login|economy|registration|cache|finishes|services|store|concept|interface|roles|styling|visualelement|dojox|types|uploader|photoswipe|requirement|algorithm|request|overload|salary|additional|military|kickass|backbone|library|bone|animation|button|detail|membership|melody|channel|metadata|package|musical|video|audio|workflow|guideline|attached|shortcut|hidden|keyboard|mask|atomic|xpath|input|scaling|common|translation|command|ferting|data|icon|test|vectera|amount|toggle|expression|vertical|approach|number|pagination|point|validators|style|firebase|aline|params|message|behavior|component|default|toolbar|ajax|framework|stage|mesocycle|pulse|arrangements|sound|version|analogue|patch|lesson|track|mechanics|civil|measurements|principels|Alice|groove|sampler|stuff|block|principels|soulfulness|set-up|relationship|children|pattern|impact|pre-production|soldering|thing|primal|bedroom|principles|software|tool|proposal|plug-ins|chords|values|methods|basicform|mindset|action|state|console|production|studio|bouncing|question|training|necessities|personal|social|needs|fertig|math|selve|ideia|rhythm|frequency|guitar|principle|8-bit|platform|element|ideas|hide/show|symbol|visualization",
	},
	{
		Key:    "Authorization_Bearer",
		Regex:  `(b|b)earer [a-zA-Z0-9_\-\.=:_\+\/]{5,100}`,
		Ignore: "prefix|authentication|token|format|your_token_here|your_key_above",
	},
	{
		Key:   "AWS_API_Key",
		Regex: `AKIA[0-9A-Z]{16}`,
	},
	{
		Key:    "Cloudinary_Basic_Auth",
		Regex:  `cloudinary://[0-9]{15}:[0-9A-Za-z]+@[a-z]+`,
		Ignore: "",
	},
	{
		Key:    "Discord_BOT_Token",
		Regex:  `((?:N|M|O)[a-zA-Z0-9]{23}\.[a-zA-Z0-9-_]{6}\.[a-zA-Z0-9-_]{27})$`,
		Ignore: "",
	},
	{
		Key:    "Firebase",
		Regex:  `[a-z0-9.-]+\.firebaseio\.com`,
		Ignore: "",
	},
	{
		Key:    "GitHub_Access_Token",
		Regex:  `([a-zA-Z0-9_-]*:[a-zA-Z0-9_-]+@github.com*)$`,
		Ignore: "",
	},
	{
		Key:    "Google_API_Key",
		Regex:  `AIza[0-9A-Za-z\-_]{35}`,
		Ignore: "",
	},
	{
		Key:    "Google_Cloud_Platform_OAuth",
		Regex:  `[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com`,
		Ignore: "",
	},
	{
		Key:    "Google_Cloud_Platform_Service_Account",
		Regex:  `"type": "service_account"`,
		Ignore: "",
	},
	{
		Key:    "Google_OAuth_Access_Token",
		Regex:  `ya29\.[0-9A-Za-z\-_]+`,
		Ignore: "",
	},
	{
		Key:    "Heroku_API_Key",
		Regex:  `[h|H][e|E][r|R][o|O][k|K][u|U].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}`,
		Ignore: "",
	},
	{
		Key:    "JSON_Web_Token",
		Regex:  `eyJ[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+`,
		Ignore: "",
	},
	{
		Key:    "MailChimp_API_Key",
		Regex:  `([0-9a-f]{32}-us[0-9]{1,2}|(mailchimp|mail-chimp|mail_chimp)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100})`,
		Ignore: "",
	},
	{
		Key:    "Mailgun_API_Key",
		Regex:  `(key-[0-9a-zA-Z]{32}|(mailgun|mail-gun|mail_gun)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100})`,
		Ignore: "",
	},
	{
		Key:    "Password_in_URL",
		Regex:  `[a-zA-Z]{3,10}://[^/ :@]{3,20}:[^/ :@]{3,20}@.{1,100}["\' ]`,
		Ignore: "",
	},
	{
		Key:    "PayPal_Braintree_Access_Token",
		Regex:  `access_token\$production\$[0-9a-z]{16}\$[0-9a-f]{32}`,
		Ignore: "",
	},
	{
		Key:    "Slack_Token",
		Regex:  `(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})`,
		Ignore: "",
	},
	{
		Key:    "Slack_Webhook",
		Regex:  `https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}`,
		Ignore: "",
	},
	{
		Key:    "Discord_Webhook",
		Regex:  `https://discord.com/api/webhooks/[0-9]{19}/[a-zA-Z0-9_-]+`,
		Ignore: "",
	},
	{
		Key:    "Square_Access_Token",
		Regex:  `sq0atp-[0-9A-Za-z\-_]{22}`,
		Ignore: "",
	},
	{
		Key:    "Square_OAuth_Secret",
		Regex:  `sq0csp-[0-9A-Za-z\-_]{43}`,
		Ignore: "",
	},
	{
		Key:    "Stripe_API_Key",
		Regex:  `sk_live_[0-9a-zA-Z]{24}`,
		Ignore: "",
	},
	{
		Key:    "Stripe_Restricted_API_Key",
		Regex:  `rk_live_[0-9a-zA-Z]{24}`,
		Ignore: "",
	},
	{
		Key:    "Twilio_API_Key",
		Regex:  `SK[0-9a-fA-F]{32}`,
		Ignore: "",
	},
	{
		Key:    "Passwords",
		Regex:  `(passwordsalt|password_salt|password-salt|passwordhash|password_hash|password-hash|keypassword|pwd|password|passwords|passwd|pass|root_password|root-password|rootpassword|rootpasswd|secret.password|gmail_password|gmail-password|gmail_username|gmail-username|login)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100}`,
		Ignore: "allow|forgot|reset|english|parola|enter|passwordflow|senha|parola|lozinka|LykilorÃ|furaha|heslo|passwort|generate|nenosiri|JelszÃ³|Adgangskode|changing|Generoitu|remember|ContraseÃ|eye-slash|Structure|company|advisory|agile|comparator|comparasion|widget|instance|request|control|finish|store|authority|editor|connect|calculation|placement|workflow|confirm|change|opaque|offscreen|text-success|text-warning|bypass|signing|translucent|generatedPass|oversome|overevery|Contraseña|ingres|Forget|repit|nueva|nova|setchange|repet|introduc|actua|distin|confirma|repeat|cadastre|retype|gerar|esqueceu|salvar|please|setuserpass|enviar|esqueci|entrar|validar|glemt|kodeord|error|indtast|haspagepassword|savepassword|checkpagepassword|needspagepassword|invalid|formmodel|element|Wachtwoord|Salasana|parole|parool|sandi|Vaihda|vahvista|nykyinen|input|genereer|bevestig|wrong|wheel|username|editpassword|cancelpassword|savepagepassword|abstract|aria-level|there|required|document|content|scope|aria-label|Oppfyller|Gjeldende|BekrÃ|strong|remove|incorrect|Establecer|définir|créer|configurar|create|passord|current|herhaal|Nieuw|Okwuntughe|Vraag|gebruik|wijzig|onthul|parol|FjalÃ|geslo|Verwenden|event|aria-roledescription|links|accesskey|landmarks|value|description|display|heading|valid|verifiera|spara|lagre|endre|antamasi|SlaptaÅ¾odis|Tallenna|vahva|Bekreft|Tilbakestille|Sterkt|Vennligst|:\"password|:'password|Tilbakestill|: \"Password|verify|watermarkPassword|Webservice|aseta|vanha|verschillen|Gehashed|secondPassword|oldPassword|firstPassword|cancel|Neues|active|Lösenord|Répéter|Pakartokite|Podaj|Powtórz|Ustaw|Upprepa|Gjenta|Skapa|Syötä|SETTINGS|LOGIN|doesExternalFilterPass|weak-password|EmailAuthProvider|weak-password|Palauta|Lösenord|Olvidé|Cancel|L\\\u00f6senord|Nouveau|L\\\xf6senord|Piilota|Fastst\\\xe4ll|L\\\xf6senord|G\\\xf6mma|P\\\xe5loggingsinformasjonen|correo|Bekr\\\xe6ft|St\\\xe6rkt|Feilet|Nåræende|N\\\xe4yt\\\xe4|\\\xc5terst\\\xe4ll|Lis\\\xe4\\\xe4|P\\\xe4ivit\\\xe4|Gl\\\xf6mt|V\\\xe6lg|Tilf\\\xf8j|Nulstil|Unohditko|opret|gentag|toista|opprett|crear|Ställ|Įveskite|computer|datasystem|R\\\xe9initialiser|=\"password|string|cancelar|contrase|mostrar|please|esconder|criar|hide|setpassword|old|bad|text|this|existing|medium|weak|lost|save|utilize|você|Recuperação|redefinir|informe|pronto|ocultar|exibir|your|update|updating|ustvari|Nastavite|nastavi|postavi|Jelszavam|minha|define|has|continuer|Promijeni|unesi|Kreiraj|imposta|continua|mano|muuda|Praegune|cari|yeni|Promijenite|email|continue|screen|shadow|picking|unable|enable|guest|Recover|bind|server|privacy|Authentication|zaloguj|ZapomniaÅ|profile|missing|editFormLabels|montrez|masquez|cambia|nuova|Inserisci|insert|show|Aktuelles|GET_USER_PASSWORD|SET_USER_PASSWORD|quickRegistrNewPasswordHidden",
	},
	{
		Key:    "ApiKey Parameters",
		Regex:  `(api(_|-|)key(_|-|)|api_token|api-token|apitoken|api_secret|api-secret|apisecret|client_key|client-key|clientkey|api_docs|api-docs|apidocs|api_key|api-key|apikey)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100}`,
		Ignore: "showCustomOnKeyup|Ingen|foobar|clienthost|msg|description|parent|processid|type|emailid|block|subject|senderip|scstatus|your|enter|api|attachment|edrDataSourceDistinct|datasource|edrMD5Distinct|you|debes|nonprimary|e-mail|return|smtp",
	},
	{
		Key:    "FTP",
		Regex:  `\bftp://[\S]{3,50}:([\S]{3,50})@[-.%\w\/:]+\b`,
		Ignore: "",
	},
	{
		Key:    "ElasticSearch",
		Regex:  `(Bearer.(private|public|search|admin)-)[A-z0-9]{1,35}`,
		Ignore: "",
	},
	{
		Key:    "AWS Cognito",
		Regex:  `(user_pools_web_client_id|user-pools-web-client-id|user_poools_id|userpoolwebclientids|user-poools-id|user_pools_id|user-pools-id|poolid|pool_id|userpoolid|userpoolwebclientid|user-pool-id|user_pool_id|identitypoolid|identity_pool_id|identity-pool-id|poolid)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100}`,
		Ignore: "",
	},
	{
		Key:    "Azure Links",
		Regex:  `[a-z0-9.-]+\.azurewebsites\.net`,
		Ignore: "",
	},
	{
		Key:    "Shopfy secret",
		Regex:  `shp(ss|at|pa)_[a-fA-F0-9]{32}`,
		Ignore: "",
	},
	{
		Key:    "graphql URls",
		Regex:  `(graphql|graphql_api_url|graphql-api-url|graphql_url|graphql-url|graphql_api|graphql-api)((:|=| : | = )(\"|'| \"| \')|(\"|'|' |\" )(:|=)(\"|'| '| \"))[0-9A-Za-z_-]{2,100}`,
		Ignore: "",
	},
	{
		Key:    "Social Media",
		Regex:  `(http(s|):\/\/)(www\.|)(instagram.com|youtube.com|twitter.com|facebook.com|pinterest.com)\/[A-Za-z0-9_-]{1,31}`,
		Ignore: "widgets|/pin/create|/pinmarklet|/images|pinit|/share|/plugin",
	},
}

func main() {
	results := make(chan string)
	var wg sync.WaitGroup

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				processFile(path, results)
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error to search files:", err)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func shouldIgnoreMatch(match string, ignoreWords string) bool {
	ignoreList := strings.Split(ignoreWords, "|")
	for _, word := range ignoreList {
		if strings.Contains(match, word) {
			return true
		}
	}
	return false
}

func processFile(path string, results chan<- string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error to read file '%s': %v\n", path, err)
		return
	}

	for _, regexCfg := range regexConfigurations {
		regex, err := regexp.Compile(regexCfg.Regex)
		if err != nil {
			fmt.Printf("Error on regex '%s': %v\n", regexCfg.Key, err)
			continue
		}

		matches := regex.FindAllString(string(content), -1)

		for _, match := range matches {
			if !shouldIgnoreMatch(match, regexCfg.Ignore) {
				results <- match
				if len(matches) > 0 {
					fmt.Println(color.WhiteString("============================"))
					fmt.Println(color.GreenString("Key: " + regexCfg.Key))
					fmt.Println(color.GreenString("File: " + path))
					fmt.Println(color.MagentaString("Regex: " + regexCfg.Regex))
					fmt.Println(color.CyanString("Match:"))
					for _, match := range matches {
						fmt.Println(color.CyanString(match))
					}
					fmt.Println()
				}
			}
		}
	}
}
