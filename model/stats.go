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
 *  stats.go - Classe amb les estadístiques
 */

package model

import (
  "fmt"
  "log"
)


type Stats struct {
  db *Database
}


func NewStats( db *Database ) *Stats {

  ret:= Stats{
    db : db,
  }

  return &ret
  
} // end NewStats


func (self *Stats) GetNumEntries() int64 {

  ret,err:= self.db.GetNumEntries ()
  if err != nil { log.Fatal ( err ) }

  return ret
  
} // end GetNumEntries


func (self *Stats) GetNumFiles() int64 {
  fmt.Println ( "TODO Stats.GetNumFiles !" )
  return 0
} // end GetNumFiles
