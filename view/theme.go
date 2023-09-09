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
 *  theme.go - Tema d'imgteka.
 */

package view

import (
  "image/color"
  
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/theme"
)




type ImgtekaTheme struct {
}


var _ fyne.Theme= (*ImgtekaTheme)(nil)


func (self *ImgtekaTheme) Color(
  
  name    fyne.ThemeColorName,
  variant fyne.ThemeVariant,
  
) color.Color {
  return theme.DefaultTheme ().Color ( name, variant )
} // end Color


func (self *ImgtekaTheme) Icon( name fyne.ThemeIconName ) fyne.Resource {
	return theme.DefaultTheme ().Icon ( name )
} // end Icon


func (self *ImgtekaTheme) Font( style fyne.TextStyle ) fyne.Resource {

  var ret fyne.Resource= theme.DefaultTheme ().Font ( style )
  if !style.Symbol {
    if style.Monospace {
      /*
      if style.Bold {
        if style.Italic {
        } else {
        }
      } else if style.Italic {
      } else {
      }
      */
    } else {
      if style.Bold {
        if style.Italic {
          ret= resourceNotoSansBoldItalicTtf
        } else {
          ret= resourceNotoSansBoldTtf
        }
      } else if style.Italic {
        ret= resourceNotoSansItalicTtf
      } else {
        ret= resourceNotoSansJPRegularTtf
      }
    }
  } else {
    ret= resourceNotoSansJPRegularTtf
  }
  
  return ret
  
} // end Font


func (self *ImgtekaTheme) Size( name fyne.ThemeSizeName ) float32 {
	return theme.DefaultTheme ().Size ( name )
} // end Size
