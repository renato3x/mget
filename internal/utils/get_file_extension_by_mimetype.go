package utils

import (
  "mime"
  "strings"
)

func GetFileExtensionByMimetype(mimetype string) string {
  pureMime := strings.Split(mimetype, ";")[0]
  
  exts, err := mime.ExtensionsByType(pureMime)
  if err != nil || len(exts) == 0 {
    return ".bin"
  }
  return exts[0]
}
