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
 *  gg.go - Tipus de fitxer ROM de Game Gear.
 */

package file_type

import (
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




const _GG_BANK_SIZE= 16384

const _GG_HEADER= "TMR SEGA"

const (
  _GG_REGION_SMS_JAPAN        = 0
  _GG_REGION_SMS_EXPORT       = 1
  _GG_REGION_GG_JAPAN         = 2
  _GG_REGION_GG_EXPORT        = 3
  _GG_REGION_GG_INTERNATIONAL = 4
  _GG_REGION_UNK              = 5
)


type _GG_Metadata struct {

  Checksum    uint16
  ProductCode int
  Version     int
  Region      int
  RomSize     int
}


type GG struct {
}


func (self *GG) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una ROM de Game Gear" )
} // end GetImage


func (self *GG) GetMetadata(file_name string) (string,error) {

  // Obri
  fd,err:= os.Open ( file_name )
  if err != nil { return "",err }
  defer fd.Close ()
  
  // Comprova grandària
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  size:= info.Size ()
  if ( size == 0 || (size%_GG_BANK_SIZE) != 0 ) {
    return "",errors.New (
      "La grandària del fitxer no és correspon amb el d'una ROM de Game Gear" )
  }
  nbanks:= size/_GG_BANK_SIZE;

  // Llig banks on pot està la capçalera.
  read_nbanks:= 1
  if nbanks > 1 { read_nbanks= 2 }
  var data_mem [2*_GG_BANK_SIZE]byte
  mem:= data_mem[:read_nbanks*_GG_BANK_SIZE]
  n,err:= fd.Read ( mem )
  if err != nil { return "",err }
  if n != read_nbanks*_GG_BANK_SIZE {
    return "",errors.New ( "Error llegint les dades" )
  }

  // Localitza capçalera
  header_pos:= -1
  if _GG_HEADER == string(mem[0x1ff0:0x1ff0+8]) {
    header_pos= 0x1ff0
  } else if _GG_HEADER == string(mem[0x3ff0:0x3ff0+8]) {
    header_pos= 0x3ff0
  } else if nbanks > 1 && _GG_HEADER == string(mem[0x7ff0:0x7ff0+8]) {
    header_pos= 0x7ff0
  }

  // Si no hi ha capçalera torna cadena buida
  if header_pos == -1 { return "",nil }
  mem= mem[header_pos:]

  // Checksum
  checksum:= uint16(mem[0xa]) | (uint16(mem[0xb])<<8)

  // Product code
  product_code:= int(mem[0xc]&0xf) +
    int(mem[0xc]>>4)*10 +
    int(mem[0xd]&0xf)*100 +
    int(mem[0xd]>>4)*1000 +
    int(mem[0xe]>>4)*10000

  // Versió
  version:= int(mem[0xe]&0xf)
  
  // Region code
  var region int
  switch mem[0xf]>>4 {
  case 3:
    region= _GG_REGION_SMS_JAPAN
  case 4:
    region= _GG_REGION_SMS_EXPORT
  case 5:
    region= _GG_REGION_GG_JAPAN
  case 6:
    region= _GG_REGION_GG_EXPORT
  case 7:
    region= _GG_REGION_GG_INTERNATIONAL
  default:
    region= _GG_REGION_UNK
  }

  // Rom size
  var rom_size int
  switch mem[0xf]&0xf {
  case 0x0:
    rom_size= 256
  case 0x1:
    rom_size= 512
  case 0x2:
    rom_size= 1024
  case 0xa:
    rom_size= 8
  case 0xb:
    rom_size= 16
  case 0xc:
    rom_size= 32
  case 0xd:
    rom_size= 48
  case 0xe:
    rom_size= 64
  case 0xf:
    rom_size= 128
  default:
    rom_size= -1
  }

  // Metadades
  md:= _GG_Metadata{checksum, product_code, version, region, rom_size}

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *GG) GetName() string { return "ROM de Game Gear" }
func (self *GG) GetShortName() string { return "GG" }
func (self *GG) IsImage() bool { return false }


func (self *GG) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  if meta_data == "" { return v }
  
  // Parseja
  md:= _GG_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[GG] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Checksum
  kv:= &KeyValue{"Checksum",fmt.Sprintf ( "%04x", md.Checksum )}
  v= append(v,kv)

  // Product code
  kv= &KeyValue{"Codi",fmt.Sprintf ( "%d", md.ProductCode )}
  v= append(v,kv)

  // Version
  kv= &KeyValue{"Versió",fmt.Sprintf ( "%d", md.Version )}
  v= append(v,kv)

  // Region
  if md.Region != _GG_REGION_UNK {
    var text string
    switch md.Region {
    case _GG_REGION_SMS_JAPAN:
      text= "SMS Japó"
    case _GG_REGION_SMS_EXPORT:
      text= "SMS Exportació"
    case _GG_REGION_GG_JAPAN:
      text= "GG Japó"
    case _GG_REGION_GG_EXPORT:
      text= "GG Exportació"
    case _GG_REGION_GG_INTERNATIONAL:
      text= "GG Internacional"
    }
    kv= &KeyValue{"Regió",text}
    v= append(v,kv)
  }

  // Grandària
  if md.RomSize != -1 {
    text:= fmt.Sprintf ( "%d KB", md.RomSize )
    kv= &KeyValue{"Grandària (segons capçalera)",text}
    v= append(v,kv)
  }
  
  return v
  
} // end ParseMetadata
