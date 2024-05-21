import unittest

from leafnode import LeafNode

class TestLeafNode(unittest.TestCase):
  def test_constructor(self):
    node = LeafNode("a", "Click me!", {"href": "https://www.google.com"})
    self.assertEqual(node.tag, "a")
    self.assertEqual(node.value, "Click me!")
    self.assertEqual(node.props, {"href": "https://www.google.com"})

  def test_to_html(self):
    node = LeafNode("p", "This is a paragraph of text.")
    self.assertEqual(node.to_html(), '<p>This is a paragraph of text.</p>')
    node = LeafNode("blockquote", "This is a paragraph of text.")
    self.assertEqual(node.to_html(), '<blockquote>This is a paragraph of text.</blockquote>')
    node = LeafNode("code", "This is a paragraph of text.")
    self.assertEqual(node.to_html(), '<code>This is a paragraph of text.</code>')
    node = LeafNode("a", "Click me!", {"href": "https://www.google.com"})
    self.assertEqual(node.to_html(), '<a href="https://www.google.com">Click me!</a>')
    node = LeafNode("img", None, {"src": "url/of/image.jpg", "alt": "Description of image"})
    self.assertEqual(node.to_html(), '<img src="url/of/image.jpg" alt="Description of image">')
    node = LeafNode("b", "This is a paragraph of text.")
    self.assertEqual(node.to_html(), '<b>This is a paragraph of text.</b>')
    node = LeafNode("i", "This is a paragraph of text.")
    self.assertEqual(node.to_html(), '<i>This is a paragraph of text.</i>')
    for header in ["h1", "h2", "h3", "h4", "h5", "h6"]:
      node = LeafNode(header, "This is a paragraph of text.")
      self.assertEqual(node.to_html(), f'<{header}>This is a paragraph of text.</{header}>')
    node = LeafNode(value="Normal text")
    self.assertEqual(node.to_html(), 'Normal text')


