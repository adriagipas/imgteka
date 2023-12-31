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
 *  platforms.go - Pestanya per a gestionar les plataformes.
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

func showEditPlatform (
  
  plat     Platform,
  main_win fyne.Window,
  list_win *List,
  dv       *DetailsViewer,
  list     *widget.List,
  
) {

  // Nom llarg
  name:= widget.NewEntry ()
  name.Text= plat.GetName ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )

  // Color
  var mcolor color.Color
  mcolor= plat.GetColor ()
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
  dialog.ShowForm ( "Edita plataforma", "Aplica", "Cancel·la", items,
    func(b bool){
      if !b { return }
      if err:= plat.Update ( name.Text, mcolor ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
        list_win.Update ()
        dv.Update ()
      }
    }, main_win )
  
} // end showEditPlatform


func showNewPlatform (
  model    DataModel,
  main_win fyne.Window,
  list     *widget.List,
) {

  // Nom curt
  shortname:= widget.NewEntry ()
  shortname.Validator= validation.NewRegexp ( `^[A-Za-z0-9]{1,3}$`,
    "el nom curt sols pot contindre majúscules, un mínim d'un caràcter"+
      " i un màxim de tres" )

  // Nom llarg
  name:= widget.NewEntry ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )

  // Color
  var mcolor color.Color
  mcolor= color.Black
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
    widget.NewFormItem ( "Nom curt", shortname ) ,
    widget.NewFormItem ( "Nom", name ),
    widget.NewFormItem ( "Color", color_box ),
  }
  dialog.ShowForm ( "Plataforma nova", "Afegeix", "Cancel·la", items,
    func(b bool){
      if !b { return }
      if err:= model.AddPlatform (
        shortname.Text, name.Text, mcolor ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
      }
    }, main_win )
  
} // end showNewPlatform


func createPlatformItemTemplate () fyne.CanvasObject {

  // Text
  label:= newPlatformLabel ( "ABC", color.Black )
  name:= widget.NewLabel ( "Template Platform Name" )
  box:= container.NewHBox ( label, name )

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
  
} // end createPlatformItemTemplate


func updatePlatformItem (
  
  co       fyne.CanvasObject,
  model    DataModel,
  list_win *List,
  dv       *DetailsViewer,
  id       int,
  list     *widget.List,
  main_win fyne.Window,
  
) {

  // Prepara
  plats:= model.GetPlatformIDs ()
  plat:= model.GetPlatform ( plats[id] )
  box:= co.(*fyne.Container).Objects[0].(*fyne.Container)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Etiqueta
  label:= newPlatformLabel ( plat.GetShortName (), plat.GetColor () )
  box.Objects[0]= label

  // Nom
  text:= fmt.Sprintf ( "%s (%d)", plat.GetName (), plat.GetNumEntries () )
  box.Objects[1].(*widget.Label).SetText ( text )

  // Esborra
  but_del:= but_box.Objects[1].(*widget.Button)
  if plat.GetNumEntries () > 0 {
    but_del.Disable ()
    but_del.OnTapped= func() {}
  } else {
    but_del.Enable ()
    but_del.OnTapped= func() {
      dialog.ShowConfirm ( "Esborra plataforma",
        "Està segur que vol esborrar la plataforma?",
        func(ok bool) {
          if ok {
            if err:= model.RemovePlatform ( plats[id] ); err != nil {
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
    showEditPlatform ( plat, main_win, list_win, dv, list )
  }
  
} // end updatePlatformItem




/****************/
/* PART PÚBLICA */
/****************/

func NewPlatformsManager (
  
  model    DataModel,
  list_win *List,
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
    return len(model.GetPlatformIDs ())
  }
  list.CreateItem= func() fyne.CanvasObject {
    return createPlatformItemTemplate ()
  }
  list.UpdateItem= func( id widget.ListItemID, w fyne.CanvasObject ) {
    updatePlatformItem ( w, model, list_win, dv, id, list, main_win )
  }
  
  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Nova Plataforma",
    theme.ContentAddIcon (), func(){
      showNewPlatform ( model, main_win, list )
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewPlatformManager
