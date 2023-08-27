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
 *  entries.go - Gestió de les entrades. Manté una "cache" de les
 *               entrades i es comunica amb la base de dades.
 */

package model

import (
  "errors"
  "fmt"
  "image"
  "log"
  "os"
  "strings"
)




/****************/
/* PART PRIVADA */
/****************/

func (self *Entries) add(
  
  id          int64,
  name        string,
  platform_id int,
  cover_id    int64,
  
) {

  self.ids= append ( self.ids, id )
  self.v[id]= &Entry{
    entries  : self,
    id       : id,
    name     : name,
    platform : platform_id,
    cover    : cover_id,
  }
  
} // end add


func (self *Entries) reset() {

  // Reseteja
  self.ids= self.ids[:0]
  self.v= make(map[int64]*Entry)
  
  // Carrega
  if err:= self.db.LoadEntries ( self ); err != nil {
    log.Fatal ( err )
  }
  
} // end reset




/****************/
/* PART PÚBLICA */
/****************/

type Entries struct {
  db    *Database
  plats *Platforms
  dirs  *Dirs
  ids   []int64
  v     map[int64]*Entry
}


func NewEntries ( db *Database, plats *Platforms, dirs *Dirs ) *Entries {

  ret:= Entries{
    db    : db,
    dirs  : dirs,
    plats : plats,
    ids   : nil,
    v     : nil,
  }
  ret.reset ()

  return &ret
  
} // end NewEntries


func (self *Entries) Add( name string, platform_id int ) error {

  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }

  // Intenta registrar
  // --> Intenta transacció
  if err:= self.db.RegisterEntryWithoutCommit (
    name, platform_id ); err != nil {
    return fmt.Errorf ( "No s'ha pogut registrar la nova entrada: %s", err )
  }
  // --> Intenta creació directòri
  plat:= self.plats.GetPlatform ( platform_id )
  dir_path,err:= self.dirs.GetEntryFolder ( plat.GetShortName (), name )
  if err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  // --> Finalitza transacció
  err= self.db.CommitLastTransaction ()
  if err != nil {
    err2:= os.Remove ( dir_path )
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  
  // Reseteja
  self.reset ()
  
  return nil
  
} // end Add


func (self *Entries) Get( id int64 ) *Entry {
  return self.v[id]
} // end Get


func (self *Entries) GetIDs() []int64 {
  return self.ids
} // end GetIDs


func (self *Entries) Remove( id int64 ) error {

  // Comprova que no té fitxers.
  e:= self.v[id]
  if e == nil {
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }
  if len(e.GetFileIDs ()) > 0 {
    return errors.New ( "No es pot esborrar una entrada amb fitxers" )
  }
  
  // Elimina
  // --> Intenta transacció
  if err:= self.db.DeleteEntryWithoutCommit ( id ); err != nil {
    return fmt.Errorf ( "No s'ha pogut esborrar l'entrada: %s", err )
  }
  // --> Intenta eliminar directori
  plat:= self.plats.GetPlatform ( e.GetPlatformID () )
  dir_path,err:= self.dirs.GetEntryFolder ( plat.GetShortName (), e.GetName () )
  if err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  err= os.Remove ( dir_path )
  if err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  // --> Finalitza transacció
  err= self.db.CommitLastTransaction ()
  if err != nil { log.Fatal ( err ) }
  
  // Reseteja
  self.reset ()

  return nil
  
} // end Remove


type Entry struct {

  // Part bàsica
  entries  *Entries
  id       int64
  name     string
  platform int
  cover    int64
  
}


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
  fmt.Println ( "TODO Entry.GetLabelIDs !" )
  return nil
}
