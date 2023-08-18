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
 *  main.go - Implementa el punt d'entrada a l'aplicació gràfica.
 */

package view

import (
  "fmt"
  "os"
  
  "github.com/diamondburned/gotk4/pkg/gtk/v4"
  "github.com/diamondburned/gotk4/pkg/gio/v2"
)


func startup ( app *gtk.Application ) {
  window:= gtk.NewApplicationWindow ( app )
  window.SetTitle ( "imgteka" )
  window.SetChild ( gtk.NewLabel ( "Prova!" ) )
  window.SetDefaultSize ( 400, 300 )
  window.Show ()
}


func activate ( app *gtk.Application ) {
  app.ActiveWindow ().Present ()
}


func Run() error {
  
  // Crea
  app:= gtk.NewApplication("com.github.adriagipas.imgteka",
    gio.ApplicationFlagsNone)
  app.ConnectStartup ( func() { startup ( app ) } )
  app.ConnectActivate ( func() { activate ( app ) } )
  
  // Executa
  if ecode:= app.Run ( os.Args ); ecode > 0 {
    return fmt.Errorf ( "Unable to run GTK application (Error Code: %d)",
      ecode )
  }

  return nil
  
}
