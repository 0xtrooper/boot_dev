from htmlnode import HTMLNode

class ParentNode(HTMLNode):
  def __init__(self, tag, children, props=None):
    super().__init__(tag=tag, children=children, props=props)
    
  def to_html(self):
    if self.tag == None:
      raise ValueError("no tag set")
    
    if self.children == None:
      raise ValueError("no children set")
    
    html = f"<{self.tag}>"

    for c in self.children:
      html = html + c.to_html()

    html = f"{html}</{self.tag}>"
    return html
