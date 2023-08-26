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
 *  platforms.go - Gestió de les plataformes. Manté una "cache" de les
 *                 plataformes i es comunica amb la base de dades.
 */

package model

import (
  "errors"
  "fmt"
  "image/color"
  "log"
  "strings"
)




/****************/
/* PART PRIVADA */
/****************/

// Utilitzat per Database per afegir una platforma
func (self *Platforms) add(
  id         int,
  short_name string,
  name       string,
  r,g,b      uint8,
) {

  // Reformata short_name
  for i:= 0; i < 3-len(short_name); i++ {
    short_name= short_name+" "
  }
  
  // Afegeix
  self.ids= append ( self.ids, id )
  self.v[id]= &Platform{
    plats : self,
    id : id,
    short_name : short_name,
    name : name,
    color: color.RGBA{r,g,b,255},
  }
  
} // end add


func (self *Platforms) reset() {

  // Reseteja
  self.ids= self.ids[:0]
  self.v= make(map[int]*Platform)

  // Carrega
  if err:= self.db.LoadPlatforms ( self ); err != nil {
    log.Fatal ( err )
  }
  
} // end load




/****************/
/* PART PÚBLICA */
/****************/

type Platforms struct {
  db  *Database
  ids []int
  v   map[int]*Platform
}


func NewPlatforms ( db *Database ) *Platforms {

  ret:= Platforms{
    db  : db,
    ids : nil,
    v   : nil,
  }
  ret.reset ()
  
  return &ret
  
} // end NewPlatforms


func (self *Platforms) GetIDs() []int { return self.ids }
func (self *Platforms) GetPlatform( id int ) *Platform { return self.v[id] }


func (self *Platforms) Add(
  short_name string,
  name       string,
  c          color.Color,
) error {

  // Processa nom curt
  short_name= strings.TrimSpace ( short_name )
  short_name= strings.ToUpper ( short_name )
  if len(short_name) == 0 {
    return errors.New ( "No s'ha especificat un nom curt" )
  }
  if len(short_name) > 3 {
    return fmt.Errorf ( "El nom curt no pot superar els 3 caràcters: '%s'",
      short_name )
  }

  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }

  // Color
  r,g,b,_:= c.RGBA ()
  r8:= uint8((float32(r)/65535.0)*255.0 + 0.5)
  g8:= uint8((float32(g)/65535.0)*255.0 + 0.5)
  b8:= uint8((float32(b)/65535.0)*255.0 + 0.5)

  // Registra i reseteja
  if err:= self.db.RegisterPlatform (
    short_name, name, r8, g8, b8 ); err != nil {
    return fmt.Errorf ( "No s'ha pogut registrar la nova plataforma: %s", err )
  }
  self.reset ()
  
  return nil
  
} // end Add


func (self *Platforms) GetNumEntriesPlatform( id int ) int64 {

  ret,err:= self.db.GetPlatformNumEntries ( id )
  if err != nil { log.Fatal ( err ) }

  return ret
  
} // end GetNumEntriesPlatform


func (self *Platforms) Remove( id int ) error {

  // Comprova que és una plataforma que no s'utilitza
  plat:= self.v[id]
  if plat == nil {
    return fmt.Errorf ( "La plataforma indicada (%d) no existeix", id )
  }
  if plat.GetNumEntries() > 0 {
    return errors.New ( "No es pot esborrar la plataforma perquè està en ús" )
  }

  // Elimina i reseteja.
  if err:= self.db.DeletePlatform ( id ); err != nil {
    return fmt.Errorf ( "No s'ha pogut esborrar la plataforma: %s", err )
  }
  self.reset ()
  
  return nil
  
} // end Remove

func (self *Platforms) UpdatePlatform(
  id    int,
  name  string,
  r,g,b uint8,
) error {

  if err:= self.db.UpdatePlatform ( id, name, r, g, b ); err != nil {
    return fmt.Errorf ( "No s'ha pogut actualitzar la plataforma: %s", err )
  }

  return nil
  
} // end UpdatePlatform


type Platform struct {
  plats      *Platforms
  id         int // Identificador intern
  short_name string
  name       string
  color      color.Color
}


func (self *Platform) GetName() string { return self.name }
func (self *Platform) GetShortName() string { return self.short_name }
func (self *Platform) GetColor() color.Color { return self.color }


func (self *Platform) GetNumEntries() int64 {
  return self.plats.GetNumEntriesPlatform ( self.id )
} // end GetNumEntries


func (self *Platform) Update( name string, c color.Color ) error {

  // Processa nom
  name= strings.TrimSpace ( name )
  if len(name) == 0 {
    return errors.New ( "No s'ha especificat un nom" )
  }

  // Color
  r,g,b,_:= c.RGBA ()
  r8:= uint8((float32(r)/65535.0)*255.0 + 0.5)
  g8:= uint8((float32(g)/65535.0)*255.0 + 0.5)
  b8:= uint8((float32(b)/65535.0)*255.0 + 0.5)

  // Intenta actualitzar en la base de dades
  if err:= self.plats.UpdatePlatform ( self.id, name, r8, g8, b8 ); err != nil {
    return err
  }
  
  // Modifica els atributs "cached"
  self.color= c
  self.name= name
  
  return nil
  
} // end Update
