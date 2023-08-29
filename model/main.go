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
  "image/color"
  "log"

  "github.com/adriagipas/imgteka/model/file_type"
  "github.com/adriagipas/imgteka/view"
)




type Model struct {
  dirs    *Dirs
  db      *Database
  plats   *Platforms
  labels  *Labels
  files   *Files
  entries *Entries
  stats   *Stats
}


func New() (*Model,error) {

  // Crea objectes
  dirs:= NewDirs ()
  db,err:= NewDatabase ( dirs )
  if err != nil { return nil,err }
  plats:= NewPlatforms ( db )
  labels:= NewLabels ( db )
  files:= NewFiles ( db, plats, dirs )
  entries:= NewEntries ( db, plats, labels, files, dirs )
  stats:= NewStats ( db )
  
  // Crea model
  ret:= Model{
    dirs    : dirs,
    db      : db,
    plats   : plats,
    labels  : labels,
    files   : files,
    entries : entries,
    stats   : stats,
  }
  
  return &ret,nil
  
} // end New


func (self *Model) Close() {
  self.db.Close ()
} // end Close


func (self *Model) RootEntries() []int64 {
  return self.entries.GetIDs ()
} // end RootEntries


func (self *Model) GetEntry( id int64 ) view.Entry {
  return self.entries.Get ( id )
} // end GetEntry


func (self *Model) GetPlatformIDs() []int {
  return self.plats.GetIDs ()
} // end GetPlatformIDs


func (self *Model) GetPlatform( id int ) view.Platform {
  return self.plats.GetPlatform ( id )
} // end GetPlatform


func (self *Model) GetFile( id int64 ) view.File {
  return self.files.Get ( id )
} // end GetFile


func (self *Model) GetLabelIDs() []int {
  return self.labels.GetIDs ()
} // end GetLabelIDs


func (self *Model) GetLabel( id int ) view.Label {
  return self.labels.Get ( id )
} // end GetLabelInfo


func (self *Model) GetStats() view.Stats {
  return self.stats
} // end GetStates


func (self *Model) GetFileTypeIDs() []int {
  return file_type.GetIDs ()
} // end GetFileTypeIDs


func (self *Model) GetFileTypeName(id int) string {
  
  ft,err:= file_type.Get ( id )
  if err != nil {
    log.Fatal ( err )
  }

  return ft.GetName ()
  
} // end GetFileTypeName


func (self *Model) AddPlatform(
  short_name string,
  name       string,
  c          color.Color,
) error {
  return self.plats.Add ( short_name, name, c )
} // end AddPlatform


func (self *Model) AddLabel( name string, c color.Color ) error {
  return self.labels.Add ( name, c )
} // end AddLabel


func (self *Model) AddEntry( name string, platform_id int ) error {
  return self.entries.Add ( name, platform_id )
} // end AddEntry


func (self *Model) RemovePlatform( id int ) error {
  return self.plats.Remove ( id )
} // end RemovePlatform


func (self *Model) RemoveLabel( id int ) error {
  return self.labels.Remove ( id )
} // end RemoveLabel


func (self *Model) RemoveEntry( id int64 ) error {
  return self.entries.Remove ( id )
} // end RemoveEntry
