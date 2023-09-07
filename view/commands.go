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
 *  commands.go - Pestanya per a gestionar els comandaments.
 */

package view

import (
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)


/****************/
/* PART PRIVADA */
/****************/

/****************/
/* PART PÚBLICA */
/****************/

func NewCommandsManager(
  
  model    DataModel,
  main_win fyne.Window,
  
) fyne.CanvasObject {

  tids:= model.GetFileTypeIDs ()
  entries:= make([]*widget.Entry,len(tids))
  form:= widget.NewForm ()
  for i,tid:= range model.GetFileTypeIDs () {
    text:= model.GetFileTypeName ( tid )
    entry:= widget.NewEntry ()
    entry.Text= model.GetFileTypeCommand ( tids[i] )
    form.Append ( text, entry )
    entries[i]= entry
  }
  form.SubmitText= "Aplica"
  form.OnSubmit= func() {
    for i:= 0; i < len(tids); i++ {
      model.SetFileTypeCommand ( tids[i], entries[i].Text )
    }
  }
  form.Refresh ()
  
  return container.NewVScroll ( form )
  
} // end NewCommandsManager
