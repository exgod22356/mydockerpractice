package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
NewPipe function
create a pipe to store the command
*/
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootURL string, mntURL string, volume string) error {
	if err := CreateReadOnlyLayer(rootURL); err != nil {
		return fmt.Errorf("create readonly layer error :%v", err)
	}
	if err := CreateWriteLayer(rootURL); err != nil {
		return fmt.Errorf("create writeonly layer error :%v", err)
	}
	if err := CreateMountPoint(rootURL, mntURL); err != nil {
		return fmt.Errorf("create mount point error :%v", err)
	}
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			if err := MountVolume(rootURL, mntURL, volumeURLs); err != nil {
				return fmt.Errorf("MountVolumeError : %v", err)
			}
			fmt.Printf("%v\n", volumeURLs)
		} else {
			return fmt.Errorf("volume error: %v is not correct", volumeURLs)
		}
	}
	return nil
}

func CreateReadOnlyLayer(rootURL string) error {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		return fmt.Errorf("path exist error:%v", err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			return fmt.Errorf("mkdir root error: %v", err)
		}
		if _, err := exec.Command("tar", "-vxf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			return fmt.Errorf("tar error: %v", err)
		}
	}
	return nil
}

func CreateWriteLayer(rootURL string) error {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		return fmt.Errorf("mkdir error :%v", err)
	}
	workURL := rootURL + "workLayer/"
	if err := os.Mkdir(workURL, 0777); err != nil {
		return fmt.Errorf("mkdir error :%v", err)
	}
	return nil
}

func CreateMountPoint(rootURL string, mntURL string) error {
	fmt.Println("mounting ")
	if err := os.Mkdir(mntURL, 0777); err != nil {
		return fmt.Errorf("mkdir error :%v", err)
	}
	dirs := "lowerdir=" + rootURL + "busybox,upperdir=" + rootURL + "writeLayer,workdir=" + rootURL + "workLayer"
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "overlay", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create mount point error : %v", err)
	}
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func volumeUrlExtract(volume string) []string {
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}

func MountVolume(rootURL string, mntURL string, volumeURLs []string) error {
	parentURL := volumeURLs[0]
	if err := os.Mkdir(parentURL, 0777); err != nil {
		if exist, _ := PathExists(parentURL); exist == false {
			return fmt.Errorf("error in mkdir parentURL: %v", err)
		}
	}
	containerURL := volumeURLs[1]
	containerVolumeURL := mntURL + containerURL
	fmt.Println(containerVolumeURL)
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		return fmt.Errorf("error in mkdir containerVolumeURL: %v", err)
	}
	cmd := exec.Command("mount", "--bind", parentURL, containerVolumeURL)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Mount volume error :%v", err)
	}
	return nil
}

func DeleteWorkSpace(rootURL string, mntURL string, volume string) error {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			if err := DeleteVolumeMountPoint(mntURL, volumeURLs[1]); err != nil {
				return fmt.Errorf("delete mount volumn error: %v", err)
			}
		}
	}
	if err := DeleteMountPoint(rootURL, mntURL); err != nil {
		return fmt.Errorf("DeleteMountPoint error: %v", err)
	}
	if err := DeleteWriteLayer(rootURL); err != nil {
		return fmt.Errorf("DeleteWriteLayer error: %v", err)
	}
	return nil
}

func DeleteMountPoint(rootURL string, mntURL string) error {
	fmt.Println(mntURL)
	cmd := exec.Command("umount", "-v", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Run umount error :%v", err)
	}
	cmd = exec.Command("umount", "-v", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Run umount error :%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		return fmt.Errorf("remove %v error:%v", mntURL, err)
	}
	return nil
}

func DeleteWriteLayer(rootURL string) error {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		return fmt.Errorf("remove %v error:%v", writeURL, err)
	}
	workURL := rootURL + "workLayer/"
	if err := os.RemoveAll(workURL); err != nil {
		return fmt.Errorf("remove %v error:%v", writeURL, err)
	}
	busyboxURL := rootURL + "busybox/"
	if err := os.RemoveAll(busyboxURL); err != nil {
		return fmt.Errorf("remove %v error:%v", writeURL, err)
	}
	return nil
}

func DeleteVolumeMountPoint(mntURL string, volumnURL string) error {
	cmd := exec.Command("umount", "-v", mntURL+volumnURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("delete mount volume error : %v", err)
	}
	return nil
}
