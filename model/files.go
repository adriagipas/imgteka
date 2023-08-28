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
  "os"

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




// FILES ///////////////////////////////////////////////////////////////////////

type Files struct {
  db *Database
  v  map[int64]*File
}


func NewFiles ( db *Database ) *Files {

  ret:= Files{
    db  : db,
    v   : nil,
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

  // Comprova existeix
  pb.Set ( "Comprova que existeix...", 0.1 )
  f,err:= os.Open ( path )
  if err != nil {
    return fmt.Errorf ( "No s'ha pogut obrir el fitxer '%s': %s", path, err )
  }
  defer f.Close ()

  // Comprova tipus i obté metadades
  pb.Set ( "Comprova tipus i obté metadades...", 0.2 )
  ft,err:= file_type.Get ( ftype )
  if err != nil { return err }
  md,err:= ft.GetMetadata ( f )
  if err != nil { return err }
  fmt.Println ( "Metadata", md )
  
  // Calcula MD5
  pb.Set ( "Calcula MD5...", 0.3 )
  md5,err:= calcMD5 ( f )
  if err != nil { return err }
  fmt.Println ( "MD5", md5 )

  // Calcula SHA1
  pb.Set ( "Calcula SHA1...", 0.4 )
  sha1,err:= calcSHA1 ( f )
  if err != nil { return err }
  fmt.Println ( "SHA1", sha1 )
  
  return nil
  
} // end Add


func (self *Files) Get( id int64 ) *File {

  ret,ok:= self.v[id]
  if !ok {
    name,entry_id,file_type,size,md5,sha1,json,last_check:= 
      self.db.GetFile ( id )
    ret= NewFile ( id, name, entry_id, file_type, size, md5,
      sha1, json, last_check )
    self.v[id]= ret
  }
  
  return ret
  
} // end Get




// FILE ////////////////////////////////////////////////////////////////////////

type File struct {
  id        int64
  name      string
  entry     int64
  file_type int
  size      int64
  md5       string
  sha1      string
  // TODO!! algunes coses com last_check i json
}


func NewFile(
  
  id         int64,
  name       string,
  entry      int64,
  file_type  int,
  size       int64,
  md5        string,
  sha1       string,
  json       string,
  last_check int64,
  
) *File {

  ret:= File{
    id        : id,
    name      : name,
    entry     : entry,
    file_type : file_type,
    size      : size,
    md5       : md5,
    sha1      : sha1,
  }

  return &ret
  
} // end NewFile


func (self *File) GetMetadata() []view.StringPair {
  fmt.Println ( "TODO File.GetMetaData !" )
  return make([]view.StringPair,0)
} // end GetMetaData


func (self *File) GetName() string { return self.name }
func (self *File) GetTypeID() int { return self.file_type }


type MetadataValue struct {
  key   string
  value string
}


func (self *MetadataValue) GetKey() string { return self.key }
func (self *MetadataValue) GetValue() string { return self.value }
