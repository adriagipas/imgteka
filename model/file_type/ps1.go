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
 *  ps1.go - Tipus de fitxer imatge CD de PlayStation 1.
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

const (
  _PS1_REGION_EUROPE  = 0
  _PS1_REGION_JAPAN   = 1
  _PS1_REGION_AMERICA = 2
)

type _PS1_Metadata struct {

  Cd     _CD_Metadata // Metadades a nivell d'imatge de CD
  Iso    _ISO_Metadata // Metadades a nivell d'ISO
  Region int
  Id     string
  
}


func (self *_PS1_Metadata) readRegion( cd cdread.CD ) error {

  // Obté reader track 0.
  tr,err:= cd.TrackReader ( 0, 0, 0 )
  if err != nil { return err }
  defer tr.Close ()

  // Llig llicència del sector 4.
  var buf [70]byte
  if err:= tr.Seek ( 4 ); err != nil { return err }
  if nr,err:= tr.Read ( buf[:] ); err != nil {
    return err
  } else if nr != len(buf) {
    return errors.New ( "no s'ha pogut llegir la llicència" )
  }

  // Comprova llicència
  tmp:= string(buf[:60])
  if tmp != "          Licensed  by          Sony Computer Entertainment " {
    return fmt.Errorf ( "la imatge no conté la llicència esperada: '%s'", tmp )
  }
  region:= string(buf[60:])
  if region == "Euro pe   " {
    self.Region= _PS1_REGION_EUROPE
  } else if region == "Amer  ica " {
    self.Region= _PS1_REGION_AMERICA
  } else if region[:4] == "Inc." {
    self.Region= _PS1_REGION_JAPAN
  } else {
    return fmt.Errorf ( "regió de CD de PlayStation desconegut: '%s'", region )
  }
  
  return nil
  
} // end readRegion


func (self *_PS1_Metadata) readId( iso *cdread.ISO ) error {

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

  // Llig la primera línia.
  f,err:= it.GetFileReader ()
  if err != nil { return err }
  defer f.Close ()
  scan:= bufio.NewScanner ( bufio.NewReader ( f ) )
  if !scan.Scan () {
    return errors.New ( "error mentre es processava SYSTEM.CNF" )
  }
  line:= scan.Text ()

  // Processa línia
  wrong_format:= fmt.Errorf ( "format SYSTEM.CNF incorrecte: '%s'", line )
  toks:= strings.Split ( line, "=" )
  if len(toks) != 2 { return wrong_format }
  toks= strings.Split ( toks[1], "cdrom:" )
  if len(toks) != 2 { return wrong_format }
  toks= strings.Split ( toks[1], ";" )
  id:= toks[0]
  id= strings.ReplaceAll ( id, ".", "" )
  id= strings.ReplaceAll ( id, "_", "-" )
  id= strings.ReplaceAll ( id, "\\", "" )
  self.Id= id
  
  return nil
  
} // end readId


func (self *_PS1_Metadata) Init( cd cdread.CD, iso *cdread.ISO ) error {

  if err:= self.readRegion ( cd ); err != nil {
    return err
  }
  if err:= self.readId ( iso ); err != nil {
    return err
  }
  
  return nil
  
} // end Init




/****************/
/* PART PÚBLICA */
/****************/

// Fitxers auxiliars que s'ha de guardar en el mateix lloc.
type PS1_Aux struct {
}


func (self *PS1_Aux) GetImage( file_name string ) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una imatge de CD de PlayStation" )
} // end GetImage


func (self *PS1_Aux) GetMetadata(file_name string) (string,error) {
  return "",nil
} // end GetMetadata


func (self *PS1_Aux) GetName() string {
  return "Fitxer auxiliar d'imatge de CD de PlayStation" }
func (self *PS1_Aux) GetShortName() string { return "PS1" }
func (self *PS1_Aux) IsImage() bool { return false }


func (self *PS1_Aux) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {
  return v
} // end ParseMetadata


type PS1 struct {
}


func (self *PS1) GetImage( file_name string ) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge una imatge de CD de PlayStation" )
} // end GetImage


func (self *PS1) GetMetadata(file_name string) (string,error) {

  // Intenta obrir el CD
  cd,err:= cdread.Open ( file_name )
  if err != nil { return "",err }

  // Inicialitza metadades.
  md:= _PS1_Metadata{}
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


func (self *PS1) GetName() string { return "CD de PlayStation" }
func (self *PS1) GetShortName() string { return "PS1" }
func (self *PS1) IsImage() bool { return false }


func (self *PS1) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _PS1_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[MD] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Metadades PS1
  var kv *KeyValue
  kv= &KeyValue{"Identificador",md.Id}
  v= append(v,kv)
  region:= ""
  switch md.Region {
  case _PS1_REGION_EUROPE:
    region= "Europa"
  case _PS1_REGION_JAPAN:
    region= "Japó"
  case _PS1_REGION_AMERICA:
    region= "Amèrica"
  }
  kv= &KeyValue{"Regió",region}
  v= append(v,kv)
  
  // CD i Iso
  v= md.Cd.Parse ( v )
  v= md.Iso.Parse ( v )
  
  return v
  
} // end ParseMetadata
