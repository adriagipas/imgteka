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
 *  toolbar.go - Implementa la barra de cerca i menú.
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

type Toolbar struct {
  root  *fyne.Container // Contenedor arrel
}


func NewToolbar (
  
  model      DataModel,
  main_win   fyne.Window,
  list       *List,
  dv         *DetailsViewer,
  status_bar *StatusBar,
  
) *Toolbar {

  // Crea
  ret:= Toolbar{
    root : container.NewVBox (),
  }

  // Crea barra cerca
  search_icon:= widget.NewIcon ( theme.SearchIcon () )
  search_entry:= widget.NewEntry ()
  search_entry.PlaceHolder= "Cerca...   p.e.: consulta1 + p:MD + l:Lluita"
  search_entry.OnSubmitted= func(text string) {
    model.FilterEntries ( text )
    list.Update ()
    status_bar.Update ()
    main_win.Canvas ().Focus ( list )
  }
  search_bar:= container.NewBorder ( nil, nil, search_icon, nil, search_entry )

  // Boto afegir
  add_but:= widget.NewButtonWithIcon ( "", theme.FolderNewIcon (),
    func(){
      ShowNewEntryDialog ( model, list, status_bar, main_win )
    })

  // Botó configuració
  conf_but:= widget.NewButtonWithIcon ( "", theme.SettingsIcon (),
    func(){
      RunConfigWin ( model, list, dv, main_win )
    })
  
  // Afegeix
  box:= container.NewBorder ( nil, nil, add_but, conf_but, search_bar )
  ret.root.Add ( box )
  ret.root.Add ( widget.NewSeparator () )
  
  return &ret
  
} // end NewToolBar


func (self *Toolbar) GetCanvas() fyne.CanvasObject { return self.root }

