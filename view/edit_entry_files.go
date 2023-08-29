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

func showEditFile(

  e        Entry,
  f        File,
  file_id  int64,
  main_win fyne.Window,
  list_win *List,
  list     *widget.List,
  
) {

  // Nom
  name:= widget.NewEntry ()
  name.Text= f.GetName ()
  name.Validator= validation.NewRegexp ( `^.+$`,
    "el nom ha de contindre almenys un caràcter" )

  // Dialeg
  items:= []*widget.FormItem{
    widget.NewFormItem ( "Nom", name ),
  }
  d:= dialog.NewForm ( "Edita nom fitxer", "Aplica", "Cancel·la", items,
    func(b bool){
      if !b { return }
      
      if err:= e.UpdateFileName ( file_id, name.Text ); err != nil {
        dialog.ShowError ( err, main_win )
      } else {
        list.Refresh ()
        list_win.Update ()
      }
      
    }, main_win )
  csize:= main_win.Content ().Size ()
  d.Resize ( fyne.Size{csize.Width*0.6,csize.Height*0.6} )
  d.Show ()
  
} // end showEditFile


type _AddFileProgressBar struct {
  pop  *widget.PopUp
  pb   *widget.ProgressBar
  text *widget.Label
}

func newAddFileProgressBar( win fyne.Window ) *_AddFileProgressBar {

  // Crea PopUP amb caixa buida
  pop_box:= container.NewMax ( )
  pop:= widget.NewModalPopUp ( pop_box, win.Canvas () )

  // Missatge de text
  text:= widget.NewLabel ( "" )

  // Barra de progrés
  pb:= widget.NewProgressBar ()

  // Contingut
  content:= container.NewVBox ( text, container.NewPadded ( pb ) )
  pop_box.Add ( content )

  // Mostra
  csize:= win.Content ().Size ()
  pop.Resize ( fyne.Size{csize.Width*0.4,pop.Size ().Height} )
  pop.Show ()

  // Retorna objecte
  ret:= _AddFileProgressBar{
    pop  : pop,
    pb   : pb,
    text : text,
  }

  return &ret
  
} // end newAddFileProgressBar


func (self *_AddFileProgressBar) Close() {
  self.pop.Hide ()
} // end Close


func (self *_AddFileProgressBar) Set( msg string, f float32 ) {

  self.pb.SetValue ( float64(f) )
  self.text.SetText ( msg )
  self.pop.Refresh ()
  
} // end Set


func showAddFileEntry(
  
  e         Entry,
  model     DataModel,
  main_win  fyne.Window,
  list      *widget.List,
  list_win  *List,
  dv        *DetailsViewer,
  statusbar *StatusBar,
) {

  // Obri dialeg
  d:= dialog.NewFileOpen ( func(r fyne.URIReadCloser,err error){
    if err != nil {
      dialog.ShowError ( err, main_win )
    } else if r != nil {
      
      // Obté nom i path
      uri:= r.URI ()
      
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
          if err:= e.AddFile ( uri.Path (), name.Text,
            type_text2id[type_sel.Selected], func() ProgressBar{
              return newAddFileProgressBar ( main_win )
          }); err != nil {
            dialog.ShowError ( err, main_win )
          } else {
            list.Refresh ()
            list_win.Refresh ()
            dv.Update ()
            statusbar.Update ()
          }
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
  but_edit:= widget.NewButtonWithIcon ( "", theme.DocumentCreateIcon (),
    func(){
      fmt.Println ( "Edita!" )
    })
  but_del:= widget.NewButtonWithIcon ( "", theme.DeleteIcon (),
    func(){
      fmt.Println ( "Esborra!" )
    })
  but_box:= container.NewHBox ( but_edit, but_del )
  
  return container.NewBorder ( nil, nil, nil, but_box, name )
  
} // end createFileEntryItemTemplate


func updateFileEntryItem (
  
  co        fyne.CanvasObject,
  e         Entry,
  model     DataModel,
  list_win  *List,
  dv        *DetailsViewer,
  statusbar *StatusBar,
  id        int,
  list      *widget.List,
  main_win  fyne.Window,
  
) {

  // Prepara
  files:= e.GetFileIDs ()
  f:= model.GetFile ( files[id] )
  label:= co.(*fyne.Container).Objects[0].(*widget.Label)
  but_box:= co.(*fyne.Container).Objects[1].(*fyne.Container)
  
  // Nom
  label.SetText ( f.GetName () )
  
  // Esborra
  but_del:= but_box.Objects[1].(*widget.Button)
  but_del.OnTapped= func() {
    dialog.ShowConfirm ( "Eliminar fitxer",
      "Està segur que vol eliminar aquest fitxer?",
        func(ok bool) {
          if ok {
            if err:= e.RemoveFile ( files[id] ); err != nil {
              dialog.ShowError ( err, main_win )
            } else {
              list.Refresh ()
              list_win.Refresh ()
              dv.Update ()
              statusbar.Update ()
            }
          }
        }, main_win )
  }

  // Edita
  but_edit:= but_box.Objects[0].(*widget.Button)
  but_edit.OnTapped= func() {
    showEditFile ( e, f, files[id], main_win, list_win, list )
  }
  
} // end updateFileEntryItem




/****************/
/* PART PÚBLICA */
/****************/

func NewEditEntryFiles (
  
  e         Entry,
  model     DataModel,
  list_win  *List,
  dv        *DetailsViewer,
  statusbar *StatusBar,
  main_win  fyne.Window,
  
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
    updateFileEntryItem ( w, e, model, list_win, dv, statusbar,
      id, list, main_win )
  }

  // Botonera
  but_new:= widget.NewButtonWithIcon ( "Afegeix Fitxer",
    theme.ContentAddIcon (), func(){
      showAddFileEntry ( e, model, main_win, list, list_win, dv, statusbar )
    })
  but_box:= container.NewHBox ( but_new )
  but_box= container.NewPadded ( but_box )
  
  // Crea contingut
  ret:= container.NewBorder ( but_box, nil, nil, nil, list )
  
  return ret
  
} // end NewEditEntryFiles
