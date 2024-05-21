import unittest

from htmlnode import HTMLNode

class TestHTMLNode(unittest.TestCase):
  def test_constructor(self):
    node = HTMLNode(tag="a", props={"href": "https://www.google.com", "target": "_blank"})
    self.assertEqual(node.tag, "a")
    self.assertEqual(node.props, {"href": "https://www.google.com", "target": "_blank"})

  def test_props_to_html(self):
    node = HTMLNode(tag="a", props={"href": "https://www.google.com", "target": "_blank"})
    self.assertEqual(node.props_to_html(), 'href="https://www.google.com" target="_blank"')


