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
 *  edit_entry_cover.go - Pestanya per a seleccionar la portada
 */

package view

import (
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PÚBLICA */
/****************/

func NewEditEntryCover (
  
  e         Entry,
  model     DataModel,
  dv        *DetailsViewer,
  main_win  fyne.Window,
  
) fyne.CanvasObject {

  // Caixa imatges
  img_box:= container.NewCenter ( )
  
  // Llista fitxers imatges
  list:= widget.NewList (
    func() int {return -1},
    func() fyne.CanvasObject {return nil},
    func(id widget.ListItemID,w fyne.CanvasObject){},
  )
  
  // --> Length
  list.Length= func() int {
    return len(e.GetImageFileIDs ())+1
  }
  
  // --> CreateItem
  list.CreateItem= func() fyne.CanvasObject {
    return widget.NewLabel ( "Template Cover Name" )
  }
  
  // --> UpdateItem
  list.UpdateItem= func( id widget.ListItemID, w fyne.CanvasObject ) {

    // Text entrada
    var text string
    if id == 0 {
      text= "[Sense portada]"
    } else {
      f:= model.GetFile ( e.GetImageFileIDs ()[id-1] )
      text= f.GetName ()
    }

    // Modifica
    w.(*widget.Label).SetText ( text )
    
  }

  // --> OnSelected
  list.OnSelected= func( id widget.ListItemID ) {

    // Obté identificador
    var fid int64
    if id == 0 {
      fid= -1
    } else {
      fid= e.GetImageFileIDs ()[id-1]
    }

    // Actualitza imatge
    img_box.RemoveAll ()
    if fid != -1 {
      img:= model.GetFile ( fid ).GetImage ( 250 )
      if img != nil {
        img_w:= canvas.NewImageFromImage ( img )
        img_w.FillMode= canvas.ImageFillContain
        img_w.SetMinSize ( fyne.Size{250,250} )
        img_box.Add ( img_w )
      }
    }
    
    // Actualitza portada
    if err:= e.SetCoverFileID ( fid ); err != nil {
      dialog.ShowError ( err, main_win )
    } else {
      dv.Update ()
    }
    img_box.Refresh ()
    
  }

  // Split
  split:= container.NewHSplit ( list, img_box )
  
  return split
  
} // end NewEditEntryCover
