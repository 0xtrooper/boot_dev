export function printReport(pages) {
  console.log('--------------------------------\nReport is starting:');
  const pagesArray = Object.entries(pages);
  pagesArray.sort((a, b) => b[1] - a[1]);

  pagesArray.forEach(([url, count]) => {
    console.log(`\tFound ${count} internal links to ${url}`);
  });
}

export function printReportBestPages(pages, n=20) {
  console.log(`--------------------------------\nReport of the best ${n} is starting:`);
  const pagesArray = Object.entries(pages);
  pagesArray.sort((a, b) => b[1] - a[1]);

  for (let i = 0; i < n; i++) {
    console.log(`Found ${pagesArray[i][1]} internal links to ${pagesArray[i][0]}`);
  }
}