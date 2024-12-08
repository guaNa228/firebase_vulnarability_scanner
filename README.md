# Firebase vulnurability scanner

Project is designed to parse JS bundles of the given websites and analyze them on subject of having firebase credentials, which could be used for further pentesting.

Also it includes scraping module, which allows getting urls of some ai startups, which js bundles would be further scanned.

## Project structure

```
domains
scraper
 --scraper.go
analyzer.go
fetcher.go
firebase.go
main.go
utils.go
links.txt
results.txt
```

## Scraping module

The scraping is performed using [this resource](https://www.aixploria.com). It contains urls a lot of startups and big companies, which are related to ai sphere. In current state the "latest" section of the website is used, aiming to target startups with higher priority, than already developed companies.

The range of scraping object may be widened by scanning all the ai related resources in this website, but it would require code modification.

The scraping module resolves short links of the website and then writes them to `links.txt` file.

## Domains

For this project i found a collection of domains with `.ai`, `.dev` and `.io` postfix. These domains could be used to widen the range of websites scanned, but it surely would take a lot of time and also domain filtering rules should be implemented, as it contains a lot of mirrors, and websites, which can be detected as poisoned by url only.

## Algorithm

The scanning is preformed in parallel, for every adress being scanned the goroutine is created, in which the content of all website js files(inline in HTML and fetched) are scanned for existance of firebase credentials, required for pentesting, which are `apiKey`, `databaseURL`, `projectId` and `authDomain`.

Then, when this credentials were found, program extracts this credentials in formats of either JSON or JS object.

## Results

The results of scanning are written to `results.txt` file. If the config was found for the domain the line in the file would be:

```
https://chaton.ai/: No config
```

If the configuration was found, it would be presented in results file the following way:

```
https://codeium.com/windsurf: map[apiKey:AIzaSyDsOl-1XpT5err0Tcnx8FFod1H8gVGIycY authDomain:auth.codeium.com projectId:exa2-fb170]
```

## Pentensting thoughts

As we can see, from 258 scanned domains only 6 use firebase. This leads me to thought, that implementing automated vulnarability scanner is kind of obsolete.

Also, the actual vulnurability scanner would check only common database instances like `users` and e.t.c. So, for my opinion, it would be better to perform manual pentesting using `firepwn` with the goal of getting more interesting results, but i think i could implement simple vulnarability checker by your request.

## Scanning improvement thoughts

Now the scanning is performed in the main page of the domain. Actually, this page may not contain firebase config, at the same time, when other pages of the platform may contain them, so the best existing approach would require checking subpages of the domain too.
