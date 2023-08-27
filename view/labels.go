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
 *  labels.go - Pestanya per a gestionar les etiquetes.
 */

package view

import (
  "fmt"
  "image/color"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/data/validation"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PRIVADA */
/****************/

func showEditLabel (
  
  label    Label,
  main_win fyne.Window,
  dv       *DetailsViewer,
  list     *widget.List,
  
) {

  // Nom llarg
  name:= widget.NewEntry ()
  name.Text= label.GetName ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )

  // Color
  var mcolor color.Color
  mcolor= label.GetColor ()
  color_rect:= canvas.NewRectangle ( mcolor )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  color_but:= widget.NewButton ( "Selecciona", func(){} )
  color_box:= container.NewHBox ( color_rect,
    widget.NewSeparator (), color_but )
  color_but.OnTapped= func(){
    picker:= dialog.NewColorPicker ( "Selecciona un color",
      "Selecciona un color", func(c color.Color){
        mcolor= c
        color_rect= canvas.NewRectangle ( mcolor )
        color_rect.SetMinSize ( fyne.Size{30,1} )
        color_box.Objects[0]= color_rect
      }, main_win )
    picker.Advanced= true
    picker.Show ()
  }
  
  // Dialeg
  items:= []*widget.FormItem{
    widget.NewFormItem ( "Nom", name ),
    widget.NewFormItem ( "Color", color_box ),
  }
  dialog.ShowForm ( "Edita etiqueta", "Aplica", "Cancel·la", items,
    func(b bool){
      if !b { return }
      if err:= label.Update ( name.Text, mcolor ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
        dv.Update ()
      }
    }, main_win )
  
} // end showEditLabel


func showNewLabel (
  model    DataModel,
  main_win fyne.Window,
  list     *widget.List,
) {

  // Nom llarg
  name:= widget.NewEntry ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )

  // Color
  var mcolor color.Color
  mcolor= color.White
  color_rect:= canvas.NewRectangle ( mcolor )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  color_but:= widget.NewButton ( "Selecciona", func(){} )
  color_box:= container.NewHBox ( color_rect,
    widget.NewSeparator (), color_but )
  color_but.OnTapped= func(){
    picker:= dialog.NewColorPicker ( "Selecciona un color",
      "Selecciona un color", func(c color.Color){
        mcolor= c
        color_rect= canvas.NewRectangle ( mcolor )
        color_rect.SetMinSize ( fyne.Size{30,1} )
        color_box.Objects[0]= color_rect
      }, main_win )
    picker.Advanced= true
    picker.Show ()
  }
  
  // Dialeg
  items:= []*widget.FormItem{
    widget.NewFormItem ( "Nom", name ),
    widget.NewFormItem ( "Color", color_box ),
  }
  dialog.ShowForm ( "Etiqueta nova", "Afegeix", "Cancel·la", items,
    func(b bool){
      if !b { return }
      if err:= model.AddLabel ( name.Text, mcolor ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
      }
    }, main_win )
  
} // end showNewLabel


func createLabelItemTemplate () fyne.CanvasObject {

  // Text
  color_rect:= canvas.NewRectangle ( color.Black )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  name:= widget.NewLabel ( "Template Label Name" )
  box:= container.NewHBox ( color_rect, name )

  // Botons
  but_edit:= widget.NewButtonWithIcon ( "", theme.DocumentCreateIcon (),
    func(){
      fmt.Println ( "Edita!" )
    })
  but_del:= widget.NewButtonWithIcon ( "", theme.DeleteIcon (),
    func(){
      fmt.Println ( "Esborra!" )
    })
  but_box:= container.NewHBox ( but_edit, but_del )

  return container.NewBorder ( nil, nil, nil, but_box, box )
  
} // end createLabelItemTemplate


func updateLabelItem (
  
  co       fyne.CanvasObject,
  model    DataModel,
  dv       *DetailsViewer,
  id       int,
  list     *widget.List,
  main_win fyne.Window,
  
) {

  // Prepara
  labels:= model.GetLabelIDs ()
  label:= model.GetLabel ( labels[id] )
  box:= co.(*fyne.Container).Objects[0].(*fyne.Container)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Color
  color_rect:= canvas.NewRectangle ( label.GetColor () )
  color_rect.SetMinSize ( fyne.Size{30,1} )
  box.Objects[0]= color_rect
  
  // Nom
  text:= fmt.Sprintf ( "%s (%d)", label.GetName (), label.GetNumEntries () )
  box.Objects[1].(*widget.Label).SetText ( text )
  
  // Esborra
  but_del:= but_box.Objects[1].(*widget.Button)
  if label.GetNumEntries () > 0 {
    but_del.Disable ()
    but_del.OnTapped= func() {}
  } else {
    but_del.Enable ()
    but_del.OnTapped= func() {
      dialog.ShowConfirm ( "Esborra etiqueta",
        "Està segur que vol esborrar l'etiqueta?",
        func(ok bool) {
          if ok {
            if err:= model.RemoveLabel ( labels[id] ); err != nil {
              dialog.ShowError ( err, main_win )
            } else {
              list.Refresh ()
            }
          }
        }, main_win )
    }
  }
  
  // Edita
  but_edit:= but_box.Objects[0].(*widget.Button)
  but_edit.OnTapped= func() {
    showEditLabel ( label, main_win, dv, list )
  }
  
} // end updateLabelItem




/****************/
/* PART PÚBLICA */
/****************/

func NewLabelsManager (
  
  model    DataModel,
  dv       *DetailsViewer,
  main_win fyne.Window,
  
) fyne.CanvasObject {
  
  // Llista plataformes
  list:= widget.NewList (
    func() int {return -1},
    func() fyne.CanvasObject {return nil},
    func(id widget.ListItemID,w fyne.CanvasObject){},
  )
  list.Length= func() int {
    return len(model.GetLabelIDs ())
  }
  list.CreateItem= func() fyne.CanvasObject {
    return createLabelItemTemplate ()
  }
  list.UpdateItem= func( id widget.ListItemID, w fyne.CanvasObject ) {
    updateLabelItem ( w, model, dv, id, list, main_win )
  }
  
  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Nova Etiqueta",
    theme.ContentAddIcon (), func(){
      showNewLabel ( model, main_win, list )
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewLabelManager
