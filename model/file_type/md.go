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
 *  md.go - Tipus de fitxer ROM de Mega Drive.
 */

package file_type

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "log"
  "math"
  "os"
  "strings"
  
  "github.com/adriagipas/imgteka/view"
  "golang.org/x/text/encoding/japanese"
)




/****************/
/* PART PRIVADA */
/****************/

const (
  _MD_SRAMINFO_AVAILABLE   = 0x01
  _MD_SRAMINFO_ISFORBACKUP = 0x02
  _MD_SRAMINFO_ODDBYTES    = 0x04
  _MD_SRAMINFO_EVENBYTES   = 0x08
  _MD_SRAMINFO_ISEEPROM    = 0x10
)


type _MD_Metadata struct {

  Console      string
  FirmBuild    string
  DomName      string
  IntName      string
  TypeSnumber  string
  Checksum     uint16
  RealChecksum uint16
  IO           string
  Start        uint32
  End          uint32
  StartRam     uint32
  EndRam       uint32
  SramFlags    uint8
  StartSram    uint32
  EndSram      uint32
  Modem        string
  CCodes       string
}


func _MD_b2s( data []byte ) string {

  data= bytes.TrimSuffix ( data, []byte{0x00} )
  ret:= string(data)

  return strings.TrimSpace ( ret )
  
} // end _MD_b2s


func _MD_CalcChecksum( data []byte ) uint16 {

  var ret uint16= 0
  data= data[0x200:]
  for i:= 0; i < len(data); i+= 2 {
    ret+= (uint16(data[i])<<8) | uint16(data[i+1])
  }

  return ret
  
} // end _MD_CalcChecksum


func _MD_ReadHeader( header *_MD_Metadata, data []byte ) {

  // Nom de la consola
  header.Console= _MD_b2s ( data[0x100:0x100+16] )

  // Signatura i data
  header.FirmBuild= _MD_b2s ( data[0x110:0x110+16] )

  // Nom domèstic
  buf:= data[0x120:0x120+48]
  dec:= japanese.ShiftJIS.NewDecoder ()
  if aux,err:= dec.Bytes ( buf ); err != nil {
    header.DomName= _MD_b2s ( buf )
  } else {
    header.DomName= _MD_b2s ( aux )
  }

  // Nom internacional
  buf= data[0x150:0x150+48]
  if aux,err:= dec.Bytes ( buf ); err != nil {
    header.IntName= _MD_b2s ( buf )
  } else {
    header.IntName= _MD_b2s ( aux )
  }

  // Tipus de programa i número série
  header.TypeSnumber= _MD_b2s ( data[0x180:0x180+14] )

  // Checksum
  header.Checksum= (uint16(data[0x18e])<<8) | uint16(data[0x18f])
  header.RealChecksum= _MD_CalcChecksum ( data )

  // Support I/O
  header.IO= _MD_b2s ( data[0x190:0x190+16] )

  // Inici ROM
  header.Start= (uint32(data[0x1a0])<<24) | (uint32(data[0x1a1])<<16) |
    (uint32(data[0x1a2])<<8) | uint32(data[0x1a3])
    
  // Fi ROM
  header.End= (uint32(data[0x1a4])<<24) | (uint32(data[0x1a5])<<16) |
    (uint32(data[0x1a6])<<8) | uint32(data[0x1a7])

  // Inici RAM
  header.StartRam= (uint32(data[0x1a8])<<24) | (uint32(data[0x1a9])<<16) |
    (uint32(data[0x1aa])<<8) | uint32(data[0x1ab])

  // Fi RAM
  header.EndRam= (uint32(data[0x1ac])<<24) | (uint32(data[0x1ad])<<16) |
    (uint32(data[0x1ae])<<8) | uint32(data[0x1af])
  
  // SRAM
  header.SramFlags= 0x00
  tmp:= data[0x1b0:0x1b0+4]
  if tmp[0]=='R' && tmp[1]=='A' && (tmp[2]&0xa7)==0xa0 &&
    (tmp[3]==0x20 || tmp[3]==0x40) {

    // Flags
    header.SramFlags|= _MD_SRAMINFO_AVAILABLE
    if (tmp[2]&0x40) != 0 {
      header.SramFlags|= _MD_SRAMINFO_ISFORBACKUP
    }
    if tmp[3] == 0x40 {
      header.SramFlags|= _MD_SRAMINFO_ISEEPROM
    }
    switch (tmp[2]&0x18)>>3 {
    case 0:
      header.SramFlags|= _MD_SRAMINFO_ODDBYTES|_MD_SRAMINFO_EVENBYTES
    case 2:
      header.SramFlags|= _MD_SRAMINFO_EVENBYTES
    case 3:
      header.SramFlags|= _MD_SRAMINFO_ODDBYTES
    }

    // Inici
    header.StartSram= (uint32(data[0x1b4])<<24) | (uint32(data[0x1b5])<<16) |
      (uint32(data[0x1b6])<<8) | uint32(data[0x1b7])
    
    // Fi
    header.EndSram= (uint32(data[0x1b8])<<24) | (uint32(data[0x1b9])<<16) |
      (uint32(data[0x1ba])<<8) | uint32(data[0x1bb])
    
  } else {
    header.StartSram= 0
    header.EndSram= 0
  }

  // Suport modem
  header.Modem= _MD_b2s ( data[0x1bc:0x1bc+12] )
  
  // Codi paisos
  header.CCodes= _MD_b2s ( data[0x1f0:0x1f0+3] )
  
} // end _MD_ReadHeader


func _MD_DecodeRegion( code string ) string {

  var ret string
  
  // Nou estil
  if len(code) == 1 &&
    ((code[0] >= '0' && code[0] <= '9') ||
      (code[0]>='A' && code[0]<='F' && code[0] != 'E')) {
    
    // Converteix a byte
    var val byte
    if code[0] <= '9' {
      val= byte(code[0]-'0')
    } else {
      val= byte(code[0]-'A') + 0xa
    }

    // Afegeix regions
    ret= ""
    if (val&0x01)!=0 {
      if len(ret)>0 { ret+= ", " }
      ret+= "Japó, Corea del Sud, Taiwan"
    }
    if (val&0x04)!=0 {
      if len(ret)>0 { ret+= ", " }
      ret+= "Estats Units d'Amèrica, Brasil"
    }
    if (val&0x08)!=0 {
      if len(ret)>0 { ret+= ", " }
      ret+= "Europa, Hong Kong"
    }
    
  } else { // Prova estil antinc
    ret= ""
    for _,c:= range code {
      if c == 'J' {
        if len(ret)>0 { ret+= ", " }
        ret+= "Japó"
      } else if c == 'U' {
        if len(ret)>0 { ret+= ", " }
        ret+= "Amèrica"
      } else if c == 'E' {
        if len(ret)>0 { ret+= ", " }
        ret+= "Europa"
      } else {
        ret= code
        break
      }
    }
  }
  
  return ret
  
} // _MD_DecodeRegion




/****************/
/* PART PÚBLICA */
/****************/

type MD struct {
}


func (self *MD) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una ROM de Mega Drive" )
} // end GetImage


func (self *MD) GetMetadata(fd *os.File) (string,error) {
  
  // Rebobina
  if _,err:= fd.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  
  // Comprova grandària
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  size:= info.Size ()
  if ( size == 0 || (size%2) != 0 || size < 0x200 || size >= math.MaxInt ) {
    return "",errors.New (
      "La grandària del fitxer no és correspon amb el d'una"+
        " ROM de Mega Drive / Genesis" )
  }

  // Llig tota la ROM
  mem:= make([]byte,size)
  n,err:= fd.Read ( mem )
  if err != nil { return "",err }
  if n != int(size) {
    return "",errors.New ( "Error llegint les dades" )
  }

  // Llig capçalera
  md:= _MD_Metadata{}
  _MD_ReadHeader ( &md, mem )

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *MD) GetName() string { return "ROM de Mega Drive / Genesis" }
func (self *MD) GetShortName() string { return "MD" }
func (self *MD) IsImage() bool { return false }


func (self *MD) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _MD_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[MD] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue

  // Sistema
  kv= &KeyValue{"Sistema",md.Console}
  v= append(v,kv)
  
  // Copyright i data
  kv= &KeyValue{"Copyright i data",md.FirmBuild}
  v= append(v,kv)

  // Nom domèstic
  kv= &KeyValue{"Nom domèstic",md.DomName}
  v= append(v,kv)

  // Nom internacional
  kv= &KeyValue{"Nom internacional",md.IntName}
  v= append(v,kv)

  // Número serie
  kv= &KeyValue{"Codi de sèrie",md.TypeSnumber}
  v= append(v,kv)

  // Checksum
  if md.Checksum == md.RealChecksum {
    kv= &KeyValue{"Checksum",fmt.Sprintf ( "%04x (Sí)", md.Checksum )}
  } else {
    kv= &KeyValue{"Checksum",
      fmt.Sprintf ( "%04x (No != %04x)", md.Checksum, md.RealChecksum )}
  }
  v= append(v,kv)

  // I/O
  if md.IO != "" {
    kv= &KeyValue{"Dispositius suportats",md.IO}
    v= append(v,kv)
  }

  // Inici ROM
  kv= &KeyValue{"Inici ROM",fmt.Sprintf ( "%08x", md.Start )}
  v= append(v,kv)

  // Fi ROM
  kv= &KeyValue{"Fi ROM",fmt.Sprintf ( "%08x", md.End )}
  v= append(v,kv)

  // Inici RAM
  kv= &KeyValue{"Inici RAM",fmt.Sprintf ( "%08x", md.StartRam )}
  v= append(v,kv)

  // Fi RAM
  kv= &KeyValue{"Fi RAM",fmt.Sprintf ( "%08x", md.EndRam )}
  v= append(v,kv)

  // SRAM
  if (md.SramFlags&_MD_SRAMINFO_AVAILABLE) != 0 {

    // Flags
    var val string
    if (md.SramFlags&_MD_SRAMINFO_ISEEPROM) != 0 {
      val= "EEPROM"
    } else {
      val= "RAM estàtica"
    }

    // Addicionals
    if (md.SramFlags&_MD_SRAMINFO_ISFORBACKUP) != 0  {
      val+= ", sols backup"
    }
    if (md.SramFlags&_MD_SRAMINFO_ODDBYTES) != 0 {
      val+= ", bytes imparells"
    }
    if (md.SramFlags&_MD_SRAMINFO_EVENBYTES) != 0 {
      val+= ", bytes parells"
    }

    // Afegeix
    kv= &KeyValue{"Memòria addicional",val}
    v= append(v,kv)

    // Inici RAM
    kv= &KeyValue{"Inici Mem. Add.",fmt.Sprintf ( "%08x", md.StartSram )}
    v= append(v,kv)
    
    // Fi RAM
    kv= &KeyValue{"Fi Mem. Add.",fmt.Sprintf ( "%08x", md.EndSram )}
    v= append(v,kv)
    
  }

  // Modem
  if md.Modem != "" {
    kv= &KeyValue{"Mòdem",md.Modem}
    v= append(v,kv)
  }
  
  // Codi països
  if md.CCodes != "" {
    kv= &KeyValue{"Regió",_MD_DecodeRegion ( md.CCodes )}
    v= append(v,kv)
  }
  
  return v
  
} // end ParseMetadata
