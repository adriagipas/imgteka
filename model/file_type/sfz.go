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
 *  sfz.go - Tipus de fitxer Story File Z-machine.
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

type SFZ_Metadata struct {

  Version    int8
  Release    uint16
  SerialCode string // Cadena buida vol dir que no en té
  
}


func _SFZ_CheckDataSize( data []byte ) error {

  version:= int8(data[0])
  size:= len(data)
  
  // Grandària capçalera.
  if size < 64 {
    return errors.New ( "El fitxer no conté capçalera" )
  }

  // V1/2/3 - 128KB
  if version == 1 || version == 2 || version == 3 {
    if size > 128*1024 {
      return fmt.Errorf ( "La grandària (%d) supera la grandària"+
        " màxima (128KB) per a un fitxer d'història de versió %d",
        size, version )
    }
    
    // V4/5 - 256KB
  } else if version == 4 || version == 5 {
    if size > 256*1024 {
      return fmt.Errorf ( "La grandària (%d) supera la grandària"+
        " màxima (256KB) per a un fitxer d'història de versió %d",
        size, version )
    }

    // V6/7/8 - 512KB
  } else if version == 6 || version == 7 || version == 8 {
    if size > 512*1024 {
      return fmt.Errorf ( "La grandària (%d) supera la grandària"+
        " màxima (512KB) per a un fitxer d'història de versió %d",
        size, version )
    }
    
    // Versió desconeguda
  } else {
    return fmt.Errorf ( "Versió desconeguda: %d", version )
  }
  
  return nil
  
} // end _SFZ_CheckDataSize


func SFZ_ReadMetadata( md *SFZ_Metadata, fd *os.File, size int64 ) error {


  // Comprova grandària no siga 0.
  if ( size == 0 ) {
    return errors.New ( "Grandària de fitxer 0" )
  }
  
  // Llig el fitxer d'història.
  data:= make([]byte,size)
  n,err:= fd.Read ( data )
  if err != nil { return err }
  if int64(n) != size { return errors.New ( "Error llegint les dades" ) }
  
  // Comprovacions de grandària.
  if err:= _SFZ_CheckDataSize ( data ); err != nil {
    return err
  }

  // Llig metadades bàsiques
  // --> Versió
  md.Version= int8(data[0])
  // --> Release number
  md.Release= (uint16(data[2])<<8) | uint16(data[3])
  // --> SerialCode
  if md.Version == 1 {
    md.SerialCode= ""
  } else {
    md.SerialCode= string(data[0x12:0x12+6])
  }
  
  return nil
  
} // end SFZ_ReadMetadata


func SFZ_ParseMetadata(
  
  v  []view.StringPair,
  md *SFZ_Metadata,
  
) []view.StringPair {

  var kv *KeyValue

  // Versió
  kv= &KeyValue{"Versió de la Màquina-Z",fmt.Sprintf ( "%d", md.Version )}
  v= append(v,kv)

  // Release
  kv= &KeyValue{"Versió",fmt.Sprintf ( "%d", md.Release )}
  v= append(v,kv)

  // Serial number
  if md.SerialCode != "" {
    kv= &KeyValue{"Codi de sèrie",md.SerialCode}
    v= append(v,kv)
  }
  
  return v
  
} // end SFZ_ParseMetadata




/****************/
/* PART PÚBLICA */
/****************/

type SFZ struct {
}


func (self *SFZ) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge un fitxer"+
      " d'història de la màquina Z" )
} // end GetImage


func (self *SFZ) GetMetadata(fd *os.File) (string,error) {

  // Rebobina
  if _,err:= fd.Seek ( 0, 0 ); err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }

  // Obté grandària
  info,err:= fd.Stat ()
  if err != nil {
    return "",fmt.Errorf ( "No s'ha pogut obtindre les metadades: %s", err )
  }
  size:= info.Size ()

  // Llig metadades
  md:= SFZ_Metadata{}
  if err:= SFZ_ReadMetadata ( &md, fd, size ); err != nil {
    return "",err
  }

  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *SFZ) GetName() string { return "Fitxer d'història de la Màquina Z" }
func (self *SFZ) GetShortName() string { return "SFZ" }
func (self *SFZ) IsImage() bool { return false }


func (self *SFZ) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= SFZ_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[SFZ] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  // Afegeix camps metadades
  v= SFZ_ParseMetadata ( v, &md )
  
  return v
  
} // end ParseMetadata
