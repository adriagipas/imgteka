/*
 * Copyright 2024 Adrià Giménez Pastor.
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
 *  nds.go - Tipus de fitxer ROM de Nintendo DS.
 */

package file_type

import (
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "image/color"
  "log"
  "os"
  "strings"
  
  "github.com/adriagipas/imgteka/view"
  "golang.org/x/text/encoding/unicode"
)




/****************/
/* PART PRIVADA */
/****************/

func _NDS_CRC16 ( data []byte ) uint16 {
  
  val:= [8]uint16{
    0xc0c1, 0xc181, 0xc301, 0xc601,
    0xcc01, 0xd801, 0xf001, 0xa001}
  crc:= uint32(0xffff)
  for _,b:= range data {
    crc^= uint32(b)
    for j:= 0; j < 8; j++ {
      if (crc&0x1) != 0x00 {
        crc= (crc>>1)^(uint32(val[j])<<(7-j))
      } else {
        crc>>= 1
      }
    }
  }
  
  return uint16(crc)
  
} // _NDS_CRC16


type _NDS_Color struct {
  r,g,b,a uint32
}


func (self _NDS_Color) RGBA() (r,g,b,a uint32) {
  return self.r,self.g,self.b,self.a
} // end _NDS_Color.RGBA


type _NDS_ColorModel struct {
}


func (self *_NDS_ColorModel) Convert(c color.Color) color.Color {
  return c
} // _NDS_ColorModel.Convert


type _NDS_Icon struct {

  data []byte
  pal  [16]_NDS_Color
  
}


func (self *_NDS_Icon) At( x,y int ) color.Color {

  const WIDTH = 4    // En tiles
  const HEIGHT = 4   // En tiles
  const TILE_WIDTH = 8
  const TILE_HEIGHT = 8
  const TILE_LINE_SIZE = 4
  const TILE_SIZE = TILE_LINE_SIZE*TILE_HEIGHT
  
  // Calcula offset byte
  offset:= 0x0020 + // Base en memòria
    (y/TILE_HEIGHT)*WIDTH*TILE_SIZE + // fila de tiles que toca
    (x/TILE_WIDTH)*TILE_SIZE + // Principi del tile
    (y%TILE_HEIGHT)*TILE_LINE_SIZE + // Principi de la línia
    (x%TILE_WIDTH)/2 // Cada byte té 2 píxels
  b:= self.data[offset]
  if (x%TILE_WIDTH)%2 == 1 {
    b>>=4
  } else {
    b&= 0x0f
  }

  return self.pal[b]
  
} // _NDS_Icon.At


func (self *_NDS_Icon) Bounds() image.Rectangle {
  return image.Rectangle{
    Min:image.Point{
      X:0,
      Y:0,
    },
    Max:image.Point{
      X:31,
      Y:31,
    },
  }
} // end _NDS_Icon.Bounds


func (self *_NDS_Icon) ColorModel() color.Model {
  return &_NDS_ColorModel{}
} // end ColorModel


func _NDS_NewIcon( data []byte ) *_NDS_Icon {

  // Crea Objecte
  ret:= _NDS_Icon{
    data: data,
  }

  // Crea paleta
  ret.pal[0].r= 0
  ret.pal[0].g= 0
  ret.pal[0].b= 0
  ret.pal[0].a= 0
  var color uint16
  for i:= 1; i < 16; i++ {
    color= uint16(data[0x220+i*2]) | (uint16(data[0x220+i*2+1])<<8)
    ret.pal[i].a= 0xffff
    ret.pal[i].r= uint32(uint16((float32(color&0x1f)/31.0)*65535.0))
    ret.pal[i].g= uint32(uint16((float32((color>>5)&0x1f)/31.0)*65535.0))
    ret.pal[i].b= uint32(uint16((float32((color>>10)&0x1f)/31.0)*65535.0))
  }
  
  return &ret
  
} // end _NDS_NewIcon


func _NDS_GetIconTitleData( fd *os.File ) ([]byte,error) {

  // Llig offset
  if _,err:= fd.Seek ( 0x68, 0 ); err != nil {
    return nil,errors.New ( "Error llegint l'offset de la secció Icon/Title" )
  }
  var mem_off [4]byte
  if n,err:= fd.Read ( mem_off[:] ); err != nil || n != 4 {
    return nil,errors.New ( "Error llegint l'offset de la secció Icon/Title" )
  }
  offset:= uint32(mem_off[0]) |
    (uint32(mem_off[1])<<8) |
    (uint32(mem_off[2])<<16) |
    (uint32(mem_off[3])<<24);
  if offset == 0 { // No en té
    return nil,nil
  }

  // Llig versió
  if _,err:= fd.Seek ( int64(uint64(offset)), 0 ); err != nil {
    return nil,errors.New ( "Error llegint la versió de la secció Icon/Title" )
  }
  var mem_ver [2]byte
  if n,err:= fd.Read ( mem_ver[:] ); err != nil || n != 2 {
    return nil,errors.New ( "Error llegint la versió de la secció Icon/Title" )
  }
  version:= uint16(mem_ver[0]) | (uint16(mem_ver[1])<<8)
  
  // Grandària a llegir
  var size uint32
  if version < 2 {
    size= 0x840
  } else if version < 3 {
    size= 0x940
  } else {
    size= 0x1240
  }

  // Llig contingut
  if _,err:= fd.Seek ( int64(uint64(offset)), 0 ); err != nil {
    return nil,errors.New ( "Error llegint la secció Icon/Title" )
  }
  mem:= make([]byte,size)
  if n,err:= fd.Read ( mem ); err != nil || uint32(n) != size {
    return nil,errors.New ( "Error llegint la secció Icon/Title" )
  }
  
  return mem,nil
  
} // _NDS_GetIconTitleData


type _NDS_Metadata struct {

  TitleHeader        string
  GameCode           string
  MakerCode          string
  UnitCode           uint8
  DeviceCapacity     uint32 // Mesurat en KB
  RegionCode         uint8
  RomVersion         uint8
  HeaderChecksum     uint16
  RealHeaderChecksum uint16
  TitleJapanese      string
  TitleEnglish       string
  TitleFrench        string
  TitleGerman        string
  TitleItalian       string
  TitleSpanish       string
  TitleChinese       string
  TitleKorean        string
  
}


func (self *_NDS_Metadata) InitHeaderFields( data []byte ) error {

  // Llig camps
  self.TitleHeader= BytesToStr_trim_0s ( data[:12] )
  self.GameCode= string(data[0x0c:0x10])
  self.MakerCode= string(data[0x10:0x12])
  self.UnitCode= uint8(data[0x12])
  self.DeviceCapacity= uint32(128<<uint8(data[0x14]))
  self.RegionCode= uint8(data[0x1d])
  self.RomVersion= uint8(data[0x1e])
  self.HeaderChecksum= uint16(data[0x15e]) | (uint16(data[0x15f])<<8)
  
  // Calcula checksum real
  self.RealHeaderChecksum= _NDS_CRC16(data[0:0x15e])
  
  return nil
  
} // end InitHeaderFields


func (self *_NDS_Metadata) InitTitles( data []byte ) {

  // Crea decoder
  dec:= unicode.UTF16(unicode.LittleEndian,unicode.IgnoreBOM).NewDecoder ()
  
  // Títols estàndard
  if aux,err:= dec.Bytes ( data[0x240:0x340] ); err == nil {
    self.TitleJapanese= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x340:0x440] ); err == nil {
    self.TitleEnglish= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x440:0x540] ); err == nil {
    self.TitleFrench= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x540:0x640] ); err == nil {
    self.TitleGerman= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x640:0x740] ); err == nil {
    self.TitleItalian= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x740:0x840] ); err == nil {
    self.TitleSpanish= BytesToStr_trim_0s(aux)
  }

  // Altres títols
  version:= uint16(data[0]) | (uint16(data[1])<<8)
  if version >= 2 {
    if aux,err:= dec.Bytes ( data[0x840:0x940] ); err == nil {
      self.TitleChinese= BytesToStr_trim_0s(aux)
    }
    if version >= 3 {
      if aux,err:= dec.Bytes ( data[0x940:0xa40] ); err == nil {
        self.TitleKorean= BytesToStr_trim_0s(aux)
      }
    }
  }
  
} // end InitTitles




/****************/
/* PART PÚBLICA */
/****************/

type NDS struct {
}


func (self *NDS) GetImage( file_name string ) (image.Image,error) {
  
  // Obri fitxer
  fd,err:= os.Open ( file_name )
  if err != nil { return nil,err }
  defer fd.Close ()
  
  // Obté dades
  data,err:= _NDS_GetIconTitleData ( fd )
  if err != nil { return nil,err }
  if data == nil { return nil,nil } // No té imatge

  return _NDS_NewIcon ( data ),nil
  
} // end GetImage


func (self *NDS) GetMetadata(file_name string) (string,error) {

  // Inicialitza
  md:= _NDS_Metadata{}
  
  // Obri
  fd,err:= os.Open ( file_name )
  if err != nil { return "",err }
  defer fd.Close ()

  // Comprova grandària
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'han pogut obtindre les metadades: %s", err )
  }
  size:= info.Size ()
  if size < 0x1000 {
    return "",errors.New (
      "La grandària del fitxer no és correspon amb el d'una"+
        " ROM de Nintendo DS" )
  }

  // Llig la capçalera
  mem:= make([]byte,0x1000)
  n,err:= fd.Read ( mem )
  if err != nil { return "",err }
  if n != 0x1000 {
    return "",errors.New ( "Error llegint la capçalera" )
  }
  md.InitHeaderFields ( mem )

  // Llig títols
  mem,err= _NDS_GetIconTitleData ( fd )
  if err != nil { return "",err }
  if mem != nil {
    md.InitTitles ( mem )
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *NDS) GetName() string { return "ROM de Nintendo DS" }
func (self *NDS) GetShortName() string { return "NDS" }
func (self *NDS) IsImage() bool { return true }


func (self *NDS) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _NDS_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[NDS] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue

  // Title Header
  kv= &KeyValue{"Títol capçalera",md.TitleHeader}
  v= append(v,kv)

  // Codi
  kv= &KeyValue{"Codi",md.GameCode}
  v= append(v,kv)

  // Codi fabricant
  kv= &KeyValue{"Codi fabricant",md.MakerCode}
  v= append(v,kv)

  // Dispositiu
  var device string
  switch md.UnitCode {
  case 0x00:
    device= "Nintendo DS"
  case 0x02:
    device= "Nintendo DS/DSi"
  case 0x03:
    device= "Nintendo DSi"
  }
  if device != "" {
    kv= &KeyValue{"Dispositiu",device}
    v= append(v,kv)
  }

  // Capacitat
  kv= &KeyValue{"Grandària (capçalera)",
    NumBytesToStr(uint64(md.DeviceCapacity*1024))}
  v= append(v,kv)

  // Region code
  var region string
  switch md.RegionCode {
  case 0x00:
    region= "Normal"
  case 0x80:
    region= "Xina"
  case 0x40:
    region= "Corea"
  }
  if region != "" {
    kv= &KeyValue{"Regió consola",region}
    v= append(v,kv)
  }

  // Rom versió
  kv= &KeyValue{"Versió ROM",fmt.Sprintf ( "%02x", md.RomVersion )}
  v= append(v,kv)

  // Checksum
  if md.HeaderChecksum == md.RealHeaderChecksum {
    kv= &KeyValue{"Checksum",
      fmt.Sprintf ( "%04x (Sí)", md.HeaderChecksum )}
  } else {
    kv= &KeyValue{"Checksum",
      fmt.Sprintf ( "%04x (No != %04x)",
        md.HeaderChecksum, md.RealHeaderChecksum )}
  }
  v= append(v,kv)

  // Títols
  if md.TitleJapanese != "" {
    kv= &KeyValue{"Títol (Japonès)",
      strings.Replace(md.TitleJapanese,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleEnglish != "" {
    kv= &KeyValue{"Títol (Anglès)",
      strings.Replace(md.TitleEnglish,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleFrench != "" {
    kv= &KeyValue{"Títol (Francès)",
      strings.Replace(md.TitleFrench,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleGerman != "" {
    kv= &KeyValue{"Títol (Alemany)",
      strings.Replace(md.TitleGerman,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleItalian != "" {
    kv= &KeyValue{"Títol (Italià)",
      strings.Replace(md.TitleItalian,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleSpanish != "" {
    kv= &KeyValue{"Títol (Espanyol)",
      strings.Replace(md.TitleSpanish,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleChinese != "" {
    kv= &KeyValue{"Títol (Xinès)",
      strings.Replace(md.TitleChinese,"\n"," | ",-1)}
    v= append(v,kv)
  }
  if md.TitleKorean != "" {
    kv= &KeyValue{"Títol (Coreà)",
      strings.Replace(md.TitleKorean,"\n"," | ",-1)}
    v= append(v,kv)
  }
  
  return v
  
} // end ParseMetada
