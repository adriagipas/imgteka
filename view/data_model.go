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


type File interface {

  // Torna el nom del fitxer
  GetName() string

  // Torna una cadena curta que descriu el tipus de fitxar. Per
  // exemple Diquet 3 1/2
  GetType() string

  // Torna el MD5  
  GetMD5() string

  // Torna el SHA1
  GetSHA1() string

  // Torna la grandària en bytes
  GetSize() int64
  
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


type DataModel interface {

  // Torna la llista dels identificadors (long) de tots els objectes del
  // model.
  RootEntries() []int64

  // Torna una entrada del model
  GetEntry(id int64) Entry

  // Independentment del tipus de UI, vull que per a cada entrada es
  // motre al costat un quadrat amb l'identificador de la plataforma
  // (3 lletres??) i un color. Per tant cal que el model torne eixa
  // informació gràfica.
  GetPlatformHints(id int) (idname string,color color.Color)

  // Torna el nom de la plataforma.
  GetPlatformName(id int) string
  
  // Torna el fitxer indicat
  GetFile(id int64) File

  // Torna la informació d'una etiqueta. Nom i color.
  GetLabelInfo(id int) (name string,color color.Color)
  
}


///////////// PROVA INTERFÍCIE  ///////////////////////////////////////////////

type _FakeFile struct {
  name  string
  type_ string
  md5   string
  sha1  string
  size  int64
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

type _FakePlatformHint struct {
  idname string
  color  color.Color
}

type _FakeDataModel struct {
  entries []_FakeEntry
  files   []_FakeFile
  phints  []_FakePlatformHint
  pnames  []string
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

func (self *_FakeDataModel) GetPlatformHints(id int) (string,color.Color) {
  tmp:= self.phints[id]
  return tmp.idname,tmp.color
}

func (self *_FakeDataModel) GetPlatformName(id int) string {
  return self.pnames[id]
}

func (self *_FakeDataModel) GetLabelInfo(id int) (string,color.Color) {
  return self.lnames[id],self.lcolors[id]
}

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
        size:1048576},
        _FakeFile{name:"disc1.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:1474560},
        _FakeFile{name:"disc2.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:1474560},
        _FakeFile{name:"disc3.img",
          type_:"Disquet 3½",
          md5:"9EE8071A16E26613E6BACDC5056ACCC5",
          sha1:"5FCFB8EAA946F1C4968E5B27DF6476CB41C8D3D3",
          size:147456000},
        _FakeFile{name:"manual.pdf",
          type_:"Document PDF",
          md5:"B937BE5E84BBAA799FF1E029FD4246E4",
          sha1:"7E8942C85FB9B1BF303BC0C7C786BA8C2FB594D2",
          size:9224710},
      },
      phints: []_FakePlatformHint {
      _FakePlatformHint{idname:"MD ",color:color.RGBA{69,143,217,255}},
        _FakePlatformHint{idname:"DOS",color:color.RGBA{128,128,128,255}},
      },
      pnames: []string {
      "Sega MegaDrive",
        "MS-DOS",
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
