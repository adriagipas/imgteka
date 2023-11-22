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
  "crypto/md5"
  "crypto/sha1"
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "io"
  "log"
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
  Submapper             int // 0 vol dir no submapper (-1 no especificat)
  PRG_RAM_Size          int // 0 vol dir que no s'especifica
  PRG_NVRAM_EEPROM_Size int // 0 vol dir que no hi ha o no s'especifica
  TV_System             int
  CHR_RAM_Size          int // 0 vol dir que no s'especifica o no hi ha
  CHR_NVRAM_Size        int // 0 vol dir que no hi ha o no s'especifica
  VS_PPU_Type           uint8 // Sols aplica si és _NES_CONSOLE_TYPE_VS_SYSTEM
  VS_CopyProtection     uint8 // Sols aplica si és _NES_CONSOLE_TYPE_VS_SYSTEM
  NumMiscellaneous      int // 0 vol dir que no en té o no està definit
  DefaultExpansionDev   uint8 // 0 vol dir no especificat
  RealMD5               string // Sense la capçalera
  RealSHA1              string // Sense la capçalera
  
}


func (self *_NES_Metadata) getDefaultExpansionDevice() string {

  switch self.DefaultExpansionDev {
  case 0x01:
    return "Controladors estàndards"
  case 0x02:
    return "NES Four Score"
  case 0x03:
    return "Famicom Four Players Adapter"
  case 0x04:
    return "Vs. System (1P via $4016)"
  case 0x05:
    return "Vs. System (1P via $4017)"
  case 0x07:
    return "Vs. Zapper"
  case 0x08:
    return "Zapper ($4017)"
  case 0x09:
    return "Dos Zapper"
  case 0x0a:
    return "Bandai Hyper Shot Lightgun"
  case 0x0b:
    return "Power Pad Side A"
  case 0x0c:
    return "Power Pad Side B"
  case 0x0d:
    return "Family Trainer Side A"
  case 0x0e:
    return "Family Trainer Side B"
  case 0x0f:
    return "Arkanoid Vaus Controller (NES)"
  case 0x10:
    return "Arkanoid Vaus Controller (Famicom)"
  case 0x11:
    return "2 Controladors Vaus + Famicom Data Recorder"
  case 0x12:
    return "Konami Hyper Shot"
  case 0x13:
    return "Coconuts Pachinko"
  case 0x14:
    return "Exciting Boxing Punching Bag"
  case 0x15:
    return "Jissen Mahjong Controller"
  case 0x16:
    return "Party Tap"
  case 0x17:
    return "Oeka Kids Tablet"
  case 0x18:
    return "Sunsoft Barcode Battler"
  case 0x19:
    return "Miracle Piano Keyboard"
  case 0x1a:
    return "Pokkun Moguraa"
  case 0x1b:
    return "Top Rider (Inflatable Bicycle)"
  case 0x1c:
    return "Double-Fisted"
  case 0x1d:
    return "Famicom 3D System"
  case 0x1e:
    return "Doremikko Keyboard"
  case 0x1f:
    return "R.O.B. Gyro Set"
  case 0x20:
    return "Famicom Data Recorder"
  case 0x21:
    return "ASCII Turbo File"
  case 0x22:
    return "IGS Storage Battle Box"
  case 0x23:
    return "Family BASIC Keyboard + Famicom Data Recorder"
  case 0x24:
    return "Dongda PEC-586 Keyboard"
  case 0x25:
    return "Bit Corp. Bit-79 Keyboard"
  case 0x26:
    return "Subor Keyboard"
  case 0x27:
    return "Subor Keyboard + mouse (3x8-bit protocol)"
  case 0x28:
    return "Subor Keyboard + mouse (24-bit protocol via $4016)"
  case 0x29:
    return "SNES Mouse ($4017.d0)"
  case 0x2a:
    return "Multicart"
  case 0x2b:
    return "2 Controladors SNES"
  case 0x2c:
    return "RacerMate Bicycle"
  case 0x2d:
    return "U-Force"
  case 0x2e:
    return "R.O.B. Stack-Up"
  case 0x2f:
    return "City Patrolman Lightgun"
  case 0x30:
    return "Sharp C1 Cassette Interface"
  case 0x31:
    return "Controlador estàndard amb botons invertits"
  case 0x32:
    return "Excalibor Sudoku Pad"
  case 0x33:
    return "ABL Pinball"
  case 0x34:
    return "Golden Nugget Casino extra buttons"
  case 0x36:
    return "Subor Keyboard + mouse (24-bit protocol via $4017)"
  case 0x37:
    return "Port test controller"
  case 0x38:
    return "Bandai Multi Game Player Gamepad buttons"
  case 0x39:
    return "Venom TV Dance Mat"
  case 0x3a:
    return "LG TV Remote Control"
  default:
    return fmt.Sprintf ( "UNK (%02X)", self.DefaultExpansionDev )
  }
  
} // end getDefaultExpansionDevice


func (self *_NES_Metadata) getConsole() string {

  switch self.Console {
  case _NES_CONSOLE_TYPE_REGULAR:
    return "Nintendo Entertainment System/Family Computer"
  case _NES_CONSOLE_TYPE_VS_SYSTEM:
    return "Nintendo Vs. System"
  case _NES_CONSOLE_TYPE_PLAYCHOICE10:
    return "Nintendo Playchoice 10"
  case _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_WDM:
    return "Famiclone (CPU amb suport decimal)"
  case _NES_CONSOLE_TYPE_REGULAR_FAMICLONE_EPSM:
    return "NES/Famicom (EPSM)"
  case _NES_CONSOLE_TYPE_VRT_VT01:
    return "V.R. Technology VT01"
  case _NES_CONSOLE_TYPE_VRT_VT02:
    return "V.R. Technology VT02"
  case _NES_CONSOLE_TYPE_VRT_VT03:
    return "V.R. Technology VT03"
  case _NES_CONSOLE_TYPE_VRT_VT09:
    return "V.R. Technology VT09"
  case _NES_CONSOLE_TYPE_VRT_VT32:
    return "V.R. Technology VT32"
  case _NES_CONSOLE_TYPE_VRT_VT369:
    return "V.R. Technology VT369"
  case _NES_CONSOLE_TYPE_UMC_UM6578:
    return "UMC UM6578"
  case _NES_CONSOLE_TYPE_FAMICOM_NETWORK_SYSTEM:
    return "Famicom Network System"
  default:
    return "Desconeguda"
  }
  
} // end getConsole


func (self *_NES_Metadata) getMapper() string {

  var ret string

  switch self.Mapper {
  case 0: 
    ret= "NROM"
  case 1:
    ret= "Nintendo MMC1"
    switch self.Submapper {
    case 5:
      ret+= " (Fixed PRG)"
    case 6:
      ret+= " (2ME)"
    }
  case 2:
    ret= "UxROM"
    switch self.Submapper {
    case 1:
      ret+= " (No bus conflicts)"
    case 2:
      ret+= " (AND bus conflicts)"
    }
  case 3:
    ret= "CNROM"
    switch self.Submapper {
    case 1:
      ret+= " (No bus conflicts)"
    case 2:
      ret+= " (AND bus conflicts)"
    }
  case 4:
    switch self.Submapper {
    case 1:
      ret= "Nintendo MMC6"
    case 3:
      ret= "MC-ACC"
    case 4:
      ret= "NEC MMC3"
    case 5:
      ret= "T9552"
    default:
      ret= "Nintendo MMC3"
    }
  case 5:
    ret= "Nintendo MMC5"
  case 7:
    ret= "AxROM"
    switch self.Submapper {
    case 1:
      ret+= " (No bus conflicts)"
    case 2:
      ret+= " (AND bus conflicts)"
    }
  case 9:
    ret= "Nintendo MMC2"
  case 10:
    ret= "Nintendo MMC4"
  case 11:
    ret= "Color Dreams"
  case 13:
    ret= "CPROM"
  case 16:
    switch self.Submapper {
    case 4:
      ret= "Bandai FCG-1/2"
    case 5:
      ret= "Bandai LZ93D50"
    default:
      ret= "Bandai FCG"
    }
  case 21:
    switch self.Submapper {
    case 0:
      ret= "Konami VRC4"
    case 1:
      ret= "Konami VRC4a"
    case 2:
      ret= "Konami VRC4c"
    default:
      ret= "Konami VRC2/VRC4"
    }
  case 22:
    switch self.Submapper {
    case 0:
      ret= "Konami VRC2a"
    default:
      ret= "Konami VRC2/VRC4"
    }
  case 23:
    switch self.Submapper {
    case 0:
      ret= "Konami VRC4"
    case 1:
      ret= "Konami VRC4f"
    case 2:
      ret= "Konami VRC4e"
    case 3:
      ret= "Konami VRC4b"
    default:
      ret= "Konami VRC2/VRC4"
    }
  case 24:
    ret= "Konami VRC6a"
  case 25:
    switch self.Submapper {
    case 0:
      ret= "Konami VRC4"
    case 1:
      ret= "Konami VRC4b"
    case 2:
      ret= "Konami VRC4d"
    case 3:
      ret= "Konami VRC4c"
    default:
      ret= "Konami VRC2/VRC4"
    }
  case 26:
    ret= "Konami VRC6b"
  case 28:
    ret= "Action 53"
  case 30:
    ret= "UNROM 512"
  case 32:
    ret= "Irem G101"
    if self.Submapper == 1 {
      ret+= " (Major League)"
    }
  case 33:
    ret= "Taito TC0190"
  case 34:
    switch self.Submapper {
    case 1:
      ret= "NINA-001"
    case 2:
      ret= "BNROM"
    default:
      ret= "BNROM / NINA-001"
    }
  case 48:
    ret= "Taito TC0690"
  case 61:
    ret= "NTDEC 0324 PCB"
  case 64:
    ret= "Tengen RAMBO-1"
  case 65:
    ret= "Irem H3001"
  case 66:
    ret= "GxROM"
  case 67:
    ret= "Sunsoft-3"
  case 68:
    if self.Submapper == 1 {
      ret= "Sunsoft Dual Cartridge"
    } else {
      ret= "Sunsoft-4"
    }
  case 69:
    ret= "Sunsoft FME-7"
  case 71:
    ret= "Codemasters"
    if self.Submapper == 1 {
      ret+= " (Fire Hawk)"
    }
  case 72:
    ret= "Jaleco JF-17"
  case 73:
    ret= "Konami VRC3"
  case 74:
    ret= "43-393/860908C"
  case 75:
    ret= "Konami VRC1"
  case 76:
    ret= "NAMCOT-3446"
  case 79:
    ret= "NINA-03/NINA-06"
  case 80:
    ret= "Taito X1-005"
  case 82:
    ret= "Taito X1-017"
  case 85:
    switch self.Submapper {
    case 1:
      ret= "Konami VRC7b"
    case 2:
      ret= "Konami VRC7a"
    default:
      ret= "Konami VRC7"
    }
  case 86:
    ret= "Jaleco JF-13"
  case 93:
    ret= "Sunsoft-2 IC"
  case 94:
    ret= "HVC-UN1ROM"
  case 95:
    ret= "NAMCOT-3425"
  case 97:
    ret= "Irem TAM-S1"
  case 105:
    ret= "NES-EVENT"
  case 113:
    ret= "HES NTD-8"
  case 118:
    ret= "TxSROM"
  case 119:
    ret= "TQROM"
  case 140:
    ret= "Jaleco JF-11/JF-14"
  case 154:
    ret= "NAMCOT-3453"
  case 158:
    ret= "Tengen 800037"
  case 159:
    ret= "Bandai EPROM (24C01)"
  case 166:
    ret= "SUBOR (166)"
  case 167:
    ret= "SUBOR (167)"
  case 171:
    ret= "Kaiser KS-7058"
  case 184:
    ret= "Sunsoft-1 IC"
  case 185:
    if self.Submapper >= 0 {
      ret= fmt.Sprintf ( "CNROM (w. prot. %d)", self.Submapper )
    } else {
      ret= "CNROM (w. prot.)"
    }
  case 206:
    switch self.Submapper {
    case 0:
      ret= "Namcot 118"
    case 1:
      ret= "Namcot 3407/3417/3451"
    default:
      ret= "Namcot 118 PCB variants"
    }
  case 210:
    switch self.Submapper {
    case 1:
      ret= "Namco 175"
    case 2:
      ret= "Namco 340"
    default:
      ret= "Namco 175/340"
    }
    // Sense nom, o no tan rellevants
  default:
    ret= fmt.Sprintf ( "INES Mapper %03d", self.Mapper )
    if self.Submapper > 0 {
      ret+= fmt.Sprintf ( " (%d)", self.Submapper )
    }
  }
  
  return ret
  
} // end getMapper


func (self *_NES_Metadata) getTvSystem() string {
  
  switch self.TV_System {
  case _NES_TV_SYSTEM_NTSC:
    return "NTSC"
  case _NES_TV_SYSTEM_PAL:
    return "PAL"
  case _NES_TV_SYSTEM_MULTIPLE:
    return "Multiregió"
  case _NES_TV_SYSTEM_DENDY:
    return "Dendy"
  default:
    return "Desconegut"
  }
  
} // end getTvSystem


func (self *_NES_Metadata) getVsPpuType() string {
  
  switch self.VS_PPU_Type {
  case 0x0:
    return "RP2C03B"
  case 0x1:
    return "RP2C03G"
  case 0x2:
    return "RP2C04-0001"
  case 0x3:
    return "RP2C04-0002"
  case 0x4:
    return "RP2C04-0003"
  case 0x5:
    return "RP2C04-0004"
  case 0x6:
    return "RC2C03B"
  case 0x7:
    return "RC2C03C"
  case 0x8:
    return "RC2C05-01"
  case 0x9:
    return "RC2C05-02"
  case 0xa:
    return "RC2C05-03"
  case 0xb:
    return "RC2C05-04"
  case 0xc:
    return "RC2C05-05"
  default:
    return fmt.Sprintf ( "UNK (%02X)", self.VS_PPU_Type )
  }
  
} // end getVsPpuType


func (self *_NES_Metadata) getVsCopyProtection() string {
  
  switch self.VS_CopyProtection {
  case 0x0:
    return "Vs. Unisystem (normal)"
  case 0x1:
    return "Vs. Unisystem (RBI Baseball protection)"
  case 0x2:
    return "Vs. Unisystem (TKO Boxing protection)"
  case 0x3:
    return "Vs. Unisystem (Super Xevious protection)"
  case 0x4:
    return "Vs. Unisystem (Vs. Ice Climber Japan protection)"
  case 0x5:
    return "Vs. Dual System (normal)"
  case 0x6:
    return "Vs. Dual System (Raid on Bungeling Bay protection)"
  default:
    return fmt.Sprintf ( "UNK (%02X)", self.VS_PPU_Type )
  }
  
} // end getVsCopyProtection


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
    header.Submapper= -1
    
    // Other
  } else {
    f7:= uint8(data[7])
    header.Mapper= int(uint32((f6>>4) | (f7&0xf0)))
    header.Submapper= -1
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


func _NES_set_vs_system_extra( header *_NES_Metadata, data []byte, ftype int ) {

  if ftype == _NES_TYPE_INES_2_0 &&
    header.Console == _NES_CONSOLE_TYPE_VS_SYSTEM {
    f13:= uint8(data[13])
    header.VS_PPU_Type= f13&0x0f
    header.VS_CopyProtection= f13>>4
  } else {
    header.VS_PPU_Type= 0
    header.VS_CopyProtection= 0
  }
  
} // _NES_set_vs_system_extra


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

  // Vs. System Type
  _NES_set_vs_system_extra ( header, data, ftype )

  // Miscellaneous
  if ftype == _NES_TYPE_INES_2_0 {
    header.NumMiscellaneous= int(uint8(data[14])&0x3)
  } else {
    header.NumMiscellaneous= -1
  }

  // Default Expansion Device
  if ftype == _NES_TYPE_INES_2_0 {
    header.DefaultExpansionDev= uint8(data[15])&0x3f
  } else {
    header.DefaultExpansionDev= 0
  }
  
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



func (self *NES) GetMetadata(file_name string) (string,error) {

  // Obri
  fd,err:= os.Open ( file_name )
  if err != nil { return "",err }
  defer fd.Close ()
  
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

  // Comprova grandària capçalera és coherent.
  tmp_size:= 16 + md.PRG_Size + md.CHR_Size
  if tmp_size > size {
    return "",fmt.Errorf ( "La informació de la capçalera és incoherent"+
      " amb la grandària de la ROM" )
  }

  // Calcula MD5
  h:= md5.New ()
  if _,err:= io.Copy ( h, fd ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el MD5: %s", err )
  }
  md.RealMD5= fmt.Sprintf ( "%x", h.Sum ( nil ) )

  // Calcula SHA1
  if _,err:= fd.Seek ( 16, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el SHA1: %s", err )
  }
  h2:= sha1.New ()
  if _,err:= io.Copy ( h2, fd ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut calcular el SHA1: %s", err )
  }
  md.RealSHA1= fmt.Sprintf ( "%x", h2.Sum ( nil ) )
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
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

  // Parseja
  md:= _NES_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[NES] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue

  // Real MD5
  kv= &KeyValue{"md5 (ROM)",md.RealMD5}
  v= append(v,kv)

  // Real SHA1
  kv= &KeyValue{"sha1 (ROM)",md.RealSHA1}
  v= append(v,kv)

  // PRG Size
  var prg_size string
  if md.PRG_Size%1024 == 0 {
    prg_size= fmt.Sprintf ( "%d KB", md.PRG_Size/1024 )
  } else {
    prg_size= fmt.Sprintf ( "%d B", md.PRG_Size )
  }
  kv= &KeyValue{"Grandària PRG ROM",prg_size}
  v= append(v,kv)
  
  // CHR Size
  var chr_size string
  if md.CHR_Size%1024 == 0 {
    chr_size= fmt.Sprintf ( "%d KB", md.CHR_Size/1024 )
  } else {
    chr_size= fmt.Sprintf ( "%d B", md.CHR_Size )
  }
  kv= &KeyValue{"Grandària CHR ROM",chr_size}
  v= append(v,kv)

  // Consola
  kv= &KeyValue{"Consola",md.getConsole ()}
  v= append(v,kv)

  // Mapper
  kv= &KeyValue{"Mapper",md.getMapper ()}
  v= append(v,kv)

  // TV System
  if md.TV_System != _NES_TV_SYSTEM_UNK {
    kv= &KeyValue{"Sistem TV",md.getTvSystem ()}
    v= append(v,kv)
  }
  
  // Mirroring
  var mirroring string
  switch md.Mirroring {
  case _NES_MIRRORING_HORIZONTAL:
    mirroring= "Horitzontal"
  case _NES_MIRRORING_VERTICAL:
    mirroring= "Vertical"
  case _NES_MIRRORING_FOUR_SCREEN:
    mirroring= "Quatre pantalles"
  }
  kv= &KeyValue{"Mirroring",mirroring}
  v= append(v,kv)
  
  // Sram
  if md.Sram {
    kv= &KeyValue{"RAM estàtica","Sí"}
  } else {
    kv= &KeyValue{"RAM estàtica","No"}
  }
  v= append(v,kv)

  // Trainer
  if md.Trainer {
    kv= &KeyValue{"Conté trainer","Sí"}
  } else {
    kv= &KeyValue{"Conté trainer","No"}
  }
  v= append(v,kv)

  // PRG RAM Size
  if md.PRG_RAM_Size > 0 {
    var prg_ram_size string
    if md.PRG_RAM_Size%1024 == 0 {
      prg_ram_size= fmt.Sprintf ( "%d KB", md.PRG_RAM_Size/1024 )
    } else {
      prg_ram_size= fmt.Sprintf ( "%d B", md.PRG_RAM_Size )
    }
    kv= &KeyValue{"Grandària PRG RAM",prg_ram_size}
    v= append(v,kv)
  }

  // NVRAM/EEPROM RAM Size
  if md.PRG_NVRAM_EEPROM_Size > 0 {
    var prg_nvram_size string
    if md.PRG_NVRAM_EEPROM_Size%1024 == 0 {
      prg_nvram_size= fmt.Sprintf ( "%d KB", md.PRG_NVRAM_EEPROM_Size/1024 )
    } else {
      prg_nvram_size= fmt.Sprintf ( "%d B", md.PRG_NVRAM_EEPROM_Size )
    }
    kv= &KeyValue{"Grandària PRG NVRAM/EEPROM",prg_nvram_size}
    v= append(v,kv)
  }

  // CHR RAM Size
  if md.CHR_RAM_Size > 0 {
    var chr_ram_size string
    if md.CHR_RAM_Size%1024 == 0 {
      chr_ram_size= fmt.Sprintf ( "%d KB", md.CHR_RAM_Size/1024 )
    } else {
      chr_ram_size= fmt.Sprintf ( "%d B", md.CHR_RAM_Size )
    }
    kv= &KeyValue{"Grandària CHR RAM",chr_ram_size}
    v= append(v,kv)
  }

  // CHR NVRAM Size
  if md.CHR_NVRAM_Size > 0 {
    var chr_nvram_size string
    if md.CHR_NVRAM_Size%1024 == 0 {
      chr_nvram_size= fmt.Sprintf ( "%d KB", md.CHR_NVRAM_Size/1024 )
    } else {
      chr_nvram_size= fmt.Sprintf ( "%d B", md.CHR_NVRAM_Size )
    }
    kv= &KeyValue{"Grandària CHR NVRAM",chr_nvram_size}
    v= append(v,kv)
  }
  
  // VS - Tipus PPU
  if md.Console == _NES_CONSOLE_TYPE_VS_SYSTEM {
    kv= &KeyValue{"Vs. Tipus PPU",md.getVsPpuType ()}
    v= append(v,kv)
  }

  // VS - Protecció
  if md.Console == _NES_CONSOLE_TYPE_VS_SYSTEM {
    kv= &KeyValue{"Vs. Protecció",md.getVsCopyProtection ()}
    v= append(v,kv)
  }

  // Roms addicionals
  if md.NumMiscellaneous > 0 {
    kv= &KeyValue{"Nº. ROM addicional",fmt.Sprintf("%d",md.NumMiscellaneous)}
    v= append(v,kv)
  }

  // Dispositius per defecte
  if md.DefaultExpansionDev > 0 {
    kv= &KeyValue{"Dispositius",md.getDefaultExpansionDevice ()}
    v= append(v,kv)
  }
  
  return v
  
} // end ParseMetadata
