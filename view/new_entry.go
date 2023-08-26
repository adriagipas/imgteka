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
 *  new_entry.go - Dialeg per a crear entrades noves.
 */

package view

import (
  "fmt"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/data/validation"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"
)


func ShowNewEntryDialog(
  
  model     DataModel,
  list      *List,
  statusbar *StatusBar,
  main_win  fyne.Window,
  
) {

  // Plataforma
  pids:= model.GetPlatformIDs ()
  options:= make([]string,len(pids))
  plat_text2id:= make(map[string]int)
  for i,id:= range pids {
    plat:= model.GetPlatform ( id )
    text:= fmt.Sprintf ( "%s - %s", plat.GetShortName (), plat.GetName () )
    options[i]= text
    plat_text2id[text]= id
  }
  plat_sel:= widget.NewSelect ( options, func(string){} )
  plat_sel.SetSelectedIndex ( 0 )
  
  // Nom
  name:= widget.NewEntry ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )
  
  // Dialeg
  items:= []*widget.FormItem{
    widget.NewFormItem ( "Plataforma", plat_sel ),
    widget.NewFormItem ( "Nom", name ),
  }
  d:= dialog.NewForm ( "Nova entrada", "Afegeix", "Cancel·la", items,
    func(b bool){
      if !b { return }
      fmt.Println ( "HOOOOLA", name.Text, plat_text2id[plat_sel.Selected] )
    }, main_win )
  win_size:= main_win.Content ().Size ()
  d.Resize ( fyne.Size{win_size.Width*0.4,win_size.Height*0.4} )
  d.Show ()
  
} // end ShowNewEntryDialog

