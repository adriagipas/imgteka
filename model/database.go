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
  "errors"
  "log"
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

const _CREATE_ENTRIES= `
CREATE TABLE IF NOT EXISTS ENTRIES (
       id INTEGER PRIMARY KEY,
       name TEXT NOT NULL,
       platform_id INTEGER NOT NULL,
       cover_id INTEGER DEFAULT -1,
       UNIQUE (platform_id,name),
       FOREIGN KEY (platform_id)
               REFERENCES PLATFORMS (id)
               ON DELETE CASCADE
               ON UPDATE NO ACTION,
       FOREIGN KEY (cover_id)
               REFERENCES FILES (id)
               ON UPDATE NO ACTION
);
`

const _CREATE_LABELS= `
CREATE TABLE IF NOT EXISTS LABELS (
       id INTEGER PRIMARY KEY,
       name TEXT NOT NULL UNIQUE,
       color_r INTEGER NOT NULL,
       color_g INTEGER NOT NULL,
       color_b INTEGER NOT NULL
);
`

const _CREATE_ENTRY_LABEL_PAIRS= `
CREATE TABLE IF NOT EXISTS ENTRY_LABEL_PAIRS (
       entry_id INTEGER NOT NULL,
       label_id INTEGER NOT NULL,
       UNIQUE (entry_id,label_id),
       FOREIGN KEY (entry_id)
               REFERENCES ENTRIES (id)
               ON DELETE CASCADE
               ON UPDATE NO ACTION,
       FOREIGN KEY (label_id)
               REFERENCES LABELS (id)
               ON DELETE CASCADE
               ON UPDATE NO ACTION
);
`


const _CREATE_FILES= `
CREATE TABLE IF NOT EXISTS FILES (
       id INTEGER PRIMARY KEY,
       name TEXT NOT NULL,
       entry_id INTEGER NOT NULL,
       type INTEGER NOT NULL,
       size INTEGER NOT NULL,
       md5 TEXT NOT NULL,
       sha1 TEXT NOT NULL,
       extra_json TEXT NOT NULL,
       last_check INTEGER NOT NULL,
       UNIQUE (entry_id,name),
       UNIQUE (type,name),
       FOREIGN KEY (entry_id)
               REFERENCES ENTRIES (id)
               ON DELETE CASCADE
               ON UPDATE NO ACTION
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
  if _,err:= db.Exec ( _CREATE_ENTRIES ); err != nil {
    return nil,err
  }
  if _,err:= db.Exec ( _CREATE_LABELS ); err != nil {
    return nil,err
  }
  if _,err:= db.Exec ( _CREATE_ENTRY_LABEL_PAIRS ); err != nil {
    return nil,err
  }
  if _,err:= db.Exec ( _CREATE_FILES ); err != nil {
    return nil,err
  }
  
  return db,nil
  
} // end initDatabase




/****************/
/* PART PÚBLICA */
/****************/

type Database struct {
  conn    *sql.DB
  last_tx *sql.Tx
}


func NewDatabase ( dirs *Dirs ) (*Database,error){

  // Crea objectes
  conn,err:= initDatabase ( dirs )
  if err != nil { return nil,err }

  // Crea objecte
  ret:= Database{
    conn    : conn,
    last_tx : nil,
  }

  return &ret,nil
  
} // end NewDatabase


func (self *Database) Close() {
  self.conn.Close ()
} // end Close


func (self *Database) CommitLastTransaction() error {

  if self.last_tx == nil {
    return errors.New ( "No hi ha cap transacció pendent" )
  }
  err:= self.last_tx.Commit ()
  self.last_tx= nil

  return err
  
} // end CommitLastTransaction


// NOTA!!! Aquesta funció caldrà actualitzar-la quan afegim els
// filtres de cerca.
func (self *Database) GetNumEntries() (int64,error) {

    // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT COUNT(*)
FROM ENTRIES;
` )
  if err != nil { return -1,err }
  defer rows.Close ()

  // Recorre consulta
  if !rows.Next () {
    return -1,errors.New ( "Error inesperat en Database.GetNumEntries" )
  }
  var ret int64
  err= rows.Scan ( &ret )
  if err != nil { return -1,err }
  
  return ret,rows.Err ()
  
} // end GetNumEntries


func (self *Database) RollbackLastTransaction() error {

  if self.last_tx == nil {
    return errors.New ( "No hi ha cap transacció pendent" )
  }
  err:= self.last_tx.Rollback ()
  self.last_tx= nil
  
  return err
  
} // end RollbackLastTransaction


// PLATFORMS ///////////////////////////////////////////////////////////////////

func (self *Database) DeletePlatform( id int ) error {

  _,err:= self.conn.Exec ( `
DELETE FROM PLATFORMS WHERE id=?;
`, id )

  return err
  
} // end DeletePlatform


func (self *Database) GetPlatformNumEntries( id int ) (int64,error) {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT COUNT(CASE WHEN platform_id = ? THEN id END)
FROM ENTRIES;
`, id )
  if err != nil { return -1,err }
  defer rows.Close ()

  // Recorre consulta
  if !rows.Next () {
    return -1,errors.New ( "Error inesperat en Database.GetPlatformEntries" )
  }
  var ret int64
  err= rows.Scan ( &ret )
  if err != nil { return -1,err }
  
  return ret,rows.Err ()
  
} // end GetPlatformNumEntries


func (self *Database) LoadPlatforms( plats *Platforms ) error {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT id,short_name,name,color_r,color_g,color_b
FROM PLATFORMS
ORDER BY short_name ASC;
` )
  if err != nil { return err }
  defer rows.Close ()

  // Recorre consulta
  for rows.Next () {
    var id int
    var short_name,name string
    var r,g,b int
    err= rows.Scan ( &id, &short_name, &name, &r, &g, &b )
    if err != nil { return err }
    plats.add ( id, short_name, name, uint8(r), uint8(g), uint8(b) )
  }
  
  return rows.Err ()
  
} // end LoadPlatforms


func (self *Database) RegisterPlatform(

  short_name string,
  name       string,
  r,g,b      uint8,
  
) error {

  _,err:= self.conn.Exec ( `
   INSERT INTO PLATFORMS(short_name, name, color_r, color_g, color_b)
          VALUES(?,?,?,?,?);
`, short_name, name, r, g, b )

  return err
  
} // end RegisterPlatform


func (self *Database) UpdatePlatform(
  id    int,
  name  string,
  r,g,b uint8,
) error {

  _,err:= self.conn.Exec ( `
UPDATE PLATFORMS SET name = ?, color_r = ?, color_g = ?, color_b = ?
       WHERE id = ?;
`, name, int(r), int(g), int(b), id )
  
  return err
  
} // end UpdatePlatform


// LABELS //////////////////////////////////////////////////////////////////////

func (self *Database) DeleteLabel( id int ) error {

  _,err:= self.conn.Exec ( `
DELETE FROM LABELS WHERE id=?;
`, id )
  
  return err
  
} // end DeleteLabel


func (self *Database) LoadLabels( labels *Labels ) error {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT id,name,color_r,color_g,color_b
FROM LABELS
ORDER BY name ASC;
` )
  if err != nil { return err }
  defer rows.Close ()

  // Recorre consulta
  for rows.Next () {
    var id int
    var name string
    var r,g,b int
    err= rows.Scan ( &id, &name, &r, &g, &b )
    if err != nil { return err }
    labels.add ( id, name, uint8(r), uint8(g), uint8(b) )
  }
  
  return rows.Err ()
  
} // end LoadLabels


func (self *Database) RegisterLabel(

  name  string,
  r,g,b uint8,
  
) error {

  _,err:= self.conn.Exec ( `
   INSERT INTO LABELS(name, color_r, color_g, color_b)
          VALUES(?,?,?,?);
`, name, r, g, b )

  return err
  
} // end RegisterLabel


func (self *Database) UpdateLabel(
  
  id    int,
  name  string,
  r,g,b uint8,
  
) error {

  _,err:= self.conn.Exec ( `
UPDATE LABELS SET name = ?, color_r = ?, color_g = ?, color_b = ?
       WHERE id = ?;
`, name, int(r), int(g), int(b), id )
  
  return err
  
} // end UpdateLabel


// ENTRIES /////////////////////////////////////////////////////////////////////

func (self *Database) DeleteEntryWithoutCommit( id int64 ) error {

  // Prepara
  tx,err:= self.conn.Begin ()
  if err != nil { log.Fatal ( err ) }
  stmt,err:= tx.Prepare ( `
DELETE FROM ENTRIES WHERE id=?;
` )
  if err != nil { log.Fatal ( err ) }
  defer stmt.Close ()
  
  // Elimina
  _,err= stmt.Exec ( id )
  if err != nil { tx.Rollback (); return err }
  
  // Registra transacció
  self.last_tx= tx
  
  return nil
  
} // end DeleteEntryWithoutCommit


// NOTA!!! En algun moment caldrà ficar la query.
// NOTA!!! Sols es carreguen les dades bàsiques.
func (self *Database) LoadEntries( entries *Entries ) error {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT id,name,platform_id,cover_id
FROM ENTRIES
ORDER BY name ASC;
` )
  if err != nil { return err }
  defer rows.Close ()

  // Recorre consulta
  for rows.Next () {
    var id,cover_id int64
    var name string
    var platform_id int
    err= rows.Scan ( &id, &name, &platform_id, &cover_id )
    if err != nil { return err }
    entries.add ( id, name, platform_id, cover_id )
  }
  
  return rows.Err ()
  
} // end LoadEntries


func (self *Database) RegisterEntryWithoutCommit(

  name        string,
  platform_id int,
  
) error {

  // Prepara
  tx,err:= self.conn.Begin ()
  if err != nil { log.Fatal ( err ) }
  stmt,err:= tx.Prepare ( `
   INSERT INTO ENTRIES(name, platform_id)
          VALUES(?,?);
` )
  if err != nil { log.Fatal ( err ) }
  defer stmt.Close ()

  // Inserta
  _,err= stmt.Exec ( name, platform_id )
  if err != nil { tx.Rollback (); return err }

  // Registra transacció
  self.last_tx= tx
  
  return nil
  
} // end RegisterEntryWithoutCommit


func (self *Database) UpdateEntryNameWithoutCommit(

  id          int64,
  name        string,
  
) error {

  // Prepara
  tx,err:= self.conn.Begin ()
  if err != nil { log.Fatal ( err ) }
  stmt,err:= tx.Prepare ( `
UPDATE ENTRIES SET name = ?
       WHERE id = ?;
` )
  if err != nil { log.Fatal ( err ) }
  defer stmt.Close ()
  
  // Inserta
  _,err= stmt.Exec ( name, id )
  if err != nil { tx.Rollback (); return err }
  
  // Registra transacció
  self.last_tx= tx
  
  return nil
  
} // end UpdateEntryNameWithoutCommit


// ENTRY_LABEL_PAIRS ///////////////////////////////////////////////////////////

func (self *Database) DeleteEntryLabelPair( id int64, label_id int ) error {

  _,err:= self.conn.Exec ( `
DELETE FROM ENTRY_LABEL_PAIRS
       WHERE entry_id = ? AND label_id = ?;
`, id, label_id )
  
  return err
  
} // end DeleteEntryLabelPair


func (self *Database) GetLabelNumEntries( id int ) (int64,error) {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT COUNT(CASE WHEN label_id = ? THEN rowid END)
FROM ENTRY_LABEL_PAIRS;
`, id )
  if err != nil { return -1,err }
  defer rows.Close ()

  // Recorre consulta
  if !rows.Next () {
    return -1,errors.New ( "Error inesperat en Database.GetLabelNumEntries" )
  }
  var ret int64
  err= rows.Scan ( &ret )
  if err != nil { return -1,err }
  
  return ret,rows.Err ()
  
} // end GetLabelNumEntries


func (self *Database) LoadLabelsEntry( id int64, e *Entry ) error {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT label_id
FROM ENTRY_LABEL_PAIRS
WHERE entry_id = ?;
`, id )
  if err != nil { return err }
  defer rows.Close ()

  // Recorre consulta
  for rows.Next () {
    var label_id int
    err= rows.Scan ( &label_id )
    if err != nil { return err }
    e.addLabelID ( label_id )
  }
  
  return rows.Err ()
  
} // end LoadLabelsEntry


func (self *Database) RegisterEntryLabelPair( id int64, label_id int ) error {
  
  _,err:= self.conn.Exec ( `
   INSERT INTO ENTRY_LABEL_PAIRS(entry_id, label_id)
          VALUES(?,?);
`, id, label_id )
  
  return err
  
} // end RegisterEntryLabelPair


// FILES ///////////////////////////////////////////////////////////////////////

// S'enten que ja està, no deuria fallar
func (self *Database) GetFile( id int64 ) (
  name       string,
  entry_id   int64,
  file_type  int,
  size       int64,
  md5        string,
  sha1       string,
  json       string,
  last_check int64,
) {
  
  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT name,entry_id,type,size,md5,sha1,extra_json,last_check
FROM LABELS
WHERE id = ?;
`, id )
  if err != nil { log.Fatal ( err ) }
  defer rows.Close ()

  // Recorre consulta
  if !rows.Next () {
    log.Fatal ( "Error inesperat en Database.GetFile" )
  }
  err= rows.Scan ( &name, &entry_id, &file_type, &size, &md5,
    &sha1, &json, &last_check )
  if err != nil { log.Fatal ( err ) }
  
  return
  
} // end GetFile


func (self *Database) LoadFilesEntry( id int64, e *Entry ) error {

  // Consulta base de dades
  rows,err:= self.conn.Query ( `
SELECT id
FROM FILES
WHERE entry_id = ?;
`, id )
  if err != nil { return err }
  defer rows.Close ()

  // Recorre consulta
  for rows.Next () {
    var file_id int64
    err= rows.Scan ( &file_id )
    if err != nil { return err }
    e.addFileID ( file_id )
  }
  
  return rows.Err ()
  
} // end LoadFilesEntry
