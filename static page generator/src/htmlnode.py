
class HTMLNode:
  def __init__(self, tag=None, value=None, children=None, props=None):
    self.tag = tag
    self.value = value
    self.children = children
    self.props = props

  def to_html(self):
    raise NotImplementedError

  def props_to_html(self):
    if self.tag == "a":
      href = self.props.get("href")
      target = self.props.get("target")
      return f'href="{href}" target="{target}"'
    
    raise NotImplementedError
    
  def __repr__(self):
    return f"HTMLNode(tag: {self.tag}, value: {self.value}, children: {self.children}, props: {self.props})"
