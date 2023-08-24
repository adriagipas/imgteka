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
  "fyne.io/fyne/v2/theme"
  "fyne.io/fyne/v2/widget"
)




/****************/
/* PART PRIVADA */
/****************/

func (self *DetailsViewer) newLabel ( id int ) fyne.CanvasObject {

  name,mcolor:= self.model.GetLabelInfo ( id )

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
  root  *fyne.Container // Contenedor intern que es modifica
  model DataModel
}


func NewDetailsViewer ( model DataModel ) *DetailsViewer {

  ret:= DetailsViewer{
    root : container.NewVBox (),
    model : model,
  }
  
  return &ret
  
} // end NewDetailsViewer


func (self *DetailsViewer) GetCanvas() fyne.CanvasObject { return self.root }


func (self *DetailsViewer) ViewEntry ( e Entry ) {
  
  // Neteja
  self.root.RemoveAll ()
  
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
      fmt.Println ( "EDIT BUTTON!!!!" )
    }),
    widget.NewToolbarAction ( theme.DeleteIcon (), func() {
      fmt.Println ( "DELETE BUTTON!!!!" )
    }),
  )
  
  // Afegeix
  tmp:= container.NewVBox ( container.NewHScroll ( card ), toolbar )
  self.root.Add ( tmp )
  
} // end ViewEntry


func (self *DetailsViewer) ViewFile ( f File ) {
  
  // Neteja
  self.root.RemoveAll ()


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
    f.GetType (),
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
    widget.NewToolbarAction ( theme.DeleteIcon (), func() {
      fmt.Println ( "DELETE BUTTON!!!!" )
    }),
  )
  
  // Afegeix
  tmp:= container.NewVBox ( container.NewHScroll ( card ), toolbar )
  self.root.Add ( tmp )
  
} // end ViewFile
