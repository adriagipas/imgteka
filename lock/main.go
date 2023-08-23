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
 *  main.go - Serveix per impedir que hi hasquen múltiples processos
 *            executant-se per a la mateixa sessió.
 */

package lock

import (
  "log"
  
  "github.com/godbus/dbus/v5"
)




/****************/
/* PART PRIVADA */
/****************/

var _con *dbus.Conn

const _LOCK_NAME= "org.github.adriagipas.imgteka.lock"
const _LOCK_INT= _LOCK_NAME+".interface"
const _LOCK_OBJECT= "/org/github/adriagipas/imgteka/lock"




/****************/
/* PART PÚBLICA */
/****************/

func Init () bool {

  var err error
  // Connecta
  _con,err= dbus.ConnectSessionBus ()
  if err != nil {
    log.Printf ( "no s'ha pogut establir connexió amb D-BUS: %s", err )
    return true // En aquest cas permet que continue
  }

  // Fica nom
  res,err:= _con.RequestName ( _LOCK_NAME, dbus.NameFlagDoNotQueue )
  if err != nil {
    log.Printf ( "error a l'obtindre el nom de D-BUS: %s", err )
    return true // En aquest cas permet que continue
  }

  ret:= res==dbus.RequestNameReplyPrimaryOwner

  // Senyals
  if ret {
    if err:= _con.AddMatchSignal (
      dbus.WithMatchInterface ( _LOCK_INT ),
    ); err != nil {
      log.Printf ( "no s'ha pogut registrar la regla de la interfície"+
        " en D-BUS: %s", err )
    }
  } else {
    if err:= _con.Emit ( _LOCK_OBJECT, _LOCK_INT+".ShowWin" ); err != nil {
      log.Printf ( "no s'ha pogut emetre la senyal: %s", err )
    }
    
  }
  
  return ret
  
} // end Init


func Close () {
  _con.Close ()
} // end Close


// CheckSignals es bloqueja
func CheckSignals () bool {

  c:= make ( chan *dbus.Signal, 100 )
	_con.Signal ( c )
  for v:= range c {
    if v.Name == _LOCK_INT+".ShowWin" {
      return true
    }
  }
  
  return false
  
} // end CheckSignals
