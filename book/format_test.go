package book

import (
	"errors"
	"os"
	"testing"
)

func TestFilename(t *testing.T) {

	got := Filename("This is a chapter / book")
	want := "This_is_a_chapter__book"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToMarkdownString(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	got := ToMarkdownString(c)
	want := "The Twelve-Factor App\n=====================\n\nIn the modern era, software is commonly delivered as a service: called _web apps_, or _software-as-a-service_. The twelve-factor app is a methodology for building software-as-a-service apps that:\n\n- Use **declarative** formats for setup automation, to minimize time and cost for new developers joining the project;\n- Have a **clean contract** with the underlying operating system, offering **maximum portability** between execution environments;\n- Are suitable for **deployment** on modern **cloud platforms**, obviating the need for servers and systems administration;\n- **Minimize divergence** between development and production, enabling **continuous deployment** for maximum agility;\n- And can **scale up** without significant changes to tooling, architecture, or development practices.\n\nThe twelve-factor methodology can be applied to apps written in any programming language, and which use any combination of backing services (database, queue, memory cache, etc).\n\nThe contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the [Heroku](http://www.heroku.com/) platform.\n\nThis document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and [avoiding the cost of software erosion](http://blog.heroku.com/archives/2011/6/28/the_new_heroku_4_erosion_resistance_explicit_contracts/).\n\nOur motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books _[Patterns of Enterprise Application Architecture](https://books.google.com/books/about/Patterns_of_enterprise_application_archi.html?id=FyWZt5DdvFkC)_ and _[Refactoring](https://books.google.com/books/about/Refactoring.html?id=1MsETFPD3I0C)_.\n\nAny developer building applications which run as a service. Ops engineers who deploy or manage such applications.\n\n## [I. Codebase](https://12factor.net/codebase)\n\n### One codebase tracked in revision control, many deploys\n\n## [II. Dependencies](https://12factor.net/dependencies)\n\n### Explicitly declare and isolate dependencies\n\n## [III. Config](https://12factor.net/config)\n\n### Store config in the environment\n\n## [IV. Backing services](https://12factor.net/backing-services)\n\n### Treat backing services as attached resources\n\n## [V. Build, release, run](https://12factor.net/build-release-run)\n\n### Strictly separate build and run stages\n\n## [VI. Processes](https://12factor.net/processes)\n\n### Execute the app as one or more stateless processes\n\n## [VII. Port binding](https://12factor.net/port-binding)\n\n### Export services via port binding\n\n## [VIII. Concurrency](https://12factor.net/concurrency)\n\n### Scale out via the process model\n\n## [IX. Disposability](https://12factor.net/disposability)\n\n### Maximize robustness with fast startup and graceful shutdown\n\n## [X. Dev/prod parity](https://12factor.net/dev-prod-parity)\n\n### Keep development, staging, and production as similar as possible\n\n## [XI. Logs](https://12factor.net/logs)\n\n### Treat logs as event streams\n\n## [XII. Admin processes](https://12factor.net/admin-processes)\n\n### Run admin/management tasks as one-off processes\n\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToMarkdown(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMarkdown(c, "")

	filename := "The_Twelve-Factor_App.md"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMarkdownFilename(t *testing.T) {

	filename := "ebook.md"
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMarkdown(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToHtmlString(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	got := ToHtmlString(c)
	want := "<h1>The Twelve-Factor App</h1>\n  \n\n  <header>\n    \n  </header>\n\n  <section>\n  <article>\n\n\n<p>In the modern era, software is commonly delivered as a service: called <em>web apps</em>, or <em>software-as-a-service</em>. The twelve-factor app is a methodology for building software-as-a-service apps that:</p>\n\n<ul>\n<li>Use <strong>declarative</strong> formats for setup automation, to minimize time and cost for new developers joining the project;</li>\n\n<li>Have a <strong>clean contract</strong> with the underlying operating system, offering <strong>maximum portability</strong> between execution environments;</li>\n\n<li>Are suitable for <strong>deployment</strong> on modern <strong>cloud platforms</strong>, obviating the need for servers and systems administration;</li>\n\n<li><strong>Minimize divergence</strong> between development and production, enabling <strong>continuous deployment</strong> for maximum agility;</li>\n\n<li>And can <strong>scale up</strong> without significant changes to tooling, architecture, or development practices.</li>\n</ul>\n\n<p>The twelve-factor methodology can be applied to apps written in any programming language, and which use any combination of backing services (database, queue, memory cache, etc).</p>\n</article>\n  <article>\n\n\n<p>The contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the <a href=\"http://www.heroku.com/\" target=\"_blank\">Heroku</a> platform.</p>\n\n<p>This document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and <a href=\"http://blog.heroku.com/archives/2011/6/28/the_new_heroku_4_erosion_resistance_explicit_contracts/\" target=\"_blank\">avoiding the cost of software erosion</a>.</p>\n\n<p>Our motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books <em><a href=\"https://books.google.com/books/about/Patterns_of_enterprise_application_archi.html?id=FyWZt5DdvFkC\" target=\"_blank\">Patterns of Enterprise Application Architecture</a></em> and <em><a href=\"https://books.google.com/books/about/Refactoring.html?id=1MsETFPD3I0C\" target=\"_blank\">Refactoring</a></em>.</p>\n</article>\n  <article>\n\n\n<p>Any developer building applications which run as a service. Ops engineers who deploy or manage such applications.</p>\n</article>\n</section>\n\n<section>\n  <article>\n\n\n<h2 id=\"i_codebase\"><a href=\"https://12factor.net/codebase\">I. Codebase</a></h2>\n\n<h3 id=\"one_codebase_tracked_in_revision_control_many_deploys\">One codebase tracked in revision control, many deploys</h3>\n\n<h2 id=\"ii_dependencies\"><a href=\"https://12factor.net/dependencies\">II. Dependencies</a></h2>\n\n<h3 id=\"explicitly_declare_and_isolate_dependencies\">Explicitly declare and isolate dependencies</h3>\n\n<h2 id=\"iii_config\"><a href=\"https://12factor.net/config\">III. Config</a></h2>\n\n<h3 id=\"store_config_in_the_environment\">Store config in the environment</h3>\n\n<h2 id=\"iv_backing_services\"><a href=\"https://12factor.net/backing-services\">IV. Backing services</a></h2>\n\n<h3 id=\"treat_backing_services_as_attached_resources\">Treat backing services as attached resources</h3>\n\n<h2 id=\"v_build_release_run\"><a href=\"https://12factor.net/build-release-run\">V. Build, release, run</a></h2>\n\n<h3 id=\"strictly_separate_build_and_run_stages\">Strictly separate build and run stages</h3>\n\n<h2 id=\"vi_processes\"><a href=\"https://12factor.net/processes\">VI. Processes</a></h2>\n\n<h3 id=\"execute_the_app_as_one_or_more_stateless_processes\">Execute the app as one or more stateless processes</h3>\n\n<h2 id=\"vii_port_binding\"><a href=\"https://12factor.net/port-binding\">VII. Port binding</a></h2>\n\n<h3 id=\"export_services_via_port_binding\">Export services via port binding</h3>\n\n<h2 id=\"viii_concurrency\"><a href=\"https://12factor.net/concurrency\">VIII. Concurrency</a></h2>\n\n<h3 id=\"scale_out_via_the_process_model\">Scale out via the process model</h3>\n\n<h2 id=\"ix_disposability\"><a href=\"https://12factor.net/disposability\">IX. Disposability</a></h2>\n\n<h3 id=\"maximize_robustness_with_fast_startup_and_graceful_shutdown\">Maximize robustness with fast startup and graceful shutdown</h3>\n\n<h2 id=\"x_devprod_parity\"><a href=\"https://12factor.net/dev-prod-parity\">X. Dev/prod parity</a></h2>\n\n<h3 id=\"keep_development_staging_and_production_as_similar_as_possible\">Keep development, staging, and production as similar as possible</h3>\n\n<h2 id=\"xi_logs\"><a href=\"https://12factor.net/logs\">XI. Logs</a></h2>\n\n<h3 id=\"treat_logs_as_event_streams\">Treat logs as event streams</h3>\n\n<h2 id=\"xii_admin_processes\"><a href=\"https://12factor.net/admin-processes\">XII. Admin processes</a></h2>\n\n<h3 id=\"run_adminmanagement_tasks_as_oneoff_processes\">Run admin/management tasks as one-off processes</h3>\n</article>\n</section>\n\n\n  \n\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToHtml(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToHtml(c, "")

	filename := "The_Twelve-Factor_App.html"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToHtmlFilename(t *testing.T) {

	filename := "ebook.html"
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToHtml(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToEpub(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToEpub(c, "")

	filename := "The_Twelve-Factor_App.epub"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToEpubFilename(t *testing.T) {

	filename := "ebook.epub"
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToEpub(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMobi(t *testing.T) {

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMobi(c, "")

	filename := "The_Twelve-Factor_App.mobi"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMobiFilename(t *testing.T) {

	filename := "ebook.mobi"
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMobi(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}
