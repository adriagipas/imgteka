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
 *  commands.go - Estructura que gestiona els comandaments per a obrir
 *                els diferents tipus de fitxers.
 */

package model

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"
  "sync"

  "github.com/adriagipas/imgteka/model/file_type"
)




/****************/
/* PART PÚBLICA */
/****************/

type Commands struct {
  
  dirs    *Dirs
  v       map[int]string // Mapeja identificador tipus a commandament.
  running map[string]bool // Controla fitxers en execució
  mu      sync.Mutex
  
}


func NewCommands( dirs *Dirs ) (*Commands,error) {

  ret:= Commands{
    dirs : dirs,
  }
  
  // Intenta deserialitzar
  fn,err:= dirs.GetCommandsConfName()
  if err != nil { return nil,err }
  f,err:= os.Open ( fn )
  ret.v= make(map[int]string)
  if err == nil {
    defer f.Close ()
    json_dec:= json.NewDecoder ( f )
    if err:= json_dec.Decode ( &ret.v ); err != nil {
      log.Printf ( "S'ha produit un error al decodificar '%s': %s\n",
        fn, err )
      ret.v= make(map[int]string) // Reset
    }
  }

  // Altres
  ret.running= make(map[string]bool)
  
  return &ret,nil
  
} // end NewCommands


func (self *Commands) Close() {

  // Obri fitxer
  fn,err:= self.dirs.GetCommandsConfName()
  if err != nil {
    log.Printf ( "Error inesperat en 'Commands': %s", err )
    return
  }
  f,err:= os.Create ( fn )
  if err != nil {
    log.Printf ( "No s'ha pogut crear '%s': %s", fn, err )
    return
  }
  defer f.Close ()
  
  // Serialitza
  json_enc:= json.NewEncoder ( f )
  if err:= json_enc.Encode ( self.v ); err != nil {
    log.Printf ( "No s'ha pogut desar el contingut en '%s': %s", fn, err )
    return
  }
  
} // end Close


// Cadena buida indica que no hi ha
func (self *Commands) GetCommand( type_id int ) string {
  return self.v[type_id]
} // end GetCommand


func (self *Commands) Run( type_id int, file_path string ) error {

  // Selecciona tipus
  ft,err:= file_type.Get ( type_id )
  if err != nil { return err }
  
  // Obté commandament
  cmd_name,ok:= self.v[type_id]
  if !ok {
    return fmt.Errorf ( "No s'ha especificat ningun comandament" +
      " per al tipus '%s'", ft.GetName () )
  }

  // Comprova que no estiga ja en execució. Si ho està torna sense
  // error.
  if _,ok:= self.running[file_path]; ok {
    return nil
  }
  
  // Crea commandament
  cmd:= exec.Command ( cmd_name, file_path )
  err= cmd.Start ()
  if err == nil {

    // Marca com en execució
    self.mu.Lock ()
    self.running[file_path]= true
    self.mu.Unlock ()

    // Llança fil que s'espera que acabe
    go func() {
      cmd.Wait ()
      self.mu.Lock ()
      delete(self.running,file_path)
      self.mu.Unlock ()
    }()
    
  }

  return err
  
} // end Run


// Una cadena buida indica que no hi ha commandament
func (self *Commands) SetCommand( type_id int, command string ) {

  command= strings.TrimSpace ( command )
  if command == "" {
    delete(self.v,type_id)
  } else {
    self.v[type_id]= command
  }
  
} // end SetCommand
