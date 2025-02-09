/*
 * Copyright 2025 Adrià Giménez Pastor.
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
 *  iso.go - Tipus de fitxer imatge generica CD (ISO 9660).
 */

package file_type

import (
  "encoding/json"
  "fmt"
  "image"
  "log"
  
  "github.com/adriagipas/imgcp/cdread"
  "github.com/adriagipas/imgteka/view"
)




/****************/
/* PART PRIVADA */
/****************/

type _CDISO_Metadata struct {

  Cd     _CD_Metadata // Metadades a nivell d'imatge de CD
  Iso    _ISO_Metadata // Metadades a nivell d'ISO
  
}




/****************/
/* PART PÚBLICA */
/****************/

type ISO struct {
}


func (self *ISO) GetImage( file_name string ) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una imatge de CD ISO-9660" )
} // end GetImage


func (self *ISO) GetMetadata(file_name string) (string,error) {

  // Intenta obrir el CD
  cd,err:= cdread.Open ( file_name )
  if err != nil { return "",err }

  // Inicialitza metadades.
  md:= _CDISO_Metadata{}
  md.Cd.Init ( cd )

  // Llig track ISO.
  iso,err:= cdread.ReadISO ( cd, 0, 0 )
  if err != nil { return "",err }
  md.Iso.Init ( iso )

  // Converteix  a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *ISO) GetName() string { return "CD ISO 9660" }
func (self *ISO) GetShortName() string { return "ISO" }
func (self *ISO) IsImage() bool { return false }


func (self *ISO) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _CDISO_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[ISO] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // CD i Iso
  v= md.Cd.Parse ( v )
  v= md.Iso.Parse ( v )
  
  return v
  
} // end ParseMetadata
