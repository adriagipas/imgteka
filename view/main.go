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
 *  main.go - Implementa el punt d'entrada a l'aplicació gràfica.
 */

package view

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2/container"
)


func Run() error {

  // Crea
  a:= app.NewWithID ( "org.github.adriagipas.imgteka" )
  win:= a.NewWindow ( "imgteka" )

  // Construeix elements
  model:= newFakeDataModel ()
  dv:= NewDetailsViewer ( model )
  list:= NewList ( model, dv )
  split:= container.NewHSplit ( list, dv.GetCanvas () )

  // Barra cerca i menú
  toolbar:= NewToolbar ( model )
  mbox:= container.NewBorder ( toolbar.GetCanvas (), nil, nil, nil, split )
  
  // Executa
  win.SetContent ( mbox )
  win.SetMaster ()
  win.Resize ( fyne.Size{800,600} )
  win.ShowAndRun ()
  
  return nil
  
}
