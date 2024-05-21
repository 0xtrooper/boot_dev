from textnode import TextNode
from leafnode import LeafNode
from parentnode import ParentNode


def main():
  t = TextNode("This is a text node", "bold", "https://www.boot.dev")
  print(t)

if __name__ == "__main__":
  main() 