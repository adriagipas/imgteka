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
 *  file_type.go - Tipus de fitxers.
 */

package file_type

import (
  "log"
)




/*************/
/* CONSTANTS */
/*************/

// IDENTIFICADORS
const ID_IMAGE_PNG = 0x100




/****************/
/* PART PRIVADA */
/****************/

var _IDS []int= []int{
  
  ID_IMAGE_PNG,
  
}




/****************/
/* PART PÚBLICA */
/****************/

func GetIDs() []int {
  return _IDS
} // end GetIDs


func GetName( id int ) string {
  
  switch id {
    
  case ID_IMAGE_PNG:
    return "Imatge PNG"
    
  default:
    log.Fatal ( "Tipus de fitxer desconegut:", id )
    return ""
  }
  
} // end GetName
