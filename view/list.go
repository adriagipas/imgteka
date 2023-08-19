/*
 * Copyright 2023 Adrià Giménez Pastor.
 *
 * This file is part of adriagipas/imgteka.
 *
 * adriagipas/imgteka is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * adriagipas/imgteka is distributed in the hope that it will be
 * useful, but WITHOUT ANY WARRANTY; without even the implied warranty
 * of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with adriagipas/imgteka.  If not, see <https://www.gnu.org/licenses/>.
 */
/*
 *  list.go - Implementa el llistat de les imatges.
 */

package view

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/widget"
)




/***********/
/* FACTORY */
/***********/

type _Factory struct {
}

func (self *_Factory) childUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
  switch id {
  case "":
    return []widget.TreeNodeID{"a", "b", "c"}
  case "a":
    return []widget.TreeNodeID{"a1", "a2"}
  }
  return []string{}
}


func (self *_Factory) isBranch(id widget.TreeNodeID) bool {
  return id == "" || id == "a"
}


func (self *_Factory) create(branch bool) fyne.CanvasObject {
  if branch {
    return widget.NewLabel("Branch template")
  }
  return widget.NewLabel("Leaf template")
}


func (self *_Factory) update(
  id     widget.TreeNodeID,
  branch bool,
  o      fyne.CanvasObject,
) {
  text := id
  if branch {
    text += " (branch)"
  }
  o.(*widget.Label).SetText(text)
}




/**********************/
/* FUNCIONS PÚBLIQUES */
/**********************/

func GetList() fyne.CanvasObject {

  f:= _Factory{}
  tree := widget.NewTree(
    // childUIDs
		func(id widget.TreeNodeID) []widget.TreeNodeID {
      return f.childUIDs(id)
		},
    // isBranch
		func(id widget.TreeNodeID) bool {
			return f.isBranch(id)
		},
    // create
		func(branch bool) fyne.CanvasObject {
			return f.create(branch)
		},
    // update
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			f.update(id,branch,o)
		})
  
  return tree
}

