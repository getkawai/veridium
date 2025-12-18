package escape

import "github.com/kawai-network/veridium/pkg/htmltomarkdown/marker"

var placeholderRune rune = marker.MarkerEscaping

// IMPORTANT: Only internally we assume it is only byte
var placeholderByte byte = marker.BytesMarkerEscaping[0]
