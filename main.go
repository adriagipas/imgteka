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
 *  main.go - Utilitat per a gestionar imatges de disquets, roms,
 *            cdroms, etc.
 */

package main

import (
  "log"

  "github.com/adriagipas/imgteka/lock"
  "github.com/adriagipas/imgteka/model"
  "github.com/adriagipas/imgteka/view"
)

func main() {

  // Inicialitza log.
  log.SetPrefix ( "[imgteka]" )
  log.SetFlags ( 0 )

  // Executa
  if lock.Init () {
    model,err:= model.New ()
    if err != nil {
      log.Fatal ( err )
    }
    if err:= view.Run ( model ); err != nil {
      log.Fatal ( err )
    }
    model.Close ()
    lock.Close ()
  }
  
}
