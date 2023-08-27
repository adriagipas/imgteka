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
 *  labels.go - Gestió de les etiquetes. Manté una "cache" de les
 *              etiquetes i es comunica amb la base de dades.
 */

package model

import (
  "errors"
  "fmt"
  "image/color"
  "log"
  "strings"
)




/****************/
/* PART PRIVADA */
/****************/

func (self *Labels) add(
  
  id         int,
  name       string,
  r,g,b      uint8,
  
) {

  self.ids= append ( self.ids, id )
  self.v[id]= &Label{
    labels : self,
    id     : id,
    name   : name,
    color  : color.RGBA{r,g,b,255},
  }
  
} // end add


func (self *Labels) reset() {

  // Reseteja
  self.ids= self.ids[:0]
  self.v= make(map[int]*Label)
  
  // Carrega
  if err:= self.db.LoadLabels ( self ); err != nil {
    log.Fatal ( err )
  }
  
} // end reset




/****************/
/* PART PÚBLICA */
/****************/

type Labels struct {
  db    *Database
  ids[] int
  v     map[int]*Label
}


func NewLabels ( db *Database ) *Labels {

  ret:= Labels{
    db  : db,
    ids : nil,
    v   : nil,
  }
  ret.reset ()
  
  return &ret
  
} // end NewLabels


func (self *Labels) GetIDs() []int { return self.ids }
func (self *Labels) Get( id int ) *Label { return self.v[id] }


func (self *Labels) Add( name string, c color.Color ) error {
  
  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }
  
  // Color
  r,g,b,_:= c.RGBA ()
  r8:= uint8((float32(r)/65535.0)*255.0 + 0.5)
  g8:= uint8((float32(g)/65535.0)*255.0 + 0.5)
  b8:= uint8((float32(b)/65535.0)*255.0 + 0.5)
  
  // Registra i reseteja
  if err:= self.db.RegisterLabel ( name, r8, g8, b8 ); err != nil {
    return fmt.Errorf ( "No s'ha pogut registrar la nova etiqueta: %s", err )
  }
  self.reset ()
  
  return nil

} // end Add


func (self *Labels) GetNumEntriesLabel( id int ) int64 {

  ret,err:= self.db.GetLabelNumEntries ( id )
  if err != nil { log.Fatal ( err ) }

  return ret
  
} // end GetNumEntriesLabel


func (self *Labels) Remove( id int ) error {

  // Comprova que és una plataforma que no s'utilitza
  label:= self.v[id]
  if label == nil {
    return fmt.Errorf ( "L'etiqueta indicada (%d) no existeix", id )
  }
  if label.GetNumEntries() > 0 {
    return errors.New ( "No es pot esborrar l'etiqueta perquè està en ús" )
  }

  // Elimina i reseteja.
  if err:= self.db.DeleteLabel ( id ); err != nil {
    return fmt.Errorf ( "No s'ha pogut esborrar l'etiqueta: %s", err )
  }
  self.reset ()
  
  return nil
  
} // end Remove


func (self *Labels) UpdateLabel(
  
  id    int,
  name  string,
  r,g,b uint8,
  
) error {

  if err:= self.db.UpdateLabel ( id, name, r, g, b ); err != nil {
    return fmt.Errorf ( "No s'ha pogut actualitzar l'etiqueta: %s", err )
  }
  
  return nil
  
} // end UpdateLabel




// LABEL ///////////////////////////////////////////////////////////////////////

type Label struct {
  labels *Labels
  id     int
  name   string
  color  color.Color
}


func (self *Label) GetName() string { return self.name }
func (self *Label) GetColor() color.Color { return self.color }


func (self *Label) GetNumEntries() int64 {
  return self.labels.GetNumEntriesLabel ( self.id )
} // end GetNumEntries


func (self *Label) Update( name string, c color.Color ) error {

  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }

  // Color
  r,g,b,_:= c.RGBA ()
  r8:= uint8((float32(r)/65535.0)*255.0 + 0.5)
  g8:= uint8((float32(g)/65535.0)*255.0 + 0.5)
  b8:= uint8((float32(b)/65535.0)*255.0 + 0.5)

  // Intenta actualitzar en la base de dades
  if err:= self.labels.UpdateLabel ( self.id, name, r8, g8, b8 ); err != nil {
    return err
  }
  
  // Modifica els atributs "cached"
  self.color= c
  self.name= name
  
  return nil
  
} // end Update
