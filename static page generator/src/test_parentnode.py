import unittest

from leafnode import LeafNode
from parentnode import ParentNode

class TestParentNode(unittest.TestCase):
  def test_constructor(self):
    child = LeafNode("a", "Click me!", {"href": "https://www.google.com"})
    node = ParentNode("p", child, {"href": "https://www.google.com"})
    self.assertEqual(node.tag, "p")
    self.assertEqual(node.children, child)
    self.assertEqual(node.props, {"href": "https://www.google.com"})

  def test_to_html(self):
    node = ParentNode(
      "p",
      [
        LeafNode("b", "Bold text"),
        LeafNode(None, "Normal text"),
        LeafNode("i", "italic text"),
        LeafNode(None, "Normal text"),
      ],
    )
    self.assertEqual(node.to_html(), '<p><b>Bold text</b>Normal text<i>italic text</i>Normal text</p>')
    
    node = ParentNode(
      "h3",
      [
        ParentNode(
          "p",
          [
            LeafNode("b", "Bold text"),
            LeafNode(None, "Normal text"),
            LeafNode("i", "italic text"),
            LeafNode(None, "Normal text"),
          ],
        ),
        LeafNode(None, "Normal text"),
      ]
    )
    self.assertEqual(node.to_html(), '<h3><p><b>Bold text</b>Normal text<i>italic text</i>Normal text</p>Normal text</h3>')


