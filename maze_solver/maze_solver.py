from graphics import Window
from maze import Maze


def main():
  max_width = 1300
  max_height = 650
  cell_size = 10
  num_cols = int(max_height / cell_size)
  num_rows = int(max_width / cell_size)
  margin = 5

  win = Window(num_rows * cell_size+2*margin, num_cols * cell_size+2*margin)
  m1 = Maze(margin, margin, num_rows, num_cols, cell_size, cell_size, win)
  m1.break_walls()
  # found_solution = m1.solve()
  found_solution = m1.solve_alpha_start()
  if found_solution:
    print("Maze solved")
  else:
    print("Maze not solved")
  win.wait_for_close()

if __name__ == "__main__":
  main()