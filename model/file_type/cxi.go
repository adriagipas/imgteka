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
 *  cxi.go - Tipus de fitxer CTR Executable Image de 3DS.
 */

package file_type

import(
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "image/color"
  "log"
  "strings"

  "github.com/adriagipas/imgcp/citrus"
  "github.com/adriagipas/imgteka/view"
  "golang.org/x/text/encoding"
  "golang.org/x/text/encoding/unicode"
)


/****************/
/* PART PRIVADA */
/****************/

// Pot tornar nil sense error si no té el fitxer icon
func _CXI_GetIconData( exefs *citrus.ExeFS ) ([]byte,error) {

  var mem [0x36c0]byte
  
  for i:= 0; i < len(exefs.Files); i++ {
    if strings.ToLower ( exefs.Files[i].Name ) == "icon" {

      // Obri
      fd,err:= exefs.OpenIndex ( i )
      if err != nil { return nil,err }
      defer fd.Close ()

      // Llig
      n,err:= fd.Read ( mem[:] )
      if err != nil {
        return nil,fmt.Errorf ( "Error llegint el fitxer 'icon': %s", err )
      }
      if n != len(mem) {
        return nil,errors.New ( "Error llegint el fitxer 'icon'" )
      }

      // Comprovació magic number
      if mem[0]!='S' || mem[1]!='M' || mem[2]!='D' || mem[3]!='H' {
        return nil,fmt.Errorf ( "El fitxer 'icon' no és de tipus SMDH" )
      }

      return mem[:],nil
      
    }
  }

  return nil,nil
  
} // end _CXI_GetIconData


type _CXI_Color struct {
  value uint16
}


func (self _CXI_Color) RGBA() (r,g,b,a uint32) {

  r= uint32(uint16((float32((self.value>>11)&0x1f)/31.0)*65535.0))
  g= uint32(uint16((float32((self.value>>5)&0x3f)/63.0)*65535.0))
  b= uint32(uint16((float32(self.value&0x1f)/31.0)*65535.0))
  a= 0xffff

  return
  
} // end _CXI_Color.RGBA


type _CXI_ColorModel struct {
}


func (self *_CXI_ColorModel) Convert(c color.Color) color.Color {
  return c
} // _CXI_ColorModel.Convert


type _CXI_Icon struct {

  data []byte
  
}


func (self *_CXI_Icon) At( x,y int ) color.Color {

  const WIDTH = 6    // En tiles
  const HEIGHT = 6   // En tiles
  const TILE_WIDTH = 8
  const TILE_HEIGHT = 8
  const TILE_LINE_SIZE = 2*TILE_WIDTH
  const TILE_SIZE = TILE_LINE_SIZE*TILE_HEIGHT

  // Inicialitza offset al principi del tile
  tile:= (y/TILE_HEIGHT)*WIDTH + (x/TILE_WIDTH)
  offset:= tile*TILE_SIZE

  // Calcula posició píxel dins del tile de manera recursiva
  pos_x:= x%TILE_WIDTH
  pos_y:= y%TILE_HEIGHT
  // --> Primera divisió
  offset+= ((((pos_y>>2)<<1) + (pos_x>>2))<<5)
  pos_x&= 0x3
  pos_y&= 0x3
  // --> Segona divisió
  offset+= ((((pos_y>>1)<<1) + (pos_x>>1))<<3)
  pos_x&= 0x1
  pos_y&= 0x1
  // --> Tercera divisió
  offset+= (((pos_y<<1) + pos_x)<<1)

  // Obté color
  color:= uint16(self.data[offset]) | (uint16(self.data[offset+1])<<8)

  return _CXI_Color{color}
  
} // _CXI_Icon.At


func (self *_CXI_Icon) Bounds() image.Rectangle {
  return image.Rectangle{
    Min:image.Point{
      X:0,
      Y:0,
    },
    Max:image.Point{
      X:48,
      Y:48,
    },
  }
} // end _CXI_Icon.Bounds


func (self *_CXI_Icon) ColorModel() color.Model {
  return &_CXI_ColorModel{}
} // end ColorModel


type _CXI_Title struct {

  Short     string
  Long      string
  Publisher string
  
}


func (self *_CXI_Title) Init( dec *encoding.Decoder, data []byte ) error {

  if aux,err:= dec.Bytes ( data[:0x80] ); err == nil {
    self.Short= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x80:0x180] ); err == nil {
    self.Long= BytesToStr_trim_0s(aux)
  }
  if aux,err:= dec.Bytes ( data[0x180:0x200] ); err == nil {
    self.Publisher= BytesToStr_trim_0s(aux)
  }

  return nil
  
} // end Init


func (self *_CXI_Title) Format() string {
  
  if self.Short == "" && self.Long == "" && self.Publisher == "" {
    return ""
  }
  
  ret:= fmt.Sprintf ( "%s / %s (%s)",
    strings.Replace(self.Short,"\n"," ~ ",-1),
    strings.Replace(self.Long,"\n"," ~ ",-1),
    strings.Replace(self.Publisher,"\n"," ~ ",-1),
  )
  
  return ret
  
} // end Format


const (
  _CXI_AGE_RATING_TYPE_UNUSED          = 0
  _CXI_AGE_RATING_TYPE_PENDING         = 1
  _CXI_AGE_RATING_TYPE_NO_RESTRICTIONS = 2
  _CXI_AGE_RATING_TYPE_RESTRICTIONS    = 3
)

type _CXI_AgeRating struct {

  Type   int
  MinAge int
  
}


func (self *_CXI_AgeRating) Init( val byte ) {

  if (val&0x80) != 0 {
    if (val&0x40) != 0 {
      self.Type= _CXI_AGE_RATING_TYPE_PENDING
    } else if (val&0x20) != 0 {
      self.Type= _CXI_AGE_RATING_TYPE_NO_RESTRICTIONS
    } else {
      self.Type= _CXI_AGE_RATING_TYPE_RESTRICTIONS
      self.MinAge= int(val-0x80)
    }
  } else {
    self.Type= _CXI_AGE_RATING_TYPE_UNUSED
  }
  
} // end Init


func (self *_CXI_AgeRating) Format() string {

  switch self.Type {
  case _CXI_AGE_RATING_TYPE_UNUSED:
    return ""
  case _CXI_AGE_RATING_TYPE_PENDING:
    return "Pendent"
  case _CXI_AGE_RATING_TYPE_NO_RESTRICTIONS:
    return "Sense restriccions d'edat"
  case _CXI_AGE_RATING_TYPE_RESTRICTIONS:
    if self.MinAge == 0 {
      return "Sense restriccions d'edat"
    } else {
      return fmt.Sprintf ( "+%d", self.MinAge )
    }
  default:
    return ""
  }
  
} // end Format


const (
  _CXI_REGION_LOCKOUT_JPN = 0x01
  _CXI_REGION_LOCKOUT_USA = 0x02
  _CXI_REGION_LOCKOUT_EUR = 0x04
  _CXI_REGION_LOCKOUT_AUS = 0x08
  _CXI_REGION_LOCKOUT_CHN = 0x10
  _CXI_REGION_LOCKOUT_KOR = 0x20
  _CXI_REGION_LOCKOUT_TWN = 0x40
)

const (
  _CXI_FLAGS_REQUIRE_EULA     = 0x0008
  _CXI_FLAGS_REQUIRE_RATING   = 0x0040
  _CXI_FLAGS_NEW3DS_EXCLUSIVE = 0x1000
)

type _CXI_Metadata struct {

  Header citrus.NCCH_Header

  // Títols
  Japanese   _CXI_Title
  English    _CXI_Title
  French     _CXI_Title
  German     _CXI_Title
  Italian    _CXI_Title
  Spanish    _CXI_Title
  SChinese   _CXI_Title
  Korean     _CXI_Title
  Dutch      _CXI_Title
  Portuguese _CXI_Title
  Russian    _CXI_Title
  TChinese   _CXI_Title

  // Ratings
  CERO      _CXI_AgeRating
  ESRB      _CXI_AgeRating
  USK       _CXI_AgeRating
  PEGI_GEN  _CXI_AgeRating
  PEGI_PRT  _CXI_AgeRating
  PEGI_BBFC _CXI_AgeRating
  COB       _CXI_AgeRating
  GRB       _CXI_AgeRating
  CGSRR     _CXI_AgeRating

  // Altres
  RegionLockout     uint32
  Flags             uint32
  EULA_VersionMinor uint8
  EULA_VersionMajor uint8
  
}


func (self *_CXI_Metadata) InitFromSMDH( data []byte ) error {

  // Llig títols
  dec:= unicode.UTF16(unicode.LittleEndian,unicode.IgnoreBOM).NewDecoder ()
  if err:= self.Japanese.Init ( dec, data[0x008:0x208] ); err != nil {
    return err
  }
  if err:= self.English.Init ( dec, data[0x208:0x408] ); err != nil {
    return err
  }
  if err:= self.French.Init ( dec, data[0x408:0x608] ); err != nil {
    return err
  }
  if err:= self.German.Init ( dec, data[0x608:0x808] ); err != nil {
    return err
  }
  if err:= self.Italian.Init ( dec, data[0x808:0xa08] ); err != nil {
    return err
  }
  if err:= self.Spanish.Init ( dec, data[0xa08:0xc08] ); err != nil {
    return err
  }
  if err:= self.SChinese.Init ( dec, data[0xc08:0xe08] ); err != nil {
    return err
  }
  if err:= self.Korean.Init ( dec, data[0xe08:0x1008] ); err != nil {
    return err
  }
  if err:= self.Dutch.Init ( dec, data[0x1008:0x1208] ); err != nil {
    return err
  }
  if err:= self.Portuguese.Init ( dec, data[0x1208:0x1408] ); err != nil {
    return err
  }
  if err:= self.Russian.Init ( dec, data[0x1408:0x1608] ); err != nil {
    return err
  }
  if err:= self.TChinese.Init ( dec, data[0x1608:0x1808] ); err != nil {
    return err
  }

  // Ratings
  self.CERO.Init ( data[0x2008] )
  self.ESRB.Init ( data[0x2009] )
  self.USK.Init ( data[0x200b] )
  self.PEGI_GEN.Init ( data[0x200c] )
  self.PEGI_PRT.Init ( data[0x200e] )
  self.PEGI_BBFC.Init ( data[0x200f] )
  self.COB.Init ( data[0x2010] )
  self.GRB.Init ( data[0x2011] )
  self.CGSRR.Init ( data[0x2012] )

  // Altres
  self.RegionLockout= uint32(data[0x2018]) |
    (uint32(data[0x2019])<<8) |
    (uint32(data[0x201a])<<16) |
    (uint32(data[0x201b])<<24)
  self.Flags= uint32(data[0x2028]) |
    (uint32(data[0x2029])<<8) |
    (uint32(data[0x202a])<<16) |
    (uint32(data[0x202b])<<24)
  self.EULA_VersionMinor= data[0x202c]
  self.EULA_VersionMajor= data[0x202d]
  
  return nil
  
} // end InitFromSMDH


func (self *_CXI_Metadata) ParseMetadata(
  v []view.StringPair,
) []view.StringPair {

  var kv *KeyValue
  
  // NCCH Header
  kv= &KeyValue{"Grandària capçalera",NumBytesToStr(uint64(self.Header.Size))}
  v= append(v,kv)
  kv= &KeyValue{"Identificador",fmt.Sprintf("%016x",self.Header.Id)}
  v= append(v,kv)
  kv= &KeyValue{"Codi fabricant",self.Header.MakerCode}
  v= append(v,kv)
  kv= &KeyValue{"Versió",fmt.Sprintf("%04x",self.Header.Version)}
  v= append(v,kv)
  //kv= &KeyValue{"Identificador programa", // <-- Redundant
  //  fmt.Sprintf("%016x",self.Header.ProgramId)}
  //v= append(v,kv)
  kv= &KeyValue{"Codi producte",
    strings.TrimRight ( self.Header.ProductCode, "\000" )}
  v= append(v,kv)

  // Títols
  if aux:= self.Japanese.Format (); aux != "" {
    kv= &KeyValue{"Títol (Japonès)",aux}
    v= append(v,kv)
  }
  if aux:= self.English.Format (); aux != "" {
    kv= &KeyValue{"Títol (Anglès)",aux}
    v= append(v,kv)
  }
  if aux:= self.French.Format (); aux != "" {
    kv= &KeyValue{"Títol (Francès)",aux}
    v= append(v,kv)
  }
  if aux:= self.German.Format (); aux != "" {
    kv= &KeyValue{"Títol (Alemany)",aux}
    v= append(v,kv)
  }
  if aux:= self.Italian.Format (); aux != "" {
    kv= &KeyValue{"Títol (Italià)",aux}
    v= append(v,kv)
  }
  if aux:= self.Spanish.Format (); aux != "" {
    kv= &KeyValue{"Títol (Espanyol)",aux}
    v= append(v,kv)
  }
  if aux:= self.SChinese.Format (); aux != "" {
    kv= &KeyValue{"Títol (Xinès simplificat)",aux}
    v= append(v,kv)
  }
  if aux:= self.Korean.Format (); aux != "" {
    kv= &KeyValue{"Títol (Coreà)",aux}
    v= append(v,kv)
  }
  if aux:= self.Dutch.Format (); aux != "" {
    kv= &KeyValue{"Títol (Neerlandès)",aux}
    v= append(v,kv)
  }
  if aux:= self.Portuguese.Format (); aux != "" {
    kv= &KeyValue{"Títol (Portuguès)",aux}
    v= append(v,kv)
  }
  if aux:= self.Russian.Format (); aux != "" {
    kv= &KeyValue{"Títol (Rus)",aux}
    v= append(v,kv)
  }
  if aux:= self.TChinese.Format (); aux != "" {
    kv= &KeyValue{"Títol (Xinès tradicional)",aux}
    v= append(v,kv)
  }

  // Ratings
  if aux:= self.CERO.Format (); aux != "" {
    kv= &KeyValue{"CERO",aux}
    v= append(v,kv)
  }
  if aux:= self.ESRB.Format (); aux != "" {
    kv= &KeyValue{"ESRB",aux}
    v= append(v,kv)
  }
  if aux:= self.USK.Format (); aux != "" {
    kv= &KeyValue{"USK",aux}
    v= append(v,kv)
  }
  if aux:= self.PEGI_GEN.Format (); aux != "" {
    kv= &KeyValue{"PEGI",aux}
    v= append(v,kv)
  }
  if aux:= self.PEGI_PRT.Format (); aux != "" {
    kv= &KeyValue{"PEGI (Portugal)",aux}
    v= append(v,kv)
  }
  if aux:= self.PEGI_BBFC.Format (); aux != "" {
    kv= &KeyValue{"PEGI (Regne Unit)",aux}
    v= append(v,kv)
  }
  if aux:= self.COB.Format (); aux != "" {
    kv= &KeyValue{"COB",aux}
    v= append(v,kv)
  }
  if aux:= self.GRB.Format (); aux != "" {
    kv= &KeyValue{"GRB",aux}
    v= append(v,kv)
  }
  if aux:= self.CGSRR.Format (); aux != "" {
    kv= &KeyValue{"CGSRR",aux}
    v= append(v,kv)
  }

  // Regió
  var region string
  if self.RegionLockout==0x7FFFFFFF {
    region= "Totes les regions"
  } else {
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_JPN)!=0 {
      region+= ", Japó"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_USA)!=0 {
      region+= ", Estats Units"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_EUR)!=0 {
      region+= ", Europa"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_AUS)!=0 {
      region+= ", Austràlia"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_CHN)!=0 {
      region+= ", Xina"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_KOR)!=0 {
      region+= ", Corea"
    }
    if (self.RegionLockout&_CXI_REGION_LOCKOUT_KOR)!=0 {
      region+= ", Taiwan"
    }
    if len(region)>0 { region= region[2:] }
  }
  kv= &KeyValue{"Regions",region}
  v= append(v,kv)

  // Flags
  var flags string
  if (self.Flags&_CXI_FLAGS_REQUIRE_EULA)!=0 {
    flags+= ", Requereix EULA"
  }
  if (self.Flags&_CXI_FLAGS_REQUIRE_RATING)!=0 {
    flags+= ", Requereix qualificació d'edat"
  }
  if (self.Flags&_CXI_FLAGS_NEW3DS_EXCLUSIVE)!=0 {
    flags+= ", Exclusiu New 3DS"
  }
  if len(flags)>0 {
    flags= fmt.Sprintf ( "%02x (%s)", self.Flags, flags[2:] )
  } else {
    flags= fmt.Sprintf ( "%08x", self.Flags )
  }
  kv= &KeyValue{"Flags",flags}
  v= append(v,kv)

  // Eula
  kv= &KeyValue{"EULA",
    fmt.Sprintf("%d.%d",self.EULA_VersionMajor,self.EULA_VersionMinor)}
  v= append(v,kv)
  
  return v
  
} // end ParseMetadata




/****************/
/* PART PÚBLICA */
/****************/

type CXI struct {
}


func (self *CXI) GetImage( file_name string ) (image.Image,error) {

  // Obri fitxer
  ncch,err:= citrus.NewNCCH ( file_name )
  if err != nil { return nil,err }
  
  // Comprova que és executable
  if ncch.Header.Type != citrus.NCCH_TYPE_CXI {
    return nil,errors.New ( "No és un fitxer NCCH executable (CXI)" )
  }

  // Obté ExeFS
  exefs,err:= ncch.GetExeFS ()
  if err != nil { return nil,err }
  if exefs == nil {
    return nil,errors.New ( "No s'ha trobat una partició ExeFS" )
  }

  // Obté metadades del fitxer icon
  icon_data,err:= _CXI_GetIconData ( exefs )
  if err != nil { return nil,err }
  if icon_data == nil { return nil,errors.New ( "No conté icona" ) }

  return &_CXI_Icon{data:icon_data[0x24c0:0x36c0]},nil
  
} // end GetImage


func (self *CXI) GetMetadata(file_name string) (string,error) {

  // Inicialitza
  md:= _CXI_Metadata{}

  // Obri fitxer
  ncch,err:= citrus.NewNCCH ( file_name )
  if err != nil { return "",err }
  md.Header= ncch.Header

  // Comprova que és executable
  if ncch.Header.Type != citrus.NCCH_TYPE_CXI {
    return "",errors.New ( "No és un fitxer NCCH executable (CXI)" )
  }

  // Obté ExeFS
  exefs,err:= ncch.GetExeFS ()
  if err != nil { return "",err }
  if exefs == nil {
    return "",errors.New ( "No s'ha trobat una partició ExeFS" )
  }

  // Obté metadades del fitxer icon
  if icon_data,err:= _CXI_GetIconData ( exefs ); err != nil {
    return "",err
  } else if icon_data != nil {
    if err:= md.InitFromSMDH ( icon_data ); err != nil {
      return "",err
    }
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *CXI) GetName() string { return "CTR Executable Image" }
func (self *CXI) GetShortName() string { return "CXI" }
func (self *CXI) IsImage() bool { return true }


func (self *CXI) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  md:= _CXI_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[CXI] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  v= md.ParseMetadata ( v )
  
  return v
  
} // end ParseMetadata
