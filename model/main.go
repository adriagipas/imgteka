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
 *  main.go - Implementa la llògica del sistema.
 */

package model

import (
  "fmt"
  "image"
  "image/color"

  "github.com/adriagipas/imgteka/view"
)




type Model struct {
  dirs  *Dirs
  db    *Database
  plats *Platforms
}


func New() (*Model,error) {

  // Crea objectes
  dirs:= NewDirs ()
  db,err:= NewDatabase ( dirs )
  if err != nil { return nil,err }
  plats:= NewPlatforms ( db )
  
  // Crea model
  ret:= Model{
    dirs : dirs,
    db : db,
    plats : plats,
  }
  
  return &ret,nil
  
} // end New


func (self *Model) Close() {
  self.db.Close ()
} // end Close


func (self *Model) RootEntries() []int64 {
  fmt.Println ( "TODO RootEntries !" )
  return make([]int64,0)
} // end RootEntries


func (self *Model) GetEntry( id int64 ) view.Entry {
  fmt.Println ( "TODO GetEntry !" )
  return &_Entry{}
} // end GetEntry


func (self *Model) GetPlatformIDs() []int {
  return self.plats.GetIDs ()
} // end GetPlatformIDs


func (self *Model) GetPlatform( id int ) view.Platform {
  return self.plats.GetPlatform ( id )
} // end GetPlatform


func (self *Model) GetFile( id int64 ) view.File {
  fmt.Println ( "TODO GetFile !" )
  return &_File{}
} // end GetFile


func (self *Model) GetLabelInfo( id int ) (name string,mcolor color.Color) {
  fmt.Println ( "TODO GetLabelInfo !" )
  return "TODO!",color.Black
} // end GetLabelInfo


func (self *Model) GetStats() view.Stats {
  fmt.Println ( "TODO GetFile !" )
  return &_Stats{}
}


func (self *Model) AddPlatform(
  short_name string,
  name       string,
  c          color.Color,
) error {
  return self.plats.Add ( short_name, name, c )
} // end AddPlatform


func (self *Model) RemovePlatform( id int ) error {
  return self.plats.Remove ( id )
} // end RemovePlatform


/// TODO!!!!!!!!!!!!!!!!!!!!! /////////////////////////////////////////////////
type _Entry struct {}
func (self *_Entry) GetName() string {return "Fake name"}
func (self *_Entry) GetPlatformID() int {return 0}
func (self *_Entry) GetFileIDs() []int64 {return make([]int64,0)}
func (self *_Entry) GetCover() image.Image {return nil}
func (self *_Entry) GetLabelIDs() []int {return make([]int,0)}
type _File struct{}
func (self *_File) GetName() string {return "Fake file"}
func (self *_File) GetType() string {return "Fake type"}
func (self *_File) GetMetadata() []view.StringPair {return make([]view.StringPair,0)}
type _Stats struct{}
func (self *_Stats) GetNumEntries() int64 {return 0}
func (self *_Stats) GetNumFiles() int64 {return 0}
