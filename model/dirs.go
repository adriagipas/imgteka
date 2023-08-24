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
 *  dirs.go - Gestiona el directoris on es desen els fitxers.
 */

package model

import (
  "path"
  
  "github.com/adrg/xdg"
)




/****************/
/* PART PRIVADA */
/****************/

const _ROOT_NAME= "imgteka"




/****************/
/* PART PÚBLICA */
/****************/

type Dirs struct {
  db_name *string // Nom base de dades
}


func NewDirs() *Dirs {
  
  ret:= Dirs{
    db_name : nil,
  }

  return &ret
  
} // end NewDirs


func (self *Dirs) GetDatabaseName() (string,error) {

  var ret string
  var err error
  
  if self.db_name == nil {
    path:= path.Join ( _ROOT_NAME, "database.db" )
    ret,err= xdg.DataFile ( path )
    if err == nil {
      self.db_name= &ret
    }
  } else {
    ret,err= *self.db_name,nil
  }
  
  return ret,err
  
} // end GetDatabaseName
