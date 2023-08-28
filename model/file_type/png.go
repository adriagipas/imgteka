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
 *  file_type.go - Tipus de fitxer PNG.
 */

package file_type

import (
  "fmt"
  "image/png"
  "os"
)



type PNG struct {
}


func (self *PNG) GetMetadata(fd *os.File) (string,error) {

  // Rebobina
  if _,err:= fd.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }

  // Llig l'imatge
  img,err:= png.Decode ( fd )
  if err != nil {
    return "",fmt.Errorf ( "El fitxer no és de tipus PNG: %s", err )
  }
  
  fmt.Println ("TODO PNG.GetMetadata!!!!",img.Bounds() )
  
  return "",nil
  
} // end GetMetadata
