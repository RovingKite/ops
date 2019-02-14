package lepton

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

// file system manifest
const manifest string = `(
    #64 bit elf to boot from host
    children:(kernel:(contents:(host:%v))
              #user program
              %v:(contents:(host:%v)))
    # filesystem path to elf for kernel to run
    program:/%v
    fault:t
    arguments:[%v sec third]
    environment:(USER:bobby PWD:/)
)`

const Version = "0.2"
const OpsReleaseUrl = "https://storage.googleapis.com/cli/%v/ops"

// boot loader
const BootImg string = ".staging/boot.img"

// kernel
const KernelImg string = ".staging/stage3.img"

// kernel + ELF image
const mergedImg string = ".staging/tempimage"

// final bootable image
const FinalImg string = "image"

const Mkfs string = ".staging/mkfs"

const ReleaseBaseUrl string = "https://storage.googleapis.com/nanos/release/%v"
const NightlyReleaseBaseUrl string = "https://storage.googleapis.com/nanos/release/nightly/"

const PackageBaseURL string = "https://storage.googleapis.com/packagehub/%v"
const PackageManifestURL string = "https://storage.googleapis.com/packagehub/manifest.json"
const PackageManifestFileName string = "manifest.json"

var PackagesCache string

func GetOpsHome() string {
	home, err := HomeDir()
	if err != nil {
		panic(err)
	}

	opshome := path.Join(home, ".ops")
	if _, err := os.Stat(opshome); os.IsNotExist(err) {
		os.MkdirAll(opshome, 0755)
	}

	return opshome
}

func GetPackageCache() string {
	if PackagesCache == "" {
		PackagesCache = path.Join(GetOpsHome(), "packages")
		if _, err := os.Stat(PackagesCache); os.IsNotExist(err) {
			os.MkdirAll(PackagesCache, 0755)
		}
	}
	return PackagesCache
}

func GetPackageManifestFile() string {
	return path.Join(GetPackageCache(), PackageManifestFileName)
}

func NightlyReleaseUrl() string {
	var sb strings.Builder
	sb.WriteString(NightlyReleaseBaseUrl)
	sb.WriteString("nanos_nightly_")
	sb.WriteString(runtime.GOOS)
	sb.WriteString(".tar.gz")
	return sb.String()
}

func NightlyLocalPath() string {
	var sb strings.Builder
	sb.WriteString(path.Join(GetOpsHome(), "nanos_nightly_"))
	sb.WriteString(runtime.GOOS)
	sb.WriteString(".tar.gz")
	return sb.String()
}

func LocalTimeStamp() (string, error) {

	data, err := ioutil.ReadFile(path.Join(GetOpsHome(), "timestamp"))
	// first time download?
	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}
	return string(data), nil
}

func RemoteTimeStamp() (string, error) {
	resp, err := http.Get(path.Join(NightlyReleaseBaseUrl, "timestamp"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LatestReleaseUrl() string {
	return ""
}
