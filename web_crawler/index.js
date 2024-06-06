import { argv } from  'node:process'
import { crawlPage } from "./crawl.js"
import { printReport, printReportBestPages } from "./report.js"

async function main() {
  if (argv.length < 3) throw new Error("no url. Run the web_crawler with one url")
  if (argv.length > 3) throw new Error("no url. Run the web_crawler with ONLY one url")

  let url = argv[2]
  if (!url.startsWith("http://") && !url.startsWith("https://")) url = "https://" + url
  console.log(`Starting web crawler at ${url}`)

  const linkMap = await crawlPage(url)
  printReportBestPages(linkMap)
}

await main()