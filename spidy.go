package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)



var (
	logerr = log.New(os.Stderr, "", 0)
	)

const (
	// 1 MB limit
	MaxBodySize = 1 * 1024 * 1024
)

type Config struct {
	maxDepth uint64
	reqPerSec uint64
	targetURI string

	targetDomain string
}

func parseUri(uri string) (string, error){
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	if !re.MatchString(uri){
		return "", errors.New("parseArgs: Corrupt URI")
	}

	dom := re.FindStringSubmatch(uri)[1]
	return dom, nil
}

func parseArgs() (*Config, error){
	var maxDepth, reqPerSec uint64
	var targetURI string

	flag.Uint64Var(&maxDepth, "d", 2, "Maximum depth, 0 - no limit")
	flag.Uint64Var(&reqPerSec, "r", 5, "Requests per second, 0 - no limit")
	flag.StringVar(&targetURI, "u", "", "Target Uri, required")
	flag.Parse()

	if(targetURI == ""){
		return nil, errors.New("parseArgs: Target URI is required")
	}
	targetDomain, err := parseUri(targetURI)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	cfg.maxDepth = maxDepth
	cfg.reqPerSec = reqPerSec
	cfg.targetURI = targetURI
	cfg.targetDomain = targetDomain

	return cfg, nil
}

type Spider struct {
	cfg *Config
}

func (s *Spider) Load(cfg *Config) {
	s.cfg = cfg
}

func (s *Spider) Run() chan string {
	outChan := make(chan string, 200 * 1000)

	go func() {
		defer close(outChan)

		cfg := s.cfg
		dom := cfg.targetDomain

		// "inculude subdomains" regexp
		re_str := `^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:[a-zA-Z0-9]+\.)?` + strings.Replace(dom,".","\\.", -1) + `(?:(?::[0-9]+)?(\/[^\n]*)?)?$`
		re := regexp.MustCompile(re_str)
		c := colly.NewCollector(
			// debug
			//colly.Debugger(&debug.LogDebugger{}),

			// domain + subdomains
			colly.URLFilters(re),

			colly.MaxDepth(int(cfg.maxDepth)),
			colly.MaxBodySize(MaxBodySize),

			colly.Async(true),
			// TODO: better way for caching GET requests, current is not stable
			colly.CacheDir("spidy_cache"),
		)


		delay := 1000 * time.Microsecond
		_ = c.Limit(&colly.LimitRule{
			DomainGlob: 	"*",
			// yeah, stupid implementation of req/sec limit,bad preсision
			Parallelism: 	int(cfg.reqPerSec),
			Delay:			delay,
			// periodic request protection bypass
			// RandomDelay: delay,
		})
		extensions.RandomUserAgent(c)
		extensions.Referer(c)

		c.OnRequest(func(r *colly.Request) {
			outChan <- fmt.Sprint(r.URL)
		})
		c.OnError(func(r *colly.Response, err error) {
			logerr.Println("[!] Warning:", err, r.Request.URL)
		})
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			e.Request.Visit(link)
		})

		c.Visit(cfg.targetURI)
		c.Visit("https://" + dom)
		c.Visit("http://" + dom)

		u, _ := url.Parse(cfg.targetURI)
		host := u.Host
		if host != dom {
			c.Visit("https://" + u.Host)
			c.Visit("http://" + u.Host)
		}


		c.Wait()
	}()

	return outChan
}

func main() {
	cfg, err := parseArgs()
	if err != nil {
		logerr.Fatal(err)
	}


	logerr.Println("[+] Spidy started working")
	spidy := Spider{}
	spidy.Load(cfg)

	c := spidy.Run()

	for link := range c {
		fmt.Println(link)
	}

	logerr.Println("[+] Spidy finished")
}
