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
  "fmt"
  "image"
  
  "github.com/adriagipas/imgteka/view"
)




/*************/
/* CONSTANTS */
/*************/

// IDENTIFICADORS
// MOLT IMPORTANT!!! NO REASSIGNAR IDS. PUC TORNAR-LOS EN L'ORDRE QUE
// VULLGA, PERÒ NO PUC MODIFICAR ELS IDS.
const ID_IMAGE_PNG  = 0x100
const ID_IMAGE_JPEG = 0x101

const ID_ROM_GG     = 0x200
const ID_ROM_GBC    = 0x201
const ID_ROM_MD     = 0x202
const ID_ROM_NES    = 0x203

const ID_ARCH_ZIP   = 0x300

const ID_EXE_SFZ    = 0x400
const ID_EXE_ZBLORB = 0x401




/*********/
/* TIPUS */
/*********/

type FileType interface {

  // Torna la imatge del fitxer indicat (d'acord amb aquest tipus).
  GetImage(file_name string) (image.Image,error)
  
  // Aquest mètode te dos propòstis:
  //  1) Torna en un string amb les metadades.
  //  2) Comprovar que efectivament el fitxer és del tipus indicat.
  //
  // NOTA! El 'fd' no té perquè estar apuntant al principi del fitxer,
  // però es pot i es deu rebobinar.
  GetMetadata(file_name string) (string,error)

  // Torna el nom
  GetName() string
  
  // Un nom curt sense espais i en majúscules
  GetShortName() string

  // Indica si d'aquest tipus es pot obtindre una imatge
  IsImage() bool

  // Parseja el string que conté el json i afegeix els valors al
  // StringPairs.
  ParseMetadata(v []view.StringPair,meta_data string) []view.StringPair
  
}


type KeyValue struct {
  key,value string
}
func (self *KeyValue) GetKey() string { return self.key }
func (self *KeyValue) GetValue() string { return self.value }




/****************/
/* PART PRIVADA */
/****************/

var _IDS []int= []int{
  
  ID_IMAGE_PNG,
  ID_IMAGE_JPEG,

  ID_ROM_GBC,
  ID_ROM_GG,
  ID_ROM_MD,
  ID_ROM_NES,

  ID_EXE_SFZ,
  ID_EXE_ZBLORB,
  
  ID_ARCH_ZIP,
  
}

// Tipus globals
var _vPNG PNG= PNG{}
var _vJPEG JPEG= JPEG{}
var _vGBC GBC= GBC{}
var _vGG GG= GG{}
var _vMD MD= MD{}
var _vNES NES= NES{}
var _vSFZ SFZ= SFZ{}
var _vZBlorb ZBlorb= ZBlorb{}
var _vZIP ZIP= ZIP{}




/****************/
/* PART PÚBLICA */
/****************/

func Get( id int ) (FileType,error) {
  
  switch id {
    
  case ID_IMAGE_PNG:
    return &_vPNG,nil
  case ID_IMAGE_JPEG:
    return &_vJPEG,nil

  case ID_ROM_GBC:
    return &_vGBC,nil
  case ID_ROM_GG:
    return &_vGG,nil
  case ID_ROM_MD:
    return &_vMD,nil
  case ID_ROM_NES:
    return &_vNES,nil
    
  case ID_EXE_SFZ:
    return &_vSFZ,nil
  case ID_EXE_ZBLORB:
    return &_vZBlorb,nil

  case ID_ARCH_ZIP:
    return &_vZIP,nil
    
  default:
    return nil,fmt.Errorf ( "Tipus de fitxer desconegut:", id )
  }
  
} // end Get


func GetIDs() []int {
  return _IDS
} // end GetIDs
