package types

// GraphicInfo structure for each graphic info in `GraphicInfo*.bin`, 40 bytes.
type GraphicInfo struct {
	ID     int32 // ID, it's the index of the graphic in the file, but it could be duplicated.
	Addr   int32 // File offset of the graphic data in `Graphic*.bin`.
	Len    int32 // Length of the graphic data in `Graphic*.bin`.
	OffX   int32
	OffY   int32
	Width  int32
	Height int32
	GridW  byte
	GridH  byte
	Access byte
	_      [5]byte
	MapID  int32 // Map ID, it's the index of the map in `map/*.dat`, 0 for not a map graphic.
}

// GraphicDataHeader structure for each graphic header in `Graphic*.bin`, 16 bytes.
type GraphicDataHeader struct {
	Magic   [2]byte // "RD" for a valid graphic
	Version byte    // 0 for raw data, 1 for encoded data, 2 for raw data with palette, 3 for encoded data with palette
	_       byte    //
	Width   int32   // Width of graphic, it shouldn't be trusted, use GraphicInfo.Width instead.
	Height  int32   // Height of graphic, it shouldn't be trusted, use GraphicInfo.Height instead.
	Len     int32   // Length of graphic data, it shouldn't be trusted, use GraphicInfo.Len instead.
}

// GraphicData structure for each graphic data in `Graphic*.bin`, 16(header) + x(data) bytes.
type GraphicData struct {
	GraphicDataHeader
	RawData []byte
}