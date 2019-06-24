# Hypatia
An API aggregator for Beat's microservices (synchronous and asynchronous APIs)

Hypatia is a GitHub scraper that:
- Every `REFRESH_TIME` minutes
  - Gets a summary list of  `GITHUB_ORGANIZATION`s’ repositories
  - For each of the repositories that contains `GITHUB_TAGS`
    - Gets all synchronous and asynchronous documentation files from branch `GITHUB_BRANCH`
    - Visualises synchronous APIs with the use of [RapiDoc](https://mrin9.github.io/RapiDoc/), asynchronous APIs with the use of [api2html](https://github.com/tobilg/api2html/), and indexes them for the search functionality with [Bleeve](https://github.com/blevesearch/bleve)

##### Steps to integrate your micro service API into Hypatia’s system:
- Write or generate your synchronous API documentation with [OpenApi](https://swagger.io/docs/specification/about/) or  asynchronous documentation with [AsyncApi](https://www.asyncapi.com/).
- (Optional) If the output of the tool you used is a YAML file, use a converter to JSON (both OpenApi and AsyncApi are compatible with the latter).
- Push your documentation under the path `/docs/swagger.json` for synchronous and `/docs/async.json` for asynchronous respectively. (branch: default repository branch)
- Tag your repository eligible for documentation scraping with the topic: `api-doc`.
- Hypatia will scrape Beat’s repositories in intervals, filter on the `api-doc` topic and update the API documentations

There are different ways to generate the documentation files. See [API Documentation - Tools and approaches](https://confluence.taxibeat.com/display/TECH/API+Documentation%3A+Tools+and+approaches)

---
In order to bind the resources run from the main folder
```
go-bindata -o bounddata/bound.go -pkg bounddata -fs -prefix "static/" static/...
```

you need go-binddata for it

In order to generate the templates go into the templates folder and run 
```cassandraql
hero .
```
you need hero template engine for it

---

Hypatia (born c. 350–370; died 415 AD)[[1]](https://books.google.com/books?id=79OvkQEACAAJ)[[2]](http://www-history.mcs.st-andrews.ac.uk/Biographies/Hypatia.html)
was a Hellenistic Neoplatonist philosopher, astronomer, and mathematician, who lived in Alexandria, Egypt, then part of
the Eastern Roman Empire. She was a prominent thinker of the Neoplatonic school in Alexandria where she taught
philosophy and astronomy. She is the first female mathematician whose life is reasonably well recorded.

In the twentieth century, Hypatia became seen as an icon for women's rights and a precursor to the feminist movement.
Since the late twentieth century, some portrayals have associated Hypatia's death with the destruction of the Library of
 Alexandria, despite the historical fact that the library no longer existed during Hypatia's lifetime.[[3]](https://books.google.com/books?id=3QPWDAAAQBAJ&pg=PA183)
 
 _Source_: [Wikipedia](https://en.wikipedia.org/wiki/Hypatia)

 
