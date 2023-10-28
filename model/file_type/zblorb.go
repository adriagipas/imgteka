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
 *  zblorb.go - Un format IFF per a emmagatzemar fitxers d'història Z.
 */

package file_type

import (
  "encoding/json"
  "encoding/xml"
  "errors"
  "fmt"
  "image"
  "image/jpeg"
  "image/png"
  "io"
  "log"
  "os"

  "github.com/adriagipas/imgteka/view"
)




/****************/
/* PART PRIVADA */
/****************/

type _ZBlorb_Metadata struct {

  ZCodeMetadata SFZ_Metadata
  StoryMetadata _ZBlorb_StoryMetadata
  
}

type _ZBlorb_Identification struct {
  Ifid   string `xml:"ifid"`
  Tuid   string `xml:"tuid"`
}

type _ZBlorb_Bibliographic struct {
  Title          string `xml:"title"`
  Author         string `xml:"author"`
  Language       string `xml:"language"`
  Headline       string `xml:"headline"`
  Firstpublished string `xml:"firstpublished"`
  Genre          string `xml:"genre"`
  Group          string `xml:"group"`
  Forgiveness    string `xml:"forgiveness"`
  Description    string `xml:"description"`
}

type _ZBlorb_Colophon struct {
  Generator        string `xml:"generator"`
  GeneratorVersion string `xml:"generatorversion"`
  Originated       string `xml:"originated"`
}

type _ZBlorb_Story struct {
  Identification _ZBlorb_Identification `xml:"identification"`
  Bibliographic  _ZBlorb_Bibliographic  `xml:"bibliographic"`
  Colophon       _ZBlorb_Colophon       `xml:"colophon"`
}

type _ZBlorb_StoryMetadata struct {
  Story _ZBlorb_Story `xml:"story"`
}

func _ZBlorb_ReadMetadata( md *_ZBlorb_Metadata, f io.Reader ) error {

  d:= xml.NewDecoder ( f )
  if err:= d.Decode ( &md.StoryMetadata ); err != nil {
    return err
  }
  
  return nil
  
} // end _ZBlorb_ReadMetadata


func _ZBlorb_ParseMetadata(

  v  []view.StringPair,
  md *_ZBlorb_StoryMetadata,
  
) []view.StringPair {

  var kv *KeyValue
  
  // Identification
  // --> ifid
  if md.Story.Identification.Ifid != "" {
    kv= &KeyValue{"IFID",md.Story.Identification.Ifid}
    v= append(v,kv)
  }
  // --> tuid
  if md.Story.Identification.Tuid != "" {
    kv= &KeyValue{"Interactive Fiction Database identifier",
      md.Story.Identification.Tuid}
    v= append(v,kv)
  }

  // Bibliographic
  // --> title / headline
  if md.Story.Bibliographic.Title != "" {
    var title string
    if md.Story.Bibliographic.Headline != "" {
      title= fmt.Sprintf ( "%s: %s", md.Story.Bibliographic.Title,
        md.Story.Bibliographic.Headline )
    } else {
      title= md.Story.Bibliographic.Title
    }
    kv= &KeyValue{"Títol",title}
    v= append(v,kv)
  }
  // --> author
  if md.Story.Bibliographic.Author != "" {
    kv= &KeyValue{"Autor",md.Story.Bibliographic.Author}
    v= append(v,kv)
  }
  // --> language
  if md.Story.Bibliographic.Language != "" {
    kv= &KeyValue{"Idioma",md.Story.Bibliographic.Language}
    v= append(v,kv)
  }
  // --> firstpublished
  if md.Story.Bibliographic.Firstpublished != "" {
    kv= &KeyValue{"Data publicació",md.Story.Bibliographic.Firstpublished}
    v= append(v,kv)
  }
  // --> genre
  if md.Story.Bibliographic.Genre != "" {
    kv= &KeyValue{"Gènere",md.Story.Bibliographic.Genre}
    v= append(v,kv)
  }
  // --> group
  if md.Story.Bibliographic.Group != "" {
    kv= &KeyValue{"Col·lecció",md.Story.Bibliographic.Group}
    v= append(v,kv)
  }
  // --> forgiveness
  if md.Story.Bibliographic.Forgiveness != "" {
    kv= &KeyValue{"Dificultat",md.Story.Bibliographic.Forgiveness}
    v= append(v,kv)
  }
  // --> description
  if md.Story.Bibliographic.Description != "" {
    kv= &KeyValue{"Descripció",md.Story.Bibliographic.Description}
    v= append(v,kv)
  }

  // Colophon
  // --> generator/generatorversion
  if md.Story.Colophon.Generator != "" {
    var compiler string
    if md.Story.Colophon.GeneratorVersion != "" {
      compiler= fmt.Sprintf ( "%s - %s", md.Story.Colophon.Generator,
        md.Story.Colophon.GeneratorVersion )
    } else {
      compiler= md.Story.Colophon.Generator
    }
    kv= &KeyValue{"Compilador",compiler}
    v= append(v,kv)
  }
  // --> originated
  if md.Story.Colophon.Originated != "" {
    kv= &KeyValue{"Data compilació",md.Story.Colophon.Originated}
    v= append(v,kv)
  }
  
  return v
  
} // end _ZBlorb_ParseMetadata


const (
  _ZBLORB_RESOURCE_PICT = 0
  _ZBLORB_RESOURCE_SND  = 1
  _ZBLORB_RESOURCE_DATA = 2
  _ZBLORB_RESOURCE_EXEC = 3
  _ZBLORB_RESOURCE_UNK  = 4
)

type _ZBlorb_Resource struct {
  Type   int
  Offset uint32
}

func _ZBlorb_ReadResourceIndex( r io.Reader ) ([]_ZBlorb_Resource,error) {

  // Llig índex
  data,err:= io.ReadAll ( r )
  if err != nil { return nil,err }

  // Nombre de recursos
  N:= uint32(
    (uint32(uint8(data[0]))<<24) |
      (uint32(uint8(data[1]))<<16) |
      (uint32(uint8(data[2]))<<8) |
      uint32(uint8(data[3])))
  if N == 0 {
    return nil,errors.New ( "No hi han recursos al fitxer" )
  }

  // Inicialitza recursos
  ret:= make([]_ZBlorb_Resource,N)
  for i:= 0; i < len(ret); i++ {
    ret[i].Type= _ZBLORB_RESOURCE_UNK
  }

  // Llig recursos
  data= data[4:]
  var num uint32
  for i:= uint32(0); i < N; i++ {

    // Nombre
    num= (uint32(uint8(data[4]))<<24) |
      (uint32(uint8(data[5]))<<16) |
      (uint32(uint8(data[6]))<<8) |
      uint32(uint8(data[7]))
    if num >= N {
      return nil,fmt.Errorf ( "el chunk RIdx està mal format: recurs"+
        " %d fora de rang", num )
    }

    // Offset
    ret[num].Offset= (uint32(uint8(data[8]))<<24) |
      (uint32(uint8(data[9]))<<16) |
      (uint32(uint8(data[10]))<<8) |
      uint32(uint8(data[11]))
    
    // Tipus
    if data[0]=='P' && data[1]=='i' && data[2]=='c' && data[3]=='t' {
      ret[num].Type= _ZBLORB_RESOURCE_PICT
    } else if data[0]=='S' && data[1]=='n' && data[2]=='d' && data[3]==' ' {
      ret[num].Type= _ZBLORB_RESOURCE_SND
    } else if data[0]=='D' && data[1]=='a' && data[2]=='t' && data[3]=='a' {
      ret[num].Type= _ZBLORB_RESOURCE_DATA
    } else if data[0]=='E' && data[1]=='x' && data[2]=='e' && data[3]=='c' {
      ret[num].Type= _ZBLORB_RESOURCE_EXEC
    }
    
    data= data[12:]
  }
  
  return ret,nil
  
} // end _ZBlorb_ReadResourceIndex


func _ZBlorb_GetImage ( fd *os.File, offset uint32 ) (image.Image,error) {

  // Aplica offset
  off:= int64(uint64(offset))
  new_off,err:= fd.Seek ( off, 0 )
  if err != nil { return nil,err }
  if new_off != off {
    return nil,errors.New ( "error inesperant mentre s'intentava "+
      "llegir la portada" )
  }

  // Llig capçalera chunk
  var buf_mem [8]byte
  buf:= buf_mem[:]
  n,err:= fd.Read ( buf )
  if err != nil { return nil,err }
  if n != len(buf) {
    return nil,errors.New ( "no s'ha pogut llegir la capçalera del "+
      "chunk que conté la portada" )
  }

  // Llig
  var ret image.Image= nil
  if buf[0]=='J' && buf[1]=='P' && buf[2]=='E' && buf[3]=='G' {
    ret,err= jpeg.Decode ( fd )
  } else if buf[0]=='P' && buf[1]=='N' && buf[2]=='G' && buf[3]==' ' {
    ret,err= png.Decode ( fd )
  } else {
    ret,err= nil,fmt.Errorf ( "chunk d'imatge desconegut '%c%c%c%c'",
      buf[0], buf[1], buf[2], buf[3] )
  }
  
  return ret,err
  
} // end _ZBlorb_GetImage




/****************/
/* PART PÚBLICA */
/****************/

type ZBlorb struct {
}


func (self *ZBlorb) GetImage(file_name string) (image.Image,error) {

  // Obri fitxer
  fd,err:= os.Open ( file_name )
  if err != nil { return nil,err }
  defer fd.Close ()
  
  // Llig IFF.
  iff,err:= newIFF ( fd )
  if err != nil { return nil,err }
  
  // Llig RIdx
  root,err:= iff.GetRootDirectory ()
  if err != nil { return nil,err }
  // --> Primer chunk ha de ser RIdx
  it,err:= root.Begin ()
  if err != nil { return nil,err }
  if it.GetType () != "RIdx" {
    return nil,errors.New ( "No és un fitxer en format Blorb: no "+
      "s'ha trobat el RIdx" )
  }
  // --> Llig índex
  sf,err:= it.GetFileReader ()
  if err != nil { return nil,err }
  index,err:= _ZBlorb_ReadResourceIndex ( sf )
  if err != nil { return nil,err }

  // Busca Frontispiece Chunk
  for ; err == nil && !it.End () && it.GetType () != "Fspc"; err= it.Next () {
  }
  if err != nil { return nil,err }
  if it.End () { return nil,errors.New ( "El fitxer no conté una portada" ) }

  // Obté recurs portada
  sf,err= it.GetFileReader ()
  if err != nil { return nil,err }
  data,err:= io.ReadAll ( sf )
  if err != nil { return nil,err }
  num_fspc:= (uint32(uint8(data[0]))<<24) |
    (uint32(uint8(data[1]))<<16) |
    (uint32(uint8(data[2]))<<8) |
    uint32(uint8(data[3]))
  if num_fspc > uint32(len(index)) ||
    index[num_fspc].Type != _ZBLORB_RESOURCE_PICT {
    return nil,fmt.Errorf ( "El recurs de portada %d no és vàlid", num_fspc  )
  }

  // Obté la imatge
  return _ZBlorb_GetImage ( fd, index[num_fspc].Offset )
  
} // end GetImage


func (self *ZBlorb) GetMetadata(fd *os.File) (string,error) {
  
  // Llig IFF.
  iff,err:= newIFF ( fd )
  if err != nil { return "",err }

  // Inicialitza metadades
  md:= _ZBlorb_Metadata{}
  
  // Comprova chunks
  root,err:= iff.GetRootDirectory()
  if err != nil { return "",err }
  // --> Primer chunk ha de ser RIdx
  it,err:= root.Begin ()
  if err != nil { return "",err }
  if it.GetType () != "RIdx" {
    return "",errors.New ( "No és un fitxer en format Blorb: no "+
      "s'ha trobat el RIdx" )
  }
  // --> Resta de chunks (fitxer història i metadades)
  zcod:= false
  for ; err == nil && !it.End (); err= it.Next () {
    switch typ:= it.GetType(); typ {
    case "ZCOD":
      zcod= true
      if sf,err:= it.GetFileReader (); err == nil {
        if err:= SFZ_ReadMetadata ( &md.ZCodeMetadata,
          sf, sf.Size() ); err != nil {
          return "",err
        }
      } else {
        return "",err
      }
    case "IFmd":
      if sf,err:= it.GetFileReader (); err == nil {
        if err:= _ZBlorb_ReadMetadata( &md, sf ); err != nil {
          return "",err
        }
      } else {
        return "",err
      }
    }
  }
  if err != nil { return "",err }
  if !zcod {
    return "",errors.New("No s'ha trobat cap fitxer d'història per "+
      "a la Màquina Z")
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }
  
  return string(b),nil
  
} // end GetMetadata


func (self *ZBlorb) GetName() string {
  return "Fitxer d'història Blorb (Màquina Z)"
} // end GetName
 

func (self *ZBlorb) GetShortName() string { return "ZBLORB" }
func (self *ZBlorb) IsImage() bool { return true } 


func (self *ZBlorb) ParseMetadata(
  
  v         []view.StringPair,
  meta_data string,
  
) []view.StringPair {

  // Parseja
  md:= _ZBlorb_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[ZBLORB] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }

  // Afegeix camps metadades
  v= SFZ_ParseMetadata ( v, &md.ZCodeMetadata )
  v= _ZBlorb_ParseMetadata ( v, &md.StoryMetadata )
  
  return v
  
} // end ParseMetadata
