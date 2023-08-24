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
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PRIVADA */
/****************/

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


func updatePlatformItem ( co fyne.CanvasObject, model DataModel, id int ) {

  // Prepara
  plats:= model.GetPlatformIDs ()
  plat:= model.GetPlatform ( plats[id] )
  box:= co.(*fyne.Container).Objects[0].(*fyne.Container)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Etiqueta
  label:= newPlatformLabel ( plat.GetShortName (), plat.GetColor () )
  box.Objects[0]= label

  // Nom
  text:= fmt.Sprintf ( "%s (%d)", plat.GetName (), plat.GetNumFiles () )
  box.Objects[1].(*widget.Label).SetText ( text )

  // Esborra
  but_del:= but_box.Objects[1].(*widget.Button)
  if plat.GetNumFiles () > 0 {
    but_del.Disable ()
    but_del.OnTapped= func() {}
  } else {
    but_del.OnTapped= func() {
      fmt.Printf ( "Esborra %v (%d)\n", plat, plats[id] )
    }
  }

  // Edita
  but_edit:= but_box.Objects[0].(*widget.Button)
  but_edit.OnTapped= func() {
    fmt.Printf ( "Edita %v\n", plat )
  }
  
} // end updatePlatformItem




/****************/
/* PART PÚBLICA */
/****************/

func NewPlatformsManager ( model DataModel ) fyne.CanvasObject {

  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Nova Plataforma",
    theme.ContentAddIcon (), func(){
      fmt.Println ( "Nova Plataforma" )
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Llista plataformes
  list:= widget.NewList (
    
    // length
    func() int {
      return len(model.GetPlatformIDs ())
    },

    // createItem
    func() fyne.CanvasObject {
      return createPlatformItemTemplate ()
    },

    // updateItem
    func( id widget.ListItemID, w fyne.CanvasObject ) {
      updatePlatformItem ( w, model, id )
    },
  )

  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewPlatformManager
