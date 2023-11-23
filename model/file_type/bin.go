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
 *  bin.go - Fitxer binary genèric.
 */

package file_type

import (
  "errors"
  "image"
  
  "github.com/adriagipas/imgteka/view"
)

type BIN struct {
}


func (self *BIN) GetImage( file_name string) (image.Image,error) {
  return nil,errors.New (
    "No es pot interpretar com una imatge un fitxer binari genèric" )
} // end GetImage


func (self *BIN) GetMetadata(file_name string) (string,error) {
  return "",nil
} // end GetMetadata


func (self *BIN) GetName() string { return "Fitxer binari" }
func (self *BIN) GetShortName() string { return "BIN" }
func (self *BIN) IsImage() bool { return false }


func (self *BIN) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {
  return v
} // end ParseMetadata
