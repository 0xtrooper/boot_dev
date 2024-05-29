import os
import shutil

from markdown_blocks import markdown_to_html_node

def _purge_dir(path):
  shutil.rmtree(path)
  os.mkdir(path)

def _get_content(path):
  return os.listdir(path)

def _copy_content(src, dst):
  contents = _get_content(src)
  for content in contents:
    pathToSrc = f"{src}/{content}"
    pathToDst = f"{dst}/{content}"
    isFile = os.path.isfile(pathToSrc)
    if isFile:
      shutil.copy(pathToSrc, pathToDst)
    else:
      os.mkdir(pathToDst)
      _copy_content(pathToSrc, pathToDst)

def copy_static_to_public():
  _purge_dir("./public")
  _copy_content("./static", "./public")

def extract_title(markdown):
  lines = markdown.split('\n')
  for i, line in enumerate(lines):
    if line.startswith('# '):
      title = line[2:].strip()
      modified_markdown = '\n'.join(lines[:i] + lines[i+1:])
      return title, modified_markdown
  raise Exception()

def generate_page(from_path, template_path, dest_path):
  print(f"Generating page from {from_path} to {dest_path} using {template_path}")

  src_file = open(from_path, "r")
  markdown = src_file.read()

  template_file = open(template_path, "r")
  template = template_file.read()

  title, markdown_without_title= extract_title(markdown)
  template_with_title = template.replace("{{ Title }}", title)

  html_node = markdown_to_html_node(markdown_without_title)
  html_content = html_node.to_html()
  content = template_with_title.replace("{{ Content }}", html_content)

  dest_file = open(dest_path, "w")
  dest_file.write(content)
  dest_file.close()

def generate_pages_recursive(dir_path_content, template_path, dest_dir_path):
  contents = _get_content(dir_path_content)
  for content in contents:
    content_html = content.replace(".md", ".html")
    pathToSrc = f"{dir_path_content}/{content}"
    pathToDst = f"{dest_dir_path}/{content_html}"
    isFile = os.path.isfile(pathToSrc)
    if isFile:
      generate_page(pathToSrc, template_path, pathToDst)
    else:
      os.mkdir(pathToDst)
      generate_pages_recursive(pathToSrc, template_path, pathToDst)
