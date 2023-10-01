package book

import (
	"testing"
	"time"
)

func TestBody(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Body()
	want := "<!doctype html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"utf-8\">\n\n  <title>The Twelve-Factor App </title>\n  <meta name=\"description\" content=\"A methodology for building modern, scalable, maintainable software-as-a-service apps.\">\n  <meta name=\"author\" content=\"Adam Wiggins\">\n\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n  <link rel=\"shortcut icon\" href=\"images/favicon.ico\">\n\n  <link rel=\"stylesheet\" href=\"/css/screen.css\" media=\"screen\">\n  <link rel=\"stylesheet\" href=\"/css/mobile.css\" media=\"screen\">\n\n  <script type=\"text/javascript\" src=\"//use.typekit.com/rsq7tro.js\"></script>\n  <script type=\"text/javascript\">try{Typekit.load();}catch(e){}</script>\n</head>\n<body>\n  \n\n  <header>\n    <h1><a href=\"./\" title=\"The Twelve-Factor App\">The Twelve-Factor App</a></h1>\n  </header>\n\n  <section class=\"abstract\">\n  <article>\n<h1 id=\"introduction\">Introduction</h1>\n\n<p>In the modern era, software is commonly delivered as a service: called <em>web apps</em>, or <em>software-as-a-service</em>. The twelve-factor app is a methodology for building software-as-a-service apps that:</p>\n\n<ul>\n<li>Use <strong>declarative</strong> formats for setup automation, to minimize time and cost for new developers joining the project;</li>\n\n<li>Have a <strong>clean contract</strong> with the underlying operating system, offering <strong>maximum portability</strong> between execution environments;</li>\n\n<li>Are suitable for <strong>deployment</strong> on modern <strong>cloud platforms</strong>, obviating the need for servers and systems administration;</li>\n\n<li><strong>Minimize divergence</strong> between development and production, enabling <strong>continuous deployment</strong> for maximum agility;</li>\n\n<li>And can <strong>scale up</strong> without significant changes to tooling, architecture, or development practices.</li>\n</ul>\n\n<p>The twelve-factor methodology can be applied to apps written in any programming language, and which use any combination of backing services (database, queue, memory cache, etc).</p>\n</article>\n  <article>\n<h1 id=\"background\">Background</h1>\n\n<p>The contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the <a href='http://www.heroku.com/' target='_blank'>Heroku</a> platform.</p>\n\n<p>This document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and <a href='http://blog.heroku.com/archives/2011/6/28/the_new_heroku_4_erosion_resistance_explicit_contracts/' target='_blank'>avoiding the cost of software erosion</a>.</p>\n\n<p>Our motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books <em><a href='https://books.google.com/books/about/Patterns_of_enterprise_application_archi.html?id=FyWZt5DdvFkC' target='_blank'>Patterns of Enterprise Application Architecture</a></em> and <em><a href='https://books.google.com/books/about/Refactoring.html?id=1MsETFPD3I0C' target='_blank'>Refactoring</a></em>.</p>\n</article>\n  <article>\n<h1 id=\"who_should_read_this_document\">Who should read this document?</h1>\n\n<p>Any developer building applications which run as a service. Ops engineers who deploy or manage such applications.</p>\n</article>\n</section>\n\n<section class=\"concrete\">\n  <article>\n<h1 id=\"the_twelve_factors\">The Twelve Factors</h1>\n\n<h2 id=\"i_codebase\"><a href=\"./codebase\">I. Codebase</a></h2>\n\n<h3 id=\"one_codebase_tracked_in_revision_control_many_deploys\">One codebase tracked in revision control, many deploys</h3>\n\n<h2 id=\"ii_dependencies\"><a href=\"./dependencies\">II. Dependencies</a></h2>\n\n<h3 id=\"explicitly_declare_and_isolate_dependencies\">Explicitly declare and isolate dependencies</h3>\n\n<h2 id=\"iii_config\"><a href=\"./config\">III. Config</a></h2>\n\n<h3 id=\"store_config_in_the_environment\">Store config in the environment</h3>\n\n<h2 id=\"iv_backing_services\"><a href=\"./backing-services\">IV. Backing services</a></h2>\n\n<h3 id=\"treat_backing_services_as_attached_resources\">Treat backing services as attached resources</h3>\n\n<h2 id=\"v_build_release_run\"><a href=\"./build-release-run\">V. Build, release, run</a></h2>\n\n<h3 id=\"strictly_separate_build_and_run_stages\">Strictly separate build and run stages</h3>\n\n<h2 id=\"vi_processes\"><a href=\"./processes\">VI. Processes</a></h2>\n\n<h3 id=\"execute_the_app_as_one_or_more_stateless_processes\">Execute the app as one or more stateless processes</h3>\n\n<h2 id=\"vii_port_binding\"><a href=\"./port-binding\">VII. Port binding</a></h2>\n\n<h3 id=\"export_services_via_port_binding\">Export services via port binding</h3>\n\n<h2 id=\"viii_concurrency\"><a href=\"./concurrency\">VIII. Concurrency</a></h2>\n\n<h3 id=\"scale_out_via_the_process_model\">Scale out via the process model</h3>\n\n<h2 id=\"ix_disposability\"><a href=\"./disposability\">IX. Disposability</a></h2>\n\n<h3 id=\"maximize_robustness_with_fast_startup_and_graceful_shutdown\">Maximize robustness with fast startup and graceful shutdown</h3>\n\n<h2 id=\"x_devprod_parity\"><a href=\"./dev-prod-parity\">X. Dev/prod parity</a></h2>\n\n<h3 id=\"keep_development_staging_and_production_as_similar_as_possible\">Keep development, staging, and production as similar as possible</h3>\n\n<h2 id=\"xi_logs\"><a href=\"./logs\">XI. Logs</a></h2>\n\n<h3 id=\"treat_logs_as_event_streams\">Treat logs as event streams</h3>\n\n<h2 id=\"xii_admin_processes\"><a href=\"./admin-processes\">XII. Admin processes</a></h2>\n\n<h3 id=\"run_adminmanagement_tasks_as_oneoff_processes\">Run admin/management tasks as one-off processes</h3>\n</article>\n</section>\n\n\n  <footer>\n  <div id=\"locales\"><a href=\"/cs/\">Česky (cs)</a> | <a href=\"/de/\">Deutsch (de)</a> | <a href=\"/el/\">Ελληνικά (el)</a> | <span>English (en)</span> | <a href=\"/es/\">Español (es)</a> | <a href=\"/fr/\">Français (fr)</a> | <a href=\"/it/\">Italiano (it)</a> | <a href=\"/ja/\">日本語 (ja)</a> | <a href=\"/ko/\">한국어 (ko)</a> | <a href=\"/pl/\">Polski (pl)</a> | <a href=\"/pt_br/\">Brazilian Portuguese (pt_br)</a> | <a href=\"/ru/\">Русский (ru)</a> | <a href=\"/sk/\">Slovensky (sk)</a> | <a href=\"/th/\">ภาษาไทย (th)</a> | <a href=\"/tr/\">Turkish (tr)</a> | <a href=\"/uk/\">Українська (uk)</a> | <a href=\"/vi/\">Tiếng Việt (vi)</a> | <a href=\"/zh_cn/\">简体中文 (zh_cn)</a></div>\n  <div>Written by Adam Wiggins</div>\n  <div>Last updated 2017</div>\n  <div><a href=\"https://github.com/heroku/12factor\">Sourcecode</a></div>\n  <div><a href=\"/12factor.epub\">Download ePub Book</a></div>\n  <div><a href=\"https://www.heroku.com/policy/privacy\">Privacy Policy</a></div>\n  <div class=\"cpra\"><a rel=\"nofollow\" href=\"https://www.salesforce.com/form/other/privacy-request/\"><img src=\"/images/privacy.png\" aria-label=\"Privacy Icon\" title=\"Privacy Icon\">Your Privacy Choices</a></div>\n</footer>\n</body>\n</html>\n"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestName(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "The Twelve-Factor App"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestCustomName(t *testing.T) {

	config := NewScrapeConfig()
	config.UseLinkName = true
	c := NewChapterFromURL("https://12factor.net/", "Custom Name", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "Custom Name"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestAuthor(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Author()
	want := "Adam Wiggins"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestContent(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "\n  \n\n  <header>\n    \n  </header>\n\n  <section>\n  <article>\n\n\n<p>In the modern era, software is commonly delivered as a service: called <em>web apps</em>, or <em>software-as-a-service</em>. The twelve-factor app is a methodology for building software-as-a-service apps that:</p>\n\n<ul>\n<li>Use <strong>declarative</strong> formats for setup automation, to minimize time and cost for new developers joining the project;</li>\n\n<li>Have a <strong>clean contract</strong> with the underlying operating system, offering <strong>maximum portability</strong> between execution environments;</li>\n\n<li>Are suitable for <strong>deployment</strong> on modern <strong>cloud platforms</strong>, obviating the need for servers and systems administration;</li>\n\n<li><strong>Minimize divergence</strong> between development and production, enabling <strong>continuous deployment</strong> for maximum agility;</li>\n\n<li>And can <strong>scale up</strong> without significant changes to tooling, architecture, or development practices.</li>\n</ul>\n\n<p>The twelve-factor methodology can be applied to apps written in any programming language, and which use any combination of backing services (database, queue, memory cache, etc).</p>\n</article>\n  <article>\n\n\n<p>The contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the <a href=\"http://www.heroku.com/\" target=\"_blank\">Heroku</a> platform.</p>\n\n<p>This document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and <a href=\"http://blog.heroku.com/archives/2011/6/28/the_new_heroku_4_erosion_resistance_explicit_contracts/\" target=\"_blank\">avoiding the cost of software erosion</a>.</p>\n\n<p>Our motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books <em><a href=\"https://books.google.com/books/about/Patterns_of_enterprise_application_archi.html?id=FyWZt5DdvFkC\" target=\"_blank\">Patterns of Enterprise Application Architecture</a></em> and <em><a href=\"https://books.google.com/books/about/Refactoring.html?id=1MsETFPD3I0C\" target=\"_blank\">Refactoring</a></em>.</p>\n</article>\n  <article>\n\n\n<p>Any developer building applications which run as a service. Ops engineers who deploy or manage such applications.</p>\n</article>\n</section>\n\n<section>\n  <article>\n\n\n<h2 id=\"i_codebase\"><a href=\"https://12factor.net/codebase\">I. Codebase</a></h2>\n\n<h3 id=\"one_codebase_tracked_in_revision_control_many_deploys\">One codebase tracked in revision control, many deploys</h3>\n\n<h2 id=\"ii_dependencies\"><a href=\"https://12factor.net/dependencies\">II. Dependencies</a></h2>\n\n<h3 id=\"explicitly_declare_and_isolate_dependencies\">Explicitly declare and isolate dependencies</h3>\n\n<h2 id=\"iii_config\"><a href=\"https://12factor.net/config\">III. Config</a></h2>\n\n<h3 id=\"store_config_in_the_environment\">Store config in the environment</h3>\n\n<h2 id=\"iv_backing_services\"><a href=\"https://12factor.net/backing-services\">IV. Backing services</a></h2>\n\n<h3 id=\"treat_backing_services_as_attached_resources\">Treat backing services as attached resources</h3>\n\n<h2 id=\"v_build_release_run\"><a href=\"https://12factor.net/build-release-run\">V. Build, release, run</a></h2>\n\n<h3 id=\"strictly_separate_build_and_run_stages\">Strictly separate build and run stages</h3>\n\n<h2 id=\"vi_processes\"><a href=\"https://12factor.net/processes\">VI. Processes</a></h2>\n\n<h3 id=\"execute_the_app_as_one_or_more_stateless_processes\">Execute the app as one or more stateless processes</h3>\n\n<h2 id=\"vii_port_binding\"><a href=\"https://12factor.net/port-binding\">VII. Port binding</a></h2>\n\n<h3 id=\"export_services_via_port_binding\">Export services via port binding</h3>\n\n<h2 id=\"viii_concurrency\"><a href=\"https://12factor.net/concurrency\">VIII. Concurrency</a></h2>\n\n<h3 id=\"scale_out_via_the_process_model\">Scale out via the process model</h3>\n\n<h2 id=\"ix_disposability\"><a href=\"https://12factor.net/disposability\">IX. Disposability</a></h2>\n\n<h3 id=\"maximize_robustness_with_fast_startup_and_graceful_shutdown\">Maximize robustness with fast startup and graceful shutdown</h3>\n\n<h2 id=\"x_devprod_parity\"><a href=\"https://12factor.net/dev-prod-parity\">X. Dev/prod parity</a></h2>\n\n<h3 id=\"keep_development_staging_and_production_as_similar_as_possible\">Keep development, staging, and production as similar as possible</h3>\n\n<h2 id=\"xi_logs\"><a href=\"https://12factor.net/logs\">XI. Logs</a></h2>\n\n<h3 id=\"treat_logs_as_event_streams\">Treat logs as event streams</h3>\n\n<h2 id=\"xii_admin_processes\"><a href=\"https://12factor.net/admin-processes\">XII. Admin processes</a></h2>\n\n<h3 id=\"run_adminmanagement_tasks_as_oneoff_processes\">Run admin/management tasks as one-off processes</h3>\n</article>\n</section>\n\n\n  \n\n\n"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestDelay(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Delay = 500

	config1 := NewScrapeConfig()

	start := time.Now()
	NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})
	elapsed := time.Since(start)

	got := elapsed
	want := time.Duration(500) * time.Millisecond

	if got < want {
		t.Errorf("got %v, wanted min %v", got, want)
	}

}

func TestContentImagesOnly(t *testing.T) {

	config := NewScrapeConfig()
	config.ImagesOnly = true

	c := NewChapterFromURL("https://12factor.net/codebase", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "<img src=\"https://12factor.net/images/codebase-deploys.png\" alt=\"One codebase maps to many deploys\"/>"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChapters(t *testing.T) {

	config0 := NewScrapeConfig()
	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://atomicdesign.bradfrost.com/table-of-contents/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 9

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersRSS(t *testing.T) {

	config0 := NewScrapeConfig()
	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://www.nginx.com/feed/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 10

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersSelector(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Selector = "section.concrete > article > h2 > a"

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimit(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Selector = "section.concrete > article > h2 > a"
	config0.Limit = 2

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 2

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimitOver(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Selector = "section.concrete > article > h2 > a"
	config0.Limit = 13

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestReverse(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Reverse = true

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://atomicdesign.bradfrost.com/table-of-contents/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := c.SubChapters()[0].Name()
	want := "About the Author | Atomic Design by Brad Frost"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestNotInclude(t *testing.T) {

	config := NewScrapeConfig()
	config.Selector = "section.concrete > article > h2 > a"
	config.Include = false

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := ""

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}
