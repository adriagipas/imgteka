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
 *  details_viewer.go - Implementa el visor de detalls.
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

const _DETAILS_VIEWER_EMPTY = 0
const _DETAILS_VIEWER_FILE  = 1
const _DETAILS_VIEWER_ENTRY = 2

func (self *DetailsViewer) newLabel ( id int ) fyne.CanvasObject {

  label:= self.model.GetLabel ( id )
  name,mcolor:= label.GetName (),label.GetColor ()

  gray:= color.RGBA{50,50,50,255}
  rect:= canvas.NewRectangle ( mcolor )
  rect.StrokeColor= gray
  rect.StrokeWidth= 0.5
  text:= canvas.NewText ( name, gray )
  text.TextSize= 11.0
  text.Alignment= fyne.TextAlignCenter
  text.TextStyle= fyne.TextStyle{Bold:true}
  ret:= container.NewMax ( rect, container.NewPadded ( text ) )
  
  return ret
  
} // end newLabel




/****************/
/* PART PÚBLICA */
/****************/

type DetailsViewer struct {
  
  root      *fyne.Container // Contenedor intern que es modifica
  model     DataModel
  statusbar *StatusBar
  win       fyne.Window

  // Estat
  state      int
  current_fe int64
  list       *List // El que s'utilitza si estem en mode Entry
  
}


func NewDetailsViewer (

  model     DataModel,
  statusbar *StatusBar,
  main_win  fyne.Window,

) *DetailsViewer {

  ret:= DetailsViewer{
    root      : container.NewVBox (),
    model     : model,
    statusbar : statusbar,
    win       : main_win,

    state      : _DETAILS_VIEWER_EMPTY,
    current_fe : -1,
    
  }
  
  return &ret
  
} // end NewDetailsViewer


func (self *DetailsViewer) GetCanvas() fyne.CanvasObject { return self.root }


func (self *DetailsViewer) Clean() {

  self.root.RemoveAll ()
  self.state= _DETAILS_VIEWER_EMPTY
  
} // end Clean


func (self *DetailsViewer) ViewEntry ( e_id int64, list *List ) {
  
  // Neteja
  self.Clean ()
  self.state= _DETAILS_VIEWER_ENTRY
  self.current_fe= e_id
  self.list= list

  // Obté entry
  e:= self.model.GetEntry ( e_id )
  
  // Crea card
  // --> Contingut
  label_ids:= e.GetLabelIDs ()
  labels:= make([]fyne.CanvasObject,len(label_ids))
  maxw,maxh:= float32(1),float32(1)
  for i:= 0; i < len(label_ids); i++ {
    l:= self.newLabel ( label_ids[i] )
    labels[i]= l
    size:= l.MinSize ()
    if size.Width > maxw { maxw= size.Width }
    if size.Height > maxh { maxh= size.Height }
  }
  content:= container.NewGridWrap ( fyne.Size{maxw,maxh} )
  content.Objects= labels
  text_tmp:= fmt.Sprintf (
    `**Nº Fitxers:** %d

**Etiquetes:**`,
    len(e.GetFileIDs ()),
  )
  text:= widget.NewRichTextFromMarkdown ( text_tmp )
  content= container.NewVBox ( text, content )
  // --> Card
  card:= widget.NewCard (
    e.GetName (),
    self.model.GetPlatform ( e.GetPlatformID () ).GetName (),
    content,
  )
  // --> Portada
  cover:= e.GetCover ()
  if cover != nil {
    img:= canvas.NewImageFromImage ( cover )
    img.FillMode= canvas.ImageFillContain
    img.SetMinSize ( fyne.Size{1,1} )
    card.SetImage ( img )
  }
  
  // Crea toolbar
  toolbar:= widget.NewToolbar (
    widget.NewToolbarSpacer (),
    widget.NewToolbarAction ( theme.DocumentCreateIcon (), func() {
      RunEditEntryWin ( e, self.model, list, self, self.statusbar, self.win )
    }),
    widget.NewToolbarAction ( theme.DeleteIcon (), func() {
      dialog.ShowConfirm ( "Esborra entrada",
        "Està segur que vol esborrar l'entrada?",
        func(ok bool) {
          if ok {
            if err:= self.model.RemoveEntry ( e_id ); err != nil {
              dialog.ShowError ( err, self.win )
            } else {
              self.Clean ()
              self.list.Refresh ()
              self.statusbar.Update ()
            }
          }
        }, self.win )
    }),
  )
  
  // Afegeix
  tmp:= container.NewVBox ( container.NewHScroll ( card ), toolbar )
  self.root.Add ( tmp )
  
} // end ViewEntry


func (self *DetailsViewer) ViewFile ( f_id int64 ) {
  
  // Neteja
  self.Clean ()
  self.state= _DETAILS_VIEWER_FILE
  self.current_fe= f_id

  // Obté fitxer
  f:= self.model.GetFile ( f_id )
  
  // Contingut
  text_tmp:= ""
  md:= f.GetMetadata ()
  for i:= 0; i < len(md); i++ {
    text_tmp+= fmt.Sprintf ( `
- **%s:** %s
`,
      md[i].GetKey (), md[i].GetValue () )
  }
  text:= widget.NewRichTextFromMarkdown ( text_tmp )
  
  // Crea card
  card:= widget.NewCard (
    f.GetName (),
    self.model.GetFileTypeName ( f.GetTypeID () ),
    text,
  )

  // Crea toolbar
  // NOTA!! En el futur el RUN el podem ficar sols si el tipus de
  // fitxer es pot executar.
  toolbar:= widget.NewToolbar (
    widget.NewToolbarSpacer (),
    widget.NewToolbarAction ( theme.MediaPlayIcon (), func() {
      fmt.Println ( "PLAY BUTTON!!!!" )
    }),
  )
  
  // Afegeix
  tmp:= container.NewVBox ( container.NewHScroll ( card ), toolbar )
  self.root.Add ( tmp )
  
} // end ViewFile


func (self *DetailsViewer) Update() {

  switch self.state {
    
  case _DETAILS_VIEWER_FILE:
    self.ViewFile ( self.current_fe )

  case _DETAILS_VIEWER_ENTRY:
    self.ViewEntry ( self.current_fe, self.list )
    
  }
  
} // end Update
