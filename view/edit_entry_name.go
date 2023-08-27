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
 *  edit_entry_name.go - Pestanya per a editar el nom de l'entrada.
 */

package view

import (
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/data/validation"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PÚBLICA */
/****************/

func NewEditEntryName (
  
  e        Entry,
  list_win *List,
  dv       *DetailsViewer,
  main_win fyne.Window,
  
) fyne.CanvasObject {

  // Nom
  name:= widget.NewEntry ()
  name.Text= e.GetName ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )
  name_box:= container.NewBorder ( nil, nil, widget.NewLabel ( "Nom:" ), nil,
    name )
  name_box= container.NewPadded ( name_box )
  
  // Botonera
  but_ok:= widget.NewButton ( "Aplica", func() {
    if err:= e.UpdateName ( name.Text ); err != nil {
      dialog.ShowError ( err, main_win )
    } else {
      list_win.Refresh ()
      dv.Update ()
    }
  })
  but_box:= container.NewBorder ( nil, nil, nil, but_ok )
  but_box= container.NewPadded ( but_box )

  // Crea contingut
  ret:= container.NewBorder ( name_box, but_box, nil, nil )

  return ret
  
} // end NewEditEntryName
