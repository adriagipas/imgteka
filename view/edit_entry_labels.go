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
 *  edit_entry_labels.go - Pestanya per a editar les etiquetes.
 */

package view

import (
  "fmt"
  "image/color"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PRIVADA */
/****************/

func showAddLabelEntry(
  
  e        Entry,
  model    DataModel,
  main_win fyne.Window,
  list     *widget.List,
  dv       *DetailsViewer,
  
) {
  
  // Selector d'etiquetes
  ids:= e.GetUnusedLabelIDs ()
  options:= make([]string,len(ids))
  lbl_text2id:= make(map[string]int)
  for i,id:= range ids {
    label:= model.GetLabel ( id )
    text:= label.GetName ()
    options[i]= text
    lbl_text2id[text]= id
  }
  lbl_sel:= widget.NewSelect ( options, func(string){} )
  lbl_sel.SetSelectedIndex ( 0 )

  // Dialeg
  items:= []*widget.FormItem{
    widget.NewFormItem ( "Etiqueta", lbl_sel ),
  }
  d:= dialog.NewForm ( "Afegeix etiqueta", "Afegeix", "Cancel·la", items,
    func(b bool){
      if !b { return }
      if err:= e.AddLabel ( lbl_text2id[lbl_sel.Selected] ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
        dv.Update ()
      }
    }, main_win )
  win_size:= main_win.Content ().Size ()
  d.Resize ( fyne.Size{win_size.Width*0.4,win_size.Height*0.4} )
  d.Show ()
  
} // end showAddLabelEntry


func createLabelEntryItemTemplate () fyne.CanvasObject {

  // Text
  color_rect:= canvas.NewRectangle ( color.Black )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  name:= widget.NewLabel ( "Template Label Name" )
  box:= container.NewHBox ( color_rect, name )

  // Botons
  but_del:= widget.NewButtonWithIcon ( "", theme.DeleteIcon (),
    func(){
      fmt.Println ( "Esborra!" )
    })
  but_box:= container.NewHBox ( but_del )

  return container.NewBorder ( nil, nil, nil, but_box, box )
  
} // end createLabelEntryItemTemplate


func updateLabelEntryItem (
  
  co       fyne.CanvasObject,
  e        Entry,
  model    DataModel,
  dv       *DetailsViewer,
  id       int,
  list     *widget.List,
  main_win fyne.Window,
  
) {

  // Prepara
  labels:= e.GetLabelIDs ()
  label:= model.GetLabel ( labels[id] )
  box:= co.(*fyne.Container).Objects[0].(*fyne.Container)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Color
  color_rect:= canvas.NewRectangle ( label.GetColor () )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  box.Objects[0]= color_rect
  
  // Nom
  text:= fmt.Sprintf ( "%s", label.GetName () )
  box.Objects[1].(*widget.Label).SetText ( text )
  
  // Esborra
  but_del:= but_box.Objects[0].(*widget.Button)
  but_del.OnTapped= func() {
    dialog.ShowConfirm ( "Eliminar etiqueta de l'entrada",
      "Està segur que vol eliminar aquesta etiqueta per a aquesta entrada?",
        func(ok bool) {
          if ok {
            if err:= e.RemoveLabel ( labels[id] ); err != nil {
              dialog.ShowError ( err, main_win )
            } else {
              list.Refresh ()
              dv.Update ()
            }
          }
        }, main_win )
    }
  
} // end updateLabelEntryItem  




/***************/
/* PART PÚBLIC */
/***************/

func NewEditEntryLabels (
  
  e        Entry,
  model    DataModel,
  dv       *DetailsViewer,
  main_win fyne.Window,
  
) fyne.CanvasObject {

  // Llista etiquetes
  list:= widget.NewList (
    func() int {return -1},
    func() fyne.CanvasObject {return nil},
    func(id widget.ListItemID,w fyne.CanvasObject){},
  )
  list.Length= func() int {
    return len(e.GetLabelIDs ())
  }
  list.CreateItem= func() fyne.CanvasObject {
    return createLabelEntryItemTemplate ()
  }
  list.UpdateItem= func( id widget.ListItemID, w fyne.CanvasObject ) {
    updateLabelEntryItem ( w, e, model, dv, id, list, main_win )
  }

  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Afegeix Etiqueta",
    theme.ContentAddIcon (), func(){
      if len(e.GetUnusedLabelIDs ()) > 0 {
        showAddLabelEntry ( e, model, main_win, list, dv )
      } else {
        dialog.ShowInformation ( "Afegeix Etiqueta",
          "No es poden afegir més etiquetes", main_win )
      }
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewEditEntryLabels
