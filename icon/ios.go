package icon

import (
	"fmt"
	"path"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/ziputil"
	"github.com/bitrise-tools/go-xcode/ipa"
	"github.com/bitrise-tools/go-xcode/plistutil"
)

func unzipIPA(pth string) (string, error) {
	dest, err := pathutil.NormalizedOSTempDirPath("ipa")
	if err != nil {
		return "", err
	}

	return dest, ziputil.UnZip(pth, dest)
}

func findIcons(plistData plistutil.PlistData) ([]string, error) {
	cfBundleIcons, ok := plistData.GetMapStringInterface("CFBundleIcons")
	if !ok {
		return nil, fmt.Errorf("CFBundleIcons not found")
	}

	primaryIcons, ok := cfBundleIcons.GetMapStringInterface("CFBundlePrimaryIcon")
	if !ok {
		return nil, fmt.Errorf("CFBundlePrimaryIcon not found")
	}

	iconNames, ok := primaryIcons.GetStringArray("CFBundleIconFiles")
	if !ok {
		return nil, fmt.Errorf("CFBundleIconFiles not found")
	}

	return iconNames, nil
}

// FetchIcon ...
func FetchIcon(ipaPth, projectName string) (string, error) {
	plistPth, err := ipa.UnwrapEmbeddedInfoPlist(ipaPth)
	if err != nil {
		panic(err)
	}

	infoPlistData, err := plistutil.NewPlistDataFromFile(plistPth)
	if err != nil {
		panic(err)
	}

	iconNames, err := findIcons(infoPlistData)
	if err != nil {
		panic(err)
	}

	if len(iconNames) == 0 {
		return "", nil
	}

	dest, err := unzipIPA(ipaPth)
	if err != nil {
		return "", err
	}

	//
	// @2x
	iconPth := path.Join(dest, "Payload", projectName+".app", iconNames[0]+"@2x.png")
	exists, err := pathutil.IsPathExists(iconPth)
	if err != nil {
		return "", err
	}

	//
	// if the @2x is not exists try the original one
	if !exists {
		iconPth = path.Join(dest, "Payload", projectName+".app", iconNames[0]+".png")
		exists, err = pathutil.IsPathExists(iconPth)
		if err != nil {
			return "", err
		}
	}
	return iconPth, nil

}
