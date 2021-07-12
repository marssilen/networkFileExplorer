package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"os"
	"path/filepath"
	"strings"
	"archive/zip"
	"io"
	"time"
)

func vpngoUsage() {
	fmt.Println("Server is not configured")
	fmt.Println("Usage:VPNGO id Location")
}

func msg(ctx iris.Context,code uint8,data interface{},msg string){
	ctx.JSON(iris.Map{
		"code":  code,
		"data": data,
		"message": msg,
	})
}
func invalidMsg(ctx iris.Context){
	ctx.JSON(iris.Map{
		"code":  1,
		"data": [0]int{},
		"message":"input not valid",
	})
	return
}
func invalidForm(ctx iris.Context){
	ctx.JSON(iris.Map{
		"code":  1,
		"data": [0]int{},
		"message":"form not valid",
	})
	return
}
func printErrors(err error){
	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println(err)
		return
	}

	for _, err := range err.(validator.ValidationErrors) {

		fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
		fmt.Println(err.Field())     // by passing alt name to ReportError like below
		fmt.Println(err.StructNamespace())
		fmt.Println(err.StructField())
		fmt.Println(err.Tag())
		fmt.Println(err.ActualTag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())
		fmt.Println(err.Param())
		fmt.Println()
	}
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string, userId string,child string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		fileWithUserId:=userId+"_"+child+"_"+f.Name
		//fmt.Println("filewiuserid is ",fileWithUserId)
		// Store filename/path for returning and using later on
		//fpath := filepath.Join(dest, f.Name)
		fpath := filepath.Join(dest, fileWithUserId)
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		fmt.Println("file name is:",f.Name)
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
func toShamsi(sTime time.Time,toLetter bool) string{

	//sTime,_:=time.Parse(time.RFC3339Nano,str)
	//fmt.Println(sTime)
	a := sTime.Day() 	 //date('d', time)
	b := int(sTime.Month()) 	//date('m', time)
	c := sTime.Year()	//date('Y', time)
	fmt.Println(a,b,c)
	var x,y,z int
	x , y , z = 0 , 0, 0
	/*
	   a->day
	   b->month
	   c->year
	   x->rooz
	   y->mah
	   z->sal
	*/
	if c<1 || a<1 || b<1 || a>31 || b>12 {
		return "error"
	}
	if c%4==0 {
		if b==2 && a>29 {
			return "error"
		}
		if b==4 || b==6 || b==9 || b==11 {
			if a > 30 {
				return "error"
			}
		}
	}
	if b==10 {
		x=a+9
	}
	if b==4 {
		x=a+12
	}
	if b==2 ||  b== 5 ||  b==6 {
		x=a+11
	}
	if b==1 ||  b==3 ||  b==7 ||  b==8 ||  b==9 ||  b== 11 ||  b==12 {
		x= a+10
	}
	if c%4==1 &&  b<3 {
		x++
		a++
	}
	if c%4==1 ||  c%4==2 ||  c%4==3 {
		if b == 3 {
			a--
		}
	}
	if c%4==1 ||  c%4==2 ||  c%4==3 {
		if b>3 {
			x--
			a--
		}
	}
	if b>1 &&  b<5 &&  a<20 {
		y= b+9
	}
	if b>1 &&  b<5 &&  a>19 {
		x= a-19
		y= b+10
	}
	if b>6 &&  b<11 &&  a<22 {
		y= b+9
	}
	if b>6 &&  b<11 {
		if a > 21 {
			x = a - 21
			y = b + 10
		}
	}
	if b==1 ||  b==5 ||  b==6 ||  b==11 ||  b==12 {
		if a < 21 {
			y = b + 9
		}
	}
	if b==1 ||  b==5 || b==6 ||  b==11 ||  b==12 {
		if a>20 {
			x= a-20
			y= b+10
		}
	}
	if b<4 &&  a<20 ||  b<3 {
		z= c-622
	} else{
		z= c-621
	}
	if y>12 {
		y= y-12
	}
	if toLetter {
		var month string
		switch y {
		case 1:
			month = "فروردین"
			break
		case 2:
			month = "اردیبهشت"
			break
		case 3:
			month = "خرداد"
			break
		case 4:
			month = "تیر"
			break
		case 5:
			month = "مرداد"
			break
		case 6:
			month = "شهریور"
			break
		case 7:
			month = "مهر"
			break
		case 8:
			month = "آبان"
			break
		case 9:
			month = "آذر"
			break
		case 10:
			month = "دی"
			break
		case 11:
			month = "بهمن"
			break
		case 12:
			month = "اسفند"
			break
		}
		return fmt.Sprint(x," "+ month+" ", z,sTime.Hour(),":",sTime.Minute(),":",sTime.Second())

	}
	return fmt.Sprint(x, y, z,sTime.Hour(),":",sTime.Minute(),":",sTime.Second())
}