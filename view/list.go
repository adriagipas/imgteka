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
 *  list.go - Implementa el llistat de les imatges.
 */

package view

import (
  "fmt"
  "image/color"
  "log"
  "strconv"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)




/*************/
/* CONSTANTS */
/*************/

const FILE_VIEW_TEMPLATE string= `
#### %s

- ***Tipus:*** *%s*

- ***md5:*** *%s*

- ***sha1:*** *%s*

- ***Grandària:*** *%s*`




/***********/
/* FACTORY */
/***********/

func int64_to_tnid(ids []int64,pref string) []widget.TreeNodeID {
  ret:= make([]widget.TreeNodeID,len(ids))
  for i:= 0; i < len(ids); i++ {
    ret[i]= pref + strconv.FormatInt ( ids[i], 10 )
  }
  return ret
}


func size_to_string ( size int64 ) string {

  var ret string
  
  // Bytes
  if size < 1024 {
    ret= fmt.Sprintf ( "%d B" )
  } else if size < 1024*1024 {
    ret= fmt.Sprintf ( "%d B (%.1f KB)", size, float32(size)/1024.0 )
  } else if size < 1024*1024*1024 {
    ret= fmt.Sprintf ( "%d B (%.1f MB)", size, float32(size)/(1024.0*1024.0) )
  } else {
    ret= fmt.Sprintf ( "%d B (%.1f GB)",
      size, float32(size)/(1024.0*1024.0*1024) )
  }
  
  return ret
  
} // end size_to_string


func newPlatformLabel(idname string, mcolor color.Color) fyne.CanvasObject {
  
  text:= canvas.NewText ( idname, color.Black )
  text.Alignment= fyne.TextAlignCenter
  text.TextStyle= fyne.TextStyle{Bold: true}
  text.TextSize= 10.0
  rect:= canvas.NewRectangle ( mcolor )
  rect.SetMinSize ( fyne.Size{7.5,1.0} )
  
  ret:= container.NewHBox ( rect, text )
  
  return ret
  
} // end newPlatformLabel


func newEntryView() fyne.CanvasObject {

  // Plataforma
  plat:= newPlatformLabel ( "BLO", color.Black )
  
  // Nom
  name:= widget.NewLabel ( "NAME" )
  
  return container.NewHBox ( plat, name )
  
} // end newEntryView


func newFileView() fyne.CanvasObject {

  // Descripció
  text:= fmt.Sprintf ( FILE_VIEW_TEMPLATE, "Blo", "Blo",
    "Blo", "Blo", size_to_string ( 999 ) )
  desc:= widget.NewRichTextFromMarkdown ( text )
  
  return container.NewMax ( desc )
  
} // end newFileView


type _Factory struct {
  model           DataModel
  dv              *DetailsViewer
  platform_labels map[int]fyne.CanvasObject
}


func newFactory (
  model DataModel,
  dv    *DetailsViewer,
) *_Factory {
  plat:= make(map[int]fyne.CanvasObject)
  ret:= _Factory{model,dv,plat}
  return &ret
}


func (self *_Factory) getPlatformLabel(id int) fyne.CanvasObject {

  // Busca si ja el tenim
  ret,ok:= self.platform_labels[id]
  if ok { return ret }

  // Crea el valor
  idname,color:= self.model.GetPlatformHints ( id )
  ret= newPlatformLabel ( idname, color )
  self.platform_labels[id]= ret

  return ret
  
} // end getPlarformLabel


func (self *_Factory) setEntryView (
  o fyne.CanvasObject,
  e Entry,
) {

  cont:= o.(*fyne.Container)

  // Etiqueta plataforma
  cont.Objects[0]= self.getPlatformLabel ( e.GetPlatformID () )
  
  // Nom
  name:= cont.Objects[1].(*widget.Label)
  name.SetText ( e.GetName () )
  
} // end setEntryView


func (self *_Factory) setFileView (o fyne.CanvasObject, f File) {

  cont:= o.(*fyne.Container)

  // Descripció
  text:= fmt.Sprintf ( FILE_VIEW_TEMPLATE,
    f.GetName (), f.GetType (), f.GetMD5 (), f.GetSHA1 (),
    size_to_string ( f.GetSize () ) )
  cont.Objects[0]= widget.NewRichTextFromMarkdown ( text )
  
} // end setFileView

func (self *_Factory) childUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
  if id == "" {
    return int64_to_tnid ( self.model.RootEntries (), "E" )
  } else if id[0] == 'E' {
    id,err:= strconv.ParseInt ( id[1:], 10, 64 )
    if err != nil { log.Fatal ( err ) }
    e:= self.model.GetEntry ( id )
    return int64_to_tnid ( e.GetFileIDs (), "F" )
  } else {
    return []string{}
  }
}


func (self *_Factory) isBranch(id widget.TreeNodeID) bool {
  return id == "" || id[0] == 'E'
}


func (self *_Factory) create(branch bool) fyne.CanvasObject {
  if branch {
    return newEntryView ()
  } else {
    return newFileView ()
  }
}


func (self *_Factory) update(
  id     widget.TreeNodeID,
  branch bool,
  o      fyne.CanvasObject,
) {
  if branch {
    id,err:= strconv.ParseInt ( id[1:], 10, 64 )
    if err != nil { log.Fatal ( err ) }
    e:= self.model.GetEntry ( id )
    self.setEntryView ( o, e )
  } else {
    id,err:= strconv.ParseInt ( id[1:], 10, 64 )
    if err != nil { log.Fatal ( err ) }
    f:= self.model.GetFile ( id )
    self.setFileView ( o, f )
  }
} // end update


func (self *_Factory) onSelected ( id widget.TreeNodeID ) {
  if id[0] == 'E' { // Entrada
    id,err:= strconv.ParseInt ( id[1:], 10, 64 )
    if err != nil { log.Fatal ( err ) }
    e:= self.model.GetEntry ( id )
    self.dv.ViewEntry ( e )
  } else if id[0] == 'F' { // Fitxer
    id,err:= strconv.ParseInt ( id[1:], 10, 64 )
    if err != nil { log.Fatal ( err ) }
    f:= self.model.GetFile ( id )
    self.dv.ViewFile ( f )
  } else { // ¿¿??
    log.Fatal ( "list.go - onSelected - WTF!!!" )
  }
} // onSelected




/****************/
/* PART PÚBLICA */
/****************/

type List struct {
  widget.Tree
}


func NewList ( model DataModel, dv *DetailsViewer ) *List {

  f:= newFactory ( model, dv )
  ret:= &List{}
  ret.ExtendBaseWidget ( ret )
  ret.ChildUIDs= func(id widget.TreeNodeID) []widget.TreeNodeID {
    return f.childUIDs ( id )
  }
  ret.IsBranch= func(id widget.TreeNodeID) bool {
    return f.isBranch ( id )
  }
  ret.CreateNode= func(branch bool) fyne.CanvasObject {
    return f.create ( branch )
  }
  ret.UpdateNode= func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
    f.update ( id, branch, o )
  }
  ret.OnSelected= func(id widget.TreeNodeID) {
    f.onSelected ( id )
  }
  
  return ret
  
} // end GetList
