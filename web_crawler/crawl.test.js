import { test, expect } from "@jest/globals";
import { normalizeURL, removeRepeating, getURLsFromHTML} from "./crawl.js";

test('normalizeURL', () => {  
  expect(normalizeURL("https://blog.boot.dev/path/")).toBe("blog.boot.dev/path")
  expect(normalizeURL("https://blog.boot.dev/path")).toBe("blog.boot.dev/path")
  expect(normalizeURL("http://blog.boot.dev/path/")).toBe("blog.boot.dev/path")
  expect(normalizeURL("http://blog.boot.dev/path")).toBe("blog.boot.dev/path")
  expect(() => normalizeURL("http:/blog.boot.dev/path")).toThrow("bad http url")
});

test('getURLsFromHTML', () => {
  const testCase1 = "<html><body><a href=\"https://blog.boot.dev/path\"><span>Go to Boot.dev</span></a></body></html>"
  const testCase2 = "<html><body><a href=\"/path\"><span>Go to Boot.dev</span></a></body></html>"
  expect(getURLsFromHTML(testCase1, "https://blog.boot.dev")).toEqual(["https://blog.boot.dev/path"])
  expect(getURLsFromHTML(testCase2, "https://blog.boot.dev")).toEqual(["https://blog.boot.dev/path"])  
})