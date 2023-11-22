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
 *  files.go - Gestió dels fitxers. Manté una "cache" dels fitxers
 *             consultats.
 */

package model

import (
  "crypto/md5"
  "crypto/sha1"
  "fmt"
  "io"
  "log"
  "os"
  "time"
  
  "github.com/adriagipas/imgteka/model/file_type"
  "github.com/adriagipas/imgteka/view"
)




// UTILS ///////////////////////////////////////////////////////////////////////

func calcMD5( f *os.File ) (string,error) {

  // Rebobina
  if _,err:= f.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el MD5: %s", err )
  }

  // Calcula MD5
  h:= md5.New ()
  if _,err:= io.Copy ( h, f ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el MD5: %s", err )
  }

  ret:= fmt.Sprintf ( "%x", h.Sum ( nil ) )

  return ret,nil
  
} // end calcMD5


func calcSHA1( f *os.File ) (string,error) {

  // Rebobina
  if _,err:= f.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el SHA1: %s", err )
  }

  // Calcula MD5
  h:= sha1.New ()
  if _,err:= io.Copy ( h, f ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el SHA1: %s", err )
  }
  
  ret:= fmt.Sprintf ( "%x", h.Sum ( nil ) )
  
  return ret,nil
  
} // end calcSHA1


func linkFile( oldname string, newname string ) error {

  if err:= os.Link ( oldname, newname ); err != nil {
    return fmt.Errorf ( "No s'ha pogut crear enllaç per a '%s': %s",
      oldname, err )
  }
  if err:= os.Chmod ( newname, 0400 ); err != nil {
    if err2:= os.Remove ( newname ); err2 != nil {
      log.Fatal ( err2 )
    }
    return fmt.Errorf ( "No s'han pogut canviar els permisos de '%s': %s",
      newname, err )
  }

  return nil
  
} // end linkFile


func readLink( path string ) (string,error) {
  
  info,err:= os.Lstat ( path )
  if err != nil { return path,err }
  for ;(info.Mode ()&os.ModeSymlink) != 0; {
    path,err= os.Readlink ( path )
    if err != nil { return path,err }
    info,err= os.Lstat ( path )
    if err != nil { return path,err }
  }
  
  return path,nil
  
} // end readLink




// FILES ///////////////////////////////////////////////////////////////////////

type Files struct {
  db    *Database
  plats *Platforms
  dirs  *Dirs
  cmds  *Commands
  v     map[int64]*File
}


func NewFiles (
  
  db    *Database,
  plats *Platforms,
  dirs  *Dirs,
  cmds  *Commands,

) *Files {

  ret:= Files{
    db    : db,
    plats : plats,
    dirs  : dirs,
    cmds  : cmds,
    v     : nil,
  }
  ret.v= make(map[int64]*File)
  
  return &ret
  
} // end NewFiles


func (self *Files) Add(

  e         *Entry,
  path      string,
  name      string,
  ftype     int,
  create_pb func() view.ProgressBar,

) error {
  
  // Crea barra de progress
  pb:= create_pb ()
  defer pb.Close ()

  // Comprova existeix i grandària
  pb.Set ( "Comprova que existeix...", 0.1 )
  path,err:= readLink ( path )
  if err != nil {
    return fmt.Errorf ( "No s'ha pogut accedir al fitxer '%s': %s", path, err )
  }
  f,err:= os.Open ( path )
  if err != nil {
    return fmt.Errorf ( "No s'ha pogut obrir el fitxer '%s': %s", path, err )
  }
  defer f.Close ()
  info,err:= f.Stat ()
  if err != nil {
    return fmt.Errorf ( "No s'ha pogut llegir el fitxer '%s': %s", path, err )
  }
  size:= info.Size ()
  
  // Comprova tipus i obté metadades
  pb.Set ( "Comprova tipus i obté metadades...", 0.2 )
  ft,err:= file_type.Get ( ftype )
  if err != nil { return err }
  md,err:= ft.GetMetadata ( path )
  if err != nil { return err }
  
  // Calcula MD5
  pb.Set ( "Calcula MD5...", 0.3 )
  md5,err:= calcMD5 ( f )
  if err != nil { return err }

  // Calcula SHA1
  pb.Set ( "Calcula SHA1...", 0.4 )
  sha1,err:= calcSHA1 ( f )
  if err != nil { return err }

  // Obté noms i stamp
  plat_name:= self.plats.GetPlatform ( e.GetPlatformID () ).GetShortName ()
  ename,err:= self.dirs.GetFileNameEntries ( plat_name, e.GetName (), name )
  if err != nil { return err }
  fname,err:= self.dirs.GetFileNameFiles ( ft.GetShortName (), name )
  if err != nil { return err }
  time_now:= time.Now ().Unix ()
  
  // Intenta commit
  pb.Set ( "Insereix en base de dades...", 0.6 )
  if err:= self.db.RegisterFileWithoutCommit ( name, e.GetID (), ftype,
    size, md5, sha1, md, time_now ); err != nil {
    return err
  }

  // Crea fitxers
  pb.Set ( "Desa fitxers en disc...", 0.7 )
  if err:= linkFile ( path, ename ); err != nil {
    self.db.RollbackLastTransaction ()
    return err
  }
  if err:= linkFile ( path, fname ); err != nil {
    if err2:= os.Remove ( ename ); err2 != nil { log.Fatal ( err2 )}
    self.db.RollbackLastTransaction ()
    return err
  }
  
  // Consolida commit
  if err:= self.db.CommitLastTransaction (); err != nil {
    if err2:= os.Remove ( ename ); err2 != nil {
      os.Remove ( fname )
      log.Fatal ( err2 )
    }
    if err2:= os.Remove ( fname ); err2 != nil { log.Fatal ( err2 ) }
    return err
  }
  
  return nil
  
} // end Add


func (self *Files) Get( id int64 ) *File {

  ret,ok:= self.v[id]
  if !ok {
    name,entry_id,file_type,size,md5,sha1,json,last_check:= 
      self.db.GetFile ( id )
    ret= NewFile ( self.dirs, self.cmds, id, name, entry_id,
      file_type, size, md5, sha1, json, last_check )
    self.v[id]= ret
  }
  
  return ret
  
} // end Get


func (self *Files) Remove( id int64, e *Entry ) error {

  // Obté fitxer
  f:= self.Get ( id )

  // Obté noms
  plat_name:= self.plats.GetPlatform ( e.GetPlatformID () ).GetShortName ()
  ename,err:= self.dirs.GetFileNameEntries (
    plat_name, e.GetName (), f.GetName () )
  if err != nil { return err }
  ft,err:= file_type.Get ( f.GetTypeID () )
  if err != nil { return err }
  fname,err:= self.dirs.GetFileNameFiles ( ft.GetShortName (), f.GetName () )
  if err != nil { return err }
  tname,err:= self.dirs.GetFileNameTemp ( ft.GetShortName (), f.GetName () )
  if err != nil { return err }
  defer os.Remove ( tname )
  
  // Crea fitxer temporal
  if err:= linkFile ( fname, tname ); err != nil {
    return err
  }
  
  // Intenta commit
  if err:= self.db.DeleteFileWithoutCommit ( id ); err != nil {
    return err
  }
  
  // Intenta esborrar fitxers
  err_e:= os.Remove ( ename )
  err_f:= os.Remove ( fname )
  if err_e != nil || err_f != nil {
    self.db.RollbackLastTransaction ()
    if err_e != nil {
      if err2:= linkFile ( tname, ename ); err2 != nil { log.Fatal ( err2 ) }
    }
    if err_f != nil {
      if err2:= linkFile ( tname, fname ); err2 != nil { log.Fatal ( err2 ) }
    }
    if err_f != nil {
      return err_f
    } else {
      return err_e
    }
  }
  
  // Força commit. Si falla intente recuperar
  if err:= self.db.CommitLastTransaction (); err != nil {
    if err2:= linkFile ( tname, ename ); err2 != nil { log.Fatal ( err2 ) }
    if err2:= linkFile ( tname, fname ); err2 != nil { log.Fatal ( err2 ) }
    return err
  }

  // Esborra del mapa
  delete(self.v,id)
  
  return nil
  
} // end Remove


func (self *Files) UpdateName( id int64, e *Entry, new_name string ) error {

  // Obté fitxer
  f:= self.Get ( id )

  // Obté nom entry
  plat_name:= self.plats.GetPlatform ( e.GetPlatformID () ).GetShortName ()
  old_ename,err:= self.dirs.GetFileNameEntries (
    plat_name, e.GetName (), f.GetName () )
  if err != nil { return err }
  new_ename,err:= self.dirs.GetFileNameEntries (
    plat_name, e.GetName (), new_name )
  if err != nil { return err }
  
  // Obté nom files
  ft,err:= file_type.Get ( f.GetTypeID () )
  if err != nil { return err }
  old_fname,err:= self.dirs.GetFileNameFiles ( ft.GetShortName (),
    f.GetName () )
  if err != nil { return err }
  new_fname,err:= self.dirs.GetFileNameFiles ( ft.GetShortName (), new_name )
  if err != nil { return err }

  // Intenta commit
  if err:= self.db.UpdateFileNameWithoutCommit ( id, new_name ); err != nil {
    return err
  }

  // Reanomena
  if err:= os.Rename ( old_ename, new_ename ); err != nil {
    self.db.RollbackLastTransaction ()
    return err
  }
  if err:= os.Rename ( old_fname, new_fname ); err != nil {
    self.db.RollbackLastTransaction ()
    // Intenta recuperar...
    if err2:= os.Rename ( new_ename, old_ename ); err2 != nil {
      log.Fatal ( err2 )
    }
    return err
  }

  // Intenta commit final
  if err:= self.db.CommitLastTransaction (); err != nil {
    // Intent desesperat per recuperar...
    err_e:= os.Rename ( new_ename, old_ename )
    err_f:= os.Rename ( new_fname, old_fname )
    if err_e != nil { log.Fatal ( err_e ) }
    if err_f != nil { log.Fatal ( err_f ) }
    return err
  }

  // Esborra del mapa per forçar que es torne a carregar.
  delete(self.v,id)
  
  return nil
  
} // end UpdateName
