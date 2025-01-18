package main

import (
	"regexp"
	"strings"
)

// Function to find substrings and determine variable/property names
func findSensitiveData(sourceCode string) map[string]string {
	dataMap := make(map[string]string)

	varPattern := regexp.MustCompile(`(?i)\s+([\w\d_]+)\s*=\s*['"` + "`" + `]([^'"` + "`" + `]+)['"` + "`" + `]`)
	propPattern := regexp.MustCompile(`(?is)['"` + "`" + `]?([\w\d_]+)['"` + "`" + `]?\s*:\s*['"` + "`" + `]([^'"` + "`" + `]*)['"` + "`" + `]`)

	for _, match := range varPattern.FindAllStringSubmatch(sourceCode, -1) {
		if containsSensitiveSubstring(match[1]) && isPossibleCredential(match[2]) {
			dataMap[match[1]] = match[2]
		}
	}

	for _, match := range propPattern.FindAllStringSubmatch(sourceCode, -1) {
		if containsSensitiveSubstring(match[1]) && isPossibleCredential(match[2]) {
			dataMap[match[1]] = match[2]
		}
	}

	return dataMap
}

// func isPossibleUrl(value string) bool {
// 	// List of 200 most popular domain zones
// 	domainZones := []string{
// 		".com", ".net", ".org", ".uk", ".ru", ".de", ".xyz", ".info", ".shop", ".online", ".top", ".ch", ".nl", ".store", ".fr", ".site", ".it", ".se", ".biz", ".cn", ".vip", ".pl", ".br", ".eu", ".jp", ".sbs", ".au", ".pro", ".рф", ".app", ".bond", ".ca", ".lol", ".club", ".cz", ".live", ".es", ".be", ".today", ".in", ".fun", ".click", ".tech", ".dk", ".dev", ".at", ".sk", ".life", ".icu", ".cloud",
// 		".co", ".asia", ".buzz", ".space", ".blog", ".hu", ".cyou", ".world", ".ir", ".art", ".kr", ".website", ".za", ".one", ".io", ".mobi", ".ro", ".cfd", ".work", ".link", ".fi", ".gr", ".no", ".nu", ".group", ".tr", ".ua", ".nz", ".mx", ".ar", ".tokyo", ".vn", ".us", ".id", ".cl", ".me", ".digital", ".网址", ".lat", ".pt", ".tw", ".ltd", ".studio", ".agency", ".cc", ".email", ".su", ".il", ".cat", ".design",
// 		".solutions", ".ie", ".skin", ".name", ".lt", ".bet", ".services", ".media", ".news", ".fyi", ".homes", ".ink", ".wiki", ".network", ".win", ".company", ".ovh", ".rest", ".autos", ".zone", ".tv", ".academy", ".ee", ".pics", ".wang", ".si", ".best", ".my", ".love", ".loan", ".by", ".rs", ".rocks", ".nyc", ".hr", ".global", ".sg", ".ai", ".ws", ".kz", ".africa", ".guru", ".team", ".monster", ".games", ".hk", ".ph", ".bid", ".mom", ".bio",
// 		".cam", ".berlin", ".quest", ".lv", ".page", ".social", ".llc", ".med", ".ooo", ".care", ".chat", ".makeup", ".ae", ".boats", ".tel", ".earth", ".health", ".fit", ".beauty", ".center", ".business", ".consulting", ".photography", ".support", ".systems", ".wtf", ".pk", ".yachts", ".bg", ".th", ".pe", ".help", ".foundation", ".family", ".london", ".events", ".plus", ".finance", ".technology", ".run", ".pink", ".realtor", ".church", ".gmbh", ".city", ".在线", ".education", ".expert", ".pictures", ".photo",
// 		".pizza", ".works", ".motorcycles", ".商标", ".bayern", ".capital", ".marketing", ".ng", ".baby", ".ninja", ".international", ".is", ".coffee", ".coach", ".farm", ".商城", ".lu", ".casino", ".photos", ".tools", ".travel", ".party", ".house", ".legal", ".community", ".公司", ".software", ".energy", ".cafe", ".cool", ".swiss", ".ing", ".pet", ".host", ".hair", ".rent", ".school", ".bar", ".training", ".ma", ".christmas", ".video", ".tips", ".law", ".ventures", ".college", ".land", ".yoga", ".fund", ".uz",
// 	}

// 	// Check if the value ends with a valid domain zone
// 	for _, zone := range domainZones {
// 		if strings.HasSuffix(value, zone) {
// 			return true
// 		}
// 	}

// 	return strings.HasPrefix(value, "https")
// }

func isPossibleCredential(value string) bool {
	credPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^[a-zA-Z0-9_.=-]{16,}$`), // Long alphanumeric strings
		regexp.MustCompile(`\d`),                     // Contains at least one digit
	}

	return len(value) >= 12 && credPatterns[0].MatchString(value) && credPatterns[1].MatchString(value)
}

func containsSensitiveSubstring(key string) bool {
	sensitiveSubstrings := []string{
		"api", "key", "secret", "token", "auth", "domain", "config", "project",
		"xsrf", "csrf", "jwt", "password", "passwd", "credential", "creds", "access",
		"grant", "client", "clientId", "clientSecret", "accessToken", "refreshToken",
		"signature", "encrypt", "decryption", "salt", "hash", "pin", "otp",
		"certificate", "cert", "private", "public", "endpoint", "url", "uri",
		"connection", "db", "database", "dsn", "host", "port", "username", "user",
		"login", "email", "mail", "smtp", "imap", "pop", "oauth", "sso", "session",
		"cookie", "vault", "storage", "file", "path", "filepath", "filePath", "dir",
		"directory", "folder", "location", "store", "cache", "temp", "temporary",
		"env", "environment", "variable", "vars", "setting", "settings", "option",
		"options", "param", "parameter", "parameters", "params", "property", "props",
		"prop", "data", "datum", "record", "entry", "info", "information", "metadata",
		"meta", "detail", "details", "secretKey", "accessKey", "apiKey", "apiSecret",
		"apiToken", "dbPassword", "dbUser", "dbName", "dbHost", "dbPort",
		"smtpPassword", "smtpUser", "smtpHost", "smtpPort", "oauthToken",
		"oauthSecret", "oauthKey", "sessionToken", "sessionSecret", "jwtToken",
		"jwtSecret", "encryptionKey", "decryptionKey", "publicKey", "privateKey",
		"sshKey", "sshPrivateKey", "sshPublicKey", "ftpPassword", "ftpUser",
		"ftpHost", "ftpPort", "sftpPassword", "sftpUser", "sftpHost", "sftpPort",
		"redisPassword", "redisHost", "redisPort", "cachePassword", "cacheHost",
		"cachePort", "storageKey", "storageSecret", "storageToken",
		"awsAccessKey", "awsSecretKey", "awsToken", "azureKey", "azureSecret",
		"gcpKey", "gcpSecret", "gcpToken", "stripeKey", "stripeSecret",
		"stripeToken", "paypalKey", "paypalSecret", "paypalToken",
		"facebookToken", "googleToken", "linkedinToken", "githubToken",
		"gitToken", "bitbucketToken", "slackToken", "discordToken",
		"telegramToken", "telegramApi", "telegramBot", "twilioToken",
		"twilioKey", "sendgridApi", "sendgridKey", "mailgunApi", "mailgunKey",
		"herokuApi", "herokuKey", "dockerApi", "dockerKey", "kubernetesApi",
		"kubernetesKey", "terraformKey", "ansibleKey", "jenkinsKey",
		"circleciKey", "travisciKey", "githubSecret", "gitlabSecret",
		"bitbucketSecret", "awsRegion", "awsBucket", "s3Key", "s3Secret",
		"s3Token", "s3AccessKey", "s3SecretKey", "azureSubscription",
		"azureTenant", "azureClientId", "azureClientSecret", "gcpProject",
		"gcpClientId", "gcpClientSecret", "gcpClientToken", "databaseUrl",
		"databaseUri", "databaseConnection", "databaseConfig",
		"dbConnectionString", "dbURI", "dbURL", "smtpConfig",
		"smtpCredentials", "smtpSettings", "emailConfig", "emailCredentials",
		"emailSettings", "apiCredentials", "apiConfig", "apiSettings",
		"authCredentials", "authConfig", "authSettings", "jwtConfig",
		"jwtCredentials", "jwtSettings", "oauthConfig", "oauthCredentials",
		"oauthSettings", "ssoConfig", "ssoCredentials", "ssoSettings",
		"cacheConfig", "cacheCredentials", "cacheSettings", "vaultConfig",
		"vaultCredentials", "vaultSettings", "sessionConfig",
		"sessionCredentials", "sessionSettings", "cookieConfig",
		"cookieCredentials", "cookieSettings", "sslKey", "sslCert",
		"tlsKey", "tlsCert", "encryptionSecret", "encryptionCredentials",
		"encryptionConfig", "decryptionSecret", "decryptionCredentials",
		"decryptionConfig", "accessControl", "securityToken",
		"refresh_token", "authToken", "serviceAccount", "serviceKey",
		"serviceSecret", "dbCredentials", "apiCredentials", "externalKey",
		"externalSecret", "integrationKey", "integrationSecret",
		"thirdPartyKey", "thirdPartySecret", "applicationKey",
		"applicationSecret", "clientCredentials", "serverKey",
		"serverSecret", "databaseSecret", "backupKey", "backupSecret",
		"licenseKey", "licenseSecret", "activationKey", "activationSecret",
		"billingKey", "billingSecret", "trackingToken", "monitoringToken",
		"analyticsKey", "analyticsSecret", "paymentKey", "paymentSecret",
		"notificationKey", "notificationSecret", "loggingKey",
		"loggingSecret", "reportingKey", "reportingSecret", "userToken",
		"userSecret", "accessSecret", "gatewayKey", "gatewaySecret",
		"paymentToken", "paymentKey", "paymentSecret", "webhookSecret",
		"webhookToken",
	}
	k := strings.ToLower(key)
	for _, s := range sensitiveSubstrings {
		if strings.Contains(k, s) {
			return true
		}
	}
	return false
}
