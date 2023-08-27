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
 *  entry.go - Representa una entrada.
 */

package model

import (
  "errors"
  "fmt"
  "image"
  "log"
  "strings"
)




/****************/
/* PART PRIVADA */
/****************/

func (self *Entry) addLabelID( id int ) {

  self.labels.ids= append ( self.labels.ids, id )
  self.labels.ids_map[id]= true
  
} // end addLabelID


func (self *Entry) resetLabels() {

  // Prepara
  self.labels.ids= self.labels.ids[:0]
  self.labels.ids_map= make(map[int]bool)

  // Carrega valors
  if err:= self.entries.LoadLabels ( self.id ); err != nil {
    log.Fatal ( err )
  }
  
  // Marca com carregat
  self.labels.loaded= true
  
} // end resetLabels




/****************/
/* PART PÚBLICA */
/****************/

type Entry struct {

  // Part bàsica
  entries  *Entries
  id       int64
  name     string
  platform int
  cover    int64

  // Relacionat amb les etiquetes
  labels struct {
    loaded  bool // Indica si s'ha inicialitzat
    ids     []int
    ids_map map[int]bool
    uids    []int
  }
  
}


func (self *Entry) AddLabel( id int ) error {

  // Afegeix
  if err:= self.entries.AddLabelEntry( self.id, id ); err != nil {
    return err
  }
  
  // Reseteja
  self.resetLabels ()
  
  return nil
  
} // end AddLabel


func (self *Entry) GetName() string { return self.name }
func (self *Entry) GetPlatformID() int { return self.platform }
func (self *Entry) GetFileIDs() []int64 {
  fmt.Println ( "TODO Entry.GetFileIDs !" )
  return nil
}


func (self *Entry) GetCover() image.Image {

  var ret image.Image
  if ( self.cover != -1 ) {
    fmt.Println ( "TODO Entry.GetCover !" )
    ret= nil
  } else {
    ret= nil
  }

  return ret
  
} // end GetCover


func (self *Entry) GetLabelIDs() []int {

  // Carrega si no s'ha carregat mai
  if !self.labels.loaded {
    self.resetLabels ()
  }
  
  return self.labels.ids
  
} // end GetLabelIDs


func (self *Entry) GetUnusedLabelIDs() []int {

  // Carrega si no s'ha carregat mai
  if !self.labels.loaded {
    self.resetLabels ()
  }

  // Crea el vector de unused
  self.labels.uids= self.labels.uids[:0]
  lids:= self.entries.GetLabelIDs ()
  for _,id:= range lids {
    if _,ok:= self.labels.ids_map[id]; !ok {
      self.labels.uids= append(self.labels.uids,id)
    }
  }
  
  return self.labels.uids
  
} // end GetUnusedLabelIDs


func (self *Entry) RemoveLabel( id int ) error {

  // Carrega si no s'ha carregat mai (No deuria passar)
  if !self.labels.loaded {
    self.resetLabels ()
  }
  
  // Comprova que existeix
  if _,ok:= self.labels.ids_map[id]; !ok {
    return fmt.Errorf ( "L'etiqueta indicada (%d) no format part de l'entrada",
      id )
  }
  
  // Elimina
  if err:= self.entries.RemoveLabelEntry( self.id, id ); err != nil {
    return err
  }
  
  // Reseteja
  self.resetLabels ()
  
  return nil
  
} // end RemoveLabel


func (self *Entry) UpdateName( name string ) error {

  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }
  
  // Actualitza nom
  if err:= self.entries.UpdateEntryName ( self.id, name ); err != nil {
    return err
  }
  
  // Modifica l'atribut
  self.name= name
  
  return nil
  
} // end UpdateName
