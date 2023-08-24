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
 *  data_model.go - Interfície que representa les dades que s'estan
 *                  mostrant en un moment donat en la llista.
 */

package view

import (
  "image"
  "image/color"
  "image/png"
  "os"
)


type StringPair interface {
  
  GetKey() string

  GetValue() string
}


type File interface {
  // Torna el nom del fitxer
  GetName() string

  // Torna una cadena curta que descriu el tipus de fitxar. Per
  // exemple Diquet 3 1/2
  GetType() string

  // Torna metadades associades a aquest fitxer
  GetMetadata() []StringPair
  
}


type Entry interface {

  // Torna el nom que es mostrarà en la interfície
  GetName() string

  // Torna identificador de la plataforma
  GetPlatformID() int

  // Torna els identificadors dels fitxers d'aquesta entrada. Són
  // globals, però diferents als de les entrades.
  GetFileIDs() []int64

  // Torna la imatge de la portada. Pot tornar nil
  GetCover() image.Image

  // Torna els identificadors de les etiquetes que té aquesta entrada.
  GetLabelIDs() []int
  
}


type Stats interface {

  // Torna el nombre d'entrades
  GetNumEntries() int64

  // Torna el nombre de fitxers
  GetNumFiles() int64
  
}


type Platform interface {

  // Torna el nom
  GetName() string

  // Torna el nom curt (màxim 3 lletres)
  GetShortName() string

  // Torna el color assignat a la plataforma
  GetColor() color.Color
  
}

type DataModel interface {

  // Torna la llista dels identificadors (long) de tots els objectes del
  // model.
  RootEntries() []int64

  // Torna una entrada del model
  GetEntry(id int64) Entry

  // Torna la plataforma.
  GetPlatform(id int) Platform
  
  // Torna el fitxer indicat
  GetFile(id int64) File

  // Torna la informació d'una etiqueta. Nom i color.
  GetLabelInfo(id int) (name string,color color.Color)

  // Obté estadístiques
  GetStats() Stats
  
}


///////////// PROVA INTERFÍCIE  ///////////////////////////////////////////////

type _FakeFile struct {
  name  string
  type_ string
  md5   string
  sha1  string
  size  int64
  md    []StringPair
}

func (self *_FakeFile) GetName() string {
  return self.name
}

func (self *_FakeFile) GetType() string {
  return self.type_
}

func (self *_FakeFile) GetMD5() string {
  return self.md5
}

func (self *_FakeFile) GetSHA1() string {
  return self.sha1
}

func (self *_FakeFile) GetSize() int64 {
  return self.size
}

func (self *_FakeFile) GetMetadata() []StringPair {
  return self.md
}

type _FakeEntry struct {
  name     string
  pid      int
  file_ids []int64
  cover    string // Pot ser nil
  labels   []int
}

func (self *_FakeEntry) GetName() string {
  return self.name
}

func (self *_FakeEntry) GetPlatformID() int {
  return self.pid
}

func (self *_FakeEntry) GetFileIDs() []int64 {
  return self.file_ids
}

func (self *_FakeEntry) GetCover() image.Image {
  f,err:= os.Open ( self.cover )
  if err != nil { return nil }
  ret,err:= png.Decode ( f )
  if err != nil { return nil }
  f.Close ()
  return ret
}

func (self *_FakeEntry) GetLabelIDs() []int {
  return self.labels
}

type _FakePlatform struct {
  name       string
  short_name string
  color      color.Color
}

func (self *_FakePlatform) GetName() string {
  return self.name
}

func (self *_FakePlatform) GetShortName() string {
  return self.short_name
}

func (self *_FakePlatform) GetColor() color.Color {
  return self.color
}

type _FakeDataModel struct {
  entries []_FakeEntry
  files   []_FakeFile
  plats   []_FakePlatform
  lnames  []string
  lcolors []color.Color
}

func (self *_FakeDataModel) RootEntries() []int64 {
  ret:= make([]int64,len(self.entries))
  for i:= 0; i < len(self.entries); i++ {
    ret[i]= int64(i)
  }
  return ret
}

func (self *_FakeDataModel) GetEntry(id int64) Entry {
  return &self.entries[id]
}

func (self *_FakeDataModel) GetFile(id int64) File {
  return &self.files[id]
}

func (self *_FakeDataModel) GetPlatform(id int) Platform {
  return &self.plats[id]
}

func (self *_FakeDataModel) GetLabelInfo(id int) (string,color.Color) {
  return self.lnames[id],self.lcolors[id]
}

func (self *_FakeDataModel) GetStats() Stats {
  ret:= _Stats{
    nentries : int64(len(self.entries)),
    nfiles : int64(len(self.files)),
  }
  return &ret
}

type _MetadataValue struct {
  property string
  value    string
}

func (self *_MetadataValue) GetKey() string { return self.property }
func (self *_MetadataValue) GetValue() string { return self.value }

type _Stats struct {
  nentries int64
  nfiles   int64
}

func (self *_Stats) GetNumEntries() int64 { return self.nentries }
func (self *_Stats) GetNumFiles() int64 { return self.nfiles }

func newFakeDataModel() *_FakeDataModel {
  ret:= _FakeDataModel {
    entries: []_FakeEntry {
      _FakeEntry{
        name:"Thunder Force IV",pid:0,
        file_ids: []int64{0},
        cover: "blo",
        labels: []int{3,4},
      },
        _FakeEntry{name:"Mortal Kombat",pid:1,
          file_ids: []int64 {1,2,3,4},
          cover: "/home/adria/COLJOCS/DOS/Mortal Kombat/screenshots/s1.png",
          labels: []int{2,4,0,1,3},
        },
      },
      files: []_FakeFile {
      _FakeFile{name:"Thunder Force IV (Europe).md",
        type_:"Sega Mega Drive ROM",
        md5:"9EE8071A16E26613E6BACDC5056ACCC5",
        sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
        size:1048576,
        md: []StringPair{
          &_MetadataValue{"md5","9EE8071A16E26613E6BACDC5056ACCC5"},
          &_MetadataValue{"sha1","5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3"},
          &_MetadataValue{"Grandària","1048576 B (1 MB)"},
          &_MetadataValue{"Firma i data","(C)T-18 1992.AUG"},
          &_MetadataValue{"Nom domèstic","THUNDER FORCE 4"},
          &_MetadataValue{"Tipus i número de serie","GM MK-1143 -50"},
        }},
        _FakeFile{name:"disc1.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:1474560,
          md: []StringPair{
            &_MetadataValue{"md5","9EE8071A16E26613E6BACDC5056ACCC5"},
            &_MetadataValue{"sha1","5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3"},
            &_MetadataValue{"size","1474560 B (1.4 MB)"},
          }},
        _FakeFile{name:"disc2.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:1474560,
          md: []StringPair{
            &_MetadataValue{"md5","9EE8071A16E26613E6BACDC5056ACCC5"},
            &_MetadataValue{"sha1","5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3"},
            &_MetadataValue{"size","1474560 B (1.4 MB)"},
          }},
        _FakeFile{name:"disc3.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:147456000,
          md: []StringPair{
            &_MetadataValue{"md5","9EE8071A16E26613E6BACDC5056ACCC5"},
            &_MetadataValue{"sha1","5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3"},
            &_MetadataValue{"size","147456000 B (140.6 MB)"},
          }},
        _FakeFile{name:"manual.pdf",
          type_:"Document PDF",
          md5:"B937BE5E84BBAA799FF1E029FD4246E4",
          sha1:"7E8942C85FB9B1BF303BC0C7C786BA8C2FB594D2",
          size:9224710,
          md: []StringPair{
            &_MetadataValue{"md5","B937BE5E84BBAA799FF1E029FD4246E4"},
            &_MetadataValue{"sha1","7E8942C85FB9B1BF303BC0C7C786BA8C2FB594D2"},
            &_MetadataValue{"size","9224710 B (8.8 MB)"},
          }},
      },
      plats: []_FakePlatform {
      _FakePlatform{
        name:"Sega MegaDrive",
        short_name:"MD ",
        color:color.RGBA{69,143,217,255},
      },
        _FakePlatform{
          name:"MS-DOS",
          short_name:"DOS",
          color:color.RGBA{128,128,128,255}},
      },
      lnames: []string {
      "Aventura",
        "Estratègia",
        "Lluita",
        "Tirs",
        "Perspectiva lateral",
      },
      lcolors: []color.Color{
        color.RGBA{255,200,200,255},
        color.RGBA{255,200,200,255},
        color.RGBA{255,200,200,255},
        color.RGBA{255,200,200,255},
        color.RGBA{200,200,255,255},
      },
    }
  return &ret
}
