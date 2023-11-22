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
 *  png.go - Tipus de fitxer PNG.
 */

package file_type

import (
  "encoding/json"
  "fmt"
  "image"
  "image/png"
  "log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




type _PNG_Metadata struct {
  Width  int
  Height int
}


type PNG struct {
}


func (self *PNG) GetImage( file_name string) (image.Image,error) {

  f,err:= os.Open ( file_name )
  if err != nil { return nil,err }
  defer f.Close ()

  return png.Decode ( f )
  
} // end GetImage


func (self *PNG) GetMetadata(file_name string) (string,error) {

  // Obri
  fd,err:= os.Open ( file_name )
  if err != nil { return "",err }
  defer fd.Close ()

  // Llig l'imatge
  img,err:= png.Decode ( fd )
  if err != nil {
    return "",fmt.Errorf ( "El fitxer no és de tipus PNG: %s", err )
  }

  // Metadades
  bounds:= img.Bounds ()
  rmin,rmax:= bounds.Min,bounds.Max
  md:= _PNG_Metadata{rmax.X - rmin.X, rmax.Y - rmin.Y}

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *PNG) GetName() string { return "Imatge PNG" }
func (self *PNG) GetShortName() string { return "PNG" }
func (self *PNG) IsImage() bool { return true }


func (self *PNG) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {
  
  // Parseja
  md:= _PNG_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[PNG] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Metadada
  kv:= KeyValue{"Dimensions",fmt.Sprintf ( "%d x %d", md.Width, md.Height )}
  v= append(v,&kv)

  return v
  
} // end ParseMetadata
