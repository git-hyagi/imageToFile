### Pre-reqs
~~~
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Alexander_Pozdnyakov/CentOS_8/
dnf install tesseract tesseract-devel tesseract-langpack-por tesseract-langpack-eng
dnf install gcc-c++
~~~

### "Improve" image quality
https://pt.wikipedia.org/wiki/OpenCV


### To build and run
~~~
cd server
docker build -t imagetotext:latest .
docker run -d -it -p 80:80 imagetotext:latest
~~~


### References
[gosseract](https://github.com/otiai10/gosseract/)  
[tesseract](https://github.com/tesseract-ocr/tesseract/wiki)
