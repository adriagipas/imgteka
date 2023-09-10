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
 *  zip.go - Tipus de fitxer ZIP.
 */

package file_type

import (
  "archive/zip"
  "encoding/json"
  "fmt"
  "image"
  "log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




type _ZIP_Metadata struct {

  Comment string
  NFiles  int
  
}


type ZIP struct {}


func (self *ZIP) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge un fitxer de tipus ZIP" )
} // end GetImage


func (self *ZIP) GetMetadata(fd *os.File) (string,error) {

  // Rebobina
  if _,err:= fd.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }

  // Crea ZIP reader
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  reader,err:= zip.NewReader ( fd, info.Size () )
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut crear un lector de fitxers ZIP: %s",
      err )
  }

  // Metadades
  md:= _ZIP_Metadata{
    Comment : reader.Comment,
    NFiles  : len(reader.File),
  }

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *ZIP) GetName() string {
  return "Fitxer d'emmagatzemament ZIP"
} // end GetName


func (self *ZIP) GetShortName() string { return "ZIP" }
func (self *ZIP) IsImage() bool { return false }


func (self *ZIP) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _ZIP_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[ZIP] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue
  
  // Comentari
  if md.Comment != "" {
    kv= &KeyValue{"Comentari",md.Comment}
    v= append(v,kv)
  }

  // Nombre fitxers
  kv= &KeyValue{"Nº Fitxers",fmt.Sprint ( md.NFiles )}
  v= append(v,kv)
  
  return v
  
} // end ParseMetadata
