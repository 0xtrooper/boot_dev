import time
import random

from graphics import Cell

class Maze:
  def __init__(self, x1, y1, num_rows, num_cols, cell_size_x, cell_size_y, win=None, max_time=10, seed=None):
    self.x1 = x1
    self.y1 = y1
    self.num_rows = num_rows
    self.num_cols = num_cols
    self.cell_size_x = cell_size_x
    self.cell_size_y = cell_size_y
    self.win = win
    self._cells = []
    self.delay = max_time / (self.num_cols * self.num_rows)
    self.delay = 0 if self.delay < 0.001 else self.delay
    if self.delay == 0:
      print("not using delay to allow fast execution")

    random.seed(seed)
    
    self._create_cells()


  def _create_cells(self):
    for c in range(self.num_cols):
      new_row = []
      for r in range(self.num_rows):
        new_row.append(Cell(self.win))
      self._cells.append(new_row)
    
    
    for c in range(self.num_cols):
      for r in range(self.num_rows):
        self._draw_cell(c, r, False)

  def _draw_cell(self, c, r, animate=True):
    if c > self.num_cols:
      raise Exception("col too large")
    if r > self.num_rows:
      raise Exception("row too large")

    cell = self._cells[c][r]
    x1 = self.x1 + self.cell_size_x * r
    x2 = x1 + self.cell_size_x
    y1 = self.y1 + self.cell_size_y * c
    y2 = y1 + self.cell_size_y
    cell.draw(x1, y1, x2, y2)
    if animate:
      self._animate()

  def _draw_move(self, from_c, from_r, to_c, to_r, undo=False, animate=True):
    if from_c > self.num_cols or to_c > self.num_cols:
      raise Exception("col too large")
    if from_r > self.num_rows or to_r > self.num_rows:
      raise Exception("row too large")

    from_cell = self._cells[from_c][from_r]
    to_cell = self._cells[to_c][to_r]
    from_cell.draw_move(to_cell, undo)
    if animate:
      self._animate()

  def _animate(self):
    if self.win is None:
      return
    self.win.redraw()
    if self.delay > 0:
      time.sleep(self.delay)

  def _break_entrance_and_exit(self):
    entrance_cell = self._cells[0][0]
    entrance_cell.has_top_wall = False
    self._draw_cell(0, 0)

    exit_cell = self._cells[self.num_cols-1][self.num_rows-1]
    exit_cell.has_bottom_wall = False
    self._draw_cell(self.num_cols-1, self.num_rows-1, False)

  def break_walls(self):
    self._break_entrance_and_exit()
    self._break_walls_r(0, 0)
    self._reset_cells_visited()

  def _break_walls_r(self, c, r):
    self._cells[c][r].visited = True
    while(True):
      to_visit = []
      if c > 0 and not self._cells[c-1][r].visited:
        to_visit.append("top")
      if c < self.num_cols-1 and not self._cells[c+1][r].visited:
        to_visit.append("bot")
      if r > 0 and not self._cells[c][r-1].visited:
        to_visit.append("left")
      if r < self.num_rows-1 and not self._cells[c][r+1].visited:
        to_visit.append("right")

      if len(to_visit) == 0:
        self._draw_cell(c, r, False)
        return
      
      direction = random.choice(to_visit)
      if direction == "top":
        self._cells[c][r].has_top_wall = False
        self._cells[c-1][r].has_bottom_wall = False
        self._break_walls_r(c-1, r)
      elif direction == "bot":
        self._cells[c][r].has_bottom_wall = False
        self._cells[c+1][r].has_top_wall = False
        self._break_walls_r(c+1, r)
      elif direction == "left":
        self._cells[c][r].has_left_wall = False
        self._cells[c][r-1].has_right_wall = False
        self._break_walls_r(c, r-1)
      elif direction == "right":
        self._cells[c][r].has_right_wall = False
        self._cells[c][r+1].has_left_wall = False
        self._break_walls_r(c, r+1)

    self._animate()

  def _reset_cells_visited(self):
    for c in range(self.num_cols):
      for r in range(self.num_rows):
        self._cells[c][r].visited = False

  def solve(self):
    return self._solve_r(0, 0)

  def _solve_r(self, c, r):
    if c == self.num_cols-1 and r == self.num_rows-1:
      return True

    cell = self._cells[c][r]
    cell.visited = True
    
    to_visit = []
    if not cell.has_bottom_wall:
      to_visit.append((c+1, r))
    if not cell.has_top_wall and c > 0:
      to_visit.append((c-1, r))
    if not cell.has_right_wall:
      to_visit.append((c, r+1))
    if not cell.has_left_wall:
      to_visit.append((c, r-1))    

    for (next_c, next_r) in to_visit:
      if self._cells[next_c][next_r].visited:
        continue
      
      self._draw_move(c, r, next_c, next_r)
      found = self._solve_r(next_c, next_r)
      if found:
        return True
      else:
        self._draw_move(c, r, next_c, next_r, True)
      
    return False

