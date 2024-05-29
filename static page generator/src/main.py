import os
import shutil
from textnode import TextNode


def main():
  copy_static_to_public()


def _purge_public():
  shutil.rmtree("../public")
  os.mkdir("../public")

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
  _purge_public()
  _copy_content("../static", "../public")


if __name__ == "__main__":
  main() 