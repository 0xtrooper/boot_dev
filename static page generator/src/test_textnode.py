import unittest

from textnode import TextNode, split_nodes_delimiter, extract_markdown_images, extract_markdown_links, split_nodes_image, split_nodes_link

class TestTextNode(unittest.TestCase):
  def test_constructor(self):
    node = TextNode("This is a text node", "bold", "https://www.boot.dev")
    self.assertEqual(node.text, "This is a text node")
    self.assertEqual(node.text_type, "bold")
    self.assertEqual(node.url, "https://www.boot.dev")

  def test_eq(self):
    node = TextNode("This is a text node", "bold")
    node2 = TextNode("This is a text node", "bold")
    self.assertEqual(node, node2)

  def test_repr(self):
    node = TextNode("This is a text node", "bold")
    self.assertEqual(repr(node), "TextNode(This is a text node, bold, None)")
    node2 = TextNode("This is a text node", "bold", "https://www.boot.dev")
    self.assertEqual(repr(node2), "TextNode(This is a text node, bold, https://www.boot.dev)")

  def test_split_nodes_delimiter(self):
    node = TextNode("This is text with a `code block` word", "text")
    new_nodes = split_nodes_delimiter([node], "`", "code")
    self.assertEqual(new_nodes[0], TextNode("This is text with a ", "text"))
    self.assertEqual(new_nodes[1], TextNode("code block", "code"))
    self.assertEqual(new_nodes[2], TextNode(" word", "text"))

    node = TextNode("This is text with a `code block word", "text")
    with self.assertRaises(SyntaxError) as context:
      split_nodes_delimiter([node], "`", "code")

    self.assertTrue("Invalid markdown, formatted section not closed" in str(context.exception))

  def test_extract_markdown_images(self):
    text = "This is text with an ![image](https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/zjjcJKZ.png) and ![another](https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/dfsdkjfd.png), this is not an img: [link](https://www.example.com)"
    imgs = extract_markdown_images(text)
    self.assertEqual(imgs[0][0], "image")
    self.assertEqual(imgs[0][1], "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/zjjcJKZ.png")
    self.assertEqual(imgs[1][0], "another")
    self.assertEqual(imgs[1][1], "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/dfsdkjfd.png") 
    self.assertEqual(len(imgs), 2)

  def test_extract_markdown_links(self):
    text = "This is text with a [link](https://www.example.com) and [another](https://www.example.com/another), this is not: ![image](https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/zjjcJKZ.png)"
    imgs = extract_markdown_links(text)
    self.assertEqual(imgs[0][0], "link")
    self.assertEqual(imgs[0][1], "https://www.example.com")
    self.assertEqual(imgs[1][0], "another")
    self.assertEqual(imgs[1][1], "https://www.example.com/another")
    self.assertEqual(len(imgs), 2)

  def test_split_nodes_image(self):
    node = TextNode("This is text with an ![image](https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/zjjcJKZ.png) and another ![second image](https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/3elNhQu.png)", "text")
    new_nodes = split_nodes_image([node])
    self.assertEqual(new_nodes[0], TextNode("This is text with an ", "text"))
    self.assertEqual(new_nodes[1], TextNode("image", "img", "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/zjjcJKZ.png"))
    self.assertEqual(new_nodes[2], TextNode(" and another ", "text"))
    self.assertEqual(new_nodes[3], TextNode("second image", "img", "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/3elNhQu.png"))

  def test_split_nodes_link(self):
    node = TextNode("This is text with a [link](https://www.example.com) and [another](https://www.example.com/another)", "text")
    new_nodes = split_nodes_link([node])
    self.assertEqual(new_nodes[0], TextNode("This is text with a ", "text"))
    self.assertEqual(new_nodes[1], TextNode("link", "link", "https://www.example.com"))
    self.assertEqual(new_nodes[2], TextNode(" and ", "text"))
    self.assertEqual(new_nodes[3], TextNode("another", "link", "https://www.example.com/another"))


if __name__ == "__main__":
  unittest.main()