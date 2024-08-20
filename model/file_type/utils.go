/*
 * Copyright 2023-2024 Adrià Giménez Pastor.
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
  "bytes"
  "errors"
  "fmt"
  "io"
  "os"
  "strconv"
  "strings"
  
  "github.com/adriagipas/imgcp/cdread"
  "github.com/adriagipas/imgteka/view"
)



/************/
/* FUNCIONS */
/************/

func BytesToStr_trim_0s( data []byte ) string {
  
  data= bytes.TrimRight ( data, "\000" )
  ret:= string(data)
  
  return strings.TrimSpace ( ret )
  
} // end BytesToStr_trim_0s


func NumBytesToStr(num_bytes uint64) string {
  if num_bytes > 1024*1024*1024 { // G
    val := float64(num_bytes)/(1024*1024*1024)
    return strconv.FormatFloat ( val, 'f', 1, 32 ) + "G"
  } else if num_bytes > 1024*1024 { // M
    val := float64(num_bytes)/(1024*1024)
    return strconv.FormatFloat ( val, 'f', 1, 32 ) + "M"
  } else if num_bytes > 1024 { // K
    val := float64(num_bytes)/1024
    return strconv.FormatFloat ( val, 'f', 1, 32 ) + "K"
  } else {
    return strconv.FormatUint ( num_bytes, 10 )
  }
} // end NumBytesToStr


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




/***************/
/* METADATA CD */
/***************/

type _CD_Metadata struct {

  Format      string
  NumSessions int
  NumTracks   int
  
}

func (self *_CD_Metadata) Init( cd cdread.CD ) {

  self.Format= cd.Format ()
  info:= cd.Info ()
  self.NumSessions= len(info.Sessions)
  self.NumTracks= len(info.Tracks)
  
} // end Init


func (self *_CD_Metadata) Parse( v []view.StringPair ) []view.StringPair {

  var kv *KeyValue

  // Format
  kv= &KeyValue{"Format",self.Format}
  v= append(v,kv)

  // Nº Sessions
  if self.NumSessions>1 {
    kv= &KeyValue{"Nº Sessions",fmt.Sprintf("%d",self.NumSessions)}
    v= append(v,kv)
  }

  // Nº Tracks
  if self.NumTracks>1 {
    kv= &KeyValue{"Nº Tracks",fmt.Sprintf("%d",self.NumTracks)}
    v= append(v,kv)
  }
  
  return v
  
} // end Parse




/****************/
/* METADATA ISO */
/****************/

type _ISO_Metadata struct {

  SystemIdentifier        string
  VolumeIdentifier        string
  VolumeSpaceSize         uint32 // Logical blocks
  VolumeSetSize           uint16 // number of disks
  VolumeSequenceNumber    uint16
  LogicalBlockSize        uint16
  VolumeSetIdentifier     string
  PublisherIdentifier     string
  DataPreparerIdentifier  string
  ApplicationIdentifier   string
  CopyrightFileIdentifier string
  AbstractFileIdentifier  string
  BiblioFileIdentifier    string
  VolumeCreation          cdread.ISO_DateTime
  VolumeModification      cdread.ISO_DateTime
  VolumeExpiration        cdread.ISO_DateTime
  VolumeEffective         cdread.ISO_DateTime
  
}


func (self *_ISO_Metadata) Init( iso *cdread.ISO ) {
  
  self.SystemIdentifier= iso.PrimaryVolume.SystemIdentifier
  self.VolumeIdentifier= iso.PrimaryVolume.VolumeIdentifier
  self.VolumeSpaceSize= iso.PrimaryVolume.VolumeSpaceSize
  self.VolumeSetSize= iso.PrimaryVolume.VolumeSetSize
  self.VolumeSequenceNumber= iso.PrimaryVolume.VolumeSequenceNumber
  self.LogicalBlockSize= iso.PrimaryVolume.LogicalBlockSize
  self.VolumeSetIdentifier= iso.PrimaryVolume.VolumeSetIdentifier
  self.PublisherIdentifier= iso.PrimaryVolume.PublisherIdentifier
  self.DataPreparerIdentifier= iso.PrimaryVolume.DataPreparerIdentifier
  self.ApplicationIdentifier= iso.PrimaryVolume.ApplicationIdentifier
  self.CopyrightFileIdentifier= iso.PrimaryVolume.CopyrightFileIdentifier
  self.AbstractFileIdentifier= iso.PrimaryVolume.AbstractFileIdentifier
  self.BiblioFileIdentifier= iso.PrimaryVolume.BiblioFileIdentifier
  self.VolumeCreation= iso.PrimaryVolume.VolumeCreation
  self.VolumeModification= iso.PrimaryVolume.VolumeModification
  self.VolumeExpiration= iso.PrimaryVolume.VolumeExpiration
  self.VolumeEffective= iso.PrimaryVolume.VolumeEffective
  
} // end Init


func (self *_ISO_Metadata) Parse( v []view.StringPair ) []view.StringPair {

  var kv *KeyValue
  
  // Identificadors.
  if len(self.SystemIdentifier)>0 {
    kv= &KeyValue{"Sistema",self.SystemIdentifier}
    v= append(v,kv)
  }
  if len(self.VolumeIdentifier)>0 {
    kv= &KeyValue{"Volum",self.VolumeIdentifier}
    v= append(v,kv)
  }

  // Grandària
  size:= uint64(self.VolumeSpaceSize)*uint64(self.LogicalBlockSize)
  kv= &KeyValue{"Grandària volum",NumBytesToStr(size)}
  v= append(v,kv)
  
  // Set
  if self.VolumeSetSize > 1 {
    kv= &KeyValue{"Nº Discs",fmt.Sprintf("%d",self.VolumeSetSize)}
    v= append(v,kv)
    kv= &KeyValue{"Num. Disc",fmt.Sprintf("%d",self.VolumeSequenceNumber)}
    v= append(v,kv)
    if len(self.VolumeSetIdentifier)>0 {
      kv= &KeyValue{"Conjunt",self.VolumeSetIdentifier}
      v= append(v,kv)
    }
  }

  // Altres aspectes cadenes
  if len(self.PublisherIdentifier)>0 {
    kv= &KeyValue{"Editora",self.PublisherIdentifier}
    v= append(v,kv)
  }
  if len(self.DataPreparerIdentifier)>0 {
    kv= &KeyValue{"Fabricant CD",self.DataPreparerIdentifier}
    v= append(v,kv)
  }
  if len(self.ApplicationIdentifier)>0 {
    kv= &KeyValue{"Aplicació",self.ApplicationIdentifier}
    v= append(v,kv)
  }
  if len(self.CopyrightFileIdentifier)>0 {
    kv= &KeyValue{"Copyright",self.CopyrightFileIdentifier}
    v= append(v,kv)
  }
  if len(self.AbstractFileIdentifier)>0 {
    kv= &KeyValue{"Resum",self.AbstractFileIdentifier}
    v= append(v,kv)
  }
  if len(self.BiblioFileIdentifier)>0 {
    kv= &KeyValue{"Bibliografia",self.BiblioFileIdentifier}
    v= append(v,kv)
  }

  // Dates
  // --> Assumint que no està buit
  ParseDate:= func(dt *cdread.ISO_DateTime) string {
    return fmt.Sprintf ( "%s/%s/%s (%s:%s:%s.%s GMT %d)",
      dt.Day, dt.Month, dt.Year,
      dt.Hour, dt.Minute, dt.Second, dt.HSecond,
      dt.GMT )
  }
  if !self.VolumeCreation.Empty {
    kv= &KeyValue{"Data creació",ParseDate(&self.VolumeCreation)}
    v= append(v,kv)
  }
  if !self.VolumeModification.Empty {
    kv= &KeyValue{"Data última modificació",ParseDate(&self.VolumeModification)}
    v= append(v,kv)
  }
  if !self.VolumeExpiration.Empty {
    kv= &KeyValue{"Data caducitat",ParseDate(&self.VolumeExpiration)}
    v= append(v,kv)
  }
  if !self.VolumeEffective.Empty {
    kv= &KeyValue{"Data disponibilitat",ParseDate(&self.VolumeEffective)}
    v= append(v,kv)
  }
  
  return v
  
} // end Parse
