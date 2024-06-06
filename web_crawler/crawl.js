import { JSDOM } from 'jsdom'

export function normalizeURL(rawURL) {
  let startIndex = rawURL.indexOf("://")
  if (startIndex == -1) throw new Error("bad http url: ", rawURL)
  startIndex += 3
  let endIndex = rawURL.length
  while(rawURL[endIndex-1] == "/") {
    endIndex = endIndex - 1
  }
  return rawURL.slice(startIndex, endIndex)
}

export function getURLsFromHTML(htmlBody, baseURL) {
  const jsdom = new JSDOM(htmlBody)
  // returns: "<a href="https://boot.dev">Learn Backend Development</a>""
  const linkAnchors = jsdom.window.document.querySelectorAll('a')

  const hrefs = []
  for (let linkAnchor of linkAnchors) {
    hrefs.push(linkAnchor.getAttribute('href'))
  }

  // expand relative links
  const links = hrefs.map(href => href[0] == "/" ? baseURL + href: href)
  
  return links
}

async function getReaderFromPage(url) {
  const resp = await fetch(url, { method: "GET" })

  if (resp.status > 400) throw new Error(`failed to get url: ${resp.status}`)

  const contentType = resp.headers.get("content-type")
  if (!contentType.includes("text/html")) throw new Error(`page is not a website: ${contentType}`)

  return resp.body.getReader()
}

async function getURLsFromPage(url) {
  const reader = await getReaderFromPage(url)
  const textDecoder = new TextDecoder()

  let buffer
  let page = ""
  do {
    buffer = await reader.read()
    page += textDecoder.decode(buffer.value)    
  } while (!buffer?.done)

  return getURLsFromHTML(page, url)
}

export async function crawlPage(startURL, onlyCrawlPage=true, maxCalls=1000) {
  const linksFound = {[normalizeURL(startURL)]: 1}
  const toVisit = [startURL]
  while(toVisit.length > 0 && maxCalls > 0) {
    const url = toVisit.shift()
    let newURLs = []
    try {
      newURLs = await getURLsFromPage(url)
    } catch(err) {
      console.log(`Failed to fetch "${url}: ${err.message}`)
      continue
    }

    console.log(`Fetched ${url}, ${toVisit.length < maxCalls ? toVisit.length : maxCalls} remaining`)
    maxCalls--

    for(let newURL of newURLs) {
      try {
        const normalizedURL = normalizeURL(newURL)
    
        if(linksFound[normalizedURL]) {
          linksFound[normalizedURL]++
          continue
        }
  
        linksFound[normalizedURL] = 1
  
        // only append:
        // - not visited yet (early exit check)
        // - NOT onlyCrawlPage
        // - onlyCrawlPage AND sub page of start url
        if (maxCalls > toVisit.length &&  (!onlyCrawlPage || newURL.includes(startURL))) {
          toVisit.push(newURL)
        }
      } catch(err) {
        console.log(`Failed to normalize ${newURL}, ${err.message}`)
      }
    }
  }

  return linksFound
}