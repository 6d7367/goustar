package goustar

import "strconv"
import "math"

// http://www.mkssoftware.com/docs/man4/tar.4.asp
const (
	// name of file
	USTarHeaderNameLen = 100
	// file mode
	USTarHeaderModeLen = 8
	// owner user ID
	USTarHeaderUidLen = 8
	// owner group ID
	USTarHeaderGidLen = 8
	// length of file in bytes
	USTarHeaderSizeLen = 12
	// modify time of file
	USTarHeaderMtimeLen = 12
	// checksum for header
	USTarHeaderChksumLen = 8
	// type of file
	USTarHeaderTypeflagLen = 1
	// name of linked file
	USTarHeaderLinknameLen = 100
	// USTAR indicator
	USTarHeaderMaginLen = 6
	// USTAR version
	USTarHeaderVersionLen = 2
	// owner user name
	USTarHeaderUnameLen = 32
	// owner group name
	USTarHeaderGnameLen = 32
	// device major number
	USTarHeaderDevmajorLen = 8
	// device major number
	USTarHeaderDevminorLen = 8
	// prefix for file name
	USTarHeaderPrefixLen = 155
	// data blocksize
	TarBlocksize = 512
)

type USTar struct {
	Name      []byte
	Mode      []byte
	Uid       []byte
	Gid       []byte
	Size      []byte
	Mtime     []byte
	Chksum    []byte
	Typeflag  []byte
	Linkname  []byte
	Magic     []byte
	Version   []byte
	Uname     []byte
	Gname     []byte
	Devmajor  []byte
	Devminor  []byte
	Prefix    []byte
	RawHeader []byte
	Content   []byte
}

func USTarFromRaw(raw []byte) []USTar {
	items := make([]USTar, 0)
	USTarMagic := []byte{117, 115, 116, 97, 114, 0}
	position := 0
	for {
		ustar := USTar{}
		// position = 0
		ustar.Name = raw[position:position + USTarHeaderNameLen]
		position += USTarHeaderNameLen
		// position = 100
		ustar.Mode = raw[position:position + USTarHeaderModeLen]
		position += USTarHeaderModeLen
		// position = 108
		ustar.Uid = raw[position: position + USTarHeaderUidLen]
		position += USTarHeaderUidLen
		// position = 116
		ustar.Gid = raw[position: position + USTarHeaderGidLen]
		position += USTarHeaderGidLen
		// position = 124
		ustar.Size = raw[position: position + USTarHeaderSizeLen]
		position += USTarHeaderSizeLen
		// position = 136
		ustar.Mtime = raw[position : position + USTarHeaderMtimeLen]
		position += USTarHeaderMtimeLen
		// position = 148
		ustar.Chksum = raw[position : position + USTarHeaderChksumLen]
		position += USTarHeaderChksumLen
		// position = 156
		ustar.Typeflag = raw[position : position + USTarHeaderTypeflagLen]
		position += USTarHeaderTypeflagLen
		// position = 157
		ustar.Linkname = raw[position : position + USTarHeaderLinknameLen]
		position += USTarHeaderLinknameLen
		// position = 257
		ustar.Magic = raw[position : position + USTarHeaderMaginLen]
		position += USTarHeaderMaginLen
		// position = 263
		ustar.Version = raw[position : position + USTarHeaderVersionLen]
		position += USTarHeaderVersionLen
		// position = 265
		ustar.Uname = raw[position : position + USTarHeaderUnameLen]
		position += USTarHeaderUnameLen
		// position = 297
		ustar.Gname = raw[position : position + USTarHeaderGnameLen]
		position += USTarHeaderGnameLen
		// position = 329
		ustar.Devmajor = raw[position : position + USTarHeaderDevmajorLen]
		position += USTarHeaderDevmajorLen
		// position = 337
		ustar.Devminor = raw[position : position + USTarHeaderDevminorLen]
		position += USTarHeaderDevminorLen
		// position = 345
		ustar.Prefix = raw[position : position + USTarHeaderPrefixLen]
		position += USTarHeaderPrefixLen
		// position = 500
		// pads for blocksize length
		position += 12
		// position = 512

		contentLen := ustar.GetContentLen()
		ustar.Content = raw[position : position + contentLen]

		if string(ustar.Magic) == string(USTarMagic) {
			contentBlockCount := int(math.Ceil(float64(contentLen) / float64(TarBlocksize)))
			nextItemStartPosition := position + (contentBlockCount * TarBlocksize)
			position = nextItemStartPosition

			items = append(items, ustar)
			
		} else {
			break
		}
	}

	return items

}

func (ustar USTar) GetContentLen() int {
	contentLen, err := strconv.ParseInt(string(ustar.Size[:11]), 8, 64)

	if err != nil {
		return 0
	}
	return int(contentLen)
}

func (ustar USTar) GetName() string {
	return string(ustar.Name)
}

func (ustar USTar) GetType() string {
	return string(ustar.Typeflag)
}