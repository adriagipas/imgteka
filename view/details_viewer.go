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
 *  details_viewer.go - Implementa el visor de detalls.
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

type DetailsViewer struct {
  root  *fyne.Container // Contenedor bàsic
  model DataModel
}


func NewDetailsViewer ( model DataModel ) *DetailsViewer {

  ret:= DetailsViewer{
    root : container.NewVBox (),
    model : model,
  }
  ret.root.Add ( widget.NewLabel ( "Hola" ) )
  
  return &ret
  
} // end NewDetailsViewer


func (self *DetailsViewer) GetCanvas() fyne.CanvasObject { return self.root }


func (self *DetailsViewer) ViewEntry ( e Entry ) {
  
  // Neteja
  self.root.RemoveAll ()

  // Títol
  title:= canvas.NewText ( e.GetName (), color.RGBA{60,60,60,255} )
  title.TextSize= 25
  title.Alignment= fyne.TextAlignCenter
  title.TextStyle= fyne.TextStyle{Bold:true}
  self.root.Add ( title )
  
  // Portada
  cover:= e.GetCover ()
  if cover != nil {
    img:= canvas.NewImageFromImage ( cover )
    img.FillMode= canvas.ImageFillContain
    img.SetMinSize ( fyne.Size{200,200} )
    self.root.Add ( img )
  }
  
} // end ViewEntry


func (self *DetailsViewer) ViewFile ( f File ) {
  fmt.Printf ( "VIEW_FILE: %v\n", f )

} // end ViewFile
