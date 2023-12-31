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
  "log"
  "os"
  "strings"

  "github.com/adriagipas/imgteka/view"
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
  self.v[id]= NewEntry ( self, id, name, platform_id, cover_id )
  
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
  db     *Database
  plats  *Platforms
  labels *Labels
  files  *Files
  dirs   *Dirs
  ids    []int64
  v      map[int64]*Entry
}


func NewEntries (

  db     *Database,
  plats  *Platforms,
  labels *Labels,
  files  *Files,
  dirs   *Dirs,
  
) *Entries {

  ret:= Entries{
    db     : db,
    dirs   : dirs,
    plats  : plats,
    labels : labels,
    files  : files,
    ids    : nil,
    v      : nil,
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


func (self *Entries) AddFileEntry(

  id        int64,
  path      string,
  name      string,
  file_type int,
  create_pb func() view.ProgressBar,
  
) error {
  
  // Obtindre entrada
  e,ok:= self.v[id]
  if !ok { // No deuria passar
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }

  // Afegeix
  if err:= self.files.Add ( e, path, name, file_type, create_pb ); err != nil {
    return err
  }
  
  return nil
  
} // end AddFileEntry


func (self *Entries) AddLabelEntry( id int64, label_id int ) error {

  if err:= self.db.RegisterEntryLabelPair ( id, label_id ); err != nil {
    return err
  }
  
  return nil
  
} // end AddLabelEntry


func (self *Entries) Filter( query *Query ) {

  self.db.SetQuery ( query )
  self.reset ()
  
} // end Filter


func (self *Entries) Get( id int64 ) *Entry {
  return self.v[id]
} // end Get


func (self *Entries) GetFile( id int64 ) *File {
  return self.files.Get ( id )
} // end GetFile


func (self *Entries) GetIDs() []int64 {
  return self.ids
} // end GetIDs


func (self *Entries) GetLabelIDs() []int {
  return self.labels.GetIDs ()
} // end GetLabelIDs


func (self *Entries) LoadFiles( id int64 ) error {

  // Comprova que existeix (no deuria passar que no)
  e,ok:= self.v[id]
  if !ok {
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }

  // Carrega
  if err:= self.db.LoadFilesEntry ( id, e ); err != nil {
    return err
  }
  
  return nil
  
} // end LoadFiles


func (self *Entries) LoadLabels( id int64 ) error {

  // Comprova que existeix (no deuria passar que no)
  e,ok:= self.v[id]
  if !ok {
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }

  // Carrega
  if err:= self.db.LoadLabelsEntry ( id, e ); err != nil {
    return err
  }

  return nil
  
} // end LoadLabels


func (self *Entries) Remove( id int64 ) error {

  // Comprova que no té fitxers.
  e,ok:= self.v[id]
  if !ok {
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


func (self *Entries) RemoveFileEntry( id int64, file_id int64 ) error {

  // Obtindre entrada
  e,ok:= self.v[id]
  if !ok { // No deuria passar
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }
  
  // Comprova que forma part de l'entrada
  f:= self.files.Get ( file_id )
  if f.GetEntryID () != id {
    return fmt.Errorf ( "La entrada (%id) no inclou el fitxer indicat (%d)",
      id, file_id)
  }
  
  // Elimina
  if err:= self.files.Remove ( file_id, e ); err != nil {
    return err
  }

  return nil
  
} // end RemoveFileEntry


func (self *Entries) RemoveLabelEntry( id int64, label_id int ) error {

  if err:= self.db.DeleteEntryLabelPair ( id, label_id ); err != nil {
    return err
  }
  
  return nil
  
} // end RemoveLabelEntry


func (self *Entries) SetCoverEntry( id int64, file_id int64 ) error {

  // Obtindre entrada
  _,ok:= self.v[id]
  if !ok { // No deuria passar
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }
  
  // Comprova que el fitxer pertany a l'entrada
  if file_id != -1 {
    f:= self.files.Get ( file_id )
    if f.GetEntryID () != id {
      return fmt.Errorf ( "La entrada (%id) no inclou el fitxer indicat (%d)",
        id, file_id)
    }
  }

  // Actualitza en la base de dades
  if err:= self.db.UpdateEntryCover ( id, file_id ); err != nil {
    return err
  }

  return nil
  
} // end SetCoverEntry


func (self *Entries) UpdateEntryName( id int64, name string ) error {

  // Prepara
  e:= self.v[id]
  
  // Intenta fer la transacció.
  if err:= self.db.UpdateEntryNameWithoutCommit ( id, name ); err != nil {
    return err
  }
  
  // Intenta el canvi de nom de directori
  plat:= self.plats.GetPlatform ( e.GetPlatformID () )
  dir_path,err:= self.dirs.GetEntryFolder ( plat.GetShortName (), e.GetName () )
  if err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  new_dir_path,err:= self.dirs.GetEntryFolder ( plat.GetShortName (), name )
  if err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  if err:= os.Remove ( new_dir_path ); err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  if err:= os.Rename ( dir_path, new_dir_path ); err != nil {
    err2:= self.db.RollbackLastTransaction ()
    if err2 != nil { log.Fatal ( err2 ) }
    return err
  }

  // Finalitza la transacció.
  err= self.db.CommitLastTransaction ()
  if err != nil { log.Fatal ( err ) }

  return nil
  
} // end UpdateEntryName


func (self *Entries) UpdateFileNameEntry(

  id      int64,
  file_id int64,
  name    string,

) error {

  // Obtindre entrada
  e,ok:= self.v[id]
  if !ok { // No deuria passar
    return fmt.Errorf ( "La entrada indicada (%d) no existeix", id )
  }
  
  // Comprova que forma part de l'entrada
  f:= self.files.Get ( file_id )
  if f.GetEntryID () != id {
    return fmt.Errorf ( "La entrada (%id) no inclou el fitxer indicat (%d)",
      id, file_id)
  }

  // Updateja el nom
  if err:= self.files.UpdateName ( file_id, e, name ); err != nil {
    return err
  }

  return nil
  
} // end UpdateFileNameEntry
