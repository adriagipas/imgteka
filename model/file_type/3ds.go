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
 *  3ds.go - Tipus de fitxer imatge de cartutx de 3DS
 */

package file_type

import (
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "log"
  
  "github.com/adriagipas/imgteka/view"
  "github.com/adriagipas/imgcp/citrus"
)


/****************/
/* PART PRIVADA */
/****************/

type _3DS_Metadata struct {

  Size         int64
  MediaID      uint64
  TitleVersion uint16
  CardRevision uint16
  TitleID      uint64 // Algo de CVer ??
  VersionCVer  uint16
  
  CXI          _CXI_Metadata
  
}




/****************/
/* PART PÚBLICA */
/****************/

type N3DS struct {
}


func (self *N3DS) GetImage( file_name string ) (image.Image,error) {

  // Obri fitxer
  cci,err:= citrus.NewCCI ( file_name )
  if err != nil { return nil,err }
  
  // Obté metadades CXI
  cxi,err:= cci.GetNCCHPartition ( 0 )
  if err != nil { return nil,err }

  // Comprova que és executable
  if cxi.Header.Type != citrus.NCCH_TYPE_CXI {
    return nil,errors.New ( "No conté un executable (CXI)" )
  }

  // Obté ExeFS
  exefs,err:= cxi.GetExeFS ()
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


func (self *N3DS) GetMetadata(file_name string) (string,error) {

  // Inicialitza
  md:= _3DS_Metadata{}

  // Obri fitxer
  cci,err:= citrus.NewCCI ( file_name )
  if err != nil { return "",err }

  // Emplena capçalera CCI
  md.Size= cci.Header.Size
  md.MediaID= cci.Header.MediaID
  md.TitleVersion= cci.Header.TitleVersion
  md.CardRevision= cci.Header.CardRevision
  md.TitleID= cci.Header.TitleID
  md.VersionCVer= cci.Header.VersionCVer

  // Obté metadades CXI
  cxi,err:= cci.GetNCCHPartition ( 0 )
  if err != nil { return "",err }
  md.CXI.Header= cxi.Header

  // Comprova que és executable
  if cxi.Header.Type != citrus.NCCH_TYPE_CXI {
    return "",errors.New (
      "La imatge no conté un fitxer NCCH executable (CXI)" )
  }

  // Obté ExeFS
  exefs,err:= cxi.GetExeFS ()
  if err != nil { return "",err }
  if exefs == nil {
    return "",errors.New ( "No s'ha trobat una partició ExeFS" )
  }

  // Obté metadades del fitxer icon
  if icon_data,err:= _CXI_GetIconData ( exefs ); err != nil {
    return "",err
  } else if icon_data != nil {
    if err:= md.CXI.InitFromSMDH ( icon_data ); err != nil {
      return "",err
    }
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *N3DS) GetName() string {
  return "Imatge de cartutx de Nintendo 3DS"
}
func (self *N3DS) GetShortName() string { return "3DS" }
func (self *N3DS) IsImage() bool { return true }


func (self *N3DS) ParseMetadata(

  v         []view.StringPair,
  meta_data string,
  
) []view.StringPair {

  md:= _3DS_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[3DS] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  var kv *KeyValue
  
  // Metadades cartutx
  kv= &KeyValue{"Grandària cartutx (capçalera)",
    NumBytesToStr(uint64(md.Size))}
  v= append(v,kv)
  kv= &KeyValue{"Identificador (cartutx)",fmt.Sprintf("%016x",md.MediaID)}
  v= append(v,kv)
  kv= &KeyValue{"Versió (cartutx)",fmt.Sprintf("%04x",md.TitleVersion)}
  v= append(v,kv)
  kv= &KeyValue{"Revisió",fmt.Sprintf("%d",md.CardRevision)}
  v= append(v,kv)
  kv= &KeyValue{"Title ID (CVer)",fmt.Sprintf("%016x",md.TitleID)}
  v= append(v,kv)
  kv= &KeyValue{"Versió (CVer)",fmt.Sprintf("%04x",md.VersionCVer)}
  v= append(v,kv)

  // Metadades CXI
  v= md.CXI.ParseMetadata ( v )
  
  return v
  
} // end ParseMetadata
