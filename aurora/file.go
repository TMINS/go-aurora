package aurora

import (
	"io"
	"mime/multipart"
	"os"
)

/*
	文件上传api使用的 gin 封装，操作和gin一致
*/

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Ctx) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.ar.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

func (c *Ctx) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.ar.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// SaveUploadedFile uploads the form file to specific dst.
func (c *Ctx) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (c *Ctx) perse() interface{} {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.ar.MaxMultipartMemory); err != nil {
			return nil
		}
	}

	return nil
}
