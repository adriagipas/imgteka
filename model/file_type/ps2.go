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
 *  ps2.go - Tipus de fitxer imatge DVD de PlayStation 2.
 */

package file_type

import (
  "bufio"
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "log"
  "strings"

  "github.com/adriagipas/imgcp/cdread"
  "github.com/adriagipas/imgteka/view"
)




/****************/
/* PART PRIVADA */
/****************/

type _PS2_Metadata struct {

  Cd           _CD_Metadata // Metadades a nivell d'imatge de CD
  Iso          _ISO_Metadata // Metadades a nivell d'ISO
  Id           string
  TitleVersion string
  VideoMode    string
  
}


func (self *_PS2_Metadata) parseSystemCnf( iso *cdread.ISO ) error {

  // Busca el fitxer SYSTEM.CNF
  root,err:= iso.Root ()
  if err != nil { return err }
  it,err:= root.Begin ()
  for ; err == nil && !it.End () &&
    it.Id () != "SYSTEM.CNF" && it.Id () != "SYSTEM.CNF;1"; err= it.Next () {
  }
  if err != nil {
    return err
  } else if it.End () {
    return errors.New ( "no s'ha trobat el fitxer SYSTEM.CNF" )
  }

  // Recorre el fitxer línia a línia
  f,err:= it.GetFileReader ()
  if err != nil { return err }
  defer f.Close ()
  scan:= bufio.NewScanner ( bufio.NewReader ( f ) )
  var line,id string
  wrong_format:= fmt.Errorf ( "format SYSTEM.CNF incorrecte: '%s'", line )
  var toks []string
  ok:= false
  for ; scan.Scan (); {
    line= scan.Text ()
    toks= strings.Split ( line, "=" )
    if len(toks) != 2 { return wrong_format }
    switch strings.ToUpper ( strings.TrimSpace ( toks[0] ) ) {
    case "BOOT2":
      toks= strings.Split ( toks[1], "cdrom0:" )
      if len(toks) != 2 { return wrong_format }
      toks= strings.Split ( toks[1], ";" )
      id= toks[0]
      id= strings.ReplaceAll ( id, ".", "" )
      id= strings.ReplaceAll ( id, "_", "-" )
      id= strings.ReplaceAll ( id, "\\", "" )
      self.Id= strings.ToUpper ( id )
      if len(self.Id) > 0 {
        ok= true
      }
    case "VER":
      self.TitleVersion= strings.TrimSpace ( toks[1] )
    case "VMODE":
      self.VideoMode= strings.TrimSpace ( toks[1] )
    }
  }
  if !ok {
    errors.New ( "No s'ha trobat el camp BOOT2 en SYSTEM.CNF" )
  }
  
  return nil
  
} // end parseSystemCnf


func (self *_PS2_Metadata) Init( cd cdread.CD, iso *cdread.ISO ) error {

  if err:= self.parseSystemCnf ( iso ); err != nil {
    return err
  }
  
  return nil
  
} // end Init




/****************/
/* PART PÚBLICA */
/****************/

type PS2 struct {
}


func (self *PS2) GetImage( file_name string ) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una imatge de DVD de PlayStation 2" )
} // end GetImage


func (self *PS2) GetMetadata(file_name string) (string,error) {

  // Intenta obrir el CD
  cd,err:= cdread.Open ( file_name )
  if err != nil { return "",err }

  // Inicialitza metadades.
  md:= _PS2_Metadata{}
  md.Cd.Init ( cd )

  // Llig track ISO.
  iso,err:= cdread.ReadISO ( cd, 0, 0 )
  if err != nil { return "",err }
  md.Iso.Init ( iso )

  // Metadades PSX.
  if err:= md.Init ( cd, iso ); err != nil {
    return "",err
  }
  
  // Converteix  a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *PS2) GetName() string { return "DVD de PlayStation 2" }
func (self *PS2) GetShortName() string { return "PS2" }
func (self *PS2) IsImage() bool { return false }


func (self *PS2) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _PS2_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[PS2] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Metadades PS2
  var kv *KeyValue
  kv= &KeyValue{"Identificador",md.Id}
  v= append(v,kv)
  if md.TitleVersion != "" {
    kv= &KeyValue{"Versió",md.TitleVersion}
    v= append(v,kv)
  }
  if md.VideoMode != "" {
    kv= &KeyValue{"Mode vídeo",md.VideoMode}
    v= append(v,kv)
  }
  
  // CD i Iso
  v= md.Cd.Parse ( v )
  v= md.Iso.Parse ( v )
  
  return v
  
} // end ParseMetadata
