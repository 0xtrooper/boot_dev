from leafnode import LeafNode
import re

class TextNode:
  def __init__(self, text, text_type, url=None):
    self.text = text
    self.text_type = text_type
    self.url = url

  def __eq__(self, other):
    return self.text == other.text and self.text_type == other.text_type and self.url == other.url

  def __repr__(self):
    return f"TextNode({self.text}, {self.text_type}, {self.url})"
  
def to_html_node(text_node):
  if text_node.text_type == "text":
    return LeafNode(value=text_node.text)    
  if text_node.text_type == "bold":
    return LeafNode(value=text_node.text, tag="b")
  if text_node.text_type == "italic":
    return LeafNode(value=text_node.text, tag="i")
  if text_node.text_type == "code":
    return LeafNode(value=text_node.text, tag="code")
  if text_node.text_type == "link":
    return LeafNode(value=text_node.text, tag="a", props={"href": text_node.url})
  if text_node.text_type == "image":
    return LeafNode(tag="img", props={"src": text_node.url, "alt": text_node.text})
  
  raise NotImplementedError(f"to_html_node: {text_node.text_type} not implemented")
  
def split_nodes_delimiter(old_nodes, delimiter, text_type):
  new_nodes = []
  for old_node in old_nodes:
    split_nodes = []
    sections = old_node.text.split(delimiter)
    if len(sections) % 2 == 0:
      raise SyntaxError("Invalid markdown, formatted section not closed")
    for i in range(len(sections)):
      if sections[i] == "":
        continue
      if i % 2 == 0:
        split_nodes.append(TextNode(sections[i], old_node.text_type))
      else:
        split_nodes.append(TextNode(sections[i], text_type))
    new_nodes.extend(split_nodes)
  return new_nodes

def split_nodes_image(old_nodes):
  new_nodes = []

  for old_node in old_nodes:
    split_nodes = []

    imgs = extract_markdown_images(old_node.text)

    temp_text = old_node.text
    for img in imgs:
      old_part = temp_text.split(f"![{img[0]}]({img[1]})", maxsplit=1)
      split_nodes.append(TextNode(old_part[0], old_node.text_type))
      split_nodes.append(TextNode(img[0], "img", img[1]))
      temp_text = old_part[1]

    new_nodes.extend(split_nodes)    

  return new_nodes

def split_nodes_link(old_nodes):
  new_nodes = []

  for old_node in old_nodes:
    split_nodes = []

    links = extract_markdown_links(old_node.text)

    temp_text = old_node.text
    for link in links:
      old_part = temp_text.split(f"[{link[0]}]({link[1]})", maxsplit=1)
      split_nodes.append(TextNode(old_part[0], old_node.text_type))
      split_nodes.append(TextNode(link[0], "link", link[1]))
      temp_text = old_part[1]

    new_nodes.extend(split_nodes)    

  return new_nodes

def extract_markdown_images(text):
  return re.findall(r"!\[(.*?)\]\((.*?)\)", text)

def extract_markdown_links(text):
  return re.findall(r"[^!]\[(.*?)\]\((.*?)\)", text)




    