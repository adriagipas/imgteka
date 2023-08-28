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
  "fmt"
  "time"

  "github.com/adriagipas/imgteka/view"
)




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
  file_type int,
  create_pb func() view.ProgressBar,

) error {

  // Crea barra de progress
  pb:= create_pb ()

  pb.Set ( "Missatge 1 ...", 0.1 )
  time.Sleep ( time.Second )

  pb.Set ( "Missatge 2 ...", 0.5 )
  time.Sleep ( time.Second )

  pb.Set ( "Missatge 3 ...", 0.8 )
  time.Sleep ( time.Second )

  pb.Set ( "Missatge 4 ...", 1 )
  time.Sleep ( time.Second )

  pb.Close ()
  
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
