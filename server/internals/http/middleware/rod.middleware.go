package middleware

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"server/internals/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

const ROD_KEY string = "@locals/rod"

var ROD_AGENTS []string = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0",
	"Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0",
	"Mozilla/5.0 (Linux; Android 14; 24030PN60G Build/UKQ1.231003.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/122.0.6261.119 Mobile Safari/537.36",
}

type RodMiddleware struct{}

func (self RodMiddleware) New() (*rod.Browser, *rod.Page) {
	scraper := rod.
		New().
		Timeout(time.Minute).
		MustConnect()

	utils.LOGGER.INFO.Println("Setting up Stealth")
	page := self.SetupStealth(scraper)

	return scraper, page
}

func (RodMiddleware) SetupStealth(instance *rod.Browser) *rod.Page {
	page := stealth.MustPage(instance)
	return page
}

func (RodMiddleware) HandleGuard(page *rod.Page) {
	query := "//title[contains(translate(text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'ddos-guard')]"

	utils.LOGGER.INFO.Println("Searching for Guard")
	element, err := page.ElementX(query)
	if err != nil {
		utils.LOGGER.INFO.Printf("Guard not found")
		return
	}

	if (element == nil) {
		utils.LOGGER.INFO.Printf("Guard pointer nil")
		return
	}

	content, err := element.Text()
	if err != nil {
		utils.LOGGER.INFO.Println("Guard found, but content unavailable", err)
		return
	}

	utils.LOGGER.INFO.Printf("Guard found, type: %v\n", content)
}

func (RodMiddleware) GetAgent() string {
	agents_len := len(ROD_AGENTS) - 1
	agent_idx := rand.IntN(agents_len)

	return ROD_AGENTS[agent_idx]
}

func (RodMiddleware) Metadata(page *rod.Page) {
	element := page.MustElement("#broken-image-dimensions.passed")

	for _, row := range element.MustParents("table").First().MustElements("tr:nth-child(n+2)") {
		cells := row.MustElements("td")
		key := cells[0].MustProperty("textContent")

		if strings.HasPrefix(key.String(), "User Agent") {
			fmt.Printf("\t\t%s: %t\n", key, !strings.Contains(cells[1].MustProperty("textContent").String(), "HeadlessChrome/"))
		} else if strings.HasPrefix(key.String(), "Hairline Feature") {
			// Detects support for hidpi/retina hairlines, which are CSS borders with less than 1px in width, for being physically 1px on hidpi screens.
			// Not all the machine suppports it.
			continue
		} else {
			fmt.Printf("\t\t%s: %s\n", key, cells[1].MustProperty("textContent"))
		}
	}

	page.MustScreenshot("")
}
