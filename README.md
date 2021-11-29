# AnimeImageDupCheck-Go

A Tool Written By GoLang, that help seperate duplicate image.

This Project is based on [AnimeImageDuplicateCheck](https://github.com/dark-person/AnimeImageDuplicateCheck), which is written by Python.

# Require

You need Go installed in you computer.

# Build Executable

```
go build
```

And run the executable. You should ensure that is a directory called `input/`.

Also, ensure some file is inside the `input/` directory.

## Window Special

For window, if you want the terminal not closed after completed, you may create a .bat, which content is
```
AnimeImageDepCheck.exe
pause
```
Then run the bat instead.

# Usage

Put all image to input folder. Then run this thing.
You will get the following file/directory:

```
Best/
Duplicate/
activity.log
record.txt
```

`Best/` is the directory that stores the highest quality of image.

`Duplicate/` is the directory that stores the duplicate image.
You may find out which image is duplicated in `record.txt`

`activity.log` will record all parameters and running processes. Mostly used in debug.

# For Any Bugs Found

Please inform the bug with these file:
```
activiy.log
record.txt
(Your Testing File)
```
