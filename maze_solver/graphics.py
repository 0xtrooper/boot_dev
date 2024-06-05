from tkinter import Tk, BOTH, Canvas

class Window:

  def __init__(self, w, h):
    self.width = w
    self.height = h
    self.root_widget = Tk()
    self.root_widget.title("Maze Solver")
    self.root_widget.protocol("WM_DELETE_WINDOW", self.close)
    self.root_widget.configure(background="white")
    self.root_widget.resizable(0,0)
    self.root_widget.maxsize(self.width, self.height)
    self.root_widget.minsize(self.width, self.height)
    self.canvas = Canvas(self.root_widget, width=self.width, height=self.height, background="white")
    self.canvas.pack()
    self.isRunning = False

  def redraw(self):
    self.root_widget.update_idletasks()
    self.root_widget.update()

  def wait_for_close(self):
    self.isRunning = True
    while(self.isRunning):
      self.redraw()

  def close(self):
    self.isRunning = False

  def draw_line(self, line, fill_color):
    line.draw(self.canvas, fill_color)

class Point:
  def __init__(self, x=0, y=0):
    self.x = x
    self.y = y

  def __repr__(self):
    return f"Point({self.x}, {self.y})"

class Line:
  def __init__(self, p1, p2, w=2):
    self.point_one = p1
    self.point_two = p2
    self.width = w

  def __repr__(self):
    return f"Line from {self.point_one} to {self.point_two}, width {self.width}"

  def draw(self, canvas, fill_color):
    if fill_color != "red" and fill_color != "black" and fill_color != "white" and fill_color != "gray":
      raise Exception("bad color, only red and black")
    canvas.create_line(self.point_one.x, self.point_one.y, self.point_two.x, self.point_two.y, fill=fill_color, width=self.width)

class Cell:
  def __init__(self, win=None):
    self.has_left_wall = True
    self.has_right_wall = True
    self.has_top_wall = True
    self.has_bottom_wall = True
    self._win = win
    self._p1 = None
    self._p2 = None
    self._p3 = None
    self._p4 = None
    self._c = None
    self.visited = False

  def draw(self, x1, y1, x2, y2):
    self._p1 = Point(x1, y1)
    self._p2 = Point(x1, y2)
    self._p3 = Point(x2, y2)
    self._p4 = Point(x2, y1)
    self._c = Point(int((x1+x2)/2), int((y1+y2)/2))

    if self._win is None:
      return

    l = Line(self._p1, self._p2)
    c = "black" if self.has_left_wall else "white"
    self._win.draw_line(l, c)

    l = Line(self._p2, self._p3)
    c = "black" if self.has_bottom_wall else "white"
    self._win.draw_line(l, c)

    l = Line(self._p3, self._p4)
    c = "black" if self.has_right_wall else "white"
    self._win.draw_line(l, c)
    
    l = Line(self._p4, self._p1)
    c = "black" if self.has_top_wall else "white"
    self._win.draw_line(l, c)
    

  def draw_move(self, to_cell, undo=False):
    color = "gray" if undo else "red"

    if to_cell is None:
      raise Exception("to_cell is not set")

    l = Line(self._c, to_cell._c)
    self._win.draw_line(l, color)