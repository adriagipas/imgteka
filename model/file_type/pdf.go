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
 *  pdf.go - Tipus de fitxer PDF.
 */

package file_type

import (
  "encoding/json"
  "fmt"
  "image"
  "log"
  "time"

  "github.com/adriagipas/imgteka/view"
  "seehuhn.de/go/pdf"
  
)




type _PDF_Metadata struct {

  Version      string
  Title        string
  Author       string
  Subject      string
  Keywords     string
  Creator      string
  Producer     string
  CreationDate string
  ModDate      string
  
}


type PDF struct {
}


func (self *PDF) GetImage( file_name string) (image.Image,error) {
  return nil,fmt.Errorf (
    "No es pot interpretar com una imatge un fitxer de tipus PDF" )
} // end GetImage


func (self *PDF) GetMetadata(file_name string) (string,error) {

  // Llig.
  r,err:= pdf.Open ( file_name, nil )
  if err != nil { return "",err }
  defer r.Close()

  // Crea metadades
  md:= _PDF_Metadata{}
  meta:= r.GetMeta()
  md.Version= meta.Version.String()

  // Metadades autoria
  if meta.Info != nil {
    md.Title= meta.Info.Title
    md.Author= meta.Info.Author
    md.Subject= meta.Info.Subject
    md.Keywords= meta.Info.Keywords
    md.Creator= meta.Info.Creator
    md.Producer= meta.Info.Producer
    md.CreationDate= meta.Info.CreationDate.Format(time.DateTime)
    md.ModDate= meta.Info.ModDate.Format(time.DateTime)
  }
  
  // Converteix a json
  b,err:= json.Marshal ( md )
  if err != nil { return "",err }

  return string(b),nil

} // end GetMetadata


func (self *PDF) GetName() string { return "Document PDF" }
func (self *PDF) GetShortName() string { return "PDF" }
func (self *PDF) IsImage() bool { return false }


func (self *PDF) ParseMetadata(

  v         []view.StringPair,
  meta_data string,

) []view.StringPair {

  // Parseja
  md:= _PDF_Metadata{}
  if err:= json.Unmarshal ( []byte(meta_data), &md ); err != nil {
    log.Printf ( "[PDF] no s'ha pogut parsejar '%s': %s", meta_data, err )
    return v
  }
  
  // Versió
  kv:= KeyValue{"Versió PDF",md.Version}
  v= append(v,&kv)

  // Títol
  if md.Title != "" {
    kv:= KeyValue{"Títol",md.Title}
    v= append(v,&kv)
  }

  // Author
  if md.Title != "" {
    kv:= KeyValue{"Autor",md.Author}
    v= append(v,&kv)
  }

  // Subject
  if md.Title != "" {
    kv:= KeyValue{"Tema",md.Subject}
    v= append(v,&kv)
  }

  // Paraules claus
  if md.Keywords != "" {
    kv:= KeyValue{"Paraules claus",md.Keywords}
    v= append(v,&kv)
  }

  // Creator
  if md.Creator != "" {
    kv:= KeyValue{"Creat amb",md.Creator}
    v= append(v,&kv)
  }

  // Producer
  if md.Producer != "" {
    kv:= KeyValue{"Convertit amb",md.Producer}
    v= append(v,&kv)
  }

  // CreationDate
  if md.CreationDate != "" {
    kv:= KeyValue{"Data creació",md.CreationDate}
    v= append(v,&kv)
  }

  // ModDate
  if md.ModDate != "" {
    kv:= KeyValue{"Data modificació",md.ModDate}
    v= append(v,&kv)
  }
  
  return v
  
} // end ParseMetadata
