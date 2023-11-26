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
 *  zip.go - Tipus de fitxer TAR (amb compressió).
 */

package file_type

import (
  "archive/tar"
  "compress/gzip"
  "encoding/json"
  "fmt"
  "image"
  "io"
  "log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




type _TAR_Metadata struct {

  NFiles  int
  
}


func _TAR_Open( fd *os.File ) (*tar.Reader,error) {

  // Prova sense compressió.
  ret:= tar.NewReader ( fd )
  _,err:= ret.Next ()
  if err == nil { // És pot llegir
    if _,err:= fd.Seek ( 0, 0 ); err != nil { return nil,err }
    ret= tar.NewReader ( fd )
    return ret,nil
  }

  // Prova gzip
  if _,err:= fd.Seek ( 0, 0 ); err != nil { return nil,err }
  zr,err:= gzip.NewReader ( fd )
  if err == nil { // Es pot interpretar com gzip
    ret= tar.NewReader ( zr )
    _,err:= ret.Next ()
    if err == nil { // És pot llegir
      if _,err:= fd.Seek ( 0, 0 ); err != nil { return nil,err }
      if err:= zr.Reset ( fd ); err != nil { return nil,err }
      ret= tar.NewReader ( zr )
      return ret,nil
    }
  }

  // FALTARIA bz2
  
  return nil,err
  
} // end _TAR_Open


type TAR struct {}


func (self *TAR) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge un fitxer de tipus TAR" )
} // end GetImage


func (self *TAR) GetMetadata(file_name string) (string,error) {
  
  // Obri
  fd,err:= os.Open ( file_name )
  if err != nil { return "",err }
  defer fd.Close ()

  // Obté reader.
  r,err:= _TAR_Open ( fd )
  if err != nil { return "",err }

  // Inicialitza header
  md:= _TAR_Metadata{
    NFiles : 0,
  }
  
  // Compta el nombre de fitxers.
  _,err= r.Next ()
  for ; err == nil; _,err= r.Next () {
    md.NFiles++
  }
  if err != io.EOF { return "",err }

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *TAR) GetName() string {
  return "Fitxer d'emmagatzemament TAR (Tape ARchives)"
} // end GetName


func (self *TAR) GetShortName() string { return "TAR" }
func (self *TAR) IsImage() bool { return false }


func (self *TAR) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _TAR_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[TAR] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue
  
  // Nombre fitxers
  kv= &KeyValue{"Nº Fitxers",fmt.Sprint ( md.NFiles )}
  v= append(v,kv)
  
  return v
  
} // end ParseMetadata
