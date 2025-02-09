/*
 * Copyright 2025 Adrià Giménez Pastor.
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
 *  fat12.go - Tipus de fitxer imatge de disquet formatat amb FAT12.
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

const _FAT12_SEC_SIZE = 512

type _FAT12_Metadata struct {

  OEM_Name          string
  BytesPerSector    uint16
  SectorsPerCluster uint8
  NumFATs           uint8
  TotalSectors      uint32
  Media             uint8
  SectorsPerTrack   uint16
  NumHeads          uint16
  HiddenSectors     uint32
  VolumeInfoPresent bool
  VolumeID          uint32
  VolumeLabel       string
  
}

func (self *_FAT12_Metadata) Read( data []byte ) error {

  // Comprova que és un boot sector. He relaxat les condicions.
  if data[0]!=0xe9 && (data[0]!=0xeb /*|| data[2]!=0x90*/) {
    return errors.New ( "El primer sector no és executable" )
  }

  // OEM
  self.OEM_Name= BytesToStr_trim_0s ( data[3:11] )

  // BytesPerSector
  self.BytesPerSector= (uint16(data[12])<<8) | uint16(data[11])

  // SectorsPerCluster
  self.SectorsPerCluster= data[13]

  // Num FATs
  self.NumFATs= data[16]

  // Total sectors
  self.TotalSectors= uint32((uint16(data[20])<<8) | uint16(data[19]))
  if self.TotalSectors == 0 {
    self.TotalSectors= uint32(data[32]) |
      (uint32(data[33])<<8) |
      (uint32(data[34])<<16) |
      (uint32(data[35])<<24)
  }

  // Media
  self.Media= data[21]

  // Sectors per track
  self.SectorsPerTrack= uint16(data[24]) | (uint16(data[25])<<8)

  // Nombre capçals
  self.NumHeads= uint16(data[26]) | (uint16(data[27])<<8)

  // Sectors ocults
  self.HiddenSectors= uint32(data[28]) |
    (uint32(data[29])<<8) |
    (uint32(data[30])<<16) |
    (uint32(data[31])<<24)

  // Informació del volum
  self.VolumeInfoPresent= data[38]==0x29
  if self.VolumeInfoPresent {
    self.VolumeID= uint32(data[39]) |
      (uint32(data[40])<<8) |
      (uint32(data[41])<<16) |
      (uint32(data[42])<<24)
    self.VolumeLabel= BytesToStr_trim_0s ( data[43:54] )
    fstype:= BytesToStr_trim_0s ( data[54:62] )
    if fstype != "FAT12" && fstype != "FAT" && fstype != "" {
      return fmt.Errorf ( "Tipus de format desconegut: '%s'", fstype )
    }
  }

  // Comprova la firma
  if data[510] != 0x55 || data[511] != 0xaa {
    return errors.New ( "No es tracta d'un disquet amb format FAT12" )
  }
  
  return nil
  
} // end Read




/****************/
/* PART PÚBLICA */
/****************/

type FAT12 struct {
}


func (self *FAT12) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge un disquet formatat com FAT12" )
} // end GetImage


func (self *FAT12) GetMetadata(file_name string) (string,error) {

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
  if ( size < _FAT12_SEC_SIZE  ) {
    return "",errors.New (
      "La grandària del fitxer no és correspon amb el d'un disquet" )
  }

  // Llig el primer sector
  mem:= make([]byte,_FAT12_SEC_SIZE)
  n,err:= fd.Read ( mem )
  if err != nil { return "",err }
  if n != _FAT12_SEC_SIZE {
    return "",errors.New ( "Error llegint les dades" )
  }
  
  // Llig metadades
  md:= _FAT12_Metadata{}
  if err:= md.Read ( mem ); err != nil {
    return "",err
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *FAT12) GetName() string {
  return "Disquet FAT12 (MS-DOS)"
}
func (self *FAT12) GetShortName() string { return "FAT12" }
func (self *FAT12) IsImage() bool { return false }


func (self *FAT12) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {


  // Parseja
  md:= _FAT12_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[FAT12] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue

  // Nom OEM
  if md.OEM_Name != "" {
    kv= &KeyValue{"OEM",md.OEM_Name}
    v= append(v,kv)
  }

  // Bytes per sector
  kv= &KeyValue{"Bytes per sector",fmt.Sprintf ( "%d", md.BytesPerSector )}
  v= append(v,kv)

  // Sectors per cluster
  kv= &KeyValue{"Sectors per cluster",
    fmt.Sprintf ( "%d", md.SectorsPerCluster )}
  v= append(v,kv)

  // Nombre de FATs
  kv= &KeyValue{"Nre. de FATs",fmt.Sprintf ( "%d", md.NumFATs )}
  v= append(v,kv)
  
  // Sectors totals
  kv= &KeyValue{"Nre. de sectors",fmt.Sprintf ( "%d", md.TotalSectors )}
  v= append(v,kv)

  // Media
  kv= &KeyValue{"Media",fmt.Sprintf ( "%02X", md.Media )}
  v= append(v,kv)

  // Sectors per track
  kv= &KeyValue{"Sectors per track",
    fmt.Sprintf ( "%d", md.SectorsPerTrack )}
  v= append(v,kv)

  // NumHeads
  kv= &KeyValue{"Nre. de capçals",fmt.Sprintf ( "%d", md.NumHeads )}
  v= append(v,kv)

  // HiddenSectors
  kv= &KeyValue{"Sectors ocults",fmt.Sprintf ( "%d", md.HiddenSectors )}
  v= append(v,kv)

  // Informació volum
  if md.VolumeInfoPresent {
    kv= &KeyValue{"Id. volum",fmt.Sprintf ( "%08X", md.VolumeID )}
    v= append(v,kv)
    kv= &KeyValue{"Etiqueta volum",md.VolumeLabel}
    v= append(v,kv)
  }
  
  return v
  
} // end ParseMetadata
