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
 *  psp.go - Tipus de fitxer imatge UMD de PlayStation Portable.
 */

package file_type

import (
  "bufio"
  "encoding/json"
  "errors"
  "fmt"
  "image"
  "image/png"
  "log"
  "strings"
  
  "github.com/adriagipas/imgcp/cdread"
  "github.com/adriagipas/imgteka/view"
)




/****************/
/* PART PRIVADA */
/****************/

type _PSP_Metadata struct {

  Cd     _CD_Metadata // Metadades a nivell d'imatge de CD
  Iso    _ISO_Metadata // Metadades a nivell d'ISO
  Id     string
  
}


func (self *_PSP_Metadata) readId( iso *cdread.ISO ) error {

  // Busca la carpeta PSP_GAME dir i UMD_DATA.BIN
  root,err:= iso.Root ()
  if err != nil { return err }
  var found_PSP_GAME bool= false
  var found_UMD_DATA_BIN bool= false
  it,err:= root.Begin ()
  for ; err == nil && !it.End () &&
    !(found_PSP_GAME && found_UMD_DATA_BIN); err= it.Next () {
    switch it.Id () {
      
    case "PSP_GAME":
      if (it.Flags()&cdread.FILE_FLAGS_DIRECTORY)!=0 {
        found_PSP_GAME= true
      } else {
        return errors.New("PSP_GAME no és un directori")
      }

    case "UMD_DATA.BIN":
      found_UMD_DATA_BIN= true
      
      // Llig la primera línia.
      f,err:= it.GetFileReader ()
      if err != nil { return err }
      defer f.Close ()
      scan:= bufio.NewScanner ( bufio.NewReader ( f ) )
      if !scan.Scan () {
        return errors.New ( "error mentre es processava UMD_DATA.BIN" )
      }
      line:= scan.Text ()

      // Processa línia
      wrong_format:= fmt.Errorf ( "format UMD_DATA.BIN incorrecte: '%s'", line )
      toks:= strings.Split ( line, "|" )
      if len(toks) < 2 { return wrong_format }
      id:= toks[0]
      self.Id= strings.ToUpper ( id )
      
    }
  }
  if err != nil {
    return err
  } else if it.End () {
    if !found_UMD_DATA_BIN {
      return errors.New ( "no s'ha trobat el fitxer UMD_DATA.BIN" )
    } else if !found_PSP_GAME {
      return errors.New ( "no s'ha trobat la carpeta PSP_GAME" )
    }
  }
  
  return nil
  
} // end readId


func (self *_PSP_Metadata) Init( cd cdread.CD, iso *cdread.ISO ) error {

  if err:= self.readId ( iso ); err != nil {
    return err
  }
  
  return nil
  
} // end Init




/****************/
/* PART PÚBLICA */
/****************/

type PSP struct {
}


func (self *PSP) GetImage( file_name string ) (image.Image,error) {

  // Intenta obrir el CD
  cd,err:= cdread.Open ( file_name )
  if err != nil { return nil,err }

  // Llig track ISO.
  iso,err:= cdread.ReadISO ( cd, 0, 0 )
  if err != nil { return nil,err }
  
  // Busca la carpeta PSP_GAME
  root,err:= iso.Root ()
  if err != nil { return nil,err }
  it,err:= root.Begin()
  for ; err == nil && !it.End () && it.Id () != "PSP_GAME"; err= it.Next () {
  }
  if err != nil {
    return nil,err
  } else if it.End () {
    return nil,errors.New("No s'ha trobat la carpeta PSP_GAME")
  }

  // Busca ICON0.PNG
  dir,err:= it.GetDirectory ()
  if err != nil { return nil,err }
  it,err= dir.Begin ()
  for ; err == nil && !it.End () && it.Id () != "ICON0.PNG"; err= it.Next () {
  }
  if err != nil {
    return nil,err
  } else if it.End () {
    return nil,errors.New("No s'ha trobat el fitxer ICON0.PNG")
  }

  // Carrega la icona.
  f,err:= it.GetFileReader ()
  if err != nil { return nil,err }
  defer f.Close ()

  return png.Decode ( f )
  
} // end GetImage


func (self *PSP) GetMetadata(file_name string) (string,error) {

  // Intenta obrir el CD
  cd,err:= cdread.Open ( file_name )
  if err != nil { return "",err }

  // Inicialitza metadades.
  md:= _PSP_Metadata{}
  md.Cd.Init ( cd )
  
  // Llig track ISO.
  iso,err:= cdread.ReadISO ( cd, 0, 0 )
  if err != nil { return "",err }
  md.Iso.Init ( iso )

  // Metadades PSP.
  if err:= md.Init ( cd, iso ); err != nil {
    return "",err
  }

  // Converteix  a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil
  
} // end GetMetadata


func (self *PSP) GetName() string { return "UMD de PlayStation Portable" }
func (self *PSP) GetShortName() string { return "PSP" }
func (self *PSP) IsImage() bool { return true }


func (self *PSP) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _PSP_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[PSP] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Metadades PSP
  var kv *KeyValue
  kv= &KeyValue{"Identificador",md.Id}
  v= append(v,kv)
  
  // CD i Iso
  v= md.Cd.Parse ( v )
  v= md.Iso.Parse ( v )
  
  return v
  
} // end ParseMetadata
