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
  "image/png"
  "log"
  "os"
  
  "github.com/adriagipas/imgteka/model/file_type"
  "github.com/adriagipas/imgteka/view"
  "github.com/nfnt/resize"
)




/****************/
/* PART PRIVADA */
/****************/

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


// Torna nil si no està cacheada o no s'ha pogut carregar.
func loadCachedImage( max_wh int, file_name string ) image.Image {

  f,err:= os.Open ( file_name )
  if err != nil { return nil }
  defer f.Close ()

  ret,err:= png.Decode ( f )
  if err != nil {
    log.Printf ( "No s'ha pogut carregar imatge del fitxer '%s'"+
      " de memòria cau: %s", file_name, err )
    return nil
  }
  
  return ret
  
} // end loadCachedImage


// Reescala i cachea. Torna la imatge reescalada.
func resizeAndCacheImage(

  img       image.Image,
  max_wh    int,
  file_name string,

) image.Image {

  // Reescala
  bounds:= img.Bounds ()
  width:= bounds.Max.X - bounds.Min.X
  height:= bounds.Max.Y - bounds.Min.Y
  if width >= height {
    img= resize.Resize ( uint(max_wh), 0, img, resize.NearestNeighbor )
  } else {
    img= resize.Resize ( 0, uint(max_wh), img, resize.NearestNeighbor )
  }

  // Desa en memòria cau.
  f,err:= os.Create ( file_name )
  if err != nil {
    log.Printf ( "No s'ha pogut creat el fitxer '%s': %s", file_name, err ) 
  }
  defer f.Close ()
  if err:= png.Encode ( f, img ); err != nil {
    log.Printf ( "No s'ha pogut desar la imatge en '%s': %s", file_name, err )
  }

  return img
  
} // resizeAndCacheImage




/****************/
/* PART PÚBLICA */
/****************/

type File struct {
  
  dirs         *Dirs
  cmds         *Commands
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
  cmds        *Commands,
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
    cmds         : cmds,
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


func (self *File) GetImage( max_wh int ) image.Image {

  var ret image.Image
  
  if self.file_type.IsImage () {

    // Nom en la cache
    cache_name:= fmt.Sprintf ( "%d-%s.png", self.id, self.md5)
    cache_fn,err:= self.dirs.GetCachedImageName ( max_wh, cache_name )
    if err != nil {
      log.Printf ( "Error inesperat en File.GetImage: %s", err )
      return nil
    }

    // Prova en la cache.
    if ret= loadCachedImage ( max_wh, cache_fn ); ret == nil {

      var err error
      
      // Si no està carrega original.
      ret,err= self.file_type.GetImage ( self.GetPath () )
      if err != nil {
        log.Printf ( "Error al intentar llegir la imatge de '%s': %s",
          self.GetPath (), err )
        ret= nil
      } else {
        // Si és molt gran intenta cache
        bounds:= ret.Bounds ()
        width:= bounds.Max.X - bounds.Min.X
        height:= bounds.Max.Y - bounds.Min.Y
        if width > max_wh || height > max_wh {
          ret= resizeAndCacheImage ( ret, max_wh, cache_fn )
        }
      }
      
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


func (self *File) Run() error {
  return self.cmds.Run ( self.file_type_id, self.GetPath () )
} // end Run




// METADATAVALUE ///////////////////////////////////////////////////////////////

type MetadataValue struct {
  key   string
  value string
}


func (self *MetadataValue) GetKey() string { return self.key }
func (self *MetadataValue) GetValue() string { return self.value }
