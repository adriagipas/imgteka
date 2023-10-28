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
 *  utils.go - Funcions d'utilitat.
 */

package file_type

import (
  "errors"
  "fmt"
  "io"
  "os"
)



/************/
/* FUNCIONS */
/************/

// Llig bytes d'un fitxer fent comprovacions
func ReadBytes(

  f        *os.File,
  f_begin  int64,
  f_length int64,
  buf      []byte,
  offset   int64,
  
) error {

  length := int64(len(buf))
  
  end := f_begin + f_length
  if offset < f_begin || offset >= end {
    return fmt.Errorf ( "error while reading bytes: offset (%d) is out"+
      " of bounds (offset:%d, length:%d)",
      offset, f_begin, f_length )
  }
  my_end := offset + length
  if my_end > end {
    return fmt.Errorf ( "error while reading bytes: segment "+
      "(offset:%d, length:%d) is out of bounds (offset:%d, length:%d)",
    offset, length, f_begin, f_length )
  }

  // Llig bytes
  nbytes,err := f.ReadAt ( buf, offset )
  if err != nil { return err }
  if nbytes != len(buf) {
    return errors.New("Unexpected error occurred while reading bytes")
  }
  
  return nil
  
} // ReadBytes



/******************/
/* SUBFILE READER */
/******************/

type SubfileReader struct {

  f           *os.File
  data_offset int64
  data_length int64
  pos         int64 // Posició actual
  
}


func (self *SubfileReader) Read(buf []byte) (int,error) {
  
  // Calcula el que queda
  remain := self.data_length-(self.pos-self.data_offset)
  if remain <= 0 { return 0,io.EOF }

  // Reajusta buffer
  lbuf := int64(len(buf))
  var sbuf []byte
  if lbuf > remain {
    sbuf= buf[:remain]
  } else {
    sbuf= buf
  }
  
  // Llig
  if err := ReadBytes ( self.f, self.data_offset,
    self.data_length, sbuf, self.pos ); err != nil {
    return -1,err
  }
  ret := len(sbuf)
  self.pos+= int64(ret)
  
  return ret,nil
  
} // end Read


func (self *SubfileReader) Size() int64 { return self.data_length }


func NewSubfileReader(
  
  fd          *os.File,
  data_offset int64,
  data_length int64,
  
) (*SubfileReader,error) {

  ret := SubfileReader{
    f: fd,
    data_offset: data_offset,
    data_length: data_length,
    pos: data_offset,
  }
  
  return &ret,nil
  
} // end NewSubfileReader
