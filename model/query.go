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
 *  query.go - Estructura que representa una consulta en la base de
 *             dades. Es crea a partir d'un text.
 */

package model

import (
  "strings"
  "text/scanner"
)



/****************/
/* PART PRIVADA */
/****************/

const (
  _WAIT_TOKEN     = 0
  _LABEL_WAIT_SEP = 1
  _LABEL_WAIT_VAL = 2
  _PLATF_WAIT_SEP = 3
  _PLATF_WAIT_VAL = 4
)


func addEntry( q *QueryOr, token string, typ int ) {

  // Elimina cometes
  if len(token)>2 {
    if token[0] == '"' {
      token= token[1:]
    }
    if token[len(token)-1] == '"' {
      token= token[:len(token)-1]
    }
  }
  token= strings.TrimSpace(token)

  // Inserta
  if len(token)>0 {
    q.Queries= append(q.Queries,QueryEntry{
      value : token,
      typ   : typ,
    })
  }
  
} // end NewQueryEntry



/****************/
/* PART PÚBLICA */
/****************/

const (
  QUERY_TYPE_NAME_ENTRY = 0
  QUERY_TYPE_LABEL      = 1
  QUERY_TYPE_PLATFORM   = 2
)


type QueryEntry struct {
  value string
  typ   int
}


type QueryOr struct {
  Queries []QueryEntry // S'ha de cumplir alguna
}


type Query struct {
  OrQueries []QueryOr // S'han de cumplir totes
}


func NewQuery( query_text string ) *Query {

  // Crea objecte
  ret:= Query{}
  ret.OrQueries= []QueryOr{QueryOr{
    Queries : nil,
  }}

  // Parseja text
  var s scanner.Scanner
  s.Init ( strings.NewReader ( query_text ) )
  state:= _WAIT_TOKEN
  current_oq:= &ret.OrQueries[0]
  for tok:= s.Scan (); tok != scanner.EOF; tok= s.Scan () {
    
    val:= s.TokenText ()

    // El Símbol '+' (AND) té prioritat absoluta
    if val == "+" {
      state= _WAIT_TOKEN
      if len(current_oq.Queries)>0 {
        ret.OrQueries= append(ret.OrQueries,QueryOr{Queries:nil})
        current_oq= &ret.OrQueries[len(ret.OrQueries)-1]
      }
      
    } else {
      switch state {
      case _WAIT_TOKEN: // WAIT TOKEN
        if val == "l" {
          state= _LABEL_WAIT_SEP
        } else if val == "p" {
          state= _PLATF_WAIT_SEP
        } else {
          addEntry( current_oq, val, QUERY_TYPE_NAME_ENTRY )
        }

      case _LABEL_WAIT_SEP: // LABEL WAIT SEP
        if val == ":" {
          state= _LABEL_WAIT_VAL
        } else {
          addEntry( current_oq, "l", QUERY_TYPE_NAME_ENTRY )
          state= _WAIT_TOKEN
        }

      case _LABEL_WAIT_VAL: // LABEL WAIT VAL
        addEntry( current_oq, val, QUERY_TYPE_LABEL )
        state= _WAIT_TOKEN

      case _PLATF_WAIT_SEP: // PLATF WAIT SEP
        if val == ":" {
          state= _PLATF_WAIT_VAL
        } else {
          addEntry( current_oq, "p", QUERY_TYPE_NAME_ENTRY )
          state= _WAIT_TOKEN
        }

      case _PLATF_WAIT_VAL: // PLATF WAIT VAL
        addEntry( current_oq, val, QUERY_TYPE_PLATFORM )
        state= _WAIT_TOKEN
        
      }
      
    }
    
  }
  
  // Ignora últim and si està buit
  if len(ret.OrQueries)>0 && ret.OrQueries[len(ret.OrQueries)-1].Queries==nil {
    ret.OrQueries= ret.OrQueries[:len(ret.OrQueries)-1]
  }
  
  return &ret
  
} // end NewQuery
