package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSHA256(t *testing.T) {
	result := SHA256("nozomoto")
	if hex.EncodeToString(result) != "d1f6ff48981b9f83da0c3c34a72a492397734680ea863a4d560ad9bcadd4d62e" {
		t.Error("sha256 SHA256 missing")
	}
}

func TestSHA256Threshold(t *testing.T) {
	result := SHA256Threshold(int64(100), int64(25))
	if fmt.Sprint(result) != "463168356949264781694283940034751631413079938662562256157830336031652518559744" {
		t.Error("SHA256Threshold value missing")
	}
}

func TestSHAtoBigNum(t *testing.T) {
	result, err := SHAtoBigNum("d1f6ff48981b9f83da0c3c34a72a492397734680ea863a4d560ad9bcadd4d62e")
	if err != nil {
		t.Error("error", err)
	} else if fmt.Sprint(result) != "94969791634291581163892593121492363376799753370261336941526847818439966643758" {
		t.Error("SHAtoBigNum is missing")
	}

	result, err = SHAtoBigNum("7baf29e6fd4eae1e2f79a0a0866bc2386cba1ac656cb6ac14db5439a3d5792d3")
	if err != nil {
		t.Error("error", err)
	} else if fmt.Sprint(result) != "55943967811142311813067957484704184073425717843965365978503945418199939912403" {
		t.Error("SHAtoBigNum is missing")
	}

	result, err = SHAtoBigNum("bc3b61ae3aee65c8948b457de1087bb940c0fc2e7e8c15897d5c49e1c4f69a65")
	if err != nil {
		t.Error("error", err)
	} else if fmt.Sprint(result) != "85139733677108548456271096930738765125135848115516693007492771253392393017957" {
		t.Error("SHAtoBigNum is missing")
	}

	result, err = SHAtoBigNum("0")
	if err == nil {
		t.Error("sha s length varidation")
	}
}

func TestIsOKHash(t *testing.T) {
	result := IsOKHash(10000, 2, "nozomoto")
	if result != true {
		t.Error("IsOKHash missing")
	}
	result = IsOKHash(100, 200, "nozomoto")
	if result != false {
		t.Error("IsOKHash missing")
	}

	result = IsOKHash(300, 2, "hatobus")
	if result != true {
		t.Error("IsOKHash missing")
	}

}
