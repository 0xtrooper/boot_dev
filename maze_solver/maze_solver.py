from graphics import Window
from maze import Maze
import sys


def main():
  max_width = 1200
  max_height = 600
  cell_size = 25
  num_cols = int(max_height / cell_size)
  num_rows = int(max_width / cell_size)
  margin = 5

  old_limit = sys.getrecursionlimit()
  new_limit = old_limit * 5
  sys.setrecursionlimit(new_limit)
  print(f"incresed recursion depth from {old_limit} to {new_limit}")

  win = Window(num_rows * cell_size+2*margin, num_cols * cell_size+2*margin)
  m1 = Maze(margin, margin, num_rows, num_cols, cell_size, cell_size, win)
  m1.break_walls()
  found_solution = m1.solve()
  if found_solution:
    print("Maze solved")
  else:
    print("Maze not solved")
  win.wait_for_close()

if __name__ == "__main__":
  main()