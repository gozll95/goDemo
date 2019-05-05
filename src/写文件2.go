case tar.TypeReg: // = regular file
			fmt.Println("Regular file:", name)
			data := make([]byte, header.Size)
			_, err := tarReader.Read(data)
			if err != nil {
				panic("Error reading file!!! PANIC!!!!!!")
			}

			ioutil.WriteFile(name, data, 0755)