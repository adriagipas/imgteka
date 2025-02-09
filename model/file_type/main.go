/*
 * Copyright 2023-2025 Adrià Giménez Pastor.
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
const ID_BIN        = 0x000

const ID_IMAGE_PNG  = 0x100
const ID_IMAGE_JPEG = 0x101

const ID_ROM_GG     = 0x200
const ID_ROM_GBC    = 0x201
const ID_ROM_MD     = 0x202
const ID_ROM_NES    = 0x203
const ID_ROM_NDS    = 0x204
const ID_ROM_3DS    = 0x205

const ID_ARCH_ZIP   = 0x300
const ID_ARCH_TAR   = 0x301

const ID_EXE_SFZ    = 0x400
const ID_EXE_ZBLORB = 0x401
const ID_EXE_CXI    = 0x402

const ID_AUX_CD_PS1 = 0x500

const ID_CD_PS1     = 0x600
const ID_UMD_PSP    = 0x601
const ID_DVD_PS2    = 0x602
const ID_CD_ISO     = 0x603

const ID_FLP_FAT12  = 0x700

const ID_DOC_PDF    = 0x800




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

  ID_DOC_PDF,
  
  ID_ROM_GBC,
  ID_ROM_GG,
  ID_ROM_MD,
  ID_ROM_NES,
  ID_ROM_NDS,
  ID_ROM_3DS,

  ID_EXE_CXI,
  ID_EXE_SFZ,
  ID_EXE_ZBLORB,

  ID_CD_PS1,
  ID_AUX_CD_PS1,
  ID_DVD_PS2,
  ID_UMD_PSP,
  ID_CD_ISO,
  
  ID_ARCH_ZIP,
  ID_ARCH_TAR,

  ID_FLP_FAT12,
  
  ID_BIN,
  
}

// Tipus globals
var _vPNG PNG= PNG{}
var _vJPEG JPEG= JPEG{}
var _vPDF PDF= PDF{}
var _vGBC GBC= GBC{}
var _vGG GG= GG{}
var _vMD MD= MD{}
var _vNES NES= NES{}
var _vNDS NDS= NDS{}
var _v3DS N3DS= N3DS{}
var _vCXI CXI= CXI{}
var _vSFZ SFZ= SFZ{}
var _vZBlorb ZBlorb= ZBlorb{}
var _vPS1 PS1= PS1{}
var _vPS1_Aux PS1_Aux= PS1_Aux{}
var _vPS2 PS2= PS2{}
var _vPSP PSP= PSP{}
var _vISO ISO= ISO{}
var _vZIP ZIP= ZIP{}
var _vTAR TAR= TAR{}
var _vFAT12 FAT12= FAT12{}
var _vBIN BIN= BIN{}




/****************/
/* PART PÚBLICA */
/****************/

func Get( id int ) (FileType,error) {
  
  switch id {
    
  case ID_IMAGE_PNG:
    return &_vPNG,nil
  case ID_IMAGE_JPEG:
    return &_vJPEG,nil

  case ID_DOC_PDF:
    return &_vPDF,nil
    
  case ID_ROM_3DS:
    return &_v3DS,nil
  case ID_ROM_GBC:
    return &_vGBC,nil
  case ID_ROM_GG:
    return &_vGG,nil
  case ID_ROM_MD:
    return &_vMD,nil
  case ID_ROM_NES:
    return &_vNES,nil
  case ID_ROM_NDS:
    return &_vNDS,nil
    
  case ID_EXE_CXI:
    return &_vCXI,nil
  case ID_EXE_SFZ:
    return &_vSFZ,nil
  case ID_EXE_ZBLORB:
    return &_vZBlorb,nil

  case ID_CD_PS1:
    return &_vPS1,nil
  case ID_AUX_CD_PS1:
    return &_vPS1_Aux,nil
  case ID_DVD_PS2:
    return &_vPS2,nil
  case ID_UMD_PSP:
    return &_vPSP,nil
  case ID_CD_ISO:
    return &_vISO,nil
    
  case ID_ARCH_ZIP:
    return &_vZIP,nil
  case ID_ARCH_TAR:
    return &_vTAR,nil

  case ID_FLP_FAT12:
    return &_vFAT12,nil
    
  case ID_BIN:
    return &_vBIN,nil
    
  default:
    return nil,fmt.Errorf ( "Tipus de fitxer desconegut:", id )
  }
  
} // end Get


func GetIDs() []int {
  return _IDS
} // end GetIDs
