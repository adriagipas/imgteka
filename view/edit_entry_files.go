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
 *  edit_entry_files.go - Pestanya per a editar els fitxers.
 */

package view

import (
  "fmt"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/data/validation"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PRIVADA */
/****************/

func showAddFileEntry(
  
  e        Entry,
  model    DataModel,
  main_win fyne.Window,
  list     *widget.List,
  list_win *List,
  dv       *DetailsViewer,
  
) {

  // Obri dialeg
  d:= dialog.NewFileOpen ( func(r fyne.URIReadCloser,err error){
    if err != nil {
      dialog.ShowError ( err, main_win )
    } else if r != nil {
      
      // Obté nom i path
      uri:= r.URI ()
      //path:= uri.Path ()

      // Nom
      name:= widget.NewEntry ()
      name.Text= uri.Name ()
      name.Validator= validation.NewRegexp ( `^.+$`,
        "el nom ha de contindre almenys un caràcter" )

      // Selector de tipus
      tids:= model.GetFileTypeIDs ()
      options:= make([]string,len(tids))
      type_text2id:= make(map[string]int)
      for i,id:= range tids {
        text:= model.GetFileTypeName ( id )
        options[i]= text
        type_text2id[text]= id
      }
      type_sel:= widget.NewSelect ( options, func(string){} )
      type_sel.SetSelectedIndex ( 0 )

      // Dialeg
      items:= []*widget.FormItem{
        widget.NewFormItem ( "Nom", name ),
        widget.NewFormItem ( "Tipus", type_sel ),
      }
      d2:= dialog.NewForm ( "Afegeix fitxer", "Afegeix", "Cancel·la", items,
        func(b bool){
          if !b { return }
          fmt.Println ( "AFEGIR !!!", uri.Path (), name.Text, type_text2id[type_sel.Selected] )
          /*
          if err:= e.AddLabel ( lbl_text2id[lbl_sel.Selected] ); err != nil {
            dialog.ShowError ( err, main_win )
          } else {
            list.Refresh ()
            list_win.Refresh ()
            dv.Update ()
          }
          */
        }, main_win )
      win_size:= main_win.Content ().Size ()
      d2.Resize ( fyne.Size{win_size.Width*0.4,win_size.Height*0.4} )
      d2.Show ()
      
    }
  }, main_win )
  csize:= main_win.Content ().Size ()
  d.Resize ( fyne.Size{csize.Width*0.8,csize.Height*0.8} )
  d.Show ()
  
} // end showAddFileEntry


func createFileEntryItemTemplate () fyne.CanvasObject {

  // Text
  name:= widget.NewLabel ( "Template Label Name" )
  
  // Botons
  but_del:= widget.NewButtonWithIcon ( "", theme.DeleteIcon (),
    func(){
      fmt.Println ( "Esborra!" )
    })
  but_box:= container.NewHBox ( but_del )
  
  return container.NewBorder ( nil, nil, nil, but_box, name )
  
} // end createFileEntryItemTemplate


func updateFileEntryItem (
  
  co       fyne.CanvasObject,
  e        Entry,
  model    DataModel,
  list_win *List,
  dv       *DetailsViewer,
  id       int,
  list     *widget.List,
  main_win fyne.Window,
  
) {

  // Prepara
  files:= e.GetFileIDs ()
  f:= model.GetFile ( files[id] )
  label:= co.(*fyne.Container).Objects[0].(*widget.Label)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Nom
  label.SetText ( f.GetName () )
  
  // Esborra
  but_del:= but_box.Objects[0].(*widget.Button)
  but_del.OnTapped= func() {
    dialog.ShowConfirm ( "Eliminar fitxer",
      "Està segur que vol eliminar aquest fitxer?",
        func(ok bool) {
          if ok {
            fmt.Println ( "ELIMINAT", f )
            /*
            if err:= e.RemoveLabel ( labels[id] ); err != nil {
              dialog.ShowError ( err, main_win )
            } else {
              list.Refresh ()
              list_win.Refresh ()
              dv.Update ()
            }
            */
          }
        }, main_win )
    }
  
} // end updateFileEntryItem




/****************/
/* PART PÚBLICA */
/****************/

func NewEditEntryFiles (
  
  e        Entry,
  model    DataModel,
  list_win *List,
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
    return len(e.GetFileIDs ())
  }
  list.CreateItem= func() fyne.CanvasObject {
    return createFileEntryItemTemplate ()
  }
  list.UpdateItem= func( id widget.ListItemID, w fyne.CanvasObject ) {
    updateFileEntryItem ( w, e, model, list_win, dv, id, list, main_win )
  }

  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Afegeix Fitxer",
    theme.ContentAddIcon (), func(){
      showAddFileEntry ( e, model, main_win, list, list_win, dv )
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewEditEntryFiles
