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
 *  file.go - Representa un fitxer.
 */

package model

import (
  "fmt"
  "image"
  "log"
  
  "github.com/adriagipas/imgteka/model/file_type"
  "github.com/adriagipas/imgteka/view"
)




// UTILS ///////////////////////////////////////////////////////////////////////

func size2text( size int64 ) string {

  var ret string
  
  if size < 1024 { // Bytes
    ret= fmt.Sprintf ( "%d B", size )
  } else if size < 1024*1024 { // KB
    ret= fmt.Sprintf ( "%.1f KB (%d B)", float32(size)/1024, size )
  } else if size < 1024*1024*1024 { // MB
    ret= fmt.Sprintf ( "%.1f MB (%d B)", float32(size)/(1024*1024), size )
  } else {
    ret= fmt.Sprintf ( "%.1f GB (%d B)", float32(size)/(1024*1024*1024), size )
  }

  return ret
  
} // end size2text




// FILE ////////////////////////////////////////////////////////////////////////

type File struct {
  
  dirs         *Dirs
  id           int64
  name         string
  entry        int64
  file_type_id int
  file_type    file_type.FileType
  size         int64
  md5          string
  sha1         string

  // Metadata
  md []view.StringPair
  
}


func NewFile(

  dirs        *Dirs,
  id           int64,
  name         string,
  entry        int64,
  file_type_id int,
  size         int64,
  md5          string,
  sha1         string,
  json         string,
  last_check   int64,
  
) *File {

  // Crea objecte
  ret:= File{
    dirs         : dirs,
    id           : id,
    name         : name,
    entry        : entry,
    file_type_id : file_type_id,
    size         : size,
    md5          : md5,
    sha1         : sha1,
  }
  var err error
  ret.file_type,err= file_type.Get ( file_type_id )
  if err != nil { log.Fatal ( err ) }

  // Crea metadata
  ret.md= make([]view.StringPair,3)
  ret.md[0]= &MetadataValue{"md5",md5}
  ret.md[1]= &MetadataValue{"sha1",sha1}
  ret.md[2]= &MetadataValue{"Grandària",size2text ( size )}
  ret.md= ret.file_type.ParseMetadata ( ret.md, json )
  
  return &ret
  
} // end NewFile


func (self *File) GetEntryID() int64 { return self.entry }


func (self *File) GetImage() image.Image {

  var ret image.Image
  var err error
  
  if self.file_type.IsImage () {
    ret,err= self.file_type.GetImage ( self.GetPath () )
    if err != nil {
      log.Printf ( "Error al intentar llegir la imatge de '%s': %s",
        self.GetPath (), err )
      ret= nil
    }
  } else {
    ret= nil
  }
  
  return ret
  
} // end GetImage


func (self *File) GetMetadata() []view.StringPair { return self.md }


func (self *File) GetName() string { return self.name }


func (self *File) GetPath() string {

  ret,err:= self.dirs.GetFileNameFiles (
    self.file_type.GetShortName (), self.name )
  if err != nil { log.Fatal ( err ) }

  return ret
  
} // end GetPath


func (self *File) GetTypeID() int { return self.file_type_id }


func (self *File) IsImage() bool {
  return self.file_type.IsImage ()
} // end IsImage




// METADATAVALUE ///////////////////////////////////////////////////////////////

type MetadataValue struct {
  key   string
  value string
}


func (self *MetadataValue) GetKey() string { return self.key }
func (self *MetadataValue) GetValue() string { return self.value }
