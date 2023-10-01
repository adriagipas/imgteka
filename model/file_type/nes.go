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
 *  nes.go - Tipus de fitxer .NES.
 */

package file_type

import (
  //"encoding/json"
  "errors"
  "fmt"
  "image"
  //"log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




/****************/
/* PART PRIVADA */
/****************/

const (
  _NES_TYPE_ARCHAIC_INES = 0
  _NES_TYPE_INES_0_7     = 1
  _NES_TYPE_INES         = 2
  _NES_TYPE_INES_2_0     = 3
)


const (
  _NES_MIRRORING_HORIZONTAL  = 0
  _NES_MIRRORING_VERTICAL    = 1
  _NES_MIRRORING_FOUR_SCREEN = 2
)


const (
  _NES_CONSOLE_TYPE_UNK                    = 0
  _NES_CONSOLE_TYPE_REGULAR                = 1
  _NES_CONSOLE_TYPE_VS_SYSTEM              = 2
  _NES_CONSOLE_TYPE_PLAYCHOICE10           = 3
  _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_WDM  = 4
  _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_EPSM = 5
  _NES_CONSOLE_TYPE_VRT_VT01               = 6
  _NES_CONSOLE_TYPE_VRT_VT02               = 7
  _NES_CONSOLE_TYPE_VRT_VT03               = 8
  _NES_CONSOLE_TYPE_VRT_VT09               = 9
  _NES_CONSOLE_TYPE_VRT_VT32               = 10
  _NES_CONSOLE_TYPE_VRT_VT369              = 11
  _NES_CONSOLE_TYPE_UMC_UM6578             = 12
  _NES_CONSOLE_TYPE_FAMICOM_NETWORK_SYSTEM = 13
)


const (
  _NES_TV_SYSTEM_NTSC     = 0
  _NES_TV_SYSTEM_PAL      = 1
  _NES_TV_SYSTEM_MULTIPLE = 2
  _NES_TV_SYSTEM_DENDY    = 3
  _NES_TV_SYSTEM_UNK      = 4
)


type _NES_Metadata struct {
  
  PRG_Size              int64
  CHR_Size              int64 // 0 vol dir CHR RAM
  Mirroring             int
  Sram                  bool
  Trainer               bool
  Console               int
  Mapper                int // Com no paren de crear nous ho deixe en número
  Submapper             int // 0 vol dir no submapper
  PRG_RAM_Size          int // 0 vol dir que no s'especifica
  PRG_NVRAM_EEPROM_Size int // 0 vol dir que no hi ha o no s'especifica
  TV_System             int
  CHR_RAM_Size          int // 0 vol dir que no s'especifica o no hi ha
  CHR_NVRAM_Size        int // 0 vol dir que no hi ha o no s'especifica
  
}


func _NES_calc_size_nes2_0( code uint16, pag_size int64 ) int64 {

  var ret int64
  
  if ( (code&0xf00) == 0xf00 ) {
    E:= int64((code>>2)&0x3f)
    MM:= int64(code&0x3)
    ret= (int64(1)<<E)*(2*MM + 1)
  } else {
    ret= int64(uint64(code))*pag_size
  }

  return ret
  
} // _NES_calc_size_nes2_0


func _NES_set_console_type( header *_NES_Metadata, data []byte, ftype int ) {

  f7:= uint8(data[7])

  // iNES 2.0
  if ftype == _NES_TYPE_INES_2_0 {
    switch f7&0x3 {
    case 0:
      header.Console= _NES_CONSOLE_TYPE_REGULAR
    case 1:
      header.Console= _NES_CONSOLE_TYPE_VS_SYSTEM
    case 2:
      header.Console= _NES_CONSOLE_TYPE_PLAYCHOICE10
    case 3:
      f13:= uint8(data[13])
      switch f13&0xf {
      case 0x0:
        header.Console= _NES_CONSOLE_TYPE_REGULAR
      case 0x1:
        header.Console= _NES_CONSOLE_TYPE_VS_SYSTEM
      case 0x2:
        header.Console= _NES_CONSOLE_TYPE_PLAYCHOICE10
      case 0x3:
        header.Console= _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_WDM
      case 0x4:
        header.Console= _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_EPSM
      case 0x5:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT01
      case 0x6:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT02
      case 0x7:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT03
      case 0x8:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT09
      case 0x9:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT32
      case 0xa:
        header.Console= _NES_CONSOLE_TYPE_VRT_VT369
      case 0xb:
        header.Console= _NES_CONSOLE_TYPE_UMC_UM6578
      case 0xc:
        header.Console= _NES_CONSOLE_TYPE_FAMICOM_NETWORK_SYSTEM
      default:
        header.Console= _NES_CONSOLE_TYPE_UNK
      }
    }
    
    // iNES
  } else if ftype == _NES_TYPE_INES {
    if (f7&0x01) == 0x00 {
      header.Console= _NES_CONSOLE_TYPE_REGULAR
    } else {
      header.Console= _NES_CONSOLE_TYPE_VS_SYSTEM
    }

    // Other
  } else { 
    header.Console= _NES_CONSOLE_TYPE_UNK
  }
  
} // end _NES_set_console_type


func _NES_set_mapper( header *_NES_Metadata, data []byte, ftype int ) {

  f6:= uint8(data[6])
  // NES 2.0
  if ftype == _NES_TYPE_INES_2_0 {
    f7:= uint8(data[7])
    f8:= uint8(data[8])
    header.Mapper= int(uint32(
      uint16(f6>>4) | uint16(f7&0xf0) | (uint16(f8&0xf)<<8)))
    header.Submapper= int(f8>>4)
    
    // Archaic iNES
  } else if ftype == _NES_TYPE_ARCHAIC_INES {
    header.Mapper= int(f6>>4)
    header.Submapper= 0
    
    // Other
  } else {
    f7:= uint8(data[7])
    header.Mapper= int(uint32((f6>>4) | (f7&0xf0)))
    header.Submapper= 0
  }
  
} // end _NES_set_mapper


func _NES_set_prg_ram_size( header *_NES_Metadata, data []byte, ftype int ) {

  // NES 2.0
  if ftype == _NES_TYPE_INES_2_0 {
    f10:= uint8(data[10])
    var shift int
    // PRG-RAM
    shift= int(f10&0xf)
    if shift == 0 {
      header.PRG_RAM_Size= 0
    } else {
      header.PRG_RAM_Size= 64<<shift
    }
    // PRG-NVRAM/EEPROM
    shift= int(f10>>4)
    if shift == 0 {
      header.PRG_NVRAM_EEPROM_Size= 0
    } else {
      header.PRG_NVRAM_EEPROM_Size= 64<<shift
    }
    
    // iNES
  } else if ftype == _NES_TYPE_INES {
    f8:= uint8(data[8])
    if f8 == 0 {
      header.PRG_RAM_Size= 8*1024
    } else {
      header.PRG_RAM_Size= int(uint32(f8))*8*1024
    }
    header.PRG_NVRAM_EEPROM_Size= 0

    // Other
  } else {
    header.PRG_RAM_Size= 0
    header.PRG_NVRAM_EEPROM_Size= 0
  }
  
} // _NES_set_prg_ram_size


func _NES_set_tv_system( header *_NES_Metadata, data []byte, ftype int ) {

  // NES 2.0
  if ftype == _NES_TYPE_INES_2_0 {
    f12:= uint8(data[12])
    switch f12&0x3 {
    case 0:
      header.TV_System= _NES_TV_SYSTEM_NTSC
    case 1:
      header.TV_System= _NES_TV_SYSTEM_PAL
    case 2:
      header.TV_System= _NES_TV_SYSTEM_MULTIPLE
    case 3:
      header.TV_System= _NES_TV_SYSTEM_DENDY
    }

    // iNES
  } else if ftype == _NES_TYPE_INES {
    f9:= uint8(data[9])
    if (f9&0x1) == 0x0 {
      header.TV_System= _NES_TV_SYSTEM_NTSC
    } else {
      header.TV_System= _NES_TV_SYSTEM_PAL
    }
    
    // Other
  } else {
    header.TV_System= _NES_TV_SYSTEM_UNK
  }
  
} // end _NES_set_tv_system


func _NES_set_chr_ram_size( header *_NES_Metadata, data []byte, ftype int ) {

  // NES 2.0
  if ftype == _NES_TYPE_INES_2_0 {
    f11:= uint8(data[11])
    var shift int
    // CHR-RAM
    shift= int(f11&0xf)
    if shift == 0 {
      header.CHR_RAM_Size= 0
    } else {
      header.CHR_RAM_Size= 64<<shift
    }
    // CHR-NVRAM
    shift= int(f11>>4)
    if shift == 0 {
      header.CHR_NVRAM_Size= 0
    } else {
      header.CHR_NVRAM_Size= 64<<shift
    }
    
    // Other
  } else {
    header.CHR_RAM_Size= 0
    header.CHR_NVRAM_Size= 0
  }
  
} // _NES_set_chr_ram_size


func _NES_ReadHeader( header *_NES_Metadata, data []byte, size int64 ) error {

  // Comprova signatura.
  if data[0]!='N' || data[1]!='E' || data[2]!='S' || data[3]!=0x1a {
    return errors.New ( "No és un fitxer .NES" )
  }

  // PRG i CHROM
  PRG_lsb:= uint16(uint8(data[4]))
  CHR_lsb:= uint16(uint8(data[5]))

  // Flags 6
  f6:= uint8(data[6])
  // --> Mirroring
  if (f6&0x08) != 0 {
    header.Mirroring= _NES_MIRRORING_FOUR_SCREEN
  } else if (f6&0x01) != 0 {
    header.Mirroring= _NES_MIRRORING_VERTICAL
  } else {
    header.Mirroring= _NES_MIRRORING_HORIZONTAL
  }
  // --> SRAM
  header.Sram= (f6&0x02)!=0
  // --> Trainer
  header.Trainer= (f6&0x04)!=0
  
  // Versió iNES
  f7,f9:= uint8(data[7]),uint8(data[9])
  // --> Tipus
  tmp_PRG:= PRG_lsb | (uint16(f9&0xf)<<8)
  tmp_CHR:= CHR_lsb | (uint16(f9>>4)<<8)
  PRG_size:= _NES_calc_size_nes2_0( tmp_PRG, 16*1024 )
  CHR_size:= _NES_calc_size_nes2_0( tmp_CHR, 8*1024 )
  tmp_size:= 16 + PRG_size + CHR_size
  var ftype int
  if (f7&0x0C) == 0x08 && tmp_size <= size {
    ftype= _NES_TYPE_INES_2_0
  } else if (f7&0x0C) == 0x04 {
    ftype= _NES_TYPE_ARCHAIC_INES
  } else if (f7&0x0C) == 0x00 {
    all_zero:= true
    for i:= 12; i <= 15; i++ {
      if uint8(data[i]) != 0 {
        all_zero= false
      }
    }
    if all_zero {
      ftype= _NES_TYPE_INES
    } else {
      ftype= _NES_TYPE_INES_0_7
    }
  } else {
    ftype= _NES_TYPE_INES_0_7
  }
  // --> PRG i CHR sizes
  if ftype == _NES_TYPE_INES_2_0 {
    header.PRG_Size= PRG_size
    header.CHR_Size= CHR_size
  } else {
    header.PRG_Size= int64(PRG_lsb)*16*1024
    header.CHR_Size= int64(CHR_lsb)*8*1024
  }
  if header.PRG_Size == 0 {
    return errors.New ( "El fitxer no conté pàgines PRG" )
  }

  // Console type
  _NES_set_console_type ( header, data, ftype )

  // Mapper i Submapper
  _NES_set_mapper ( header, data, ftype )

  // PRG-RAM size
  _NES_set_prg_ram_size ( header, data, ftype )

  // TV System
  _NES_set_tv_system ( header, data, ftype )

  // CHR-RAM Size
  _NES_set_chr_ram_size ( header, data, ftype )
  
  fmt.Printf ( "F7:%X F9:%X PRG:%X CHR:%X PRG:%d CHR:%d FTYPE:%d\n", f7, f9, tmp_PRG, tmp_CHR, PRG_size, CHR_size, ftype )
  
  return nil
  
} // _NES_ReadHeader




/****************/
/* PART PÚBLICA */
/****************/


type NES struct {
}


func (self *NES) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge (iNES)" )
} // end GetImage



func (self *NES) GetMetadata(fd *os.File) (string,error) {

  // Rebobina
  if _,err:= fd.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }

  // Obté grandària.
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  size:= info.Size ()
  
  // Llig capçalera
  var mem [16]byte
  buf:= mem[:]
  n,err:= fd.Read ( buf )
  if err != nil { return "",err }
  if n != 16 {
    return "",errors.New ( "Error llegint la capçalera" )
  }

  // Llig capçalera
  md:= _NES_Metadata{}
  if err:= _NES_ReadHeader ( &md, buf, size ); err != nil {
    return "",err
  }
  
  fmt.Println ( "IEEEEE", size, buf, md )
  
  return "",errors.New ( "Fent iNES format" )
} // end GetMetadata


func (self *NES) GetName() string {
  return "ROM de Nintendo Entertainment System (iNES/NES 2.0)"
} // end GetName


func (self *NES) GetShortName() string { return "NES" }
func (self *NES) IsImage() bool { return false }


func (self *NES) ParseMetadata(

  v         []view.StringPair,
  meta_data string,
  
) []view.StringPair {
  
  return v
  
} // end ParseMetadata
