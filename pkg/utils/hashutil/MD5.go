package hashutil

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/hex"
	"github.com/golang/glog"
	"io"

	"kubesphere.io/kubesphere/pkg/utils/readerutils"
)

func GetMD5(reader io.ReadCloser) (string, error) {
	md5reader := readerutils.NewMD5Reader(reader)
	data := make([]byte, bytefmt.KILOBYTE)
	for {
		_, err := md5reader.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			glog.Error(err)
			return "", err
		}
	}
	err := reader.Close()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5reader.MD5()), nil
}
