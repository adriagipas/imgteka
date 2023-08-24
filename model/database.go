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
 *  database.go - Base de dades.
 */

package model

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)




/****************/
/* PART PRIVADA */
/****************/

const _CREATE_PLATFORMS= `
CREATE TABLE IF NOT EXISTS PLATFORMS (
       id INTEGER PRIMARY KEY,
       short_name TEXT NOT NULL UNIQUE,
       name TEXT NOT NULL,
       color_r INTEGER NOT NULL,
       color_g INTEGER NOT NULL,
       color_b INTEGER NOT NULL
);
`


func initDatabase ( dirs *Dirs ) (*sql.DB,error) {

  // Nom
  db_fn,err:= dirs.GetDatabaseName ()
  if err != nil { return nil,err }

  // Connecta
  db,err:= sql.Open ( "sqlite3", db_fn )
  if err != nil { return nil,err }

  // Crea taules si cal
  if _,err:= db.Exec ( _CREATE_PLATFORMS ); err != nil {
    return nil,err
  }
  
  return db,nil
  
} // end initDatabase




/****************/
/* PART PÚBLICA */
/****************/

type Database struct {
  conn *sql.DB
}


func NewDatabase ( dirs *Dirs ) (*Database,error){

  // Crea objectes
  conn,err:= initDatabase ( dirs )
  if err != nil { return nil,err }

  // Crea objecte
  ret:= Database{
    conn : conn,
  }

  return &ret,nil
  
} // end NewDatabase


func (self *Database) Close () {
  self.conn.Close ()
} // end Close
