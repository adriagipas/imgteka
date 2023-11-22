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
 *  gbc.go - Tipus de fitxer ROM de Game Boy (Color).
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




/****************/
/* PART PRIVADA */
/****************/

const _GBC_BANK_SIZE= 16384

const (
  _GBC_GBCFLAG_ONLY_GBC = 0
  _GBC_GBCFLAG_GBC      = 1
  _GBC_GBCFLAG_GB       = 2
)

type _GBC_Metadata struct {

  Title               string // A vegades inclou també el codi del
                             // fabricant perquè és difícil de
                             // distingir.
  Manufacturer        string // A vegades està buit.
  CGBFlag             int
  OldLicense          uint8  // Byte antiga llicència (0x33 indica que
                             // s'utilitza NewLicense)
  NewLicense          string
  SGBFlag             bool   // Cert indica que suporta Super Game Boy
  Mapper              string
  RomSize             int    // Grandària de la ROM (segons capçalera) en
                             // trossos de 16KB. -1 indica que no es sap.
  RamSize             int    // Mesurat en KB. -1 indica que no es sap.
  JapaneseRom         bool   // Cert indica destinat a mercat japonés.
  Version             uint8  // Versió del joc
  Checksum            uint8  // Checksum capçalera
  GlobalChecksum      uint16
  RealChecksum        uint8
  RealGlobalChecksum  uint16
  NintendoLogo        bool   // Cert indica que la ROM el conté
  
}


func _GBC_IsUpper( b byte ) bool {
  if b >= 'A' && b <= 'Z' {
    return true
  } else {
    return false
  }
} // end _GBC_IsUpper


func _GBC_GetMapper( data []byte ) string {

  code:= data[0x147]
  switch code {
  case 0x00:
    return "ROM"
  case 0x01:
    return "MBC1"
  case 0x02:
    return "MBC1+RAM"
  case 0x03:
    return "MBC1+RAM+BATTERY"
  case 0x05:
    return "MBC2"
  case 0x06:
    return "MBC2+BATTERY"
  case 0x08:
    return "ROM+RAM"
  case 0x09:
    return "ROM+RAM+BATTERY"
  case 0x0B:
    return "MMM01"
  case 0x0C:
    return "MMM01+RAM"
  case 0x0D:
    return "MMM01+RAM+BATTERY"
  case 0x0F:
    return "MBC3+TIMER+BATTERY"
  case 0x10:
    return "MBC3+TIMER+RAM+BATTERY"
  case 0x11:
    return "MBC3"
  case 0x12:
    return "MBC3+RAM"
  case 0x13:
    return "MBC3+RAM+BATTERY"
  case 0x15:
    return "MBC4"
  case 0x16:
    return "MBC4+RAM"
  case 0x17:
    return "MBC4+RAM+BATTERY"
  case 0x19:
    return "MBC5"
  case 0x1A:
    return "MBC5+RAM"
  case 0x1B:
    return "MBC5+RAM+BATTERY"
  case 0x1C:
    return "MBC5+RUMBLE"
  case 0x1D:
    return "MBC5+RUMBLE+RAM"
  case 0x1E:
    return "MBC5+RUMBLE+RAM+BATTERY"
  case 0xFC:
    return "POCKET CAMERA"
  case 0xFD:
    return "BANDAI TAMA5"
  case 0xFE:
    return "HuC3"
  case 0xFF:
    return "HuC1+RAM+BATTERY"
  default:
    return fmt.Sprintf ( "UNK (%02X)", code )
  }
  
} // end _GBC_GetMapper


func _GBC_CalcChecksum( data []byte ) uint8 {

  var aux int= 0
  for i:= 0x134; i <= 0x14c; i++ {
    aux+= int(uint8(data[i])) + 1
  }

  return uint8((-aux)&0xFF)
  
} // end _GBC_CalcChecksum


func _GBC_CalcGlobalChecksum( data []byte ) uint16 {

  var aux int= 0
  for i:= 0; i < 0x14e; i++ {
    aux+= int(uint8(data[i]))
  }
  for i:= 0x14e+2; i < len(data); i++ {
    aux+= int(uint8(data[i]))
  }

  return uint16(aux&0xFFFF)
  
} // end _GBC_CalcGlobalChecksum


var _GBC_LOGO [48]uint8= [48]uint8{
  0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
  0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
  0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
  0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
  0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
  0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

func _GBC_CheckNintendoLogo( data []byte ) bool {

  p:= data[0x104:0x104+24]
  for i:= 0; i < 24; i++ {
    if _GBC_LOGO[i] != p[i] {
      return false
    }
  }
  
  return true
  
} // end _GBC_CheckNintendoLogo


func _GBC_ReadHeader( header *_GBC_Metadata, data []byte ) {

  // Flag CGB, títol i codi fabricant.
  // Açò és un poc heurístic.
  aux:= uint8(data[0x143])
  p:= data[0x134:]
  if (aux&0x80) != 0 {
    // --> CGB Flag
    if (aux&0xC0) != 0 {
      header.CGBFlag= _GBC_GBCFLAG_ONLY_GBC
    } else {
      header.CGBFlag= _GBC_GBCFLAG_GBC
    }
    // --> Title
    var i int
    for i= 0; i < 15 && (_GBC_IsUpper ( p[i] ) || p[i]==' '); i++ {}
    header.Title= string(p[:i])
    // --> Manufacturer
    if i <= 11 || i == 15 {
      p= data[0x13f:]
      for i= 0; i < 4 && _GBC_IsUpper ( p[i] ); i++ {}
      if i != 4 {
        header.Manufacturer= ""
      } else {
        header.Manufacturer= string(p[:4])
        if len(header.Manufacturer) < len(header.Title) {
          tmp_len:= len(header.Title)-len(header.Manufacturer)
          tmp:= header.Title[tmp_len:]
          if tmp == header.Manufacturer {
            header.Title= header.Title[:tmp_len]
          }
        }
      }
    } else {
      header.Manufacturer= ""
    }
    
  } else { // Assumiré que tots els cartutxos de GB són antics.
    // --> CGB Flag
    header.CGBFlag= _GBC_GBCFLAG_GB
    // --> Manufacturer
    header.Manufacturer= ""
    // --> Title
    var i int
    for i= 0; i < 16 && _GBC_IsUpper ( p[i] ); i++ {}
    header.Title= string(p[:i])
  }

  // Codi de la llicència
  header.OldLicense= uint8(data[0x14b])
  if header.OldLicense == 0x33 {
    header.NewLicense= string(data[0x144:0x144+2])
  } else {
    header.NewLicense= ""
  }

  // Flag SGB
  header.SGBFlag= data[0x146]==0x03

  // Mapper
  header.Mapper= _GBC_GetMapper ( data )

  // Rom size
  aux= uint8(data[0x148])
  if aux < 0x8 {
    header.RomSize= 2<<aux
  } else {
    switch aux {
    case 0x52:
      header.RomSize= 72
    case 0x53:
      header.RomSize= 80
    case 0x54:
      header.RomSize= 96
    default:
      header.RomSize= -1
    }
  }

  // Ram size
  aux= uint8(data[0x149])
  switch aux {
  case 0x00:
    header.RamSize= 0
  case 0x01:
    header.RamSize= 2
  case 0x02:
    header.RamSize= 8
  case 0x03:
    header.RamSize= 32
  default:
    header.RamSize= -1
  }

  // Japanese
  header.JapaneseRom= data[0x14a]==0x00

  // Version
  header.Version= uint8(data[0x14c])

  // Header Checksum
  header.Checksum= uint8(data[0x14d])
  header.RealChecksum= _GBC_CalcChecksum ( data )

  // Global checksum
  header.GlobalChecksum= (uint16(data[0x14e])<<8) | uint16(data[0x14f])
  header.RealGlobalChecksum= _GBC_CalcGlobalChecksum ( data )

  // Nintendo logo
  header.NintendoLogo= _GBC_CheckNintendoLogo ( data )
  
} // end _GBC_ReadHeader




/****************/
/* PART PÚBLICA */
/****************/

type GBC struct {
}


func (self *GBC) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una ROM de Game Boy (Color)" )
} // end GetImage


func (self *GBC) GetMetadata(file_name string) (string,error) {

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
  if ( size == 0 || (size%_GBC_BANK_SIZE) != 0 ) {
    return "",errors.New (
      "La grandària del fitxer no és correspon amb el d'una"+
        " ROM de Game Boy (Color)" )
  }
  nbanks:= int(size/_GBC_BANK_SIZE)
  
  // Llig el bank de la capçalera
  mem:= make([]byte,_GBC_BANK_SIZE*nbanks)
  n,err:= fd.Read ( mem )
  if err != nil { return "",err }
  if n != _GBC_BANK_SIZE*nbanks {
    return "",errors.New ( "Error llegint les dades" )
  }

  // Llig capçalera
  md:= _GBC_Metadata{}
  _GBC_ReadHeader ( &md, mem )
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *GBC) GetName() string { return "ROM de Game Boy (Color)" }
func (self *GBC) GetShortName() string { return "GBC" }
func (self *GBC) IsImage() bool { return false }


func (self *GBC) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _GBC_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[GBC] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue
  
  // Títol
  kv= &KeyValue{"Títol",md.Title}
  v= append(v,kv)

  // Fabricant
  if md.Manufacturer != "" {
    kv= &KeyValue{"Fabricant",md.Manufacturer}
    v= append(v,kv)
  }

  // Compatibilitat
  var comp string
  switch md.CGBFlag {
  case _GBC_GBCFLAG_ONLY_GBC:
    comp= "GBC"
  case _GBC_GBCFLAG_GBC:
    comp= "GBC/GB"
  default:
    comp= "GB"
  }
  kv= &KeyValue{"Compatibilitat",comp}
  v= append(v,kv)

  // Llicència
  var license string
  if md.OldLicense != 0x33 {
    license= fmt.Sprintf ( "%02x", md.OldLicense )
  } else {
    license= md.NewLicense
  }
  kv= &KeyValue{"Llicència",license}
  v= append(v,kv)

  // Suport SGB
  if md.SGBFlag {
    kv= &KeyValue{"Suport SGB","Sí"}
  } else {
    kv= &KeyValue{"Suport SGB","No"}
  }
  v= append(v,kv)

  // Mapper
  kv= &KeyValue{"Mapper",md.Mapper}
  v= append(v,kv)

  // Grandària ROM
  if md.RomSize != -1 {
    text:= fmt.Sprintf ( "%d KB", md.RomSize*16 )
    kv= &KeyValue{"Grandària (segons capçalera)",text}
    v= append(v,kv)
  }

  // Grandària SRAM
  if md.RamSize != -1 {
    text:= fmt.Sprintf ( "%d KB", md.RamSize )
    kv= &KeyValue{"Grandària SRAM",text}
    v= append(v,kv)
  }

  // ROM Japonesa
  if md.JapaneseRom {
    kv= &KeyValue{"ROM japonesa","Sí"}
  } else {
    kv= &KeyValue{"ROM japonesa","No"}
  }
  v= append(v,kv)

  // Checksum
  if md.Checksum == md.RealChecksum {
    kv= &KeyValue{"Checksum",fmt.Sprintf ( "%02x (Sí)", md.Checksum )}
  } else {
    kv= &KeyValue{"Checksum",
      fmt.Sprintf ( "%02x (No != %02x)", md.Checksum, md.RealChecksum )}
  }
  v= append(v,kv)

  // Checksum
  if md.GlobalChecksum == md.RealGlobalChecksum {
    kv= &KeyValue{"Checksum global",
      fmt.Sprintf ( "%04x (Sí)", md.GlobalChecksum )}
  } else {
    kv= &KeyValue{"Checksum global",
      fmt.Sprintf ( "%04x (No != %04x)",
        md.GlobalChecksum, md.RealGlobalChecksum )}
  }
  v= append(v,kv)
  
  // Logo nintendo
  if md.NintendoLogo {
    kv= &KeyValue{"Logo Nintendo","Sí"}
  } else {
    kv= &KeyValue{"Logo Nintendo","No"}
  }
  v= append(v,kv)
  
  return v
  
} // end ParseMetadata
