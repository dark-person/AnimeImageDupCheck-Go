package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"log"
	"os"
	"strings"

	"github.com/corona10/goimagehash"
)

const jpeg_suffix = "jpeg"
const jpg_suffix = "jpg"
const png_suffix = "png"

type ImageFile struct {
	Filename  string // example: 'example.jpg'
	Fullpath  string // example: 'imagefile/example.jpg'
	Filesize  int64  // example: 1024
	Width     int    // example: 1024
	Height    int    // example: 1024
	Directory string // example: 'imagefile'
}

type DuplicateImage struct {
	Filename  string // example: 'example_copy.jpg'
	Fullpath  string // example: 'imagefile/example_copy.jpg'
	HashValue string
}

// Input : directory string
// Output : Slice of ImageFile
func getImageLists(directory string) ([]ImageFile, error) {
	var results []ImageFile

	files, files_err := os.ReadDir(directory)
	if files_err != nil {
		log.Fatal("Cannot open Directory to get Image List, ", files_err)
	}

	for _, file := range files {

		file_lower := strings.ToLower(file.Name())

		if strings.Contains(file_lower, jpeg_suffix) || strings.Contains(file_lower, jpg_suffix) || strings.Contains(file_lower, png_suffix) {
			var tempFile ImageFile

			// Struct Initializing
			tempFile.Filename = file.Name()
			tempFile.Directory = directory
			tempFile.Fullpath = directory + "/" + file.Name()
			tempFile.Filesize = -1

			results = append(results, tempFile)
		}
	}

	fmt.Println("Got File List in the input. Start processing...")
	log.Println("getImageLists Result : ", results)
	return results, nil
}

// Input: Address of ImageFile struct
// Output: pHash value (string) of image
func analyzeImage(imageFile *ImageFile) (string, error) {
	file, file_err := os.Open(imageFile.Fullpath)
	if file_err != nil {
		log.Fatal("Cannot Open Specific Image, ", file_err)
	}

	fileinfo, _ := file.Stat()
	imageFile.Filesize = fileinfo.Size()
	defer file.Close()

	temp_image, _, image_err := image.Decode(bufio.NewReader(file))
	if image_err != nil {
		log.Fatal("Cannot Decode Specific Image, ", image_err)
	}

	imageFile.Width = temp_image.Bounds().Dx()
	imageFile.Height = temp_image.Bounds().Dy()

	hash, _ := goimagehash.PerceptionHash(temp_image)
	log.Println("analyzeImage(", imageFile.Fullpath, ") : ", hash.ToString())

	return hash.ToString(), nil
}

// Input: Slice of ImageFile
// Output: Map of ImageFile with hash vaule as key, Slice of duplicate image fullpath, error
func analyzeImages(filelist []ImageFile) (map[string]ImageFile, []DuplicateImage, error) {

	// Message For info User progress
	fmt.Println("Starting analyze image..")

	hash_map := make(map[string]ImageFile)
	var duplicate_list []DuplicateImage

	for index, file := range filelist {

		hash_value, _ := analyzeImage(&file)
		value, isExist := hash_map[hash_value]

		if !isExist {
			log.Println("New Hash Value (", hash_value, ") Detected. File added :", file.Filename)
			hash_map[hash_value] = file
		} else {
			log.Println("Hash Value (", hash_value, ")Exist. Start Comparing.. ")
			log.Println("Current Best Image (", value.Filename, ") Size: ", value.Filesize)
			log.Println("New Image (", file.Filename, ") Size: ", file.Filesize)

			if file.Filesize > value.Filesize {
				// Create New DuplicateImage Record
				var temp DuplicateImage
				temp.Filename = value.Filename
				temp.Fullpath = value.Fullpath
				temp.HashValue = hash_value

				duplicate_list = append(duplicate_list, temp)
				hash_map[hash_value] = file

				log.Println("Updated the best image (", file.Filename, ").")
			} else {
				// Create New DuplicateImage Record
				var temp DuplicateImage
				temp.Filename = file.Filename
				temp.Fullpath = file.Fullpath
				temp.HashValue = hash_value

				duplicate_list = append(duplicate_list, temp)
				log.Println("Add new duplicate image (", file.Filename, ").")
			}
		}
		// Message For info User progress
		if index%10 == 0 && index != 0 {
			fmt.Println("Analyzed Image :", index, "/", len(filelist))
		}
	}

	// Message For info User progress
	fmt.Println("Analyzed Image :", len(filelist), "/", len(filelist), ". Analyze Completed.")

	return hash_map, duplicate_list, nil
}

func moveFile(filepath, newDirectory, filename string) error {
	os.MkdirAll(newDirectory, 0755)
	return os.Rename(filepath, newDirectory+"/"+filename)
}

func main() {
	// Set Logger
	log_file, err := os.OpenFile("activity.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot Open Log File, ", err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	log.Println("=====  Log Ready. Program Start.")

	// Get Image List by directory
	image_list, input_err := getImageLists("input")
	if input_err != nil {
		log.Fatal("Cannot Open Directory", input_err)
	}

	result_map, duplicate, _ := analyzeImages(image_list)

	fmt.Println("Result is completed. ", len(duplicate), " duplicate Found. ")
	fmt.Println("Now transfer file to different folder.")

	// Start Record and move file
	record, record_err := os.OpenFile("record.txt", os.O_RDWR|os.O_CREATE, 0755)
	if record_err != nil {
		log.Fatal("Cannot open record.txt, ", record_err)
	}

	// Print and Move Result and Duplicates
	record.WriteString("Best Image:\n")
	for _, item := range result_map {
		record.WriteString(fmt.Sprintf(item.Filename) + "\n")

		err := moveFile(item.Fullpath, "Best", item.Filename)
		if err != nil {
			log.Fatal("Cannot Move Image to Best,", err)
		}
	}

	record.WriteString("================================================\n\n")

	record.WriteString("Duplicate Image:\n")
	for _, item := range duplicate {
		record.WriteString(fmt.Sprintf("%s (Duplicate of %s)\n", item.Filename, result_map[item.HashValue].Filename))

		err := moveFile(item.Fullpath, "Duplicate", item.Filename)
		if err != nil {
			log.Fatal("Cannot Move Image to Duplicate,", err)
		}
	}

	log.Println("===================================== Finish ===========================================")
}
