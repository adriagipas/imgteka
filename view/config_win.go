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
 *  config_win.go - Finestra per a configurar imgteka.
 */

package view

import (

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PÚBLICA */
/****************/

func RunConfigWin ( model DataModel, main_win fyne.Window ) {

  // Crea PopUP amb una caixa buida
  pop_box:= container.NewMax ()
  pop:= widget.NewModalPopUp ( pop_box, main_win.Canvas () )
  
  // Contingut
  // --> Pestanyes
  plats_tab:= container.NewTabItem (
    "Plataformes",
    container.NewPadded ( NewPlatformsManager ( model, main_win ) ),
  )
  tabs:= container.NewAppTabs ( plats_tab )
  
  // --> Botonera
  but_close:= widget.NewButtonWithIcon ( "Tanca", theme.CancelIcon (), func(){
    pop.Hide ()
  })
  but_box:= container.NewBorder ( widget.NewSeparator (), nil, nil, but_close )
  
  // Mostra
  content:= container.NewBorder ( nil, but_box, nil, nil, tabs )
  pop_box.Add ( content )
  csize:= main_win.Content ().Size ()
  pop.Resize ( fyne.Size{csize.Width*0.7,csize.Height*0.7} )
  pop.Show ()
  
} // end RunConfigWin
