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
)


type ProgressBar interface {

  Close()
  
  // Fixa un missatge i una fracció de completesa [0..1].
  Set(message string,fraction float32)
  
}


type StringPair interface {
  
  GetKey() string

  GetValue() string
  
}


type File interface {

  // Torna la imatge que representa el fitxer, o nil si no en té (o no
  // es pot carregar). S'indica l'ample i alt màxims.
  GetImage(max_wh int) image.Image
  
  // Torna el nom del fitxer
  GetName() string

  // Torna l'identificador del tipus
  GetTypeID() int

  // Torna metadades associades a aquest fitxer
  GetMetadata() []StringPair

  // "Executa" el fitxer.
  Run() error
  
}


type Entry interface {

  // Torna el nom que es mostrarà en la interfície
  GetName() string

  // Torna identificador de la plataforma
  GetPlatformID() int

  // Torna els identificadors dels fitxers d'aquesta entrada. Són
  // globals, però diferents als de les entrades.
  GetFileIDs() []int64

  // Torna els identificadors dels fitxers d'aquesta entrada que poden
  // ser tractats com a imatges.
  GetImageFileIDs() []int64

  // Torna la imatge de la portada. Pot tornar nil. S'indica
  // l'amplària i altura màxima.
  GetCover(max_wh int) image.Image
  
  // Torna els identificadors de les etiquetes que té aquesta entrada.
  GetLabelIDs() []int

  // Torna els identificadors de les etiquetes no emprades per aquesta
  // entrada.
  GetUnusedLabelIDs() []int

  // Afegeix (i crea) un nou fitxer.
  // path -> Path fitxer
  // name -> Nom amb el que volem registrar el fitxer
  // file_type -> Identificador tipus fitxer
  // create_pb -> Funció que crea i mostra una barra de progrés
  AddFile(path string,name string,file_type int,
    create_pb func() ProgressBar) error
  
  // Afegeix una nova etiqueta
  AddLabel(id int) error

  // Elimina el fitxer de l'entrada
  RemoveFile(id int64) error
  
  // Elimina una etiqueta de l'entrada
  RemoveLabel(id int) error

  // Fixa l'identificador del fitxer que serà la portada. -1 indica
  // que no té portada.
  SetCoverFileID(id int64) error
  
  // Actualitza el nom de l'entrada.
  UpdateName(name string) error

  // Actualitza el nom d'un fitxer
  UpdateFileName(id int64,name string) error
  
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

  // Torna el nombre d'entrades
  GetNumEntries() int64
  
  // Actualitza els atributs bàsics d'una plataforma. (No es pot
  // modificar el nom curt)
  Update(name string,c color.Color) error
  
}


type Label interface {

  // Torna el nom
  GetName() string

  // Torna el color de l'etiqueta
  GetColor() color.Color

  // Torna el nombre d'entrades
  GetNumEntries() int64

  // Actualitza els atributs bàsics d'una etiqueta.
  Update(name string,c color.Color) error
  
}


type DataModel interface {

  // Torna la llista dels identificadors (long) de tots els objectes del
  // model.
  RootEntries() []int64

  // Torna una entrada del model
  GetEntry(id int64) Entry

  // Torna els identificadors de les plataformes
  GetPlatformIDs() []int
  
  // Torna la plataforma.
  GetPlatform(id int) Platform
  
  // Torna el fitxer indicat
  GetFile(id int64) File

  // Torna els identificadors de les etiquetes
  GetLabelIDs() []int
  
  // Torna l'etiqueta
  GetLabel(id int) Label
  
  // Obté estadístiques
  GetStats() Stats

  // Obté els identificadors dels tipus de fitxer
  GetFileTypeIDs() []int

  // Obté el nom d'un tipus
  GetFileTypeName(id int) string

  // Obté el comandament per un tipus de fitxer. Si no en té,
  // aleshores cadena buida.
  GetFileTypeCommand(id int) string

  // Fixa comandament per a un tipus de fitxer. Cadena buida elimina
  // el comandament.
  SetFileTypeCommand(id int,command string)
  
  // Afegeix una nova plataforma
  AddPlatform(short_name string,name string,c color.Color) error

  // Afegeix una nova etiqueta
  AddLabel(name string,c color.Color) error

  // Afegeix una nova entrada
  AddEntry(name string,platform_id int) error

  // Elimina una plataforma
  RemovePlatform(id int) error

  // Elimina una entrada
  RemoveEntry(id int64) error
  
  // Elimina una etiqueta
  RemoveLabel(id int) error
  
}
