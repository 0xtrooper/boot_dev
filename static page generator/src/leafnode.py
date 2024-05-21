from htmlnode import HTMLNode

class LeafNode(HTMLNode):  
  def __init__(self, tag=None, value=None, props=None):
    super().__init__(tag=tag, value=value, props=props)
    
  def to_html(self):
    if self.tag == "img":
      return f'<img src="{self.props.get("src")}" alt="{self.props.get("alt")}">'
      
    if self.value == None:
      raise ValueError
    
    if self.tag == "p":
      return f"<p>{self.value}</p>"
    if self.tag == "code":
      return f"<code>{self.value}</code>"
    if self.tag == "blockquote":
      return f"<blockquote>{self.value}</blockquote>"
    if self.tag == "a":
      return f'<a href="{self.props.get("href")}">{self.value}</a>'
    if self.tag == "b":
      return f"<b>{self.value}</b>"
    if self.tag == "i":
      return f"<i>{self.value}</i>"
    if self.tag in ["h1", "h2", "h3", "h4", "h5", "h6"]:
      return f"<{self.tag}>{self.value}</{self.tag}>"
      
    # defualt case: return normal text
    return self.value