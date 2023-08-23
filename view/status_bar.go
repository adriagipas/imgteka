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
 *  status_bar.go - Implementa la barra d'estat inferior.
 */

package view

import (
  "fmt"
  "image/color"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PÚBLICA */
/****************/

type StatusBar struct {
  root     *fyne.Container // Contenedor arrel
  text_box *fyne.Container
  model    DataModel
}


func NewStatusBar ( model DataModel ) *StatusBar {

  // Text
  text_box:= container.NewMax ()

  // Barra
  bar:= container.NewBorder ( widget.NewSeparator (), nil, text_box, nil )

  // Crea
  ret:= StatusBar{
    root : bar,
    text_box : text_box,
    model : model,
  }

  // Actualitza
  ret.Update ()
  
  return &ret
  
} // end NewStatusBar


func (self *StatusBar) GetCanvas() fyne.CanvasObject { return self.root }

func (self *StatusBar) Update() {

  stats:= self.model.GetStats ()
  text:= fmt.Sprintf ( "Entrades: %d    Fitxers: %d",
    stats.GetNumEntries (), stats.GetNumFiles () )
  self.text_box.RemoveAll ()
  self.text_box.Add ( canvas.NewText ( text, color.RGBA{25,25,25,255} ) )
  
} // end Update
